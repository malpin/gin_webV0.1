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
	"strconv"
	"time"
)

func ReleasePost(c *gin.Context) {
	//获取参数
	var post model.PostMsg
	err := c.ShouldBindJSON(&post)
	if err != nil {
		//转换错误格式
		errors, ok := err.(validator.ValidationErrors)
		//判断是否是ValidationErrors之中的错误
		if ok {
			//返回
			Bean.ErrorWithMsg(c, Bean.DATA_ERROR, errors.Translate(tool.ValidatorTrans)) //全局翻译器 翻译错误信息
			return
		}
		fmt.Println("执行了")
		//返回
		Bean.Error(c, Bean.DATA_ERROR)
		return
	}

	genID := tool.GenID()
	post.PostID = genID
	uid, err := tool.GetContextUser(c)
	if err != nil {
		Bean.Error(c, Bean.DATA_ERROR)
		return
	}
	post.AuthorID = uid
	post.CreateTime = time.Now()
	post.UpdateTime = time.Now()
	post.Status = 1

	i, err := service.ReleasePost(&post)
	if err != nil {
		//返回
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	//返回
	Bean.Success(c, i)
}

// ReceivePostById 根据id获取帖子详情
func ReceivePostById(c *gin.Context) {
	param := c.Param("id")
	pid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	post, err := service.ReceivePostById(pid)
	if err != nil {
		if err == Bean.POST_ID_ERROR.MarkError {
			Bean.Error(c, Bean.POST_ID_ERROR)
			return
		}
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	Bean.Success(c, post)
}

// GetPostList 获取帖子列表
func GetPostList(c *gin.Context) {
	//获取分页参数
	page := c.Query("page")
	size := c.Query("size")
	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	s, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}

	postlist, err := service.GetPostList(p, s)
	if err != nil {
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	Bean.Success(c, postlist)
}

// GetPostList 获取帖子列表
func GetPostListOrder(c *gin.Context) {
	//默认请求参数
	p := model.ParamPostList{
		Page:  1,
		Size:  10,
		Order: model.OrderTime,
	}

	//获取请求参数
	//获取排序规则
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		Bean.Error(c, Bean.DATA_ERROR)
		return
	}
	//从redis中查询
	postlist, err := service.GetPostListOrder(&p)
	if err != nil {
		fmt.Println(err)
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	Bean.Success(c, postlist)
}
