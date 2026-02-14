package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	// Node configuration
	NodeName string

	// Server configuration
	ServerURL     string
	APIKey        string
	ServerTimeout time.Duration
	TLSVerify     bool

	// Speedtest configuration
	SpeedtestCron    string
	SpeedtestTimeout time.Duration
	RetryOnFailure   bool

	// Sync configuration
	BatchSize     int
	SyncInterval  time.Duration
	AliveInterval time.Duration

	// Database configuration
	DBPath string

	// Data retention
	RetentionDays int

	// Logging configuration
	LogLevel  string
	LogFormat string
	LogOutput string
}

// Load loads configuration from command-line arguments and environment variables
// Command-line arguments take precedence over environment variables
func Load() *Config {
	v := viper.New()

	// Set up environment variable prefix
	v.SetEnvPrefix("SPEEDTEST")
	v.AutomaticEnv()

	// Define flags using pflag
	pflag.String("node-name", getHostname(), "Human-readable node name")
	pflag.String("server-url", "", "Data server URL (HTTPS recommended)")
	pflag.String("api-key", "", "API key for authentication")
	pflag.Duration("server-timeout", 30*time.Second, "HTTP request timeout")
	pflag.Bool("tls-verify", true, "Verify TLS certificates")

	pflag.String("speedtest-cron", "*/10 * * * *", "Cron expression for measurements")
	pflag.Duration("speedtest-timeout", 120*time.Second, "Speedtest execution timeout")
	pflag.Bool("retry-on-failure", true, "Retry once if speedtest fails")

	pflag.Int("batch-size", 20, "Max measurements per sync request")
	pflag.Duration("sync-interval", 30*time.Second, "Check for unsent data interval")
	pflag.Duration("alive-interval", 60*time.Second, "Send alive signal interval")

	pflag.String("db-path", "./data/speedtest.db", "SQLite database path")

	pflag.Int("retention-days", 7, "Keep local data for N days")

	pflag.String("log-level", "info", "Log level: debug, info, warn, error")
	pflag.String("log-format", "json", "Log format: json or console")
	pflag.String("log-output", "./logs/speedtest-node.log", "Log file path")

	// Add version and help flags
	pflag.BoolP("version", "v", false, "Print version information")
	pflag.BoolP("help", "h", false, "Show help message")

	pflag.Parse()

	// Show help if requested
	if help, _ := pflag.CommandLine.GetBool("help"); help {
		fmt.Fprintf(os.Stdout, "Speedtest Node - Network Measurement Collector\n\n")
		fmt.Fprintf(os.Stdout, "Usage:\n")
		fmt.Fprintf(os.Stdout, "  speedtest-node [flags]\n\n")
		fmt.Fprintf(os.Stdout, "Flags:\n")
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stdout, "  All flags can be set via environment variables with SPEEDTEST_ prefix.\n")
		fmt.Fprintf(os.Stdout, "  Example: SPEEDTEST_NODE_NAME, SPEEDTEST_SERVER_URL, SPEEDTEST_API_KEY\n")
		fmt.Fprintf(os.Stdout, "  Note: Command-line flags take precedence over environment variables.\n")
		os.Exit(0)
	}

	// Show version if requested
	if version, _ := pflag.CommandLine.GetBool("version"); version {
		fmt.Fprintf(os.Stdout, "speedtest-node version 1.0.0\n")
		os.Exit(0)
	}

	// Bind pflag to viper (flags take precedence over env vars)
	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding flags: %v\n", err)
		os.Exit(1)
	}

	// Map environment variables to config keys
	// Viper will automatically look for SPEEDTEST_NODE_NAME, etc.
	v.BindEnv("node-name", "SPEEDTEST_NODE_NAME")
	v.BindEnv("server-url", "SPEEDTEST_SERVER_URL")
	v.BindEnv("api-key", "SPEEDTEST_SERVER_API_KEY")
	v.BindEnv("server-timeout", "SPEEDTEST_SERVER_TIMEOUT")
	v.BindEnv("tls-verify", "SPEEDTEST_TLS_VERIFY")
	v.BindEnv("speedtest-cron", "SPEEDTEST_CRON")
	v.BindEnv("speedtest-timeout", "SPEEDTEST_TIMEOUT")
	v.BindEnv("retry-on-failure", "SPEEDTEST_RETRY_ON_FAILURE")
	v.BindEnv("batch-size", "SPEEDTEST_BATCH_SIZE")
	v.BindEnv("sync-interval", "SPEEDTEST_SYNC_INTERVAL")
	v.BindEnv("alive-interval", "SPEEDTEST_ALIVE_INTERVAL")
	v.BindEnv("db-path", "SPEEDTEST_DB_PATH")
	v.BindEnv("retention-days", "SPEEDTEST_RETENTION_DAYS")
	v.BindEnv("log-level", "SPEEDTEST_LOG_LEVEL")
	v.BindEnv("log-format", "SPEEDTEST_LOG_FORMAT")
	v.BindEnv("log-output", "SPEEDTEST_LOG_OUTPUT")

	// Build config from viper
	cfg := &Config{
		NodeName:         v.GetString("node-name"),
		ServerURL:        v.GetString("server-url"),
		APIKey:           v.GetString("api-key"),
		ServerTimeout:    v.GetDuration("server-timeout"),
		TLSVerify:        v.GetBool("tls-verify"),
		SpeedtestCron:    v.GetString("speedtest-cron"),
		SpeedtestTimeout: v.GetDuration("speedtest-timeout"),
		RetryOnFailure:   v.GetBool("retry-on-failure"),
		BatchSize:        v.GetInt("batch-size"),
		SyncInterval:     v.GetDuration("sync-interval"),
		AliveInterval:    v.GetDuration("alive-interval"),
		DBPath:           v.GetString("db-path"),
		RetentionDays:    v.GetInt("retention-days"),
		LogLevel:         v.GetString("log-level"),
		LogFormat:        v.GetString("log-format"),
		LogOutput:        v.GetString("log-output"),
	}

	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// For now, we don't require strict validation
	// The application will fail gracefully if server URL or API key is missing
	return nil
}

// getHostname returns the system hostname or "unknown" if it fails
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
