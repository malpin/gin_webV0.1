package controller

import (
	"fmt"
	"gin_web/Bean"
	"gin_web/model"
	"gin_web/service"
	"gin_web/tool"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//点赞
func PostVote(c *gin.Context) {
	//获取参数
	p := new(model.ParamVoteData)
	//解析参数
	if err := c.ShouldBindJSON(p); err != nil {
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
	fmt.Println(p)
	//获取当前userID
	userid, err := tool.GetContextUser(c)
	if err != nil {
		Bean.Error(c, Bean.DATA_ERROR)
		return
	}
	//投票
	err = service.PostVote(userid, p)
	if err != nil {
		Bean.Error(c, Bean.ERR_VOTE_REPEATED)
		return
	}
	Bean.Success(c, "操作成功")

}
