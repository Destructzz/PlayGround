package server

import (
	"backend/internal/domain"
	"backend/internal/http/handlers"
	"backend/internal/http/middleware"
	"backend/internal/repo/sqlc"
	"backend/internal/service"
	"backend/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter configures Gin engine, middleware and routes.
func NewRouter(env string, pool *pgxpool.Pool, queries *sqlc.Queries) *gin.Engine {
	pkg.SetGinMode(env)
	gin.EnableJsonDecoderDisallowUnknownFields()

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.RequestLogger(),
		gin.Recovery(),
		middleware.Cors(),
	)

	tool := handlers.NewToolHandler(queries)

	r.GET("/docs", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), tool.Docs)
	r.GET("/swagger/*any", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/openapi.json", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), tool.GetOpenAPI)

	health := handlers.NewHealth(pool)
	userService := service.NewUserService(queries)
	auth := handlers.NewAuth(userService, queries)
	seedService := service.NewSeed(pool, queries)
	seed := handlers.NewSeed(seedService)
	zoneService := service.NewZone(queries)
	zone := handlers.NewZone(zoneService)
	zoneTagService := service.NewZoneTag(queries)
	zoneTag := handlers.NewZoneTag(zoneTagService)
	configurationService := service.NewComputerConfiguration(queries)
	configuration := handlers.NewComputerConfiguration(configurationService)
	serviceService := service.NewServiceService(queries)
	svc := handlers.NewService(serviceService)
	bookingService := service.NewBooking(queries)
	booking := handlers.NewBooking(bookingService)
	staffService := service.NewStaff(queries)
	staff := handlers.NewStaff(staffService)
	paymentService := service.NewPayment(queries)
	payment := handlers.NewPayment(paymentService)
	public := handlers.NewPublic(queries)

	r.GET("/healthz", health.Health)
	r.GET("/readyz", health.Ready)

	api := r.Group("/api/v1")

	// Auth routes (public)
	api.GET("/auth/:provider", auth.Begin)
	api.GET("/auth/:provider/callback", auth.Callback)
	api.GET("/auth/session", middleware.AuthOptional(queries), auth.Session)
	api.POST("/auth/logout", middleware.AuthOptional(queries), auth.Logout)

	api.GET("/ping", health.Ping)
	api.GET("/pong", health.Pong)
	api.GET("/user", auth.ListUsers)
	api.GET("/seed", seed.Get)
	api.POST("/seed", seed.Post)
	api.DELETE("/seed", seed.Delete)

	// Public catalog routes (no auth required)
	publicScope := api.Group("/public")
	publicScope.GET("/home", public.Home)
	publicScope.GET("/lounge", public.Lounge)
	publicScope.GET("/event", public.Event)
	publicScope.GET("/gaming", public.Gaming)

	// Zone CRUD (admin/internal)
	zoneScope := api.Group("/zone")
	zoneScope.POST("", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), zone.Create)
	zoneScope.GET("", zone.Get)
	zoneScope.GET("/:id", zone.GetById)
	zoneScope.DELETE("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), zone.Delete)
	zoneScope.PATCH("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), zone.Patch)

	zoneTagScope := api.Group("/zone-tag")
	zoneTagScope.POST("", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), zoneTag.Create)
	zoneTagScope.GET("", zoneTag.Get)
	zoneTagScope.GET("/:id", zoneTag.GetByID)
	zoneTagScope.DELETE("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), zoneTag.Delete)
	zoneTagScope.PATCH("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), zoneTag.Patch)

	configurationScope := api.Group("/configuration")
	configurationScope.POST("", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), configuration.Create)
	configurationScope.GET("", configuration.Get)
	configurationScope.GET("/:id", configuration.GetByID)
	configurationScope.DELETE("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), configuration.Delete)
	configurationScope.PATCH("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), configuration.Patch)

	// Service CRUD
	serviceScope := api.Group("/service")
	serviceScope.POST("", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), svc.Create)
	serviceScope.GET("", svc.Get)
	serviceScope.GET("/:id", svc.GetById)
	serviceScope.DELETE("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), svc.Delete)
	serviceScope.PATCH("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), svc.Patch)

	// Booking CRUD (write operations require auth)
	bookingScope := api.Group("/booking")
	bookingScope.POST("", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), booking.Create)
	bookingScope.GET("", booking.Get)
	bookingScope.GET("/:id", booking.GetById)
	bookingScope.DELETE("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), booking.Delete)
	bookingScope.PATCH("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), booking.Patch)

	// Staff CRUD
	staffScope := api.Group("/staff")
	staffScope.POST("", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), staff.Create)
	staffScope.GET("", staff.Get)
	staffScope.GET("/:id", staff.GetById)
	staffScope.DELETE("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), staff.Delete)
	staffScope.PATCH("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), staff.Patch)

	// Payment CRUD
	paymentScope := api.Group("/payment")
	paymentScope.POST("", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), payment.Create)
	paymentScope.GET("", payment.Get)
	paymentScope.GET("/:id", payment.GetById)
	paymentScope.DELETE("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), payment.Delete)
	paymentScope.PATCH("/:id", middleware.AuthRequiredWithRole(queries, domain.RoleAdmin), payment.Patch)

	return r
}
