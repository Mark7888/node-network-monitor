package services

import (
	"context"
	"time"

	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"

	"go.uber.org/zap"
)

// NodeTracker monitors node status and updates them accordingly
type NodeTracker struct {
	db     db.Database
	config *config.Config
	ctx    context.Context
	cancel context.CancelFunc
}

// NewNodeTracker creates a new node tracker
func NewNodeTracker(database db.Database, cfg *config.Config) *NodeTracker {
	ctx, cancel := context.WithCancel(context.Background())
	return &NodeTracker{
		db:     database,
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins monitoring node status
func (nt *NodeTracker) Start() {
	logger.Log.Info("Starting node status tracker",
		zap.Duration("check_interval", nt.config.Node.StatusCheckInterval),
		zap.Duration("alive_timeout", nt.config.Node.AliveTimeout),
		zap.Duration("inactive_timeout", nt.config.Node.InactiveTimeout),
	)

	// Run once immediately
	nt.checkNodeStatus()

	ticker := time.NewTicker(nt.config.Node.StatusCheckInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				nt.checkNodeStatus()
			case <-nt.ctx.Done():
				ticker.Stop()
				logger.Log.Info("Node status tracker stopped")
				return
			}
		}
	}()
}

// Stop stops the node tracker
func (nt *NodeTracker) Stop() {
	logger.Log.Info("Stopping node status tracker")
	nt.cancel()
}

// checkNodeStatus checks and updates node statuses
func (nt *NodeTracker) checkNodeStatus() {
	err := nt.db.UpdateNodeStatus(nt.config.Node.AliveTimeout, nt.config.Node.InactiveTimeout)
	if err != nil {
		logger.Log.Error("Failed to update node status", zap.Error(err))
	}
}
