package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init() (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB falled ,", zap.Error(err))
		return err
	}
	//设置与数据库的最大打开连接数。
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	//最大闲置数
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))

	return
}
func Close() {
	err := db.Close()
	if err != nil {
		zap.L().Fatal("mysql close err: ", zap.Error(err))
	}
}
