package Auto

import "github.com/gin-gonic/gin"

func RouterSet(filetype string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("filetype", filetype)
	}
}
