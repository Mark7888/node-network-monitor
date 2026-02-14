package models

import "time"

// LoginRequest represents admin login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Username  string    `json:"username"`
}

// RefreshTokenResponse represents the response after token refresh
type RefreshTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// DashboardSummary represents the dashboard summary data
type DashboardSummary struct {
	TotalNodes          int                `json:"total_nodes"`
	ActiveNodes         int                `json:"active_nodes"`
	UnreachableNodes    int                `json:"unreachable_nodes"`
	TotalMeasurements   int64              `json:"total_measurements"`
	MeasurementsLast24h int64              `json:"measurements_last_24h"`
	LastMeasurement     *time.Time         `json:"last_measurement,omitempty"`
	AverageStats24h     *DashboardStats24h `json:"average_stats_24h,omitempty"`
}

// DashboardStats24h represents average statistics for the last 24 hours
type DashboardStats24h struct {
	DownloadMbps float64 `json:"download_mbps"`
	UploadMbps   float64 `json:"upload_mbps"`
	PingMs       float64 `json:"ping_ms"`
	JitterMs     float64 `json:"jitter_ms"`
	PacketLoss   float64 `json:"packet_loss"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status        string `json:"status"`
	Database      string `json:"database"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Version       string `json:"version"`
}
