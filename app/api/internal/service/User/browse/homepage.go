package browse

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetHome(c *gin.Context) {
	username := c.Param("username")
	var count int64
	if err := configs.Db.Model(&User.User{}).Where("name = ?", username).Count(&count).Error; err != nil {
		configs.Logger.Error("GetHome", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadGateway, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "找不到该用户",
		})
		return
	}
	// 验证用户是否存在

	var user User.User
	res := configs.Db.Where("name = ?", username).First(&user)
	if res.Error != nil {
		configs.Logger.Error("GetHome", zap.Error(res.Error))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  res.Error.Error(),
		})
		return
	}
	thePath := user.AvatarURL
	fileInfo, err := os.Stat(thePath)
	if thePath == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户还未上传头像",
		})
		return
	}
	if _, err := os.Stat(thePath); os.IsNotExist(err) {
		configs.Logger.Error("GetHome", zap.Error(res.Error))
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 验证头像是否存在

	file, err := os.Open(thePath)
	if err != nil {
		configs.Logger.Error("GetHome", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "文件打开失败",
		})
	}
	defer func() {
		if err := file.Close(); err != nil {
			configs.Logger.Error("GetHome", zap.Error(err))
		}
	}()
	// 打开头像文件

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		configs.Logger.Error("读取文件类型失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "无法读取文件",
		})
		return
	}
	_, _ = file.Seek(0, 0)
	contentType := http.DetectContentType(buffer[:n])
	// 检测头像类型

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.Header("Cache-Control", "public, max-age=3600")
	c.Header("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))
	// 设置请求头

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		configs.Logger.Error("传输图片失败", zap.Error(err))
		return
	}
}
