package middleware

import (
	"context"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/syahidfrd/go-boilerplate/entity"
)

// GenerateCorrelationID will search for a correlation header and set a request-level
// correlation id into the context. If no header is found, a new UUID will be generated.
func (m *MiddlewareManager) GenerateCorrelationID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := c.Request().Header.Get(entity.HeaderXCorrelationID)
			if correlationID == "" {
				correlationID = uuid.NewV4().String()
			}

			c.Request().Header.Set(entity.HeaderXCorrelationID, correlationID)
			newReq := c.Request().WithContext(context.WithValue(c.Request().Context(), entity.CorrelationContextKey, correlationID))
			c.SetRequest(newReq)
			return next(c)
		}
	}
}
