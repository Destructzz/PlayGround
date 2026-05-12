package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type PlaceService struct {
	queries *sqlc.Queries
}

func NewPlaceService(queries *sqlc.Queries) *PlaceService {
	return &PlaceService{queries: queries}
}

func (s *PlaceService) CreatePlace(ctx context.Context, req domain.CreatePlaceRequest) (sqlc.ZonePlace, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var configID pgtype.Int8
	if req.ConfigurationID != nil {
		configID = pgtype.Int8{Int64: *req.ConfigurationID, Valid: true}
	}

	return s.queries.CreateZonePlace(ctx, sqlc.CreateZonePlaceParams{
		ZoneID:          req.ZoneID,
		Label:           req.Label,
		ConfigurationID: configID,
		SortOrder:       req.SortOrder,
		IsActive:        isActive,
	})
}

func (s *PlaceService) GetPlaces(ctx context.Context) ([]sqlc.ZonePlace, error) {
	return s.queries.GetAllZonePlaces(ctx)
}

func (s *PlaceService) GetPlaceByID(ctx context.Context, id int64) (sqlc.ZonePlace, error) {
	return s.queries.GetZonePlaceByID(ctx, id)
}

func (s *PlaceService) DeleteByID(ctx context.Context, id int64) (int64, error) {
	return s.queries.DeleteZonePlace(ctx, id)
}

func (s *PlaceService) PatchByID(ctx context.Context, id int64, req domain.PatchPlaceRequest) (sqlc.ZonePlace, error) {
	params := sqlc.UpdateZonePlaceParams{
		ID: id,
	}

	if req.Label != nil {
		params.Label = pgtype.Text{String: *req.Label, Valid: true}
	}
	if req.ConfigurationID != nil {
		params.ConfigurationID = pgtype.Int8{Int64: *req.ConfigurationID, Valid: true}
	}
	if req.SortOrder != nil {
		params.SortOrder = pgtype.Int4{Int32: *req.SortOrder, Valid: true}
	}
	if req.IsActive != nil {
		params.IsActive = pgtype.Bool{Bool: *req.IsActive, Valid: true}
	}

	return s.queries.UpdateZonePlace(ctx, params)
}
