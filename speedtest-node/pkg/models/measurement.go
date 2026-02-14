package models

import "time"

// Measurement represents a complete speedtest measurement
type Measurement struct {
	ID        int64     `json:"-"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"-"`

	// Ping data
	Ping *PingData `json:"ping,omitempty"`

	// Download data
	Download *TransferData `json:"download,omitempty"`

	// Upload data
	Upload *TransferData `json:"upload,omitempty"`

	// Network info
	PacketLoss float64    `json:"packet_loss"`
	ISP        string     `json:"isp"`
	Interface  *Interface `json:"interface,omitempty"`

	// Server info
	Server *Server `json:"server,omitempty"`

	// Result info
	Result *Result `json:"result,omitempty"`

	// Sync status (not sent to server)
	Sent   bool       `json:"-"`
	SentAt *time.Time `json:"-"`
}

// PingData represents ping measurement data
type PingData struct {
	Jitter  float64 `json:"jitter"`
	Latency float64 `json:"latency"`
	Low     float64 `json:"low"`
	High    float64 `json:"high"`
}

// TransferData represents download or upload measurement data
type TransferData struct {
	Bandwidth int64    `json:"bandwidth"` // bytes per second
	Bytes     int64    `json:"bytes"`
	Elapsed   int      `json:"elapsed"` // milliseconds
	Latency   *Latency `json:"latency,omitempty"`
}

// Latency represents latency measurements
type Latency struct {
	IQM    float64 `json:"iqm"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Jitter float64 `json:"jitter"`
}

// Interface represents network interface information
type Interface struct {
	InternalIP string `json:"internal_ip"`
	Name       string `json:"name"`
	MacAddr    string `json:"mac_addr"`
	IsVPN      bool   `json:"is_vpn"`
	ExternalIP string `json:"external_ip"`
}

// Server represents the speedtest server information
type Server struct {
	ID       int    `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Country  string `json:"country"`
	IP       string `json:"ip"`
}

// Result represents the speedtest result metadata
type Result struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// FailedMeasurement represents a failed speedtest attempt
type FailedMeasurement struct {
	ID           int64      `json:"-"`
	Timestamp    time.Time  `json:"timestamp"`
	CreatedAt    time.Time  `json:"-"`
	ErrorMessage string     `json:"error_message"`
	RetryCount   int        `json:"retry_count"`
	Sent         bool       `json:"-"`
	SentAt       *time.Time `json:"-"`
}

// AliveRequest represents the alive/keepalive signal request
type AliveRequest struct {
	NodeID    string    `json:"node_id"`
	NodeName  string    `json:"node_name"`
	Timestamp time.Time `json:"timestamp"`
}

// AliveResponse represents the server response to alive signal
type AliveResponse struct {
	Status     string    `json:"status"`
	ServerTime time.Time `json:"server_time"`
}

// MeasurementsRequest represents a batch of measurements to send to server
type MeasurementsRequest struct {
	NodeID       string         `json:"node_id"`
	NodeName     string         `json:"node_name"`
	Measurements []*Measurement `json:"measurements"`
}

// MeasurementsResponse represents the server response to measurements
type MeasurementsResponse struct {
	Status   string `json:"status"`
	Received int    `json:"received"`
	Failed   int    `json:"failed"`
}

// FailedMeasurementsRequest represents a batch of failed measurements to send to server
type FailedMeasurementsRequest struct {
	NodeID      string               `json:"node_id"`
	NodeName    string               `json:"node_name"`
	FailedTests []*FailedMeasurement `json:"failed_tests"`
}

// FailedMeasurementsResponse represents the server response to failed measurements
type FailedMeasurementsResponse struct {
	Status   string `json:"status"`
	Received int    `json:"received"`
}

// SpeedtestResult represents the raw output from speedtest CLI
type SpeedtestResult struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
		Low     float64 `json:"low"`
		High    float64 `json:"high"`
	} `json:"ping"`
	Download struct {
		Bandwidth int64 `json:"bandwidth"`
		Bytes     int64 `json:"bytes"`
		Elapsed   int   `json:"elapsed"`
		Latency   struct {
			IQM    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"download"`
	Upload struct {
		Bandwidth int64 `json:"bandwidth"`
		Bytes     int64 `json:"bytes"`
		Elapsed   int   `json:"elapsed"`
		Latency   struct {
			IQM    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"upload"`
	PacketLoss float64 `json:"packetLoss"`
	ISP        string  `json:"isp"`
	Interface  struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVPN      bool   `json:"isVpn"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		ID       int    `json:"id"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		IP       string `json:"ip"`
	} `json:"server"`
	Result struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"result"`
}
