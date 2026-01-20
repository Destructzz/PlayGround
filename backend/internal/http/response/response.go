package response

import (
	"net/http"
	"time"

	"backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string    `json:"status" example:"ready"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type MessageResponse struct {
	Message   string    `json:"message" example:"pong"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

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
