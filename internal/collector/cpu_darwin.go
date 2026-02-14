//go:build darwin
// +build darwin

package collector

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// GatherCPUInfo collects CPU-related information on macOS.
func (s *SystemInfo) GatherCPUInfo(ctx context.Context) error {
	if err := s.gatherCPUCores(ctx); err != nil {
		return err
	}

	if err := s.calculateCPUSpeed(ctx); err != nil {
		return err
	}

	if err := s.calculateCPUUsage(ctx); err != nil {
		return err
	}

	if err := s.gatherLoadAverage(ctx); err != nil {
		return err
	}

	return nil
}

// gatherCPUCores retrieves the number of CPU cores on macOS.
func (s *SystemInfo) gatherCPUCores(ctx context.Context) error {
	coresStr, err := execCommand(ctx, "sysctl -n hw.physicalcpu")
	if err != nil {
		return fmt.Errorf("getting cpu cores: %w", err)
	}

	cores, err := strconv.Atoi(coresStr)
	if err != nil {
		return fmt.Errorf("%w: cpu cores", ErrParseFailure)
	}

	s.CPUCores = cores
	return nil
}

// calculateCPUSpeed calculates CPU speed on macOS.
func (s *SystemInfo) calculateCPUSpeed(ctx context.Context) error {
	// Get logical CPU count
	s.CPUThreads = runtime.NumCPU()

	// Try to get CPU frequency in Hz
	freqStr := execCommandSafe(ctx, "sysctl -n hw.cpufrequency")
	if freqStr != "" {
		freq, err := strconv.ParseFloat(freqStr, 64)
		if err == nil {
			// Convert Hz to GHz
			s.CPUHz = freq / 1000000000.0
			return nil
		}
	}

	// Fallback: try to get max frequency for Apple Silicon
	maxFreqStr := execCommandSafe(ctx, "sysctl -n hw.cpufrequency_max")
	if maxFreqStr != "" {
		freq, err := strconv.ParseFloat(maxFreqStr, 64)
		if err == nil {
			// Convert Hz to GHz
			s.CPUHz = freq / 1000000000.0
			return nil
		}
	}

	// Fallback: parse from CPU brand string
	brandStr := execCommandSafe(ctx, "sysctl -n machdep.cpu.brand_string")
	if brandStr != "" {
		// Try to parse frequency from brand string (e.g., "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz")
		if strings.Contains(brandStr, "GHz") {
			parts := strings.Split(brandStr, "@")
			if len(parts) == 2 {
				ghzPart := strings.TrimSpace(parts[1])
				ghzPart = strings.TrimSuffix(ghzPart, "GHz")
				ghzPart = strings.TrimSpace(ghzPart)
				speed, err := strconv.ParseFloat(ghzPart, 64)
				if err == nil {
					s.CPUHz = speed
					return nil
				}
			}
		}

		// For Apple Silicon without frequency in brand string, use a reasonable default
		if strings.Contains(brandStr, "Apple") {
			// Apple Silicon doesn't expose frequency, use a placeholder
			s.CPUHz = 0.0
			return nil
		}
	}

	// If all methods fail, set to 0 (unknown)
	s.CPUHz = 0.0
	return nil
}

// calculateCPUUsage calculates current CPU usage percentage on macOS.
func (s *SystemInfo) calculateCPUUsage(ctx context.Context) error {
	// Use top command to get CPU usage
	output, err := execCommand(ctx, "ps -A -o %cpu | awk '{s+=$1} END {print s}'")
	if err != nil {
		return fmt.Errorf("getting cpu usage: %w", err)
	}

	cpuUsed, err := strconv.ParseFloat(output, 64)
	if err != nil {
		return fmt.Errorf("%w: cpu usage", ErrParseFailure)
	}

	if s.CPUThreads == 0 {
		return fmt.Errorf("%w: cpu threads not initialized", ErrInvalidOutput)
	}

	cpuUsage := cpuUsed / float64(s.CPUThreads)
	s.CPUUsage = fmt.Sprintf("%.2f %%", cpuUsage)

	return nil
}

// gatherLoadAverage retrieves system load average on macOS.
func (s *SystemInfo) gatherLoadAverage(ctx context.Context) error {
	loadAvg, err := execCommand(ctx, "sysctl -n vm.loadavg | awk '{print $2 \" \" $3 \" \" $4}'")
	if err != nil {
		return fmt.Errorf("getting load average: %w", err)
	}

	s.LoadAvg = loadAvg
	return nil
}
