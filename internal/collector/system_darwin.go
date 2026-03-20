//go:build darwin
// +build darwin

package collector

import (
	"context"
	"os"
	"strings"
)

func gatherSystem(ctx context.Context, info *SystemInfo) {
	// Hostname
	info.Hostname = execCommandSafe(ctx, "hostname -s")

	// Current user
	info.User = execCommandSafe(ctx, "whoami")
	if info.User == "" {
		info.User = os.Getenv("USER")
	}

	// OS: "macOS Sonoma 14.2.1" style
	name := execCommandSafe(ctx, "sw_vers -productName")
	version := execCommandSafe(ctx, "sw_vers -productVersion")
	if name == "" {
		name = "macOS"
	}
	info.OS = name + " " + version

	// Kernel: "Darwin 23.2.0 arm64"
	info.Kernel = execCommandSafe(ctx, "uname -srm")

	// Shell from environment
	shell := os.Getenv("SHELL")
	if shell != "" {
		// Show only the binary name's basename
		parts := strings.Split(shell, "/")
		info.Shell = parts[len(parts)-1]
	}
}
