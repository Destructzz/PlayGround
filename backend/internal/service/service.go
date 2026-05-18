package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type ServiceService struct {
	queries *sqlc.Queries
}

func NewServiceService(queries *sqlc.Queries) *ServiceService {
	return &ServiceService{queries: queries}
}

func (s *ServiceService) CreateService(ctx context.Context, dto domain.CreateServiceRequest) (sqlc.Service, error) {
	currency := dto.Currency
	if currency == "" {
		currency = "RUB"
	}

	createParams := sqlc.CreateServiceParams{
		Name: dto.Name,
		ZoneID: dto.ZoneID,
		Duration: int32(dto.Duration),
		Currency: currency,
		Description: pgtype.Text{
			String: dto.Description,
			Valid: dto.Description != "",
		},
	}

	if dto.IsActive != nil {
		createParams.IsActive = *dto.IsActive
	}

	if dto.DetailsJSON != "" {
		createParams.DetailsJson = []byte(dto.DetailsJSON)
	} else {
		createParams.DetailsJson = []byte("{}")
	}

	var priceStr string
	switch p := dto.Price.(type) {
	case string:
		priceStr = p
	case float64:
		priceStr = fmt.Sprintf("%.2f", p)
	case float32:
		priceStr = fmt.Sprintf("%.2f", p)
	case int:
		priceStr = fmt.Sprintf("%d", p)
	case int64:
		priceStr = fmt.Sprintf("%d", p)
	case int32:
		priceStr = fmt.Sprintf("%d", p)
	default:
		priceStr = fmt.Sprintf("%v", p)
	}

	var parsedPrice pgtype.Numeric
	if err := parsedPrice.Scan(priceStr); err != nil {
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
		var priceStr string
		switch p := (*dto.Price).(type) {
		case string:
			priceStr = p
		case float64:
			priceStr = fmt.Sprintf("%.2f", p)
		case float32:
			priceStr = fmt.Sprintf("%.2f", p)
		case int:
			priceStr = fmt.Sprintf("%d", p)
		case int64:
			priceStr = fmt.Sprintf("%d", p)
		case int32:
			priceStr = fmt.Sprintf("%d", p)
		default:
			priceStr = fmt.Sprintf("%v", p)
		}

		var parsedPrice pgtype.Numeric
		if err := parsedPrice.Scan(priceStr); err != nil {
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

	if dto.DetailsJSON != nil {
		patchParams.DetailsJson = []byte(*dto.DetailsJSON)
	}

	return s.queries.PatchService(ctx, patchParams)
}

func (s *ServiceService) DeleteService(ctx context.Context, id int64) (int64, error) {
	return s.queries.DeleteService(ctx, id)
}
