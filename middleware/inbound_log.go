package middleware

import (
	"github.com/labstack/echo"
	"github.com/syahidfrd/go-boilerplate/utils"
)

func (m *MiddlewareManager) InboundLog(ctx echo.Context, reqBody, resBody []byte) {
	m.logger.Infow("INBOUND LOG",
		"cid", utils.GetCID(ctx.Request().Context()),
		"ip", ctx.Request().RemoteAddr,
		"method", ctx.Request().Method,
		"header", ctx.Request().Header,
		"path", ctx.Path(),
		"body", utils.CompactJSON(reqBody),
		"status", ctx.Response().Status,
		"response", utils.CompactJSON(resBody),
	)
}
