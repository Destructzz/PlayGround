package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type ComputerConfigurationService struct {
	queries *sqlc.Queries
}

func NewComputerConfiguration(queries *sqlc.Queries) *ComputerConfigurationService {
	return &ComputerConfigurationService{queries: queries}
}

func (s *ComputerConfigurationService) Create(ctx context.Context, dto domain.CreateComputerConfigurationRequest) (sqlc.ComputerConfiguration, error) {
	specs, err := json.Marshal(dto.SpecsJSON)
	if err != nil {
		return sqlc.ComputerConfiguration{}, err
	}
	if len(specs) == 0 {
		specs = []byte("[]")
	}

	return s.queries.CreateComputerConfiguration(ctx, sqlc.CreateComputerConfigurationParams{
		ZoneTagsID: dto.ZoneTagID,
		SpecsJson:  specs,
	})
}

func (s *ComputerConfigurationService) List(ctx context.Context) ([]sqlc.ComputerConfiguration, error) {
	return s.queries.ListComputerConfigurations(ctx)
}

func (s *ComputerConfigurationService) GetByID(ctx context.Context, id int64) (sqlc.ComputerConfiguration, error) {
	return s.queries.GetComputerConfigurationByID(ctx, id)
}

func (s *ComputerConfigurationService) Delete(ctx context.Context, id int64) (int64, error) {
	return s.queries.DeleteComputerConfiguration(ctx, id)
}

func (s *ComputerConfigurationService) Patch(ctx context.Context, id int64, dto domain.PatchComputerConfigurationRequest) (sqlc.ComputerConfiguration, error) {
	params := sqlc.PatchComputerConfigurationParams{ID: id}

	if dto.ZoneTagID != nil {
		params.ZoneTagsID = pgtype.Int8{
			Int64: *dto.ZoneTagID,
			Valid: true,
		}
	}

	if dto.SpecsJSON != nil {
		specs, err := json.Marshal(*dto.SpecsJSON)
		if err != nil {
			return sqlc.ComputerConfiguration{}, err
		}
		params.SpecsJson = specs
	}

	return s.queries.PatchComputerConfiguration(ctx, params)
}
