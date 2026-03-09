package response

import (
	"time"
)

type DeleteZoneResponse struct {
	ID        int64     `json:"id" example:"1"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ZoneResponse struct {
	Zone      ZoneDoc   `json:"zone"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ZoneListResponse struct {
	Zones     []ZoneDoc `json:"zones"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ZoneDoc struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" example:"Main Hall"`
	ZoneType    string    `json:"zone_type" example:"game"`
	Capacity    int32     `json:"capacity" example:"20"`
	Description *string   `json:"description,omitempty" example:"Large hall"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-19T15:37:27.514667373Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-01-19T15:37:27.514667373Z"`
}
