package domain

type CreatePlaceRequest struct {
	ZoneID          int64  `json:"zone_id" binding:"required"`
	Label           string `json:"label" binding:"required"`
	ConfigurationID *int64 `json:"configuration_id" binding:"omitempty"`
	SortOrder       int32  `json:"sort_order" binding:"omitempty"`
	IsActive        *bool  `json:"is_active" binding:"omitempty"`
}

type PatchPlaceRequest struct {
	Label           *string `json:"label" binding:"omitempty"`
	ConfigurationID *int64  `json:"configuration_id" binding:"omitempty"`
	SortOrder       *int32  `json:"sort_order" binding:"omitempty"`
	IsActive        *bool   `json:"is_active" binding:"omitempty"`
}
