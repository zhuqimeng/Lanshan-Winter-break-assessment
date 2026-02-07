package dao

import (
	"io"
	"net/http"
	"os"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"
	"zhihu/utils/files"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetMdFile(c *gin.Context) {
	filetype := c.Param("filetype")
	thePath := c.Param("url")
	if filetype == "article" {
		thePath = "Storage/Document/Article/" + thePath
	} else if filetype == "question" {
		thePath = "Storage/Document/Question/" + thePath
	} else if filetype == "answer" {
		thePath = "Storage/Document/Answer/" + thePath
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
		return
	}
	fileInfo, err := os.Stat(thePath)
	if os.IsNotExist(err) {
		configs.Logger.Error("GetMdFile", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 验证文件是否存在

	file, err := os.Open(thePath)
	if err != nil {
		configs.Logger.Error("GetMdFile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "文件打开失败",
		})
	}
	defer func() {
		if err := file.Close(); err != nil {
			configs.Logger.Error("GetMdFile", zap.Error(err))
		}
	}()

	if !files.IsMarkdown(file.Name()) {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "非法的文件类型",
		})
		return
	}
	contentType := "text/markdown; charset=utf-8"
	files.HeaderSet(c, fileInfo, contentType)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		configs.Logger.Error("传输文件失败", zap.Error(err), zap.String("file", thePath))
		return
	}
}

func GetTop(c *gin.Context) {
	var (
		questions []Document.Question
		result    []gin.H
		orderBy   string
	)
	sortBy := c.Query("OrderBy")
	if sortBy == "like" {
		orderBy = "like_num DESC"
	} else if sortBy == "time" {
		orderBy = "created_at DESC"
	} else {
		orderBy = "like_num DESC"
	}
	if err := configs.Db.Model(&Document.Question{}).Order(orderBy).Find(&questions).Error; err != nil {
		configs.Logger.Error("GetTop", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	for _, q := range questions {
		result = append(result, gin.H{
			"author":   q.Username,
			"title":    q.Title,
			"url":      q.URL,
			"like_num": q.LikeNum,
			"time":     q.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
