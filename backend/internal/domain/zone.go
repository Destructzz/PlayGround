package domain

import (
	"encoding/json"
	"time"

	"backend/internal/repo/sqlc"
)

type Zone struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	ZoneType    sqlc.ZoneType   `json:"zone_type"`
	ZoneTagID   *int32          `json:"zone_tag_id"`
	Capacity    int32           `json:"capacity"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	DetailsJSON json.RawMessage `json:"details_json" swaggertype:"object"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateZoneRequest struct {
	Name        string        `json:"name" binding:"required"`
	Type        sqlc.ZoneType `json:"type" binding:"required,oneof=game event lounge sys"`
	ZoneTagID   *int32        `json:"zone_tag_id"`
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
