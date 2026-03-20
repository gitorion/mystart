//go:build darwin
// +build darwin

package collector

import (
	"context"
	"runtime"
	"strconv"
	"strings"
)

func gatherCPU(ctx context.Context, info *SystemInfo) {
	// CPU model
	info.CPUModel = execCommandSafe(ctx, "sysctl -n machdep.cpu.brand_string")

	// Physical cores
	coresStr := execCommandSafe(ctx, "sysctl -n hw.physicalcpu")
	if n, err := strconv.Atoi(coresStr); err == nil {
		info.CPUCores = n
	}

	// Logical cores (threads)
	info.CPUThreads = runtime.NumCPU()

	// CPU frequency (Hz → GHz)
	for _, key := range []string{"hw.cpufrequency", "hw.cpufrequency_max"} {
		freqStr := execCommandSafe(ctx, "sysctl -n "+key)
		if freq, err := strconv.ParseFloat(freqStr, 64); err == nil && freq > 0 {
			info.CPUHz = freq / 1e9
			break
		}
	}
	// Fallback: parse GHz from brand string (e.g. "Intel … @ 2.60GHz")
	if info.CPUHz == 0 && strings.Contains(info.CPUModel, "@") {
		parts := strings.Split(info.CPUModel, "@")
		ghz := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(parts[1]), "GHz"))
		if v, err := strconv.ParseFloat(ghz, 64); err == nil {
			info.CPUHz = v
		}
	}

	// CPU usage via top (one sample)
	out := execCommandSafe(ctx, "top -l 1 -n 0 | grep -E '^CPU usage'")
	if out != "" {
		// "CPU usage: 12.34% user, 5.67% sys, 82.00% idle"
		var user, sys float64
		fields := strings.Fields(out)
		for i, f := range fields {
			switch f {
			case "user,":
				if i > 0 {
					pct := strings.TrimSuffix(fields[i-1], "%")
					user, _ = strconv.ParseFloat(pct, 64)
				}
			case "sys,":
				if i > 0 {
					pct := strings.TrimSuffix(fields[i-1], "%")
					sys, _ = strconv.ParseFloat(pct, 64)
				}
			}
		}
		info.CPUUsage = user + sys
	}

	// Load average
	lavg := execCommandSafe(ctx, "sysctl -n vm.loadavg")
	// Format: "{ 1.23 0.89 0.72 }"
	lavg = strings.Trim(lavg, "{ }")
	parts := strings.Fields(lavg)
	if len(parts) >= 3 {
		info.LoadAvg1, _ = strconv.ParseFloat(parts[0], 64)
		info.LoadAvg5, _ = strconv.ParseFloat(parts[1], 64)
		info.LoadAvg15, _ = strconv.ParseFloat(parts[2], 64)
	}
}
