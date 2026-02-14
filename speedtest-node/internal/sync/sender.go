package sync

import (
	"encoding/json"
	"fmt"
	"mark7888/speedtest-node/pkg/models"

	"go.uber.org/zap"
)

// Sender handles sending measurements and failed measurements to the server
type Sender struct {
	client   *Client
	nodeID   string
	nodeName string
	logger   *zap.Logger
}

// NewSender creates a new sender
func NewSender(client *Client, nodeID, nodeName string, logger *zap.Logger) *Sender {
	return &Sender{
		client:   client,
		nodeID:   nodeID,
		nodeName: nodeName,
		logger:   logger,
	}
}

// SendMeasurements sends a batch of measurements to the server
func (s *Sender) SendMeasurements(measurements []*models.Measurement) error {
	if len(measurements) == 0 {
		return nil
	}

	s.logger.Info("Sending measurements to server", zap.Int("count", len(measurements)))

	request := &models.MeasurementsRequest{
		NodeID:       s.nodeID,
		NodeName:     s.nodeName,
		Measurements: measurements,
	}

	respData, err := s.client.Post("/api/v1/measurements", request)
	if err != nil {
		s.logger.Error("Failed to send measurements", zap.Error(err))
		return err
	}

	var response models.MeasurementsResponse
	if err := json.Unmarshal(respData, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	s.logger.Info("Measurements sent successfully",
		zap.Int("received", response.Received),
		zap.Int("failed", response.Failed),
	)

	return nil
}

// SendFailedMeasurements sends a batch of failed measurements to the server
func (s *Sender) SendFailedMeasurements(failed []*models.FailedMeasurement) error {
	if len(failed) == 0 {
		return nil
	}

	s.logger.Info("Sending failed measurements to server", zap.Int("count", len(failed)))

	request := &models.FailedMeasurementsRequest{
		NodeID:      s.nodeID,
		NodeName:    s.nodeName,
		FailedTests: failed,
	}

	respData, err := s.client.Post("/api/v1/measurements/failed", request)
	if err != nil {
		s.logger.Error("Failed to send failed measurements", zap.Error(err))
		return err
	}

	var response models.FailedMeasurementsResponse
	if err := json.Unmarshal(respData, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	s.logger.Info("Failed measurements sent successfully", zap.Int("received", response.Received))

	return nil
}
