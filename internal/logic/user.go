package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/logger"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 返回加密后的密码
func hashPassword(password string) (string, error) {
	toHash := []byte(password)
	hashedPas, err := bcrypt.GenerateFromPassword(toHash, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashedPasString := string(hashedPas)
	return hashedPasString, nil
}

// ComparePassword 返回的error非空表示密码不匹配
func comparePassword(hashedPass string, password string) error {
	hashedPassByte := []byte(hashedPass)
	passwordByte := []byte(password)
	return bcrypt.CompareHashAndPassword(hashedPassByte, passwordByte)
}

func CreateUser(username string, password string) (int64, error) {
	//  还要查询用户是不是已经存在
	hashedPassword, err := hashPassword(password)
	if err != nil {
		logger.L.Error("密码加密失败：", zap.Error(err))
		return -1, err
	}
	return mysql.DealUser().CreateUser(username, hashedPassword)
}

// CheckUser 校验用户用户名和密码是否正确，正确返回用户id
func CheckUser(username, password string) (bool, int64, error) {
	todo := mysql.DealUser()
	//  从数据库获取密码
	DBpassword, id, err := todo.QueryLoginParams(username)
	if err != nil {
		return false, -1, err
	}
	//  密码匹配错误，非运行异常
	if err := comparePassword(DBpassword, password); err != nil {
		return false, -1, nil
	}
	return true, id, nil
}

// GetUserExistByName 返回true表示用户存在，false表示不存在
func GetUserExistByName(username string) (bool, error) {
	return mysql.DealUser().QueryUserExistByName(username)
}

func GetUserExistByID(id int64) (bool, error) {
	return mysql.DealUser().QueryUserExistByID(id)
}

func GetUserInformation(toQueryUserID int64, userID int64) (*model.UserMSG, error) {
	userMsg, err := mysql.DealUser().QueryUserByID(toQueryUserID)
	if err != nil {
		logger.L.Error("查询用户信息失败：", zap.Error(err))
		return nil, err
	}
	followed, err := GetIsFollowed(userID, toQueryUserID)
	if err != nil {
		logger.L.Error("查询是否关注失败：", zap.Error(err))
		return nil, err
	}
	userMsg.IsFollow = followed
	return userMsg, nil
}
