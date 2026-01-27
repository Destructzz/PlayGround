package response

import (
	"backend/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func withRequestID(c *gin.Context, payload gin.H) {
	if rid, ok := middleware.RequestIDFromContext(c); ok {
		payload["request_id"] = rid
	}
}
