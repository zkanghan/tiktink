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
	ReleaseComment(req *model.CommentActionReq, userID string) (*model.CommentMSG, error)
	DeleteComment(videoID string, commentID string, userID string) (bool, error)
	GetCommentList(req model.CommentListReq, currentUserID string) ([]*model.CommentMSG, error)
}

var _ commentFunc = &commentDealer{}

func NewCommentDealer(ctx *tracer.TraceCtx) *commentDealer {
	return &commentDealer{
		Context: ctx,
	}
}

// ReleaseComment 指针类型的方法接收者修改才是有效的
func (c *commentDealer) ReleaseComment(req *model.CommentActionReq, userID string) (*model.CommentMSG, error) {

	comment, err := mysql.NewCommentDealer().CreateComment(req.VideoID, userID, req.CommentText)
	if err != nil {
		return nil, err
	}
	userMsg, err := mysql.NewUserDealer().QueryUserByID(userID)
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
func (c *commentDealer) DeleteComment(videoID string, commentID string, userID string) (bool, error) {

	affectRows, err := mysql.NewCommentDealer().DeleteComment(videoID, commentID, userID)
	if err != nil {
		return false, err
	}
	return affectRows > 0, nil
}

func (c *commentDealer) GetCommentList(req model.CommentListReq, currentUserID string) ([]*model.CommentMSG, error) {

	commentMsgs, err := mysql.NewCommentDealer().QueryCommentList(req.VideoID, req.PageCount)
	if err != nil {
		return nil, err
	}

	// 获取需要的用户id
	var toUserIDs []string
	for _, comment := range commentMsgs {
		toUserIDs = append(toUserIDs, comment.UserID)
	}
	//  获取评论者在userID关注了哪些
	followedUsers, err := mysql.NewRelationDealer().QueryListIsFollow(currentUserID, toUserIDs)
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
