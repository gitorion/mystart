//go:build linux
// +build linux

package collector

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"
)

// GatherCPUInfo collects CPU-related information on Linux.
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

// gatherCPUCores retrieves the number of CPU cores.
func (s *SystemInfo) gatherCPUCores(ctx context.Context) error {
	coresStr, err := execCommand(ctx, `cat /proc/cpuinfo | grep "cpu cores" | head -n 1 | awk '{print $4}'`)
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

// calculateCPUSpeed calculates average CPU speed across all threads.
func (s *SystemInfo) calculateCPUSpeed(ctx context.Context) error {
	output, err := execCommand(ctx, "cat /proc/cpuinfo | grep MHz | awk '{print $4}'")
	if err != nil {
		return fmt.Errorf("getting cpu speed: %w", err)
	}

	var totalSpeed float64
	threadCount := 0

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		speed, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			continue
		}
		totalSpeed += speed
		threadCount++
	}

	if threadCount == 0 {
		return fmt.Errorf("%w: no CPU speed information found", ErrInvalidOutput)
	}

	s.CPUThreads = threadCount
	// Convert MHz to GHz
	s.CPUHz = (totalSpeed / float64(threadCount)) / 1000.0

	return nil
}

// calculateCPUUsage calculates current CPU usage percentage.
func (s *SystemInfo) calculateCPUUsage(ctx context.Context) error {
	cpuUsedStr, err := execCommand(ctx, `ps -eo pcpu | awk '{tot=tot+$1} END {print tot}'`)
	if err != nil {
		return fmt.Errorf("getting cpu usage: %w", err)
	}

	cpuUsed, err := strconv.ParseFloat(cpuUsedStr, 64)
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

// gatherLoadAverage retrieves system load average.
func (s *SystemInfo) gatherLoadAverage(ctx context.Context) error {
	loadAvg, err := execCommand(ctx, `cat /proc/loadavg | awk '{print $1 " " $2 " " $3 " " $4}'`)
	if err != nil {
		return fmt.Errorf("getting load average: %w", err)
	}

	s.LoadAvg = loadAvg
	return nil
}
