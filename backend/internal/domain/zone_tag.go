package domain

type CreateZoneTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type PatchZoneTagRequest struct {
	Name *string `json:"name" binding:"omitempty"`
}
