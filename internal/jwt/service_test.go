package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewService(secretKey)

	assert.NotNil(t, service)
	assert.Equal(t, []byte(secretKey), service.secretKey)
}

func TestService_GenerateToken_Success(t *testing.T) {
	service := NewService("test-secret-key")
	userID := int64(123)

	token, err := service.GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Token should have 3 parts separated by dots
	parts := strings.Split(token, ".")
	assert.Equal(t, 3, len(parts))
}

func TestService_GenerateToken_DifferentUsers(t *testing.T) {
	service := NewService("test-secret-key")

	token1, err1 := service.GenerateToken(1)
	token2, err2 := service.GenerateToken(2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, token1, token2)
}

func TestService_ValidateToken_Success(t *testing.T) {
	service := NewService("test-secret-key")
	userID := int64(123)

	// Generate token
	token, err := service.GenerateToken(userID)
	assert.NoError(t, err)

	// Validate token
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.True(t, claims.ExpiresAt > time.Now().Unix())
}

func TestService_ValidateToken_InvalidToken(t *testing.T) {
	service := NewService("test-secret-key")

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "malformed token",
			token: "invalid.token.format",
		},
		{
			name:  "random string",
			token: "completely-invalid-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestService_ValidateToken_WrongSecretKey(t *testing.T) {
	// Generate token with one secret
	service1 := NewService("secret-key-1")
	token, err := service1.GenerateToken(123)
	assert.NoError(t, err)

	// Try to validate with different secret
	service2 := NewService("secret-key-2")
	claims, err := service2.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestService_ValidateToken_ExpiredToken(t *testing.T) {
	service := NewService("test-secret-key")

	// Create token with custom expiration (already expired)
	claims := &Claims{
		UserID: 123,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
			IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(service.secretKey)
	assert.NoError(t, err)

	// Validate expired token
	validatedClaims, err := service.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestService_TokenRoundTrip(t *testing.T) {
	service := NewService("test-secret-key")
	originalUserID := int64(456)

	// Generate token
	token, err := service.GenerateToken(originalUserID)
	assert.NoError(t, err)

	// Validate token
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)

	// Verify user ID matches
	assert.Equal(t, originalUserID, claims.UserID)

	// Verify expiration is in the future
	assert.True(t, claims.ExpiresAt > time.Now().Unix())

	// Verify issued at is in the past
	assert.True(t, claims.IssuedAt <= time.Now().Unix())
}

func TestClaims_ExpirationTime(t *testing.T) {
	service := NewService("test-secret-key")

	beforeGeneration := time.Now()
	token, err := service.GenerateToken(123)
	afterGeneration := time.Now()

	assert.NoError(t, err)

	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)

	// Token should expire in approximately 24 hours
	expectedExpiration := beforeGeneration.Add(24 * time.Hour).Unix()
	actualExpiration := claims.ExpiresAt

	// Allow for a few seconds difference due to execution time
	assert.InDelta(t, expectedExpiration, actualExpiration, 10)

	// Token should be issued between before and after generation
	assert.True(t, claims.IssuedAt >= beforeGeneration.Unix())
	assert.True(t, claims.IssuedAt <= afterGeneration.Unix())
}
