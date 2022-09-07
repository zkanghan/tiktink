package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/tracer"
)

type commentDealer struct {
	Context *tracer.TraceCtx
}

type commentFunc interface {
	ReleaseComment(req *model.CommentActionReq, userID int64) (*model.CommentMSG, error)
	DeleteComment(videoID int64, commentID int64, userID int64) (bool, error)
	GetCommentList(videoID int64, userID int64) ([]*model.CommentMSG, error)
}

var _ commentFunc = &commentDealer{}

func NewCommentDealer(ctx *tracer.TraceCtx) *commentDealer {
	return &commentDealer{
		Context: ctx,
	}
}

// ReleaseComment 指针类型的方法接收者修改才是有效的
func (c *commentDealer) ReleaseComment(req *model.CommentActionReq, userID int64) (*model.CommentMSG, error) {
	c.Context.TraceCaller() // save context message

	comment, err := mysql.NewCommentDealer(c.Context).CreateComment(req.VideoID, userID, req.CommentText)
	if err != nil {
		return nil, err
	}
	userMsg, err := mysql.NewUserDealer(c.Context).QueryUserByID(userID)
	if err != nil {
		return nil, err
	}
	//  format the time
	createDate := comment.CreateAt.Format("2006-01-02")
	commentMsg := &model.CommentMSG{
		CommentID:  comment.ID,
		UserMSG:    *userMsg,
		Content:    comment.Content,
		CreateDate: createDate,
	}
	return commentMsg, nil
}

// DeleteComment 评论删除成功返回true
func (c *commentDealer) DeleteComment(videoID int64, commentID int64, userID int64) (bool, error) {
	c.Context.TraceCaller() //save context message

	affectRows, err := mysql.NewCommentDealer(c.Context).DeleteComment(videoID, commentID, userID)
	if err != nil {
		return false, err
	}
	return affectRows > 0, nil
}

func (c *commentDealer) GetCommentList(videoID int64, userID int64) ([]*model.CommentMSG, error) {
	c.Context.TraceCaller()

	commentMsgs, err := mysql.NewCommentDealer(c.Context).QueryCommentList(videoID)
	if err != nil {
		return nil, err
	}
	//  todo: 把循环去掉改为一次查询
	for _, comment := range commentMsgs {
		followed, err := NewRelationDealer(c.Context).GetIsFollowed(userID, comment.UserMSG.ID)
		if err != nil {
			return nil, err
		}
		comment.UserMSG.IsFollow = followed
		comment.CreateDate = comment.CreateDate[0:10]
	}
	return commentMsgs, nil
}
