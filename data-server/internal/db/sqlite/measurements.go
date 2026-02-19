package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InsertMeasurement inserts or updates a measurement
func (s *SQLiteDB) InsertMeasurement(m *models.Measurement) error {
	ctx, cancel := withTimeout()
	defer cancel()

	// Convert boolean to integer for SQLite
	isVPN := 0
	if m.InterfaceIsVPN != nil && *m.InterfaceIsVPN {
		isVPN = 1
	}

	query := `
		INSERT INTO measurements (
			node_id, timestamp, created_at,
			ping_jitter, ping_latency, ping_low, ping_high,
			download_bandwidth, download_bytes, download_elapsed,
			download_latency_iqm, download_latency_low, download_latency_high, download_latency_jitter,
			upload_bandwidth, upload_bytes, upload_elapsed,
			upload_latency_iqm, upload_latency_low, upload_latency_high, upload_latency_jitter,
			packet_loss, isp,
			interface_internal_ip, interface_name, interface_mac, interface_is_vpn, interface_external_ip,
			server_id, server_host, server_port, server_name, server_location, server_country, server_ip,
			result_id, result_url
		) VALUES (
			?, ?, CURRENT_TIMESTAMP,
			?, ?, ?, ?,
			?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?,
			?, ?, ?, ?,
			?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?
		)
		ON CONFLICT (node_id, timestamp) DO UPDATE SET
			ping_jitter = excluded.ping_jitter,
			ping_latency = excluded.ping_latency,
			ping_low = excluded.ping_low,
			ping_high = excluded.ping_high,
			download_bandwidth = excluded.download_bandwidth,
			download_bytes = excluded.download_bytes,
			download_elapsed = excluded.download_elapsed,
			download_latency_iqm = excluded.download_latency_iqm,
			download_latency_low = excluded.download_latency_low,
			download_latency_high = excluded.download_latency_high,
			download_latency_jitter = excluded.download_latency_jitter,
			upload_bandwidth = excluded.upload_bandwidth,
			upload_bytes = excluded.upload_bytes,
			upload_elapsed = excluded.upload_elapsed,
			upload_latency_iqm = excluded.upload_latency_iqm,
			upload_latency_low = excluded.upload_latency_low,
			upload_latency_high = excluded.upload_latency_high,
			upload_latency_jitter = excluded.upload_latency_jitter,
			packet_loss = excluded.packet_loss,
			isp = excluded.isp,
			interface_internal_ip = excluded.interface_internal_ip,
			interface_name = excluded.interface_name,
			interface_mac = excluded.interface_mac,
			interface_is_vpn = excluded.interface_is_vpn,
			interface_external_ip = excluded.interface_external_ip,
			server_id = excluded.server_id,
			server_host = excluded.server_host,
			server_port = excluded.server_port,
			server_name = excluded.server_name,
			server_location = excluded.server_location,
			server_country = excluded.server_country,
			server_ip = excluded.server_ip,
			result_id = excluded.result_id,
			result_url = excluded.result_url
	`

	_, err := s.db.ExecContext(ctx, query,
		m.NodeID.String(), m.Timestamp,
		m.PingJitter, m.PingLatency, m.PingLow, m.PingHigh,
		m.DownloadBandwidth, m.DownloadBytes, m.DownloadElapsed,
		m.DownloadLatencyIqm, m.DownloadLatencyLow, m.DownloadLatencyHigh, m.DownloadLatencyJitter,
		m.UploadBandwidth, m.UploadBytes, m.UploadElapsed,
		m.UploadLatencyIqm, m.UploadLatencyLow, m.UploadLatencyHigh, m.UploadLatencyJitter,
		m.PacketLoss, m.ISP,
		m.InterfaceInternalIP, m.InterfaceName, m.InterfaceMacAddr, isVPN, m.InterfaceExternalIP,
		m.ServerID, m.ServerHost, m.ServerPort, m.ServerName, m.ServerLocation, m.ServerCountry, m.ServerIP,
		m.ResultID, m.ResultURL,
	)

	if err != nil {
		return fmt.Errorf("failed to insert measurement: %w", err)
	}

	return nil
}

// InsertFailedMeasurement inserts a failed measurement record
func (s *SQLiteDB) InsertFailedMeasurement(nodeID uuid.UUID, timestamp time.Time, errorMessage string, retryCount int) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := s.builder.
		Insert("failed_measurements").
		Columns("node_id", "timestamp", "error_message", "retry_count", "created_at").
		Values(nodeID.String(), timestamp, errorMessage, retryCount, time.Now().UTC()).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert failed measurement: %w", err)
	}

	return nil
}

