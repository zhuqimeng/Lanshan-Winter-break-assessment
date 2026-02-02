package router

import (
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/middleware/Auth"
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

	userR := r.Group("/user")
	{
		userR.GET("/:username/homepage/:filename", User.GetHome)
		userR.GET("/:username/articles", DocumentDao.GetUserArticles)
		userR.GET("/:username/questions", DocumentDao.GetUserQuestions)
	}
	// 浏览用户路由

	broR := r.Group("/browse")
	{
		broR.GET("/:filetype/:url", DocumentDao.GetMdFile)
	}
	// 浏览网页路由

	if err := r.Run(":8080"); err != nil {
		configs.Logger.Fatal("Run error", zap.Error(err))
	}
}
