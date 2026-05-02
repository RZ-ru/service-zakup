// internal/middleware/user_context.go
package middleware

import "github.com/gin-gonic/gin"

func UserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")

		if userID == "" {
			userID = "anonymous" // временно
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
