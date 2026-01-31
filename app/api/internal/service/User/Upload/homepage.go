package Upload

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"
	"zhihu/utils/files"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HomeUpd(c *gin.Context) {
	username, _ := c.Get("username")
	filetype := c.Param("filetype")
	var file *multipart.FileHeader
	var err error
	if filetype == "avatar" {
		file, err = c.FormFile("avatar")
	} else if filetype == "profile" {
		file, err = c.FormFile("profile")
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "错误的请求信息",
		})
		return
	}
	if err != nil {
		configs.Sugar.Error("HomeUpd", err)
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
	var allowedTypes map[string]bool
	if filetype == "avatar" {
		allowedTypes = map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}
	} else {
		allowedTypes = map[string]bool{
			"text/markdown": true,
		}
	}

	fileHeader, err := file.Open()
	if err != nil {
		configs.Sugar.Error("HomeUpd", "文件打开失败", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "无法打开文件",
		})
		return
	}
	defer func() {
		if err := fileHeader.Close(); err != nil {
			configs.Sugar.Error("HomeUpd", "文件关闭失败", err)
		}
	}()

	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		configs.Sugar.Error("HomeUpd", "文件读取失败", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "无法读取文件",
		})
		return
	}

	contentType := mimetype.Detect(buffer).String()
	if files.IsMarkdown(file.Filename) {
		contentType = "text/markdown"
	}

	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":  http.StatusBadGateway,
			"error": "不支持的文件类型",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s_%d%s", username, timestamp, ext)
	var thePath string
	if filetype == "avatar" {
		thePath = fmt.Sprintf("Storage/User/Avatar/%s", filename)
	} else {
		thePath = fmt.Sprintf("Storage/User/Profile/%s", filename)
	}
	if err := c.SaveUploadedFile(file, thePath); err != nil {
		configs.Logger.Error("HomeUpd", zap.Any("username", username), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "文件保存失败",
		})
		return
	}

	// 开始数据库事务
	var URL string
	if filetype == "avatar" {
		URL = "AvatarURL"
	} else {
		URL = "ProfileURL"
	}
	tx := configs.Db.Begin()
	if err := tx.Model(&User.User{}).Where("name = ?", username).Update(URL, thePath).Error; err != nil {
		tx.Rollback()
		configs.Logger.Error("HomeUpd", zap.Any("username", username), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "文件上传到数据库失败",
		})
	}
	tx.Commit()

	configs.Logger.Info("HomeUpd", zap.Any("username", username), zap.String("url", thePath))
	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "上传成功",
		"username": username,
		"filename": filename,
		"size":     file.Size,
		"url":      thePath,
	})
}
