package middleware

import "github.com/gin-gonic/gin"

func UserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")

		if userID == "" {
			userID = "anonymous" // для теста
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
