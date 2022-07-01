package controller

import (
	"gin_web/service"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	//获取校验参数

	service.SingUp()
}