// GetMeasurementsByNode retrieves measurements for a specific node
func (s *SQLiteDB) GetMeasurementsByNode(nodeID uuid.UUID, from, to *time.Time, page, limit int, status string) ([]models.Measurement, int, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Default to "all" if status is empty
	if status == "" {
		status = "all"
	}

	// Build WHERE conditions
	whereConditions := sq.And{sq.Eq{"node_id": nodeID.String()}}
	if from != nil {
		whereConditions = append(whereConditions, sq.GtOrEq{"timestamp": *from})
	}
	if to != nil {
		whereConditions = append(whereConditions, sq.LtOrEq{"timestamp": *to})
	}

	// Get total count based on status
	var total int
	var err error

	if status == "failed" {
		countQuery, countArgs, _ := s.builder.
			Select("COUNT(*)").
			From("failed_measurements").
			Where(whereConditions).
			ToSql()
		err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	} else if status == "successful" {
		countQuery, countArgs, _ := s.builder.
			Select("COUNT(*)").
			From("measurements").
			Where(whereConditions).
			ToSql()
		err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	} else {
		// Count both tables
		countQuery1, countArgs1, _ := s.builder.
			Select("COUNT(*)").
			From("measurements").
			Where(whereConditions).
			ToSql()
		countQuery2, countArgs2, _ := s.builder.
			Select("COUNT(*)").
			From("failed_measurements").
			Where(whereConditions).
			ToSql()

		var count1, count2 int
		_ = s.db.QueryRowContext(ctx, countQuery1, countArgs1...).Scan(&count1)
		_ = s.db.QueryRowContext(ctx, countQuery2, countArgs2...).Scan(&count2)
		total = count1 + count2
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count measurements: %w", err)
	}

	// Build select query based on status
	var rows *sql.Rows
	if status == "failed" {
		selectQuery, selectArgs, _ := s.builder.
			Select("id", "node_id", "timestamp", "created_at", "error_message").
			From("failed_measurements").
			Where(whereConditions).
			OrderBy("timestamp DESC").
			Limit(uint64(limit)).
			Offset(uint64((page - 1) * limit)).
			ToSql()
		rows, err = s.db.QueryContext(ctx, selectQuery, selectArgs...)
	} else if status == "successful" {
		selectQuery, selectArgs, _ := s.builder.
			Select("*").
			From("measurements").
			Where(whereConditions).
			OrderBy("timestamp DESC").
			Limit(uint64(limit)).
			Offset(uint64((page - 1) * limit)).
			ToSql()
		rows, err = s.db.QueryContext(ctx, selectQuery, selectArgs...)
	} else {
		// Union both tables - need to handle this specially
		query := `
			SELECT * FROM (
				SELECT
					id, node_id, timestamp, created_at,
					ping_jitter, ping_latency, ping_low, ping_high,
					download_bandwidth, download_bytes, download_elapsed,
					download_latency_iqm, download_latency_low, download_latency_high, download_latency_jitter,
					upload_bandwidth, upload_bytes, upload_elapsed,
					upload_latency_iqm, upload_latency_low, upload_latency_high, upload_latency_jitter,
					packet_loss, isp,
					interface_internal_ip, interface_name, interface_mac, interface_is_vpn, interface_external_ip,
					server_id, server_host, server_port, server_name, server_location, server_country, server_ip,
					result_id, result_url,
					0 as is_failed,
					NULL as error_message
				FROM measurements
				WHERE ` + s.buildWhereClause(whereConditions) + `
				UNION ALL
				SELECT
					id, node_id, timestamp, created_at,
					NULL, NULL, NULL, NULL,
					NULL, NULL, NULL,
					NULL, NULL, NULL, NULL,
					NULL, NULL, NULL,
					NULL, NULL, NULL, NULL,
					NULL, NULL,
					NULL, NULL, NULL, NULL, NULL,
					NULL, NULL, NULL, NULL, NULL, NULL, NULL,
					NULL, NULL,
					1 as is_failed,
					error_message
				FROM failed_measurements
				WHERE ` + s.buildWhereClause(whereConditions) + `
			) combined
			ORDER BY timestamp DESC
			LIMIT ? OFFSET ?
		`

		args := append(append(s.extractArgs(whereConditions), s.extractArgs(whereConditions)...), limit, (page-1)*limit)
		rows, err = s.db.QueryContext(ctx, query, args...)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to query measurements: %w", err)
	}
	defer rows.Close()

	var measurements []models.Measurement
	for rows.Next() {
		var m models.Measurement
		var nodeIDStr string
		var isFailedInt int
		var isVPNInt *int

		if status == "failed" {
			// For failed measurements, only scan these fields
			err := rows.Scan(
				&m.ID, &nodeIDStr, &m.Timestamp, &m.CreatedAt,
				&m.ErrorMessage,
			)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to scan failed measurement: %w", err)
			}
			m.IsFailed = true
		} else if status == "successful" {
			// For successful measurements, scan all fields
			err := rows.Scan(
				&m.ID, &nodeIDStr, &m.Timestamp, &m.CreatedAt,
				&m.PingJitter, &m.PingLatency, &m.PingLow, &m.PingHigh,
				&m.DownloadBandwidth, &m.DownloadBytes, &m.DownloadElapsed,
				&m.DownloadLatencyIqm, &m.DownloadLatencyLow, &m.DownloadLatencyHigh, &m.DownloadLatencyJitter,
				&m.UploadBandwidth, &m.UploadBytes, &m.UploadElapsed,
				&m.UploadLatencyIqm, &m.UploadLatencyLow, &m.UploadLatencyHigh, &m.UploadLatencyJitter,
				&m.PacketLoss, &m.ISP,
				&m.InterfaceInternalIP, &m.InterfaceName, &m.InterfaceMacAddr, &isVPNInt, &m.InterfaceExternalIP,
				&m.ServerID, &m.ServerHost, &m.ServerPort, &m.ServerName, &m.ServerLocation, &m.ServerCountry, &m.ServerIP,
				&m.ResultID, &m.ResultURL,
			)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to scan measurement: %w", err)
			}

			// Convert integer to boolean for is_vpn
			if isVPNInt != nil {
				isVPN := *isVPNInt == 1
				m.InterfaceIsVPN = &isVPN
			}
			m.IsFailed = false
		} else {
			// Union query - scan all fields
			err := rows.Scan(
				&m.ID, &nodeIDStr, &m.Timestamp, &m.CreatedAt,
				&m.PingJitter, &m.PingLatency, &m.PingLow, &m.PingHigh,
				&m.DownloadBandwidth, &m.DownloadBytes, &m.DownloadElapsed,
				&m.DownloadLatencyIqm, &m.DownloadLatencyLow, &m.DownloadLatencyHigh, &m.DownloadLatencyJitter,
				&m.UploadBandwidth, &m.UploadBytes, &m.UploadElapsed,
				&m.UploadLatencyIqm, &m.UploadLatencyLow, &m.UploadLatencyHigh, &m.UploadLatencyJitter,
				&m.PacketLoss, &m.ISP,
				&m.InterfaceInternalIP, &m.InterfaceName, &m.InterfaceMacAddr, &isVPNInt, &m.InterfaceExternalIP,
				&m.ServerID, &m.ServerHost, &m.ServerPort, &m.ServerName, &m.ServerLocation, &m.ServerCountry, &m.ServerIP,
				&m.ResultID, &m.ResultURL,
				&isFailedInt, &m.ErrorMessage,
			)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to scan measurement: %w", err)
			}

			// Convert integer to boolean
			if isVPNInt != nil {
				isVPN := *isVPNInt == 1
				m.InterfaceIsVPN = &isVPN
			}
			m.IsFailed = isFailedInt == 1
		}

		m.NodeID, _ = uuid.Parse(nodeIDStr)
		measurements = append(measurements, m)
	}

	return measurements, total, nil
}

