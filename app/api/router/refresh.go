package router

import (
	"net/http"
	"time"
	"zhihu/app/api/configs"
	"zhihu/utils/tokens"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		configs.Logger.Error("refresh", zap.Error(err))
		return
	}
	username, err := tokens.CheckToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": err.Error(),
		})
		configs.Logger.Error("refresh", zap.Error(err))
		return
	}
	token, err := tokens.MakeToken(username, time.Now().Add(2*time.Hour))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		configs.Logger.Error("refresh", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":             http.StatusOK,
		"new_access_token": token,
	})
	configs.Logger.Info("refresh", zap.String("status", "success"))
}
