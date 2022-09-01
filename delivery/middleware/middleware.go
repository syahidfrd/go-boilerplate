package middleware

import (
	"github.com/syahidfrd/go-boilerplate/utils/jwt"
	"github.com/syahidfrd/go-boilerplate/utils/logger"
)

// Middleware ...
type Middleware struct {
	jwtSvc jwt.JWTService
	logger logger.Logger
}

// NewMiddleware will create new Middleware object
func NewMiddleware(jwtSvc jwt.JWTService, logger logger.Logger) *Middleware {
	return &Middleware{
		jwtSvc: jwtSvc,
		logger: logger,
	}
}
