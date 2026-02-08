package api

import (
	"health-check-service/config"
	"health-check-service/internal/handler"
	"health-check-service/internal/repository/mongodb"
	"health-check-service/internal/scheduler"

	"github.com/gin-gonic/gin"
)

// Router represents the API router
type Router struct {
	engine    *gin.Engine
	cfg       *config.Config
	scheduler *scheduler.Scheduler
}

// NewRouter creates a new router
func NewRouter(cfg *config.Config, scheduler *scheduler.Scheduler) *Router {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.Use(CORSMiddleware())

	return &Router{
		engine:    engine,
		cfg:       cfg,
		scheduler: scheduler,
	}
}

// Setup sets up all routes
func (r *Router) Setup(
	serviceRepo *mongodb.ServiceRepository,
	recipientRepo *mongodb.RecipientRepository,
	healthLogRepo *mongodb.HealthLogRepository,
) {
	// Initialize handlers
	serviceHandler := handler.NewServiceHandler(serviceRepo, r.scheduler, r.cfg)
	recipientHandler := handler.NewRecipientHandler(recipientRepo)
	healthHandler := handler.NewHealthHandler(healthLogRepo, serviceRepo)

	// API routes
	api := r.engine.Group("/api")
	{
		// Services CRUD
		services := api.Group("/services")
		{
			services.GET("", serviceHandler.GetAll)
			services.GET("/:id", serviceHandler.GetByID)
			services.POST("", serviceHandler.Create)
			services.PUT("/:id", serviceHandler.Update)
			services.DELETE("/:id", serviceHandler.Delete)
		}

		// Recipients CRUD
		recipients := api.Group("/recipients")
		{
			recipients.GET("", recipientHandler.GetAll)
			recipients.GET("/:id", recipientHandler.GetByID)
			recipients.POST("", recipientHandler.Create)
			recipients.PUT("/:id", recipientHandler.Update)
			recipients.DELETE("/:id", recipientHandler.Delete)
		}

		// Health logs
		health := api.Group("/health-logs")
		{
			health.GET("", healthHandler.GetRecentLogs)
			health.GET("/:serviceId", healthHandler.GetLogsByServiceID)
		}

		// Dashboard
		api.GET("/dashboard/stats", healthHandler.GetDashboardStats)
	}

	// Serve static files for dashboard
	r.engine.Static("/dashboard", "./web/dashboard")
	r.engine.StaticFile("/", "./web/dashboard/index.html")
}

// Run starts the server
func (r *Router) Run() error {
	return r.engine.Run(":" + r.cfg.AppPort)
}

// Engine returns the gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
