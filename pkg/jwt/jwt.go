package jwt

import (
	"tiktink/pkg/logger"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

const (
	TokenExpireDuration = time.Hour * 24 * 365
)

var mySecret = []byte("我来到我看见我记录")

func keyFunc(token *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

type MyClaims struct {
	UserID             int64  `json:"id"`
	Username           string `json:"username"`
	jwt.StandardClaims        //内嵌匿名字段，但如果给这个结构体加上名称的话表示内嵌的是一个独立的结构体
	//匿名结构体实现继承的效果，有名字的结构体实现组合效果。如果内嵌多个匿名结构体则可实现类似于多重继承的效果
}

func GenToken(userID int64, userName string) (aToken string, err error) {
	// 建立自己的token字段
	c := MyClaims{
		userID,
		userName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    viper.GetString("app.name"),
		},
	}
	// 加密并获得的完整编码后的aToken
	if aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret); err != nil {
		return "", err
	}

	return
}

// ParseToken 解析token返回包含用户消息的结构体
func ParseToken(tokenString string) (*MyClaims, bool, error) {
	// 解析token
	var claims = new(MyClaims) //解析好的存放在mc中

	//keyFunc 用来根据token值返回密钥
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil { //解析过程出错
		logger.PrintLog("解析token错误:", err)
		return nil, false, err
	}
	if token.Valid { // 校验token是否合法
		return claims, true, nil
	}
	return nil, false, nil
}

// RefreshToken 后期改为双token的鉴权模式
func RefreshToken(aToken, rToken string) (newAToken string, err error) {
	// 由于refresh token 不携带额外参数，直接使用parse而不是parseWithClaims
	// refresh Token 过期直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return "", err
	}

	var claims = new(MyClaims)
	//  从aToken解析数据绑定到claims上
	_, err = jwt.ParseWithClaims(aToken, claims, keyFunc)

	v, _ := err.(*jwt.ValidationError)

	//不是过期错误
	if v.Errors != jwt.ValidationErrorExpired {
		return "", err
	}
	return GenToken(claims.UserID, claims.Username)
}
