package model

import "time"

type User struct {
	UserId     int64     `db:"user_id"`
	Username   string    `db:"username"`
	Password   string    `db:"password"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

func NewUser(userId int64, username string, password string, createTime time.Time, updateTime time.Time) *User {
	return &User{UserId: userId, Username: username, Password: password, CreateTime: createTime, UpdateTime: updateTime}
}
