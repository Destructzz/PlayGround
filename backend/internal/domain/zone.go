package domain

import (
	"backend/internal/repo/sqlc"
)

type CreateZoneRequest struct {
	Name        string        `json:"name" binding:"required"`
	Type        sqlc.ZoneType `json:"type" binding:"required,oneof=game event lounge sys"`
	ZoneTagID   int32         `json:"zone_tag_id" binding:"required"`
	Capacity    int           `json:"capacity" binding:"required,min=1"`
	Description string        `json:"description"`
	IsActive    *bool         `json:"is_active"`
	DetailsJSON string        `json:"details_json"`
}

type PatchZoneRequest struct {
	Name        *string        `json:"name" binding:"omitempty"`
	Type        *sqlc.ZoneType `json:"type" binding:"omitempty,oneof=game event lounge sys"`
	ZoneTagID   *int32         `json:"zone_tag_id" binding:"omitempty"`
	Capacity    *int           `json:"capacity" binding:"omitempty,min=1"`
	Description *string        `json:"description"`
	IsActive    *bool          `json:"is_active"`
	DetailsJSON *string        `json:"details_json"`
}
