package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func react(c *gin.Context) {
	tips :=
		`若您已经有账号，请点击 /login 登录;
还没有账号？ 点击 /register 注册。
当然您也可以以游客身份访问此网站。`
	c.JSON(http.StatusOK, gin.H{
		"tips": tips,
	})
}
