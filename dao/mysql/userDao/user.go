package userDao

import (
	"errors"
	"fmt"
	"gin_web/Bean"
	"gin_web/dao/mysql"
	"gin_web/model"
	"go.uber.org/zap"
	"time"
)

// AddUser 添加用户
func AddUser(user *model.User) (int64, error) {
	sql := "insert into user(user_id,username,password,create_time,update_time) values(?,?,?,?,?)"
	result, err := mysql.MysqlDB.Exec(sql, user.UserId, user.Username, user.Password, user.CreateTime, user.UpdateTime)
	if err != nil {
		zap.L().Error("userDao AddUser 添加用户执行 失败了", zap.Error(err))
		return -1, err
	}
	//插入成功会返回自增的id
	id, err := result.LastInsertId()
	if err != nil {
		zap.L().Error(fmt.Sprintf("userDao AddUser 添加用户 失败了 ,用户id为:%d", user.UserId), zap.Error(err))
		return -1, err
	}
	zap.L().Debug(fmt.Sprintf("userDao AddUser 添加用户 成功了 ,用户id为:%d", user.UserId), zap.Error(err))
	return id, nil
}

// login 用户登录
func Login(u *model.User) (err error) {
	p := u.Password
	var user model.User
	sql := "select user_id,username,password from user where username=? and password=?"
	err = mysql.MysqlDB.Get(&user, sql, u.Username, u.Password)
	if err != nil {
		zap.L().Debug("userDao Login 用户登录执行 失败了", zap.Error(err))
		return err
	}
	if p != u.Password {
		zap.L().Debug(fmt.Sprintf("用户id:%d 的用户 在%s时登录失败了", u.UserId, time.Now()), zap.Error(err))
		return Bean.ErrorInvalidPassword
	}

	zap.L().Debug(fmt.Sprintf("用户id:%d 的用户 在%s时登录成功了", u.UserId, time.Now()), zap.Error(err))
	return nil
}

//根据用户名删除用户
func DeleteUserById() {

}

//更新用户信息
func UpdateUser() {

}

// FandUserByName 根据用户名查找用户是否存在
func FandUserByName(username string) error {
	sql := "select count(user_id) from user  where username = ?"
	var count int
	err := mysql.MysqlDB.Get(&count, sql, username)
	if err != nil {
		zap.L().Error("userDao FandUserByName 根据用户名查找用户是否存在 查询失败", zap.Error(err))
		return err
	}
	if count > 0 {
		return errors.New("用户存在了")
	}
	return nil
}

//根据id查找用户
func FindUserById() {

}

//查询所有用户信息
func FindAllUser() {

}
