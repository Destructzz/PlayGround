package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type ServiceService struct {
	queries *sqlc.Queries
}

func NewServiceService(queries *sqlc.Queries) *ServiceService {
	return &ServiceService{queries: queries}
}

func (s *ServiceService) CreateService(ctx context.Context, dto domain.CreateServiceRequest) (sqlc.Service, error) {
	createParams := sqlc.CreateServiceParams{
		Name: dto.Name,
		ZoneID: dto.ZoneID,
		Duration: int32(dto.Duration),
		// We will set this securely below by using price.Scan()
		// Price: ... 
		Currency: dto.Currency,
		Description: pgtype.Text{
			String: dto.Description,
			Valid: dto.Description != "",
		},
	}

	if dto.IsActive != nil {
		createParams.IsActive = *dto.IsActive
	}

	createParams.DetailsJson = []byte("{}")

	var parsedPrice pgtype.Numeric
	if err := parsedPrice.Scan(dto.Price); err != nil {
		return sqlc.Service{}, err
	}
	createParams.Price = parsedPrice

	return s.queries.CreateService(ctx, createParams)
}

func (s *ServiceService) ListServices(ctx context.Context) ([]sqlc.Service, error) {
	return s.queries.ListServices(ctx)
}

func (s *ServiceService) GetServiceByID(ctx context.Context, id int64) (sqlc.Service, error) {
	return s.queries.GetServiceByID(ctx, id)
}

func (s *ServiceService) UpdateService(ctx context.Context, id int64, dto domain.UpdateServiceRequest) (sqlc.Service, error) {
	patchParams := sqlc.PatchServiceParams{ID: id}

	if dto.Name != nil {
		patchParams.Name = pgtype.Text{
			String: *dto.Name,
			Valid: true,
		}
	}

	if dto.ZoneID != nil {
		patchParams.ZoneID = pgtype.Int8{
			Int64: *dto.ZoneID,
			Valid: true,
		}
	}

	if dto.Duration != nil {
		patchParams.Duration = pgtype.Int4{
			Int32: int32(*dto.Duration),
			Valid: true,
		}
	}

	if dto.Price != nil {
		var parsedPrice pgtype.Numeric
		if err := parsedPrice.Scan(*dto.Price); err != nil {
			return sqlc.Service{}, err
		}
		patchParams.Price = parsedPrice
	}

	if dto.Currency != nil {
		patchParams.Currency = pgtype.Text{
			String: *dto.Currency,
			Valid: true,
		}
	}

	if dto.Description != nil {
		patchParams.Description = pgtype.Text{
			String: *dto.Description,
			Valid: true,
		}
	}

	if dto.IsActive != nil {
		patchParams.IsActive = pgtype.Bool{
			Bool: *dto.IsActive,
			Valid: true,
		}
	}

	return s.queries.PatchService(ctx, patchParams)
}

func (s *ServiceService) DeleteService(ctx context.Context, id int64) (int64, error) {
	return s.queries.DeleteService(ctx, id)
}
