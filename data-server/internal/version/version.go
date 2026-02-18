package version

// Version is the application version, can be set at build time using ldflags
// Example: go build -ldflags "-X mark7888/speedtest-data-server/internal/version.Version=v1.0.0"
var Version = "development"

// Get returns the current application version
func Get() string {
	return Version
}
