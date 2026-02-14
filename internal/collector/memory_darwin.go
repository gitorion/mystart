//go:build darwin
// +build darwin

package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// GatherMemoryInfo collects memory and swap information on macOS.
func (s *SystemInfo) GatherMemoryInfo(ctx context.Context) error {
	if err := s.gatherRAMInfo(ctx); err != nil {
		return err
	}

	if err := s.gatherSwapInfo(ctx); err != nil {
		return err
	}

	return nil
}

// gatherRAMInfo collects RAM usage information on macOS.
func (s *SystemInfo) gatherRAMInfo(ctx context.Context) error {
	// Get total memory in bytes
	memTotalStr, err := execCommand(ctx, "sysctl -n hw.memsize")
	if err != nil {
		return fmt.Errorf("getting memory total: %w", err)
	}

	memTotal, err := strconv.ParseFloat(memTotalStr, 64)
	if err != nil {
		return fmt.Errorf("%w: memory total", ErrParseFailure)
	}

	// Convert bytes to GB
	s.MemTotal = fmt.Sprintf("%.2fG", memTotal/1024/1024/1024)

	// Get memory usage using vm_stat
	vmStatOutput, err := execCommand(ctx, "vm_stat")
	if err != nil {
		return fmt.Errorf("getting vm_stat: %w", err)
	}

	// Parse vm_stat output
	pageSize := 4096.0 // Default page size
	var pagesActive, pagesWired, pagesCompressed float64

	for _, line := range strings.Split(vmStatOutput, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value := strings.TrimSuffix(fields[len(fields)-1], ".")
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}

		if strings.Contains(line, "Pages active") {
			pagesActive = val
		} else if strings.Contains(line, "Pages wired") {
			pagesWired = val
		} else if strings.Contains(line, "Pages occupied by compressor") {
			pagesCompressed = val
		}
	}

	// Calculate used memory (active + wired + compressed)
	usedBytes := (pagesActive + pagesWired + pagesCompressed) * pageSize
	s.MemUsed = fmt.Sprintf("%.2fG", usedBytes/1024/1024/1024)

	return nil
}

// gatherSwapInfo collects swap usage information on macOS.
func (s *SystemInfo) gatherSwapInfo(ctx context.Context) error {
	// Get swap usage using sysctl
	swapUsageOutput, err := execCommand(ctx, "sysctl -n vm.swapusage")
	if err != nil {
		// Swap might not be available
		s.SwapTotal = "0.00G"
		s.SwapUsed = "0.00G"
		return nil
	}

	// Parse output like: "total = 2048.00M  used = 671.75M  free = 1376.25M"
	var totalSwap, usedSwap float64
	var totalUnit, usedUnit string

	parts := strings.Split(swapUsageOutput, " ")
	for i, part := range parts {
		if part == "total" && i+2 < len(parts) {
			value := parts[i+2]
			if strings.HasSuffix(value, "M") {
				totalSwap, _ = strconv.ParseFloat(strings.TrimSuffix(value, "M"), 64)
				totalUnit = "M"
			} else if strings.HasSuffix(value, "G") {
				totalSwap, _ = strconv.ParseFloat(strings.TrimSuffix(value, "G"), 64)
				totalUnit = "G"
			}
		} else if part == "used" && i+2 < len(parts) {
			value := parts[i+2]
			if strings.HasSuffix(value, "M") {
				usedSwap, _ = strconv.ParseFloat(strings.TrimSuffix(value, "M"), 64)
				usedUnit = "M"
			} else if strings.HasSuffix(value, "G") {
				usedSwap, _ = strconv.ParseFloat(strings.TrimSuffix(value, "G"), 64)
			}
		}
	}

	// Convert to GB
	if totalUnit == "M" {
		totalSwap = totalSwap / 1024
	}
	if usedUnit == "M" {
		usedSwap = usedSwap / 1024
	}

	s.SwapTotal = fmt.Sprintf("%.2fG", totalSwap)
	s.SwapUsed = fmt.Sprintf("%.2fG", usedSwap)

	return nil
}
