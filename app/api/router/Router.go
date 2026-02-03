package router

import (
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/middleware/Auth"
	"zhihu/app/api/internal/middleware/Auto"
	"zhihu/app/api/internal/service/Document/Answer"
	"zhihu/app/api/internal/service/Document/Comment"
	DocumentDao "zhihu/app/api/internal/service/Document/dao"
	"zhihu/app/api/internal/service/User"
	"zhihu/app/api/internal/service/User/Sign"
	"zhihu/app/api/internal/service/User/Upload"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Router() {
	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20
	r.GET("", react)
	r.GET("/ping", pong)
	r.GET("/refresh", refresh)
	// 基本响应和token刷新

	r.POST("/register", Sign.Register)
	r.GET("/login", Sign.Login)
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
		userR.GET("homepage/:filename", User.GetHome)
		userR.GET("articles", Auto.RouterSet("article"), DocumentDao.GetUserInfo)
		userR.GET("questions", Auto.RouterSet("question"), DocumentDao.GetUserInfo)
		userR.GET("answers", Auto.RouterSet("answer"), DocumentDao.GetUserInfo)
		userR.GET("comments", Auto.RouterSet("comment"), DocumentDao.GetUserInfo)
	}
	// 浏览用户路由

	broR := r.Group("/browse/:filetype/:url")
	{
		broR.GET("", DocumentDao.GetMdFile)
		broR.GET("comment", Comment.Read)
		broR.POST("comment/create", Auth.TokenChecker(), Comment.Create)
		broR.GET("answer", Answer.Read)
		broR.POST("answer/create", Auth.TokenChecker(), Answer.Create)
	}
	// 浏览网页路由

	if err := r.Run(":8080"); err != nil {
		configs.Logger.Fatal("Run error", zap.Error(err))
	}
}
