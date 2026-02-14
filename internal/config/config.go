package config

import "time"

// Unicode emojis and symbols.
const (
	ThumbUp    = "\U0001F44D"
	StopEmoji  = "\U0001F6D1"
	PointRight = "\U0001F449"
	CrossMark  = "\u274C"
)

// HostTasks maps hostnames to their descriptions.
var HostTasks = map[string]string{
	"titan":    "VM - Dockerised Plex Media and Transmission Server",
	"hyperion": "LXC - Reverse Proxy by Caddy",
	"dione":    "LXC - Uptime monitoring by Uptime Kuma",
	"mimas":    "LXC - Metrics and Logging by Grafana/Prometheus",
	"backup":   "LXC - File backup",
	"dad":      "VM - Dockerised Plex Media Server and SMB share",
	"saturn":   "Local Proxmox Server",
	"neptune":  "Remote Proxmox Server",
}

// Config holds application configuration.
type Config struct {
	CommandTimeout time.Duration
	OverallTimeout time.Duration
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		CommandTimeout: 5 * time.Second,
		OverallTimeout: 30 * time.Second,
	}
}
