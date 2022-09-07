package logic

import (
	"bytes"
	"fmt"
	"os"
	"tiktink/internal/cos"
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/snowid"
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
	PublishVideo(video *model.PublishVideoReq, userID int64) error
	GetIsVideoExist(videoID int64) (bool, error)
	GetVideoList(userID int64, authorID int64) ([]*model.VideoMSG, error)
}

var _ videoFunc = &videoDealer{}

func NewVideoDealer(ctx *tracer.TraceCtx) *videoDealer {
	return &videoDealer{
		Context: ctx,
	}
}

//  文件上传之后生成缩略图并上传           返回缩略图在云上的访问路径
func generateCover(videoPath string, coverID string, ctx *tracer.TraceCtx) (string, error) {
	ctx.TraceCaller()
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
	coverURL, err := cos.PublishFileToServer(buf, coverDst, ctx)
	if err != nil {
		return "", err
	}
	return coverURL, err
}

func (v *videoDealer) PublishVideo(video *model.PublishVideoReq, userID int64) error {
	v.Context.TraceCaller()
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
	videoURL, err := cos.PublishFileToServer(videoFile, videoServeDst, v.Context)
	if err != nil {
		zap.L().Error("视频文件上传服务器失败：", zap.Error(err))
		return err
	}
	//  再最后试一下，生成缩略图并上传
	coverURL, err := generateCover(videoURL, coverID, v.Context)
	if err != nil {
		zap.L().Error("生成视频缩略图失败：", zap.Error(err))
		return err
	}
	//  路径存储到数据库
	videoModel := &model.Video{
		AuthorID: userID,
		Title:    video.Title,
		PlayURL:  videoURL,
		VideoKey: videoID,
		CoverURL: coverURL,
		ImageKey: coverID,
	}
	err = mysql.NewVideoDealer(v.Context).PublishVideo(videoModel)
	if err != nil {
		zap.L().Error("视频存储数据库出错：", zap.Error(err))
		return err
	}
	return nil
}

func (v *videoDealer) GetIsVideoExist(videoID int64) (bool, error) {
	v.Context.TraceCaller()
	return mysql.NewVideoDealer(v.Context).QueryVideoExist(videoID)
}

func (v *videoDealer) GetVideoList(userID int64, authorID int64) ([]*model.VideoMSG, error) {
	v.Context.TraceCaller()
	videoMsgs, err := mysql.NewVideoDealer(v.Context).QueryVideoByAuthorID(authorID)
	if err != nil {
		return nil, err
	}
	//  todo: 把循环去掉改为一次查询

	for _, video := range videoMsgs {
		followed, err := NewRelationDealer(v.Context).GetIsFollowed(userID, authorID)
		if err != nil {
			return nil, err
		}
		liked, err := NewFavoriteDealer(v.Context).GetIsLiked(userID, video.ID)
		if err != nil {
			return nil, err
		}
		video.IsFollow = followed
		video.IsFavorite = liked
	}
	return videoMsgs, nil
}
