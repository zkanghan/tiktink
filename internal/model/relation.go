package model

import "tiktink/internal/code"

type Follow struct {
	UserID   string `gorm:"column:user_id;not null"`
	ToUserID string `gorm:"column:to_user_id;not null"`
}

func (Follow) TableName() string {
	return "follow"
}

//  以下为接口请求和响应的结构体

type FollowActionReq struct {
	ToUserID   string `form:"to_user_id" binding:"required"`
	ActionType int8   `form:"action_type" binding:"required,oneof=1 2"`
}

type FollowListReq struct {
	UserID    string `form:"user_id" binding:"required"`
	PageCount int64  `form:"pn" binding:"required,min=0"`
}

type FollowListResp struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	UserList   []*UserMSG   `json:"user_list"`
}
