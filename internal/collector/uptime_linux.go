//go:build linux
// +build linux

package collector

import (
	"context"
	"strconv"
	"strings"
)

func gatherUptime(ctx context.Context, info *SystemInfo) {
	// /proc/uptime: "123456.78 234567.89" – first field is seconds since boot
	raw := execCommandSafe(ctx, "cat /proc/uptime")
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return
	}
	// Drop fractional part
	secStr := strings.SplitN(fields[0], ".", 2)[0]
	total, err := strconv.Atoi(secStr)
	if err != nil {
		return
	}
	info.UptimeDays = total / 86400
	info.UptimeHours = (total % 86400) / 3600
	info.UptimeMinutes = (total % 3600) / 60
	info.UptimeSeconds = total % 60
}
