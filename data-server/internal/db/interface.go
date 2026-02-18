package db

import (
	"time"

	"mark7888/speedtest-data-server/pkg/models"

	"github.com/google/uuid"
)

// Database defines the interface that all database implementations must satisfy
type Database interface {
	// Connection management
	Close() error
	Ping() error
	SafePing() error
	Migrate() error

	// API Keys
	CreateAPIKey(name, plainKey, createdBy string) (*models.APIKey, error)
	GetAPIKeyByID(id uuid.UUID) (*models.APIKey, error)
	GetAllAPIKeys() ([]models.APIKey, error)
	GetEnabledAPIKeys() ([]models.APIKey, error)
	UpdateAPIKeyEnabled(id uuid.UUID, enabled bool) error
	DeleteAPIKey(id uuid.UUID) error
	UpdateAPIKeyLastUsed(id uuid.UUID) error
	VerifyAPIKey(plainKey string) (*models.APIKey, error)

	// Nodes
	UpsertNode(nodeID uuid.UUID, nodeName string) error
	GetNodeByID(nodeID uuid.UUID) (*models.Node, error)
	GetAllNodes(status string, page, limit int) ([]models.Node, int, error)
	GetNodeWithStats(nodeID uuid.UUID) (*models.NodeWithStats, error)
	UpdateNodeStatus(aliveTimeout, inactiveTimeout time.Duration) error
	GetNodeCounts() (total, active, unreachable, inactive int, err error)
	ArchiveNode(nodeID uuid.UUID, archived bool) error
	SetNodeFavorite(nodeID uuid.UUID, favorite bool) error
	DeleteNode(nodeID uuid.UUID) error

	// Measurements
	InsertMeasurement(m *models.Measurement) error
	GetMeasurementsByNode(nodeID uuid.UUID, from, to *time.Time, page, limit int, status string) ([]models.Measurement, int, error)
	InsertFailedMeasurement(nodeID uuid.UUID, timestamp time.Time, errorMessage string, retryCount int) error
	GetAggregatedMeasurements(nodeIDs []uuid.UUID, from, to time.Time, interval string) ([]models.AggregatedMeasurement, error)
	GetMeasurementCounts() (total int64, last24h int64, lastTimestamp *time.Time, err error)
	GetLast24hStats() (*models.DashboardStats24h, error)
	CleanupOldMeasurements(retentionDays int) (int64, error)
	CleanupOldFailedMeasurements(retentionDays int) (int64, error)
}
