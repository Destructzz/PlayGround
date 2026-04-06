package service

import (
	"context"
	"time"

	"backend/internal/domain"
	"backend/internal/repo/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type StaffService struct {
	queries *sqlc.Queries
}

func NewStaff(queries *sqlc.Queries) *StaffService {
	return &StaffService{queries: queries}
}

func (s *StaffService) CreateStaff(ctx context.Context, dto domain.CreateStaffRequest) (sqlc.Staff, error) {
	var userID pgtype.UUID
	if err := userID.Scan(dto.UserID); err != nil {
		return sqlc.Staff{}, err
	}

	hireDate, err := time.Parse("2006-01-02", dto.HireDate)
	if err != nil {
		return sqlc.Staff{}, err
	}

	params := sqlc.CreateStaffParams{
		UserID:   userID,
		Position: dto.Position,
		HireDate: pgtype.Date{
			Time:  hireDate,
			Valid: true,
		},
		Phone: pgtype.Text{
			String: dto.Phone,
			Valid:  dto.Phone != "",
		},
		Email: pgtype.Text{
			String: dto.Email,
			Valid:  dto.Email != "",
		},
	}

	if dto.IsActive != nil {
		params.IsActive = *dto.IsActive
	}

	return s.queries.CreateStaff(ctx, params)
}

func (s *StaffService) ListStaff(ctx context.Context) ([]sqlc.Staff, error) {
	return s.queries.ListStaff(ctx)
}

func (s *StaffService) GetStaffByID(ctx context.Context, id int64) (sqlc.Staff, error) {
	return s.queries.GetStaffByID(ctx, id)
}

func (s *StaffService) PatchStaff(ctx context.Context, id int64, dto domain.PatchStaffRequest) (sqlc.Staff, error) {
	params := sqlc.PatchStaffParams{ID: id}

	if dto.Position != nil {
		params.Position = sqlc.NullPositionType{
			PositionType: *dto.Position,
			Valid:        true,
		}
	}

	if dto.HireDate != nil {
		t, err := time.Parse("2006-01-02", *dto.HireDate)
		if err != nil {
			return sqlc.Staff{}, err
		}
		params.HireDate = pgtype.Date{
			Time:  t,
			Valid: true,
		}
	}

	if dto.Phone != nil {
		params.Phone = pgtype.Text{
			String: *dto.Phone,
			Valid:  true,
		}
	}

	if dto.Email != nil {
		params.Email = pgtype.Text{
			String: *dto.Email,
			Valid:  true,
		}
	}

	if dto.IsActive != nil {
		params.IsActive = pgtype.Bool{
			Bool:  *dto.IsActive,
			Valid: true,
		}
	}

	return s.queries.PatchStaff(ctx, params)
}

func (s *StaffService) DeleteStaff(ctx context.Context, id int64) (int64, error) {
	return s.queries.DeleteStaff(ctx, id)
}
