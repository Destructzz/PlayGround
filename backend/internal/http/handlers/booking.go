package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"backend/internal/service"
	"backend/pkg"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Booking struct {
	bookingService *service.BookingService
}

func NewBooking(bookingService *service.BookingService) *Booking {
	return &Booking{bookingService: bookingService}
}

// Create создает бронирование.
// @Summary     Create booking
// @Description Creates a booking from JSON payload (user derived from session)
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreateBookingRequest true "Booking payload"
// @Success     201 {object} response.BookingResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     409 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/booking [post]
func (b *Booking) Create(c *gin.Context) {
	user, ok := pkg.UserFromContext(c)
	if !ok {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusUnauthorized),
			response.WithError("unauthorized", "Authentication required to create a booking", nil),
		).JSON(c)
		return
	}

	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	booking, err := b.bookingService.CreateBooking(c.Request.Context(), user.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookingOverlap):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusConflict),
				response.WithError("booking_overlap", "A booking already exists for this resource and time window", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrPlaceUnavailable):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusConflict),
				response.WithError("place_unavailable", "This place is already booked for the selected time", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrCapacityExceeded):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusConflict),
				response.WithError("capacity_exceeded", "Participants exceed available capacity for this time slot", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrInvalidBookingRange):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("invalid_booking_range", "Start time must be earlier than end time", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrServiceZoneMismatch):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("service_zone_mismatch", "Selected service does not belong to the selected zone", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrPlaceZoneMismatch):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("place_zone_mismatch", "Selected place does not belong to the selected zone", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrInactiveService):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusConflict),
				response.WithError("inactive_service", "Selected service is inactive", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrInactivePlace):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusConflict),
				response.WithError("inactive_place", "Selected place is inactive", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrInvalidServiceWindow):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("invalid_service_window", "Booking duration must match the selected service", nil),
			).JSON(c)
			return
		case errors.Is(err, service.ErrPastBooking):
			response.NewResponseBuilder(
				response.WithStatus(http.StatusBadRequest),
				response.WithError("past_booking", "Cannot create a booking in the past", nil),
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
			}
		}

		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("creation_failed", "Failed to create booking", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusCreated),
		response.WithData("booking", booking),
	).JSON(c)
}

// Get возвращает все бронирования.
// @Summary     List bookings
// @Description Returns all bookings
// @Tags        bookings
// @Produce     json
// @Success     200 {object} response.BookingListResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/booking [get]
func (b *Booking) Get(c *gin.Context) {
	bookings, err := b.bookingService.ListBookings(c.Request.Context())
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("bookings", bookings),
	).JSON(c)
}

// MyBookings возвращает бронирования текущего пользователя.
// @Summary     List my bookings
// @Description Returns bookings for the authenticated user
// @Tags        bookings
// @Produce     json
// @Success     200 {object} response.BookingListResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/booking/my [get]
func (b *Booking) MyBookings(c *gin.Context) {
	user, ok := pkg.UserFromContext(c)
	if !ok {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusUnauthorized),
			response.WithError("unauthorized", "Authentication required to get your bookings", nil),
		).JSON(c)
		return
	}

	bookings, err := b.bookingService.ListBookingsByUserID(c.Request.Context(), user.ID)
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("bookings", bookings),
	).JSON(c)
}

// GetById возвращает бронирование по id.
// @Summary     Get booking by id
// @Description Returns booking by id from path param
// @Tags        bookings
// @Produce     json
// @Param       id path int64 true "Booking ID"
// @Success     200 {object} response.BookingResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/booking/{id} [get]
func (b *Booking) GetById(c *gin.Context) {
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

	booking, err := b.bookingService.GetBookingByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "booking not found", nil),
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
		response.WithData("booking", booking),
	).JSON(c)
}

// Patch обновляет бронирование по id.
// @Summary     Patch booking
// @Description Partially updates booking fields by id
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Param       id path int64 true "Booking ID"
// @Param       payload body domain.PatchBookingRequest true "Booking patch payload"
// @Success     200 {object} response.BookingResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/booking/{id} [patch]
func (b *Booking) Patch(c *gin.Context) {
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

	var req domain.PatchBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("update_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	booking, err := b.bookingService.PatchBooking(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "booking not found", nil),
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
		response.WithData("booking", booking),
	).JSON(c)
}

// Delete удаляет бронирование по id.
// @Summary     Delete booking
// @Description Deletes booking by id
// @Tags        bookings
// @Produce     json
// @Param       id path int64 true "Booking ID"
// @Success     200 {object} response.DeleteBookingResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/booking/{id} [delete]
func (b *Booking) Delete(c *gin.Context) {
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

	_, err = b.bookingService.DeleteBooking(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", fmt.Sprintf("booking with %d ID doesn't exist", id), nil),
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
