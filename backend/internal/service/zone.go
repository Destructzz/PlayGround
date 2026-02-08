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
	return z.queries.CreateZone(
		ctx,
		sqlc.CreateZoneParams{
			Name:     dto.Name,
			ZoneType: dto.Type,
			Capacity: int32(dto.Capacity),
			Description: pgtype.Text{
				String: dto.Description,
				Valid:  dto.Description != "",
			},
		},
	)
}

func (z *ZoneService) GetZones(ctx context.Context) ([]sqlc.Zone, error) {
	return z.queries.ListZones(ctx)
}

func (z *ZoneService) GetZoneByID(ctx context.Context, id int64) (sqlc.Zone, error) {
	return z.queries.GetZoneByID(ctx, id)
}
