package dao

import (
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckUser(c *gin.Context, username string) bool {
	var count int64
	if err := configs.Db.Model(&User.User{}).Where("name = ?", username).Count(&count).Error; err != nil {
		configs.Logger.Error("CheckUser", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return false
	}
	if count == 0 {
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "找不到该用户",
		})
		return false
	}
	return true
}

func GetUserArticles(c *gin.Context) {
	username := c.Param("username")
	if ok := CheckUser(c, username); !ok {
		return
	}
	var articles []Document.Article
	result := configs.Db.Model(&Document.Article{}).Where("username = ?", username).Order("created_at DESC").Find(&articles)
	if result.Error != nil {
		configs.Logger.Error("GetArticles", zap.Error(result.Error))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  result.Error.Error(),
		})
		return
	}

	var response []gin.H
	for _, article := range articles {
		response = append(response, gin.H{
			"id":        article.ID,
			"title":     article.Title,
			"url":       article.URL,
			"createdAt": article.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	configs.Logger.Info("GetArticles", zap.String("username", username))
	c.JSON(http.StatusOK, gin.H{
		"total":    len(articles),
		"articles": response,
	})
}

func GetUserQuestions(c *gin.Context) {
	username := c.Param("username")
	if ok := CheckUser(c, username); !ok {
		return
	}
	var questions []Document.Question
	result := configs.Db.Model(&Document.Question{}).Where("username = ?", username).Order("created_at DESC").Find(&questions)
	if result.Error != nil {
		configs.Logger.Error("GetQuestions", zap.Error(result.Error))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  result.Error.Error(),
		})
		return
	}

	var response []gin.H
	for _, question := range questions {
		response = append(response, gin.H{
			"id":        question.ID,
			"title":     question.Title,
			"url":       question.URL,
			"createdAt": question.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	configs.Logger.Info("GetQuestions", zap.String("username", username))
	c.JSON(http.StatusOK, gin.H{
		"total":     len(questions),
		"questions": response,
	})
}
