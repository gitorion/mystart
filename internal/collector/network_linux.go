//go:build linux
// +build linux

package collector

import (
	"context"
	"strings"
)

func gatherNetwork(ctx context.Context, info *SystemInfo) {
	// IPv4 via routing table lookup
	ipv4Out := execCommandSafe(ctx, "ip route get 8.8.8.8 2>/dev/null")
	if ipv4Out != "" {
		// Look for "src <ip>"
		fields := strings.Fields(ipv4Out)
		for i, f := range fields {
			if f == "src" && i+1 < len(fields) {
				info.IPv4 = fields[i+1]
				break
			}
		}
	}

	// IPv6 via routing table lookup
	ipv6Out := execCommandSafe(ctx, "ip -6 route get 2001:4860:4860::8888 2>/dev/null")
	if ipv6Out != "" {
		fields := strings.Fields(ipv6Out)
		for i, f := range fields {
			if f == "src" && i+1 < len(fields) {
				info.IPv6 = fields[i+1]
				break
			}
		}
	}
}
