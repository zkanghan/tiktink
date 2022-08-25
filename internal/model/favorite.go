package model

import "tiktink/internal/code"

type Favorite struct {
	VideoID int64 `gorm:"column:video_id;not null"`
	UserID  int64 `gorm:"column:user_id;not null"`
}

func (Favorite) TableName() string {
	return "favorite"
}

type FavoriteActionReq struct {
	VideoID    int64 `form:"video_id" binding:"required"`
	ActionType int8  `form:"action_type" binding:"required,oneof=1 2"`
}

type FavoriteListReq struct {
	UserID int64 `form:"user_id" binding:"required"`
}

type FavoriteListResp struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	VideoList  []*VideoMSG  `json:"video_list"`
}
