package Upload

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AvatarUpd(c *gin.Context) {
	username, _ := c.Get("username")
	file, err := c.FormFile("avatar")
	if err != nil {
		configs.Sugar.Error("AvatarUpd", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  err.Error(),
		})
		return
	}

	if file.Size > 5<<20 {
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "文件太大，最大支持 5MB",
		})
		return
	}
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	fileHeader, err := file.Open()
	if err != nil {
		configs.Sugar.Error("AvatarUpd", "文件打开失败", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "无法打开文件",
		})
		return
	}
	defer func() {
		if err := fileHeader.Close(); err != nil {
			configs.Sugar.Error("AvatarUpd", "文件关闭失败", err)
		}
	}()

	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		configs.Sugar.Error("AvatarUpd", "文件读取失败", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "无法读取文件",
		})
		return
	}

	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":  http.StatusBadGateway,
			"error": "不支持的文件类型",
		})
		return
	}

	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s_%d_%s", username, timestamp, ext)
	thePath := fmt.Sprintf("Storage/User/Avatar/%s", filename)
	if err := c.SaveUploadedFile(file, thePath); err != nil {
		configs.Logger.Error("AvatarUpd", zap.Any("username", username), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "文件保存失败",
		})
		return
	}

	// 开始数据库事务
	tx := configs.Db.Begin()
	if err := tx.Model(&User.User{}).Where("name = ?", username).Update("AvatarURL", thePath).Error; err != nil {
		tx.Rollback()
		configs.Logger.Error("AvatarUpd", zap.Any("username", username), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "头像上传到数据库失败",
		})
	}
	tx.Commit()

	configs.Logger.Info("AvatarUpd", zap.Any("username", username), zap.String("url", thePath))
	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "上传成功",
		"username": username,
		"filename": filename,
		"size":     file.Size,
		"url":      thePath,
	})
}
