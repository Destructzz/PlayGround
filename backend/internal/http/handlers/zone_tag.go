package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"backend/internal/service"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type ZoneTag struct {
	service *service.ZoneTagService
}

func NewZoneTag(service *service.ZoneTagService) *ZoneTag {
	return &ZoneTag{service: service}
}

// Create создает тег зоны.
// @Summary     Create zone tag
// @Description Creates a zone tag from JSON payload
// @Tags        zone-tags
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreateZoneTagRequest true "Zone tag payload"
// @Success     201 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone-tag [post]
func (h *ZoneTag) Create(c *gin.Context) {
	var req domain.CreateZoneTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("bad_request", fmt.Sprintf("zone tag %s is already exist", req.Name), nil),
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
		response.WithStatus(http.StatusCreated),
		response.WithData("zone_tag", result),
	).JSON(c)
}

// Get возвращает все теги зон.
// @Summary     List zone tags
// @Description Returns all zone tags
// @Tags        zone-tags
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone-tag [get]
func (h *ZoneTag) Get(c *gin.Context) {
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
		response.WithData("zone_tags", items),
	).JSON(c)
}

// GetByID возвращает тег зоны по id.
// @Summary     Get zone tag by id
// @Description Returns zone tag by path id
// @Tags        zone-tags
// @Produce     json
// @Param       id path int true "Zone tag ID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone-tag/{id} [get]
func (h *ZoneTag) GetByID(c *gin.Context) {
	id, ok := parseInt32Param(c, "id")
	if !ok {
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "Zone tag not found", nil),
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
		response.WithData("zone_tag", item),
	).JSON(c)
}

// Delete удаляет тег зоны.
// @Summary     Delete zone tag
// @Description Deletes zone tag by path id
// @Tags        zone-tags
// @Produce     json
// @Param       id path int true "Zone tag ID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone-tag/{id} [delete]
func (h *ZoneTag) Delete(c *gin.Context) {
	id, ok := parseInt32Param(c, "id")
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

// Patch обновляет тег зоны.
// @Summary     Patch zone tag
// @Description Updates zone tag by path id
// @Tags        zone-tags
// @Accept      json
// @Produce     json
// @Param       id path int true "Zone tag ID"
// @Param       payload body domain.PatchZoneTagRequest true "Zone tag patch payload"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone-tag/{id} [patch]
func (h *ZoneTag) Patch(c *gin.Context) {
	id, ok := parseInt32Param(c, "id")
	if !ok {
		return
	}

	var req domain.PatchZoneTagRequest
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
				response.WithError("not_found", "Zone tag not found", nil),
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
		response.WithData("zone_tag", item),
	).JSON(c)
}

func parseInt32Param(c *gin.Context, name string) (int32, bool) {
	rawID := c.Param(name)
	if rawID == "" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parametrs", nil),
		).JSON(c)
		return 0, false
	}

	id, err := strconv.ParseInt(rawID, 10, 32)
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid id", nil),
		).JSON(c)
		return 0, false
	}

	return int32(id), true
}
