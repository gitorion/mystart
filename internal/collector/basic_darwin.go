//go:build darwin
// +build darwin

package collector

import (
	"context"
	"fmt"

	"github.com/orion/mystart/internal/config"
)

// GatherBasicInfo collects basic system identification information on macOS.
func (s *SystemInfo) GatherBasicInfo(ctx context.Context) error {
	var err error

	// Get macOS version
	s.Uname, err = execCommand(ctx, "sw_vers -productVersion")
	if err != nil {
		return fmt.Errorf("getting OS version: %w", err)
	}
	s.Uname = "macOS " + s.Uname

	// Get OS name
	productName, err := execCommand(ctx, "sw_vers -productName")
	if err != nil {
		s.Distro = "macOS"
	} else {
		s.Distro = productName
	}

	s.Host, err = execCommand(ctx, "hostname -s")
	if err != nil {
		return fmt.Errorf("getting hostname: %w", err)
	}

	s.User, err = execCommand(ctx, "whoami")
	if err != nil {
		return fmt.Errorf("getting username: %w", err)
	}

	// Set host task description
	if task, ok := config.HostTasks[s.Host]; ok {
		s.HostTask = task
	}

	return nil
}
