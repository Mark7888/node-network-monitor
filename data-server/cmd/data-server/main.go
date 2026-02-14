package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mark7888/speedtest-data-server/internal/api"
	"mark7888/speedtest-data-server/internal/auth"
	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/internal/services"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Initialize(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Log.Info("Starting data server",
		zap.String("version", "1.0.0"),
		zap.String("mode", cfg.Server.Mode),
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// Connect to database
	database, err := db.New(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		logger.Log.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Hash admin password if not already hashed (for first run)
	if len(cfg.Admin.Password) < 60 { // bcrypt hashes are 60 characters
		hashedPassword, err := auth.HashPassword(cfg.Admin.Password)
		if err != nil {
			logger.Log.Fatal("Failed to hash admin password", zap.Error(err))
		}
		cfg.Admin.Password = hashedPassword
		logger.Log.Info("Admin password hashed successfully")
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiry)

	// Setup router
	router := api.SetupRouter(cfg, database, jwtManager)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.API.Timeout,
		WriteTimeout: cfg.API.Timeout,
		IdleTimeout:  2 * cfg.API.Timeout,
	}

	// Start background services
	nodeTracker := services.NewNodeTracker(database, cfg)
	nodeTracker.Start()
	defer nodeTracker.Stop()

	cleanupService := services.NewCleanupService(database, cfg)
	cleanupService.Start()
	defer cleanupService.Stop()

	// Start server in a goroutine
	go func() {
		logger.Log.Info("Server listening", zap.String("addr", addr))

		var err error
		if cfg.Server.TLSEnabled {
			logger.Log.Info("Starting server with TLS",
				zap.String("cert", cfg.Server.TLSCert),
				zap.String("key", cfg.Server.TLSKey),
			)
			err = server.ListenAndServeTLS(cfg.Server.TLSCert, cfg.Server.TLSKey)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited")
}
