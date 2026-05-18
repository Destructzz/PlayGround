package domain

import (
	"backend/internal/repo/sqlc"
)

type CreateBookingRequest struct {
	ZoneID       int64              `json:"zone_id" binding:"required"`
	ServiceID    int64              `json:"service_id" binding:"omitempty"`
	PlaceID      *int64             `json:"place_id" binding:"omitempty"`
	StartTime    string             `json:"start_time" binding:"required"`
	EndTime      string             `json:"end_time" binding:"required"`
	Participants int                `json:"participants" binding:"required,min=1"`
	Status       sqlc.BookingStatus `json:"status" binding:"omitempty,oneof=created confirmed canceled completed"`
	ContactName  string             `json:"contact_name"`
	ContactEmail string             `json:"contact_email"`
	ContactPhone string             `json:"contact_phone"`
	DetailsJSON  string             `json:"details_json"`
}

type PatchBookingRequest struct {
	ZoneID       *int64              `json:"zone_id" binding:"omitempty"`
	ServiceID    *int64              `json:"service_id" binding:"omitempty"`
	StartTime    *string             `json:"start_time" binding:"omitempty"`
	EndTime      *string             `json:"end_time" binding:"omitempty"`
	Participants *int                `json:"participants" binding:"omitempty,min=1"`
	Status       *sqlc.BookingStatus `json:"status" binding:"omitempty,oneof=created confirmed canceled completed"`
	ContactName  *string             `json:"contact_name" binding:"omitempty"`
	ContactEmail *string             `json:"contact_email" binding:"omitempty"`
	ContactPhone *string             `json:"contact_phone" binding:"omitempty"`
}
