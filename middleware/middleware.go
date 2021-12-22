package middleware

import "github.com/syahidfrd/go-boilerplate/utils/logger"

// MiddlewareManager ...
type MiddlewareManager struct {
	logger logger.Logger
}

// NewMiddlewareManager will create new an MiddlewareManager object
func NewMiddlewareManager(logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		logger: logger,
	}
}
