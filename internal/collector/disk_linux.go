//go:build linux
// +build linux

package collector

import (
	"bufio"
	"context"
	"strconv"
	"strings"
)

// skipFSTypesLinux lists virtual/pseudo filesystem types to ignore.
var skipFSTypesLinux = map[string]bool{
	"tmpfs":       true,
	"devtmpfs":    true,
	"sysfs":       true,
	"proc":        true,
	"cgroup":      true,
	"cgroup2":     true,
	"pstore":      true,
	"bpf":         true,
	"tracefs":     true,
	"debugfs":     true,
	"securityfs":  true,
	"overlay":     true,
	"fuse.lxcfs":  true,
	"squashfs":    true,
	"hugetlbfs":   true,
	"mqueue":      true,
	"fusectl":     true,
	"fuse":        true,
}

func gatherDisk(ctx context.Context, info *SystemInfo) {
	// df -k with filesystem type (-T flag)
	out := execCommandSafe(ctx, "df -kT")
	if out == "" {
		// Fall back to df -k without type column
		out = execCommandSafe(ctx, "df -k")
	}

	var mounts []DiskMount
	seen := make(map[string]bool)

	sc := bufio.NewScanner(strings.NewReader(out))
	firstLine := true
	hasType := strings.Contains(strings.SplitN(out, "\n", 2)[0], "Type")
	for sc.Scan() {
		line := sc.Text()
		if firstLine {
			firstLine = false
			continue
		}
		fields := strings.Fields(line)

		var device, fsType, mountPoint string
		var totalKStr, usedKStr string

		if hasType && len(fields) >= 7 {
			device = fields[0]
			fsType = fields[1]
			totalKStr = fields[2]
			usedKStr = fields[3]
			mountPoint = fields[len(fields)-1]
		} else if len(fields) >= 6 {
			device = fields[0]
			fsType = ""
			totalKStr = fields[1]
			usedKStr = fields[2]
			mountPoint = fields[len(fields)-1]
		} else {
			continue
		}

		// Skip virtual filesystems
		if skipFSTypesLinux[fsType] {
			continue
		}

		// Only include real block devices
		if !strings.HasPrefix(device, "/dev/") {
			continue
		}

		// Skip boot partition and small system mounts (< 100 MB)
		if strings.HasPrefix(mountPoint, "/boot") ||
			strings.HasPrefix(mountPoint, "/sys") ||
			strings.HasPrefix(mountPoint, "/proc") ||
			strings.HasPrefix(mountPoint, "/run/") {
			continue
		}

		if seen[mountPoint] {
			continue
		}
		seen[mountPoint] = true

		total1K, err1 := strconv.ParseFloat(totalKStr, 64)
		used1K, err2 := strconv.ParseFloat(usedKStr, 64)
		if err1 != nil || err2 != nil || total1K < 102400 { // < 100 MB
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

	sortMounts(mounts)
	info.DiskMounts = mounts
}

func sortMounts(mounts []DiskMount) {
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
