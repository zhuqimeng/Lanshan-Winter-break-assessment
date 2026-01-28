package router

import (
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/middleware/Auth"
	"zhihu/app/api/internal/service/User/Sign"
	"zhihu/app/api/internal/service/User/Upload"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Router() {
	r := gin.Default()
	r.POST("/register", Sign.Register)
	r.GET("/login", Sign.Login)
	r.GET("/ping", Auth.TokenChecker(), pong)
	r.GET("/refresh", refresh)
	updR := r.Group("/update")
	updR.Use(Auth.TokenChecker())
	{
		updR.POST("/password", Upload.UpdPwd)
	}
	if err := r.Run(":8080"); err != nil {
		configs.Logger.Fatal("Run error", zap.Error(err))
	}
}
