//go:build darwin
// +build darwin

package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// GatherUptimeInfo collects system uptime information on macOS.
func (s *SystemInfo) GatherUptimeInfo(ctx context.Context) error {
	// Get boot time
	bootTimeStr, err := execCommand(ctx, "sysctl -n kern.boottime")
	if err != nil {
		return fmt.Errorf("getting boot time: %w", err)
	}

	// Parse boot time: "{ sec = 1234567890, usec = 0 }" format
	// Extract seconds since epoch
	parts := strings.Split(bootTimeStr, ",")
	if len(parts) < 1 {
		return fmt.Errorf("%w: boot time format", ErrParseFailure)
	}

	secPart := strings.TrimSpace(parts[0])
	secPart = strings.TrimPrefix(secPart, "{ sec = ")
	secPart = strings.TrimSpace(secPart)

	bootTime, err := strconv.ParseInt(secPart, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: boot time", ErrParseFailure)
	}

	// Get current time
	currentTimeStr, err := execCommand(ctx, "date +%s")
	if err != nil {
		return fmt.Errorf("getting current time: %w", err)
	}

	currentTime, err := strconv.ParseInt(currentTimeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: current time", ErrParseFailure)
	}

	// Calculate uptime in seconds
	uptime := int(currentTime - bootTime)

	// Calculate time components
	s.UptimeDays = uptime / 60 / 60 / 24
	s.UptimeHours = (uptime / 60 / 60) % 24
	s.UptimeMinutes = (uptime / 60) % 60
	s.UptimeSeconds = uptime % 60

	return nil
}
