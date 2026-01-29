package router

import (
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/middleware/Auth"
	"zhihu/app/api/internal/service/User/Sign"
	"zhihu/app/api/internal/service/User/Upload"
	"zhihu/app/api/internal/service/User/browse"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Router() {
	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20
	r.GET("/ping", Auth.TokenChecker(), pong)
	r.GET("/refresh", refresh)
	// 基本响应和token刷新

	r.POST("/register", Sign.Register)
	r.GET("/login", Sign.Login)
	// 用户登录注册

	updR := r.Group("/update")
	updR.Use(Auth.TokenChecker())
	{
		updR.POST("/password", Upload.UpdPwd)
		updR.POST("/avatar", Upload.AvatarUpd)
	}
	// 上传文件路由

	getR := r.Group("/browse")
	getR.Use(Auth.TokenChecker())
	{
		getR.GET("/homepage/:username", browse.GetHome)
	}
	// 浏览网站路由

	if err := r.Run(":8080"); err != nil {
		configs.Logger.Fatal("Run error", zap.Error(err))
	}
}
