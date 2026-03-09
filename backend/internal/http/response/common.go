package response

import (
	"backend/internal/http/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseBuilder struct {
	status  int
	payload gin.H
}

type ResponseOption func(*ResponseBuilder) error

func NewResponseBuilder(options ...ResponseOption) *ResponseBuilder {
	rb := &ResponseBuilder{
		status: http.StatusOK,
		payload: gin.H{
			"timestamp": time.Now().UTC(),
		},
	}

	for _, opt := range options {
		_ = opt(rb) // Игнорируем ошибку, поскольку наши опции сейчас не падают
	}

	return rb
}

func WithData(name string, obj any) ResponseOption {
	return func(rb *ResponseBuilder) error {
		if name != "" {
			rb.payload[name] = obj
		}
		return nil
	}
}

func WithStatus(status int) ResponseOption {
	return func(rb *ResponseBuilder) error {
		rb.status = status
		return nil
	}
}

func WithError(code string, msg string, details any) ResponseOption {
	return func(rb *ResponseBuilder) error {
		rb.payload["code"] = code
		rb.payload["message"] = msg
		if details != nil {
			rb.payload["details"] = details
		}
		return nil
	}
}

func (rb *ResponseBuilder) JSON(c *gin.Context) {
	withRequestID(c, rb.payload)
	c.JSON(rb.status, rb.payload)
}

func withRequestID(c *gin.Context, payload gin.H) {
	if rid, ok := middleware.RequestIDFromContext(c); ok {
		payload["request_id"] = rid
	}
}

type ParamResponse struct {
	Param     string    `json:"param" example:"123"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}
