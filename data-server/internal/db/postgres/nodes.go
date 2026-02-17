package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UpsertNode creates or updates a node (used for alive signals and self-registration)
func (p *PostgresDB) UpsertNode(nodeID uuid.UUID, nodeName string) error {
	ctx, cancel := withTimeout()
	defer cancel()

	now := time.Now().UTC()

	// PostgreSQL UPSERT using ON CONFLICT
	query := `
		INSERT INTO nodes (id, name, first_seen, last_seen, last_alive, status)
		VALUES ($1, $2, NOW(), NOW(), NOW(), 'active')
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			last_seen = $3,
			last_alive = $3,
			status = 'active',
			updated_at = $3
	`

	_, err := p.db.ExecContext(ctx, query, nodeID, nodeName, now)
	if err != nil {
		return fmt.Errorf("failed to upsert node: %w", err)
	}

	return nil
}

// GetNodeByID retrieves a node by ID
func (p *PostgresDB) GetNodeByID(nodeID uuid.UUID) (*models.Node, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Select("id", "name", "first_seen", "last_seen", "last_alive", "status", "created_at", "updated_at").
		From("nodes").
		Where(sq.Eq{"id": nodeID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var node models.Node
	err = p.db.QueryRowContext(ctx, query, args...).Scan(
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
func (p *PostgresDB) GetAllNodes(status string, page, limit int) ([]models.Node, int, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Build base query
	selectQuery := p.builder.
		Select("id", "name", "first_seen", "last_seen", "last_alive", "status", "created_at", "updated_at").
		From("nodes")

	countQuery := p.builder.Select("COUNT(*)").From("nodes")

	// Apply status filter if provided
	if status != "" {
		selectQuery = selectQuery.Where(sq.Eq{"status": status})
		countQuery = countQuery.Where(sq.Eq{"status": status})
	}

	// Get total count
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var total int
	err = p.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count nodes: %w", err)
	}

	// Get nodes with pagination
	selectSQL, selectArgs, err := selectQuery.
		OrderBy("last_alive DESC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := p.db.QueryContext(ctx, selectSQL, selectArgs...)
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
func (p *PostgresDB) GetNodeWithStats(nodeID uuid.UUID) (*models.NodeWithStats, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get node
	node, err := p.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}

	nodeWithStats := &models.NodeWithStats{
		Node: *node,
	}

	// Get measurement count
	countQuery, countArgs, err := p.builder.
		Select("COUNT(*)").
		From("measurements").
		Where(sq.Eq{"node_id": nodeID}).
		ToSql()
	if err != nil {
		logger.Log.Warn("Failed to build measurement count query", zap.Error(err))
	} else {
		err = p.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&nodeWithStats.MeasurementCount)
		if err != nil {
			logger.Log.Warn("Failed to get measurement count", zap.Error(err))
		}
	}

	// Get failed test count
	failedCountQuery, failedCountArgs, err := p.builder.
		Select("COUNT(*)").
		From("failed_measurements").
		Where(sq.Eq{"node_id": nodeID}).
		ToSql()
	if err != nil {
		logger.Log.Warn("Failed to build failed count query", zap.Error(err))
	} else {
		err = p.db.QueryRowContext(ctx, failedCountQuery, failedCountArgs...).Scan(&nodeWithStats.FailedTestCount)
		if err != nil {
			logger.Log.Warn("Failed to get failed test count", zap.Error(err))
		}
	}

	// Get statistics
	stats := &models.NodeStatistics{}
	statsQuery := `
		SELECT
			COALESCE(AVG(download_bandwidth) / 125000.0, 0) as avg_download_mbps,
			COALESCE(AVG(upload_bandwidth) / 125000.0, 0) as avg_upload_mbps,
			COALESCE(AVG(ping_latency), 0) as avg_ping_ms,
			COALESCE(AVG(ping_jitter), 0) as avg_jitter_ms,
			COALESCE(AVG(packet_loss), 0) as avg_packet_loss
		FROM measurements
		WHERE node_id = $1
	`
	err = p.db.QueryRowContext(ctx, statsQuery, nodeID).Scan(
		&stats.AvgDownloadMbps,
		&stats.AvgUploadMbps,
		&stats.AvgPingMs,
		&stats.AvgJitterMs,
		&stats.AvgPacketLoss,
	)
	if err != nil {
		logger.Log.Warn("Failed to get node statistics", zap.Error(err))
	}

	// Get success rate for last 24 hours
	past24h := time.Now().UTC().Add(-24 * time.Hour)
	successRateQuery := `
		SELECT
			(SELECT COUNT(*) FROM measurements WHERE node_id = $1 AND timestamp >= $2) as success_count,
			(SELECT COUNT(*) FROM failed_measurements WHERE node_id = $1 AND timestamp >= $2) as failed_count
	`
	err = p.db.QueryRowContext(ctx, successRateQuery, nodeID, past24h).Scan(
		&stats.SuccessCount24h,
		&stats.FailedCount24h,
	)
	if err != nil {
		logger.Log.Warn("Failed to get 24h measurement counts", zap.Error(err))
	} else {
		// Calculate success rate
		totalCount := stats.SuccessCount24h + stats.FailedCount24h
		if totalCount > 0 {
			stats.SuccessRate24h = (float64(stats.SuccessCount24h) / float64(totalCount)) * 100
		} else {
			stats.SuccessRate24h = 0
		}
	}

	nodeWithStats.Statistics = stats

	// Get latest measurement
	latestMeasurement := &models.MeasurementSummary{}
	latestQuery := `
		SELECT
			timestamp,
			COALESCE(download_bandwidth / 125000.0, 0) as download_mbps,
			COALESCE(upload_bandwidth / 125000.0, 0) as upload_mbps,
			COALESCE(ping_latency, 0) as ping_ms
		FROM measurements
		WHERE node_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`
	err = p.db.QueryRowContext(ctx, latestQuery, nodeID).Scan(
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
func (p *PostgresDB) UpdateNodeStatus(aliveTimeout, inactiveTimeout time.Duration) error {
	ctx, cancel := withTimeout()
	defer cancel()

	now := time.Now().UTC()
	unreachableThreshold := now.Add(-aliveTimeout)
	inactiveThreshold := now.Add(-inactiveTimeout)

	// Update to unreachable
	query1 := `
		UPDATE nodes
		SET status = 'unreachable', updated_at = NOW()
		WHERE status = 'active' AND last_alive < $1
	`

	result, err := p.db.ExecContext(ctx, query1, unreachableThreshold)
	if err != nil {
		return fmt.Errorf("failed to update unreachable nodes: %w", err)
	}

	unreachableCount, _ := result.RowsAffected()
	if unreachableCount > 0 {
		logger.Log.Info("Updated nodes to unreachable", zap.Int64("count", unreachableCount))
	}

	// Update to inactive
	query2 := `
		UPDATE nodes
		SET status = 'inactive', updated_at = NOW()
		WHERE status IN ('active', 'unreachable') AND last_alive < $1
	`

	result, err = p.db.ExecContext(ctx, query2, inactiveThreshold)
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
func (p *PostgresDB) GetNodeCounts() (total, active, unreachable, inactive int, err error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'active') as active,
			COUNT(*) FILTER (WHERE status = 'unreachable') as unreachable,
			COUNT(*) FILTER (WHERE status = 'inactive') as inactive
		FROM nodes
	`

	err = p.db.QueryRowContext(ctx, query).Scan(&total, &active, &unreachable, &inactive)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get node counts: %w", err)
	}

	return total, active, unreachable, inactive, nil
}
