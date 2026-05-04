package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("secret")

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Забираем заголовок
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing Authorization header",
			})
			return
		}

		// 2. Убираем "Bearer "
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Парсим токен
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// 4. Достаём claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid claims",
			})
			return
		}

		// 5. Достаём user_id
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id",
			})
			return
		}

		// 6. Достаём role
		role, _ := claims["role"].(string)

		// 7. Кладём в контекст
		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}
