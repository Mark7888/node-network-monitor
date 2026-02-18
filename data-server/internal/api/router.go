package api

import (
	"net/http"
	"time"

	"mark7888/speedtest-data-server/internal/api/handlers"
	"mark7888/speedtest-data-server/internal/api/middleware"
	"mark7888/speedtest-data-server/internal/auth"
	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// SetupRouter configures and returns the Gin router
func SetupRouter(cfg *config.Config, database db.Database, jwtManager *auth.JWTManager) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.API.RateLimit)
	rateLimiter.Cleanup()

	// Health check endpoint (no auth required)
	// Uses SafePing to prevent DDoS attacks by caching ping results
	router.GET("/health", func(c *gin.Context) {
		dbStatus := "connected"
		if err := database.SafePing(); err != nil {
			dbStatus = "disconnected"
		}

		c.JSON(http.StatusOK, models.HealthResponse{
			Status:        "healthy",
			Database:      dbStatus,
			UptimeSeconds: int64(time.Since(startTime).Seconds()),
			Version:       "1.0.0",
		})
	})

	// Create handlers
	nodeHandler := handlers.NewNodeHandler(database)
	measurementHandler := handlers.NewMeasurementHandler(database)
	adminHandler := handlers.NewAdminHandler(database, jwtManager, cfg)
	apiKeyHandler := handlers.NewAPIKeyHandler(database)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Node API (requires API key authentication)
		nodeAPI := v1.Group("/node")
		nodeAPI.Use(middleware.APIKeyAuth(database))
		nodeAPI.Use(middleware.RateLimit(rateLimiter))
		{
			nodeAPI.POST("/alive", nodeHandler.HandleAlive)
		}

		// Measurements API (requires API key authentication)
		measurementsAPI := v1.Group("/measurements")
		measurementsAPI.Use(middleware.APIKeyAuth(database))
		measurementsAPI.Use(middleware.RateLimit(rateLimiter))
		{
			measurementsAPI.POST("", measurementHandler.HandleSubmitMeasurements)
			measurementsAPI.POST("/failed", measurementHandler.HandleSubmitFailedMeasurements)
		}

		// Admin API
		adminAPI := v1.Group("/admin")
		{
			// Login endpoint (no auth required)
			adminAPI.POST("/login", adminHandler.HandleLogin)

			// Protected admin endpoints (require JWT)
			protected := adminAPI.Group("")
			protected.Use(middleware.JWTAuth(jwtManager))
			protected.Use(middleware.RateLimit(rateLimiter))
			{
				// Token refresh
				protected.POST("/refresh", adminHandler.HandleRefreshToken)

				// Dashboard
				protected.GET("/dashboard", adminHandler.HandleGetDashboard)

				// Nodes
				protected.GET("/nodes", adminHandler.HandleListNodes)
				protected.GET("/nodes/:id", adminHandler.HandleGetNodeDetails)
				protected.GET("/nodes/:id/measurements", adminHandler.HandleGetNodeMeasurements)

				// API Keys
				protected.GET("/api-keys", apiKeyHandler.HandleListAPIKeys)
				protected.POST("/api-keys", apiKeyHandler.HandleCreateAPIKey)
				protected.PATCH("/api-keys/:id", apiKeyHandler.HandleUpdateAPIKey)
				protected.DELETE("/api-keys/:id", apiKeyHandler.HandleDeleteAPIKey)
			}
		}

		// Measurements aggregation (admin only)
		measurementsAdminAPI := v1.Group("/admin/measurements")
		measurementsAdminAPI.Use(middleware.JWTAuth(jwtManager))
		measurementsAdminAPI.Use(middleware.RateLimit(rateLimiter))
		{
			measurementsAdminAPI.GET("/aggregate", adminHandler.HandleGetAggregatedMeasurements)
		}
	}

	return router
}
