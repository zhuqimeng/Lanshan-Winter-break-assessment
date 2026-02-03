package Comment

import (
	"fmt"
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func checkLink(c *gin.Context, link, filetype string) bool {
	var count int64
	switch filetype {
	case "article":
		if err := configs.Db.Model(&Document.Article{}).Where("url = ?", link).Count(&count).Error; err != nil {
			configs.Logger.Error("checkLink", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			return false
		}
	case "question":
		if err := configs.Db.Model(&Document.Question{}).Where("url = ?", link).Count(&count).Error; err != nil {
			configs.Logger.Error("checkLink", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			return false
		}
	case "answer":
		if err := configs.Db.Model(&Document.Answer{}).Where("url = ?", link).Count(&count).Error; err != nil {
			configs.Logger.Error("checkLink", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			return false
		}
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
		return false
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
		return false
	}
	return true
}

func Read(c *gin.Context) {
	filetype := c.Param("filetype")
	link := "/browse/" + filetype + "/" + c.Param("url")
	var comments []Document.Comment
	var result []gin.H
	if !checkLink(c, link, filetype) {
		return
	}
	if err := configs.Db.Model(&Document.Comment{}).Where("link = ?", link).Order("created_at DESC").Find(&comments).Error; err != nil {
		configs.Logger.Error("CommentRead", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	if len(comments) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "这里还没有评论，期待你的足迹~~",
		})
		return
	}
	for _, comment := range comments {
		result = append(result, gin.H{
			"username":  comment.Username,
			"content":   comment.Content,
			"createdAt": comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		fmt.Println(comment.Content)
	}
	c.JSON(http.StatusOK, gin.H{
		"total":    len(comments),
		"comments": result,
	})
}
