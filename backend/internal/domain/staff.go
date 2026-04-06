package domain

import (
	"backend/internal/repo/sqlc"
)

type CreateStaffRequest struct {
	UserID   string            `json:"user_id" binding:"required"`
	Position sqlc.PositionType `json:"position" binding:"required,oneof=admin seller operator tech"`
	HireDate string            `json:"hire_date" binding:"required"`
	Phone    string            `json:"phone"`
	Email    string            `json:"email"`
	IsActive *bool             `json:"is_active"`
}

type PatchStaffRequest struct {
	Position *sqlc.PositionType `json:"position" binding:"omitempty,oneof=admin seller operator tech"`
	HireDate *string            `json:"hire_date" binding:"omitempty"`
	Phone    *string            `json:"phone"`
	Email    *string            `json:"email"`
	IsActive *bool              `json:"is_active"`
}
