package Bean

import "errors"

var (
	ErrorUserExist       = errors.New("用户已经存在了")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)
