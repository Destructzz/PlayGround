package domain

type ComputerSpecificationEntry struct {
	Title string `json:"title" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type CreateComputerConfigurationRequest struct {
	ZoneTagID int64                        `json:"zone_tag_id" binding:"required"`
	SpecsJSON []ComputerSpecificationEntry `json:"specs_json" binding:"required,min=1,dive"`
}

type PatchComputerConfigurationRequest struct {
	ZoneTagID *int64                        `json:"zone_tag_id" binding:"omitempty"`
	SpecsJSON *[]ComputerSpecificationEntry `json:"specs_json" binding:"omitempty"`
}
