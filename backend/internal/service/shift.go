package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidShiftTime      = errors.New("invalid shift time")
	ErrInvalidShiftTimeRange = errors.New("invalid shift time range")
	ErrShiftOverlap          = errors.New("shift overlaps with an existing shift for this zone tag")
)

type ShiftService struct {
	queries *sqlc.Queries
}

func NewShift(queries *sqlc.Queries) *ShiftService {
	return &ShiftService{queries: queries}
}

func (s *ShiftService) CreateShift(ctx context.Context, userID pgtype.UUID, dto domain.CreateShiftRequest) (sqlc.GetShiftByIDRow, error) {
	startTime, err := parseShiftTime(dto.StartTime)
	if err != nil {
		return sqlc.GetShiftByIDRow{}, err
	}

	endTime, err := parseShiftTime(dto.EndTime)
	if err != nil {
		return sqlc.GetShiftByIDRow{}, err
	}

	if err := validateShiftWindow(startTime, endTime); err != nil {
		return sqlc.GetShiftByIDRow{}, err
	}

	var zoneTagID pgtype.Int8
	if dto.ZoneTagID != nil {
		zoneTagID = pgtype.Int8{Int64: *dto.ZoneTagID, Valid: true}
	}

	// Check for overlapping shifts within the same zone_tag
	if zoneTagID.Valid {
		overlaps, overlapErr := s.queries.HasOverlappingShift(ctx, sqlc.HasOverlappingShiftParams{
			ZoneTagID: zoneTagID,
			ExcludeID: 0,
			StartTime: pgtype.Timestamptz{Time: startTime, Valid: true},
			EndTime:   pgtype.Timestamptz{Time: endTime, Valid: true},
		})
		if overlapErr != nil {
			return sqlc.GetShiftByIDRow{}, overlapErr
		}
		if overlaps {
			return sqlc.GetShiftByIDRow{}, ErrShiftOverlap
		}
	}

	createdShift, err := s.queries.CreateShift(ctx, sqlc.CreateShiftParams{
		UserID:    userID,
		ZoneTagID: zoneTagID,
		StartTime: pgtype.Timestamptz{Time: startTime, Valid: true},
		EndTime:   pgtype.Timestamptz{Time: endTime, Valid: true},
		Note:      nullableText(dto.Note),
	})
	if err != nil {
		return sqlc.GetShiftByIDRow{}, err
	}

	return s.queries.GetShiftByID(ctx, createdShift.ID)
}

func (s *ShiftService) ListShifts(ctx context.Context) ([]sqlc.ListShiftsRow, error) {
	return s.queries.ListShifts(ctx)
}

func (s *ShiftService) GetShiftByID(ctx context.Context, id int64) (sqlc.GetShiftByIDRow, error) {
	return s.queries.GetShiftByID(ctx, id)
}

func (s *ShiftService) ListShiftsByZoneTagID(ctx context.Context, zoneTagID int64) ([]sqlc.ListShiftsByZoneTagIDRow, error) {
	return s.queries.ListShiftsByZoneTagID(ctx, pgtype.Int8{Int64: zoneTagID, Valid: true})
}

func (s *ShiftService) GetServiceDurationsByZoneTagID(ctx context.Context, zoneTagID int64) ([]int32, error) {
	return s.queries.GetServiceDurationsByZoneTag(ctx, pgtype.Int4{Int32: int32(zoneTagID), Valid: true})
}

func (s *ShiftService) PatchShift(ctx context.Context, id int64, dto domain.PatchShiftRequest) (sqlc.GetShiftByIDRow, error) {
	params := sqlc.PatchShiftParams{ID: id}

	var startTime *time.Time
	if dto.StartTime != nil {
		parsedStartTime, err := parseShiftTime(*dto.StartTime)
		if err != nil {
			return sqlc.GetShiftByIDRow{}, err
		}

		startTime = &parsedStartTime
		params.StartTime = pgtype.Timestamptz{Time: parsedStartTime, Valid: true}
	}

	var endTime *time.Time
	if dto.EndTime != nil {
		parsedEndTime, err := parseShiftTime(*dto.EndTime)
		if err != nil {
			return sqlc.GetShiftByIDRow{}, err
		}

		endTime = &parsedEndTime
		params.EndTime = pgtype.Timestamptz{Time: parsedEndTime, Valid: true}
	}

	if startTime != nil && endTime != nil {
		if err := validateShiftWindow(*startTime, *endTime); err != nil {
			return sqlc.GetShiftByIDRow{}, err
		}
	}

	if dto.ZoneTagID != nil {
		params.ZoneTagID = pgtype.Int8{Int64: *dto.ZoneTagID, Valid: true}
	}

	if dto.Note != nil {
		params.Note = pgtype.Text{String: *dto.Note, Valid: true}
	}

	// Check for overlapping shifts if zone_tag_id is set
	if params.ZoneTagID.Valid {
		checkStart := params.StartTime
		checkEnd := params.EndTime
		// If times not being patched, fetch existing to validate window against new zone_tag
		if !checkStart.Valid || !checkEnd.Valid {
			existing, fetchErr := s.queries.GetShiftByID(ctx, id)
			if fetchErr != nil {
				return sqlc.GetShiftByIDRow{}, fetchErr
			}
			if !checkStart.Valid {
				checkStart = existing.Shift.StartTime
			}
			if !checkEnd.Valid {
				checkEnd = existing.Shift.EndTime
			}
		}

		overlaps, overlapErr := s.queries.HasOverlappingShift(ctx, sqlc.HasOverlappingShiftParams{
			ZoneTagID: params.ZoneTagID,
			ExcludeID: id,
			StartTime: checkStart,
			EndTime:   checkEnd,
		})
		if overlapErr != nil {
			return sqlc.GetShiftByIDRow{}, overlapErr
		}
		if overlaps {
			return sqlc.GetShiftByIDRow{}, ErrShiftOverlap
		}
	}

	updatedShift, err := s.queries.PatchShift(ctx, params)
	if err != nil {
		return sqlc.GetShiftByIDRow{}, err
	}

	return s.queries.GetShiftByID(ctx, updatedShift.ID)
}

func (s *ShiftService) DeleteShift(ctx context.Context, id int64) (int64, error) {
	return s.queries.DeleteShift(ctx, id)
}

func parseShiftTime(value string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidShiftTime, value)
	}

	return parsedTime, nil
}

func validateShiftWindow(startTime, endTime time.Time) error {
	if !endTime.After(startTime) {
		return ErrInvalidShiftTimeRange
	}

	return nil
}

func nullableText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: *value, Valid: true}
}
