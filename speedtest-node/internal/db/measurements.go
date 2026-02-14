package db

import (
	"database/sql"
	"mark7888/speedtest-node/pkg/models"
	"time"

	"go.uber.org/zap"
)

// InsertMeasurement stores a measurement in the database
func (db *DB) InsertMeasurement(m *models.Measurement) error {
	query := `
		INSERT INTO measurements (
			timestamp, ping_jitter, ping_latency, ping_low, ping_high,
			download_bandwidth, download_bytes, download_elapsed,
			download_latency_iqm, download_latency_low, download_latency_high, download_latency_jitter,
			upload_bandwidth, upload_bytes, upload_elapsed,
			upload_latency_iqm, upload_latency_low, upload_latency_high, upload_latency_jitter,
			packet_loss, isp,
			interface_internal_ip, interface_name, interface_mac, interface_is_vpn, interface_external_ip,
			server_id, server_host, server_port, server_name, server_location, server_country, server_ip,
			result_id, result_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var pingJitter, pingLatency, pingLow, pingHigh sql.NullFloat64
	if m.Ping != nil {
		pingJitter = sql.NullFloat64{Float64: m.Ping.Jitter, Valid: true}
		pingLatency = sql.NullFloat64{Float64: m.Ping.Latency, Valid: true}
		pingLow = sql.NullFloat64{Float64: m.Ping.Low, Valid: true}
		pingHigh = sql.NullFloat64{Float64: m.Ping.High, Valid: true}
	}

	var downloadBandwidth, downloadBytes sql.NullInt64
	var downloadElapsed sql.NullInt64
	var downloadLatencyIQM, downloadLatencyLow, downloadLatencyHigh, downloadLatencyJitter sql.NullFloat64
	if m.Download != nil {
		downloadBandwidth = sql.NullInt64{Int64: m.Download.Bandwidth, Valid: true}
		downloadBytes = sql.NullInt64{Int64: m.Download.Bytes, Valid: true}
		downloadElapsed = sql.NullInt64{Int64: int64(m.Download.Elapsed), Valid: true}
		if m.Download.Latency != nil {
			downloadLatencyIQM = sql.NullFloat64{Float64: m.Download.Latency.IQM, Valid: true}
			downloadLatencyLow = sql.NullFloat64{Float64: m.Download.Latency.Low, Valid: true}
			downloadLatencyHigh = sql.NullFloat64{Float64: m.Download.Latency.High, Valid: true}
			downloadLatencyJitter = sql.NullFloat64{Float64: m.Download.Latency.Jitter, Valid: true}
		}
	}

	var uploadBandwidth, uploadBytes sql.NullInt64
	var uploadElapsed sql.NullInt64
	var uploadLatencyIQM, uploadLatencyLow, uploadLatencyHigh, uploadLatencyJitter sql.NullFloat64
	if m.Upload != nil {
		uploadBandwidth = sql.NullInt64{Int64: m.Upload.Bandwidth, Valid: true}
		uploadBytes = sql.NullInt64{Int64: m.Upload.Bytes, Valid: true}
		uploadElapsed = sql.NullInt64{Int64: int64(m.Upload.Elapsed), Valid: true}
		if m.Upload.Latency != nil {
			uploadLatencyIQM = sql.NullFloat64{Float64: m.Upload.Latency.IQM, Valid: true}
			uploadLatencyLow = sql.NullFloat64{Float64: m.Upload.Latency.Low, Valid: true}
			uploadLatencyHigh = sql.NullFloat64{Float64: m.Upload.Latency.High, Valid: true}
			uploadLatencyJitter = sql.NullFloat64{Float64: m.Upload.Latency.Jitter, Valid: true}
		}
	}

	var interfaceInternalIP, interfaceName, interfaceMAC, interfaceExternalIP sql.NullString
	var interfaceIsVPN sql.NullBool
	if m.Interface != nil {
		interfaceInternalIP = sql.NullString{String: m.Interface.InternalIP, Valid: true}
		interfaceName = sql.NullString{String: m.Interface.Name, Valid: true}
		interfaceMAC = sql.NullString{String: m.Interface.MacAddr, Valid: true}
		interfaceIsVPN = sql.NullBool{Bool: m.Interface.IsVPN, Valid: true}
		interfaceExternalIP = sql.NullString{String: m.Interface.ExternalIP, Valid: true}
	}

	var serverID, serverPort sql.NullInt64
	var serverHost, serverName, serverLocation, serverCountry, serverIP sql.NullString
	if m.Server != nil {
		serverID = sql.NullInt64{Int64: int64(m.Server.ID), Valid: true}
		serverPort = sql.NullInt64{Int64: int64(m.Server.Port), Valid: true}
		serverHost = sql.NullString{String: m.Server.Host, Valid: true}
		serverName = sql.NullString{String: m.Server.Name, Valid: true}
		serverLocation = sql.NullString{String: m.Server.Location, Valid: true}
		serverCountry = sql.NullString{String: m.Server.Country, Valid: true}
		serverIP = sql.NullString{String: m.Server.IP, Valid: true}
	}

	var resultID, resultURL sql.NullString
	if m.Result != nil {
		resultID = sql.NullString{String: m.Result.ID, Valid: true}
		resultURL = sql.NullString{String: m.Result.URL, Valid: true}
	}

	_, err := db.conn.Exec(query,
		m.Timestamp, pingJitter, pingLatency, pingLow, pingHigh,
		downloadBandwidth, downloadBytes, downloadElapsed,
		downloadLatencyIQM, downloadLatencyLow, downloadLatencyHigh, downloadLatencyJitter,
		uploadBandwidth, uploadBytes, uploadElapsed,
		uploadLatencyIQM, uploadLatencyLow, uploadLatencyHigh, uploadLatencyJitter,
		m.PacketLoss, m.ISP,
		interfaceInternalIP, interfaceName, interfaceMAC, interfaceIsVPN, interfaceExternalIP,
		serverID, serverHost, serverPort, serverName, serverLocation, serverCountry, serverIP,
		resultID, resultURL,
	)

	if err != nil {
		db.logger.Error("Failed to insert measurement", zap.Error(err))
		return err
	}

	db.logger.Debug("Measurement inserted successfully", zap.Time("timestamp", m.Timestamp))
	return nil
}

// GetUnsentMeasurements retrieves unsent measurements with a limit
func (db *DB) GetUnsentMeasurements(limit int) ([]*models.Measurement, error) {
	query := `
		SELECT 
			id, timestamp, created_at,
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
		WHERE sent = 0
		ORDER BY timestamp ASC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var measurements []*models.Measurement
	for rows.Next() {
		m := &models.Measurement{}
		var pingJitter, pingLatency, pingLow, pingHigh sql.NullFloat64
		var downloadBandwidth, downloadBytes, downloadElapsed sql.NullInt64
		var downloadLatencyIQM, downloadLatencyLow, downloadLatencyHigh, downloadLatencyJitter sql.NullFloat64
		var uploadBandwidth, uploadBytes, uploadElapsed sql.NullInt64
		var uploadLatencyIQM, uploadLatencyLow, uploadLatencyHigh, uploadLatencyJitter sql.NullFloat64
		var interfaceInternalIP, interfaceName, interfaceMAC, interfaceExternalIP sql.NullString
		var interfaceIsVPN sql.NullBool
		var serverID, serverPort sql.NullInt64
		var serverHost, serverName, serverLocation, serverCountry, serverIP sql.NullString
		var resultID, resultURL sql.NullString

		err := rows.Scan(
			&m.ID, &m.Timestamp, &m.CreatedAt,
			&pingJitter, &pingLatency, &pingLow, &pingHigh,
			&downloadBandwidth, &downloadBytes, &downloadElapsed,
			&downloadLatencyIQM, &downloadLatencyLow, &downloadLatencyHigh, &downloadLatencyJitter,
			&uploadBandwidth, &uploadBytes, &uploadElapsed,
			&uploadLatencyIQM, &uploadLatencyLow, &uploadLatencyHigh, &uploadLatencyJitter,
			&m.PacketLoss, &m.ISP,
			&interfaceInternalIP, &interfaceName, &interfaceMAC, &interfaceIsVPN, &interfaceExternalIP,
			&serverID, &serverHost, &serverPort, &serverName, &serverLocation, &serverCountry, &serverIP,
			&resultID, &resultURL,
		)
		if err != nil {
			return nil, err
		}

		// Reconstruct nested structures
		if pingJitter.Valid {
			m.Ping = &models.PingData{
				Jitter:  pingJitter.Float64,
				Latency: pingLatency.Float64,
				Low:     pingLow.Float64,
				High:    pingHigh.Float64,
			}
		}

		if downloadBandwidth.Valid {
			m.Download = &models.TransferData{
				Bandwidth: downloadBandwidth.Int64,
				Bytes:     downloadBytes.Int64,
				Elapsed:   int(downloadElapsed.Int64),
			}
			if downloadLatencyIQM.Valid {
				m.Download.Latency = &models.Latency{
					IQM:    downloadLatencyIQM.Float64,
					Low:    downloadLatencyLow.Float64,
					High:   downloadLatencyHigh.Float64,
					Jitter: downloadLatencyJitter.Float64,
				}
			}
		}

		if uploadBandwidth.Valid {
			m.Upload = &models.TransferData{
				Bandwidth: uploadBandwidth.Int64,
				Bytes:     uploadBytes.Int64,
				Elapsed:   int(uploadElapsed.Int64),
			}
			if uploadLatencyIQM.Valid {
				m.Upload.Latency = &models.Latency{
					IQM:    uploadLatencyIQM.Float64,
					Low:    uploadLatencyLow.Float64,
					High:   uploadLatencyHigh.Float64,
					Jitter: uploadLatencyJitter.Float64,
				}
			}
		}

		if interfaceInternalIP.Valid {
			m.Interface = &models.Interface{
				InternalIP: interfaceInternalIP.String,
				Name:       interfaceName.String,
				MacAddr:    interfaceMAC.String,
				IsVPN:      interfaceIsVPN.Bool,
				ExternalIP: interfaceExternalIP.String,
			}
		}

		if serverID.Valid {
			m.Server = &models.Server{
				ID:       int(serverID.Int64),
				Host:     serverHost.String,
				Port:     int(serverPort.Int64),
				Name:     serverName.String,
				Location: serverLocation.String,
				Country:  serverCountry.String,
				IP:       serverIP.String,
			}
		}

		if resultID.Valid {
			m.Result = &models.Result{
				ID:  resultID.String,
				URL: resultURL.String,
			}
		}

		measurements = append(measurements, m)
	}

	return measurements, nil
}

