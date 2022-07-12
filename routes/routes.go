package routes

import (
	"gin_web/controller"
	_ "gin_web/docs"
	"gin_web/logger"
	"gin_web/tool"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	//登录
	r.POST("/login", controller.UserLogin)
	//登出
	r.POST("/loginout", controller.UserLoginOut)
	//注册
	r.POST("/signup", controller.UserSignUp)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(tool.VerifyToken()) //认证中间件
	//获取社区列表
	r.GET("/communityList", controller.GetCommunityList)
	//按id获取社区详情
	r.GET("/community/:id", controller.GetIntroductionById)
	//发布帖子
	r.POST("/releasePost", controller.ReleasePost)
	//根据id获取帖子详情
	r.GET("/post/:id", controller.ReceivePostById)
	//获取帖子列表
	r.GET("/postList", controller.GetPostListOrder)
	//根据社区获取帖子信息
	//r.GET("/postListCommunity", controller.GetPostListOrderByCommunity)
	//点赞
	r.POST("/vote", controller.PostVote)
	return r
}
