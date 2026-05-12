package domain

import (
	"backend/internal/repo/sqlc"
	"backend/pkg"
	"time"
)

type CreateShiftRequest struct {
	StartTime string  `json:"start_time" binding:"required"`
	EndTime   string  `json:"end_time" binding:"required"`
	ZoneTagID *int64  `json:"zone_tag_id" binding:"omitempty"`
	Note      *string `json:"note"`
}

type PatchShiftRequest struct {
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	ZoneTagID *int64  `json:"zone_tag_id"`
	Note      *string `json:"note"`
}

type ShiftViewForAdmin struct {
	ID        int64         `json:"id"`
	UserID    string        `json:"user_id"`
	ZoneTagID *int64        `json:"zone_tag_id,omitempty"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Note      *string       `json:"note,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	User      ShiftUserView `json:"user"`
}

type ShiftViewForUser struct {
	ID        int64         `json:"id"`
	ZoneTagID *int64        `json:"zone_tag_id,omitempty"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Note      *string       `json:"note,omitempty"`
}

type ShiftUserView struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
}

func NewShiftViewForAdmin(shift sqlc.Shift, user sqlc.User) ShiftViewForAdmin {
	var zoneTagID *int64
	if shift.ZoneTagID.Valid {
		zoneTagID = &shift.ZoneTagID.Int64
	}

	return ShiftViewForAdmin{
		ID:        shift.ID,
		UserID:    pkg.UUIDString(shift.UserID),
		ZoneTagID: zoneTagID,
		StartTime: shift.StartTime.Time,
		EndTime:   shift.EndTime.Time,
		Note:      pkg.TextPointer(shift.Note),
		CreatedAt: shift.CreatedAt.Time,
		UpdatedAt: shift.UpdatedAt.Time,
		User: ShiftUserView{
			ID:        pkg.UUIDString(user.ID),
			FullName:  user.FullName,
			Email:     user.Email,
			AvatarURL: pkg.TextPointer(user.AvatarUrl),
			Phone:     pkg.TextPointer(user.Phone),
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
	}
}

func NewShiftViewForUser(shift sqlc.Shift) ShiftViewForUser {
	var zoneTagID *int64
	if shift.ZoneTagID.Valid {
		zoneTagID = &shift.ZoneTagID.Int64
	}

	return ShiftViewForUser{
		ID:        shift.ID,
		ZoneTagID: zoneTagID,
		StartTime: shift.StartTime.Time,
		EndTime:   shift.EndTime.Time,
		Note:      pkg.TextPointer(shift.Note),
	}
}

func NewShiftListForAdmin(rows []sqlc.ListShiftsRow) []ShiftViewForAdmin {
	shifts := make([]ShiftViewForAdmin, 0, len(rows))
	for _, row := range rows {
		shifts = append(shifts, NewShiftViewForAdmin(row.Shift, row.User))
	}

	return shifts
}

func NewShiftListForUser(rows []sqlc.ListShiftsRow) []ShiftViewForUser {
	shifts := make([]ShiftViewForUser, 0, len(rows))
	for _, row := range rows {
		shifts = append(shifts, NewShiftViewForUser(row.Shift))
	}

	return shifts
}
