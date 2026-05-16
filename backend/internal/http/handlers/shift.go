package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"backend/internal/repo/sqlc"
	"backend/internal/service"
	"backend/pkg"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"math"
	"time"
)

type Shift struct {
	shiftService *service.ShiftService
}

func NewShift(shiftService *service.ShiftService) *Shift {
	return &Shift{shiftService: shiftService}
}

// Create creates a shift for the authenticated admin user.
// @Summary     Create shift
// @Description Creates a shift from JSON payload and attaches the current authenticated admin user
// @Tags        shifts
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreateShiftRequest true "Shift payload"
// @Success     201 {object} response.ShiftResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/shift [post]
func (s *Shift) Create(c *gin.Context) {
	user, ok := pkg.UserFromContext(c)
	if !ok {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusUnauthorized),
			response.WithError("unauthorized", "Authentication required", nil),
		).JSON(c)
		return
	}

	var req domain.CreateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	shift, err := s.shiftService.CreateShift(c.Request.Context(), user.ID, req)
	if err != nil {
		handleShiftWriteError(c, err, "Failed to create shift")
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusCreated),
		response.WithData("shift", domain.NewShiftViewForAdmin(shift.Shift, shift.User)),
	).JSON(c)
}

// Get returns all shifts.
// @Summary     List shifts
// @Description Returns all shifts with creator user data
// @Tags        shifts
// @Produce     json
// @Success     200 {object} response.ShiftListResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/shift [get]
func (s *Shift) Get(c *gin.Context) {
	shifts, err := s.shiftService.ListShifts(c.Request.Context())
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	user, ok := pkg.UserFromContext(c)
	if ok && user.Role == sqlc.RoleAdmin {
		response.NewResponseBuilder(
			response.WithData("shifts", domain.NewShiftListForAdmin(shifts)),
		).JSON(c)
	} else {
		now := time.Now()
		shiftViews := make([]domain.ShiftViewForUser, 0)
		for _, shift := range shifts {
			if shift.Shift.EndTime.Time.After(now) {
				shiftViews = append(shiftViews, domain.NewShiftViewForUser(shift.Shift))
			}
		}
		response.NewResponseBuilder(
			response.WithData("shifts", shiftViews),
		).JSON(c)
	}
}

// GetByID returns shift by id.
// @Summary     Get shift by id
// @Description Returns shift by id with creator user data
// @Tags        shifts
// @Produce     json
// @Param       id path int64 true "Shift ID"
// @Success     200 {object} response.ShiftResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/shift/{id} [get]
func (s *Shift) GetByID(c *gin.Context) {
	id, ok := pkg.ParseIDParam(c)
	if !ok {
		return
	}

	shift, err := s.shiftService.GetShiftByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "shift not found", nil),
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

	user, ok := pkg.UserFromContext(c)
	if ok && user.Role == sqlc.RoleAdmin {
		response.NewResponseBuilder(
			response.WithData("shift", domain.NewShiftViewForAdmin(shift.Shift, shift.User)),
		).JSON(c)
	} else {
		response.NewResponseBuilder(
			response.WithData("shift", domain.NewShiftViewForUser(shift.Shift)),
		).JSON(c)
	}
}

