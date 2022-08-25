package model

import (
	"mime/multipart"
	"tiktink/internal/code"
)

// Video 与数据库交互的video模型
type Video struct {
	ID       int64  `gorm:"column:video_id;not null"`
	AuthorID int64  `gorm:"column:author_id;not null"`
	PlayURL  string `gorm:"column:play_url;not null"`
	VideoKey string `gorm:"column:video_key;not null"`
	CoverURL string `gorm:"column:cover_url;not null"`
	ImageKey string `gorm:"column:image_key;not null"`
	Title    string `gorm:"column:title;not null"`
}

func (v Video) TableName() string {
	return "videos"
}

type PublishVideoReq struct {
	Data  *multipart.FileHeader `form:"data" binding:"required"`
	Title string                `form:"title" binding:"required"`
}

type PublishListReq struct {
	UserID int64 `form:"user_id" binding:"required"`
}

type PublishListResp struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	VideoList  []*VideoMSG  `json:"video_list"`
}

type VideoMSG struct {
	ID            int64 `json:"id" gorm:"column:video_id"`
	UserMSG       `json:"author"`
	PlayURL       string `json:"play_url" gorm:"play_url"`
	CoverURL      string `json:"cover_url" gorm:"cover_url"`
	FavoriteCount int64  `json:"favorite_count" gorm:"favorite_count"`
	CommentCount  int64  `json:"comment_count" gorm:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title" gorm:"title"`
}
