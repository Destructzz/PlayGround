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

type Payment struct {
	paymentService *service.PaymentService
}

func NewPayment(paymentService *service.PaymentService) *Payment {
	return &Payment{paymentService: paymentService}
}

// Create создает платёж.
// @Summary     Create payment
// @Description Creates a payment from JSON payload
// @Tags        payments
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreatePaymentRequest true "Payment payload"
// @Success     201 {object} response.PaymentResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/payment [post]
func (p *Payment) Create(c *gin.Context) {
	var req domain.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	payment, err := p.paymentService.CreatePayment(c.Request.Context(), req)
	if err != nil {
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
			response.WithError("creation_failed", "Failed to create payment", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusCreated),
		response.WithData("payment", payment),
	).JSON(c)
}

// Get возвращает все платежи.
// @Summary     List payments
// @Description Returns all payments
// @Tags        payments
// @Produce     json
// @Success     200 {object} response.PaymentListResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/payment [get]
func (p *Payment) Get(c *gin.Context) {
	payments, err := p.paymentService.ListPayments(c.Request.Context())
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("payments", payments),
	).JSON(c)
}

// GetById возвращает платёж по id.
// @Summary     Get payment by id
// @Description Returns payment by id from path param
// @Tags        payments
// @Produce     json
// @Param       id path int64 true "Payment ID"
// @Success     200 {object} response.PaymentResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/payment/{id} [get]
func (p *Payment) GetById(c *gin.Context) {
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

	payment, err := p.paymentService.GetPaymentByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "payment not found", nil),
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
		response.WithData("payment", payment),
	).JSON(c)
}

// Patch обновляет платёж по id.
// @Summary     Patch payment
// @Description Partially updates payment fields by id
// @Tags        payments
// @Accept      json
// @Produce     json
// @Param       id path int64 true "Payment ID"
// @Param       payload body domain.PatchPaymentRequest true "Payment patch payload"
// @Success     200 {object} response.PaymentResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/payment/{id} [patch]
func (p *Payment) Patch(c *gin.Context) {
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

	var req domain.PatchPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("update_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	payment, err := p.paymentService.PatchPayment(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "payment not found", nil),
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
		response.WithData("payment", payment),
	).JSON(c)
}

// Delete удаляет платёж по id.
// @Summary     Delete payment
// @Description Deletes payment by id
// @Tags        payments
// @Produce     json
// @Param       id path int64 true "Payment ID"
// @Success     200 {object} response.DeletePaymentResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/payment/{id} [delete]
func (p *Payment) Delete(c *gin.Context) {
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

	_, err = p.paymentService.DeletePayment(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", fmt.Sprintf("payment with %d ID doesn't exist", id), nil),
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
