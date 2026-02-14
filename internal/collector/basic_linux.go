//go:build linux
// +build linux

package collector

import (
	"context"
	"fmt"

	"github.com/orion/mystart/internal/config"
)

// GatherBasicInfo collects basic system identification information on Linux.
func (s *SystemInfo) GatherBasicInfo(ctx context.Context) error {
	var err error

	s.Uname, err = execCommand(ctx, `cat /etc/os-release | grep "PRETTY" | cut -d'=' -f2-`)
	if err != nil {
		return fmt.Errorf("getting uname: %w", err)
	}

	s.Distro, err = execCommand(ctx, `grep -Po "(?<=^ID=).+" /etc/os-release | sed 's/"//g'`)
	if err != nil {
		return fmt.Errorf("getting distro: %w", err)
	}

	s.Host, err = execCommand(ctx, "uname -n")
	if err != nil {
		return fmt.Errorf("getting hostname: %w", err)
	}

	s.User, err = execCommand(ctx, "id -un")
	if err != nil {
		return fmt.Errorf("getting username: %w", err)
	}

	// Set host task description
	if task, ok := config.HostTasks[s.Host]; ok {
		s.HostTask = task
	}

	return nil
}
