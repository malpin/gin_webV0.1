package service

import (
	"gin_web/Bean"
	"gin_web/dao/mysql/userDao"
	"gin_web/model"
	"gin_web/tool"
	"time"
)

// SingUp 添加用户
func SingUp(p model.ParamsSignUp) (int64, error) {
	//根据用户名查找用户是否存在
	err := userDao.FandUserByName(p.Username)
	if err != nil {
		return -1, Bean.ErrorUserExist
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

// Login 用户登录
func Login(p model.ParamsLogin) error {
	u := model.User{
		Username: p.Username,
		Password: tool.Md5(p.Password),
	}
	err := userDao.Login(&u)
	if err != nil {
		return err
	}
	return nil
}
