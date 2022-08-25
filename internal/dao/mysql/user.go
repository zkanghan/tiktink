package mysql

import "tiktink/internal/model"

type userFunc interface {
	CreateUser(Username, password string) (int64, error)
	QueryLoginParams(username string) (string, int64, error)
	QueryNameByID(id int64) (string, error)
	QueryUserExistByID(id int64) (bool, error)
	QueryUserByID(userId int64) (*model.UserMSG, error)
	QueryUserExistByName(username string) (bool, error)
}

type userDealer struct{}

func (u userDealer) QueryUserByID(userId int64) (*model.UserMSG, error) {
	userMsg := new(model.UserMSG)
	err := db.Raw("select `user_id`,`user_name`,`follow_count`,`follower_count` "+
		"from `users` where `user_id` = ?", userId).Scan(userMsg).Error
	if err != nil {
		return nil, err
	}
	return userMsg, nil
}

func (u userDealer) QueryUserExistByName(username string) (bool, error) {
	res := new(int8)
	err := db.Raw("select 1 from users where  user_name = ? limit 1", username).Scan(res).Error
	if err != nil {
		return false, err
	}
	return *res == 1, nil
}

// QueryUserExistByID 混用会导致性能下降
func (u userDealer) QueryUserExistByID(id int64) (bool, error) {
	res := new(int8)
	err := db.Raw("select 1 from users where  user_id = ? limit 1", id).Scan(res).Error
	if err != nil {
		return false, err
	}
	return *res == 1, nil
}

func (u userDealer) QueryNameByID(id int64) (string, error) {
	userName := new(string)
	err := db.Raw("select id from users where user_name = ?", id).Scan(userName).Error
	if err != nil {
		return "", err
	}
	return *userName, nil
}

func DealUser() userFunc {
	return &userDealer{}
}

// CreateUser 用户注册，返回用户id
func (u userDealer) CreateUser(Username, password string) (int64, error) {
	user := &model.User{
		UserName: Username,
		Password: password,
	}
	if err := db.Create(user).Error; err != nil {
		return -1, err
	}
	return user.ID, nil
}

// QueryLoginParams 查询用户id 和密码
func (u userDealer) QueryLoginParams(username string) (string, int64, error) {
	password := new(string)
	id := new(int64)
	rows, err := db.Raw("select password,user_id from users where user_name = ?", username).Rows()
	defer rows.Close()
	if err != nil {
		return "", -1, err
	}
	for rows.Next() {
		if err := rows.Scan(password, id); err != nil {
			return "", -1, err
		}
	}
	return *password, *id, nil
}
