package logic

import (
	"bytes"
	"fmt"
	"os"
	"tiktink/internal/cos"
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/snowid"
	"tiktink/pkg/tools"
	"tiktink/pkg/tracer"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

const (
	filePreVideo    string = "video/"
	fileSuffixVideo string = ".mp4"
	filePreImage    string = "image/"
	fileSuffixImage string = ".jpg"
)

type videoDealer struct {
	Context *tracer.TraceCtx
}

type videoFunc interface {
	PublishVideo(video *model.PublishVideoReq, userID string) error
	GetIsVideoExist(videoID string) (bool, error)
	GetVideoList(userID string, req model.PublishListReq) ([]*model.VideoMSG, error)
}

func NewVideoDealer(ctx *tracer.TraceCtx) videoFunc {
	return &videoDealer{
		Context: ctx,
	}
}

//  文件上传之后生成缩略图并上传           返回缩略图在云上的访问路径
func generateCover(videoPath string, coverID string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	coverDst := filePreImage + coverID + fileSuffixImage
	coverURL, err := cos.PublishFileToServer(buf, coverDst)
	if err != nil {
		return "", err
	}
	return coverURL, err
}

func (v *videoDealer) PublishVideo(video *model.PublishVideoReq, userID string) error {
	//  生成视频唯一id，并拼接对象存储唯一key
	videoID := snowid.GenID()
	coverID := snowid.GenID()
	videoServeDst := filePreVideo + videoID + fileSuffixVideo
	// 打开文件
	videoFile, err := video.Data.Open()
	if err != nil {
		zap.L().Error("打开文件失败：", zap.Error(err))
		return err
	}
	defer videoFile.Close()
	//  把文件上传到腾讯云
	videoURL, err := cos.PublishFileToServer(videoFile, videoServeDst)
	if err != nil {
		zap.L().Error("视频文件上传服务器失败：", zap.Error(err))
		return err
	}
	//  再最后试一下，生成缩略图并上传
	coverURL, err := generateCover(videoURL, coverID)
	if err != nil {
		zap.L().Error("生成视频缩略图失败：", zap.Error(err))
		return err
	}
	//  路径存储到数据库
	videoModel := &model.Video{
		AuthorID: userID,
		Title:    video.Title,
		PlayURL:  videoURL,
		VideoID:  videoID,
		CoverURL: coverURL,
		ImageID:  coverID,
	}
	err = mysql.NewVideoDealer().PublishVideo(videoModel)
	if err != nil {
		zap.L().Error("视频存储数据库出错：", zap.Error(err))
		return err
	}
	return nil
}

func (v *videoDealer) GetIsVideoExist(videoID string) (bool, error) {
	return mysql.NewVideoDealer().QueryVideoExist(videoID)
}

func (v *videoDealer) GetVideoList(currentUserID string, req model.PublishListReq) ([]*model.VideoMSG, error) {
	videoMsgs, err := mysql.NewVideoDealer().QueryVideoByAuthorID(req.UserID, req.PageNumber)
	if err != nil {
		return []*model.VideoMSG{}, err
	}

	// 组装视频id切片
	var videoIDs []string
	for _, video := range videoMsgs {
		videoIDs = append(videoIDs, video.VideoID)
	}
	//  查询当前用户是否点赞了哪些视频列表中的视频
	likedVideoIDs, err := mysql.NewFavoriteDealer().QueryListIsLiked(currentUserID, videoIDs)
	if err != nil {
		return []*model.VideoMSG{}, err
	}
	// 转map便于查询
	likedVideoIDsMap := tools.SliceIntToSet(likedVideoIDs)
	// 判断当前用户是否关注视频作者
	followed, err := NewRelationDealer(v.Context).GetIsFollowed(currentUserID, req.UserID)
	if err != nil {
		return []*model.VideoMSG{}, err
	}
	for _, video := range videoMsgs {
		_, liked := likedVideoIDsMap[video.VideoID]
		video.IsFavorite = liked
		video.IsFollow = followed
	}
	return videoMsgs, nil
}