// Helper functions for building WHERE clauses manually
func (s *SQLiteDB) buildWhereClause(conditions sq.And) string {
	parts := []string{}
	for _, cond := range conditions {
		switch v := cond.(type) {
		case sq.Eq:
			for col := range v {
				parts = append(parts, fmt.Sprintf("%s = ?", col))
			}
		case sq.GtOrEq:
			for col := range v {
				parts = append(parts, fmt.Sprintf("%s >= ?", col))
			}
		case sq.LtOrEq:
			for col := range v {
				parts = append(parts, fmt.Sprintf("%s <= ?", col))
			}
		}
	}
	return strings.Join(parts, " AND ")
}

func (s *SQLiteDB) extractArgs(conditions sq.And) []interface{} {
	args := []interface{}{}
	for _, cond := range conditions {
		switch v := cond.(type) {
		case sq.Eq:
			for _, val := range v {
				args = append(args, val)
			}
		case sq.GtOrEq:
			for _, val := range v {
				args = append(args, val)
			}
		case sq.LtOrEq:
			for _, val := range v {
				args = append(args, val)
			}
		}
	}
	return args
}

// GetAggregatedMeasurements retrieves aggregated measurements for charting
func (s *SQLiteDB) GetAggregatedMeasurements(nodeIDs []uuid.UUID, from, to time.Time, interval string, hideArchived bool) ([]models.AggregatedMeasurement, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get database-specific date truncation function
	truncFunc := getDateTruncSQL(interval)

	// Build WHERE conditions
	args := []interface{}{from, to}
	whereClause := "m.timestamp >= ? AND m.timestamp <= ?"

	if len(nodeIDs) > 0 {
		placeholders := []string{}
		for _, nodeID := range nodeIDs {
			placeholders = append(placeholders, "?")
			args = append(args, nodeID.String())
		}
		whereClause += fmt.Sprintf(" AND m.node_id IN (%s)", strings.Join(placeholders, ","))
	}

	// Add archived filter if requested
	if hideArchived {
		whereClause += " AND n.archived = 0"
	}

	// Build full query with raw SQL for date truncation
	query := fmt.Sprintf(`
		SELECT
			%s as time_bucket,
			m.node_id,
			n.name as node_name,
			COALESCE(AVG(m.download_bandwidth) / 125000.0, 0) as avg_download_mbps,
			COALESCE(AVG(m.upload_bandwidth) / 125000.0, 0) as avg_upload_mbps,
			COALESCE(AVG(m.ping_latency), 0) as avg_ping_ms,
			COALESCE(AVG(m.ping_jitter), 0) as avg_jitter_ms,
			COALESCE(AVG(m.packet_loss), 0) as avg_packet_loss,
			COALESCE(MIN(m.download_bandwidth) / 125000.0, 0) as min_download_mbps,
			COALESCE(MAX(m.download_bandwidth) / 125000.0, 0) as max_download_mbps,
			COUNT(*) as sample_count,
			0 as has_failures
		FROM measurements m
		JOIN nodes n ON m.node_id = n.id
		WHERE %s
		GROUP BY time_bucket, m.node_id, n.name
		ORDER BY time_bucket, m.node_id
	`, truncFunc, whereClause)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query aggregated measurements: %w", err)
	}
	defer rows.Close()

	var results []models.AggregatedMeasurement
	for rows.Next() {
		var agg models.AggregatedMeasurement
		var timeBucketStr string
		var nodeIDStr string
		var hasFailures int

		err := rows.Scan(
			&timeBucketStr,
			&nodeIDStr,
			&agg.NodeName,
			&agg.AvgDownloadMbps,
			&agg.AvgUploadMbps,
			&agg.AvgPingMs,
			&agg.AvgJitterMs,
			&agg.AvgPacketLoss,
			&agg.MinDownloadMbps,
			&agg.MaxDownloadMbps,
			&agg.SampleCount,
			&hasFailures,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan aggregated measurement: %w", err)
		}

		// Parse the string timestamp
		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeBucketStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time_bucket: %w", err)
		}
		agg.Timestamp = parsedTime
		agg.NodeID, _ = uuid.Parse(nodeIDStr)

		results = append(results, agg)
	}

	return results, nil
}

