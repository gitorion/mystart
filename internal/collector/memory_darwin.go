//go:build darwin
// +build darwin

package collector

import (
	"context"
	"strconv"
	"strings"
)

func gatherMemory(ctx context.Context, info *SystemInfo) {
	// Total RAM (bytes → GB)
	memStr := execCommandSafe(ctx, "sysctl -n hw.memsize")
	if memBytes, err := strconv.ParseFloat(memStr, 64); err == nil {
		info.MemTotalGB = memBytes / 1024 / 1024 / 1024
	}

	// Used RAM from vm_stat
	vmStat := execCommandSafe(ctx, "vm_stat")
	pageSize := 4096.0
	var active, wired, compressed float64
	for _, line := range strings.Split(vmStat, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSuffix(fields[len(fields)-1], "."), 64)
		if err != nil {
			continue
		}
		switch {
		case strings.Contains(line, "Pages active"):
			active = val
		case strings.Contains(line, "Pages wired"):
			wired = val
		case strings.Contains(line, "Pages occupied by compressor"):
			compressed = val
		}
	}
	info.MemUsedGB = (active + wired + compressed) * pageSize / 1024 / 1024 / 1024

	// Swap from sysctl vm.swapusage
	// Output: "total = 2048.00M  used = 512.00M  free = 1536.00M  (encrypted)"
	swap := execCommandSafe(ctx, "sysctl -n vm.swapusage")
	parseSwapField := func(keyword string) float64 {
		idx := strings.Index(swap, keyword)
		if idx < 0 {
			return 0
		}
		rest := swap[idx+len(keyword):]
		rest = strings.TrimSpace(rest)
		if i := strings.IndexAny(rest, " \t("); i > 0 {
			rest = rest[:i]
		}
		var mul float64 = 1
		if strings.HasSuffix(rest, "G") {
			mul = 1
			rest = strings.TrimSuffix(rest, "G")
		} else if strings.HasSuffix(rest, "M") {
			mul = 1.0 / 1024
			rest = strings.TrimSuffix(rest, "M")
		}
		v, _ := strconv.ParseFloat(rest, 64)
		return v * mul
	}
	info.SwapTotalGB = parseSwapField("total = ")
	info.SwapUsedGB = parseSwapField("used = ")
}
