package dao

import (
	"io"
	"net/http"
	"os"
	"zhihu/app/api/configs"
	"zhihu/utils/files"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetMdFile(c *gin.Context) {
	filetype := c.Param("filetype")
	thePath := c.Param("url")
	if filetype == "article" {
		thePath = "Storage/Document/Article/" + thePath
	} else if filetype == "question" {
		thePath = "Storage/Document/Question/" + thePath
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
		})
		return
	}
	fileInfo, err := os.Stat(thePath)
	if os.IsNotExist(err) {
		configs.Logger.Error("GetMdFile", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 验证文件是否存在

	file, err := os.Open(thePath)
	if err != nil {
		configs.Logger.Error("GetMdFile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "文件打开失败",
		})
	}
	defer func() {
		if err := file.Close(); err != nil {
			configs.Logger.Error("GetMdFile", zap.Error(err))
		}
	}()

	if !files.IsMarkdown(file.Name()) {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "非法的文件类型",
		})
		return
	}
	contentType := "text/markdown; charset=utf-8"
	files.HeaderSet(c, fileInfo, contentType)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		configs.Logger.Error("传输文件失败", zap.Error(err))
		return
	}
}
