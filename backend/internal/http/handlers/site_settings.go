package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SiteSettings struct {
	service *service.SiteSettingsService
}

func NewSiteSettings(service *service.SiteSettingsService) *SiteSettings {
	return &SiteSettings{service: service}
}

// Get возвращает настройки сайта.
// @Summary     Get site settings
// @Description Returns the global site configuration settings
// @Tags        settings
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/settings [get]
func (h *SiteSettings) Get(c *gin.Context) {
	item, err := h.service.Get(c.Request.Context())
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("settings", item),
	).JSON(c)
}

// Upsert обновляет или вставляет настройки сайта.
// @Summary     Update site settings
// @Description Creates or updates the global site configuration settings
// @Tags        settings
// @Accept      json
// @Produce     json
// @Param       payload body domain.UpdateSiteSettingsRequest true "Settings payload"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/settings [post]
func (h *SiteSettings) Upsert(c *gin.Context) {
	var req domain.UpdateSiteSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("update_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	item, err := h.service.Upsert(c.Request.Context(), req)
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("update_failed", "Failed to update settings", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("settings", item),
	).JSON(c)
}
