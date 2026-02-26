package postgres

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
func (p *PostgresDB) InsertMeasurement(m *models.Measurement) error {
	ctx, cancel := withTimeout()
	defer cancel()

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
			$1, $2, NOW(),
			$3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16,
			$17, $18, $19, $20,
			$21, $22,
			$23, $24, $25, $26, $27,
			$28, $29, $30, $31, $32, $33, $34,
			$35, $36
		)
		ON CONFLICT (node_id, timestamp) DO UPDATE SET
			ping_jitter = EXCLUDED.ping_jitter,
			ping_latency = EXCLUDED.ping_latency,
			ping_low = EXCLUDED.ping_low,
			ping_high = EXCLUDED.ping_high,
			download_bandwidth = EXCLUDED.download_bandwidth,
			download_bytes = EXCLUDED.download_bytes,
			download_elapsed = EXCLUDED.download_elapsed,
			download_latency_iqm = EXCLUDED.download_latency_iqm,
			download_latency_low = EXCLUDED.download_latency_low,
			download_latency_high = EXCLUDED.download_latency_high,
			download_latency_jitter = EXCLUDED.download_latency_jitter,
			upload_bandwidth = EXCLUDED.upload_bandwidth,
			upload_bytes = EXCLUDED.upload_bytes,
			upload_elapsed = EXCLUDED.upload_elapsed,
			upload_latency_iqm = EXCLUDED.upload_latency_iqm,
			upload_latency_low = EXCLUDED.upload_latency_low,
			upload_latency_high = EXCLUDED.upload_latency_high,
			upload_latency_jitter = EXCLUDED.upload_latency_jitter,
			packet_loss = EXCLUDED.packet_loss,
			isp = EXCLUDED.isp,
			interface_internal_ip = EXCLUDED.interface_internal_ip,
			interface_name = EXCLUDED.interface_name,
			interface_mac = EXCLUDED.interface_mac,
			interface_is_vpn = EXCLUDED.interface_is_vpn,
			interface_external_ip = EXCLUDED.interface_external_ip,
			server_id = EXCLUDED.server_id,
			server_host = EXCLUDED.server_host,
			server_port = EXCLUDED.server_port,
			server_name = EXCLUDED.server_name,
			server_location = EXCLUDED.server_location,
			server_country = EXCLUDED.server_country,
			server_ip = EXCLUDED.server_ip,
			result_id = EXCLUDED.result_id,
			result_url = EXCLUDED.result_url
	`

	_, err := p.db.ExecContext(ctx, query,
		m.NodeID, m.Timestamp,
		m.PingJitter, m.PingLatency, m.PingLow, m.PingHigh,
		m.DownloadBandwidth, m.DownloadBytes, m.DownloadElapsed,
		m.DownloadLatencyIqm, m.DownloadLatencyLow, m.DownloadLatencyHigh, m.DownloadLatencyJitter,
		m.UploadBandwidth, m.UploadBytes, m.UploadElapsed,
		m.UploadLatencyIqm, m.UploadLatencyLow, m.UploadLatencyHigh, m.UploadLatencyJitter,
		m.PacketLoss, m.ISP,
		m.InterfaceInternalIP, m.InterfaceName, m.InterfaceMacAddr, m.InterfaceIsVPN, m.InterfaceExternalIP,
		m.ServerID, m.ServerHost, m.ServerPort, m.ServerName, m.ServerLocation, m.ServerCountry, m.ServerIP,
		m.ResultID, m.ResultURL,
	)

	if err != nil {
		return fmt.Errorf("failed to insert measurement: %w", err)
	}

	return nil
}

// InsertFailedMeasurement inserts a failed measurement record
func (p *PostgresDB) InsertFailedMeasurement(nodeID uuid.UUID, timestamp time.Time, errorMessage string, retryCount int) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Insert("failed_measurements").
		Columns("node_id", "timestamp", "error_message", "retry_count", "created_at").
		Values(nodeID, timestamp, errorMessage, retryCount, sq.Expr("NOW()")).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert failed measurement: %w", err)
	}

	return nil
}

// GetMeasurementsByNode retrieves measurements for a specific node
func (p *PostgresDB) GetMeasurementsByNode(nodeID uuid.UUID, from, to *time.Time, page, limit int, status string) ([]models.Measurement, int, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Default to "all" if status is empty
	if status == "" {
		status = "all"
	}

	// Build WHERE conditions
	whereConditions := sq.And{sq.Eq{"node_id": nodeID}}
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
		countQuery, countArgs, _ := p.builder.
			Select("COUNT(*)").
			From("failed_measurements").
			Where(whereConditions).
			ToSql()
		err = p.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	} else if status == "successful" {
		countQuery, countArgs, _ := p.builder.
			Select("COUNT(*)").
			From("measurements").
			Where(whereConditions).
			ToSql()
		err = p.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	} else {
		// Count both tables
		countQuery1, countArgs1, _ := p.builder.
			Select("COUNT(*)").
			From("measurements").
			Where(whereConditions).
			ToSql()
		countQuery2, countArgs2, _ := p.builder.
			Select("COUNT(*)").
			From("failed_measurements").
			Where(whereConditions).
			ToSql()

		var count1, count2 int
		_ = p.db.QueryRowContext(ctx, countQuery1, countArgs1...).Scan(&count1)
		_ = p.db.QueryRowContext(ctx, countQuery2, countArgs2...).Scan(&count2)
		total = count1 + count2
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count measurements: %w", err)
	}

	// Build select query based on status
	var rows *sql.Rows
	if status == "failed" {
		query := `
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
				true as is_failed,
				error_message
			FROM failed_measurements
			WHERE ` + p.buildWhereClause(whereConditions) + `
			ORDER BY timestamp DESC
			LIMIT $` + fmt.Sprintf("%d", len(p.extractArgs(whereConditions))+1) +
			` OFFSET $` + fmt.Sprintf("%d", len(p.extractArgs(whereConditions))+2)

		args := append(p.extractArgs(whereConditions), limit, (page-1)*limit)
		rows, err = p.db.QueryContext(ctx, query, args...)
	} else if status == "successful" {
		query := `
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
				false as is_failed,
				NULL as error_message
			FROM measurements
			WHERE ` + p.buildWhereClause(whereConditions) + `
			ORDER BY timestamp DESC
			LIMIT $` + fmt.Sprintf("%d", len(p.extractArgs(whereConditions))+1) +
			` OFFSET $` + fmt.Sprintf("%d", len(p.extractArgs(whereConditions))+2)

		args := append(p.extractArgs(whereConditions), limit, (page-1)*limit)
		rows, err = p.db.QueryContext(ctx, query, args...)
	} else {
		// Union both tables
		numArgs := len(p.extractArgs(whereConditions))
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
					false as is_failed,
					NULL as error_message
				FROM measurements
				WHERE ` + p.buildWhereClause(whereConditions) + `
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
					true as is_failed,
					error_message
				FROM failed_measurements
				WHERE ` + p.buildWhereClauseWithOffset(whereConditions, numArgs) + `
			) combined
			ORDER BY timestamp DESC
			LIMIT $` + fmt.Sprintf("%d", numArgs*2+1) +
			` OFFSET $` + fmt.Sprintf("%d", numArgs*2+2)

		args := append(append(p.extractArgs(whereConditions), p.extractArgs(whereConditions)...), limit, (page-1)*limit)
		rows, err = p.db.QueryContext(ctx, query, args...)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to query measurements: %w", err)
	}
	defer rows.Close()

	var measurements []models.Measurement
	for rows.Next() {
		var m models.Measurement
		err := rows.Scan(
			&m.ID, &m.NodeID, &m.Timestamp, &m.CreatedAt,
			&m.PingJitter, &m.PingLatency, &m.PingLow, &m.PingHigh,
			&m.DownloadBandwidth, &m.DownloadBytes, &m.DownloadElapsed,
			&m.DownloadLatencyIqm, &m.DownloadLatencyLow, &m.DownloadLatencyHigh, &m.DownloadLatencyJitter,
			&m.UploadBandwidth, &m.UploadBytes, &m.UploadElapsed,
			&m.UploadLatencyIqm, &m.UploadLatencyLow, &m.UploadLatencyHigh, &m.UploadLatencyJitter,
			&m.PacketLoss, &m.ISP,
			&m.InterfaceInternalIP, &m.InterfaceName, &m.InterfaceMacAddr, &m.InterfaceIsVPN, &m.InterfaceExternalIP,
			&m.ServerID, &m.ServerHost, &m.ServerPort, &m.ServerName, &m.ServerLocation, &m.ServerCountry, &m.ServerIP,
			&m.ResultID, &m.ResultURL,
			&m.IsFailed, &m.ErrorMessage,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan measurement: %w", err)
		}
		measurements = append(measurements, m)
	}

	return measurements, total, nil
}

