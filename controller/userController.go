package controller

import (
	"errors"
	"fmt"
	"gin_web/Bean"
	"gin_web/model"
	"gin_web/service"
	"gin_web/tool"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserSignUp 用户注册
func UserSignUp(c *gin.Context) {
	//获取校验参数
	var p model.ParamsSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Debug("UserSignUp Parse data error", zap.Error(err))
		//转换错误格式
		errors, ok := err.(validator.ValidationErrors)
		//判断是否是ValidationErrors之中的错误
		if !ok {
			Bean.Error(c, Bean.DATA_ERROR)
			return
		}
		fmt.Println("执行了---")
		Bean.ErrorWithMsg(c, Bean.DATA_ERROR, errors.Translate(tool.ValidatorTrans)) //翻译错误信息
		return
	}
	//业务处理
	userID, err := service.SingUp(p)
	if err != nil {
		zap.L().Error("service SingUp error 添加失败了", zap.Error(err))

		if errors.Is(err, Bean.ErrorUserExist) {
			//返回
			Bean.Error(c, Bean.ADMIN_USERNAME_EXIST)
			return
		}
		//返回
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	//返回结果
	Bean.Success(c, userID)
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	//获取参数
	var p model.ParamsLogin
	if err := c.ShouldBindJSON(&p); err != nil { //ShouldBindJSON(&p)
		zap.L().Error("UserLogin Parse data error", zap.Error(err))
		//转换错误格式
		errors, ok := err.(validator.ValidationErrors)
		//判断是否是ValidationErrors之中的错误
		if !ok {
			Bean.Error(c, Bean.DATA_ERROR)
			return
		}
		Bean.ErrorWithMsg(c, Bean.DATA_ERROR, errors.Translate(tool.ValidatorTrans)) //翻译错误信息
		return
	}
	//业务处理
	err := service.Login(p)
	if err != nil {
		zap.L().Error("service UserLogin error 用户登录失败了", zap.Error(err))
		if errors.Is(err, Bean.ErrorInvalidPassword) {
			Bean.Error(c, Bean.ADMIN_PASSWORD_ERROR)
			return
		}
		//返回
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	//返回结果
	Bean.Success(c, "登录成功了")
}
