package response

import (
	"time"
)

type DeleteStaffResponse struct {
	ID        int64     `json:"id" example:"1"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type StaffResponse struct {
	Staff     StaffDoc  `json:"staff"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type StaffListResponse struct {
	Staff     []StaffDoc `json:"staff"`
	Timestamp time.Time  `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string     `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type StaffDoc struct {
	ID        int64     `json:"id" example:"1"`
	UserID    string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Position  string    `json:"position" example:"admin"`
	HireDate  string    `json:"hire_date" example:"2026-01-15"`
	Phone     *string   `json:"phone,omitempty" example:"+79001234567"`
	Email     *string   `json:"email,omitempty" example:"staff@example.com"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2026-01-19T15:37:27.514667373Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-01-19T15:37:27.514667373Z"`
}
