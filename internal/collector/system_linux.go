//go:build linux
// +build linux

package collector

import (
	"context"
	"os"
	"strings"
)

func gatherSystem(ctx context.Context, info *SystemInfo) {
	// Hostname
	info.Hostname = execCommandSafe(ctx, "hostname -s")
	if info.Hostname == "" {
		info.Hostname = execCommandSafe(ctx, "uname -n")
	}

	// Current user
	info.User = execCommandSafe(ctx, "id -un")
	if info.User == "" {
		info.User = os.Getenv("USER")
	}

	// OS: read PRETTY_NAME from /etc/os-release
	pretty := execCommandSafe(ctx, `grep -Po '(?<=^PRETTY_NAME=).+' /etc/os-release | tr -d '"'`)
	if pretty == "" {
		pretty = "Linux"
	}
	info.OS = pretty

	// Kernel: "Linux 5.15.0-91-generic x86_64"
	info.Kernel = execCommandSafe(ctx, "uname -srm")

	// Shell from environment
	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		info.Shell = parts[len(parts)-1]
	}
}
