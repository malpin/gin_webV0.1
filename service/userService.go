package service

import (
	"fmt"
	"gin_web/Bean"
	"gin_web/dao/mysql/userDao"
	"gin_web/model"
	"gin_web/tool"
	"go.uber.org/zap"
	"time"
)

// SingUp 添加用户
func SingUp(p *model.ParamsSignUp) (int64, error) {
	//根据用户名查找用户是否存在
	err := userDao.FandUserByName(p.Username)
	if err != nil {
		return -1, err
	}
	//生成一个id
	id := tool.GenID()
	//创建一个user对象
	user := model.NewUser(id, p.Username, tool.Md5(p.Password), time.Now(), time.Now())
	//添加用户
	addUserId, err := userDao.AddUser(user)
	if err != nil {
		return -1, err
	}
	return addUserId, nil
}

// Login 用户登录
func Login(p *model.ParamsLogin) (user model.User, err error) {
	//创建user对象,把密码加密
	u := model.User{
		Username: p.Username,
		Password: tool.Md5(p.Password),
	}
	user, err = userDao.Login(&u)
	if err != nil {
		return user, err
	}

	//生成token
	token, err := tool.OutputToken(user.UserId, user.Username)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("用户id:%d ,用户名%s ,在%s时登录 jwt 生成Token失败了", user.UserId, user.Username, time.Now()), zap.Error(err))
		return user, Bean.SYSTEM_BUSY.MarkError
	}
	user.UserToken = token

	return user, nil
}
