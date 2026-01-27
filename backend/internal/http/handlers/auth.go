package handlers

import (
	"backend/internal/http/middleware"
	"backend/internal/http/response"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
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
	setProviderQuery(c, provider)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Callback completes OAuth flow and returns user profile.
// @Summary     OAuth callback
// @Description Completes OAuth flow and returns user profile
// @Tags        auth
// @Param       provider path string true "OAuth provider" Enums(google)
// @Produce     json
// @Success     200 {object} response.AuthResponse
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
		response.ErrorUnauthorized(c, "auth_failed", "Authentication failed", nil)
		return
	}

	response.Auth(c, response.AuthUserResponse{
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Name:      user.Name,
		Provider:  user.Provider,
	})
}

func ensureProvider(c *gin.Context) (string, bool) {
	provider := c.Param("provider")
	if provider == "" {
		provider = c.Query("provider")
	}
	if provider == "" {
		response.ErrorBadRequest(c, "provider_required", "Provider is required", nil)
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
