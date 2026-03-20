//go:build darwin
// +build darwin

package collector

import (
	"bufio"
	"context"
	"strconv"
	"strings"
)

// skipDevicesDarwin lists filesystem types to ignore on macOS.
var skipFSTypesDarwin = map[string]bool{
	"devfs":    true,
	"autofs":   true,
	"tmpfs":    true,
	"nullfs":   true,
	"fdesc":    true,
	"kernfs":   true,
	"procfs":   true,
}

func gatherDisk(ctx context.Context, info *SystemInfo) {
	// df -k: 1 K-block units, portable across macOS versions
	out := execCommandSafe(ctx, "df -k")
	if out == "" {
		return
	}

	seen := make(map[string]bool) // deduplicate by mount point
	var mounts []DiskMount

	sc := bufio.NewScanner(strings.NewReader(out))
	firstLine := true
	for sc.Scan() {
		line := sc.Text()
		if firstLine {
			firstLine = false
			continue // skip header
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		device := fields[0]
		mountPoint := fields[len(fields)-1]

		// Skip virtual / system filesystems
		if skipFSTypesDarwin[device] {
			continue
		}
		if !strings.HasPrefix(device, "/dev/") {
			continue
		}

		// Show /, the APFS data volume, and external volumes under /Volumes/
		allowed := mountPoint == "/" ||
			mountPoint == "/System/Volumes/Data" ||
			strings.HasPrefix(mountPoint, "/Volumes/")
		if !allowed {
			continue
		}

		if seen[mountPoint] {
			continue
		}
		seen[mountPoint] = true

		total1K, err1 := strconv.ParseFloat(fields[1], 64)
		used1K, err2 := strconv.ParseFloat(fields[2], 64)
		if err1 != nil || err2 != nil || total1K == 0 {
			continue
		}

		totalGB := total1K / 1024 / 1024
		usedGB := used1K / 1024 / 1024
		pct := usedGB / totalGB * 100

		mounts = append(mounts, DiskMount{
			Path:        mountPoint,
			UsedGB:      usedGB,
			TotalGB:     totalGB,
			UsedPercent: pct,
		})
	}

	// Sort: root first, then alphabetical
	sortMounts(mounts)
	info.DiskMounts = mounts
}

func sortMounts(mounts []DiskMount) {
	// Simple insertion sort: "/" first, then others alphabetically
	for i := 1; i < len(mounts); i++ {
		for j := i; j > 0; j-- {
			a, b := mounts[j-1], mounts[j]
			if a.Path == "/" {
				break
			}
			if b.Path < a.Path {
				mounts[j-1], mounts[j] = mounts[j], mounts[j-1]
			} else {
				break
			}
		}
	}
}
