package dao

import (
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Infos interface {
	Print() gin.H
}

func checkUser(c *gin.Context, username string) bool {
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

func foundError(c *gin.Context, err error) {
	configs.Logger.Error("GetUserInfo", zap.Error(err))
	c.JSON(http.StatusBadRequest, gin.H{
		"code": http.StatusBadRequest,
		"msg":  err.Error(),
	})
}

func GetUserInfo(c *gin.Context) {
	username := c.Param("username")
	if !checkUser(c, username) {
		return
	}
	filetype := c.Query("filetype")
	var file []Infos
	switch filetype {
	case "article":
		var articles []Document.Article
		result := configs.Db.Model(&Document.Article{}).Where("username = ?", username).Order("created_at DESC").Find(&articles)
		if result.Error != nil {
			foundError(c, result.Error)
			return
		}
		file = make([]Infos, len(articles))
		for i, article := range articles {
			file[i] = &article
		}
	case "question":
		var questions []Document.Question
		result := configs.Db.Model(&Document.Question{}).Where("username = ?", username).Order("created_at DESC").Find(&questions)
		if result.Error != nil {
			foundError(c, result.Error)
			return
		}
		file = make([]Infos, len(questions))
		for i, question := range questions {
			file[i] = &question
		}
	case "answer":
		var answers []Document.Answer
		result := configs.Db.Model(&Document.Answer{}).Where("username = ?", username).Order("created_at DESC").Find(&answers)
		if result.Error != nil {
			foundError(c, result.Error)
			return
		}
		file = make([]Infos, len(answers))
		for i, answer := range answers {
			file[i] = &answer
		}
	case "comment":
		var comments []Document.Comment
		result := configs.Db.Model(&Document.Comment{}).Where("username = ?", username).Order("created_at DESC").Find(&comments)
		if result.Error != nil {
			foundError(c, result.Error)
			return
		}
		file = make([]Infos, len(comments))
		for i, comment := range comments {
			file[i] = &comment
		}
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "不存在的页面",
		})
		return
	}
	if len(file) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "这里空空如也哦~~",
		})
		return
	}
	var response []gin.H
	for _, v := range file {
		response = append(response, v.Print())
	}
	configs.Logger.Info("GetUserInfo", zap.String("username", username))
	c.JSON(http.StatusOK, gin.H{
		"total": len(file),
		"data":  response,
	})
}
