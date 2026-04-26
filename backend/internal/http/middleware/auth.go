package middleware

import (
	"backend/internal/http/response"
	"backend/internal/repo/sqlc"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

const (
	sessionCookieName = "pg_session"
	sessionUserKey    = "session_user"
	sessionIDKey      = "session_id"
)

// AuthRequired is middleware that enforces an authenticated session.
// It reads the session cookie, validates it against the DB, and injects the user into the context.
func AuthRequired(queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, sessionID, ok := resolveSession(c, queries)
		if !ok {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusUnauthorized),
				response.WithError("unauthorized", "Authentication required", nil),
			).JSON(c)
			c.Abort()
			return
		}

		c.Set(sessionUserKey, user)
		c.Set(sessionIDKey, sessionID)
		c.Next()
	}
}

// AuthOptional resolves the session if present but does not enforce it.
// Protected routes should use AuthRequired instead.
func AuthOptional(queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, sessionID, ok := resolveSession(c, queries)
		if ok {
			c.Set(sessionUserKey, user)
			c.Set(sessionIDKey, sessionID)
		}
		c.Next()
	}
}

// UserFromContext retrieves the authenticated user from the Gin context.
func UserFromContext(c *gin.Context) (sqlc.User, bool) {
	val, exists := c.Get(sessionUserKey)
	if !exists {
		return sqlc.User{}, false
	}
	user, ok := val.(sqlc.User)
	return user, ok
}

// SessionIDFromContext retrieves the session ID from the Gin context.
func SessionIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(sessionIDKey)
	if !exists {
		return uuid.UUID{}, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// SetSessionCookie sets the browser cookie for the session.
func SetSessionCookie(c *gin.Context, sessionID string) {
	c.SetCookie(
		sessionCookieName,
		sessionID,
		7*24*60*60, // 7 days
		"/",
		"",    // domain auto
		false, // secure = false for local dev
		true,  // httpOnly
	)
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(c *gin.Context) {
	c.SetCookie(
		sessionCookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}

func resolveSession(c *gin.Context, queries *sqlc.Queries) (sqlc.User, uuid.UUID, bool) {
	cookie, err := c.Cookie(sessionCookieName)
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
