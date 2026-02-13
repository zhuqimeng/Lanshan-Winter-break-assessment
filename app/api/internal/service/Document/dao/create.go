package dao

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"
	"zhihu/app/api/internal/service/User/Follow"
	"zhihu/utils/Strings"
	"zhihu/utils/files"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Create(c *gin.Context) {
	username := c.GetString("username")
	filetype := c.Param("filetype")
	var file *multipart.FileHeader
	var err error
	if filetype == "article" {
		file, err = c.FormFile("article")
	} else if filetype == "question" {
		file, err = c.FormFile("question")
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "错误的请求信息",
		})
		return
	}
	if err != nil {
		configs.Sugar.Error("DocUpd", err)
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
		"text/markdown": true,
	}
	fileHeader, err := file.Open()
	if err != nil {
		configs.Sugar.Error("DocUpd", "文件打开失败", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "无法打开文件",
		})
		return
	}
	defer func() {
		if err = fileHeader.Close(); err != nil {
			configs.Sugar.Error("DocUpd", "文件关闭失败", err)
		}
	}()

	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		configs.Sugar.Error("DocUpd", "文件读取失败", err)
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

	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d-%s%s", timestamp, username, ext)
	var thePath string
	if filetype == "article" {
		thePath = fmt.Sprintf("Storage/Document/Article/%s", filename)
	} else {
		thePath = fmt.Sprintf("Storage/Document/Question/%s", filename)
	}
	if err = c.SaveUploadedFile(file, thePath); err != nil {
		configs.Logger.Error("DocUpd", zap.Any("username", username), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "文件保存失败",
		})
		return
	}
	// 文件本地储存

	if filetype == "article" {
		mdText, _ := io.ReadAll(fileHeader)
		plainText := Strings.MdToPlainText(string(mdText))
		summary, err := configs.Llm.Summarize(plainText, 500)
		if err != nil {
			configs.Sugar.Error("Article Summary err", zap.Any("username", username), zap.Error(err))
		}
		article := &Document.Article{
			Username: username,
			Title:    file.Filename,
			URL:      "/browse/article/" + filename,
			Summary:  summary,
		}
		result := configs.Db.Create(article)
		if result.Error != nil {
			configs.Sugar.Error("DocUpd", zap.String("username", username), zap.Error(err))
			c.JSON(http.StatusBadGateway, gin.H{
				"code": http.StatusBadGateway,
				"msg":  "文件保存到数据库失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "文章创建成功",
			"article": gin.H{
				"username":   username,
				"id":         article.ID,
				"title":      article.Title,
				"url":        article.URL,
				"summary":    article.Summary,
				"created_at": article.CreatedAt,
			},
		})
		content, _ := json.Marshal(article)
		go Follow.AddFeedToFollower(username, string(content))
	} else {
		question := &Document.Question{
			Username: username,
			Title:    file.Filename,
			URL:      "/browse/question/" + filename,
		}
		result := configs.Db.Create(question)
		if result.Error != nil {
			configs.Sugar.Error("DocUpd", zap.String("username", username), zap.Error(err))
			c.JSON(http.StatusBadGateway, gin.H{
				"code": http.StatusBadGateway,
				"msg":  "文件保存到数据库失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "问题创建成功",
			"question": gin.H{
				"username":   username,
				"id":         question.ID,
				"title":      question.Title,
				"url":        question.URL,
				"created_at": question.CreatedAt,
			},
		})
		content, _ := json.Marshal(question)
		go Follow.AddFeedToFollower(username, string(content))
	}
}
