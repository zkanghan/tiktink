package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
)

func ReleaseComment(req *model.CommentActionReq, userID int64) (*model.CommentMSG, error) {
	comment, err := mysql.DealComment().CreateComment(req.VideoID, userID, req.CommentText)
	if err != nil {
		return nil, err
	}
	userMsg, err := mysql.DealUser().QueryUserByID(userID)
	if err != nil {
		return nil, err
	}
	//  格式化时间
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
func DeleteComment(videoID int64, commentID int64, userID int64) (bool, error) {
	affectRows, err := mysql.DealComment().DeleteComment(videoID, commentID, userID)
	if err != nil {
		return false, err
	}
	return affectRows > 0, nil
}

func GetCommentList(videoID int64, userID int64) ([]*model.CommentMSG, error) {
	commentMsgs, err := mysql.DealComment().QueryCommentList(videoID)
	if err != nil {
		return nil, err
	}
	for _, comment := range commentMsgs {
		followed, err := GetIsFollowed(userID, comment.UserMSG.ID)
		if err != nil {
			return nil, err
		}
		comment.UserMSG.IsFollow = followed
		comment.CreateDate = comment.CreateDate[0:10]
	}
	return commentMsgs, nil
}
