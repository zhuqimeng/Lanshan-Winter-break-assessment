package Answer

import (
	"net/http"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/Document"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func checkQuestion(c *gin.Context, link string) bool {
	var count int64
	if err := configs.Db.Model(&Document.Question{}).Where("url = ?", link).Count(&count).Error; err != nil {
		configs.Logger.Error("AnsRead", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器故障",
		})
		return false
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "该问题不存在或已被删除",
		})
		return false
	}
	return true
}

func Read(c *gin.Context) {
	if filetype := c.Param("filetype"); filetype != "question" {
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "只能看到问题的回答",
		})
		return
	}
	link := "/browse/question/" + c.Param("url")
	if !checkQuestion(c, link) {
		return
	}
	var answers []Document.Answer
	if err := configs.Db.Model(&Document.Answer{}).Where("link = ?", link).Order("created_at DESC").Find(&answers).Error; err != nil {
		configs.Logger.Error("AnsRead", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "服务器故障",
		})
		return
	}
	if len(answers) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "该问题还没有回答，快来做第一个吃螃蟹的人吧~~",
		})
		return
	}
	var result []gin.H
	for _, answer := range answers {
		result = append(result, gin.H{
			"url":       answer.URL,
			"username":  answer.Username,
			"createdAt": answer.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"total":   len(answers),
		"answers": result,
	})
}
