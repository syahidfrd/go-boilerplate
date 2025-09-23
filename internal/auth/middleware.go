package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/jwt"
	"github.com/syahidfrd/go-boilerplate/internal/pkg/render"
)

type contextKey string

// UserIDKey is the context key used to store user ID in request context
const UserIDKey contextKey = "user_id"

// JWTMiddleware provides JWT authentication middleware functionality
type JWTMiddleware struct {
	jwtService *jwt.Service
}

// NewJWTMiddleware creates a new JWT middleware with the provided JWT service
func NewJWTMiddleware(jwtService *jwt.Service) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

// Authenticate validates JWT tokens and adds user ID to request context
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Ctx(ctx).Warn().Msg("missing authorization header")
			render.JSON(w, http.StatusUnauthorized, map[string]string{"message": "authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Ctx(ctx).Warn().Msg("invalid authorization header format")
			render.JSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid authorization header format"})
			return
		}

		token := parts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			log.Ctx(ctx).Warn().Msgf("invalid token: %s", err.Error())
			render.JSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid token"})
			return
		}

		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
