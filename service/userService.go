package service

import (
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
		return -1, err
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
