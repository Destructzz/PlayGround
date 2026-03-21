package domain

import (
	"backend/internal/repo/sqlc"
)

type CreateBookingRequest struct {
	UserID       string             `json:"user_id" binding:"required"`
	ZoneID       int64              `json:"zone_id" binding:"required"`
	ServiceID    int64              `json:"service_id" binding:"required"`
	StartTime    string             `json:"start_time" binding:"required"`
	EndTime      string             `json:"end_time" binding:"required"`
	Participants int                `json:"participants" binding:"required,min=1"`
	Status       sqlc.BookingStatus `json:"status" binding:"required,oneof=created confirmed canceled completed"`
}

type PatchBookingRequest struct {
	UserID       *string             `json:"user_id" binding:"omitempty"`
	ZoneID       *int64              `json:"zone_id" binding:"omitempty"`
	ServiceID    *int64              `json:"service_id" binding:"omitempty"`
	StartTime    *string             `json:"start_time" binding:"omitempty"`
	EndTime      *string             `json:"end_time" binding:"omitempty"`
	Participants *int                `json:"participants" binding:"omitempty,min=1"`
	Status       *sqlc.BookingStatus `json:"status" binding:"omitempty,oneof=created confirmed canceled completed"`
}
