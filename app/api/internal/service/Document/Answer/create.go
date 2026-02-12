package Answer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"
	"zhihu/app/api/internal/service/User/Follow"
	"zhihu/utils/files"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Create(c *gin.Context) {
	username := c.GetString("username")
	link := "/browse/question/" + c.Param("url")
	if filetype := c.Param("filetype"); filetype != "question" {
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "只能对问题回答",
		})
		return
	}
	if !checkQuestion(c, link) {
		return
	}
	file, err := c.FormFile("answer")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
		})
		return
	}
	if file.Size > 5<<20 {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "文件太大，最大支持 5MB",
		})
		return
	}
	if !files.IsMarkdown(file.Filename) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "不支持的文件类型",
		})
		return
	}
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d-%s%s", timestamp, username, ext)
	thePath := fmt.Sprintf("Storage/Document/Answer/%s", filename)
	if err := c.SaveUploadedFile(file, thePath); err != nil {
		configs.Logger.Error("AnsUpd", zap.String("username", username), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "文件保存失败",
		})
		return
	}
	answer := &Document.Answer{
		Username: username,
		Link:     link,
		URL:      "/browse/answer/" + filename,
	}
	if err = configs.Db.Create(answer).Error; err != nil {
		configs.Logger.Error("AnsUpd", zap.String("username", username), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "文件保存到数据库失败",
		})
		return
	}
	configs.Logger.Info("AnsUpd", zap.String("username", username))
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": answer,
	})
	content, _ := json.Marshal(answer)
	go Follow.AddFeedToFollower(username, string(content))
}
