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
		corsMiddleware(),
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Static("/docs", "./static")
	r.StaticFile("/openapi.json", "./docs/swagger.json")

	health := handlers.NewHealth(pool)
	userService := service.NewUserService(queries)
	auth := handlers.NewAuth(userService, queries)
	zoneService := service.NewZone(queries)
	zone := handlers.NewZone(zoneService)
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
	api.POST("/auth/dev-login", auth.DevLogin)

	api.GET("/ping", health.Ping)
	api.GET("/pong", health.Pong)
	api.GET("/user", auth.ListUsers)

	// Public catalog routes (no auth required)
	publicScope := api.Group("/public")
	publicScope.GET("/home", public.Home)
	publicScope.GET("/lounge", public.Lounge)
	publicScope.GET("/event", public.Event)
	publicScope.GET("/gaming", public.Gaming)

	// Zone CRUD (admin/internal)
	zoneScope := api.Group("/zone")
	zoneScope.POST("", zone.Create)
	zoneScope.GET("", zone.Get)
	zoneScope.GET("/:id", zone.GetById)
	zoneScope.DELETE("/:id", zone.Delete)
	zoneScope.PATCH("/:id", zone.Patch)

	// Service CRUD
	serviceScope := api.Group("/service")
	serviceScope.POST("", svc.Create)
	serviceScope.GET("", svc.Get)
	serviceScope.GET("/:id", svc.GetById)
	serviceScope.DELETE("/:id", svc.Delete)
	serviceScope.PATCH("/:id", svc.Patch)

	// Booking CRUD (write operations require auth)
	bookingScope := api.Group("/booking")
	bookingScope.POST("", middleware.AuthRequired(queries), booking.Create)
	bookingScope.GET("", booking.Get)
	bookingScope.GET("/:id", booking.GetById)
	bookingScope.DELETE("/:id", middleware.AuthRequired(queries), booking.Delete)
	bookingScope.PATCH("/:id", middleware.AuthRequired(queries), booking.Patch)

	// Staff CRUD
	staffScope := api.Group("/staff")
	staffScope.POST("", staff.Create)
	staffScope.GET("", staff.Get)
	staffScope.GET("/:id", staff.GetById)
	staffScope.DELETE("/:id", staff.Delete)
	staffScope.PATCH("/:id", staff.Patch)

	// Payment CRUD
	paymentScope := api.Group("/payment")
	paymentScope.POST("", payment.Create)
	paymentScope.GET("", payment.Get)
	paymentScope.GET("/:id", payment.GetById)
	paymentScope.DELETE("/:id", payment.Delete)
	paymentScope.PATCH("/:id", payment.Patch)

	return r
}

func setGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		return
	}

	gin.SetMode(gin.DebugMode)
}

// corsMiddleware enables CORS for frontend dev server.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
