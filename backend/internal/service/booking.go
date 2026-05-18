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
	ErrInvalidBookingRange  = errors.New("booking time range is invalid")
	ErrServiceZoneMismatch  = errors.New("service does not belong to zone")
	ErrPlaceZoneMismatch    = errors.New("place does not belong to zone")
	ErrInactiveService      = errors.New("service is inactive")
	ErrInactivePlace        = errors.New("place is inactive")
	ErrInvalidServiceWindow = errors.New("booking duration must match service duration")
	ErrPastBooking          = errors.New("cannot book in the past")
)

type BookingService struct {
	queries *sqlc.Queries
}

func NewBooking(queries *sqlc.Queries) *BookingService {
	return &BookingService{queries: queries}
}

// CreateBooking creates a booking with the provided session user ID.
// It enforces overlap, place availability, and capacity invariants.
// For lounge and event zones (no dedicated place), exclusive overlap check is
// skipped — only capacity is enforced, since multiple groups can co-exist.
func (b *BookingService) CreateBooking(ctx context.Context, userID pgtype.UUID, dto domain.CreateBookingRequest) (sqlc.Booking, error) {
	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		return sqlc.Booking{}, err
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		return sqlc.Booking{}, err
	}

	if !startTime.Before(endTime) {
		return sqlc.Booking{}, ErrInvalidBookingRange
	}

	if startTime.Before(time.Now()) {
		return sqlc.Booking{}, ErrPastBooking
	}

	zone, err := b.queries.GetZoneByID(ctx, dto.ZoneID)
	if err != nil {
		return sqlc.Booking{}, err
	}

	// For lounge and event zones, if no ServiceID is provided, automatically find or create a free service.
	if zone.ZoneType == sqlc.ZoneTypeLounge || zone.ZoneType == sqlc.ZoneTypeEvent {
		if dto.ServiceID == 0 {
			services, err := b.queries.ListServicesByZoneID(ctx, zone.ID)
			if err == nil && len(services) > 0 {
				dto.ServiceID = services[0].ID
			} else {
				// Create default free service for the zone
				price := pgtype.Numeric{}
				_ = price.Scan("0.00")
				createdService, err := b.queries.CreateService(ctx, sqlc.CreateServiceParams{
					Name:     zone.Name + " Booking",
					ZoneID:   zone.ID,
					Duration: 60, // 1 hour standard slot duration
					Price:    price,
					Currency: "RUB",
					Description: pgtype.Text{
						String: "Бесплатное бронирование",
						Valid:  true,
					},
					IsActive:    true,
					DetailsJson: []byte("{}"),
				})
				if err != nil {
					return sqlc.Booking{}, err
				}
				dto.ServiceID = createdService.ID
			}
		}
	} else if dto.ServiceID == 0 {
		return sqlc.Booking{}, errors.New("service_id is required for gaming zones")
	}

	serviceEntity, err := b.queries.GetServiceByID(ctx, dto.ServiceID)
	if err != nil {
		return sqlc.Booking{}, err
	}
	if serviceEntity.ZoneID != dto.ZoneID {
		return sqlc.Booking{}, ErrServiceZoneMismatch
	}
	if !serviceEntity.IsActive {
		return sqlc.Booking{}, ErrInactiveService
	}

	// For lounge/event zones, we bypass the strict duration restriction because the booking duration is flexible.
	if zone.ZoneType == sqlc.ZoneTypeGame {
		serviceDuration := time.Duration(serviceEntity.Duration) * time.Minute
		if !startTime.Add(serviceDuration).Equal(endTime) {
			return sqlc.Booking{}, ErrInvalidServiceWindow
		}
	}

	startTZ := pgtype.Timestamptz{Time: startTime, Valid: true}
	endTZ := pgtype.Timestamptz{Time: endTime, Valid: true}

	// For game zones without a dedicated place: check exclusive overlap.
	// For lounge/event zones: skip overlap — multiple groups share capacity.
	if dto.PlaceID == nil && zone.ZoneType == sqlc.ZoneTypeGame {
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
	}

	// Check place availability if place_id is provided
	if dto.PlaceID != nil {
		placeEntity, err := b.queries.GetZonePlaceByID(ctx, *dto.PlaceID)
		if err != nil {
			return sqlc.Booking{}, err
		}
		if placeEntity.ZoneID != dto.ZoneID {
			return sqlc.Booking{}, ErrPlaceZoneMismatch
		}
		if !placeEntity.IsActive {
			return sqlc.Booking{}, ErrInactivePlace
		}

		placeOverlap, err := b.queries.CheckPlaceBookingOverlap(ctx, sqlc.CheckPlaceBookingOverlapParams{
			PlaceID:   pgtype.Int8{Int64: *dto.PlaceID, Valid: true},
			EndTime:   endTZ,
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

func (b *BookingService) ListBookingsByUserID(ctx context.Context, userID pgtype.UUID) ([]sqlc.Booking, error) {
	return b.queries.ListBookingsByUser(ctx, userID)
}

func (b *BookingService) GetBookingByID(ctx context.Context, id int64) (sqlc.Booking, error) {
	return b.queries.GetBookingByID(ctx, id)
}

func (b *BookingService) PatchBooking(ctx context.Context, id int64, dto domain.PatchBookingRequest) (sqlc.Booking, error) {
	currentBooking, err := b.queries.GetBookingByID(ctx, id)
	if err != nil {
		return sqlc.Booking{}, err
	}

	if currentBooking.Status == sqlc.BookingStatusCanceled || 
	   currentBooking.Status == sqlc.BookingStatusCompleted || 
	   currentBooking.StartTime.Time.Before(time.Now()) {
		return sqlc.Booking{}, errors.New("cannot edit archived booking")
	}

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

	if dto.ContactName != nil {
		params.ContactName = pgtype.Text{String: *dto.ContactName, Valid: true}
	}
	if dto.ContactEmail != nil {
		params.ContactEmail = pgtype.Text{String: *dto.ContactEmail, Valid: true}
	}
	if dto.ContactPhone != nil {
		params.ContactPhone = pgtype.Text{String: *dto.ContactPhone, Valid: true}
	}

	return b.queries.PatchBooking(ctx, params)
}

func (b *BookingService) DeleteBooking(ctx context.Context, id int64) (int64, error) {
	return b.queries.DeleteBooking(ctx, id)
}
