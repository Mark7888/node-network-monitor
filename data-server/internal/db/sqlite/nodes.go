package sqlite

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
func (s *SQLiteDB) UpsertNode(nodeID uuid.UUID, nodeName string, nodeLocation *string) error {
	ctx, cancel := withTimeout()
	defer cancel()

	now := time.Now().UTC()

	if nodeLocation != nil {
		// Location explicitly provided: insert/update including the location column.
		// A non-nil pointer to an empty string clears the value; a non-empty string sets it.
		var locationVal interface{}
		if *nodeLocation != "" {
			locationVal = *nodeLocation
		}

		query := `
			INSERT INTO nodes (id, name, location, first_seen, last_seen, last_alive, status)
			VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'active')
			ON CONFLICT (id) DO UPDATE SET
				name = excluded.name,
				location = excluded.location,
				last_seen = ?,
				last_alive = ?,
				status = 'active',
				updated_at = ?
		`
		_, err := s.db.ExecContext(ctx, query, nodeID.String(), nodeName, locationVal, now, now, now)
		if err != nil {
			return fmt.Errorf("failed to upsert node: %w", err)
		}
	} else {
		// No location provided: do not touch the existing location value.
		query := `
			INSERT INTO nodes (id, name, first_seen, last_seen, last_alive, status)
			VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'active')
			ON CONFLICT (id) DO UPDATE SET
				name = excluded.name,
				last_seen = ?,
				last_alive = ?,
				status = 'active',
				updated_at = ?
		`
		_, err := s.db.ExecContext(ctx, query, nodeID.String(), nodeName, now, now, now)
		if err != nil {
			return fmt.Errorf("failed to upsert node: %w", err)
		}
	}

	return nil
}

