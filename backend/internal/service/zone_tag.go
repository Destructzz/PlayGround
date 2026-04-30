package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type ZoneTagService struct {
	queries *sqlc.Queries
}

func NewZoneTag(queries *sqlc.Queries) *ZoneTagService {
	return &ZoneTagService{queries: queries}
}

func (s *ZoneTagService) Create(ctx context.Context, dto domain.CreateZoneTagRequest) (sqlc.ZoneTag, error) {
	return s.queries.CreateZoneTag(ctx, dto.Name)
}

func (s *ZoneTagService) List(ctx context.Context) ([]sqlc.ZoneTag, error) {
	return s.queries.ListZoneTags(ctx)
}

func (s *ZoneTagService) GetByID(ctx context.Context, id int32) (sqlc.ZoneTag, error) {
	return s.queries.GetZoneTagByID(ctx, id)
}

func (s *ZoneTagService) Delete(ctx context.Context, id int32) (int64, error) {
	return s.queries.DeleteZoneTag(ctx, id)
}

func (s *ZoneTagService) Patch(ctx context.Context, id int32, dto domain.PatchZoneTagRequest) (sqlc.ZoneTag, error) {
	params := sqlc.PatchZoneTagParams{ID: id}

	if dto.Name != nil {
		params.Name = pgtype.Text{
			String: *dto.Name,
			Valid:  true,
		}
	}

	return s.queries.PatchZoneTag(ctx, params)
}
