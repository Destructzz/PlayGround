package server

import (
	"backend/internal/http/handlers"
	"backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRouter configures Gin engine, middleware and routes.
func NewRouter(env string, pool *pgxpool.Pool) *gin.Engine {
	setGinMode(env)

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		gin.Logger(),
		gin.Recovery(),
	)

	health := handlers.NewHealth(pool)

	r.GET("/healthz", health.Health)
	r.GET("/readyz", health.Ready)

	api := r.Group("/api/v1")
	api.GET("/ping", health.Ping)
	api.GET("/pong", health.Pong)

	return r
}

func setGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		return
	}

	gin.SetMode(gin.DebugMode)
}
