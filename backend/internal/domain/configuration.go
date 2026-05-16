package domain

import "encoding/json"

type ComputerSpecificationEntry struct {
	Title string `json:"title" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type CreateComputerConfigurationRequest struct {
	ZoneTagID int64                        `json:"zone_tag_id" binding:"required"`
	SpecsJSON []ComputerSpecificationEntry `json:"specs_json" binding:"required,min=1,dive"`
}

type ComputerConfiguration struct {
	ID        int64           `json:"id"`
	ZoneTagID int64           `json:"zone_tag_id"`
	SpecsJSON json.RawMessage `json:"specs_json" swaggertype:"object"`
}

type PatchComputerConfigurationRequest struct {
	ZoneTagID *int64                        `json:"zone_tag_id" binding:"omitempty"`
	SpecsJSON *[]ComputerSpecificationEntry `json:"specs_json" binding:"omitempty"`
}
