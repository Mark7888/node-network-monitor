package services

import (
	"context"
	"time"

	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"

	"go.uber.org/zap"
)

// CleanupService handles periodic data cleanup based on retention policies
type CleanupService struct {
	db     db.Database
	config *config.Config
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCleanupService creates a new cleanup service
func NewCleanupService(database db.Database, cfg *config.Config) *CleanupService {
	ctx, cancel := context.WithCancel(context.Background())
	return &CleanupService{
		db:     database,
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins the cleanup service
func (cs *CleanupService) Start() {
	logger.Log.Info("Starting data cleanup service",
		zap.Duration("cleanup_interval", cs.config.Retention.CleanupInterval),
		zap.Int("measurement_retention_days", cs.config.Retention.MeasurementsDays),
		zap.Int("failed_retention_days", cs.config.Retention.FailedDays),
	)

	ticker := time.NewTicker(cs.config.Retention.CleanupInterval)

	go func() {
		// Run once after a delay to avoid startup load
		time.Sleep(1 * time.Minute)
		cs.runCleanup()

		for {
			select {
			case <-ticker.C:
				cs.runCleanup()
			case <-cs.ctx.Done():
				ticker.Stop()
				logger.Log.Info("Data cleanup service stopped")
				return
			}
		}
	}()
}

// Stop stops the cleanup service
func (cs *CleanupService) Stop() {
	logger.Log.Info("Stopping data cleanup service")
	cs.cancel()
}

// runCleanup performs the actual cleanup
func (cs *CleanupService) runCleanup() {
	logger.Log.Info("Running data cleanup")

	// Cleanup old measurements
	deletedMeasurements, err := cs.db.CleanupOldMeasurements(cs.config.Retention.MeasurementsDays)
	if err != nil {
		logger.Log.Error("Failed to cleanup measurements", zap.Error(err))
	} else if deletedMeasurements > 0 {
		logger.Log.Info("Cleaned up measurements", zap.Int64("deleted", deletedMeasurements))
	}

	// Cleanup old failed measurements
	deletedFailed, err := cs.db.CleanupOldFailedMeasurements(cs.config.Retention.FailedDays)
	if err != nil {
		logger.Log.Error("Failed to cleanup failed measurements", zap.Error(err))
	} else if deletedFailed > 0 {
		logger.Log.Info("Cleaned up failed measurements", zap.Int64("deleted", deletedFailed))
	}

	logger.Log.Info("Data cleanup completed")
}
