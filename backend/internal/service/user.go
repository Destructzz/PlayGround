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
