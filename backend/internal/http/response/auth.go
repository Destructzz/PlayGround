package response

import (
	"time"
)

type AuthUserResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"user@example.com"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Name      string `json:"name" example:"Ada Lovelace"`
	Provider  string `json:"provider" example:"google"`
	Role      string `json:"role" example:"admin"`
}

type AuthResponse struct {
	User      AuthUserResponse `json:"user"`
	Timestamp time.Time        `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string           `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}
