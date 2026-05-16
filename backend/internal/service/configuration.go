package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"backend/pkg"
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

func (s *ComputerConfigurationService) Create(ctx context.Context, dto domain.CreateComputerConfigurationRequest) (domain.ComputerConfiguration, error) {
	specs, err := json.Marshal(dto.SpecsJSON)
	if err != nil {
		return domain.ComputerConfiguration{}, err
	}
	if len(specs) == 0 {
		specs = []byte("[]")
	}

	config, err := s.queries.CreateComputerConfiguration(ctx, sqlc.CreateComputerConfigurationParams{
		ZoneTagsID: dto.ZoneTagID,
		SpecsJson:  specs,
	})
	if err != nil {
		return domain.ComputerConfiguration{}, err
	}
	return mapConfigToDomain(config), nil
}

func (s *ComputerConfigurationService) List(ctx context.Context) ([]domain.ComputerConfiguration, error) {
	items, err := s.queries.ListComputerConfigurations(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ComputerConfiguration, 0, len(items))
	for _, item := range items {
		result = append(result, mapConfigToDomain(item))
	}
	return result, nil
}

func (s *ComputerConfigurationService) GetByID(ctx context.Context, id int64) (domain.ComputerConfiguration, error) {
	item, err := s.queries.GetComputerConfigurationByID(ctx, id)
	if err != nil {
		return domain.ComputerConfiguration{}, err
	}
	return mapConfigToDomain(item), nil
}

func (s *ComputerConfigurationService) Delete(ctx context.Context, id int64) (int64, error) {
	return s.queries.DeleteComputerConfiguration(ctx, id)
}

func (s *ComputerConfigurationService) Patch(ctx context.Context, id int64, dto domain.PatchComputerConfigurationRequest) (domain.ComputerConfiguration, error) {
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
			return domain.ComputerConfiguration{}, err
		}
		params.SpecsJson = specs
	}

	config, err := s.queries.PatchComputerConfiguration(ctx, params)
	if err != nil {
		return domain.ComputerConfiguration{}, err
	}
	return mapConfigToDomain(config), nil
}

func mapConfigToDomain(c sqlc.ComputerConfiguration) domain.ComputerConfiguration {
	return domain.ComputerConfiguration{
		ID:        c.ID,
		ZoneTagID: c.ZoneTagsID,
		SpecsJSON: pkg.JSONOrArray(c.SpecsJson),
	}
}
