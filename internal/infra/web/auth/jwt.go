package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)


type JWTConfig struct {
	Secret 			string
	AcessTokenExp 	time.Duration
	RefreshTokenExp time.Duration
}

type JWTService struct {
	config JWTConfig
}


func NewJWTService(config JWTConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}


type Claims struct {
	UserID 		uuid.UUID `json:"user_id"`
	Dispositive string    `json:"user_dispositive"`
	jwt.RegisteredClaims
}


func (s *JWTService) GenerateToken(userID uuid.UUID, dispositive string) (string, error) {
	expirationTime := time.Now().Add(s.config.AcessTokenExp)
	claims := &Claims{
		UserID: userID,
		Dispositive: dispositive,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.config.Secret))
}