// getDateTruncSQL returns the SQL for date truncation based on interval for SQLite
func getDateTruncSQL(interval string) string {
	switch interval {
	case "5m":
		// Truncate to 5-minute blocks
		return "strftime('%Y-%m-%d %H:', timestamp) || printf('%02d:00', (CAST(strftime('%M', timestamp) AS INTEGER) / 5) * 5)"
	case "15m":
		// Truncate to 15-minute blocks
		return "strftime('%Y-%m-%d %H:', timestamp) || printf('%02d:00', (CAST(strftime('%M', timestamp) AS INTEGER) / 15) * 15)"
	case "1h":
		return "strftime('%Y-%m-%d %H:00:00', timestamp)"
	case "6h":
		// Truncate to 6-hour blocks (00:00, 06:00, 12:00, 18:00)
		return "strftime('%Y-%m-%d', timestamp) || ' ' || printf('%02d:00:00', (CAST(strftime('%H', timestamp) AS INTEGER) / 6) * 6)"
	case "1d":
		return "strftime('%Y-%m-%d 00:00:00', timestamp)"
	default:
		return "strftime('%Y-%m-%d %H:00:00', timestamp)"
	}
}

// GetMeasurementCounts retrieves measurement counts
func (s *SQLiteDB) GetMeasurementCounts() (total int64, last24h int64, lastTimestamp *time.Time, err error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get total count
	countQuery, countArgs, _ := s.builder.
		Select("COUNT(*)").
		From("measurements").
		ToSql()
	err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count measurements: %w", err)
	}

	// Get last 24h count
	past24h := time.Now().UTC().Add(-24 * time.Hour)
	last24hQuery, last24hArgs, _ := s.builder.
		Select("COUNT(*)").
		From("measurements").
		Where(sq.GtOrEq{"created_at": past24h}).
		ToSql()
	err = s.db.QueryRowContext(ctx, last24hQuery, last24hArgs...).Scan(&last24h)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count last 24h measurements: %w", err)
	}

	// Get last measurement timestamp
	var ts time.Time
	lastQuery, lastArgs, _ := s.builder.
		Select("timestamp").
		From("measurements").
		OrderBy("timestamp DESC").
		Limit(1).
		ToSql()
	err = s.db.QueryRowContext(ctx, lastQuery, lastArgs...).Scan(&ts)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, nil, fmt.Errorf("failed to get last measurement timestamp: %w", err)
	}
	if err == nil {
		lastTimestamp = &ts
	}

	return total, last24h, lastTimestamp, nil
}

