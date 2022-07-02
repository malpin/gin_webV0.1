package mysql

import (
	"fmt"
	"gin_web/settings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var MysqlDB *sqlx.DB

func Init(config *settings.MysqlConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Dbname,
	)
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	MysqlDB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB falled ,", zap.Error(err))
		return err
	}
	//设置与数据库的最大打开连接数。
	MysqlDB.SetMaxOpenConns(config.MaxOpenConns)
	//最大闲置数
	MysqlDB.SetMaxIdleConns(config.MaxIdleConns)

	return
}
func Close() {
	err := MysqlDB.Close()
	if err != nil {
		zap.L().Fatal("mysql close err: ", zap.Error(err))
	}
}
