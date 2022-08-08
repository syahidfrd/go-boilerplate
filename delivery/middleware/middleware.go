package middleware

import "github.com/syahidfrd/go-boilerplate/utils/logger"

// Middleware ...
type Middleware struct {
	logger logger.Logger
}

// NewMiddlewareManager will create new an MiddlewareManager object
func NewMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}
