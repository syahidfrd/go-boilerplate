package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Service provides JWT token generation and validation functionality
type Service struct {
	secretKey []byte
}

// Claims represents JWT token claims with user ID
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// NewService creates a new JWT service with the provided secret key
func NewService(secretKey string) *Service {
	return &Service{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken creates a JWT token for the given user ID with 24-hour expiration
func (s *Service) GenerateToken(userID int64) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT token, returning the claims if valid
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	return claims, nil
}
