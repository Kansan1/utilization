package middleware

import (
	"net/http"
	"utilization-backend/src/api/result"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否登录
		token := c.GetHeader("Authorization")
		// 如果用户未登录，返回未授权的响应
		if token == "" {
			result.Fail(
				c,
				http.StatusUnauthorized,
				"未登录")
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
