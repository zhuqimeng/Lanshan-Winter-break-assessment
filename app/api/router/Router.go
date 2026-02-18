package router

import (
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/middleware/Auth"
	"zhihu/app/api/internal/service/Document/Answer"
	"zhihu/app/api/internal/service/Document/Comment"
	DocumentDao "zhihu/app/api/internal/service/Document/dao"
	"zhihu/app/api/internal/service/Message"
	"zhihu/app/api/internal/service/User"
	"zhihu/app/api/internal/service/User/Follow"
	"zhihu/app/api/internal/service/User/Sign"
	"zhihu/app/api/internal/service/User/Upload"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Router(hub *Message.Hub) {
	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20

	r.GET("", react)
	r.GET("/top", DocumentDao.GetTop)
	r.GET("/feed", Auth.TokenChecker(), Follow.ShowFeeds)
	// 默认主页和关注动态

	r.GET("/search", DocumentDao.SearchKeyword)
	r.GET("/chat", Auth.TokenChecker(), hub.HandleWebSocket)
	// 文章搜索以及用户私信

	r.GET("/ping", pong)
	r.POST("/refresh", refresh)
	// 基本响应和token刷新

	r.POST("/register", Sign.Register)
	r.POST("/login", Sign.Login)
	// 用户登录注册

	updR := r.Group("/update")
	updR.Use(Auth.TokenChecker())
	{
		updR.POST("/password", Upload.UpdPwd)
		updR.POST("/homepage/:filetype", Upload.HomeUpd)
		updR.POST("/new/:filetype", DocumentDao.Create)
	}
	// 上传文件路由

	userR := r.Group("/user/:username")
	{
		userR.GET("", DocumentDao.GetUserInfo)
		userR.GET("homepage", User.GetHome)
		userR.GET("following", Follow.GetFollowing)
		userR.GET("follower", Follow.GetFollower)
	}
	// 浏览用户路由

	broR := r.Group("/browse/:filetype/:url")
	{
		broR.GET("", DocumentDao.GetMdFile)
		broR.GET("comment", Comment.Read)
		broR.POST("comment/create", Auth.TokenChecker(), Auth.FrequencyChecker("comment"), Comment.Create)
		broR.GET("answer", Answer.Read)
		broR.POST("answer/create", Auth.TokenChecker(), Answer.Create)
	}
	// 浏览文件路由

	r.PUT("/like/:filetype/:url", Auth.TokenChecker(), Auth.FrequencyChecker("like"), DocumentDao.ChangeLike)
	r.PUT("/follow/:username", Auth.TokenChecker(), Auth.FrequencyChecker("follow"), Follow.OnFollow)
	r.PUT("/unfollow/:username", Auth.TokenChecker(), Auth.FrequencyChecker("follow"), Follow.OffFollow)
	// 功能按键

	if err := r.Run(":8080"); err != nil {
		configs.Logger.Fatal("Run error", zap.Error(err))
	}
}
