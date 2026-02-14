package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InsertMeasurement inserts or updates a measurement
func (db *DB) InsertMeasurement(m *models.Measurement) error {
	ctx, cancel := withTimeout()
	defer cancel()

	nowSQL := db.getNowSQL()
	query := fmt.Sprintf(`
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
			$1, $2, %s,
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
	`, nowSQL)

	_, err := db.ExecContext(ctx, query,
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

// GetMeasurementsByNode retrieves measurements for a specific node
func (db *DB) GetMeasurementsByNode(nodeID uuid.UUID, from, to *time.Time, page, limit int) ([]models.Measurement, int, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Build query based on time filters
	var whereClauses []string
	var args []interface{}
	argPos := 1

	whereClauses = append(whereClauses, fmt.Sprintf("node_id = $%d", argPos))
	args = append(args, nodeID)
	argPos++

	if from != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("timestamp >= $%d", argPos))
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("timestamp <= $%d", argPos))
		args = append(args, *to)
		argPos++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := "SELECT COUNT(*) FROM measurements WHERE " + whereClause
	var total int
	err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count measurements: %w", err)
	}

	// Get measurements
	selectQuery := `
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
			result_id, result_url
		FROM measurements
		WHERE ` + whereClause + `
		ORDER BY timestamp DESC
		LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := db.QueryContext(ctx, selectQuery, args...)
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
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan measurement: %w", err)
		}
		measurements = append(measurements, m)
	}

	return measurements, total, nil
}

// InsertFailedMeasurement inserts a failed measurement record
func (db *DB) InsertFailedMeasurement(nodeID uuid.UUID, timestamp time.Time, errorMessage string, retryCount int) error {
	ctx, cancel := withTimeout()
	defer cancel()

	nowSQL := db.getNowSQL()
	query := fmt.Sprintf(`
		INSERT INTO failed_measurements (node_id, timestamp, error_message, retry_count, created_at)
		VALUES ($1, $2, $3, $4, %s)
	`, nowSQL)

	_, err := db.ExecContext(ctx, query, nodeID, timestamp, errorMessage, retryCount)
	if err != nil {
		return fmt.Errorf("failed to insert failed measurement: %w", err)
	}

	return nil
}

// GetAggregatedMeasurements retrieves aggregated measurements for charting
func (db *DB) GetAggregatedMeasurements(nodeIDs []uuid.UUID, from, to time.Time, interval string) ([]models.AggregatedMeasurement, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Determine truncation based on interval
	var truncFunc string
	switch interval {
	case "1h":
		truncFunc = "date_trunc('hour', timestamp)"
	case "6h":
		truncFunc = "date_trunc('hour', timestamp) - (EXTRACT(HOUR FROM timestamp)::int % 6) * INTERVAL '1 hour'"
	case "1d":
		truncFunc = "date_trunc('day', timestamp)"
	default:
		truncFunc = "date_trunc('hour', timestamp)"
	}

	// Build query
	var whereClauses []string
	var args []interface{}
	argPos := 1

	whereClauses = append(whereClauses, fmt.Sprintf("timestamp >= $%d", argPos))
	args = append(args, from)
	argPos++

	whereClauses = append(whereClauses, fmt.Sprintf("timestamp <= $%d", argPos))
	args = append(args, to)
	argPos++

	if len(nodeIDs) > 0 {
		placeholders := make([]string, len(nodeIDs))
		for i, nodeID := range nodeIDs {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, nodeID)
			argPos++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("m.node_id IN (%s)", strings.Join(placeholders, ",")))
	}

	whereClause := strings.Join(whereClauses, " AND ")

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
			COUNT(*) as sample_count
		FROM measurements m
		JOIN nodes n ON m.node_id = n.id
		WHERE %s
		GROUP BY time_bucket, m.node_id, n.name
		ORDER BY time_bucket, m.node_id
	`, truncFunc, whereClause)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query aggregated measurements: %w", err)
	}
	defer rows.Close()

	var results []models.AggregatedMeasurement
	for rows.Next() {
		var agg models.AggregatedMeasurement
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
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan aggregated measurement: %w", err)
		}
		results = append(results, agg)
	}

	return results, nil
}

// GetMeasurementCounts retrieves measurement counts
func (db *DB) GetMeasurementCounts() (total int64, last24h int64, lastTimestamp *time.Time, err error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Get total count
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM measurements").Scan(&total)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count measurements: %w", err)
	}

	// Get last 24h count
	past24h := time.Now().UTC().Add(-24 * time.Hour)
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM measurements WHERE created_at >= $1",
		past24h,
	).Scan(&last24h)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count last 24h measurements: %w", err)
	}

	// Get last measurement timestamp
	var ts time.Time
	err = db.QueryRowContext(ctx,
		"SELECT timestamp FROM measurements ORDER BY timestamp DESC LIMIT 1",
	).Scan(&ts)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, nil, fmt.Errorf("failed to get last measurement timestamp: %w", err)
	}
	if err == nil {
		lastTimestamp = &ts
	}

	return total, last24h, lastTimestamp, nil
}

// GetLast24hStats retrieves average statistics for the last 24 hours
func (db *DB) GetLast24hStats() (*models.DashboardStats24h, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	past24h := time.Now().UTC().Add(-24 * time.Hour)

	stats := &models.DashboardStats24h{}
	err := db.QueryRowContext(ctx, `
		SELECT
			COALESCE(AVG(download_bandwidth) / 125000.0, 0) as avg_download_mbps,
			COALESCE(AVG(upload_bandwidth) / 125000.0, 0) as avg_upload_mbps,
			COALESCE(AVG(ping_latency), 0) as avg_ping_ms,
			COALESCE(AVG(ping_jitter), 0) as avg_jitter_ms,
			COALESCE(AVG(packet_loss), 0) as avg_packet_loss
		FROM measurements
		WHERE timestamp >= $1
	`, past24h).Scan(
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
func (db *DB) CleanupOldMeasurements(retentionDays int) (int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	result, err := db.ExecContext(ctx,
		"DELETE FROM measurements WHERE timestamp < $1",
		cutoffDate,
	)
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
func (db *DB) CleanupOldFailedMeasurements(retentionDays int) (int64, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	cutoffDate := time.Now().UTC().AddDate(0, 0, -retentionDays)

	result, err := db.ExecContext(ctx,
		"DELETE FROM failed_measurements WHERE timestamp < $1",
		cutoffDate,
	)
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
