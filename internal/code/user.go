package code

type ResCode int64

const Success ResCode = 0

const (
	NeedLogin ResCode = 10000 + iota
	InvalidToken
	InvalidParam
	WrongPassword
	UserExist
	UserNotExist
	ServeBusy
)

var msgMAP = map[ResCode]string{
	//  成功响应
	Success: "success",
	//=======用户===============================
	NeedLogin:     "请登录后重试",
	InvalidToken:  "无效token",
	InvalidParam:  "参数错误",
	WrongPassword: "用户名或密码错误",
	ServeBusy:     "服务繁忙",
	UserExist:     "用户名已存在",
	UserNotExist:  "用户不存在",
	//=======视频=============================================
	InvalidFile:   "文件格式错误",
	VideoNotExist: "视频不存在",
	//=================社交=============
	RepeatFollow:   "无法重复关注",
	RepeatUnFollow: "无法重复取消关注",
	FollowSelf:     "不能取关或关注自己",
	// 点赞模块
	RepeatLiked:   "重复点赞",
	RepeatUnLiked: "重复取消关注",
	//评论模块
	EmptyComment:    "评论为空",
	CommentNotExist: "评论不存在",
}

func (r ResCode) MSG() string {
	return msgMAP[r]
}
