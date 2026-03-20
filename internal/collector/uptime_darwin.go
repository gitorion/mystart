//go:build darwin
// +build darwin

package collector

import (
	"context"
	"strconv"
	"strings"
)

func gatherUptime(ctx context.Context, info *SystemInfo) {
	// kern.boottime: "{ sec = 1234567890, usec = 123456 } Mon Jan 01 00:00:00 2024"
	bootStr := execCommandSafe(ctx, "sysctl -n kern.boottime")
	nowStr := execCommandSafe(ctx, "date +%s")

	var bootSec, nowSec int64

	// Extract seconds from boottime
	if idx := strings.Index(bootStr, "sec = "); idx >= 0 {
		rest := bootStr[idx+6:]
		if end := strings.IndexAny(rest, ",} "); end > 0 {
			rest = rest[:end]
		}
		bootSec, _ = strconv.ParseInt(strings.TrimSpace(rest), 10, 64)
	}

	nowSec, _ = strconv.ParseInt(nowStr, 10, 64)

	if bootSec > 0 && nowSec > bootSec {
		total := int(nowSec - bootSec)
		info.UptimeDays = total / 86400
		info.UptimeHours = (total % 86400) / 3600
		info.UptimeMinutes = (total % 3600) / 60
		info.UptimeSeconds = total % 60
	}
}
