package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	if username, ok := c.Get("username"); ok {
		c.JSON(http.StatusOK, gin.H{
			"message":  "pong",
			"username": username,
		})
	}
}
