package handlers

import (
	"net/http"

	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MeasurementHandler handles measurement-related endpoints
type MeasurementHandler struct {
	db *db.DB
}

// NewMeasurementHandler creates a new measurement handler
func NewMeasurementHandler(database *db.DB) *MeasurementHandler {
	return &MeasurementHandler{
		db: database,
	}
}

// HandleSubmitMeasurements handles measurement submissions from nodes
// POST /api/v1/measurements
func (h *MeasurementHandler) HandleSubmitMeasurements(c *gin.Context) {
	var req models.MeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid measurements request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Ensure node exists
	err := h.db.UpsertNode(req.NodeID, req.NodeName)
	if err != nil {
		logger.Log.Error("Failed to upsert node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to register node",
		})
		return
	}

	// Process each measurement
	received := len(req.Measurements)
	inserted := 0
	updated := 0
	failed := 0

	for _, detail := range req.Measurements {
		measurement := convertToMeasurement(req.NodeID, &detail)
		err := h.db.InsertMeasurement(measurement)
		if err != nil {
			logger.Log.Error("Failed to insert measurement",
				zap.Error(err),
				zap.String("node_id", req.NodeID.String()),
				zap.Time("timestamp", detail.Timestamp),
			)
			failed++
		} else {
			// We can't easily tell if it was inserted or updated without checking beforehand
			// For simplicity, count as inserted
			inserted++
		}
	}

	logger.Log.Info("Measurements processed",
		zap.String("node_id", req.NodeID.String()),
		zap.Int("received", received),
		zap.Int("inserted", inserted),
		zap.Int("failed", failed),
	)

	c.JSON(http.StatusOK, models.MeasurementResponse{
		Status:   "ok",
		Received: received,
		Inserted: inserted,
		Updated:  updated,
		Failed:   failed,
	})
}

// HandleSubmitFailedMeasurements handles failed measurement submissions
// POST /api/v1/measurements/failed
func (h *MeasurementHandler) HandleSubmitFailedMeasurements(c *gin.Context) {
	var req models.FailedMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid failed measurements request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Ensure node exists
	err := h.db.UpsertNode(req.NodeID, req.NodeName)
	if err != nil {
		logger.Log.Error("Failed to upsert node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to register node",
		})
		return
	}

	// Process each failed test
	received := len(req.FailedTests)
	for _, failedTest := range req.FailedTests {
		err := h.db.InsertFailedMeasurement(
			req.NodeID,
			failedTest.Timestamp,
			failedTest.ErrorMessage,
			failedTest.RetryCount,
		)
		if err != nil {
			logger.Log.Error("Failed to insert failed measurement",
				zap.Error(err),
				zap.String("node_id", req.NodeID.String()),
			)
		}
	}

	logger.Log.Info("Failed measurements recorded",
		zap.String("node_id", req.NodeID.String()),
		zap.Int("count", received),
	)

	c.JSON(http.StatusOK, models.FailedMeasurementResponse{
		Status:   "ok",
		Received: received,
	})
}

// convertToMeasurement converts MeasurementDetail to Measurement model
func convertToMeasurement(nodeID uuid.UUID, detail *models.MeasurementDetail) *models.Measurement {
	m := &models.Measurement{
		NodeID:    nodeID,
		Timestamp: detail.Timestamp,
	}

	// Ping
	if detail.Ping != nil {
		m.PingJitter = &detail.Ping.Jitter
		m.PingLatency = &detail.Ping.Latency
		m.PingLow = &detail.Ping.Low
		m.PingHigh = &detail.Ping.High
	}

	// Download
	if detail.Download != nil {
		m.DownloadBandwidth = &detail.Download.Bandwidth
		m.DownloadBytes = &detail.Download.Bytes
		m.DownloadElapsed = &detail.Download.Elapsed

		if detail.Download.Latency != nil {
			m.DownloadLatencyIqm = &detail.Download.Latency.Iqm
			m.DownloadLatencyLow = &detail.Download.Latency.Low
			m.DownloadLatencyHigh = &detail.Download.Latency.High
			m.DownloadLatencyJitter = &detail.Download.Latency.Jitter
		}
	}

	// Upload
	if detail.Upload != nil {
		m.UploadBandwidth = &detail.Upload.Bandwidth
		m.UploadBytes = &detail.Upload.Bytes
		m.UploadElapsed = &detail.Upload.Elapsed

		if detail.Upload.Latency != nil {
			m.UploadLatencyIqm = &detail.Upload.Latency.Iqm
			m.UploadLatencyLow = &detail.Upload.Latency.Low
			m.UploadLatencyHigh = &detail.Upload.Latency.High
			m.UploadLatencyJitter = &detail.Upload.Latency.Jitter
		}
	}

	// Packet loss
	m.PacketLoss = detail.PacketLoss

	// ISP
	m.ISP = detail.ISP

	// Interface
	if detail.Interface != nil {
		m.InterfaceInternalIP = &detail.Interface.InternalIP
		m.InterfaceName = &detail.Interface.Name
		m.InterfaceMacAddr = &detail.Interface.MacAddr
		m.InterfaceIsVPN = &detail.Interface.IsVPN
		m.InterfaceExternalIP = &detail.Interface.ExternalIP
	}

	// Server
	if detail.Server != nil {
		m.ServerID = &detail.Server.ID
		m.ServerHost = &detail.Server.Host
		m.ServerPort = &detail.Server.Port
		m.ServerName = &detail.Server.Name
		m.ServerLocation = &detail.Server.Location
		m.ServerCountry = &detail.Server.Country
		m.ServerIP = &detail.Server.IP
	}

	// Result
	if detail.Result != nil {
		m.ResultID = &detail.Result.ID
		m.ResultURL = &detail.Result.URL
	}

	return m
}
