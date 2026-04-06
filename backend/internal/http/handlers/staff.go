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

type Staff struct {
	staffService *service.StaffService
}

func NewStaff(staffService *service.StaffService) *Staff {
	return &Staff{staffService: staffService}
}

// Create создает сотрудника.
// @Summary     Create staff
// @Description Creates a staff member from JSON payload
// @Tags        staff
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreateStaffRequest true "Staff payload"
// @Success     201 {object} response.StaffResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/staff [post]
func (s *Staff) Create(c *gin.Context) {
	var req domain.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	staff, err := s.staffService.CreateStaff(c.Request.Context(), req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				response.NewResponseBuilder(
					response.WithStatus(http.StatusBadRequest),
					response.WithError("bad_request", "staff member with this user_id already exists", nil),
				).JSON(c)
				return
			case "23503":
				response.NewResponseBuilder(
					response.WithStatus(http.StatusBadRequest),
					response.WithError("bad_request", fmt.Sprintf("referenced entity does not exist: %s", pgErr.ConstraintName), nil),
				).JSON(c)
				return
			}
		}

		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("creation_failed", "Failed to create staff member", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusCreated),
		response.WithData("staff", staff),
	).JSON(c)
}

// Get возвращает всех сотрудников.
// @Summary     List staff
// @Description Returns all staff members
// @Tags        staff
// @Produce     json
// @Success     200 {object} response.StaffListResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/staff [get]
func (s *Staff) Get(c *gin.Context) {
	staffList, err := s.staffService.ListStaff(c.Request.Context())
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("staff", staffList),
	).JSON(c)
}

// GetById возвращает сотрудника по id.
// @Summary     Get staff by id
// @Description Returns staff member by id from path param
// @Tags        staff
// @Produce     json
// @Param       id path int64 true "Staff ID"
// @Success     200 {object} response.StaffResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/staff/{id} [get]
func (s *Staff) GetById(c *gin.Context) {
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

	staff, err := s.staffService.GetStaffByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "staff member not found", nil),
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
		response.WithData("staff", staff),
	).JSON(c)
}

// Patch обновляет сотрудника по id.
// @Summary     Patch staff
// @Description Partially updates staff fields by id
// @Tags        staff
// @Accept      json
// @Produce     json
// @Param       id path int64 true "Staff ID"
// @Param       payload body domain.PatchStaffRequest true "Staff patch payload"
// @Success     200 {object} response.StaffResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/staff/{id} [patch]
func (s *Staff) Patch(c *gin.Context) {
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

	var req domain.PatchStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("update_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	staff, err := s.staffService.PatchStaff(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "staff member not found", nil),
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
		response.WithData("staff", staff),
	).JSON(c)
}

// Delete удаляет сотрудника по id (soft delete).
// @Summary     Delete staff
// @Description Soft-deletes staff member by id
// @Tags        staff
// @Produce     json
// @Param       id path int64 true "Staff ID"
// @Success     200 {object} response.DeleteStaffResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/staff/{id} [delete]
func (s *Staff) Delete(c *gin.Context) {
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

	_, err = s.staffService.DeleteStaff(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", fmt.Sprintf("staff with %d ID doesn't exist", id), nil),
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
