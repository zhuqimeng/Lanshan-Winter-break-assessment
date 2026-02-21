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
	"zhihu/utils/randoms"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Infos interface {
	Print() gin.H
}

func checkUser(c *gin.Context, username string) bool {
	if exists, err := configs.UserBf.Exists(context.Background(), username); err != nil {
		configs.Logger.Error("user bloom err", zap.Error(err))
		// 布隆过滤器故障时直接用数据库查询
	} else if !exists {
		return false
	}
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
	if !checkUser(c, username) {
		return
	}
	filetype := c.Query("filetype")
	// 1. 尝试从 Redis 获取缓存
	ctx := context.Background() // 这里需要注意使用无超时限制的上下文，因为后面需要异步处理。
	cacheKey := fmt.Sprintf("user:%s:%s", username, filetype)
	singleFlightKey := fmt.Sprintf("sf:%s", cacheKey)
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
	// 使用 singleFlight 防止缓存击穿
	result, err, shared := configs.RequestGroup.Do(singleFlightKey, func() (interface{}, error) {
		var file []Infos
		switch filetype {
		case "article":
			var articles []Document.Article
			result := configs.Db.Model(&Document.Article{}).Where("username = ?", username).Order("created_at DESC").Find(&articles)
			if result.Error != nil {
				return nil, result.Error
			}
			file = make([]Infos, len(articles))
			for i, article := range articles {
				file[i] = &article
			}
		case "question":
			var questions []Document.Question
			result := configs.Db.Model(&Document.Question{}).Where("username = ?", username).Order("created_at DESC").Find(&questions)
			if result.Error != nil {
				return nil, result.Error
			}
			file = make([]Infos, len(questions))
			for i, question := range questions {
				file[i] = &question
			}
		case "answer":
			var answers []Document.Answer
			result := configs.Db.Model(&Document.Answer{}).Where("username = ?", username).Order("created_at DESC").Find(&answers)
			if result.Error != nil {
				return nil, result.Error
			}
			file = make([]Infos, len(answers))
			for i, answer := range answers {
				file[i] = &answer
			}
		case "comment":
			var comments []Document.Comment
			result := configs.Db.Model(&Document.Comment{}).Where("username = ?", username).Order("created_at DESC").Find(&comments)
			if result.Error != nil {
				return nil, result.Error
			}
			file = make([]Infos, len(comments))
			for i, comment := range comments {
				file[i] = &comment
			}
		default:
			return nil, fmt.Errorf("不合法的文件类型")
		}
		if len(file) == 0 {
			// 返回空列表，不报错
			return []gin.H{}, nil
		}
		var response []gin.H
		for _, v := range file {
			response = append(response, v.Print())
		}
		return response, nil
	})
	if err != nil {
		if err.Error() == "不合法的文件类型" {
			c.JSON(http.StatusNotFound, gin.H{
				"code": http.StatusNotFound,
				"msg":  err.Error(),
			})
		} else {
			configs.Logger.Error("GetUserInfo", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
			})
		}
		return
	}
	response, ok := result.([]gin.H)
	if !ok {
		configs.Logger.Error("GetUserInfo from singleFlight", zap.String("username", username))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "数据处理错误",
		})
		return
	}
	if shared {
		configs.Logger.Info("GetUserInfo from shared singleFlight", zap.String("username", username), zap.String("filetype", filetype))
	}
	// 3. 将结构写入缓存，异步处理，不阻塞响应
	go func() {
		responseJSON, err := json.Marshal(response)
		if err != nil {
			configs.Logger.Error("Marshal response", zap.Error(err))
			return
		}
		err = configs.Cli.Set(ctx, cacheKey, responseJSON, time.Duration(10+randoms.GetRandomNumber(1, 5))*time.Minute).Err()
		if err != nil {
			configs.Logger.Error("Set response", zap.Error(err))
		} else {
			configs.Logger.Info("Set response", zap.String("cacheKey", cacheKey))
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"total": len(response),
		"data":  response,
	})
}
