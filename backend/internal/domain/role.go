package domain

import "backend/internal/repo/sqlc"

const (
	RoleAdmin  sqlc.Role = "admin"
	RoleSeller sqlc.Role = "seller"
	RoleClient sqlc.Role = "client"
)