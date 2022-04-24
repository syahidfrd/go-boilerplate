package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/labstack/echo"
	"github.com/syahidfrd/go-boilerplate/utils"
)

func (m *MiddlewareManager) InboundLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		start := time.Now()
		req := c.Request()
		res := c.Response()

		// Request body
		reqBody := []byte{}
		if c.Request().Body != nil {
			reqBody, _ = ioutil.ReadAll(c.Request().Body)
		}
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		if err = next(c); err != nil {
			c.Error(err)
		}

		// Latency in ms
		latency := float64(time.Since(start).Nanoseconds()/1e4) / 100.0

		m.logger.Infow("INBOUND LOG",
			"cid", utils.GetCID(req.Context()),
			"ip", req.RemoteAddr,
			"method", req.Method,
			"user_agent", req.UserAgent(),
			"path", req.URL.Path,
			"body", utils.CompactJSON(reqBody),
			"status", res.Status,
			"latency", latency,
		)

		return
	}
}
