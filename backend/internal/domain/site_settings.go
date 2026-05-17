package domain

import "encoding/json"

type SiteSettings struct {
	ID               int32           `json:"id"`
	SettingsJSON     json.RawMessage `json:"settings_json" swaggertype:"array,number"`
	GalleryItemsJSON json.RawMessage `json:"gallery_items_json" swaggertype:"array,object"`
}

type UpdateSiteSettingsRequest struct {
	SettingsJSON     json.RawMessage `json:"settings_json" binding:"required"`
	GalleryItemsJSON json.RawMessage `json:"gallery_items_json" binding:"omitempty"`
}
