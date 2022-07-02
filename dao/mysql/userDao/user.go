package userDao

import (
	"errors"
	"fmt"
	"gin_web/dao/mysql"
	"gin_web/model"
	"go.uber.org/zap"
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
