package domain

import (
	"backend/internal/repo/sqlc"
)

type CreateZoneRequest struct {
	Name        string        `json:"name" binding:"required"`
	Type        sqlc.ZoneType `json:"type" binding:"required,oneof=game event vip lounge sys"`
	Capacity    int           `json:"capacity" binding:"required,min=1"`
	Description string        `json:"description"`
	IsActive    *bool         `json:"is_active"`
}