// Helper functions for building WHERE clauses manually
func (p *PostgresDB) buildWhereClause(conditions sq.And) string {
	return p.buildWhereClauseWithOffset(conditions, 0)
}

func (p *PostgresDB) buildWhereClauseWithOffset(conditions sq.And, offset int) string {
	parts := []string{}
	argNum := offset + 1
	for _, cond := range conditions {
		switch v := cond.(type) {
		case sq.Eq:
			for col := range v {
				parts = append(parts, fmt.Sprintf("%s = $%d", col, argNum))
				argNum++
			}
		case sq.GtOrEq:
			for col := range v {
				parts = append(parts, fmt.Sprintf("%s >= $%d", col, argNum))
				argNum++
			}
		case sq.LtOrEq:
			for col := range v {
				parts = append(parts, fmt.Sprintf("%s <= $%d", col, argNum))
				argNum++
			}
		}
	}
	return strings.Join(parts, " AND ")
}

func (p *PostgresDB) extractArgs(conditions sq.And) []interface{} {
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
func (p *PostgresDB) GetAggregatedMeasurements(nodeIDs []uuid.UUID, from, to time.Time, interval string, hideArchived bool) ([]models.AggregatedMeasurement, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get database-specific date truncation function
	truncFunc, err := getDateTruncSQL(interval)
	if err != nil {
		return nil, fmt.Errorf("GetAggregatedMeasurements: %w", err)
	}

	// Build WHERE conditions
	whereBuilder := p.builder.
		Select().
		Where(sq.GtOrEq{"timestamp": from}).
		Where(sq.LtOrEq{"timestamp": to})

	if len(nodeIDs) > 0 {
		whereBuilder = whereBuilder.Where(sq.Eq{"m.node_id": nodeIDs})
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
			false as has_failures
		FROM measurements m
		JOIN nodes n ON m.node_id = n.id
		WHERE m.timestamp >= $1 AND m.timestamp <= $2
	`, truncFunc)

	args := []interface{}{from, to}
	argPos := 3

	if len(nodeIDs) > 0 {
		placeholders := []string{}
		for _, nodeID := range nodeIDs {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argPos))
			args = append(args, nodeID)
			argPos++
		}
		query += fmt.Sprintf(" AND m.node_id IN (%s)", strings.Join(placeholders, ","))
	}

	// Add archived filter if requested
	if hideArchived {
		query += " AND n.archived = false"
	}

	query += `
		GROUP BY time_bucket, m.node_id, n.name
		ORDER BY time_bucket, m.node_id
	`

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query aggregated measurements: %w", err)
	}
	defer rows.Close()

	var results []models.AggregatedMeasurement
	for rows.Next() {
		var agg models.AggregatedMeasurement
		var hasFailures bool

		err := rows.Scan(
			&agg.Timestamp,
			&agg.NodeID,
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

		results = append(results, agg)
	}

	return results, nil
}

// getDateTruncSQL returns the SQL for date truncation based on interval.
// Returns an error for unrecognised intervals to prevent unvalidated strings from
// being interpolated into a raw SQL query.
func getDateTruncSQL(interval string) (string, error) {
	switch interval {
	case "5m":
		return "date_trunc('hour', timestamp) + INTERVAL '5 min' * floor(EXTRACT(MINUTE FROM timestamp)::int / 5)", nil
	case "15m":
		return "date_trunc('hour', timestamp) + INTERVAL '15 min' * floor(EXTRACT(MINUTE FROM timestamp)::int / 15)", nil
	case "1h":
		return "date_trunc('hour', timestamp)", nil
	case "6h":
		return "date_trunc('hour', timestamp) - (EXTRACT(HOUR FROM timestamp)::int % 6) * INTERVAL '1 hour'", nil
	case "1d":
		return "date_trunc('day', timestamp)", nil
	default:
		return "", fmt.Errorf("invalid interval %q: must be one of 5m, 15m, 1h, 6h, 1d", interval)
	}
}

// GetMeasurementCounts retrieves measurement counts
func (p *PostgresDB) GetMeasurementCounts() (total int64, last24h int64, lastTimestamp *time.Time, err error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get total count
	countQuery, countArgs, _ := p.builder.
		Select("COUNT(*)").
		From("measurements").
		ToSql()
	err = p.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count measurements: %w", err)
	}

	// Get last 24h count
	past24h := time.Now().UTC().Add(-24 * time.Hour)
	last24hQuery, last24hArgs, _ := p.builder.
		Select("COUNT(*)").
		From("measurements").
		Where(sq.GtOrEq{"created_at": past24h}).
		ToSql()
	err = p.db.QueryRowContext(ctx, last24hQuery, last24hArgs...).Scan(&last24h)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count last 24h measurements: %w", err)
	}

	// Get last measurement timestamp
	var ts time.Time
	lastQuery, lastArgs, _ := p.builder.
		Select("timestamp").
		From("measurements").
		OrderBy("timestamp DESC").
		Limit(1).
		ToSql()
	err = p.db.QueryRowContext(ctx, lastQuery, lastArgs...).Scan(&ts)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, nil, fmt.Errorf("failed to get last measurement timestamp: %w", err)
	}
	if err == nil {
		lastTimestamp = &ts
	}

	return total, last24h, lastTimestamp, nil
}

// GetLast24hStats retrieves average statistics for the last 24 hours
func (p *PostgresDB) GetLast24hStats() (*models.DashboardStats24h, error) {
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
		WHERE m.timestamp >= $1 AND n.archived = false
	`

	err := p.db.QueryRowContext(ctx, query, past24h).Scan(
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
func (p *PostgresDB) CleanupOldMeasurements(retentionDays int) (int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	cutoffDate := time.Now().UTC().AddDate(0, 0, -retentionDays)

	query, args, err := p.builder.
		Delete("measurements").
		Where(sq.Lt{"timestamp": cutoffDate}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := p.db.ExecContext(ctx, query, args...)
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
func (p *PostgresDB) CleanupOldFailedMeasurements(retentionDays int) (int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	cutoffDate := time.Now().UTC().AddDate(0, 0, -retentionDays)

	query, args, err := p.builder.
		Delete("failed_measurements").
		Where(sq.Lt{"timestamp": cutoffDate}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := p.db.ExecContext(ctx, query, args...)
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
