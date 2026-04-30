package handlers

import (
	"backend/internal/http/response"
	"backend/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Seed struct {
	service *service.SeedService
}

func NewSeed(service *service.SeedService) *Seed {
	return &Seed{service: service}
}

// Get возвращает текущее состояние seed-данных.
// @Summary     Get seed state
// @Description Returns current seed data snapshot in non-production environments
// @Tags        seed
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/seed [get]
func (h *Seed) Get(c *gin.Context) {
	if !seedAllowed(c) {
		return
	}

	snapshot, err := h.service.Get(c.Request.Context())
	if err != nil {
		zap.L().Warn("seed get failed", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("seed_error", "Failed to load seed data", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("seed", snapshot),
	).JSON(c)
}

// Post создает базовые seed-данные для gaming.
// @Summary     Apply seed
// @Description Creates base gaming seed data without duplicates in non-production environments
// @Tags        seed
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/seed [post]
func (h *Seed) Post(c *gin.Context) {
	if !seedAllowed(c) {
		return
	}

	snapshot, err := h.service.Apply(c.Request.Context())
	if err != nil {
		zap.L().Warn("seed apply failed", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("seed_error", "Failed to apply seed data", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("seed", snapshot),
	).JSON(c)
}

// Delete удаляет базовые seed-данные для gaming.
// @Summary     Delete seed
// @Description Deletes base gaming seed data in non-production environments
// @Tags        seed
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/seed [delete]
func (h *Seed) Delete(c *gin.Context) {
	if !seedAllowed(c) {
		return
	}

	snapshot, err := h.service.Delete(c.Request.Context())
	if err != nil {
		zap.L().Warn("seed delete failed", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("seed_error", "Failed to delete seed data", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("seed", snapshot),
	).JSON(c)
}

func seedAllowed(c *gin.Context) bool {
	if os.Getenv("APP_ENV") == "production" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusForbidden),
			response.WithError("forbidden", "Seed endpoint is not available in production", nil),
		).JSON(c)
		return false
	}

	return true
}
