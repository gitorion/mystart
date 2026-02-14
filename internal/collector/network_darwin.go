//go:build darwin
// +build darwin

package collector

import (
	"context"
	"strings"

	"github.com/orion/mystart/internal/config"
)

// GatherNetworkInfo collects network configuration information on macOS.
func (s *SystemInfo) GatherNetworkInfo(ctx context.Context) error {
	s.gatherIPv4(ctx)
	s.gatherIPv6(ctx)
	return nil
}

// gatherIPv4 retrieves the IPv4 address on macOS.
func (s *SystemInfo) gatherIPv4(ctx context.Context) {
	// Get the default route interface
	iface := execCommandSafe(ctx, "route -n get default 2>/dev/null | grep interface | awk '{print $2}'")
	if iface == "" {
		s.IPv4 = config.CrossMark
		return
	}

	// Get IP address for that interface
	ipv4 := execCommandSafe(ctx, "ipconfig getifaddr "+iface)
	if ipv4 == "" {
		// Try alternative method
		ipv4 = execCommandSafe(ctx, "ifconfig "+iface+" | grep 'inet ' | awk '{print $2}'")
	}

	if ipv4 == "" {
		s.IPv4 = config.CrossMark
	} else {
		s.IPv4 = strings.TrimSpace(ipv4)
	}
}

// gatherIPv6 retrieves the IPv6 address on macOS.
func (s *SystemInfo) gatherIPv6(ctx context.Context) {
	// Get the default route interface
	iface := execCommandSafe(ctx, "route -n get default 2>/dev/null | grep interface | awk '{print $2}'")
	if iface == "" {
		s.IPv6 = config.CrossMark
		return
	}

	// Get IPv6 address for that interface (excluding link-local)
	ipv6 := execCommandSafe(ctx, "ifconfig "+iface+" | grep 'inet6 ' | grep -v 'fe80:' | head -1 | awk '{print $2}'")

	if ipv6 == "" {
		s.IPv6 = config.CrossMark
	} else {
		s.IPv6 = strings.TrimSpace(ipv6)
	}
}
