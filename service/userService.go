package service

import (
	"errors"
	"gin_web/dao/mysql/userDao"
	"gin_web/model"
	"gin_web/tool"
	"time"
)

var (
	ErrorUserExist       = errors.New("用户已经存在了")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
)

// SingUp 添加用户
func SingUp(p model.ParamsSignUp) (int64, error) {
	//根据用户名查找用户是否存在
	err := userDao.FandUserByName(p.Username)
	if err != nil {
		return -1, ErrorUserExist
	}
	//生成一个id
	id := tool.GenID()
	user := model.NewUser(id, p.Username, tool.Md5(p.Password), time.Now(), time.Now())
	//添加用户
	addUserId, err := userDao.AddUser(user)
	if err != nil {
		return -1, err
	}
	return addUserId, nil
}
