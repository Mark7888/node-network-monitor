package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"mark7888/speedtest-node/internal/config"
	"mark7888/speedtest-node/internal/db"
	"mark7888/speedtest-node/internal/logger"
	"mark7888/speedtest-node/internal/scheduler"
	"mark7888/speedtest-node/internal/speedtest"
	"mark7888/speedtest-node/internal/sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Ensure log directory exists
	if cfg.LogOutput != "" {
		logDir := filepath.Dir(cfg.LogOutput)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:         cfg.LogLevel,
		Format:        cfg.LogFormat,
		Output:        cfg.LogOutput,
		OutputConsole: cfg.LogOutputConsole,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting speedtest-node",
		zap.String("node_name", cfg.NodeName),
		zap.String("version", "1.0.0"),
	)

	// Initialize database
	database, err := db.New(cfg.DBPath, log)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.Close()

	// Get or generate node ID
	nodeID, err := database.GetConfig("node_id")
	if err != nil {
		log.Fatal("Failed to get node_id from config", zap.Error(err))
	}

	if nodeID == "" {
		nodeID = uuid.New().String()
		if err := database.SetConfig("node_id", nodeID); err != nil {
			log.Fatal("Failed to store node_id", zap.Error(err))
		}
		log.Info("Generated new node ID", zap.String("node_id", nodeID))
	} else {
		log.Info("Using existing node ID", zap.String("node_id", nodeID))
	}

	// Initialize speedtest executor
	executor := speedtest.NewExecutor(cfg.SpeedtestTimeout, cfg.RetryOnFailure, log)

	// Initialize sync client (only if server URL and API key are provided)
	var sender *sync.Sender
	var aliveSender *sync.AliveSender

	if cfg.ServerURL != "" && cfg.APIKey != "" {
		client := sync.NewClient(cfg.ServerURL, cfg.APIKey, cfg.ServerTimeout, cfg.TLSVerify, log)
		sender = sync.NewSender(client, nodeID, cfg.NodeName, log)
		aliveSender = sync.NewAliveSender(client, nodeID, cfg.NodeName, cfg.NodeLocation, log)
		log.Info("Sync client initialized", zap.String("server_url", cfg.ServerURL))
	} else {
		log.Warn("Server URL or API key not provided, running in offline mode")
	}

	// Initialize scheduler
	sched, err := scheduler.New(
		cfg.SpeedtestCron,
		executor,
		database,
		sender,
		aliveSender,
		cfg.SyncInterval,
		cfg.AliveInterval,
		cfg.RetentionDays,
		cfg.BatchSize,
		log,
	)
	if err != nil {
		log.Fatal("Failed to initialize scheduler", zap.Error(err))
	}

	// Start scheduler
	sched.Start()

	log.Info("Speedtest-node is running")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info("Received shutdown signal")

	// Graceful shutdown
	sched.Stop()

	log.Info("Speedtest-node stopped")
}
