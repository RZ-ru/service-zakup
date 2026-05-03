// internal/middleware/logger.go
package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		requestID := c.GetString("request_id")
		userID := c.GetString("user_id")

		fmt.Printf(
			"[%s] %s %s | user=%s | status=%d | %s\n",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			userID,
			c.Writer.Status(),
			duration,
		)
	}
}
