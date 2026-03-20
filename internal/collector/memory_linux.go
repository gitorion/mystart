//go:build linux
// +build linux

package collector

import (
	"bufio"
	"context"
	"strconv"
	"strings"
)

func gatherMemory(ctx context.Context, info *SystemInfo) {
	raw := execCommandSafe(ctx, "cat /proc/meminfo")
	vals := parseMeminfo(raw)

	kbToGB := func(kb float64) float64 { return kb / 1024 / 1024 }

	info.MemTotalGB = kbToGB(vals["MemTotal"])
	avail := vals["MemAvailable"]
	info.MemUsedGB = kbToGB(vals["MemTotal"] - avail)

	info.SwapTotalGB = kbToGB(vals["SwapTotal"])
	info.SwapUsedGB = kbToGB(vals["SwapTotal"] - vals["SwapFree"])
}

func parseMeminfo(raw string) map[string]float64 {
	m := make(map[string]float64)
	sc := bufio.NewScanner(strings.NewReader(raw))
	for sc.Scan() {
		line := sc.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		valStr := strings.Fields(strings.TrimSpace(parts[1]))
		if len(valStr) == 0 {
			continue
		}
		v, err := strconv.ParseFloat(valStr[0], 64)
		if err == nil {
			m[key] = v
		}
	}
	return m
}
