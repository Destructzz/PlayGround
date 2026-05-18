package handlers

import (
	"backend/internal/http/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	service *service.StatsService
}

func NewStatsHandler(service *service.StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

// GetAdminStats returns dashboard statistics
// @Summary Get admin dashboard statistics
// @Description Returns summary, revenue and booking statistics for the admin dashboard
// @Tags Admin
// @Produce json
// @Success 200 {object} domain.AdminStatsResponse
// @Router /api/v1/admin/stats [get]
func (h *StatsHandler) GetAdminStats(c *gin.Context) {
	stats, err := h.service.GetAdminStats(c.Request.Context())
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("internal", "Failed to get admin stats", nil),
		).JSON(c)
		return
	}

	c.JSON(http.StatusOK, stats)
}
