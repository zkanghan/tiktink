package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/tracer"

	"gorm.io/gorm"
)

const pageRows int64 = 20

func NewRelationDealer(ctx *tracer.TraceCtx) followFunc {
	return &followDealer{
		Context: ctx,
	}
}

type followFunc interface {
	QueryIsFollow(UserID int64, ToUserID int64) (bool, error)
	DoFollow(UserID int64, ToUserID int64) error
	DoCancelFollow(UserID int64, ToUserID int64) error
	QueryFollowList(req *model.FollowListReq) ([]*model.UserMSG, error)
	QueryFansList(req *model.FollowListReq) ([]*model.UserMSG, error)
}

type followDealer struct {
	Context *tracer.TraceCtx
}

func (f *followDealer) QueryFansList(req *model.FollowListReq) ([]*model.UserMSG, error) {
	f.Context.TraceCaller()
	userMSGs := new([]*model.UserMSG)
	begin := (req.PageCount - 1) * pageRows //展示记录的起点
	end := pageRows                         //展示记录的终点
	err := db.Raw("SELECT `users`.user_id`,`user_name`,`follow_count`,`follower_count`FROM `users`"+
		"WHERE `users`.`id` "+
		"IN (SELECT `user_id` FROM `follow` WHERE `follow`.`to_user_id` = ?) "+
		"LIMIT ?,?", req.UserID, begin, end).Scan(userMSGs).Error
	if err != nil {
		return nil, err
	}
	return *userMSGs, err
}

func (f *followDealer) QueryFollowList(req *model.FollowListReq) ([]*model.UserMSG, error) {
	f.Context.TraceCaller()
	userMSGs := new([]*model.UserMSG)
	begin := (req.PageCount - 1) * pageRows //展示记录的起点
	end := pageRows                         //展示记录的终点
	err := db.Raw("SELECT `user_id`,`user_name`,`follow_count`,`follower_count`FROM `users` "+
		"WHERE `users`.`id` "+
		"IN (SELECT `to_user_id` FROM `follow` WHERE `follow`.`user_id` = ?) "+
		"LIMIT ?,?", req.UserID, begin, end).Scan(userMSGs).Error
	if err != nil {
		return nil, err
	}
	return *userMSGs, nil
}

func (f *followDealer) DoFollow(UserID int64, ToUserID int64) error {
	f.Context.TraceCaller()
	followModel := model.Follow{
		UserID:   UserID,
		ToUserID: ToUserID,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(followModel).Error; err != nil {
			return err
		}
		if err := tx.Exec("update users set follow_count = follow_count+1 where id = ?", UserID).Error; err != nil {
			return err
		}
		if err := tx.Exec("update users set follower_count = follower_count+1 where id = ?", ToUserID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (f *followDealer) DoCancelFollow(UserID int64, ToUserID int64) error {
	f.Context.TraceCaller()
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? and to_user_id = ?", UserID, ToUserID).Delete(&model.Follow{}).Error; err != nil {
			return err
		}
		//粉丝列表减少一位
		if err := tx.Exec("update users set follower_count = follower_count-1 where id = ?", ToUserID).Error; err != nil {
			return err
		}
		//关注列表减少一位
		if err := tx.Exec("update users set follow_count = follow_count-1 where id = ?", UserID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (f *followDealer) QueryIsFollow(UserID int64, ToUserID int64) (bool, error) {
	f.Context.TraceCaller()
	res := new(int8)
	err := db.Raw("select 1 from follow where user_id = ? and to_user_id = ? limit 1", UserID, ToUserID).Scan(res).Error
	return *res == 1, err
}