// MarkMeasurementsAsSent marks measurements as sent
func (db *DB) MarkMeasurementsAsSent(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE measurements SET sent = 1, sent_at = ? WHERE id IN (?` + repeatPlaceholder(len(ids)-1) + `)`
	args := make([]interface{}, len(ids)+1)
	args[0] = time.Now()
	for i, id := range ids {
		args[i+1] = id
	}

	_, err := db.conn.Exec(query, args...)
	return err
}

// DeleteMeasurementsBefore deletes measurements older than the given time
func (db *DB) DeleteMeasurementsBefore(before time.Time) error {
	_, err := db.conn.Exec("DELETE FROM measurements WHERE timestamp < ? AND sent = 1", before)
	return err
}

// InsertFailedMeasurement stores a failed measurement attempt
func (db *DB) InsertFailedMeasurement(timestamp time.Time, errorMsg string, retryCount int) error {
	_, err := db.conn.Exec(
		"INSERT INTO failed_measurements (timestamp, error_message, retry_count) VALUES (?, ?, ?)",
		timestamp, errorMsg, retryCount,
	)
	if err != nil {
		db.logger.Error("Failed to insert failed measurement", zap.Error(err))
		return err
	}
	db.logger.Debug("Failed measurement recorded", zap.Time("timestamp", timestamp))
	return nil
}

// GetUnsentFailedMeasurements retrieves unsent failed measurements with a limit
func (db *DB) GetUnsentFailedMeasurements(limit int) ([]*models.FailedMeasurement, error) {
	query := `
		SELECT id, timestamp, created_at, error_message, retry_count
		FROM failed_measurements
		WHERE sent = 0
		ORDER BY timestamp ASC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var failed []*models.FailedMeasurement
	for rows.Next() {
		f := &models.FailedMeasurement{}
		err := rows.Scan(&f.ID, &f.Timestamp, &f.CreatedAt, &f.ErrorMessage, &f.RetryCount)
		if err != nil {
			return nil, err
		}
		failed = append(failed, f)
	}

	return failed, nil
}

// MarkFailedMeasurementsAsSent marks failed measurements as sent
func (db *DB) MarkFailedMeasurementsAsSent(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE failed_measurements SET sent = 1, sent_at = ? WHERE id IN (?` + repeatPlaceholder(len(ids)-1) + `)`
	args := make([]interface{}, len(ids)+1)
	args[0] = time.Now()
	for i, id := range ids {
		args[i+1] = id
	}

	_, err := db.conn.Exec(query, args...)
	return err
}

// DeleteFailedMeasurementsBefore deletes failed measurements older than the given time
func (db *DB) DeleteFailedMeasurementsBefore(before time.Time) error {
	_, err := db.conn.Exec("DELETE FROM failed_measurements WHERE timestamp < ? AND sent = 1", before)
	return err
}

// repeatPlaceholder returns a string with n repeated ", ?" for SQL IN clauses
func repeatPlaceholder(n int) string {
	if n <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < n; i++ {
		result += ", ?"
	}
	return result
}
