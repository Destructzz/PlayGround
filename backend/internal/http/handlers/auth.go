package handlers

import (
	"backend/internal/domain"
	"backend/internal/http/middleware"
	"backend/internal/http/response"
	"backend/internal/repo/sqlc"
	"backend/internal/service"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

type Auth struct {
	userService *service.UserService
	queries     *sqlc.Queries
}

func NewAuth(userService *service.UserService, queries *sqlc.Queries) *Auth {
	return &Auth{userService: userService, queries: queries}
}

// Begin starts OAuth flow for the given provider.
// @Summary     Start OAuth
// @Description Redirects to provider authorization page
// @Tags        auth
// @Param       provider path string true "OAuth provider" Enums(google)
// @Success     307 {string} string "Redirect to provider"
// @Failure     400 {object} response.ErrorResponse
// @Router      /api/v1/auth/{provider} [get]
func (a *Auth) Begin(c *gin.Context) {
	provider, ok := ensureProvider(c)
	if !ok {
		return
	}

	// Preserve return_to for post-auth redirect
	returnTo := c.Query("return_to")
	if returnTo != "" {
		c.SetCookie("pg_return_to", returnTo, 600, "/", "", false, true)
	}

	setProviderQuery(c, provider)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Callback completes OAuth flow, creates a session, and redirects to frontend.
// @Summary     OAuth callback
// @Description Completes OAuth flow and redirects with session cookie
// @Tags        auth
// @Param       provider path string true "OAuth provider" Enums(google)
// @Success     302 {string} string "Redirect to frontend"
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Router      /api/v1/auth/{provider}/callback [get]
func (a *Auth) Callback(c *gin.Context) {
	provider, ok := ensureProvider(c)
	if !ok {
		return
	}
	setProviderQuery(c, provider)

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		logger := zap.L()
		if rid, ok := middleware.RequestIDFromContext(c); ok {
			logger = logger.With(zap.String("request_id", rid))
		}
		logger.Warn("auth callback error", zap.String("provider", provider), zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(401),
			response.WithError("auth_failed", "Authentication failed", nil),
		).JSON(c)
		return
	}

	req := domain.UpsertUserRequest{
		GoogleID:  user.UserID,
		FullName:  user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}

	dbUser, err := a.userService.UpsertUser(c.Request.Context(), req)
	if err != nil {
		zap.L().Warn("failed to upsert user", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("db_error", "Failed to save user", nil),
		).JSON(c)
		return
	}

	// Create session
	session, err := a.queries.CreateSession(c.Request.Context(), dbUser.ID)
	if err != nil {
		zap.L().Warn("failed to create session", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("session_error", "Failed to create session", nil),
		).JSON(c)
		return
	}

	// Set session cookie
	sessionIDBytes, _ := session.ID.Value()
	sessionIDStr := ""
	if v, ok := sessionIDBytes.([16]byte); ok {
		var pgUUID pgtype.UUID
		pgUUID.Bytes = v
		pgUUID.Valid = true
		sessionIDStr = pgUUID.String()
	}
	if sessionIDStr == "" {
		sessionIDStr = session.ID.String()
	}
	middleware.SetSessionCookie(c, sessionIDStr)

	// Redirect to return_to or frontend root
	returnTo, _ := c.Cookie("pg_return_to")
	c.SetCookie("pg_return_to", "", -1, "/", "", false, true)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	redirectTarget := frontendURL + "/"
	if returnTo != "" {
		redirectTarget = frontendURL + returnTo
	}

	c.Redirect(http.StatusFound, redirectTarget)
}

// Session returns the current session state.
// @Summary     Get session
// @Description Returns current authenticated user or signed-out state
// @Tags        auth
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/auth/session [get]
func (a *Auth) Session(c *gin.Context) {
	user, ok := middleware.UserFromContext(c)
	if !ok {
		response.NewResponseBuilder(
			response.WithData("authenticated", false),
			response.WithData("user", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("authenticated", true),
		response.WithData("user", response.AuthUserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			AvatarURL: user.AvatarUrl.String,
			Name:      user.FullName,
			Provider:  "google",
		}),
	).JSON(c)
}

// Logout clears the session.
// @Summary     Logout
// @Description Clears session cookie and deletes server session
// @Tags        auth
// @Success     204
// @Router      /api/v1/auth/logout [post]
func (a *Auth) Logout(c *gin.Context) {
	sessionID, ok := middleware.SessionIDFromContext(c)
	if ok {
		var pgSessionID pgtype.UUID
		_ = pgSessionID.Scan(sessionID.String())
		_ = a.queries.DeleteSession(c.Request.Context(), pgSessionID)
	}
	middleware.ClearSessionCookie(c)
	c.Status(http.StatusNoContent)
}

// DevLogin seeds a session for automated QA in non-production environments.
// @Summary     Dev login
// @Description Creates a test user session for QA (non-production only)
// @Tags        auth
// @Success     200 {object} map[string]interface{}
// @Failure     403 {object} response.ErrorResponse
// @Router      /api/v1/auth/dev-login [post]
func (a *Auth) DevLogin(c *gin.Context) {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusForbidden),
			response.WithError("forbidden", "Dev login is not available in production", nil),
		).JSON(c)
		return
	}

	// Upsert a deterministic test user
	req := domain.UpsertUserRequest{
		GoogleID:  "dev-test-google-id",
		FullName:  "Dev Test User",
		Email:     "dev@playground.local",
		AvatarURL: "",
	}

	dbUser, err := a.userService.UpsertUser(c.Request.Context(), req)
	if err != nil {
		zap.L().Warn("dev-login upsert failed", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("db_error", "Failed to create dev user", nil),
		).JSON(c)
		return
	}

	session, err := a.queries.CreateSession(c.Request.Context(), dbUser.ID)
	if err != nil {
		zap.L().Warn("dev-login session failed", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("session_error", "Failed to create dev session", nil),
		).JSON(c)
		return
	}

	middleware.SetSessionCookie(c, session.ID.String())

	response.NewResponseBuilder(
		response.WithData("authenticated", true),
		response.WithData("user", response.AuthUserResponse{
			ID:        dbUser.ID.String(),
			Email:     dbUser.Email,
			AvatarURL: dbUser.AvatarUrl.String,
			Name:      dbUser.FullName,
			Provider:  "dev",
		}),
	).JSON(c)
}

// ListUsers возвращает список пользователей.
// @Summary     List users
// @Description Returns all active users
// @Tags        users
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/user [get]
func (a *Auth) ListUsers(c *gin.Context) {
	users, err := a.userService.ListUsers(c.Request.Context())
	if err != nil {
		zap.L().Warn("database error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("database_fault", "some problems while using database", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("users", users),
	).JSON(c)
}

func ensureProvider(c *gin.Context) (string, bool) {
	provider := c.Param("provider")
	if provider == "" {
		provider = c.Query("provider")
	}
	if provider == "" {
		response.NewResponseBuilder(
			response.WithStatus(400),
			response.WithError("provider_required", "Provider is required", nil),
		).JSON(c)
		return "", false
	}
	return provider, true
}

func setProviderQuery(c *gin.Context, provider string) {
	values := c.Request.URL.Query()
	if values.Get("provider") == "" {
		values.Set("provider", provider)
		c.Request.URL.RawQuery = values.Encode()
	}

	if c.Request.Form == nil {
		c.Request.Form = url.Values{}
	}
	if c.Request.Form.Get("provider") == "" {
		c.Request.Form.Set("provider", provider)
	}
}
