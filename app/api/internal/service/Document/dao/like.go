package dao

import (
	"errors"
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type likeFile interface {
	UpdLike(bool)
}

func updateLike(username, url, filetype string, status bool) error {
	var (
		likeReq User.LikeUrlUser
		file    likeFile
	)
	switch filetype {
	case "article":
		var term Document.Article
		if err := configs.Db.Model(&Document.Article{}).Where("url = ?", url).First(&term).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("该文件不存在或已被删除。")
			}
			configs.Logger.Error("updateLike", zap.Error(err))
			return err
		}
		file = &term
	case "question":
		var term Document.Question
		if err := configs.Db.Model(&Document.Question{}).Where("url = ?", url).First(&term).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("该文件不存在或已被删除。")
			}
			configs.Logger.Error("updateLike", zap.Error(err))
			return err
		}
		file = &term
	case "answer":
		var term Document.Answer
		if err := configs.Db.Model(&Document.Answer{}).Where("url = ?", url).First(&term).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("该文件不存在或已被删除。")
			}
			configs.Logger.Error("updateLike", zap.Error(err))
			return err
		}
		file = &term
	default:
		return errors.New("不存在的 URL")
	}
	if err := configs.Db.Model(&User.LikeUrlUser{}).Where("username = ? AND url = ?", username, url).First(&likeReq).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if status {
				likeReq = User.LikeUrlUser{
					Username: username,
					Url:      url,
					Status:   false,
				}
			} else {
				return errors.New("你还未点赞过，无法取消点赞。")
			}
		} else {
			configs.Logger.Error("updateLike", zap.Error(err))
			return err
		}
	}
	if status == likeReq.Status {
		return errors.New("重复的操作。")
	}
	likeReq.Status = status
	file.UpdLike(status)
	if err := configs.Db.Save(&likeReq).Error; err != nil {
		configs.Logger.Error("updateLike", zap.Error(err))
		return err
	}
	if err := configs.Db.Save(file).Error; err != nil {
		configs.Logger.Error("updateLike", zap.Error(err))
		return err
	}
	configs.Logger.Info("updateLike", zap.String("user", username), zap.String("url", url), zap.Bool("status", status))
	return nil
}

func ChangeLike(c *gin.Context) {
	username := c.GetString("username")
	fileType := c.Param("filetype")
	url := "/browse/" + fileType + "/" + c.Param("url")
	status := c.Query("status")
	if status == "" {
		status = "add"
	}
	var err error
	if status == "add" {
		err = updateLike(username, url, fileType, true)
	} else if status == "del" {
		err = updateLike(username, url, fileType, false)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
