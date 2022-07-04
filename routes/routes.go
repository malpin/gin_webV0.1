package routes

import (
	"gin_web/controller"
	"gin_web/logger"
	"gin_web/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.POST("/login", controller.UserLogin)
	r.POST("/signup", controller.UserSignUp)
	r.GET("/ping", tool.JWTAuthMiddleware(), func(c *gin.Context) {
		//测试
		c.String(http.StatusOK, "pang")
	})

	return r
}
