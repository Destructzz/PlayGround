package handlers

import (
	"backend/internal/http/response"
	"backend/internal/repo/sqlc"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Public struct {
	queries *sqlc.Queries
}

func NewPublic(queries *sqlc.Queries) *Public {
	return &Public{queries: queries}
}

// Home returns summary data for the homepage cards.
// @Summary     Home catalog
// @Description Returns homepage summaries for gaming/lounge/event
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/home [get]
func (p *Public) Home(c *gin.Context) {
	zones, err := p.queries.ListZonesWithServices(c.Request.Context())
	if err != nil {
		zap.L().Warn("home catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load home catalog", nil),
		).JSON(c)
		return
	}

	// Group zones by type for frontend cards
	gaming := make([]catalogZone, 0)
	lounge := make([]catalogZone, 0)
	event := make([]catalogZone, 0)

	for _, z := range zones {
		cz := catalogZone{
			ID:          z.ID,
			Name:        z.Name,
			ZoneType:    string(z.ZoneType),
			ZoneTagID:   z.ZoneTagID,
			Capacity:    z.Capacity,
			Description: z.Description.String,
			IsActive:    z.IsActive,
			DetailsJSON: parseJSON(z.DetailsJson),
		}
		switch z.ZoneType {
		case sqlc.ZoneTypeGame:
			gaming = append(gaming, cz)
		case sqlc.ZoneTypeLounge:
			lounge = append(lounge, cz)
		case sqlc.ZoneTypeEvent:
			event = append(event, cz)
		}
	}

	response.NewResponseBuilder(
		response.WithData("gaming", gaming),
		response.WithData("lounge", lounge),
		response.WithData("event", event),
	).JSON(c)
}

// Lounge returns the lounge catalog for the /lounge page.
// @Summary     Lounge catalog
// @Description Returns lounge zones with services and availability
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/lounge [get]
func (p *Public) Lounge(c *gin.Context) {
	zones, err := p.queries.ListZonesByType(c.Request.Context(), sqlc.ZoneTypeLounge)
	if err != nil {
		zap.L().Warn("lounge catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load lounge catalog", nil),
		).JSON(c)
		return
	}

	items := make([]catalogZoneWithServices, 0, len(zones))
	for _, z := range zones {
		services, err := p.queries.ListServicesByZoneID(c.Request.Context(), z.ID)
		if err != nil {
			zap.L().Warn("services lookup error", zap.Int64("zone_id", z.ID), zap.Error(err))
			services = []sqlc.Service{}
		}

		svcItems := make([]catalogService, 0, len(services))
		for _, s := range services {
			svcItems = append(svcItems, catalogService{
				ID:          s.ID,
				Name:        s.Name,
				Duration:    s.Duration,
				Price:       s.Price.Int.String(),
				Currency:    s.Currency,
				Description: s.Description.String,
				DetailsJSON: parseJSON(s.DetailsJson),
			})
		}

		items = append(items, catalogZoneWithServices{
			catalogZone: catalogZone{
				ID:          z.ID,
				Name:        z.Name,
				ZoneType:    string(z.ZoneType),
				ZoneTagID:   z.ZoneTagID,
				Capacity:    z.Capacity,
				Description: z.Description.String,
				IsActive:    z.IsActive,
				DetailsJSON: parseJSON(z.DetailsJson),
			},
			Services: svcItems,
		})
	}

	response.NewResponseBuilder(
		response.WithData("zones", items),
	).JSON(c)
}

// Event returns the event catalog for the /event page.
// @Summary     Event catalog
// @Description Returns event zones with services and availability
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/event [get]
func (p *Public) Event(c *gin.Context) {
	zones, err := p.queries.ListZonesByType(c.Request.Context(), sqlc.ZoneTypeEvent)
	if err != nil {
		zap.L().Warn("event catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load event catalog", nil),
		).JSON(c)
		return
	}

	items := make([]catalogZoneWithServices, 0, len(zones))
	for _, z := range zones {
		services, err := p.queries.ListServicesByZoneID(c.Request.Context(), z.ID)
		if err != nil {
			services = []sqlc.Service{}
		}

		svcItems := make([]catalogService, 0, len(services))
		for _, s := range services {
			svcItems = append(svcItems, catalogService{
				ID:          s.ID,
				Name:        s.Name,
				Duration:    s.Duration,
				Price:       s.Price.Int.String(),
				Currency:    s.Currency,
				Description: s.Description.String,
				DetailsJSON: parseJSON(s.DetailsJson),
			})
		}

		items = append(items, catalogZoneWithServices{
			catalogZone: catalogZone{
				ID:          z.ID,
				Name:        z.Name,
				ZoneType:    string(z.ZoneType),
				ZoneTagID:   z.ZoneTagID,
				Capacity:    z.Capacity,
				Description: z.Description.String,
				IsActive:    z.IsActive,
				DetailsJSON: parseJSON(z.DetailsJson),
			},
			Services: svcItems,
		})
	}

	response.NewResponseBuilder(
		response.WithData("zones", items),
	).JSON(c)
}

