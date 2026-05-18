package service

import (
	"context"

	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"backend/pkg"

	"github.com/jackc/pgx/v5/pgtype"
)

type ZoneService struct {
	queries *sqlc.Queries
}

func NewZone(queries *sqlc.Queries) *ZoneService {
	return &ZoneService{queries: queries}
}

func (z *ZoneService) CreateZone(ctx context.Context, dto domain.CreateZoneRequest) (domain.Zone, error) {
	isActive := true
	if dto.IsActive != nil {
		isActive = *dto.IsActive
	}
	detailsJSON := []byte("{}")
	if dto.DetailsJSON != "" {
		detailsJSON = []byte(dto.DetailsJSON)
	}

	var zoneTagID pgtype.Int4
	if dto.ZoneTagID != nil {
		zoneTagID = pgtype.Int4{Int32: *dto.ZoneTagID, Valid: true}
	}

	zone, err := z.queries.CreateZone(
		ctx,
		sqlc.CreateZoneParams{
			Name:        dto.Name,
			ZoneType:    dto.Type,
			ZoneTagID:   zoneTagID,
			Capacity:    int32(dto.Capacity),
			Description: pgtype.Text{
				String: dto.Description,
				Valid:  dto.Description != "",
			},
			IsActive:    isActive,
			DetailsJson: detailsJSON,
		},
	)
	if err != nil {
		return domain.Zone{}, err
	}
	return mapZoneToDomain(zone), nil
}

func (z *ZoneService) GetZones(ctx context.Context) ([]domain.Zone, error) {
	zones, err := z.queries.ListZones(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Zone, 0, len(zones))
	for _, zone := range zones {
		result = append(result, mapZoneToDomain(zone))
	}
	return result, nil
}

func (z *ZoneService) GetZoneByID(ctx context.Context, id int64) (domain.Zone, error) {
	zone, err := z.queries.GetZoneByID(ctx, id)
	if err != nil {
		return domain.Zone{}, err
	}
	return mapZoneToDomain(zone), nil
}

func (z *ZoneService) DeleteByID(ctx context.Context, id int64) (int64, error) {
	return z.queries.DeleteZone(ctx, id)
}

func (z *ZoneService) PatchByID(ctx context.Context, id int64, dto domain.PatchZoneRequest) (domain.Zone, error) {
	params := sqlc.PatchZoneParams{ID: id}
	if dto.Name != nil {
		params.Name = pgtype.Text{
			String: *dto.Name,
			Valid:  true,
		}
	}
	if dto.Type != nil {
		params.ZoneType = sqlc.NullZoneType{
			ZoneType: *dto.Type,
			Valid:    true,
		}
	}
	if dto.ZoneTagID != nil {
		params.ZoneTagID = pgtype.Int4{
			Int32: *dto.ZoneTagID,
			Valid: true,
		}
	}
	if dto.Capacity != nil {
		params.Capacity = pgtype.Int4{
			Int32: int32(*dto.Capacity),
			Valid: true,
		}
	}
	if dto.Description != nil {
		params.Description = pgtype.Text{
			String: *dto.Description,
			Valid:  true,
		}
	}
	if dto.IsActive != nil {
		params.IsActive = pgtype.Bool{
			Bool:  *dto.IsActive,
			Valid: true,
		}
	}
	if dto.DetailsJSON != nil {
		params.DetailsJson = []byte(*dto.DetailsJSON)
	}

	zone, err := z.queries.PatchZone(ctx, params)
	if err != nil {
		return domain.Zone{}, err
	}
	return mapZoneToDomain(zone), nil
}

func mapZoneToDomain(z sqlc.Zone) domain.Zone {
	var zoneTagID *int32
	if z.ZoneTagID.Valid {
		val := z.ZoneTagID.Int32
		zoneTagID = &val
	}

	return domain.Zone{
		ID:          z.ID,
		Name:        z.Name,
		ZoneType:    z.ZoneType,
		ZoneTagID:   zoneTagID,
		Capacity:    z.Capacity,
		Description: z.Description.String,
		IsActive:    z.IsActive,
		DetailsJSON: pkg.JSONOrObject(z.DetailsJson),
		CreatedAt:   z.CreatedAt.Time,
		UpdatedAt:   z.UpdatedAt.Time,
	}
}