// GetLast24hStats retrieves average statistics for the last 24 hours
func (s *SQLiteDB) GetLast24hStats() (*models.DashboardStats24h, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	past24h := time.Now().UTC().Add(-24 * time.Hour)

	stats := &models.DashboardStats24h{}

	query := `
		SELECT
			COALESCE(AVG(m.download_bandwidth) / 125000.0, 0) as avg_download_mbps,
			COALESCE(AVG(m.upload_bandwidth) / 125000.0, 0) as avg_upload_mbps,
			COALESCE(AVG(m.ping_latency), 0) as avg_ping_ms,
			COALESCE(AVG(m.ping_jitter), 0) as avg_jitter_ms,
			COALESCE(AVG(m.packet_loss), 0) as avg_packet_loss
		FROM measurements m
		INNER JOIN nodes n ON m.node_id = n.id
		WHERE m.timestamp >= ? AND n.archived = 0
	`

	err := s.db.QueryRowContext(ctx, query, past24h).Scan(
		&stats.DownloadMbps,
		&stats.UploadMbps,
		&stats.PingMs,
		&stats.JitterMs,
		&stats.PacketLoss,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get last 24h stats: %w", err)
	}

	return stats, nil
}

// CleanupOldMeasurements removes measurements older than the retention period
func (s *SQLiteDB) CleanupOldMeasurements(retentionDays int) (int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	cutoffDate := time.Now().UTC().AddDate(0, 0, -retentionDays)

	query, args, err := s.builder.
		Delete("measurements").
		Where(sq.Lt{"timestamp": cutoffDate}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old measurements: %w", err)
	}

	deleted, _ := result.RowsAffected()
	if deleted > 0 {
		logger.Log.Info("Cleaned up old measurements",
			zap.Int64("deleted", deleted),
			zap.Time("cutoff_date", cutoffDate),
		)
	}

	return deleted, nil
}

// CleanupOldFailedMeasurements removes failed measurements older than the retention period
func (s *SQLiteDB) CleanupOldFailedMeasurements(retentionDays int) (int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	cutoffDate := time.Now().UTC().AddDate(0, 0, -retentionDays)

	query, args, err := s.builder.
		Delete("failed_measurements").
		Where(sq.Lt{"timestamp": cutoffDate}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old failed measurements: %w", err)
	}

	deleted, _ := result.RowsAffected()
	if deleted > 0 {
		logger.Log.Info("Cleaned up old failed measurements",
			zap.Int64("deleted", deleted),
			zap.Time("cutoff_date", cutoffDate),
		)
	}

	return deleted, nil
}
