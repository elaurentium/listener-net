package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthService interface {
	GenerateToken(userID uuid.UUID, dispositive string) (string, error)
}

type authService struct {
	secretKey []byte
}

func NewAuthService() AuthService {
	return &authService{
		secretKey: []byte("1123213123123-241241241241"),
	}
}

func (a *authService) GenerateToken(userID uuid.UUID, dispositive string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":          userID.String(),
		"user_dispositive": dispositive,
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}
