//go:build linux
// +build linux

package collector

import (
	"context"
	"fmt"
	"strconv"
)

// GatherMemoryInfo collects memory and swap information on Linux.
func (s *SystemInfo) GatherMemoryInfo(ctx context.Context) error {
	if err := s.gatherRAMInfo(ctx); err != nil {
		return err
	}

	if err := s.gatherSwapInfo(ctx); err != nil {
		return err
	}

	return nil
}

// gatherRAMInfo collects RAM usage information.
func (s *SystemInfo) gatherRAMInfo(ctx context.Context) error {
	// Memory total
	memTotalStr, err := execCommand(ctx, `cat /proc/meminfo | grep MemTotal | awk '{print $2}'`)
	if err != nil {
		return fmt.Errorf("getting memory total: %w", err)
	}

	memTotal, err := strconv.ParseFloat(memTotalStr, 64)
	if err != nil {
		return fmt.Errorf("%w: memory total", ErrParseFailure)
	}

	// Convert KB to GB
	s.MemTotal = fmt.Sprintf("%.2fG", memTotal/1024/1024)

	// Memory available
	memAvailStr, err := execCommand(ctx, `cat /proc/meminfo | grep MemAvailable | awk '{print $2}'`)
	if err != nil {
		return fmt.Errorf("getting memory available: %w", err)
	}

	memAvail, err := strconv.ParseFloat(memAvailStr, 64)
	if err != nil {
		return fmt.Errorf("%w: memory available", ErrParseFailure)
	}

	// Calculate used memory
	memUsed := (memTotal - memAvail) / 1024 / 1024
	s.MemUsed = fmt.Sprintf("%.2fG", memUsed)

	return nil
}

// gatherSwapInfo collects swap usage information.
func (s *SystemInfo) gatherSwapInfo(ctx context.Context) error {
	// Swap total
	swapTotalStr, err := execCommand(ctx, `cat /proc/meminfo | grep SwapTotal | awk '{print $2}'`)
	if err != nil {
		return fmt.Errorf("getting swap total: %w", err)
	}

	swapTotal, err := strconv.ParseFloat(swapTotalStr, 64)
	if err != nil {
		return fmt.Errorf("%w: swap total", ErrParseFailure)
	}

	// Convert KB to GB
	s.SwapTotal = fmt.Sprintf("%.2fG", swapTotal/1024/1024)

	// Swap free
	swapFreeStr, err := execCommand(ctx, `cat /proc/meminfo | grep SwapFree | awk '{print $2}'`)
	if err != nil {
		return fmt.Errorf("getting swap free: %w", err)
	}

	swapFree, err := strconv.ParseFloat(swapFreeStr, 64)
	if err != nil {
		return fmt.Errorf("%w: swap free", ErrParseFailure)
	}

	// Calculate used swap
	swapUsed := (swapTotal - swapFree) / 1024 / 1024
	s.SwapUsed = fmt.Sprintf("%.2fG", swapUsed)

	return nil
}
