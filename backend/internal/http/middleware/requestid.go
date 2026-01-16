package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"
const requestIDHeader = "X-Request-ID"

// RequestID injects a request id into the context and response headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if rid == "" {
			rid = uuid.New().String()
		}

		c.Set(requestIDKey, rid)
		c.Writer.Header().Set(requestIDHeader, rid)

		c.Next()
	}
}

// RequestIDFromContext retrieves the request id stored by RequestID middleware.
func RequestIDFromContext(c *gin.Context) (string, bool) {
	if value, exists := c.Get(requestIDKey); exists {
		if rid, ok := value.(string); ok {
			return rid, true
		}
	}
	return "", false
}
