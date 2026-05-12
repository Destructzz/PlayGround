package pkg

import (
	"backend/internal/repo/sqlc"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

const (
	SessionCookieName = "pg_session"
	SessionUserKey    = "session_user"
	SessionIDKey      = "session_id"
)

func ResolveSession(c *gin.Context, queries *sqlc.Queries) (sqlc.User, uuid.UUID, bool) {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil || cookie == "" {
		return sqlc.User{}, uuid.UUID{}, false
	}

	sessionUUID, err := uuid.Parse(cookie)
	if err != nil {
		return sqlc.User{}, uuid.UUID{}, false
	}

	var pgSessionID pgtype.UUID
	if err := pgSessionID.Scan(sessionUUID.String()); err != nil {
		return sqlc.User{}, uuid.UUID{}, false
	}

	row, err := queries.GetSessionByID(context.Background(), pgSessionID)
	if err != nil {
		zap.L().Debug("session lookup failed", zap.Error(err))
		return sqlc.User{}, uuid.UUID{}, false
	}

	user := sqlc.User{
		ID:        row.UserID,
		GoogleID:  row.GoogleID,
		FullName:  row.FullName,
		Email:     row.Email,
		AvatarUrl: row.AvatarUrl,
		Phone:     row.Phone,
		Role:      row.Role,
	}

	return user, sessionUUID, true
}

// UserFromContext retrieves the authenticated user from the Gin context.
func UserFromContext(c *gin.Context) (sqlc.User, bool) {
	val, exists := c.Get(SessionUserKey)
	if !exists {
		return sqlc.User{}, false
	}
	user, ok := val.(sqlc.User)
	return user, ok
}

// SessionIDFromContext retrieves the session ID from the Gin context.
func SessionIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(SessionIDKey)
	if !exists {
		return uuid.UUID{}, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
