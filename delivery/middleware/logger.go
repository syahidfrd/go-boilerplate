package middleware

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/syahidfrd/go-boilerplate/utils"
)

func (m *Middleware) Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			reqBody := []byte{}
			if req.Body != nil {
				reqBody, _ = ioutil.ReadAll(req.Body)
			}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

			if err := next(c); err != nil {
				c.Error(err)
			}

			m.logger.Infow("INBOUND LOG",
				"request_id", utils.GetReqID(req.Context()),
				"remote_ip", c.RealIP(),
				"host", req.Host,
				"uri", req.RequestURI,
				"method", req.Method,
				"user_agent", req.UserAgent(),
				"body", utils.CompactJSON(reqBody),
				"status", res.Status,
				"latency", float64(time.Since(start).Nanoseconds()/1e4)/100.0,
				"bytes_in", req.Header.Get(echo.HeaderContentLength),
				"bytes_out", strconv.FormatInt(res.Size, 10),
			)

			return nil
		}
	}
}
