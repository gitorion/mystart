package collector

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Error variables with Err prefix following Uber style guide.
var (
	ErrCommandFailed   = errors.New("command execution failed")
	ErrInvalidOutput   = errors.New("invalid command output")
	ErrParseFailure    = errors.New("failed to parse value")
	ErrServiceNotFound = errors.New("service not found")
)

// SystemInfo holds all collected system information.
type SystemInfo struct {
	// System identification
	Uname    string
	Distro   string
	Host     string
	User     string
	HostTask string

	// CPU information
	CPUCores   int
	CPUThreads int
	CPUHz      float64
	CPUUsage   string
	LoadAvg    string

	// Memory information
	MemTotal  string
	MemUsed   string
	SwapTotal string
	SwapUsed  string

	// Disk information
	DiskUse      string
	DiskSize     string
	DiskPoolSize string
	DiskPoolUsed string

	// Network information
	IPv4 string
	IPv6 string

	// User session information
	LastLog        string
	ThisLog        string
	ProcessesUser  string
	ProcessesAll   string
	ActiveSessions string
	Users          int

	// Uptime information
	UptimeDays    int
	UptimeHours   int
	UptimeMinutes int
	UptimeSeconds int

	// Fan information (for specific hosts)
	Fan1 string
	Fan2 string
	Fan3 string
	Fan4 string
	Fan5 string
	Fan6 string

	// VPN/Transmission information (for specific hosts)
	NordAddr        string
	TransAddr       string
	VPNCheck        string
	TranskickStatus string
}

// NewSystemInfo creates and initializes a new SystemInfo instance.
func NewSystemInfo() *SystemInfo {
	return &SystemInfo{}
}

// execCommand executes a shell command with context and returns output.
func execCommand(ctx context.Context, command string) (string, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %v (stderr: %s)", ErrCommandFailed, err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// execCommandSafe executes a command and returns empty string on error.
// Used for optional system information that may not be available.
func execCommandSafe(ctx context.Context, command string) string {
	output, err := execCommand(ctx, command)
	if err != nil {
		return ""
	}
	return output
}
