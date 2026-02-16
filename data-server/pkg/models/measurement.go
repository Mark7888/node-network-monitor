package models

import (
	"time"

	"github.com/google/uuid"
)

// Measurement represents a complete speedtest measurement
type Measurement struct {
	ID        int64     `json:"id" db:"id"`
	NodeID    uuid.UUID `json:"node_id" db:"node_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Ping metrics
	PingJitter  *float64 `json:"ping_jitter,omitempty" db:"ping_jitter"`
	PingLatency *float64 `json:"ping_latency,omitempty" db:"ping_latency"`
	PingLow     *float64 `json:"ping_low,omitempty" db:"ping_low"`
	PingHigh    *float64 `json:"ping_high,omitempty" db:"ping_high"`

	// Download metrics
	DownloadBandwidth     *int64   `json:"download_bandwidth,omitempty" db:"download_bandwidth"`
	DownloadBytes         *int64   `json:"download_bytes,omitempty" db:"download_bytes"`
	DownloadElapsed       *int     `json:"download_elapsed,omitempty" db:"download_elapsed"`
	DownloadLatencyIqm    *float64 `json:"download_latency_iqm,omitempty" db:"download_latency_iqm"`
	DownloadLatencyLow    *float64 `json:"download_latency_low,omitempty" db:"download_latency_low"`
	DownloadLatencyHigh   *float64 `json:"download_latency_high,omitempty" db:"download_latency_high"`
	DownloadLatencyJitter *float64 `json:"download_latency_jitter,omitempty" db:"download_latency_jitter"`

	// Upload metrics
	UploadBandwidth     *int64   `json:"upload_bandwidth,omitempty" db:"upload_bandwidth"`
	UploadBytes         *int64   `json:"upload_bytes,omitempty" db:"upload_bytes"`
	UploadElapsed       *int     `json:"upload_elapsed,omitempty" db:"upload_elapsed"`
	UploadLatencyIqm    *float64 `json:"upload_latency_iqm,omitempty" db:"upload_latency_iqm"`
	UploadLatencyLow    *float64 `json:"upload_latency_low,omitempty" db:"upload_latency_low"`
	UploadLatencyHigh   *float64 `json:"upload_latency_high,omitempty" db:"upload_latency_high"`
	UploadLatencyJitter *float64 `json:"upload_latency_jitter,omitempty" db:"upload_latency_jitter"`

	// Network info
	PacketLoss          *float64 `json:"packet_loss,omitempty" db:"packet_loss"`
	ISP                 *string  `json:"isp,omitempty" db:"isp"`
	InterfaceInternalIP *string  `json:"interface_internal_ip,omitempty" db:"interface_internal_ip"`
	InterfaceName       *string  `json:"interface_name,omitempty" db:"interface_name"`
	InterfaceMacAddr    *string  `json:"interface_mac,omitempty" db:"interface_mac"`
	InterfaceIsVPN      *bool    `json:"interface_is_vpn,omitempty" db:"interface_is_vpn"`
	InterfaceExternalIP *string  `json:"interface_external_ip,omitempty" db:"interface_external_ip"`

	// Server info
	ServerID       *int    `json:"server_id,omitempty" db:"server_id"`
	ServerHost     *string `json:"server_host,omitempty" db:"server_host"`
	ServerPort     *int    `json:"server_port,omitempty" db:"server_port"`
	ServerName     *string `json:"server_name,omitempty" db:"server_name"`
	ServerLocation *string `json:"server_location,omitempty" db:"server_location"`
	ServerCountry  *string `json:"server_country,omitempty" db:"server_country"`
	ServerIP       *string `json:"server_ip,omitempty" db:"server_ip"`

	// Result info
	ResultID  *string `json:"result_id,omitempty" db:"result_id"`
	ResultURL *string `json:"result_url,omitempty" db:"result_url"`

	// Failed measurement info
	IsFailed     bool    `json:"is_failed" db:"is_failed"`
	ErrorMessage *string `json:"error_message,omitempty" db:"error_message"`
}

// MeasurementRequest represents the JSON structure from the speedtest node
type MeasurementRequest struct {
	NodeID       uuid.UUID           `json:"node_id" binding:"required"`
	NodeName     string              `json:"node_name" binding:"required"`
	Measurements []MeasurementDetail `json:"measurements" binding:"required,min=1"`
}

// MeasurementDetail represents a single measurement from the JSON
type MeasurementDetail struct {
	Timestamp  time.Time        `json:"timestamp" binding:"required"`
	Ping       *PingMetrics     `json:"ping"`
	Download   *TransferMetrics `json:"download"`
	Upload     *TransferMetrics `json:"upload"`
	PacketLoss *float64         `json:"packet_loss"`
	ISP        *string          `json:"isp"`
	Interface  *InterfaceInfo   `json:"interface"`
	Server     *ServerInfo      `json:"server"`
	Result     *ResultInfo      `json:"result"`
}

// PingMetrics contains ping test results
type PingMetrics struct {
	Jitter  float64 `json:"jitter"`
	Latency float64 `json:"latency"`
	Low     float64 `json:"low"`
	High    float64 `json:"high"`
}

// TransferMetrics contains download/upload test results
type TransferMetrics struct {
	Bandwidth int64           `json:"bandwidth"`
	Bytes     int64           `json:"bytes"`
	Elapsed   int             `json:"elapsed"`
	Latency   *LatencyMetrics `json:"latency"`
}

// LatencyMetrics contains latency details for transfers
type LatencyMetrics struct {
	Iqm    float64 `json:"iqm"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Jitter float64 `json:"jitter"`
}

