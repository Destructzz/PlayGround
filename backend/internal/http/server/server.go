package server

import (
	"backend/internal/http/handlers"
	"backend/internal/http/middleware"
	"backend/internal/repo/sqlc"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter configures Gin engine, middleware and routes.
func NewRouter(env string, pool *pgxpool.Pool, queries *sqlc.Queries) *gin.Engine {
	setGinMode(env)
	gin.EnableJsonDecoderDisallowUnknownFields()

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
	userService := service.NewUserService(queries)
	auth := handlers.NewAuth(userService)
	zoneService := service.NewZone(queries)
	zone := handlers.NewZone(zoneService)
	serviceService := service.NewServiceService(queries)
	svc := handlers.NewService(serviceService)
	bookingService := service.NewBooking(queries, serviceService)
	booking := handlers.NewBooking(bookingService)

	r.GET("/healthz", health.Health)
	r.GET("/readyz", health.Ready)

	api := r.Group("/api/v1")
	api.GET("/auth/:provider", auth.Begin)
	api.GET("/auth/:provider/callback", auth.Callback)
	api.GET("/ping", health.Ping)
	api.GET("/pong", health.Pong)

	zoneScope := api.Group("/zone")

	zoneScope.POST("", zone.Create)
	zoneScope.GET("", zone.Get)
	zoneScope.GET("/:id", zone.GetById)
	zoneScope.DELETE("/:id", zone.Delete)
	zoneScope.PATCH("/:id", zone.Patch)

	serviceScope := api.Group("/service")

	serviceScope.POST("", svc.Create)
	serviceScope.GET("", svc.Get)
	serviceScope.GET("/:id", svc.GetById)
	serviceScope.DELETE("/:id", svc.Delete)
	serviceScope.PATCH("/:id", svc.Patch)

	bookingScope := api.Group("/booking")

	bookingScope.POST("", booking.Create)
	bookingScope.GET("", booking.Get)
	bookingScope.GET("/:id", booking.GetById)
	bookingScope.DELETE("/:id", booking.Delete)
	bookingScope.PATCH("/:id", booking.Patch)

	return r
}

func setGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		return
	}

	gin.SetMode(gin.DebugMode)
}
