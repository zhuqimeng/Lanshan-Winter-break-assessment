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
			"msg":  "该用户不存在或被封禁",
		})
		return false
	}
	return true
}

func GetUserInfo(c *gin.Context) {
	username := c.Param("username")
	if !CheckUser(c, username) {
		return
	}
	filetype := c.GetString("filetype")
	var response []gin.H
	switch filetype {
	case "article":
		var articles []Document.Article
		result := configs.Db.Model(&Document.Article{}).Where("username = ?", username).Order("created_at DESC").Find(&articles)
		if result.Error != nil {
			configs.Logger.Error("GetUserInfo", zap.Error(result.Error))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  result.Error.Error(),
			})
			return
		}
		if len(articles) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "该用户还未发布过文章",
			})
			return
		}
		for _, article := range articles {
			response = append(response, gin.H{
				"id":        article.ID,
				"title":     article.Title,
				"url":       article.URL,
				"createdAt": article.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		configs.Logger.Info("GetUserInfo", zap.String("username", username))
		c.JSON(http.StatusOK, gin.H{
			"total":    len(articles),
			"articles": response,
		})
	case "question":
		var questions []Document.Question
		result := configs.Db.Model(&Document.Question{}).Where("username = ?", username).Order("created_at DESC").Find(&questions)
		if result.Error != nil {
			configs.Logger.Error("GetUserInfo", zap.Error(result.Error))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  result.Error.Error(),
			})
			return
		}
		if len(questions) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "该用户还未提过问题",
			})
			return
		}
		for _, question := range questions {
			response = append(response, gin.H{
				"id":        question.ID,
				"title":     question.Title,
				"url":       question.URL,
				"createdAt": question.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		configs.Logger.Info("GetUserInfo", zap.String("username", username))
		c.JSON(http.StatusOK, gin.H{
			"total":     len(questions),
			"questions": response,
		})
	case "answer":
		var answers []Document.Answer
		result := configs.Db.Model(&Document.Answer{}).Where("username = ?", username).Order("created_at DESC").Find(&answers)
		if result.Error != nil {
			configs.Logger.Error("GetUserInfo", zap.Error(result.Error))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  result.Error.Error(),
			})
			return
		}
		if len(answers) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "该用户还未做过回答",
			})
			return
		}
		for _, answer := range answers {
			response = append(response, gin.H{
				"id":           answer.ID,
				"question_url": answer.Link,
				"answer_url":   answer.URL,
				"createdAt":    answer.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		configs.Logger.Info("GetUserInfo", zap.String("username", username))
		c.JSON(http.StatusOK, gin.H{
			"total":   len(answers),
			"answers": response,
		})
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "不存在的页面",
		})
	}
}
