package Auth

import (
	"errors"
	"net/http"
	"zhihu/utils/tokens"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "需要登录才能访问",
			})
			c.Abort()
			return
		}
		username, err := tokens.CheckToken(authHeader)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Token已过期，请重新登录",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Token无效",
					"error":   err.Error(),
				})
			}
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
