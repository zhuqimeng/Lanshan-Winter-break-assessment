package Follow

import (
	"encoding/json"
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetFollowing(c *gin.Context) {
	username := c.Param("username")
	var (
		followings   []User.Relation
		followingIDs []string
	)
	if err := configs.Db.Model(&User.Relation{}).Where("follower = ?", username).Find(&followings).Error; err != nil {
		configs.Logger.Error("GetFollowing FindFollowing err", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
		return
	}
	for _, v := range followings {
		followingIDs = append(followingIDs, v.Username)
	}
	c.JSON(http.StatusOK, gin.H{
		"data": followingIDs,
	})
}

func getFollowerStrings(username string) (error, []string) {
	var (
		followers []User.Relation
		result    []string
	)
	if err := configs.Db.Model(&User.Relation{}).Where("username = ?", username).Find(&followers).Error; err != nil {
		return err, []string{}
	}
	for _, v := range followers {
		result = append(result, v.Follower)
	}
	return nil, result
}

func GetFollower(c *gin.Context) {
	username := c.Param("username")
	err, followerIDs := getFollowerStrings(username)
	if err != nil {
		configs.Logger.Error("GetFollower GetFollower err", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": followerIDs,
	})
}

func getPost(username string) (error, []string) {
	var (
		articles  []Document.Article
		questions []Document.Question
		answers   []Document.Answer
		result    []string
	)
	err := configs.Db.Where("username = ?", username).Order("created_at DESC").Limit(20).Find(&articles).Error
	if err != nil {
		return err, []string{}
	}
	for _, v := range articles {
		term, _ := json.Marshal(v)
		result = append(result, string(term))
	}
	err = configs.Db.Where("username = ?", username).Order("created_at DESC").Limit(20).Find(&questions).Error
	if err != nil {
		return err, []string{}
	}
	for _, v := range questions {
		term, _ := json.Marshal(v)
		result = append(result, string(term))
	}
	err = configs.Db.Where("username = ?", username).Order("created_at DESC").Limit(20).Find(&answers).Error
	if err != nil {
		return err, []string{}
	}
	for _, v := range answers {
		term, _ := json.Marshal(v)
		result = append(result, string(term))
	}
	return nil, result
}
