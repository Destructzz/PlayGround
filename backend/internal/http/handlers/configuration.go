package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"backend/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ComputerConfiguration struct {
	service *service.ComputerConfigurationService
}

func NewComputerConfiguration(service *service.ComputerConfigurationService) *ComputerConfiguration {
	return &ComputerConfiguration{service: service}
}

// Create создает конфигурацию компьютера.
// @Summary     Create configuration
// @Description Creates a computer configuration from JSON payload
// @Tags        configurations
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreateComputerConfigurationRequest true "Configuration payload"
// @Success     201 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/configuration [post]
func (h *ComputerConfiguration) Create(c *gin.Context) {
	var req domain.CreateComputerConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	item, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("creation_failed", "Failed to create configuration", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusCreated),
		response.WithData("configuration", item),
	).JSON(c)
}

// Get возвращает все конфигурации компьютеров.
// @Summary     List configurations
// @Description Returns all computer configurations
// @Tags        configurations
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/configuration [get]
func (h *ComputerConfiguration) Get(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("configurations", items),
	).JSON(c)
}

// GetByID возвращает конфигурацию компьютера по id.
// @Summary     Get configuration by id
// @Description Returns computer configuration by path id
// @Tags        configurations
// @Produce     json
// @Param       id path int true "Configuration ID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/configuration/{id} [get]
func (h *ComputerConfiguration) GetByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "Configuration not found", nil),
			).JSON(c)
			return
		}

		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("configuration", item),
	).JSON(c)
}

// Delete удаляет конфигурацию компьютера.
// @Summary     Delete configuration
// @Description Deletes computer configuration by path id
// @Tags        configurations
// @Produce     json
// @Param       id path int true "Configuration ID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/configuration/{id} [delete]
func (h *ComputerConfiguration) Delete(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}

	deletedID, err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("id", deletedID),
	).JSON(c)
}

// Patch обновляет конфигурацию компьютера.
// @Summary     Patch configuration
// @Description Updates computer configuration by path id
// @Tags        configurations
// @Accept      json
// @Produce     json
// @Param       id path int true "Configuration ID"
// @Param       payload body domain.PatchComputerConfigurationRequest true "Configuration patch payload"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/configuration/{id} [patch]
func (h *ComputerConfiguration) Patch(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}

	var req domain.PatchComputerConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("patch_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	item, err := h.service.Patch(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "Configuration not found", nil),
			).JSON(c)
			return
		}

		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("configuration", item),
	).JSON(c)
}

func parseInt64Param(c *gin.Context, name string) (int64, bool) {
	rawID := c.Param(name)
	if rawID == "" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parametrs", nil),
		).JSON(c)
		return 0, false
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid id", nil),
		).JSON(c)
		return 0, false
	}

	return id, true
}
