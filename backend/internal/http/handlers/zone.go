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

type Zone struct {
	zoneService *service.ZoneService
}

func NewZone(zoneService *service.ZoneService) *Zone {
	return &Zone{zoneService: zoneService}
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
// @Router      /api/v1/zone [post]
func (z *Zone) Create(c *gin.Context) {
	var req domain.CreateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	zone, err := z.zoneService.CreateZone(c.Request.Context(), req)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				response.NewResponseBuilder(
					response.WithStatus(http.StatusBadRequest),
					response.WithError("bad_request", fmt.Sprintf("zone %s is already exist", req.Name), nil),
				).JSON(c)
				return
			}
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
		response.WithData("zone", zone),
	).JSON(c)
}

// Get возвращает все зоны.
// @Summary     List zones
// @Description Returns all zones
// @Tags        zones
// @Produce     json
// @Success     200 {object} response.ZoneListResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone [get]
func (z *Zone) Get(c *gin.Context) {
	zones, err := z.zoneService.GetZones(c.Request.Context())

	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("zones", zones),
	).JSON(c)
}

// GetById возвращает зону по id.
// @Summary     Get zone by id
// @Description Returns zone id from path param
// @Tags        zones
// @Produce     json
// @Param       id path int64 true "Zone ID"
// @Success     200 {object} response.ZoneResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone/{id} [get]
func (z *Zone) GetById(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parametrs", nil),
		).JSON(c)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		zap.L().Warn("invalid id param", zap.String("value", rawID))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "invalid id", nil),
		).JSON(c)
		return
	}

	zone, err := z.zoneService.GetZoneByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("not_found", "zone not found", nil),
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
		response.WithData("zone", zone),
	).JSON(c)
}

// Delete удаляет зону по id.
// @Summary     Delete zone
// @Description Deletes zone by id
// @Tags        zones
// @Produce     json
// @Param       id path int64 true "Zone ID"
// @Success     200 {object} response.DeleteZoneResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone/{id} [delete]
func (z *Zone) Delete(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parametrs", nil),
		).JSON(c)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		zap.L().Warn("invalid id param", zap.String("value", rawID))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "invalid id", nil),
		).JSON(c)
		return
	}

	_, err = z.zoneService.DeleteByID(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("bad_request", fmt.Sprintf("zone with %d ID doesn't exist", id), nil),
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
		response.WithData("id", id),
	).JSON(c)
}

// Patch обновляет зону по id.
// @Summary     Patch zone
// @Description Partially updates zone fields by id
// @Tags        zones
// @Accept      json
// @Produce     json
// @Param       id path int64 true "Zone ID"
// @Param       payload body domain.PatchZoneRequest true "Zone patch payload"
// @Success     200 {object} response.ZoneResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/zone/{id} [patch]
func (z *Zone) Patch(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parametrs", nil),
		).JSON(c)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		zap.L().Warn("invalid id param", zap.String("value", rawID))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "invalid id", nil),
		).JSON(c)
		return
	}

	var req domain.PatchZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	zone, err := z.zoneService.PatchByID(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("not_found", "zone not found", nil),
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
		response.WithData("zone", zone),
	).JSON(c)
}
