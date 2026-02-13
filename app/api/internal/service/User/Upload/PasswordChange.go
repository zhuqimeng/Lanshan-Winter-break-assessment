package Upload

import (
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"
	"zhihu/utils/Strings"
	"zhihu/utils/randoms"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UpdPwd(c *gin.Context) {
	username, _ := c.Get("username")
	var user User.User
	res := configs.Db.Where("name = ?", username).First(&user)
	if res.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "用户不存在或已注销",
		})
		configs.Logger.Error("UpdPwd", zap.Any("username", username), zap.Error(res.Error))
		return
	}
	var req User.UpdPwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Strings.VerifyPassword(user.Password, req.OldPassword) == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码验证错误！"})
		return
	}
	salt, err := randoms.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash := Strings.HashPassword(req.NewPassword + salt)
	user.Password = salt + "_" + hash
	configs.Db.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	configs.Logger.Info("UpdPwd", zap.Any("username", username))
}
