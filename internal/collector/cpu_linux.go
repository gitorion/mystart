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
	// Read CPU stats from /proc/stat (most reliable method)
	// First snapshot
	stat1, err := execCommand(ctx, `cat /proc/stat | grep "^cpu " | awk '{print $2" "$3" "$4" "$5" "$6" "$7" "$8}'`)
	if err != nil {
		return fmt.Errorf("reading /proc/stat: %w", err)
	}

	// Sleep briefly to get a delta
	_, _ = execCommand(ctx, "sleep 0.2")

	// Second snapshot
	stat2, err := execCommand(ctx, `cat /proc/stat | grep "^cpu " | awk '{print $2" "$3" "$4" "$5" "$6" "$7" "$8}'`)
	if err != nil {
		return fmt.Errorf("reading /proc/stat: %w", err)
	}

	// Parse both snapshots
	fields1 := strings.Fields(stat1)
	fields2 := strings.Fields(stat2)

	if len(fields1) < 4 || len(fields2) < 4 {
		return fmt.Errorf("%w: invalid /proc/stat format", ErrInvalidOutput)
	}

	// Calculate totals and idle for both snapshots
	var total1, idle1, total2, idle2 float64

	for i := 0; i < len(fields1) && i < 7; i++ {
		val, _ := strconv.ParseFloat(fields1[i], 64)
		total1 += val
		if i == 3 { // idle is the 4th field
			idle1 = val
		}
	}

	for i := 0; i < len(fields2) && i < 7; i++ {
		val, _ := strconv.ParseFloat(fields2[i], 64)
		total2 += val
		if i == 3 { // idle is the 4th field
			idle2 = val
		}
	}

	// Calculate the delta
	totalDelta := total2 - total1
	idleDelta := idle2 - idle1

	if totalDelta == 0 {
		s.CPUUsage = "0.00%"
		return nil
	}

	// CPU usage percentage
	cpuUsage := 100.0 * (totalDelta - idleDelta) / totalDelta
	s.CPUUsage = fmt.Sprintf("%.2f%%", cpuUsage)

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
