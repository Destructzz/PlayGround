package service

import (
	"context"
	"time"

	"backend/internal/domain"
	"backend/internal/repo/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type BookingService struct {
	queries *sqlc.Queries
}

func NewBooking(queries *sqlc.Queries) *BookingService {
	return &BookingService{queries: queries}
}

func (b *BookingService) CreateBooking(ctx context.Context, dto domain.CreateBookingRequest) (sqlc.Booking, error) {
	var userID pgtype.UUID
	if err := userID.Scan(dto.UserID); err != nil {
		return sqlc.Booking{}, err
	}

	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		return sqlc.Booking{}, err
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		return sqlc.Booking{}, err
	}

	return b.queries.CreateBooking(ctx, sqlc.CreateBookingParams{
		UserID:    userID,
		ZoneID:    dto.ZoneID,
		ServiceID: dto.ServiceID,
		StartTime: pgtype.Timestamptz{
			Time:  startTime,
			Valid: true,
		},
		EndTime: pgtype.Timestamptz{
			Time:  endTime,
			Valid: true,
		},
		Participants: int32(dto.Participants),
		Status:       dto.Status,
	})
}

func (b *BookingService) ListBookings(ctx context.Context) ([]sqlc.Booking, error) {
	return b.queries.ListBookings(ctx)
}

func (b *BookingService) GetBookingByID(ctx context.Context, id int64) (sqlc.Booking, error) {
	return b.queries.GetBookingByID(ctx, id)
}

func (b *BookingService) PatchBooking(ctx context.Context, id int64, dto domain.PatchBookingRequest) (sqlc.Booking, error) {
	params := sqlc.PatchBookingParams{ID: id}

	if dto.UserID != nil {
		var userID pgtype.UUID
		if err := userID.Scan(*dto.UserID); err != nil {
			return sqlc.Booking{}, err
		}
		params.UserID = userID
	}

	if dto.ZoneID != nil {
		params.ZoneID = pgtype.Int8{
			Int64: *dto.ZoneID,
			Valid: true,
		}
	}

	if dto.ServiceID != nil {
		params.ServiceID = pgtype.Int8{
			Int64: *dto.ServiceID,
			Valid: true,
		}
	}

	if dto.StartTime != nil {
		t, err := time.Parse(time.RFC3339, *dto.StartTime)
		if err != nil {
			return sqlc.Booking{}, err
		}
		params.StartTime = pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	if dto.EndTime != nil {
		t, err := time.Parse(time.RFC3339, *dto.EndTime)
		if err != nil {
			return sqlc.Booking{}, err
		}
		params.EndTime = pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	if dto.Participants != nil {
		params.Participants = pgtype.Int4{
			Int32: int32(*dto.Participants),
			Valid: true,
		}
	}

	if dto.Status != nil {
		params.Status = sqlc.NullBookingStatus{
			BookingStatus: *dto.Status,
			Valid:         true,
		}
	}

	return b.queries.PatchBooking(ctx, params)
}

func (b *BookingService) DeleteBooking(ctx context.Context, id int64) (int64, error) {
	return b.queries.DeleteBooking(ctx, id)
}
