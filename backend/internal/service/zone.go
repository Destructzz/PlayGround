package service

import (
	"context"

	"backend/internal/domain"
	"backend/internal/repo/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type ZoneService struct {
	queries *sqlc.Queries
}

func NewZone(queries *sqlc.Queries) *ZoneService {
	return &ZoneService{queries: queries}
}

func (z *ZoneService) CreateZone(ctx context.Context, dto domain.CreateZoneRequest) (sqlc.Zone, error) {
	isActive := true
	if dto.IsActive != nil {
		isActive = *dto.IsActive
	}
	detailsJSON := []byte("{}")
	if dto.DetailsJSON != "" {
		detailsJSON = []byte(dto.DetailsJSON)
	}

	return z.queries.CreateZone(
		ctx,
		sqlc.CreateZoneParams{
			Name:     dto.Name,
			ZoneType: dto.Type,
			ZoneTagID: dto.ZoneTagID,
			Capacity: int32(dto.Capacity),
			Description: pgtype.Text{
				String: dto.Description,
				Valid:  dto.Description != "",
			},
			IsActive:    isActive,
			DetailsJson: detailsJSON,
		},
	)
}

func (z *ZoneService) GetZones(ctx context.Context) ([]sqlc.Zone, error) {
	return z.queries.ListZones(ctx)
}

func (z *ZoneService) GetZoneByID(ctx context.Context, id int64) (sqlc.Zone, error) {
	return z.queries.GetZoneByID(ctx, id)
}

func (z *ZoneService) DeleteByID(ctx context.Context, id int64) (int64, error) {
	return z.queries.DeleteZone(ctx, id)
}

func (z *ZoneService) PatchByID(ctx context.Context, id int64, dto domain.PatchZoneRequest) (sqlc.Zone, error) {
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

	return z.queries.PatchZone(ctx, params)
}
