package response

import (
	"time"
)

type MessageResponse struct {
	Message   string    `json:"message" example:"pong"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}
