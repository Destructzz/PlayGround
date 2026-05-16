package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const gamingSeedKey = "gaming_base_v1"

type SeedService struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

type SeedSnapshot struct {
	SeedKey  string         `json:"seed_key"`
	Seeded   bool           `json:"seeded"`
	ZoneTags []sqlc.ZoneTag `json:"zone_tags"`
	Configurations []domain.ComputerConfiguration `json:"configurations"`

	Zones []SeedZoneSnapshot `json:"zones"`
}

type SeedZoneSnapshot struct {
	Zone     sqlc.Zone        `json:"zone"`
	Places   []sqlc.ZonePlace `json:"places"`
	Services []sqlc.Service   `json:"services"`
}

type seedTagDefinition struct {
	Name        string
	ZoneName    string
	Capacity    int32
	Description string
	ZoneDetails map[string]any
	ServiceName string
	Price       string
	Duration    int32
	ServiceDesc string
	Places      []string
	Specs       []map[string]string
}

func NewSeed(pool *pgxpool.Pool, queries *sqlc.Queries) *SeedService {
	return &SeedService{pool: pool, queries: queries}
}

func (s *SeedService) Get(ctx context.Context) (SeedSnapshot, error) {
	tags, err := s.queries.ListZoneTags(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}

	seedTags := make([]sqlc.ZoneTag, 0, len(seedDefinitions))
	tagByName := make(map[string]sqlc.ZoneTag, len(tags))
	for _, tag := range tags {
		tagByName[tag.Name] = tag
	}
	for _, definition := range seedDefinitions {
		if tag, ok := tagByName[definition.Name]; ok {
			seedTags = append(seedTags, tag)
		}
	}

	configs, err := s.queries.ListComputerConfigurations(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}

	zones, err := s.queries.ListZonesByType(ctx, sqlc.ZoneTypeGame)
	if err != nil {
		return SeedSnapshot{}, err
	}

	zoneByName := make(map[string]sqlc.Zone, len(zones))
	for _, zone := range zones {
		zoneByName[zone.Name] = zone
	}

	seedZoneSnapshots := make([]SeedZoneSnapshot, 0, len(seedDefinitions))
	seedConfigurations := make([]domain.ComputerConfiguration, 0, len(seedDefinitions))

	for _, definition := range seedDefinitions {
		specsJSON, err := json.Marshal(definition.Specs)
		if err != nil {
			return SeedSnapshot{}, err
		}

		tag, ok := tagByName[definition.Name]
		if !ok {
			continue
		}

		for _, cfg := range configs {
			if cfg.ZoneTagsID == int64(tag.ID) && specsEqual(cfg.SpecsJson, specsJSON) {
				seedConfigurations = append(seedConfigurations, domain.ComputerConfiguration{
					ID:        cfg.ID,
					ZoneTagID: cfg.ZoneTagsID,
					SpecsJSON: json.RawMessage(cfg.SpecsJson),
				})
				break
			}
		}

		zone, ok := zoneByName[definition.ZoneName]
		if !ok || !isSeedZone(zone.DetailsJson) {
			continue
		}

		places, err := s.queries.ListZonePlaces(ctx, zone.ID)
		if err != nil {
			return SeedSnapshot{}, err
		}

		services, err := s.queries.ListServicesByZoneID(ctx, zone.ID)
		if err != nil {
			return SeedSnapshot{}, err
		}

		seedZoneSnapshots = append(seedZoneSnapshots, SeedZoneSnapshot{
			Zone:     zone,
			Places:   places,
			Services: services,
		})
	}

	return SeedSnapshot{
		SeedKey:        gamingSeedKey,
		Seeded:         len(seedZoneSnapshots) == len(seedDefinitions),
		ZoneTags:       seedTags,
		Configurations: seedConfigurations,
		Zones:          seedZoneSnapshots,
	}, nil
}

