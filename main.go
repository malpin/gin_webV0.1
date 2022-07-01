package main

import (
	"context"
	"fmt"
	"gin_web/dao/mysql"
	"gin_web/dao/redis"
	"gin_web/logger"
	"gin_web/routes"
	"gin_web/settings"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//1,加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("configuration file err:%v\n", err)

		return
	}

	//2.使用zap记录相关日志
	//settings.Conf.LogConfig 来自配置文件初始化时候的全局变量,这个变量用了结构体保存配置,在配置变化时候会重新改变结构体
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed ,err:%v\n", err)
		return
	}
	defer zap.L().Sync() //Sync 调用底层Core的 Sync 方法，刷新所有缓冲的日志条目。应用程序应注意在退出之前调用 Sync。
	zap.L().Debug("logger init success")
	//3.初始化mysql
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	defer mysql.Close()

	//4.初始化redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()

	//5.注册路由
	router := routes.Setup()
	//6.启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: router,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Fatal("Server exiting")
}
