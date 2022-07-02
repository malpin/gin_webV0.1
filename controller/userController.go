package controller

import (
	"gin_web/model"
	"gin_web/service"
	"gin_web/tool"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

func SignUp(c *gin.Context) {
	//获取校验参数
	var p model.ParamsSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("SignUp Parse data error", zap.Error(err))
		//转换错误格式
		errors, ok := err.(validator.ValidationErrors)
		//判断是否是ValidationErrors之中的错误
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": "提交的数据有误!!", //不是就直接返回了
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": errors.Translate(tool.ValidatorTrans), //翻译错误
		})
		return
	}
	//业务处理
	service.SingUp()
	//返回
	c.JSON(http.StatusOK, gin.H{
		"msg": "成功了",
	})
}