func (s *SeedService) Apply(ctx context.Context) (SeedSnapshot, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	tags, err := q.ListZoneTags(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}

	tagByName := make(map[string]sqlc.ZoneTag, len(tags))
	for _, tag := range tags {
		tagByName[tag.Name] = tag
	}

	configs, err := q.ListComputerConfigurations(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}

	zones, err := q.ListZones(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}

	zoneByName := make(map[string]sqlc.Zone, len(zones))
	for _, zone := range zones {
		zoneByName[zone.Name] = zone
	}

	for _, definition := range seedDefinitions {
		tag, ok := tagByName[definition.Name]
		if !ok {
			tag, err = q.CreateZoneTag(ctx, definition.Name)
			if err != nil {
				return SeedSnapshot{}, err
			}
			tagByName[definition.Name] = tag
		}

		configSpecs, err := json.Marshal(definition.Specs)
		if err != nil {
			return SeedSnapshot{}, err
		}

		zoneDetails, err := json.Marshal(definition.ZoneDetails)
		if err != nil {
			return SeedSnapshot{}, err
		}

		configuration, found := findSeedConfiguration(configs, int64(tag.ID), configSpecs)
		if !found {
			configuration, err = q.CreateComputerConfiguration(ctx, sqlc.CreateComputerConfigurationParams{
				ZoneTagsID: int64(tag.ID),
				SpecsJson:  configSpecs,
			})
			if err != nil {
				return SeedSnapshot{}, err
			}
			configs = append(configs, configuration)
		}

		zone, ok := zoneByName[definition.ZoneName]
		if ok {
			zone, err = q.PatchZone(ctx, sqlc.PatchZoneParams{
				ID: zone.ID,
				Name: pgtype.Text{
					String: definition.ZoneName,
					Valid:  true,
				},
				ZoneType: sqlc.NullZoneType{
					ZoneType: sqlc.ZoneTypeGame,
					Valid:    true,
				},
				ZoneTagID: pgtype.Int4{
					Int32: tag.ID,
					Valid: true,
				},
				Capacity: pgtype.Int4{
					Int32: definition.Capacity,
					Valid: true,
				},
				Description: pgtype.Text{
					String: definition.Description,
					Valid:  definition.Description != "",
				},
				IsActive: pgtype.Bool{
					Bool:  true,
					Valid: true,
				},
				DetailsJson: zoneDetails,
			})
			if err != nil {
				return SeedSnapshot{}, err
			}
			zoneByName[definition.ZoneName] = zone
		} else {
			zone, err = q.CreateZone(ctx, sqlc.CreateZoneParams{
				Name:      definition.ZoneName,
				ZoneType:  sqlc.ZoneTypeGame,
				ZoneTagID: tag.ID,
				Capacity:  definition.Capacity,
				Description: pgtype.Text{
					String: definition.Description,
					Valid:  definition.Description != "",
				},
				IsActive:    true,
				DetailsJson: zoneDetails,
			})
			if err != nil {
				return SeedSnapshot{}, err
			}
			zoneByName[definition.ZoneName] = zone
		}

		services, err := q.ListServicesByZoneID(ctx, zone.ID)
		if err != nil {
			return SeedSnapshot{}, err
		}

		if _, ok := findSeedService(services, definition); !ok {
			price := pgtype.Numeric{}
			if err := price.Scan(definition.Price); err != nil {
				return SeedSnapshot{}, err
			}

			serviceDetails, err := json.Marshal(map[string]any{"seed_key": gamingSeedKey})
			if err != nil {
				return SeedSnapshot{}, err
			}

			if _, err := q.CreateService(ctx, sqlc.CreateServiceParams{
				Name:     definition.ServiceName,
				ZoneID:   zone.ID,
				Duration: definition.Duration,
				Price:    price,
				Currency: "RUB",
				Description: pgtype.Text{
					String: definition.ServiceDesc,
					Valid:  definition.ServiceDesc != "",
				},
				IsActive:    true,
				DetailsJson: serviceDetails,
			}); err != nil {
				return SeedSnapshot{}, err
			}
		}

		places, err := q.ListZonePlaces(ctx, zone.ID)
		if err != nil {
			return SeedSnapshot{}, err
		}

		placeByLabel := make(map[string]sqlc.ZonePlace, len(places))
		for _, place := range places {
			placeByLabel[place.Label] = place
		}

		for index, label := range definition.Places {
			if _, ok := placeByLabel[label]; ok {
				continue
			}

			if _, err := q.CreateZonePlace(ctx, sqlc.CreateZonePlaceParams{
				ZoneID: zone.ID,
				Label:  label,
				ConfigurationID: pgtype.Int8{
					Int64: configuration.ID,
					Valid: true,
				},
				SortOrder: int32(index + 1),
				IsActive:  true,
			}); err != nil {
				return SeedSnapshot{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return SeedSnapshot{}, err
	}

	return s.Get(ctx)
}

func (s *SeedService) Delete(ctx context.Context) (SeedSnapshot, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)
	snapshot, err := s.Get(ctx)
	if err != nil {
		return SeedSnapshot{}, err
	}

	for _, zoneSnapshot := range snapshot.Zones {
		for _, place := range zoneSnapshot.Places {
			if _, err := q.DeleteZonePlace(ctx, place.ID); err != nil {
				return SeedSnapshot{}, err
			}
		}

		for _, svc := range zoneSnapshot.Services {
			if _, err := q.DeleteService(ctx, svc.ID); err != nil {
				return SeedSnapshot{}, err
			}
		}

		if _, err := q.DeleteZone(ctx, zoneSnapshot.Zone.ID); err != nil {
			return SeedSnapshot{}, err
		}
	}

	for _, cfg := range snapshot.Configurations {
		if _, err := q.DeleteComputerConfiguration(ctx, cfg.ID); err != nil {
			return SeedSnapshot{}, err
		}
	}

	for _, tag := range snapshot.ZoneTags {
		if _, err := q.DeleteZoneTag(ctx, tag.ID); err != nil {
			return SeedSnapshot{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return SeedSnapshot{}, err
	}

	return SeedSnapshot{
		SeedKey: gamingSeedKey,
		Seeded:  false,
	}, nil
}

func findSeedConfiguration(configs []sqlc.ComputerConfiguration, zoneTagID int64, specs []byte) (sqlc.ComputerConfiguration, bool) {
	for _, cfg := range configs {
		if cfg.ZoneTagsID == zoneTagID && specsEqual(cfg.SpecsJson, specs) {
			return cfg, true
		}
	}

	return sqlc.ComputerConfiguration{}, false
}

func findSeedService(services []sqlc.Service, definition seedTagDefinition) (sqlc.Service, bool) {
	for _, service := range services {
		if service.Name == definition.ServiceName && service.Duration == definition.Duration && service.Description.String == definition.ServiceDesc {
			return service, true
		}
	}

	return sqlc.Service{}, false
}

func isSeedZone(details []byte) bool {
	var payload map[string]any
	if err := json.Unmarshal(details, &payload); err != nil {
		return false
	}

	seedKey, _ := payload["seed_key"].(string)
	return seedKey == gamingSeedKey
}

func specsEqual(left []byte, right []byte) bool {
	var leftValue any
	var rightValue any

	if err := json.Unmarshal(left, &leftValue); err != nil {
		return false
	}
	if err := json.Unmarshal(right, &rightValue); err != nil {
		return false
	}

	leftJSON, _ := json.Marshal(leftValue)
	rightJSON, _ := json.Marshal(rightValue)
	return bytes.Equal(leftJSON, rightJSON)
}

var seedDefinitions = []seedTagDefinition{
	{
		Name:        "Standard",
		ZoneName:    "Standard Arena",
		Capacity:    10,
		Description: "Надежные сборки для комфортной игры в любые современные проекты.",
		ZoneDetails: map[string]any{"prefix": "БАЗОВАЯ", "title": "МОЩЬ", "hexColor": "#10b981", "seed_key": gamingSeedKey},
		ServiceName: "Standard",
		Price:       "150",
		Duration:    60,
		ServiceDesc: "RTX 4060 Ti • 165Hz",
		Places:      buildLabels("ПК ", 10),
		Specs: []map[string]string{
			{"title": "Видеокарта", "value": "RTX 4060 Ti"},
			{"title": "Процессор", "value": "Intel Core i5 13400F"},
			{"title": "Оперативная память", "value": "16GB DDR4"},
			{"title": "Монитор", "value": "AOC 165Hz IPS"},
			{"title": "Периферия", "value": "HyperX Series"},
		},
	},
	{
		Name:        "VIP Class",
		ZoneName:    "VIP Arena",
		Capacity:    6,
		Description: "Премиальная зона для киберспортивной игры на максималках.",
		ZoneDetails: map[string]any{"prefix": "БЕСКОМПРОМИССНАЯ", "title": "МОЩЬ", "hexColor": "#22d3ee", "seed_key": gamingSeedKey},
		ServiceName: "VIP",
		Price:       "250",
		Duration:    60,
		ServiceDesc: "RTX 4090 • 360Hz",
		Places:      buildLabels("VIP ", 6),
		Specs: []map[string]string{
			{"title": "Видеокарта", "value": "RTX 4090 24GB"},
			{"title": "Процессор", "value": "Intel Core i9 14900K"},
			{"title": "Оперативная память", "value": "64GB DDR5 XMP"},
			{"title": "Монитор", "value": "Zowie 360Hz 1ms"},
			{"title": "Периферия", "value": "Logitech G PRO X"},
		},
	},
	{
		Name:        "Bootcamp",
		ZoneName:    "Bootcamp Room",
		Capacity:    5,
		Description: "Изолированная тренировочная комната для команды.",
		ZoneDetails: map[string]any{"prefix": "КОМАНДНАЯ", "title": "СИНЕРГИЯ", "hexColor": "#f97316", "seed_key": gamingSeedKey},
		ServiceName: "Bootcamp",
		Price:       "200",
		Duration:    60,
		ServiceDesc: "На человека • от 5 чел.",
		Places:      buildLabels("BC ", 5),
		Specs: []map[string]string{
			{"title": "Видеокарта", "value": "RTX 4080 Super"},
			{"title": "Процессор", "value": "AMD Ryzen 7 7800X3D"},
			{"title": "Оперативная память", "value": "32GB DDR5 6000MHz"},
			{"title": "Монитор", "value": "Alienware 240Hz OLED"},
			{"title": "Периферия", "value": "Razer Esports"},
		},
	},
}

func buildLabels(prefix string, count int) []string {
	labels := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		labels = append(labels, prefix+strconv.Itoa(i))
	}
	return labels
}
