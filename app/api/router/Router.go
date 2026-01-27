package router

import (
	"log"
	"zhihu/app/api/internal/middleware/Auth"
	"zhihu/app/api/internal/service/User/Sign"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	r.POST("/register", Sign.Register)
	r.GET("/login", Sign.Login)
	r.GET("/ping", Auth.TokenChecker(), pong)
	r.GET("/refresh", refresh)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
