package service

import (
	"backend/internal/repo/sqlc"
	"backend/pkg"
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type PublicService struct {
	queries *sqlc.Queries
}

func NewPublicService(queries *sqlc.Queries) *PublicService {
	return &PublicService{queries: queries}
}

type PublicCatalogZone struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	ZoneType    string          `json:"zone_type"`
	ZoneTagID   int32           `json:"zone_tag_id"`
	Capacity    int32           `json:"capacity"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	DetailsJSON json.RawMessage `json:"details_json"`
}

type PublicCatalogService struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Duration    int32           `json:"duration"`
	Price       string          `json:"price"`
	Currency    string          `json:"currency"`
	Description string          `json:"description"`
	DetailsJSON json.RawMessage `json:"details_json"`
}

type PublicCatalogZoneWithServices struct {
	PublicCatalogZone
	Services []PublicCatalogService `json:"services"`
}

type HomeCatalog struct {
	Gaming []PublicCatalogZone `json:"gaming"`
	Lounge []PublicCatalogZone `json:"lounge"`
	Event  []PublicCatalogZone `json:"event"`
}

type GamingZoneTag struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type GamingConfiguration struct {
	ID        int64           `json:"id"`
	ZoneTagID int64           `json:"zone_tag_id"`
	SpecsJSON json.RawMessage `json:"specs_json"`
}

type GamingPlace struct {
	ID              int64           `json:"id"`
	Label           string          `json:"label"`
	ConfigurationID *int64          `json:"configuration_id,omitempty"`
	SortOrder       int32           `json:"sort_order"`
	IsActive        bool            `json:"is_active"`
	Specs           json.RawMessage `json:"specs,omitempty"`
}

type GamingZone struct {
	PublicCatalogZone
	Places   []GamingPlace          `json:"places"`
	Services []PublicCatalogService `json:"services"`
}

type GamingCatalog struct {
	ZoneTags       []GamingZoneTag       `json:"zone_tags"`
	Zones          []GamingZone          `json:"zones"`
	Configurations []GamingConfiguration `json:"configurations"`
}

type GamingAvailabilityBooking struct {
	BookingID int64     `json:"booking_id"`
	PlaceID   int64     `json:"place_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

type GamingAvailability struct {
	ZoneID   int64                       `json:"zone_id"`
	Date     string                      `json:"date"`
	Bookings []GamingAvailabilityBooking `json:"bookings"`
}

type LoungeSlot struct {
	Hour              int    `json:"hour"`
	Label             string `json:"label"`
	BookedParticipants int32  `json:"booked_participants"`
	Remaining         int32  `json:"remaining"`
	Available         bool   `json:"available"`
}

type LoungeAvailability struct {
	ZoneID   int64        `json:"zone_id"`
	Capacity int32        `json:"capacity"`
	Date     string       `json:"date"`
	Slots    []LoungeSlot `json:"slots"`
}

func (s *PublicService) HomeCatalog(ctx context.Context) (HomeCatalog, error) {
	zones, err := s.queries.ListZonesWithServices(ctx)
	if err != nil {
		return HomeCatalog{}, err
	}

	result := HomeCatalog{
		Gaming: make([]PublicCatalogZone, 0),
		Lounge: make([]PublicCatalogZone, 0),
		Event:  make([]PublicCatalogZone, 0),
	}

	for _, zone := range zones {
		mapped := mapCatalogZone(zone.ID, zone.Name, string(zone.ZoneType), zone.ZoneTagID, zone.Capacity, zone.Description, zone.IsActive, zone.DetailsJson)

		switch zone.ZoneType {
		case sqlc.ZoneTypeGame:
			result.Gaming = append(result.Gaming, mapped)
		case sqlc.ZoneTypeLounge:
			result.Lounge = append(result.Lounge, mapped)
		case sqlc.ZoneTypeEvent:
			result.Event = append(result.Event, mapped)
		}
	}

	return result, nil
}

func (s *PublicService) LoungeCatalog(ctx context.Context) ([]PublicCatalogZoneWithServices, error) {
	return s.catalogByZoneType(ctx, sqlc.ZoneTypeLounge)
}

func (s *PublicService) EventCatalog(ctx context.Context) ([]PublicCatalogZoneWithServices, error) {
	return s.catalogByZoneType(ctx, sqlc.ZoneTypeEvent)
}

func (s *PublicService) GamingCatalog(ctx context.Context) (GamingCatalog, error) {
	var displayedTagIDs []int32
	settings, err := s.queries.GetSiteSettings(ctx)
	if err == nil && len(settings.SettingsJson) > 0 {
		_ = json.Unmarshal(settings.SettingsJson, &displayedTagIDs)
	}

	allTags, err := s.queries.ListZoneTags(ctx)
	if err != nil {
		return GamingCatalog{}, err
	}

	var zoneTags []sqlc.ZoneTag
	if len(displayedTagIDs) > 0 {
		tagMap := make(map[int32]sqlc.ZoneTag, len(allTags))
		for _, tag := range allTags {
			tagMap[int32(tag.ID)] = tag
		}
		for _, id := range displayedTagIDs {
			if tag, ok := tagMap[id]; ok {
				zoneTags = append(zoneTags, tag)
			}
		}
	} else {
		for i := 0; i < len(allTags) && i < 3; i++ {
			zoneTags = append(zoneTags, allTags[i])
		}
	}

	tagOrder := make(map[int32]int, len(zoneTags))
	tagIDs := make([]int64, 0, len(zoneTags))
	tagItems := make([]GamingZoneTag, 0, len(zoneTags))
	for index, tag := range zoneTags {
		tagOrder[tag.ID] = index
		tagIDs = append(tagIDs, int64(tag.ID))
		tagItems = append(tagItems, GamingZoneTag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	zones, err := s.queries.ListZonesByType(ctx, sqlc.ZoneTypeGame)
	if err != nil {
		return GamingCatalog{}, err
	}

	configs := make([]sqlc.ComputerConfiguration, 0)
	if len(tagIDs) > 0 {
		configs, err = s.queries.ListComputerConfigurationsByZoneTagIDs(ctx, tagIDs)
		if err != nil {
			configs = make([]sqlc.ComputerConfiguration, 0)
		}
	}

	configItems := make([]GamingConfiguration, 0, len(configs))
	configMap := make(map[int64]json.RawMessage, len(configs))
	for _, cfg := range configs {
		specs := pkg.JSONOrArray(cfg.SpecsJson)
		configMap[cfg.ID] = specs
		configItems = append(configItems, GamingConfiguration{
			ID:        cfg.ID,
			ZoneTagID: cfg.ZoneTagsID,
			SpecsJSON: specs,
		})
	}

	items := make([]GamingZone, 0, len(zones))
	for _, zone := range zones {
		if !zone.ZoneTagID.Valid {
			continue
		}
		if _, ok := tagOrder[zone.ZoneTagID.Int32]; !ok {
			continue
		}

		places, err := s.queries.ListActiveZonePlaces(ctx, zone.ID)
		if err != nil {
			places = make([]sqlc.ZonePlace, 0)
		}

		services, err := s.queries.ListServicesByZoneID(ctx, zone.ID)
		if err != nil {
			services = make([]sqlc.Service, 0)
		}

		placeItems := make([]GamingPlace, 0, len(places))
		for _, place := range places {
			item := GamingPlace{
				ID:        place.ID,
				Label:     place.Label,
				SortOrder: place.SortOrder,
				IsActive:  place.IsActive,
			}

			if place.ConfigurationID.Valid {
				item.ConfigurationID = &place.ConfigurationID.Int64
				if specs, ok := configMap[place.ConfigurationID.Int64]; ok {
					item.Specs = specs
				}
			}

			placeItems = append(placeItems, item)
		}

		items = append(items, GamingZone{
			PublicCatalogZone: mapCatalogZone(zone.ID, zone.Name, string(zone.ZoneType), zone.ZoneTagID, zone.Capacity, zone.Description, zone.IsActive, zone.DetailsJson),
			Places:            placeItems,
			Services:          mapCatalogServices(services),
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

	return GamingCatalog{
		ZoneTags:       tagItems,
		Zones:          items,
		Configurations: configItems,
	}, nil
}

func (s *PublicService) GamingAvailability(ctx context.Context, zoneID int64, selectedDate time.Time) (GamingAvailability, error) {
	dateStart := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.Add(24 * time.Hour)

	rows, err := s.queries.ListZonePlaceBookingsForDate(ctx, sqlc.ListZonePlaceBookingsForDateParams{
		ZoneID:    zoneID,
		DateStart: pgtype.Timestamptz{Time: dateStart, Valid: true},
		DateEnd:   pgtype.Timestamptz{Time: dateEnd, Valid: true},
	})
	if err != nil {
		return GamingAvailability{}, err
	}

	items := make([]GamingAvailabilityBooking, 0, len(rows))
	for _, row := range rows {
		if !row.PlaceID.Valid {
			continue
		}

		items = append(items, GamingAvailabilityBooking{
			BookingID: row.ID,
			PlaceID:   row.PlaceID.Int64,
			StartTime: row.StartTime.Time,
			EndTime:   row.EndTime.Time,
			Status:    string(row.Status),
		})
	}

	return GamingAvailability{
		ZoneID:   zoneID,
		Date:     dateStart.Format("2006-01-02"),
		Bookings: items,
	}, nil
}

func (s *PublicService) LoungeAvailability(ctx context.Context, zoneID int64, selectedDate time.Time) (LoungeAvailability, error) {
	zone, err := s.queries.GetZoneByID(ctx, zoneID)
	if err != nil {
		return LoungeAvailability{}, err
	}

	// Working hours: 10:00–23:00 Moscow time (UTC+3)
	loc := time.FixedZone("MSK", 3*3600)
	dayStart := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), 0, 0, 0, 0, loc)
	dayEnd := dayStart.Add(24 * time.Hour)

	rows, err := s.queries.ListLoungeBookingsForDate(ctx, sqlc.ListLoungeBookingsForDateParams{
		ZoneID:  zoneID,
		DateEnd: pgtype.Timestamptz{Time: dayEnd, Valid: true},
		DateStart: pgtype.Timestamptz{Time: dayStart, Valid: true},
	})
	if err != nil {
		return LoungeAvailability{}, err
	}

	const openHour = 10
	const closeHour = 23
	slots := make([]LoungeSlot, 0, closeHour-openHour)

	for h := openHour; h < closeHour; h++ {
		slotStart := time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), h, 0, 0, 0, loc)
		slotEnd := slotStart.Add(time.Hour)

		var booked int32
		for _, row := range rows {
			rowStart := row.StartTime.Time
			rowEnd := row.EndTime.Time
			// Overlap: row starts before slot ends AND row ends after slot starts
			if rowStart.Before(slotEnd) && rowEnd.After(slotStart) {
				booked += row.Participants
			}
		}

		remaining := zone.Capacity - booked
		if remaining < 0 {
			remaining = 0
		}

		slots = append(slots, LoungeSlot{
			Hour:              h,
			Label:             slotStart.Format("15:04") + "–" + slotEnd.Format("15:04"),
			BookedParticipants: booked,
			Remaining:         remaining,
			Available:         remaining > 0,
		})
	}

	return LoungeAvailability{
		ZoneID:   zoneID,
		Capacity: zone.Capacity,
		Date:     dayStart.Format("2006-01-02"),
		Slots:    slots,
	}, nil
}

func (s *PublicService) catalogByZoneType(ctx context.Context, zoneType sqlc.ZoneType) ([]PublicCatalogZoneWithServices, error) {
	zones, err := s.queries.ListZonesByType(ctx, zoneType)
	if err != nil {
		return nil, err
	}

	items := make([]PublicCatalogZoneWithServices, 0, len(zones))
	for _, zone := range zones {
		services, err := s.queries.ListServicesByZoneID(ctx, zone.ID)
		if err != nil {
			services = make([]sqlc.Service, 0)
		}

		items = append(items, PublicCatalogZoneWithServices{
			PublicCatalogZone: mapCatalogZone(zone.ID, zone.Name, string(zone.ZoneType), zone.ZoneTagID, zone.Capacity, zone.Description, zone.IsActive, zone.DetailsJson),
			Services:          mapCatalogServices(services),
		})
	}

	return items, nil
}

func mapCatalogZone(id int64, name string, zoneType string, zoneTagID pgtype.Int4, capacity int32, description pgtype.Text, isActive bool, detailsJSON []byte) PublicCatalogZone {
	var tagID int32
	if zoneTagID.Valid {
		tagID = zoneTagID.Int32
	}
	return PublicCatalogZone{
		ID:          id,
		Name:        name,
		ZoneType:    zoneType,
		ZoneTagID:   tagID,
		Capacity:    capacity,
		Description: pkg.TextValue(description),
		IsActive:    isActive,
		DetailsJSON: pkg.JSONOrObject(detailsJSON),
	}
}

func mapCatalogServices(services []sqlc.Service) []PublicCatalogService {
	items := make([]PublicCatalogService, 0, len(services))
	for _, item := range services {
		items = append(items, PublicCatalogService{
			ID:          item.ID,
			Name:        item.Name,
			Duration:    item.Duration,
			Price:       pkg.NumericToString(item.Price),
			Currency:    item.Currency,
			Description: pkg.TextValue(item.Description),
			DetailsJSON: pkg.JSONOrObject(item.DetailsJson),
		})
	}

	return items
}
