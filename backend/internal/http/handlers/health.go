package handlers

import (
	"backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Health provides readiness and liveness endpoints.
type Health struct {
	pool *pgxpool.Pool
}

// NewHealth constructs a Health handler.
func NewHealth(pool *pgxpool.Pool) *Health {
	return &Health{pool: pool}
}

// Health returns basic liveness information.
// @Summary     Liveness check
// @Description Returns service status and timestamp
// @Tags        health
// @Produce     json
// @Success     200 {object} response.HealthResponse
// @Router      /healthz [get]
func (h *Health) Health(c *gin.Context) {
	response.NewResponseBuilder(
		response.WithData("status", "ok"),
	).JSON(c)
}

// Ready checks database connectivity and reports readiness.
// @Summary     Readiness check
// @Description Pings database and returns readiness status
// @Tags        health
// @Produce     json
// @Success     200 {object} response.HealthResponse
// @Router      /readyz [get]
func (h *Health) Ready(c *gin.Context) {
	ctx := c.Request.Context()
	status := "ready"
	if err := h.pool.Ping(ctx); err != nil {
		status = "degraded"
	}

	response.NewResponseBuilder(
		response.WithData("status", status),
	).JSON(c)
}

// Ping echoes pong with a timestamp.
// @Summary     Ping
// @Description Returns pong with timestamp
// @Tags        health
// @Produce     json
// @Success     200 {object} response.MessageResponse
// @Router      /api/v1/ping [get]
func (h *Health) Ping(c *gin.Context) {
	response.NewResponseBuilder(
		response.WithData("message", "pong"),
	).JSON(c)
}

// Pong echoes ping with a timestamp.
// @Summary     Pong
// @Description Returns ping with timestamp
// @Tags        health
// @Produce     json
// @Success     200 {object} response.MessageResponse
// @Router      /api/v1/pong [get]
func (h *Health) Pong(c *gin.Context) {
	response.NewResponseBuilder(
		response.WithData("message", "ping"),
	).JSON(c)
}
