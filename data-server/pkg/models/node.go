package models

import (
	"time"

	"github.com/google/uuid"
)

// NodeStatus represents the operational status of a node
type NodeStatus string

const (
	NodeStatusActive      NodeStatus = "active"
	NodeStatusUnreachable NodeStatus = "unreachable"
	NodeStatusInactive    NodeStatus = "inactive"
)

// Node represents a speedtest measurement node
type Node struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	FirstSeen time.Time  `json:"first_seen" db:"first_seen"`
	LastSeen  time.Time  `json:"last_seen" db:"last_seen"`
	LastAlive time.Time  `json:"last_alive" db:"last_alive"`
	Status    NodeStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// NodeWithStats extends Node with statistical information
type NodeWithStats struct {
	Node
	MeasurementCount  int64               `json:"total_measurements,omitempty"`
	FailedTestCount   int64               `json:"failed_test_count,omitempty"`
	LatestMeasurement *MeasurementSummary `json:"latest_measurement,omitempty"`
	Statistics        *NodeStatistics     `json:"statistics,omitempty"`
}

// NodeStatistics contains aggregated statistics for a node
type NodeStatistics struct {
	AvgDownloadMbps float64 `json:"avg_download_mbps"`
	AvgUploadMbps   float64 `json:"avg_upload_mbps"`
	AvgPingMs       float64 `json:"avg_ping_ms"`
	AvgJitterMs     float64 `json:"avg_jitter_ms"`
	AvgPacketLoss   float64 `json:"avg_packet_loss"`
	SuccessRate24h  float64 `json:"success_rate_24h"`  // Success rate for last 24 hours (0-100)
	SuccessCount24h int64   `json:"success_count_24h"` // Successful measurements in last 24h
	FailedCount24h  int64   `json:"failed_count_24h"`  // Failed measurements in last 24h
}

// MeasurementSummary provides a brief summary of a measurement
type MeasurementSummary struct {
	Timestamp    time.Time `json:"timestamp"`
	DownloadMbps float64   `json:"download_mbps"`
	UploadMbps   float64   `json:"upload_mbps"`
	PingMs       float64   `json:"ping_ms"`
}

// AliveRequest represents a node alive/registration request
type AliveRequest struct {
	NodeID    uuid.UUID `json:"node_id" binding:"required"`
	NodeName  string    `json:"node_name" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// AliveResponse represents the response to an alive request
type AliveResponse struct {
	Status         string    `json:"status"`
	ServerTime     time.Time `json:"server_time"`
	NodeRegistered bool      `json:"node_registered"`
}
