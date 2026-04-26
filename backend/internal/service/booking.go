package service

import (
	"context"
	"errors"
	"time"

	"backend/internal/domain"
	"backend/internal/repo/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBookingOverlap       = errors.New("booking overlaps with existing booking")
	ErrPlaceUnavailable     = errors.New("gaming place is unavailable for this time")
	ErrCapacityExceeded     = errors.New("participants exceed zone capacity for this time")
)

type BookingService struct {
	queries *sqlc.Queries
}

func NewBooking(queries *sqlc.Queries) *BookingService {
	return &BookingService{queries: queries}
}

// CreateBooking creates a booking with the provided session user ID.
// It enforces overlap, place availability, and capacity invariants.
func (b *BookingService) CreateBooking(ctx context.Context, userID pgtype.UUID, dto domain.CreateBookingRequest) (sqlc.Booking, error) {
	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		return sqlc.Booking{}, err
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		return sqlc.Booking{}, err
	}

	startTZ := pgtype.Timestamptz{Time: startTime, Valid: true}
	endTZ := pgtype.Timestamptz{Time: endTime, Valid: true}

	// Check zone/service overlap
	overlapCount, err := b.queries.CheckZoneBookingOverlap(ctx, sqlc.CheckZoneBookingOverlapParams{
		ZoneID:    dto.ZoneID,
		ServiceID: dto.ServiceID,
		EndTime:   endTZ,
		StartTime: startTZ,
	})
	if err != nil {
		return sqlc.Booking{}, err
	}
	if overlapCount > 0 {
		return sqlc.Booking{}, ErrBookingOverlap
	}

	// Check place availability if place_id is provided
	if dto.PlaceID != nil {
		placeOverlap, err := b.queries.CheckPlaceBookingOverlap(ctx, sqlc.CheckPlaceBookingOverlapParams{
			PlaceID: pgtype.Int8{Int64: *dto.PlaceID, Valid: true},
			EndTime: endTZ,
			StartTime: startTZ,
		})
		if err != nil {
			return sqlc.Booking{}, err
		}
		if placeOverlap > 0 {
			return sqlc.Booking{}, ErrPlaceUnavailable
		}
	}

	// Check capacity
	bookedCount, err := b.queries.CountActiveBookingsForZone(ctx, sqlc.CountActiveBookingsForZoneParams{
		ZoneID:    dto.ZoneID,
		EndTime:   endTZ,
		StartTime: startTZ,
	})
	if err != nil {
		return sqlc.Booking{}, err
	}

	zone, err := b.queries.GetZoneByID(ctx, dto.ZoneID)
	if err != nil {
		return sqlc.Booking{}, err
	}

	if int32(bookedCount)+int32(dto.Participants) > zone.Capacity {
		return sqlc.Booking{}, ErrCapacityExceeded
	}

	status := dto.Status
	if status == "" {
		status = sqlc.BookingStatusCreated
	}

	detailsJSON := []byte("{}")
	if dto.DetailsJSON != "" {
		detailsJSON = []byte(dto.DetailsJSON)
	}

	var placeID pgtype.Int8
	if dto.PlaceID != nil {
		placeID = pgtype.Int8{Int64: *dto.PlaceID, Valid: true}
	}

	return b.queries.CreateBooking(ctx, sqlc.CreateBookingParams{
		UserID:       userID,
		ZoneID:       dto.ZoneID,
		ServiceID:    dto.ServiceID,
		PlaceID:      placeID,
		StartTime:    startTZ,
		EndTime:      endTZ,
		Participants: int32(dto.Participants),
		Status:       status,
		ContactName:  dto.ContactName,
		ContactEmail: dto.ContactEmail,
		ContactPhone: dto.ContactPhone,
		DetailsJson:  detailsJSON,
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
