package speedtest

import (
	"context"
	"fmt"
	"mark7888/speedtest-node/pkg/models"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

// Executor handles speedtest execution
type Executor struct {
	timeout        time.Duration
	retryOnFailure bool
	logger         *zap.Logger
}

// NewExecutor creates a new speedtest executor
func NewExecutor(timeout time.Duration, retryOnFailure bool, logger *zap.Logger) *Executor {
	return &Executor{
		timeout:        timeout,
		retryOnFailure: retryOnFailure,
		logger:         logger,
	}
}

// Run executes a speedtest and returns the measurement
// If enabled, it will retry once on failure
func (e *Executor) Run() (*models.Measurement, error) {
	e.logger.Info("Starting speedtest")

	// Try to run speedtest
	measurement, err := e.executeSpeedtest()
	if err != nil {
		e.logger.Warn("Speedtest failed", zap.Error(err))

		// Retry if enabled
		if e.retryOnFailure {
			e.logger.Info("Retrying speedtest after 5 seconds")
			time.Sleep(5 * time.Second)

			measurement, err = e.executeSpeedtest()
			if err != nil {
				e.logger.Error("Speedtest failed after retry", zap.Error(err))
				return nil, fmt.Errorf("speedtest failed after retry: %w", err)
			}
		} else {
			return nil, err
		}
	}

	e.logger.Info("Speedtest completed successfully",
		zap.Float64("download_mbps", float64(measurement.Download.Bandwidth)*8/1_000_000),
		zap.Float64("upload_mbps", float64(measurement.Upload.Bandwidth)*8/1_000_000),
		zap.Float64("ping_ms", measurement.Ping.Latency),
		zap.Float64("packet_loss", measurement.PacketLoss),
	)

	return measurement, nil
}

// executeSpeedtest runs the speedtest CLI command
func (e *Executor) executeSpeedtest() (*models.Measurement, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// Run speedtest command
	cmd := exec.CommandContext(ctx, "speedtest", "--format", "json", "--accept-license", "--accept-gdpr")

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("speedtest timeout after %v", e.timeout)
		}
		return nil, fmt.Errorf("speedtest execution failed: %w", err)
	}

	// Parse the output
	measurement, err := ParseResult(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse speedtest result: %w", err)
	}

	return measurement, nil
}
