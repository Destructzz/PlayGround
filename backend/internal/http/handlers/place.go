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

type Place struct {
	placeService *service.PlaceService
}

func NewPlace(placeService *service.PlaceService) *Place {
	return &Place{placeService: placeService}
}

// Create создает место.
// @Summary     Create place
// @Description Creates a place from JSON payload
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreatePlaceRequest true "Place payload"
// @Success     200 {object} response.PlaceResponse
// @Failure     400 {object} response.ErrorResponse
// @Router      /api/v1/place [post]
func (p *Place) Create(c *gin.Context) {
	var req domain.CreatePlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure is not allowed", br),
		).JSON(c)
		return
	}

	place, err := p.placeService.CreatePlace(c.Request.Context(), req)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				response.NewResponseBuilder(
					response.WithStatus(http.StatusBadRequest),
					response.WithError("bad_request", fmt.Sprintf("place %s already exists", req.Label), nil),
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
		response.WithData("place", place),
	).JSON(c)
}

// Get возвращает все места.
// @Summary     List places
// @Description Returns all places
// @Tags        places
// @Produce     json
// @Success     200 {object} response.PlaceListResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/place [get]
func (p *Place) Get(c *gin.Context) {
	places, err := p.placeService.GetPlaces(c.Request.Context())

	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("places", places),
	).JSON(c)
}

// GetById возвращает место по id.
// @Summary     Get place by id
// @Description Returns place by id from path param
// @Tags        places
// @Produce     json
// @Param       id path int64 true "Place ID"
// @Success     200 {object} response.PlaceResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/place/{id} [get]
func (p *Place) GetById(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parameters", nil),
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

	place, err := p.placeService.GetPlaceByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "place not found", nil),
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
		response.WithData("place", place),
	).JSON(c)
}

// Delete удаляет место по id.
// @Summary     Delete place
// @Description Deletes place by id
// @Tags        places
// @Produce     json
// @Param       id path int64 true "Place ID"
// @Success     200 {object} response.DeletePlaceResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/place/{id} [delete]
func (p *Place) Delete(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parameters", nil),
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

	_, err = p.placeService.DeleteByID(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("bad_request", fmt.Sprintf("place with %d ID doesn't exist", id), nil),
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

// Patch обновляет место по id.
// @Summary     Patch place
// @Description Partially updates place fields by id
// @Tags        places
// @Accept      json
// @Produce     json
// @Param       id path int64 true "Place ID"
// @Param       payload body domain.PatchPlaceRequest true "Place patch payload"
// @Success     200 {object} response.PlaceResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/place/{id} [patch]
func (p *Place) Patch(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parameters", nil),
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

	var req domain.PatchPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure is not allowed", br),
		).JSON(c)
		return
	}

	place, err := p.placeService.PatchByID(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "place not found", nil),
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
		response.WithData("place", place),
	).JSON(c)
}