// Gaming returns the gaming catalog for the /gaming page.
// @Summary     Gaming catalog
// @Description Returns gaming zones with places, configurations, and availability
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/gaming [get]
func (p *Public) Gaming(c *gin.Context) {
	zoneTags, err := p.queries.ListUserZoneTags(c.Request.Context())
	if err != nil {
		zap.L().Warn("zone tags lookup error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load gaming tags", nil),
		).JSON(c)
		return
	}

	tagOrder := make(map[int32]int, len(zoneTags))
	tagIDs := make([]int64, 0, len(zoneTags))
	tagItems := make([]gamingZoneTag, 0, len(zoneTags))
	for index, tag := range zoneTags {
		tagOrder[tag.ID] = index
		tagIDs = append(tagIDs, int64(tag.ID))
		tagItems = append(tagItems, gamingZoneTag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	zones, err := p.queries.ListZonesByType(c.Request.Context(), sqlc.ZoneTypeGame)
	if err != nil {
		zap.L().Warn("gaming catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load gaming catalog", nil),
		).JSON(c)
		return
	}

	configs := []sqlc.ComputerConfiguration{}
	if len(tagIDs) > 0 {
		configs, err = p.queries.ListComputerConfigurationsByZoneTagIDs(c.Request.Context(), tagIDs)
		if err != nil {
			configs = []sqlc.ComputerConfiguration{}
		}
	}

	configItems := make([]gamingConfiguration, 0, len(configs))
	configMap := make(map[int64]json.RawMessage, len(configs))
	for _, cfg := range configs {
		parsedSpecs := parseJSONArray(cfg.SpecsJson)
		configMap[cfg.ID] = parsedSpecs
		configItems = append(configItems, gamingConfiguration{
			ID:        cfg.ID,
			ZoneTagID: cfg.ZoneTagsID,
			SpecsJSON: parsedSpecs,
		})
	}

	items := make([]gamingZone, 0, len(zones))
	for _, z := range zones {
		if _, ok := tagOrder[z.ZoneTagID]; !ok {
			continue
		}

		places, err := p.queries.ListActiveZonePlaces(c.Request.Context(), z.ID)
		if err != nil {
			places = []sqlc.ZonePlace{}
		}

		services, err := p.queries.ListServicesByZoneID(c.Request.Context(), z.ID)
		if err != nil {
			services = []sqlc.Service{}
		}

		placeItems := make([]gamingPlace, 0, len(places))
		for _, pl := range places {
			gp := gamingPlace{
				ID:        pl.ID,
				Label:     pl.Label,
				SortOrder: pl.SortOrder,
				IsActive:  pl.IsActive,
			}
			if pl.ConfigurationID.Valid {
				gp.ConfigurationID = &pl.ConfigurationID.Int64
				if specs, ok := configMap[pl.ConfigurationID.Int64]; ok {
					gp.Specs = specs
				}
			}
			placeItems = append(placeItems, gp)
		}

		svcItems := make([]catalogService, 0, len(services))
		for _, s := range services {
			svcItems = append(svcItems, catalogService{
				ID:          s.ID,
				Name:        s.Name,
				Duration:    s.Duration,
				Price:       s.Price.Int.String(),
				Currency:    s.Currency,
				Description: s.Description.String,
				DetailsJSON: parseJSON(s.DetailsJson),
			})
		}

		items = append(items, gamingZone{
			catalogZone: catalogZone{
				ID:          z.ID,
				Name:        z.Name,
				ZoneType:    string(z.ZoneType),
				ZoneTagID:   z.ZoneTagID,
				Capacity:    z.Capacity,
				Description: z.Description.String,
				IsActive:    z.IsActive,
				DetailsJSON: parseJSON(z.DetailsJson),
			},
			Places:   placeItems,
			Services: svcItems,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		leftOrder := tagOrder[items[i].ZoneTagID]
		rightOrder := tagOrder[items[j].ZoneTagID]
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}

		return items[i].ID < items[j].ID
	})

	response.NewResponseBuilder(
		response.WithData("zone_tags", tagItems),
		response.WithData("zones", items),
		response.WithData("configurations", configItems),
	).JSON(c)
}

// --- DTO types for public catalog responses ---

type catalogZone struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	ZoneType    string          `json:"zone_type"`
	ZoneTagID   int32           `json:"zone_tag_id"`
	Capacity    int32           `json:"capacity"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	DetailsJSON json.RawMessage `json:"details_json"`
}

type catalogService struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Duration    int32           `json:"duration"`
	Price       string          `json:"price"`
	Currency    string          `json:"currency"`
	Description string          `json:"description"`
	DetailsJSON json.RawMessage `json:"details_json"`
}

type catalogZoneWithServices struct {
	catalogZone
	Services []catalogService `json:"services"`
}

type gamingZoneTag struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type gamingConfiguration struct {
	ID        int64           `json:"id"`
	ZoneTagID int64           `json:"zone_tag_id"`
	SpecsJSON json.RawMessage `json:"specs_json"`
}

type gamingPlace struct {
	ID              int64           `json:"id"`
	Label           string          `json:"label"`
	ConfigurationID *int64          `json:"configuration_id,omitempty"`
	SortOrder       int32           `json:"sort_order"`
	IsActive        bool            `json:"is_active"`
	Specs           json.RawMessage `json:"specs,omitempty"`
}

type gamingZone struct {
	catalogZone
	Places   []gamingPlace    `json:"places"`
	Services []catalogService `json:"services"`
}

func parseJSON(b []byte) json.RawMessage {
	if len(b) == 0 {
		return json.RawMessage("{}")
	}
	return json.RawMessage(b)
}

func parseJSONArray(b []byte) json.RawMessage {
	if len(b) == 0 {
		return json.RawMessage("[]")
	}
	return json.RawMessage(b)
}
