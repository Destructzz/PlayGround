package response

import (
	"net/http"
	"time"

	"backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context, status string) {
	payload := gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
	}

	if rid, ok := middleware.RequestIDFromContext(c); ok {
		payload["request_id"] = rid
	}

	c.JSON(http.StatusOK, payload)
}

func Message(c *gin.Context, message string) {
	payload := gin.H{
		"message":   message,
		"timestamp": time.Now().UTC(),
	}

	if rid, ok := middleware.RequestIDFromContext(c); ok {
		payload["request_id"] = rid
	}

	c.JSON(http.StatusOK, payload)
}
