package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Zone struct{}

func NewZone() *Zone {
	return &Zone{}
}

// Create создает зону.
// @Summary     Create zone
// @Description Creates a zone from JSON payload
// @Tags        zones
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreateZoneRequest true "Zone payload"
// @Success     200 {object} response.ZoneResponse
// @Failure     400 {object} response.ErrorResponse
// @Router      /api/v1/zones [post]
func (z *Zone) Create(c *gin.Context) {
	var req domain.CreateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.Error(c, http.StatusBadRequest, "creation_failed", "This data structure are not allowed", br)
		return
	}

	response.Zone(c, req)
}
