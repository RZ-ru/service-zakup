package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 15 * time.Minute

type AuthService struct {
	secret []byte
}

func NewAuthService() (*AuthService, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	return &AuthService{
		secret: []byte(secret),
	}, nil
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id required")
	}

	role := "user"
	if userID == "admin" {
		role = "admin"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	})

	return token.SignedString(s.secret)
}
