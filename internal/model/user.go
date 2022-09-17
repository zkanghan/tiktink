package model

import (
	"tiktink/internal/code"
)

// User 与数据库交互的User模型
type User struct {
	UserID   int64  `gorm:"column:user_id"`
	UserName string `gorm:"column:user_name;not null"`
	Password string `gorm:"column:password;not null"`
}

func (u User) TableName() string {
	return "users"
}

type UserMSG struct {
	UserID        int64  `json:"id" gorm:"column:user_id"`
	Name          string `json:"name" gorm:"column:user_name"`
	FollowCount   int64  `json:"follow_count" gorm:"column:follow_count"`
	FollowerCount int64  `json:"follower_count" gorm:"column:follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

//  下面为响应和请求的参数模型

type UserRequest struct {
	UserName string `binding:"required,min=1,max=20" form:"username"`
	Password string `binding:"required,min=1,max=20" form:"password"`
}

type LoginResponse struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	UserID     int64        `json:"user_id"`
	Token      string       `json:"token"`
}

type RegisterResponse struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	UserID     int64        `json:"user_id"`
}

type UserInfoRequest struct {
	UserID int64 `binding:"required" form:"user_id"`
}

// UserInfoResponse 匿名字段实现继承的效果
type UserInfoResponse struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	UserMSG    `json:"user"`
}
