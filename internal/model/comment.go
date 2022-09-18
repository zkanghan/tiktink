package model

import (
	"tiktink/internal/code"
	"time"
)

// Comment 与数据库交互的结构体
type Comment struct {
	CommentID string    `gorm:"column:comment_id"`
	AuthorID  string    `gorm:"column:author_id;not null"`
	VideoID   string    `gorm:"column:video_id;not null"`
	Content   string    `gorm:"column:content;not null"`
	CreateAt  time.Time `gorm:"column:create_date;autoCreateTime;not null"`
}

func (Comment) TableName() string {
	return "comments"
}

// 以下为请求与响应的结构体

type CommentActionReq struct {
	VideoID     string `form:"video_id" binding:"required"`
	ActionType  int8   `form:"action_type" binding:"required,oneof=1 2"`
	CommentText string `form:"comment_text"`
	CommentID   string `form:"comment_id"`
}

type CommentActionResp struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	Comment    *CommentMSG  `json:"comment"`
}

type CommentListReq struct {
	VideoID string `form:"video_id" binding:"required"`
}

type CommentListResp struct {
	StatusCode  code.ResCode  `json:"status_code"`
	StatusMsg   string        `json:"status_msg"`
	CommentList []*CommentMSG `json:"comment_list"`
}

type CommentMSG struct {
	CommentID  string `json:"id"`
	UserMSG    `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
