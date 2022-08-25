package logic

import (
	"bytes"
	"fmt"
	"os"
	"tiktink/internal/cos"
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/snowid"

	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

const (
	filePreVideo    string = "video/"
	fileSuffixVideo string = ".mp4"
	filePreImage    string = "image/"
	fileSuffixImage string = ".jpg"
)

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

func PublishVideo(c *gin.Context, video *model.PublishVideoReq, userID int64) error {
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
		VideoKey: videoID,
		CoverURL: coverURL,
		ImageKey: coverID,
	}
	err = mysql.DealVideo().PublishVideo(videoModel)
	if err != nil {
		zap.L().Error("视频存储数据库出错：", zap.Error(err))
		return err
	}
	return nil
}

func GetIsVideoExist(videoID int64) (bool, error) {
	return mysql.DealVideo().QueryVideoExist(videoID)
}

func GetVideoList(userID int64, authorID int64) ([]*model.VideoMSG, error) {
	videoMsgs, err := mysql.DealVideo().QueryVideoByAuthorID(authorID)
	if err != nil {
		return nil, err
	}
	for _, video := range videoMsgs {
		followed, err := GetIsFollowed(userID, authorID)
		if err != nil {
			return nil, err
		}
		liked, err := GetIsLiked(userID, video.ID)
		if err != nil {
			return nil, err
		}
		video.IsFollow = followed
		video.IsFavorite = liked
	}
	return videoMsgs, nil
}
