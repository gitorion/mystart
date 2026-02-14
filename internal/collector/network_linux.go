//go:build linux
// +build linux

package collector

import (
	"context"

	"github.com/orion/mystart/internal/config"
)

// GatherNetworkInfo collects network configuration information on Linux.
func (s *SystemInfo) GatherNetworkInfo(ctx context.Context) error {
	s.gatherIPv4(ctx)
	s.gatherIPv6(ctx)
	return nil
}

// gatherIPv4 retrieves the IPv4 address.
func (s *SystemInfo) gatherIPv4(ctx context.Context) {
	checkIPv4 := execCommandSafe(ctx, "ip route get 8.8.8.8 2>/dev/null")
	if checkIPv4 == "" {
		s.IPv4 = config.CrossMark
		return
	}

	ipv4 := execCommandSafe(ctx, "ip route get 8.8.8.8 | grep src | awk '{print $7}'")
	if ipv4 == "" {
		s.IPv4 = config.CrossMark
	} else {
		s.IPv4 = ipv4
	}
}

// gatherIPv6 retrieves the IPv6 address.
func (s *SystemInfo) gatherIPv6(ctx context.Context) {
	checkIPv6 := execCommandSafe(ctx, "ip route get 2001:4860:4860::8888 2>/dev/null")
	if checkIPv6 == "" {
		s.IPv6 = config.CrossMark
		return
	}

	ipv6 := execCommandSafe(ctx, "ip route get 2001:4860:4860::8888 | grep src | awk '{print $11}'")
	if ipv6 == "" {
		s.IPv6 = config.CrossMark
	} else {
		s.IPv6 = ipv6
	}
}
