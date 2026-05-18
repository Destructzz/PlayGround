package domain

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	ZoneID      int64  `json:"zone_id" binding:"required"`
	Duration    int64  `json:"duration" binding:"required"`
	Price       any    `json:"price" binding:"required"`
	Currency    string `json:"currency" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
	DetailsJSON string `json:"details_json" binding:"omitempty"`
}

type UpdateServiceRequest struct {
	Name        *string `json:"name" binding:"omitempty"`
	ZoneID      *int64  `json:"zone_id" binding:"omitempty"`
	Duration    *int64  `json:"duration" binding:"omitempty"`
	Price       *any    `json:"price" binding:"omitempty"`
	Currency    *string `json:"currency" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
	IsActive    *bool   `json:"is_active" binding:"omitempty"`
	DetailsJSON *string `json:"details_json" binding:"omitempty"`
}