// InterfaceInfo contains network interface details
type InterfaceInfo struct {
	InternalIP string `json:"internal_ip"`
	Name       string `json:"name"`
	MacAddr    string `json:"mac_addr"`
	IsVPN      bool   `json:"is_vpn"`
	ExternalIP string `json:"external_ip"`
}

// ServerInfo contains speedtest server details
type ServerInfo struct {
	ID       int    `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Country  string `json:"country"`
	IP       string `json:"ip"`
}

// ResultInfo contains speedtest result details
type ResultInfo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// MeasurementResponse represents the response after submitting measurements
type MeasurementResponse struct {
	Status   string `json:"status"`
	Received int    `json:"received"`
	Inserted int    `json:"inserted"`
	Updated  int    `json:"updated"`
	Failed   int    `json:"failed"`
}

// FailedMeasurement represents a failed speedtest attempt
type FailedMeasurement struct {
	ID           int64     `json:"id" db:"id"`
	NodeID       uuid.UUID `json:"node_id" db:"node_id"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
	ErrorMessage *string   `json:"error_message,omitempty" db:"error_message"`
	RetryCount   int       `json:"retry_count" db:"retry_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// FailedMeasurementRequest represents failed test submissions
type FailedMeasurementRequest struct {
	NodeID      uuid.UUID          `json:"node_id" binding:"required"`
	NodeName    string             `json:"node_name" binding:"required"`
	FailedTests []FailedTestDetail `json:"failed_tests" binding:"required,min=1"`
}

// FailedTestDetail represents a single failed test
type FailedTestDetail struct {
	Timestamp    time.Time `json:"timestamp" binding:"required"`
	ErrorMessage string    `json:"error_message" binding:"required"`
	RetryCount   int       `json:"retry_count"`
}

// FailedMeasurementResponse represents the response to failed test submission
type FailedMeasurementResponse struct {
	Status   string `json:"status"`
	Received int    `json:"received"`
}

// AggregatedMeasurement represents aggregated measurement data for charts
type AggregatedMeasurement struct {
	Timestamp       time.Time `json:"timestamp" db:"time_bucket"`
	NodeID          uuid.UUID `json:"node_id" db:"node_id"`
	NodeName        string    `json:"node_name" db:"node_name"`
	AvgDownloadMbps float64   `json:"avg_download_mbps" db:"avg_download_mbps"`
	AvgUploadMbps   float64   `json:"avg_upload_mbps" db:"avg_upload_mbps"`
	AvgPingMs       float64   `json:"avg_ping_ms" db:"avg_ping_ms"`
	AvgJitterMs     float64   `json:"avg_jitter_ms" db:"avg_jitter_ms"`
	AvgPacketLoss   float64   `json:"avg_packet_loss" db:"avg_packet_loss"`
	MinDownloadMbps float64   `json:"min_download_mbps" db:"min_download_mbps"`
	MaxDownloadMbps float64   `json:"max_download_mbps" db:"max_download_mbps"`
	SampleCount     int       `json:"sample_count" db:"sample_count"`
}

// AggregationRequest represents query parameters for aggregated data
type AggregationRequest struct {
	NodeIDs  []uuid.UUID `form:"node_ids"`
	From     time.Time   `form:"from" binding:"required"`
	To       time.Time   `form:"to" binding:"required"`
	Interval string      `form:"interval" binding:"required,oneof=5m 15m 1h 6h 1d"`
}

// AggregationResponse represents the response with aggregated data
type AggregationResponse struct {
	Data         []AggregatedMeasurement `json:"data"`
	Interval     string                  `json:"interval"`
	TotalSamples int                     `json:"total_samples"`
}
