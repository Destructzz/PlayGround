package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger logs HTTP requests using zap.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		path := c.Request.URL.Path
		route := c.FullPath()
		if route == "" {
			route = path
		}

		logger := zap.L()
		if rid, ok := RequestIDFromContext(c); ok {
			logger = logger.With(zap.String("request_id", rid))
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("route", route),
			zap.String("host", c.Request.Host),
			zap.String("scheme", c.Request.URL.Scheme),
			zap.String("proto", c.Request.Proto),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.Int("status", status),
			zap.Int64("duration_ms", latency.Milliseconds()),
			zap.Int64("bytes_in", c.Request.ContentLength),
			zap.Int("bytes_out", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case status >= 500:
			logger.Error("http request", fields...)
		case status >= 400:
			logger.Warn("http request", fields...)
		default:
			logger.Info("http request", fields...)
		}
	}
}
