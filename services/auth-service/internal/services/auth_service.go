package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	secret []byte
}

func NewAuthService() *AuthService {
	return &AuthService{
		secret: []byte("secret"),
	}
}

func (s *AuthService) GenerateToken(userID string) (string, error) {

	role := "user"
	if userID == "admin" {
		role = "admin"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString(s.secret)
}
