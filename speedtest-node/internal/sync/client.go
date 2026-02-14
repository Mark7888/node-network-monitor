package sync

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Client handles HTTP communication with the data server
type Client struct {
	serverURL string
	apiKey    string
	timeout   time.Duration
	tlsVerify bool
	client    *http.Client
	logger    *zap.Logger
}

// NewClient creates a new sync client
func NewClient(serverURL, apiKey string, timeout time.Duration, tlsVerify bool, logger *zap.Logger) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !tlsVerify,
		},
	}

	return &Client{
		serverURL: serverURL,
		apiKey:    apiKey,
		timeout:   timeout,
		tlsVerify: tlsVerify,
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		logger: logger,
	}
}

// Post sends a POST request to the given endpoint
func (c *Client) Post(endpoint string, payload interface{}) ([]byte, error) {
	// Marshal payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	url := c.serverURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
