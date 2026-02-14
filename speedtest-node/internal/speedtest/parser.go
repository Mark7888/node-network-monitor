package speedtest

import (
	"encoding/json"
	"mark7888/speedtest-node/pkg/models"
	"time"
)

// ParseResult parses the speedtest JSON output into a Measurement model
func ParseResult(data []byte) (*models.Measurement, error) {
	var result models.SpeedtestResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, result.Timestamp)
	if err != nil {
		return nil, err
	}

	measurement := &models.Measurement{
		Timestamp:  timestamp,
		PacketLoss: result.PacketLoss,
		ISP:        result.ISP,
	}

	// Ping data
	measurement.Ping = &models.PingData{
		Jitter:  result.Ping.Jitter,
		Latency: result.Ping.Latency,
		Low:     result.Ping.Low,
		High:    result.Ping.High,
	}

	// Download data
	measurement.Download = &models.TransferData{
		Bandwidth: result.Download.Bandwidth,
		Bytes:     result.Download.Bytes,
		Elapsed:   result.Download.Elapsed,
		Latency: &models.Latency{
			IQM:    result.Download.Latency.IQM,
			Low:    result.Download.Latency.Low,
			High:   result.Download.Latency.High,
			Jitter: result.Download.Latency.Jitter,
		},
	}

	// Upload data
	measurement.Upload = &models.TransferData{
		Bandwidth: result.Upload.Bandwidth,
		Bytes:     result.Upload.Bytes,
		Elapsed:   result.Upload.Elapsed,
		Latency: &models.Latency{
			IQM:    result.Upload.Latency.IQM,
			Low:    result.Upload.Latency.Low,
			High:   result.Upload.Latency.High,
			Jitter: result.Upload.Latency.Jitter,
		},
	}

	// Interface
	measurement.Interface = &models.Interface{
		InternalIP: result.Interface.InternalIP,
		Name:       result.Interface.Name,
		MacAddr:    result.Interface.MacAddr,
		IsVPN:      result.Interface.IsVPN,
		ExternalIP: result.Interface.ExternalIP,
	}

	// Server
	measurement.Server = &models.Server{
		ID:       result.Server.ID,
		Host:     result.Server.Host,
		Port:     result.Server.Port,
		Name:     result.Server.Name,
		Location: result.Server.Location,
		Country:  result.Server.Country,
		IP:       result.Server.IP,
	}

	// Result
	measurement.Result = &models.Result{
		ID:  result.Result.ID,
		URL: result.Result.URL,
	}

	return measurement, nil
}
