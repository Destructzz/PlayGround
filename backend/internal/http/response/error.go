package response

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Code      string    `json:"code" example:"auth_failed"`
	Message   string    `json:"message" example:"Authentication failed"`
	Details   any       `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type BindError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
	Value any    `json:"value,omitempty"`
}

func ParseBindError(err error) BindError {
	// 1️⃣ Ошибки валидации
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make([]FieldError, 0, len(ve))

		for _, fe := range ve {
			fields = append(fields, FieldError{
				Field: fe.Field(),
				Rule:  fe.Tag(),
				Value: fe.Value(),
			})
		}

		return BindError{
			Code:    "validation_failed",
			Message: "Validation failed",
			Fields:  fields,
		}
	}

	// 2️⃣ Лишние поля
	if strings.HasPrefix(err.Error(), "json: unknown field") {
		return BindError{
			Code:    "unknown_field",
			Message: err.Error(),
		}
	}

	// 3️⃣ Остальное (битый JSON, неверные типы)
	return BindError{
		Code:    "invalid_json",
		Message: "Invalid request payload",
	}
}
