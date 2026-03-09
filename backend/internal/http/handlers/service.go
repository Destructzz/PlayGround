package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/response"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Service struct {
	serviceService *service.ServiceService
}

func (s *Service) Create(c *gin.Context) {
	var dto domain.CreateServiceRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		br := response.ParseBindError(err)
		response.NewResponseBuilder(
			response.WithStatus(400),
			response.WithError("creation_failed", "This data structure are not allowed", br),
		).JSON(c)
		return
	}

	result, err := s.serviceService.CreateService(c, dto)
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(500),
			response.WithError("creation_failed", "Failed to create service", err),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("service", result),
		response.WithStatus(201),
	).JSON(c)
}
