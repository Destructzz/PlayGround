package response

import (
	"backend/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Zone(c *gin.Context, dto domain.CreateZoneRequest) {
	payload := gin.H{
		"zone":      dto,
		"timestamp": time.Now().UTC(),
	}

	withRequestID(c, payload)
	c.JSON(http.StatusOK, payload)
}

type ZoneResponse struct {
	Zone      domain.CreateZoneRequest `json:"zone"`
	Timestamp time.Time                `json:"timestamp"`
	RequestID string                   `json:"request_id,omitempty"`
}
