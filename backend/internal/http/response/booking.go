package response

import (
	"time"
)

type DeleteBookingResponse struct {
	ID        int64     `json:"id" example:"1"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type BookingResponse struct {
	Booking   BookingDoc `json:"booking"`
	Timestamp time.Time  `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string     `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type BookingListResponse struct {
	Bookings  []BookingDoc `json:"bookings"`
	Timestamp time.Time    `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string       `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type BookingDoc struct {
	ID           int64     `json:"id" example:"1"`
	UserID       string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ZoneID       int64     `json:"zone_id" example:"1"`
	ServiceID    int64     `json:"service_id" example:"1"`
	StartTime    time.Time `json:"start_time" example:"2026-03-20T10:00:00Z"`
	EndTime      time.Time `json:"end_time" example:"2026-03-20T11:00:00Z"`
	Participants int32     `json:"participants" example:"4"`
	TotalPrice   string    `json:"total_price" example:"1500.00"`
	Status       string    `json:"status" example:"created"`
	CreatedAt    time.Time `json:"created_at" example:"2026-01-19T15:37:27.514667373Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2026-01-19T15:37:27.514667373Z"`
}