// GetByZoneTagID returns the shifts for a zone tag.
// @Summary     Get shifts by zone tag id
// @Description Returns all upcoming and active shifts for a given zone_tag_id
// @Tags        shifts
// @Produce     json
// @Param       zone_tag_id path int64 true "Zone Tag ID"
// @Success     200 {object} response.ShiftListResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/shift/zone-tag/{zone_tag_id} [get]
func (s *Shift) GetByZoneTagID(c *gin.Context) {
	rawID := c.Param("zone_tag_id")
	if rawID == "" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parameters", nil),
		).JSON(c)
		return
	}

	zoneTagID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || zoneTagID <= 0 {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "invalid zone_tag_id", nil),
		).JSON(c)
		return
	}

	shifts, err := s.shiftService.ListShiftsByZoneTagID(c.Request.Context(), zoneTagID)
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	now := time.Now()
	durations, _ := s.shiftService.GetServiceDurationsByZoneTagID(c.Request.Context(), zoneTagID)
	var step int32 = 60 // Default to 60 minutes
	if len(durations) > 0 {
		step = durations[0]
		for _, d := range durations {
			if d < step {
				step = d
			}
		}
	}

	shiftViews := make([]domain.ShiftViewForUser, 0)
	for _, shift := range shifts {
		// Show shift only if it hasn't ended yet
		if shift.Shift.EndTime.Time.After(now) {
			sView := domain.NewShiftViewForUser(shift.Shift)

			// If shift has started, adjust start_time to the next available slot
			if sView.StartTime.Before(now) {
				diff := now.Sub(sView.StartTime)
				intervals := int64(math.Ceil(diff.Minutes() / float64(step)))
				sView.StartTime = sView.StartTime.Add(time.Duration(intervals*int64(step)) * time.Minute)
			}

			// Only add if there's still time left in the shift after adjustment
			if sView.StartTime.Before(sView.EndTime) {
				shiftViews = append(shiftViews, sView)
			}
		}
	}
	response.NewResponseBuilder(
		response.WithData("shifts", shiftViews),
	).JSON(c)
}

// Patch updates shift by id.
// @Summary     Patch shift
// @Description Partially updates shift fields by id
// @Tags        shifts
// @Accept      json
// @Produce     json
// @Param       id path int64 true "Shift ID"
// @Param       payload body domain.PatchShiftRequest true "Shift patch payload"
// @Success     200 {object} response.ShiftResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/shift/{id} [patch]
func (s *Shift) Patch(c *gin.Context) {
	id, ok := pkg.ParseIDParam(c)
	if !ok {
		return
	}

	var req domain.PatchShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("update_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	shift, err := s.shiftService.PatchShift(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "shift not found", nil),
			).JSON(c)
			return
		}

		handleShiftWriteError(c, err, "Failed to update shift")
		return
	}

	response.NewResponseBuilder(
		response.WithData("shift", domain.NewShiftViewForAdmin(shift.Shift, shift.User)),
	).JSON(c)
}

// Delete soft-deletes shift by id.
// @Summary     Delete shift
// @Description Soft-deletes shift by id
// @Tags        shifts
// @Produce     json
// @Param       id path int64 true "Shift ID"
// @Success     200 {object} response.DeleteShiftResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/shift/{id} [delete]
func (s *Shift) Delete(c *gin.Context) {
	id, ok := pkg.ParseIDParam(c)
	if !ok {
		return
	}

	_, err := s.shiftService.DeleteShift(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", fmt.Sprintf("shift with %d ID doesn't exist", id), nil),
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

func handleShiftWriteError(c *gin.Context, err error, message string) {
	switch {
	case errors.Is(err, service.ErrInvalidShiftTime):
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("bad_request", "start_time and end_time must be valid RFC3339 timestamps", nil),
		).JSON(c)
		return
	case errors.Is(err, service.ErrInvalidShiftTimeRange):
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("bad_request", "end_time must be after start_time", nil),
		).JSON(c)
		return
	case errors.Is(err, service.ErrShiftOverlap):
		response.NewResponseBuilder(
			response.WithStatus(http.StatusConflict),
			response.WithError("conflict", "shift time window overlaps with an existing shift for this zone tag", nil),
		).JSON(c)
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("bad_request", fmt.Sprintf("referenced entity does not exist: %s", pgErr.ConstraintName), nil),
			).JSON(c)
			return
		case "23514":
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("bad_request", "end_time must be after start_time", nil),
			).JSON(c)
			return
		}
	}

	zap.L().Warn("database error", zap.Error(err))
	response.NewResponseBuilder(
		response.WithStatus(http.StatusInternalServerError),
		response.WithError("database_fault", message, nil),
	).JSON(c)
}
