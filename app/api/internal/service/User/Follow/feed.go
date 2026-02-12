package Follow

import (
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func pushFeedAToB(A, B string) {
	err, content := getPost(A)
	if err != nil {
		configs.Logger.Error("push feed err", zap.String("A", A), zap.String("B", B), zap.Error(err))
		return
	}
	for _, v := range content {
		term := User.FeedItem{
			Username: B,
			Content:  v,
			Author:   A,
		}
		if err = configs.Db.Create(&term).Error; err != nil {
			configs.Logger.Error("push feed err", zap.String("A", A), zap.String("B", B), zap.Error(err))
			return
		}
	}
	configs.Logger.Info("push feed success", zap.String("A", A), zap.String("B", B))
}

func ShowFeeds(c *gin.Context) {
	username := c.GetString("username")
	var (
		response []gin.H
		feeds    []User.FeedItem
	)
	if err := configs.Db.Where("username = ?", username).Order("created_at DESC").Limit(1000).Find(&feeds).Error; err != nil {
		configs.Logger.Error("show feed err", zap.String("username", username), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
		return
	}
	for _, v := range feeds {
		response = append(response, gin.H{
			"content": v.Content,
			"author":  v.Author,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}

func AddFeedToFollower(username, content string) {
	err, followers := getFollowerStrings(username)
	if err != nil {
		configs.Logger.Error("addFeedToFollower err", zap.String("username", username), zap.Error(err))
		return
	}
	for _, v := range followers {
		term := User.FeedItem{
			Username: v,
			Content:  content,
			Author:   username,
		}
		if err = configs.Db.Create(&term).Error; err != nil {
			configs.Logger.Error("addFeedToFollower err", zap.String("username", username), zap.Error(err))
			return
		}
	}
	configs.Logger.Info("addFeedToFollower success", zap.String("username", username))
}
