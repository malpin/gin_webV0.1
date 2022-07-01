package service

import (
	"gin_web/dao/mysql/userDao"
	"gin_web/tool"
)

func SingUp() {
	//根据用户名查找用户是否存在
	userDao.FandUserByName()
	//生成一个id
	tool.GenID()
	//添加用户
	userDao.AddUser()
}
