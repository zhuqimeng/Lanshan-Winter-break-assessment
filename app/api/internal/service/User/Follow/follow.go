package Follow

import (
	"errors"
	"net/http"
	"time"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func checkFollow(username, follower string) (error, bool) {
	var count int64
	if err := configs.Db.Model(&User.Relation{}).Where("username = ? AND follower = ?", username, follower).Count(&count).Error; err != nil {
		return err, false
	}
	return nil, count > 0
}

func checkUser(username string) (error, bool) {
	var count int64
	if err := configs.Db.Model(&User.User{}).Where("name = ?", username).Count(&count).Error; err != nil {
		return err, false
	}
	return nil, count > 0
}

func OnFollow(c *gin.Context) {
	follower := c.GetString("username")
	username := c.Param("username")
	if username == follower {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "不能关注自己",
		})
		return
	}
	err, ok := checkUser(username)
	if err != nil {
		configs.Logger.Error("OnFollow err", zap.String("username", username), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "该用户不存在或已被封禁。",
		})
		return
	}
	err, ok = checkFollow(username, follower)
	if err != nil {
		configs.Logger.Error("OnFollow CheckFollow err", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
		return
	}
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "已经关注过了哦！",
		})
		return
	}
	var newRelation User.Relation
	if err = configs.Db.Unscoped().Where("username = ? AND follower = ?", username, follower).First(&newRelation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRelation = User.Relation{
				Follower: follower,
				Username: username,
			}
			if err = configs.Db.Create(&newRelation).Error; err != nil {
				configs.Logger.Error("OnFollow Create err", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
				})
				return
			}
			configs.Logger.Info("OnFollow Create success", zap.String("username", username))
			c.JSON(http.StatusOK, gin.H{
				"msg": "关注成功。",
			})
			go pushFeedAToB(username, follower)
			return
			// 首次关注
		}
		configs.Logger.Error("OnFollow", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
		return
	}
	if err = configs.Db.Unscoped().Model(&User.Relation{}).Where("username = ? AND follower = ?", username, follower).Update("deleted_at", nil).Error; err != nil {
		configs.Logger.Error("OnFollow Update err", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "关注成功。",
	})
	// 再次关注
}

func OffFollow(c *gin.Context) {
	follower := c.GetString("username")
	username := c.Param("username")
	if username == follower {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的操作。",
		})
		return
	}
	err, ok := checkUser(username)
	if err != nil {
		configs.Logger.Error("OffFollow err", zap.String("username", username), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "该用户不存在或已被封禁。",
		})
		return
	}
	err, ok = checkFollow(username, follower)
	if err != nil {
		configs.Logger.Error("OffFollow CheckFollow err", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的操作。",
		})
		return
	}
	if err = configs.Db.Model(&User.Relation{}).Where("username = ? AND follower = ?", username, follower).Update("deleted_at", time.Now()).Error; err != nil {
		configs.Logger.Error("OffFollow Update err", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "操作成功。",
	})
	// 软删除
}
