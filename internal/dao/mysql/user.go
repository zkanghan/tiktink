package mysql

import (
	"database/sql"
	"tiktink/internal/model"
	"tiktink/pkg/snowid"
	"tiktink/pkg/tracer"

	"github.com/pkg/errors"
)

type userFunc interface {
	CreateUser(Username, password string) (string, error)
	QueryLoginParams(username string) (string, string, error)
	QueryNameByID(id string) (string, error)
	QueryUserExistByID(id string) (bool, error)
	QueryUserByID(userId string) (*model.UserMSG, error)
	QueryUserExistByName(username string) (bool, error)
}

type userDealer struct{}

func NewUserDealer() userFunc {
	return &userDealer{}
}

func (u *userDealer) QueryUserByID(userId string) (*model.UserMSG, error) {
	userMsg := new(model.UserMSG)
	err := db.Raw("select `user_id`,`user_name`,`follow_count`,`follower_count` from `users` where `user_id` = ?", userId).Scan(userMsg).Error
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(userId))
	}
	return userMsg, nil
}

func (u *userDealer) QueryUserExistByName(username string) (bool, error) {
	res := new(int8)
	err := db.Raw("select 1 from users where  user_name = ? limit 1", username).Scan(res).Error
	if err != nil {
		return false, errors.Wrap(err, tracer.FormatParam(username))
	}
	return *res == 1, nil
}

// QueryUserExistByID 混用会导致性能下降
func (u *userDealer) QueryUserExistByID(id string) (bool, error) {
	res := new(int8)
	err := db.Raw("select 1 from users where  user_id = ? limit 1", id).Scan(res).Error
	if err != nil {
		return false, errors.Wrap(err, tracer.FormatParam(id))
	}
	return *res == 1, nil
}

func (u *userDealer) QueryNameByID(id string) (username string, err error) {
	userName := new(string)
	err = db.Raw("select user_id from users where user_name = ?", id).Scan(userName).Error
	if err != nil {
		return "", errors.Wrap(err, tracer.FormatParam(id))
	}
	return *userName, nil
}

// CreateUser 用户注册，返回用户id
func (u *userDealer) CreateUser(Username, password string) (userID string, err error) {
	user := &model.User{
		UserID:   snowid.GenID(),
		UserName: Username,
		Password: password,
	}
	if err = db.Create(user).Error; err != nil {
		return "", errors.Wrap(err, tracer.FormatParam(Username, password))
	}
	return user.UserID, nil
}

// QueryLoginParams 查询用户id 和密码
func (u *userDealer) QueryLoginParams(username string) (password string, userID string, err error) {
	var rows *sql.Rows
	rows, err = db.Raw("select password,user_id from users where user_name = ?", username).Rows()
	if err != nil {
		return "", "", errors.Wrap(err, tracer.FormatParam(username))
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&password, &userID); err != nil {
			return "", "", errors.Wrap(err, tracer.FormatParam(username))
		}
	}
	return password, userID, nil
}