// GetNodeByID retrieves a node by ID
func (s *SQLiteDB) GetNodeByID(nodeID uuid.UUID) (*models.Node, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := s.builder.
		Select("id", "name", "location", "first_seen", "last_seen", "last_alive", "status", "archived", "favorite", "created_at", "updated_at").
		From("nodes").
		Where(sq.Eq{"id": nodeID.String()}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var node models.Node
	var idStr string
	var archived, favorite int

	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&idStr,
		&node.Name,
		&node.Location,
		&node.FirstSeen,
		&node.LastSeen,
		&node.LastAlive,
		&node.Status,
		&archived,
		&favorite,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("node not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	node.ID, _ = uuid.Parse(idStr)
	node.Archived = archived != 0
	node.Favorite = favorite != 0

	return &node, nil
}

// GetAllNodes retrieves all nodes with optional status filter
func (s *SQLiteDB) GetAllNodes(status string, page, limit int) ([]models.Node, int, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Validate status against known values to prevent unexpected filter behaviour.
	// Empty string means "no filter" (return all nodes).
	if status != "" && status != "active" && status != "unreachable" && status != "inactive" {
		return nil, 0, fmt.Errorf("invalid status %q: must be one of active, unreachable, inactive", status)
	}

	// Build base query
	selectQuery := s.builder.
		Select("id", "name", "location", "first_seen", "last_seen", "last_alive", "status", "archived", "favorite", "created_at", "updated_at").
		From("nodes")

	countQuery := s.builder.Select("COUNT(*)").From("nodes")

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
	err = s.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total)
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

	rows, err := s.db.QueryContext(ctx, selectSQL, selectArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query nodes: %w", err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		var idStr string
		var archived, favorite int

		err := rows.Scan(
			&idStr,
			&node.Name,
			&node.Location,
			&node.FirstSeen,
			&node.LastSeen,
			&node.LastAlive,
			&node.Status,
			&archived,
			&favorite,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan node: %w", err)
		}

		node.ID, _ = uuid.Parse(idStr)
		node.Archived = archived != 0
		node.Favorite = favorite != 0
		nodes = append(nodes, node)
	}

	return nodes, total, nil
}

// GetNodeWithStats retrieves a node with statistics
func (s *SQLiteDB) GetNodeWithStats(nodeID uuid.UUID) (*models.NodeWithStats, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get node
	node, err := s.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}

	nodeWithStats := &models.NodeWithStats{
		Node: *node,
	}

	// Get measurement count
	countQuery, countArgs, err := s.builder.
		Select("COUNT(*)").
		From("measurements").
		Where(sq.Eq{"node_id": nodeID.String()}).
		ToSql()
	if err != nil {
		logger.Log.Warn("Failed to build measurement count query", zap.Error(err))
	} else {
		err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&nodeWithStats.MeasurementCount)
		if err != nil {
			logger.Log.Warn("Failed to get measurement count", zap.Error(err))
		}
	}

	// Get failed test count
	failedCountQuery, failedCountArgs, err := s.builder.
		Select("COUNT(*)").
		From("failed_measurements").
		Where(sq.Eq{"node_id": nodeID.String()}).
		ToSql()
	if err != nil {
		logger.Log.Warn("Failed to build failed count query", zap.Error(err))
	} else {
		err = s.db.QueryRowContext(ctx, failedCountQuery, failedCountArgs...).Scan(&nodeWithStats.FailedTestCount)
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
		WHERE node_id = ?
	`
	err = s.db.QueryRowContext(ctx, statsQuery, nodeID.String()).Scan(
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
			(SELECT COUNT(*) FROM measurements WHERE node_id = ? AND timestamp >= ?) as success_count,
			(SELECT COUNT(*) FROM failed_measurements WHERE node_id = ? AND timestamp >= ?) as failed_count
	`
	err = s.db.QueryRowContext(ctx, successRateQuery, nodeID.String(), past24h, nodeID.String(), past24h).Scan(
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
	latestQuery, latestArgs, err := s.builder.
		Select(
			"timestamp",
			"COALESCE(download_bandwidth / 125000.0, 0) as download_mbps",
			"COALESCE(upload_bandwidth / 125000.0, 0) as upload_mbps",
			"COALESCE(ping_latency, 0) as ping_ms",
		).
		From("measurements").
		Where(sq.Eq{"node_id": nodeID.String()}).
		OrderBy("timestamp DESC").
		Limit(1).
		ToSql()
	if err != nil {
		logger.Log.Warn("Failed to build latest measurement query", zap.Error(err))
	} else {
		err = s.db.QueryRowContext(ctx, latestQuery, latestArgs...).Scan(
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
	}

	return nodeWithStats, nil
}

// UpdateNodeStatus updates the status of nodes based on last_alive timestamp
func (s *SQLiteDB) UpdateNodeStatus(aliveTimeout, inactiveTimeout time.Duration) error {
	ctx, cancel := withTimeout()
	defer cancel()

	now := time.Now().UTC()
	unreachableThreshold := now.Add(-aliveTimeout)
	inactiveThreshold := now.Add(-inactiveTimeout)

	// Update to unreachable
	query1, args1, err := s.builder.
		Update("nodes").
		Set("status", "unreachable").
		Set("updated_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.And{
			sq.Eq{"status": "active"},
			sq.Lt{"last_alive": unreachableThreshold},
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query1, args1...)
	if err != nil {
		return fmt.Errorf("failed to update unreachable nodes: %w", err)
	}

	unreachableCount, _ := result.RowsAffected()
	if unreachableCount > 0 {
		logger.Log.Info("Updated nodes to unreachable", zap.Int64("count", unreachableCount))
	}

	// Update to inactive
	query2, args2, err := s.builder.
		Update("nodes").
		Set("status", "inactive").
		Set("updated_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.And{
			sq.Or{
				sq.Eq{"status": "active"},
				sq.Eq{"status": "unreachable"},
			},
			sq.Lt{"last_alive": inactiveThreshold},
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err = s.db.ExecContext(ctx, query2, args2...)
	if err != nil {
		return fmt.Errorf("failed to update inactive nodes: %w", err)
	}

	inactiveCount, _ := result.RowsAffected()
	if inactiveCount > 0 {
		logger.Log.Info("Updated nodes to inactive", zap.Int64("count", inactiveCount))
	}

	return nil
}

// GetNodeCounts returns counts of nodes by status (excluding archived nodes)
func (s *SQLiteDB) GetNodeCounts() (total, active, unreachable, inactive int, err error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// SQLite doesn't support FILTER, so we use CASE WHEN.
	// COALESCE is required because SUM() returns NULL (not 0) on an empty table,
	// which cannot be scanned into an int.
	query := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0) as active,
			COALESCE(SUM(CASE WHEN status = 'unreachable' THEN 1 ELSE 0 END), 0) as unreachable,
			COALESCE(SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END), 0) as inactive
		FROM nodes
		WHERE archived = 0
	`

	err = s.db.QueryRowContext(ctx, query).Scan(&total, &active, &unreachable, &inactive)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get node counts: %w", err)
	}

	return total, active, unreachable, inactive, nil
}

// ArchiveNode sets the archived status of a node
func (s *SQLiteDB) ArchiveNode(nodeID uuid.UUID, archived bool) error {
	ctx, cancel := withTimeout()
	defer cancel()

	archivedInt := 0
	if archived {
		archivedInt = 1
	}

	query, args, err := s.builder.
		Update("nodes").
		Set("archived", archivedInt).
		Set("updated_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.Eq{"id": nodeID.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to archive node: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("node not found")
	}

	logger.Log.Info("Node archived status updated",
		zap.String("node_id", nodeID.String()),
		zap.Bool("archived", archived),
	)

	return nil
}

// SetNodeFavorite sets the favorite status of a node
func (s *SQLiteDB) SetNodeFavorite(nodeID uuid.UUID, favorite bool) error {
	ctx, cancel := withTimeout()
	defer cancel()

	favoriteInt := 0
	if favorite {
		favoriteInt = 1
	}

	query, args, err := s.builder.
		Update("nodes").
		Set("favorite", favoriteInt).
		Set("updated_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.Eq{"id": nodeID.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to set node favorite: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("node not found")
	}

	logger.Log.Info("Node favorite status updated",
		zap.String("node_id", nodeID.String()),
		zap.Bool("favorite", favorite),
	)

	return nil
}

// DeleteNode deletes a node and all its associated measurements
func (s *SQLiteDB) DeleteNode(nodeID uuid.UUID) error {
	ctx, cancel := withTimeout()
	defer cancel()

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	nodeIDStr := nodeID.String()

	// Delete measurements (cascade should handle this, but explicit is safer)
	_, err = tx.ExecContext(ctx, "DELETE FROM measurements WHERE node_id = ?", nodeIDStr)
	if err != nil {
		return fmt.Errorf("failed to delete measurements: %w", err)
	}

	// Delete failed measurements
	_, err = tx.ExecContext(ctx, "DELETE FROM failed_measurements WHERE node_id = ?", nodeIDStr)
	if err != nil {
		return fmt.Errorf("failed to delete failed measurements: %w", err)
	}

	// Delete node
	result, err := tx.ExecContext(ctx, "DELETE FROM nodes WHERE id = ?", nodeIDStr)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("node not found")
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Log.Info("Node deleted successfully",
		zap.String("node_id", nodeID.String()),
	)

	return nil
}
