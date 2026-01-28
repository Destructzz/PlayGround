package response

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Zone(c *gin.Context, zone sqlc.Zone) {
	payload := gin.H{
		"zone":      zone,
		"timestamp": time.Now().UTC(),
	}

	withRequestID(c, payload)
	c.JSON(http.StatusCreated, payload)
}

type ZoneResponse struct {
	Zone      domain.CreateZoneRequest `json:"zone"`
	Timestamp time.Time                `json:"timestamp"`
	RequestID string                   `json:"request_id,omitempty"`
}
