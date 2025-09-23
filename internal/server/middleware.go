package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Middleware represents a function that wraps an HTTP handler
type Middleware func(http.Handler) http.Handler

// chainMiddleware applies multiple middleware functions to an HTTP handler in sequence
func chainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}

// responseWriter wraps http.ResponseWriter to capture response status
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// wrapResponseWriter creates a new responseWriter wrapper
func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// Status returns the HTTP status code of the response
func (rw *responseWriter) Status() int {
	return rw.status
}

// WriteHeader captures the status code before writing the header
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// logSeverity determines the log level based on HTTP status code
func logSeverity(statusCode int) zerolog.Level {
	switch {
	case statusCode >= 500:
		return zerolog.ErrorLevel
	case statusCode >= 400:
		return zerolog.ErrorLevel
	case statusCode >= 300:
		return zerolog.WarnLevel
	case statusCode >= 200:
		return zerolog.InfoLevel
	default:
		return zerolog.DebugLevel
	}
}

// logFields represents the structured fields for HTTP request logging
type logFields struct {
	RemoteIP   string
	Host       string
	UserAgent  string
	Method     string
	Path       string
	Body       string
	StatusCode int
	Latency    float64
}

// MarshalZerologObject implements zerolog object marshaling for structured logging
func (l *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("remote_ip", l.RemoteIP).
		Str("host", l.Host).
		Str("user_agent", l.UserAgent).
		Str("method", l.Method).
		Str("path", l.Path).
		Str("body", l.Body).
		Int("status_code", l.StatusCode).
		Float64("latency", l.Latency)
}

// loggerMiddleware creates HTTP request logging middleware with optional filtering
func loggerMiddleware(filter func(w http.ResponseWriter, r *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check filter
			if filter != nil && filter(w, r) {
				next.ServeHTTP(w, r)
				return
			}

			// Start timer
			start := time.Now()

			// Read request body
			var buf []byte
			if r.Body != nil {
				buf, _ = io.ReadAll(r.Body)

				// Restore the io.ReadCloser to its original state
				r.Body = io.NopCloser(bytes.NewBuffer(buf))
			}

			// Wraps an http.ResponseWriter, returning a proxy that allows you to
			// hook into various parts of the response process.
			ww := wrapResponseWriter(w)
			next.ServeHTTP(ww, r)

			dur := float64(time.Since(start).Nanoseconds()/1e4) / 100.0

			// Create log fields
			fields := &logFields{
				RemoteIP:   r.RemoteAddr,
				Host:       r.Host,
				UserAgent:  r.UserAgent(),
				Method:     r.Method,
				Path:       r.URL.Path,
				Body:       formatReqBody(r, buf),
				StatusCode: ww.Status(),
				Latency:    dur,
			}

			sev := logSeverity(ww.Status())
			logEntry := log.Ctx(r.Context()).WithLevel(sev).EmbedObject(fields)
			logEntry.Msg("http request")
		})
	}
}

// formatReqBody formats request body for logging, compacting JSON if valid
func formatReqBody(r *http.Request, data []byte) string {
	// Skip logging body for sensitive endpoints
	if strings.Contains(r.URL.Path, "/auth/") {
		return "[REDACTED]"
	}

	var js map[string]any
	if json.Unmarshal(data, &js) != nil {
		return string(data)
	}

	result := new(bytes.Buffer)
	if err := json.Compact(result, data); err != nil {
		log.Ctx(r.Context()).Error().Msgf("error compacting body request json: %s", err.Error())
		return ""
	}

	return result.String()
}

// realIPMiddleware extracts real client IP from various headers
func realIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rip := realIP(r); rip != "" {
			r.RemoteAddr = rip
		}

		next.ServeHTTP(w, r)
	})
}

// realIP extracts the real client IP address from request headers
func realIP(r *http.Request) string {
	trueClientIP := http.CanonicalHeaderKey("True-Client-IP")
	xForwardedFor := http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP := http.CanonicalHeaderKey("X-Real-IP")

	var ip string
	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

// recoverMiddleware handles panics and logs them with stack traces
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				log.Ctx(r.Context()).
					Error().
					Err(err).
					Bytes("stack", debug.Stack()).
					Msg("panic recover")

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware generates or extracts request ID and adds it to context
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestIDHeader := "X-Request-Id"
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
			r.Header.Set(requestIDHeader, requestID)
		}

		w.Header().Set(requestIDHeader, requestID)

		ctx := log.With().
			Str("request_id", r.Header.Get(requestIDHeader)).
			Logger().
			WithContext(r.Context())

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware handles Cross-Origin Resource Sharing (CORS) headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, HEAD, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
