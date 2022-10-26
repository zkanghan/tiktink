package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/tracer"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

const followPageRows int64 = 20

type followFunc interface {
	QueryIsFollow(UserID string, ToUserID string) (bool, error)
	DoFollow(UserID string, ToUserID string) error
	DoCancelFollow(UserID string, ToUserID string) error
	QueryFollowList(req *model.FollowListReq) ([]*model.UserMSG, error)
	QueryFansList(req *model.FollowListReq) ([]*model.UserMSG, error)
	QueryListIsFollow(userID string, toUsers []string) ([]string, error)
}

type followDealer struct{}

func NewRelationDealer() followFunc {
	return &followDealer{}
}

func (f *followDealer) QueryFansList(req *model.FollowListReq) ([]*model.UserMSG, error) {
	userMSGs := new([]*model.UserMSG)
	offset := (req.PageCount - 1) * followPageRows //展示记录的起点
	count := followPageRows                        //展示记录的终点
	err := db.Raw("SELECT `users`.`user_id`,`user_name`,`follow_count`,`follower_count` from `users` WHERE `users`.`user_id` IN (SELECT `user_id` FROM `follow` WHERE `follow`.`to_user_id` = ?) LIMIT ?,?",
		req.UserID, offset, count).Scan(userMSGs).Error
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(req))
	}
	return *userMSGs, err
}

func (f *followDealer) QueryFollowList(req *model.FollowListReq) ([]*model.UserMSG, error) {
	userMSGs := new([]*model.UserMSG)
	offset := (req.PageCount - 1) * followPageRows //展示记录的起点
	count := followPageRows                        //展示记录的终点
	err := db.Raw("SELECT `user_id`,`user_name`,`follow_count`,`follower_count`FROM `users` "+
		"WHERE `users`.`user_id` "+
		"IN (SELECT `to_user_id` FROM `follow` WHERE `follow`.`user_id` = ?) "+
		"LIMIT ?,?", req.UserID, offset, count).Scan(userMSGs).Error
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(req))
	}
	return *userMSGs, nil
}

func (f *followDealer) DoFollow(UserID string, ToUserID string) error {
	followModel := model.Follow{
		UserID:   UserID,
		ToUserID: ToUserID,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(followModel).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(UserID, ToUserID))
		}
		if err := tx.Exec("update users set follow_count = follow_count+1 where user_id = ?", UserID).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(UserID, ToUserID))
		}
		if err := tx.Exec("update users set follower_count = follower_count+1 where user_id = ?", ToUserID).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(UserID, ToUserID))
		}
		return nil
	})
}

func (f *followDealer) DoCancelFollow(UserID string, ToUserID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? and to_user_id = ?", UserID, ToUserID).Delete(&model.Follow{}).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(UserID, ToUserID))
		}
		//粉丝列表减少一位
		if err := tx.Exec("update users set follower_count = follower_count-1 where user_id = ?", ToUserID).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(UserID, ToUserID))
		}
		//关注列表减少一位
		if err := tx.Exec("update users set follow_count = follow_count-1 where user_id = ?", UserID).Error; err != nil {
			return errors.Wrap(err, tracer.FormatParam(UserID, ToUserID))
		}
		return nil
	})
}

func (f *followDealer) QueryIsFollow(UserID string, ToUserID string) (bool, error) {
	res := new(int8)
	err := db.Raw("select 1 from follow where user_id = ? and to_user_id = ? limit 1", UserID, ToUserID).Scan(res).Error
	return *res == 1, errors.Wrap(err, tracer.FormatParam(UserID, ToUserID))
}

// QueryListIsFollow 返回toUsers中user关注的部分
func (f *followDealer) QueryListIsFollow(userID string, toUsers []string) ([]string, error) {
	var followedUsers []string
	err := db.Raw("select to_user_id from follow where user_id = ? and to_user_id in ?", userID, toUsers).Scan(&followedUsers).Error
	if err != nil {
		return []string{}, errors.Wrap(err, tracer.FormatParam(userID, toUsers))
	}
	return followedUsers, nil
}
