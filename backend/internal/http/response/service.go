package response

import (
	"time"
)

type DeleteServiceResponse struct {
	ID        int64     `json:"id" example:"1"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ServiceResponse struct {
	Service   ServiceDoc `json:"service"`
	Timestamp time.Time  `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string     `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ServiceListResponse struct {
	Services  []ServiceDoc `json:"services"`
	Timestamp time.Time    `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string       `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ServiceDoc struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" example:"VIP Gaming"`
	ZoneID      int64     `json:"zone_id" example:"1"`
	Duration    int32     `json:"duration" example:"60"`
	Price       string    `json:"price" example:"15.00"`
	Currency    string    `json:"currency" example:"USD"`
	Description *string   `json:"description,omitempty" example:"1 hour of VIP PS5"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-19T15:37:27.514667373Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-01-19T15:37:27.514667373Z"`
}
