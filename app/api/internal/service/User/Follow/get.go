package Follow

import (
	"net/http"
	"zhihu/app/api/configs"
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

func GetFollowerStrings(username string) (error, []string) {
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
	err, followerIDs := GetFollowerStrings(username)
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
