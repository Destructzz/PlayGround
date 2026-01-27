package server

import (
	"backend/internal/http/handlers"
	"backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter configures Gin engine, middleware and routes.
func NewRouter(env string, pool *pgxpool.Pool) *gin.Engine {
	setGinMode(env)

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.RequestLogger(),
		gin.Recovery(),
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/docs", "./static")
	r.StaticFile("/openapi.json", "./docs/swagger.json")

	health := handlers.NewHealth(pool)
	auth := handlers.NewAuth()

	r.GET("/healthz", health.Health)
	r.GET("/readyz", health.Ready)
	api := r.Group("/api/v1")
	api.GET("/auth/:provider", auth.Begin)
	api.GET("/auth/:provider/callback", auth.Callback)
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
