package Auth

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"zhihu/app/api/configs"

	"github.com/gin-gonic/gin"
)

func FrequencyChecker(RequestType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		ctx := c.Request.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		key := fmt.Sprintf("%s:limit:%s", RequestType, username)
		now := time.Now().UnixMilli()
		window := 3 * time.Minute
		windowMs := window.Milliseconds()
		limit := 15
		// 三分钟内只允许操作十五次

		luaScript := `
			local key = KEYS[1]
			local now = tonumber(ARGV[1])
			local window = tonumber(ARGV[2])
			local limit = tonumber(ARGV[3])
			
			-- 移除窗口外的记录
			redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
			
			-- 获取当前数量
			local count = redis.call('ZCARD', key)
			
			if count >= limit then
				return 0
			end
			
			-- 添加当前请求并设置过期时间
			redis.call('ZADD', key, now, now)
			redis.call('EXPIRE', key, window / 1000 + 60)
			return 1
		`
		// Lua脚本保证原子性

		result, err := configs.Cli.Eval(ctx, luaScript, []string{key}, now, windowMs, limit).Int()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		if result == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "过于频繁地操作，请稍后再试。",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
