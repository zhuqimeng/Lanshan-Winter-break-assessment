package Sign

import (
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"
	"zhihu/app/api/internal/service/User/dao"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(c *gin.Context) {
	var req User.CreateUserReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dao.CreateUser(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
	configs.Logger.Info("register success", zap.String("username", req.Name))
}
