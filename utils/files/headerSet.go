package files

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HeaderSet(c *gin.Context, fileInfo os.FileInfo, contentType string) {
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.Header("Cache-Control", "public, max-age=3600")
	c.Header("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))
}
