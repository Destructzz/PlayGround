package response

import "time"

type DeleteShiftResponse struct {
	ID        int64     `json:"id" example:"1"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ShiftResponse struct {
	Shift     ShiftDoc  `json:"shift"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ShiftListResponse struct {
	Shifts    []ShiftDoc `json:"shifts"`
	Timestamp time.Time  `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string     `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ShiftDoc struct {
	ID        int64        `json:"id" example:"1"`
	UserID    string       `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ZoneTagID *int64       `json:"zone_tag_id,omitempty" example:"1"`
	StartTime time.Time    `json:"start_time" example:"2026-03-22T09:00:00Z"`
	EndTime   time.Time    `json:"end_time" example:"2026-03-22T18:00:00Z"`
	Note      *string      `json:"note,omitempty" example:"Front desk day shift"`
	CreatedAt time.Time    `json:"created_at" example:"2026-01-19T15:37:27.514667373Z"`
	UpdatedAt time.Time    `json:"updated_at" example:"2026-01-19T15:37:27.514667373Z"`
	User      ShiftUserDoc `json:"user"`
}

type ShiftUserDoc struct {
	ID        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FullName  string  `json:"full_name" example:"Ada Lovelace"`
	Email     string  `json:"email" example:"ada@example.com"`
	AvatarURL *string `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	Phone     *string `json:"phone,omitempty" example:"+79991234567"`
	Role      string  `json:"role" example:"admin"`
	IsActive  bool    `json:"is_active" example:"true"`
}
