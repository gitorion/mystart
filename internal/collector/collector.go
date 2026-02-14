package collector

import (
	"context"
	"fmt"
)

// GatherAll collects all system information.
// This is the main entry point for gathering system data.
func (s *SystemInfo) GatherAll(ctx context.Context) error {
	if err := s.GatherBasicInfo(ctx); err != nil {
		return fmt.Errorf("gathering basic info: %w", err)
	}

	if err := s.GatherCPUInfo(ctx); err != nil {
		return fmt.Errorf("gathering CPU info: %w", err)
	}

	if err := s.GatherMemoryInfo(ctx); err != nil {
		return fmt.Errorf("gathering memory info: %w", err)
	}

	if err := s.GatherDiskInfo(ctx); err != nil {
		return fmt.Errorf("gathering disk info: %w", err)
	}

	if err := s.GatherNetworkInfo(ctx); err != nil {
		return fmt.Errorf("gathering network info: %w", err)
	}

	if err := s.GatherUserInfo(ctx); err != nil {
		return fmt.Errorf("gathering user info: %w", err)
	}

	if err := s.GatherUptimeInfo(ctx); err != nil {
		return fmt.Errorf("gathering uptime info: %w", err)
	}

	// Gather optional host-specific information
	s.GatherHostSpecificInfo(ctx)

	return nil
}
