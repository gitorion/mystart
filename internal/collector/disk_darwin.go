//go:build darwin
// +build darwin

package collector

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/orion/mystart/internal/config"
)

// GatherDiskInfo collects disk usage information on macOS.
func (s *SystemInfo) GatherDiskInfo(ctx context.Context) error {
	if err := s.gatherRootDiskInfo(ctx); err != nil {
		return err
	}

	// Disk pool calculation is optional and may fail
	s.calculateDiskPool(ctx)

	return nil
}

// gatherRootDiskInfo retrieves root filesystem usage.
func (s *SystemInfo) gatherRootDiskInfo(ctx context.Context) error {
	diskUse, err := execCommand(ctx, `df -h / | tail -1 | awk '{print $5}'`)
	if err != nil {
		return fmt.Errorf("getting disk usage: %w", err)
	}
	s.DiskUse = diskUse

	diskSize, err := execCommand(ctx, `df -h / | tail -1 | awk '{print $2}'`)
	if err != nil {
		return fmt.Errorf("getting disk size: %w", err)
	}
	s.DiskSize = diskSize

	return nil
}

// calculateDiskPool calculates total disk pool size and usage from disk* devices.
func (s *SystemInfo) calculateDiskPool(ctx context.Context) {
	output := execCommandSafe(ctx, "df -h | grep /dev/disk")
	if output == "" {
		s.DiskPoolSize = config.CrossMark
		s.DiskPoolUsed = config.CrossMark
		return
	}

	var totalDisk, usedDisk float64

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		sizeStr := fields[1]
		usedStr := fields[2]

		// Check if size is in terabytes
		if strings.HasSuffix(sizeStr, "Ti") || strings.HasSuffix(sizeStr, "T") {
			size, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSuffix(sizeStr, "Ti"), "T"), 64)
			if err == nil {
				totalDisk += size
			}

			used, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSuffix(usedStr, "Ti"), "T"), 64)
			if err == nil {
				usedDisk += used
			}
		}
	}

	if totalDisk > 0 {
		s.DiskPoolSize = fmt.Sprintf("%.1fTB", totalDisk)
	} else {
		s.DiskPoolSize = config.CrossMark
	}

	if usedDisk > 0 {
		s.DiskPoolUsed = fmt.Sprintf("%.1fTB", usedDisk)
	} else {
		s.DiskPoolUsed = config.CrossMark
	}
}
