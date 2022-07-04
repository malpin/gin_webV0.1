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
	"time"
)

// UserSignUp 用户注册
func UserSignUp(c *gin.Context) {
	//获取校验参数
	var p model.ParamsSignUp
	//解析参数
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Warn("Parse data error", zap.Error(err))
		//转换错误格式
		errors, ok := err.(validator.ValidationErrors)
		//判断是否是ValidationErrors之中的错误
		if ok {
			Bean.ErrorWithMsg(c, Bean.DATA_ERROR, errors.Translate(tool.ValidatorTrans)) //全局翻译器 翻译错误信息
			return
		}
		Bean.Error(c, Bean.DATA_ERROR)
		return
	}
	//业务处理
	userID, err := service.SingUp(&p)
	if err != nil {
		zap.L().Warn("service SingUp error 添加失败了", zap.Error(err))
		if errors.Is(err, Bean.USERNAME_EXIST.MarkError) {
			//返回
			Bean.Error(c, Bean.USERNAME_EXIST)
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
	//解析参数
	if err := c.ShouldBindJSON(&p); err != nil { //ShouldBindJSON(&p)
		zap.L().Warn("Parse data error", zap.Error(err))
		//转换错误格式
		errors, ok := err.(validator.ValidationErrors)
		//判断是否是ValidationErrors之中的错误
		if ok {
			Bean.ErrorWithMsg(c, Bean.DATA_ERROR, errors.Translate(tool.ValidatorTrans)) //全局翻译器 翻译错误信息
			return
		}
		Bean.Error(c, Bean.DATA_ERROR)
		return
	}
	//业务处理
	user, err := service.Login(&p)
	if err != nil || user.UserToken == "" {
		zap.L().Warn(fmt.Sprintf("用户id:%d ,用户名%s ,在%s时登录失败了", user.UserId, user.Username, time.Now()), zap.Error(err))
		if errors.Is(err, Bean.PASSWORD_ERROR.MarkError) {
			Bean.Error(c, Bean.PASSWORD_ERROR)
			return
		}
		//返回
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}

	Bean.Success(c, user)
}
