package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/syahidfrd/go-boilerplate/entity"
)

// RequestID will search for a correlation header and set a request-level
// correlation id into the context. If no header is found, a new UUID will be generated.
func (m *Middleware) RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(entity.RequestIDHeader)
			if requestID == "" {
				requestID = uuid.NewV4().String()
			}

			ctx := c.Request().Context()
			newReq := c.Request().WithContext(context.WithValue(ctx, entity.RequestIDKey, requestID))
			c.SetRequest(newReq)
			c.Request().Header.Set(entity.RequestIDHeader, requestID)

			return next(c)
		}
	}
}
