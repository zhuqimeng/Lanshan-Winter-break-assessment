package Sign

import (
	"net/http"
	"time"
	"zhihu/app/api/internal/model/User"
	"zhihu/app/api/internal/service/User/dao"
	"zhihu/utils/tokens"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req User.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dao.ReadUser(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := tokens.MakeToken(req.Name, time.Now().Add(2*time.Hour))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refreshToken, err := tokens.MakeToken(req.Name, time.Now().Add(7*24*time.Hour))
	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}
