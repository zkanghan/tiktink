package model

import "tiktink/internal/code"

type Favorite struct {
	VideoID string `gorm:"column:video_id;not null"`
	UserID  string `gorm:"column:user_id;not null"`
}

func (Favorite) TableName() string {
	return "favorites"
}

type FavoriteActionReq struct {
	VideoID    string `form:"video_id" binding:"required"`
	ActionType int8   `form:"action_type" binding:"required,oneof=1 2"`
}

type FavoriteListReq struct {
	UserID     string `form:"user_id" binding:"required"`
	PageNumber int    `form:"pn" binding:"required"`
}

type FavoriteListResp struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	VideoList  []*VideoMSG  `json:"video_list"`
}

type FavoriteRedis struct {
	UserID  string `json:"UserID"`
	VideoID string `json:"VideoID"`
	Status  string `json:"Status"`
}
