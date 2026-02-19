package Comment

import (
	"context"
	"fmt"
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateCommentReq struct {
	Comment string `json:"comment" binding:"required"`
}

func Create(c *gin.Context) {
	username := c.GetString("username")
	filetype := c.Param("filetype")
	link := "/browse/" + filetype + "/" + c.Param("url")
	if !checkLink(c, link, filetype) {
		return
	}
	var req CreateCommentReq
	if err := c.ShouldBind(&req); err != nil {
		configs.Logger.Error("CreateComment", zap.Error(err), zap.String("username", username))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment := &Document.Comment{
		Content:  req.Comment,
		Username: username,
		Link:     link,
	}
	if err := configs.Db.Create(comment).Error; err != nil {
		configs.Logger.Error("CreateComment", zap.Error(err), zap.String("username", username))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	configs.Logger.Info("CreateComment", zap.String("username", username))
	c.JSON(http.StatusOK, gin.H{"comment": comment})
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s:comment", username)
	// 异步删除缓存
	go func() {
		err := configs.Cli.Del(ctx, cacheKey).Err()
		if err != nil {
			configs.Logger.Error("Failed to delete cache", zap.Any("username", username), zap.Error(err))
		} else {
			configs.Logger.Info("Deleted cache", zap.Any("username", username))
		}
	}()
}
