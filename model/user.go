package model

import "time"

type User struct {
	UserId     int64     `db:"user_id" json:"userid,string"`
	Username   string    `db:"username" json:"username"`
	Password   string    `db:"password" json:"password,omitempty"`
	CreateTime time.Time `db:"create_time" json:"createTime,omitempty"`
	UpdateTime time.Time `db:"update_time" json:"updateTime,omitempty"`
	UserToken  string    `db:"_" json:"token,omitempty"`
}

func NewUser(userId int64, username string, password string, createTime time.Time, updateTime time.Time) *User {
	return &User{UserId: userId, Username: username, Password: password, CreateTime: createTime, UpdateTime: updateTime}
}
