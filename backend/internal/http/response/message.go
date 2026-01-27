package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageResponse struct {
	Message   string    `json:"message" example:"pong"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

func Message(c *gin.Context, message string) {
	payload := gin.H{
		"message":   message,
		"timestamp": time.Now().UTC(),
	}

	withRequestID(c, payload)
	c.JSON(http.StatusOK, payload)
}
