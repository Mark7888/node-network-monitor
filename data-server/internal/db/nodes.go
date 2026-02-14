package db

import (
	"database/sql"
	"fmt"
	"time"

	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UpsertNode creates or updates a node (used for alive signals and self-registration)
func (db *DB) UpsertNode(nodeID uuid.UUID, nodeName string) error {
	ctx, cancel := withTimeout()
	defer cancel()

	now := time.Now().UTC()
	nowSQL := db.getNowSQL()
	query := fmt.Sprintf(`
		INSERT INTO nodes (id, name, first_seen, last_seen, last_alive, status)
		VALUES ($1, $2, %s, %s, %s, 'active')
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			last_seen = $3,
			last_alive = $3,
			status = 'active',
			updated_at = $3
	`, nowSQL, nowSQL, nowSQL)

	_, err := db.ExecContext(ctx, query, nodeID, nodeName, now)
	if err != nil {
		return fmt.Errorf("failed to upsert node: %w", err)
	}

	return nil
}

// GetNodeByID retrieves a node by ID
func (db *DB) GetNodeByID(nodeID uuid.UUID) (*models.Node, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query := `
		SELECT id, name, first_seen, last_seen, last_alive, status, created_at, updated_at
		FROM nodes
		WHERE id = $1
	`

	var node models.Node
	err := db.QueryRowContext(ctx, query, nodeID).Scan(
		&node.ID,
		&node.Name,
		&node.FirstSeen,
		&node.LastSeen,
		&node.LastAlive,
		&node.Status,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("node not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	return &node, nil
}

// GetAllNodes retrieves all nodes with optional status filter
func (db *DB) GetAllNodes(status string, page, limit int) ([]models.Node, int, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Build query based on status filter
	var countQuery, selectQuery string
	var args []interface{}

	if status != "" {
		countQuery = "SELECT COUNT(*) FROM nodes WHERE status = $1"
		selectQuery = `
			SELECT id, name, first_seen, last_seen, last_alive, status, created_at, updated_at
			FROM nodes
			WHERE status = $1
			ORDER BY last_alive DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{status, limit, (page - 1) * limit}
	} else {
		countQuery = "SELECT COUNT(*) FROM nodes"
		selectQuery = `
			SELECT id, name, first_seen, last_seen, last_alive, status, created_at, updated_at
			FROM nodes
			ORDER BY last_alive DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, (page - 1) * limit}
	}

	// Get total count
	var total int
	var err error
	if status != "" {
		err = db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	} else {
		err = db.QueryRowContext(ctx, countQuery).Scan(&total)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count nodes: %w", err)
	}

	// Get nodes
	rows, err := db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query nodes: %w", err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(
			&node.ID,
			&node.Name,
			&node.FirstSeen,
			&node.LastSeen,
			&node.LastAlive,
			&node.Status,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, node)
	}

	return nodes, total, nil
}

// GetNodeWithStats retrieves a node with statistics
func (db *DB) GetNodeWithStats(nodeID uuid.UUID) (*models.NodeWithStats, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get node
	node, err := db.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}

	nodeWithStats := &models.NodeWithStats{
		Node: *node,
	}

	// Get measurement count
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM measurements WHERE node_id = $1",
		nodeID,
	).Scan(&nodeWithStats.MeasurementCount)
	if err != nil {
		logger.Log.Warn("Failed to get measurement count", zap.Error(err))
	}

	// Get failed test count
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM failed_measurements WHERE node_id = $1",
		nodeID,
	).Scan(&nodeWithStats.FailedTestCount)
	if err != nil {
		logger.Log.Warn("Failed to get failed test count", zap.Error(err))
	}

	// Get statistics
	stats := &models.NodeStatistics{}
	err = db.QueryRowContext(ctx, `
		SELECT
			COALESCE(AVG(download_bandwidth) / 125000.0, 0) as avg_download_mbps,
			COALESCE(AVG(upload_bandwidth) / 125000.0, 0) as avg_upload_mbps,
			COALESCE(AVG(ping_latency), 0) as avg_ping_ms,
			COALESCE(AVG(ping_jitter), 0) as avg_jitter_ms,
			COALESCE(AVG(packet_loss), 0) as avg_packet_loss
		FROM measurements
		WHERE node_id = $1
	`, nodeID).Scan(
		&stats.AvgDownloadMbps,
		&stats.AvgUploadMbps,
		&stats.AvgPingMs,
		&stats.AvgJitterMs,
		&stats.AvgPacketLoss,
	)
	if err != nil {
		logger.Log.Warn("Failed to get node statistics", zap.Error(err))
	} else {
		nodeWithStats.Statistics = stats
	}

	// Get latest measurement
	latestMeasurement := &models.MeasurementSummary{}
	err = db.QueryRowContext(ctx, `
		SELECT
			timestamp,
			COALESCE(download_bandwidth / 125000.0, 0) as download_mbps,
			COALESCE(upload_bandwidth / 125000.0, 0) as upload_mbps,
			COALESCE(ping_latency, 0) as ping_ms
		FROM measurements
		WHERE node_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`, nodeID).Scan(
		&latestMeasurement.Timestamp,
		&latestMeasurement.DownloadMbps,
		&latestMeasurement.UploadMbps,
		&latestMeasurement.PingMs,
	)
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Warn("Failed to get latest measurement", zap.Error(err))
	} else if err == nil {
		nodeWithStats.LatestMeasurement = latestMeasurement
	}

	return nodeWithStats, nil
}

// UpdateNodeStatus updates the status of nodes based on last_alive timestamp
func (db *DB) UpdateNodeStatus(aliveTimeout, inactiveTimeout time.Duration) error {
	ctx, cancel := withTimeout()
	defer cancel()

	now := time.Now().UTC()
	unreachableThreshold := now.Add(-aliveTimeout)
	inactiveThreshold := now.Add(-inactiveTimeout)

	nowSQL := db.getNowSQL()

	// Update to unreachable
	query1 := fmt.Sprintf(`
		UPDATE nodes
		SET status = 'unreachable', updated_at = %s
		WHERE status = 'active' AND last_alive < $1
	`, nowSQL)

	result, err := db.ExecContext(ctx, query1, unreachableThreshold)
	if err != nil {
		return fmt.Errorf("failed to update unreachable nodes: %w", err)
	}

	unreachableCount, _ := result.RowsAffected()
	if unreachableCount > 0 {
		logger.Log.Info("Updated nodes to unreachable", zap.Int64("count", unreachableCount))
	}

	// Update to inactive
	query2 := fmt.Sprintf(`
		UPDATE nodes
		SET status = 'inactive', updated_at = %s
		WHERE status IN ('active', 'unreachable') AND last_alive < $1
	`, nowSQL)

	result, err = db.ExecContext(ctx, query2, inactiveThreshold)
	if err != nil {
		return fmt.Errorf("failed to update inactive nodes: %w", err)
	}

	inactiveCount, _ := result.RowsAffected()
	if inactiveCount > 0 {
		logger.Log.Info("Updated nodes to inactive", zap.Int64("count", inactiveCount))
	}

	return nil
}

// GetNodeCounts returns counts of nodes by status
func (db *DB) GetNodeCounts() (total, active, unreachable, inactive int, err error) {
	ctx, cancel := withTimeout()
	defer cancel()

	err = db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'active') as active,
			COUNT(*) FILTER (WHERE status = 'unreachable') as unreachable,
			COUNT(*) FILTER (WHERE status = 'inactive') as inactive
		FROM nodes
	`).Scan(&total, &active, &unreachable, &inactive)

	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get node counts: %w", err)
	}

	return total, active, unreachable, inactive, nil
}
