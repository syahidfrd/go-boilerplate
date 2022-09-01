package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: secretKey,
		issuer:    "go-boilerplate",
	}
}

func (s *jwtService) GenerateToken(ctx context.Context, userID int64) (token string, err error) {
	claims := &jwtCustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = t.SignedString([]byte(s.secretKey))
	return
}

func (s *jwtService) ValidateToken(ctx context.Context, tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}
