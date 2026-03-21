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

type Service struct {
	serviceService *service.ServiceService
}

func NewService(serviceService *service.ServiceService) *Service {
	return &Service{serviceService: serviceService}
}

// Create создает сервис.
// @Summary     Create service
// @Description Creates a service from JSON payload
// @Tags        services
// @Accept      json
// @Produce     json
// @Param       payload body domain.CreateServiceRequest true "Service payload"
// @Success     200 {object} response.ServiceResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/service [post]
func (s *Service) Create(c *gin.Context) {
	var dto domain.CreateServiceRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	result, err := s.serviceService.CreateService(c.Request.Context(), dto)
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("creation_failed", "Failed to create service", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("service", result),
		response.WithStatus(http.StatusCreated),
	).JSON(c)
}

// Get возвращает все сервисы.
// @Summary     List services
// @Description Returns all services
// @Tags        services
// @Produce     json
// @Success     200 {object} response.ServiceListResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/service [get]
func (s *Service) Get(c *gin.Context) {
	result, err := s.serviceService.ListServices(c.Request.Context())

	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithStatus(http.StatusOK),
		response.WithData("services", result),
	).JSON(c)
}

// GetById возвращает сервис по id.
// @Summary     Get service by id
// @Description Returns service id from path param
// @Tags        services
// @Produce     json
// @Param       id path int64 true "Service ID"
// @Success     200 {object} response.ServiceResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/service/{id} [get]
func (s *Service) GetById(c *gin.Context) {
	rawID := c.Param("id")

	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid service ID", nil),
		).JSON(c)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		zap.L().Warn("invalid id", zap.Error(err))	
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid service ID", nil),
		).JSON(c)
		return
	}

	service, err := s.serviceService.GetServiceByID(c.Request.Context(), id)

	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "service not found", nil),
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
		response.WithStatus(http.StatusOK),
		response.WithData("service", service),
	).JSON(c)
}

// Patch обновляет сервис по id.
// @Summary     Patch service
// @Description Partially updates service fields by id
// @Tags        services
// @Accept      json
// @Produce     json
// @Param       id path int64 true "Service ID"
// @Param       payload body domain.UpdateServiceRequest true "Service patch payload"
// @Success     201 {object} response.ServiceResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/service/{id} [patch]
func (s *Service) Patch(c *gin.Context){
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid service ID", nil),
		).JSON(c)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		zap.L().Warn("invalid id", zap.Error(err))	
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid service ID", nil),
		).JSON(c)
		return
	}

	var dto domain.UpdateServiceRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	result, err := s.serviceService.UpdateService(c.Request.Context(), id, dto)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "service not found", nil),
			).JSON(c)
			return
		}
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("update_failed", "Failed to update service", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("service", result),
		response.WithStatus(http.StatusCreated),
	).JSON(c)
}

// Delete удаляет сервис по id.
// @Summary     Delete service
// @Description Deletes service by id
// @Tags        services
// @Produce     json
// @Param       id path int64 true "Service ID"
// @Success     200 {object} response.DeleteServiceResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/service/{id} [delete]
func (s *Service) Delete(c *gin.Context){
	rawID := c.Param("id")
	if rawID == "" {
		zap.L().Warn("missing id param")
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid service ID", nil),
		).JSON(c)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		zap.L().Warn("invalid id", zap.Error(err))	
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_id", "Invalid service ID", nil),
		).JSON(c)
		return
	}

	result, err := s.serviceService.DeleteService(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			zap.L().Warn("incorrect id", zap.Error(err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusNotFound),
				response.WithError("not_found", "service not found", nil),
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
		response.WithStatus(http.StatusOK),
		response.WithData("id", result),
	).JSON(c)
}



