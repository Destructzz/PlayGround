package response

import (
	"backend/internal/repo/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateZone(c *gin.Context, zone sqlc.Zone) {
	payload := gin.H{
		"zone":      zone,
		"timestamp": time.Now().UTC(),
	}

	withRequestID(c, payload)
	c.JSON(http.StatusCreated, payload)
}

func GetZoneByID(c *gin.Context, zone sqlc.Zone) {
	payload := gin.H{
		"zone":      zone,
		"timestamp": time.Now().UTC(),
	}

	withRequestID(c, payload)
	c.JSON(http.StatusOK, payload)
}

func GetZone(c *gin.Context, zones []sqlc.Zone) {
	payload := gin.H{
		"zones":     zones,
		"timestamp": time.Now().UTC(),
	}

	withRequestID(c, payload)
	c.JSON(http.StatusOK, payload)
}

type ZoneResponse struct {
	Zone      ZoneDoc   `json:"zone"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ZoneListResponse struct {
	Zones     []ZoneDoc `json:"zones"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ZoneDoc struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" example:"Main Hall"`
	ZoneType    string    `json:"zone_type" example:"game"`
	Capacity    int32     `json:"capacity" example:"20"`
	Description *string   `json:"description,omitempty" example:"Large hall"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-19T15:37:27.514667373Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-01-19T15:37:27.514667373Z"`
}
