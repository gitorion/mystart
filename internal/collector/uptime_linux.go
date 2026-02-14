//go:build linux
// +build linux

package collector

import (
	"context"
	"fmt"
	"strconv"
)

// GatherUptimeInfo collects system uptime information on Linux.
func (s *SystemInfo) GatherUptimeInfo(ctx context.Context) error {
	uptimeStr, err := execCommand(ctx, "cut -d. -f1 /proc/uptime")
	if err != nil {
		return fmt.Errorf("getting uptime: %w", err)
	}

	uptime, err := strconv.Atoi(uptimeStr)
	if err != nil {
		return fmt.Errorf("%w: uptime", ErrParseFailure)
	}

	// Calculate time components
	s.UptimeDays = uptime / 60 / 60 / 24
	s.UptimeHours = (uptime / 60 / 60) % 24
	s.UptimeMinutes = (uptime / 60) % 60
	s.UptimeSeconds = uptime % 60

	return nil
}
