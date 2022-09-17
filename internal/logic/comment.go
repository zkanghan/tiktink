package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/tools"
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
		CommentID:  comment.CommentID,
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

// TODO:限制一次查询的评论数
func (c *commentDealer) GetCommentList(videoID int64, userID int64) ([]*model.CommentMSG, error) {
	c.Context.TraceCaller()

	commentMsgs, err := mysql.NewCommentDealer(c.Context).QueryCommentList(videoID)
	if err != nil {
		return nil, err
	}

	// 获取需要的用户id
	var toUserIDs []int64
	for _, comment := range commentMsgs {
		toUserIDs = append(toUserIDs, comment.UserID)
	}
	//  获取评论者在userID关注了哪些
	followedUsers, err := mysql.NewRelationDealer(c.Context).QueryListIsFollow(userID, toUserIDs)
	if err != nil {
		return []*model.CommentMSG{}, err
	}
	set := tools.SliceIntToSet(followedUsers) //转化成map方便查询
	for _, comment := range commentMsgs {
		_, followed := set[comment.UserID] //存在该用户ID表示已关注
		comment.UserMSG.IsFollow = followed
		comment.CreateDate = comment.CreateDate[0:10]
	}
	return commentMsgs, nil
}
