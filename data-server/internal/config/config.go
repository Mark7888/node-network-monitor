package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Admin     AdminConfig
	JWT       JWTConfig
	Node      NodeConfig
	Retention RetentionConfig
	API       APIConfig
	Logging   LoggingConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host       string
	Port       int
	Mode       string // debug, release
	TLSEnabled bool
	TLSCert    string
	TLSKey     string
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Type               string // postgres or sqlite
	Path               string // SQLite database file path
	Host               string
	Port               int
	User               string
	Password           string
	Name               string
	SSLMode            string
	MaxConnections     int
	MaxIdle            int
	ConnectionLifetime time.Duration
}

// AdminConfig holds admin user configuration
type AdminConfig struct {
	Username string
	Password string
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// NodeConfig holds node monitoring configuration
type NodeConfig struct {
	AliveTimeout        time.Duration
	InactiveTimeout     time.Duration
	StatusCheckInterval time.Duration
}

// RetentionConfig holds data retention configuration
type RetentionConfig struct {
	MeasurementsDays int
	FailedDays       int
	CleanupInterval  time.Duration
}

// APIConfig holds API-related configuration
type APIConfig struct {
	RateLimit int
	Timeout   time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level         string
	Format        string
	Output        string
	OutputConsole bool
}

// Load loads configuration from command-line arguments and environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Define flags
	// Server
	flag.StringVar(&cfg.Server.Host, "host", getEnv("SERVER_HOST", "0.0.0.0"), "Server host address")
	flag.IntVar(&cfg.Server.Port, "port", getEnvInt("SERVER_PORT", 8080), "Server port")
	flag.StringVar(&cfg.Server.Mode, "mode", getEnv("SERVER_MODE", "release"), "Server mode: debug, release")
	flag.BoolVar(&cfg.Server.TLSEnabled, "tls-enabled", getEnvBool("TLS_ENABLED", false), "Enable HTTPS/TLS")
	flag.StringVar(&cfg.Server.TLSCert, "tls-cert", getEnv("TLS_CERT", ""), "Path to TLS certificate file")
	flag.StringVar(&cfg.Server.TLSKey, "tls-key", getEnv("TLS_KEY", ""), "Path to TLS key file")

	// Database
	flag.StringVar(&cfg.Database.Type, "db-type", getEnv("DB_TYPE", "postgres"), "Database type: postgres, sqlite")
	flag.StringVar(&cfg.Database.Path, "db-path", getEnv("DB_PATH", "./data/speedtest.db"), "SQLite database file path")
	flag.StringVar(&cfg.Database.Host, "db-host", getEnv("DB_HOST", "localhost"), "PostgreSQL host")
	flag.IntVar(&cfg.Database.Port, "db-port", getEnvInt("DB_PORT", 5432), "PostgreSQL port")
	flag.StringVar(&cfg.Database.User, "db-user", getEnv("DB_USER", "speedtest"), "PostgreSQL user")
	flag.StringVar(&cfg.Database.Password, "db-password", getEnv("DB_PASSWORD", ""), "PostgreSQL password")
	flag.StringVar(&cfg.Database.Name, "db-name", getEnv("DB_NAME", "speedtest2"), "PostgreSQL database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", getEnv("DB_SSLMODE", "require"), "PostgreSQL SSL mode")
	flag.IntVar(&cfg.Database.MaxConnections, "db-max-connections", getEnvInt("DB_MAX_CONNECTIONS", 25), "Max database connections")
	flag.IntVar(&cfg.Database.MaxIdle, "db-max-idle", getEnvInt("DB_MAX_IDLE", 5), "Max idle connections")
	flag.DurationVar(&cfg.Database.ConnectionLifetime, "db-connection-lifetime", getEnvDuration("DB_CONNECTION_LIFETIME", 5*time.Minute), "Connection lifetime")

	// Admin
	flag.StringVar(&cfg.Admin.Username, "admin-username", getEnv("ADMIN_USERNAME", "admin"), "Admin username")
	flag.StringVar(&cfg.Admin.Password, "admin-password", getEnv("ADMIN_PASSWORD", ""), "Admin password")

	// JWT
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", getEnv("JWT_SECRET", ""), "JWT secret key")
	flag.DurationVar(&cfg.JWT.Expiry, "jwt-expiry", getEnvDuration("JWT_EXPIRY", 24*time.Hour), "JWT token expiration")

	// Node
	flag.DurationVar(&cfg.Node.AliveTimeout, "alive-timeout", getEnvDuration("ALIVE_TIMEOUT", 2*time.Minute), "Node alive signal timeout")
	flag.DurationVar(&cfg.Node.InactiveTimeout, "inactive-timeout", getEnvDuration("INACTIVE_TIMEOUT", 1*time.Hour), "Node inactive timeout")
	flag.DurationVar(&cfg.Node.StatusCheckInterval, "status-check-interval", getEnvDuration("STATUS_CHECK_INTERVAL", 30*time.Second), "Node status check interval")

	// Retention
	flag.IntVar(&cfg.Retention.MeasurementsDays, "retention-measurements", getEnvInt("RETENTION_MEASUREMENTS", 365), "Keep measurements for N days")
	flag.IntVar(&cfg.Retention.FailedDays, "retention-failed", getEnvInt("RETENTION_FAILED", 90), "Keep failed measurements for N days")
	flag.DurationVar(&cfg.Retention.CleanupInterval, "cleanup-interval", getEnvDuration("CLEANUP_INTERVAL", 24*time.Hour), "Data cleanup interval")

	// API
	flag.IntVar(&cfg.API.RateLimit, "rate-limit", getEnvInt("RATE_LIMIT", 100), "Requests per minute per API key")
	flag.DurationVar(&cfg.API.Timeout, "timeout", getEnvDuration("API_TIMEOUT", 30*time.Second), "API request timeout")

	// Logging
	flag.StringVar(&cfg.Logging.Level, "log-level", getEnv("LOG_LEVEL", "info"), "Log level: debug, info, warn, error")
	flag.StringVar(&cfg.Logging.Format, "log-format", getEnv("LOG_FORMAT", "json"), "Log format: json or console")
	flag.StringVar(&cfg.Logging.Output, "log-output", getEnv("LOG_OUTPUT", "./logs/data-server.log"), "Log file path")
	flag.BoolVar(&cfg.Logging.OutputConsole, "log-output-console", getEnvBool("LOG_OUTPUT_CONSOLE", false), "Also output logs to console")

	flag.Parse()

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	// Validate database type
	if c.Database.Type != "postgres" && c.Database.Type != "sqlite" {
		return fmt.Errorf("database type must be 'postgres' or 'sqlite'")
	}

	// PostgreSQL requires password
	if c.Database.Type == "postgres" && c.Database.Password == "" {
		return fmt.Errorf("database password is required for PostgreSQL (--db-password or DB_PASSWORD)")
	}

	if c.Admin.Password == "" {
		return fmt.Errorf("admin password is required (--admin-password or ADMIN_PASSWORD)")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required (--jwt-secret or JWT_SECRET)")
	}
	if c.Server.TLSEnabled && (c.Server.TLSCert == "" || c.Server.TLSKey == "") {
		return fmt.Errorf("TLS certificate and key are required when TLS is enabled")
	}
	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	if c.Database.Type == "sqlite" {
		return c.Database.Path
	}
	// PostgreSQL
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// Helper functions to get environment variables with defaults

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
