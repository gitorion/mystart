//go:build linux
// +build linux

package collector

import (
	"bufio"
	"context"
	"strconv"
	"strings"
)

func gatherCPU(ctx context.Context, info *SystemInfo) {
	// CPU model
	info.CPUModel = execCommandSafe(ctx,
		`grep -m1 "model name" /proc/cpuinfo | cut -d: -f2 | sed 's/^ //'`)

	// Physical cores
	coresStr := execCommandSafe(ctx,
		`grep -m1 "cpu cores" /proc/cpuinfo | awk '{print $4}'`)
	if n, err := strconv.Atoi(coresStr); err == nil && n > 0 {
		info.CPUCores = n
	}

	// Logical cores + average frequency
	mhzOut := execCommandSafe(ctx, `grep "cpu MHz" /proc/cpuinfo | awk '{print $4}'`)
	if mhzOut != "" {
		var total float64
		var count int
		sc := bufio.NewScanner(strings.NewReader(mhzOut))
		for sc.Scan() {
			if v, err := strconv.ParseFloat(sc.Text(), 64); err == nil {
				total += v
				count++
			}
		}
		if count > 0 {
			info.CPUThreads = count
			info.CPUHz = (total / float64(count)) / 1000.0
		}
	}
	if info.CPUCores == 0 && info.CPUThreads > 0 {
		info.CPUCores = info.CPUThreads
	}

	// CPU usage: two /proc/stat snapshots with a small sleep between them
	cpuUsage := func() float64 {
		parse := func(line string) (total, idle float64) {
			fields := strings.Fields(line)
			for i := 1; i < len(fields) && i <= 8; i++ {
				v, _ := strconv.ParseFloat(fields[i], 64)
				total += v
				if i == 4 {
					idle = v
				}
			}
			return
		}
		stat1 := execCommandSafe(ctx, `grep "^cpu " /proc/stat`)
		execCommandSafe(ctx, "sleep 0.2")
		stat2 := execCommandSafe(ctx, `grep "^cpu " /proc/stat`)
		t1, i1 := parse(stat1)
		t2, i2 := parse(stat2)
		dt := t2 - t1
		if dt == 0 {
			return 0
		}
		return 100.0 * (dt - (i2 - i1)) / dt
	}
	info.CPUUsage = cpuUsage()

	// Load average from /proc/loadavg
	lavg := execCommandSafe(ctx, "cat /proc/loadavg")
	parts := strings.Fields(lavg)
	if len(parts) >= 3 {
		info.LoadAvg1, _ = strconv.ParseFloat(parts[0], 64)
		info.LoadAvg5, _ = strconv.ParseFloat(parts[1], 64)
		info.LoadAvg15, _ = strconv.ParseFloat(parts[2], 64)
	}
}
