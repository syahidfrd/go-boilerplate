package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syahidfrd/go-boilerplate/internal/jwt"
)

func TestNewJWTMiddleware(t *testing.T) {
	mockJWTService := new(MockJWTService)
	middleware := NewJWTMiddleware(mockJWTService)

	assert.NotNil(t, middleware)
	assert.Equal(t, mockJWTService, middleware.jwtService)
}

func TestJWTMiddleware_Authenticate_Success(t *testing.T) {
	mockJWTService := new(MockJWTService)
	middleware := NewJWTMiddleware(mockJWTService)

	// Mock valid token validation
	claims := &jwt.Claims{UserID: 123}
	mockJWTService.On("ValidateToken", "valid-token").Return(claims, nil)

	// Create test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify user ID is in context
		userID, ok := GetUserIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, int64(123), userID)
		w.WriteHeader(http.StatusOK)
	})

	// Create request with valid authorization header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()

	// Execute middleware
	handler := middleware.Authenticate(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockJWTService.AssertExpectations(t)
}

func TestJWTMiddleware_Authenticate_MissingAuthHeader(t *testing.T) {
	mockJWTService := new(MockJWTService)
	middleware := NewJWTMiddleware(mockJWTService)

	// Create test handler that should not be called
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// Create request without authorization header
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Execute middleware
	handler := middleware.Authenticate(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "authorization header required")
}

func TestJWTMiddleware_Authenticate_InvalidAuthHeaderFormat(t *testing.T) {
	mockJWTService := new(MockJWTService)
	middleware := NewJWTMiddleware(mockJWTService)

	// Create test handler that should not be called
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"Only Bearer", "Bearer"},
		{"Missing Bearer", "token123"},
		{"Wrong prefix", "Basic token123"},
		{"Too many parts", "Bearer token123 extra"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tc.header)
			rr := httptest.NewRecorder()

			handler := middleware.Authenticate(nextHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
			assert.Contains(t, rr.Body.String(), "invalid authorization header format")
		})
	}
}

func TestJWTMiddleware_Authenticate_InvalidToken(t *testing.T) {
	mockJWTService := new(MockJWTService)
	middleware := NewJWTMiddleware(mockJWTService)

	// Mock invalid token validation
	mockJWTService.On("ValidateToken", "invalid-token").Return(nil, assert.AnError)

	// Create test handler that should not be called
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// Create request with invalid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	// Execute middleware
	handler := middleware.Authenticate(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid token")
	mockJWTService.AssertExpectations(t)
}

func TestGetUserIDFromContext_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, int64(123))

	userID, ok := GetUserIDFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, int64(123), userID)
}

func TestGetUserIDFromContext_NotFound(t *testing.T) {
	ctx := context.Background()

	userID, ok := GetUserIDFromContext(ctx)

	assert.False(t, ok)
	assert.Equal(t, int64(0), userID)
}

func TestGetUserIDFromContext_WrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, "not-an-int")

	userID, ok := GetUserIDFromContext(ctx)

	assert.False(t, ok)
	assert.Equal(t, int64(0), userID)
}
