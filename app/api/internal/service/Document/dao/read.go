package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	// 1. 尝试从 Redis 获取缓存
	ctx := context.Background() // 这里需要注意使用无超时限制的上下文，因为后面需要异步处理。
	cacheKey := fmt.Sprintf("user:%s:%s", username, filetype)
	cacheData, err := configs.Cli.Get(ctx, cacheKey).Result()
	if err == nil {
		var response []gin.H
		if json.Unmarshal([]byte(cacheData), &response) == nil {
			configs.Logger.Info("GetUserInfo from cache", zap.String("username", username), zap.String("filetype", filetype))
			c.JSON(http.StatusOK, gin.H{
				"total": len(response),
				"data":  response,
			})
			return
		}
		configs.Logger.Warn("Failed to parse cache data", zap.String("username", username), zap.String("filetype", filetype))
	}
	// 2. 缓存未命中，从数据库查询
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
	// 3. 将结构写入缓存，异步处理，不阻塞响应
	go func() {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			configs.Logger.Error("Marshal response", zap.Error(err))
			return
		}
		err = configs.Cli.Set(ctx, cacheKey, responseJSON, 10*time.Minute).Err()
		if err != nil {
			configs.Logger.Error("Set response", zap.Error(err))
		} else {
			configs.Logger.Info("Set response", zap.String("cacheKey", cacheKey))
		}
	}()
	configs.Logger.Info("GetUserInfo from database", zap.String("username", username))
	c.JSON(http.StatusOK, gin.H{
		"total": len(file),
		"data":  response,
	})
}
