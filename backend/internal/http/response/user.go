package response

import "time"

type UserDoc struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FullName  string    `json:"full_name" example:"Иван Иванов"`
	Email     string    `json:"email" example:"ivan@playground.local"`
	Phone     string    `json:"phone" example:"+79991112233"`
	Role      string    `json:"role" example:"seller"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2026-05-18T10:00:00Z"`
}

type SearchUsersResponse struct {
	Users     []UserDoc `json:"users"`
	Timestamp time.Time `json:"timestamp" example:"2026-05-18T10:00:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type ListSellersResponse struct {
	Sellers   []UserDoc `json:"sellers"`
	Timestamp time.Time `json:"timestamp" example:"2026-05-18T10:00:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type SetUserRoleResponse struct {
	User      UserDoc   `json:"user"`
	Timestamp time.Time `json:"timestamp" example:"2026-05-18T10:00:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}
