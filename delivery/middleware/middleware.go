package middleware

import "github.com/syahidfrd/go-boilerplate/utils/logger"

// Middleware ...
type Middleware struct {
	logger logger.Logger
}

// NewMiddleware will create new an Middleware object
func NewMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}
