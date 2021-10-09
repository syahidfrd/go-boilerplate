package middleware

import (
	"context"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type ContextKey string

const (
	HeaderXCorrelationID  string     = "X-Correlation-ID"
	CorrelationContextKey ContextKey = "cid"
)

// GenerateCorrelationID will search for a correlation header and set a request-level
// correlation id into the context. If no header is found, a new UUID will be generated.
func GenerateCorrelationID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := c.Request().Header.Get(HeaderXCorrelationID)
			if correlationID == "" {
				correlationID = uuid.NewV4().String()
			}

			c.Request().Header.Set(HeaderXCorrelationID, correlationID)
			newReq := c.Request().WithContext(context.WithValue(c.Request().Context(), CorrelationContextKey, correlationID))
			c.SetRequest(newReq)
			return next(c)
		}
	}
}
