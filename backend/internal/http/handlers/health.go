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
func (h *Health) Health(c *gin.Context) {
	response.Health(c, "ok")
}

// Ready checks database connectivity and reports readiness.
func (h *Health) Ready(c *gin.Context) {
	ctx := c.Request.Context()
	status := "ready"
	if err := h.pool.Ping(ctx); err != nil {
		status = "degraded"
	}

	response.Health(c, status)
}

// Ping echoes pong with a timestamp.
func (h *Health) Ping(c *gin.Context) {
	response.Message(c, "pong")
}

// Pong echoes ping with a timestamp.
func (h *Health) Pong(c *gin.Context) {
	response.Message(c, "ping")
}
