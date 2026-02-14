package sync

import (
	"encoding/json"
	"fmt"
	"mark7888/speedtest-node/pkg/models"
	"time"

	"go.uber.org/zap"
)

// AliveSender handles sending alive/keepalive signals to the server
type AliveSender struct {
	client   *Client
	nodeID   string
	nodeName string
	logger   *zap.Logger
}

// NewAliveSender creates a new alive sender
func NewAliveSender(client *Client, nodeID, nodeName string, logger *zap.Logger) *AliveSender {
	return &AliveSender{
		client:   client,
		nodeID:   nodeID,
		nodeName: nodeName,
		logger:   logger,
	}
}

// SendAlive sends an alive signal to the server
func (a *AliveSender) SendAlive() error {
	request := &models.AliveRequest{
		NodeID:    a.nodeID,
		NodeName:  a.nodeName,
		Timestamp: time.Now().UTC(),
	}

	respData, err := a.client.Post("/api/v1/node/alive", request)
	if err != nil {
		a.logger.Error("Failed to send alive signal", zap.Error(err))
		return err
	}

	var response models.AliveResponse
	if err := json.Unmarshal(respData, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	a.logger.Debug("Alive signal sent successfully", zap.String("status", response.Status))

	return nil
}
