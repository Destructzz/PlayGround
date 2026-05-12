package response

import "time"

type PlaceResponse struct {
	Place PlaceDoc `json:"place"`
}

type PlaceListResponse struct {
	Places []PlaceDoc `json:"places"`
}

type DeletePlaceResponse struct {
	ID int64 `json:"id"`
}

type PlaceDoc struct {
	ID              int64     `json:"id" example:"1"`
	ZoneID          int64     `json:"zone_id" example:"2"`
	Label           string    `json:"label" example:"PC-1"`
	ConfigurationID *int64    `json:"configuration_id,omitempty" example:"1"`
	SortOrder       int32     `json:"sort_order" example:"1"`
	IsActive        bool      `json:"is_active" example:"true"`
	CreatedAt       time.Time `json:"created_at" example:"2026-01-19T15:37:27.514Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2026-01-19T15:37:27.514Z"`
}

