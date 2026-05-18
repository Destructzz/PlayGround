package service

import (
	"context"

	"backend/internal/domain"
	"backend/internal/repo/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	queries *sqlc.Queries
}

func NewUserService(queries *sqlc.Queries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) UpsertUser(ctx context.Context, req domain.UpsertUserRequest) (sqlc.User, error) {
	return s.queries.UpsertUser(ctx, sqlc.UpsertUserParams{
		GoogleID: pgtype.Text{
			String: req.GoogleID,
			Valid:  req.GoogleID != "",
		},
		FullName: req.FullName,
		Email:    req.Email,
		AvatarUrl: pgtype.Text{
			String: req.AvatarURL,
			Valid:  req.AvatarURL != "",
		},
	})
}

func (s *UserService) ListUsers(ctx context.Context) ([]sqlc.User, error) {
	return s.queries.ListUsers(ctx)
}

func (s *UserService) PatchUser(ctx context.Context, userID pgtype.UUID, req domain.PatchUserRequest) (sqlc.User, error) {
	params := sqlc.PatchUserParams{ID: userID}
	
	if req.FullName != nil {
		params.FullName = pgtype.Text{String: *req.FullName, Valid: true}
	}
	if req.Phone != nil {
		params.Phone = pgtype.Text{String: *req.Phone, Valid: true}
	}
	
	return s.queries.PatchUser(ctx, params)
}

func (s *UserService) SetUserRole(ctx context.Context, userID pgtype.UUID, role sqlc.Role) (sqlc.User, error) {
	return s.queries.SetUserRole(ctx, sqlc.SetUserRoleParams{
		ID:   userID,
		Role: role,
	})
}

func (s *UserService) SearchUsersByEmail(ctx context.Context, query string) ([]sqlc.User, error) {
	return s.queries.SearchUsersByEmail(ctx, pgtype.Text{String: query, Valid: true})
}

func (s *UserService) ListSellers(ctx context.Context) ([]sqlc.User, error) {
	return s.queries.ListSellers(ctx)
}
