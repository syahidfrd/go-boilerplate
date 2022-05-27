package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/syahidfrd/go-boilerplate/entity"
)

// GenerateCID will search for a correlation header and set a request-level
// correlation id into the context. If no header is found, a new UUID will be generated.
func (m *MiddlewareManager) GenerateCID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cid := c.Request().Header.Get(entity.HeaderXCorrelationID)
			if cid == "" {
				cid = uuid.NewV4().String()
			}

			newReq := c.Request().WithContext(context.WithValue(c.Request().Context(), entity.CorrelationContextKey, cid))
			c.SetRequest(newReq)
			c.Request().Header.Set(entity.HeaderXCorrelationID, cid)
			return next(c)
		}
	}
}
