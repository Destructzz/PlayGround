package middleware

import (
	"backend/internal/http/response"
	"backend/internal/repo/sqlc"
	"backend/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


// AuthRequired is middleware that enforces an authenticated session.
// It reads the session cookie, validates it against the DB, and injects the user into the context.
func AuthRequired(queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, sessionID, ok := pkg.ResolveSession(c, queries)
		if !ok {
			cookie, err := c.Cookie(pkg.SessionCookieName)
			zap.L().Warn("session not found", zap.Any("user", user), zap.Any("sessionID", sessionID), zap.Any("cookie", cookie), zap.Any("cookie_err", err))
			response.NewResponseBuilder(
				response.WithStatus(http.StatusUnauthorized),
				response.WithError("unauthorized", "Authentication required", nil),
			).JSON(c)
			c.Abort()
			return
		}else{
			cookie, err := c.Cookie(pkg.SessionCookieName)
			zap.L().Warn("session not found", zap.Any("user", user), zap.Any("sessionID", sessionID), zap.Any("cookie", cookie), zap.Any("cookie_err", err))
		}

		c.Set(pkg.SessionUserKey, user)
		c.Set(pkg.SessionIDKey, sessionID)
		c.Next()
	}
}

func AuthRequiredWithRole(queries *sqlc.Queries, roles ...sqlc.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, sessionID, ok := pkg.ResolveSession(c, queries)
		if !ok {
			response.NewResponseBuilder(
				response.WithStatus(http.StatusUnauthorized),
				response.WithError("unauthorized", "Authentication required", nil),
			).JSON(c)
			c.Abort()
			return
		}

		// if user.Role != role {
		// 	response.NewResponseBuilder(
		// 		response.WithStatus(http.StatusUnauthorized),
		// 		response.WithError("unauthorized", "Authorization required", nil),
		// 	).JSON(c)
		// 	c.Abort()
		// 	return
		// }

		for _, r := range roles {
			if user.Role == r {
				c.Set(pkg.SessionUserKey, user)
				c.Set(pkg.SessionIDKey, sessionID)
				c.Next()
				return
			}
		}

		response.NewResponseBuilder(
			response.WithStatus(http.StatusUnauthorized),
			response.WithError("unauthorized", "Authorization required", nil),
		).JSON(c)
		c.Abort()
	}
}

// AuthOptional resolves the session if present but does not enforce it.
// Protected routes should use AuthRequired instead.
func AuthOptional(queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, sessionID, ok := pkg.ResolveSession(c, queries)
		if ok {
			c.Set(pkg.SessionUserKey, user)
			c.Set(pkg.SessionIDKey, sessionID)
		}
		c.Next()
	}
}


// SetSessionCookie sets the browser cookie for the session.
func SetSessionCookie(c *gin.Context, sessionID string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     pkg.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   isSecureRequest(c),
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     pkg.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecureRequest(c),
		SameSite: http.SameSiteLaxMode,
	})
}

func isSecureRequest(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}

	return strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}

