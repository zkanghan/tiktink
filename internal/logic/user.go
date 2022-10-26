package logic

import (
	"tiktink/internal/dao/mysql"
	"tiktink/internal/model"
	"tiktink/pkg/tracer"

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

type userDealer struct {
	Context *tracer.TraceCtx
}

type userFunc interface {
	CreateUser(username string, password string) (string, error)
	CheckUser(username, password string) (bool, string, error)
	GetUserExistByName(username string) (bool, error)
	GetUserExistByID(id string) (bool, error)
	GetUserInformation(toQueryUserID string, userID string) (*model.UserMSG, error)
}

func NewUserDealer(ctx *tracer.TraceCtx) userFunc {
	return &userDealer{
		Context: ctx,
	}
}

func (u *userDealer) CreateUser(username string, password string) (string, error) {
	//  还要查询用户是不是已经存在
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	return mysql.NewUserDealer().CreateUser(username, hashedPassword)
}

// CheckUser 校验用户用户名和密码是否正确，正确返回用户id
func (u *userDealer) CheckUser(username, password string) (bool, string, error) {
	todo := mysql.NewUserDealer()
	//  从数据库获取密码
	DBpassword, id, err := todo.QueryLoginParams(username)
	if err != nil {
		return false, "", err
	}
	//  密码匹配错误，非运行异常
	if err := comparePassword(DBpassword, password); err != nil {
		return false, "", nil
	}
	return true, id, nil
}

// GetUserExistByName 返回true表示用户存在，false表示不存在
func (u *userDealer) GetUserExistByName(username string) (bool, error) {
	return mysql.NewUserDealer().QueryUserExistByName(username)
}

func (u *userDealer) GetUserExistByID(id string) (bool, error) {
	return mysql.NewUserDealer().QueryUserExistByID(id)
}

func (u *userDealer) GetUserInformation(toQueryUserID string, userID string) (*model.UserMSG, error) {
	userMsg, err := mysql.NewUserDealer().QueryUserByID(toQueryUserID)
	if err != nil {
		return nil, err
	}
	relationer := NewRelationDealer(u.Context)
	followed, err := relationer.GetIsFollowed(userID, toQueryUserID)
	if err != nil {
		return nil, err
	}
	userMsg.IsFollow = followed
	return userMsg, nil
}
