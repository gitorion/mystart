//go:build darwin
// +build darwin

package collector

import (
	"context"
	"strings"
)

func gatherNetwork(ctx context.Context, info *SystemInfo) {
	// Determine the default route interface
	iface := execCommandSafe(ctx, "route -n get default 2>/dev/null | awk '/interface:/{print $2}'")

	if iface != "" {
		// IPv4
		ipv4 := execCommandSafe(ctx, "ipconfig getifaddr "+iface)
		if ipv4 == "" {
			ipv4 = execCommandSafe(ctx, "ifconfig "+iface+" | awk '/inet /{print $2}'")
		}
		info.IPv4 = strings.TrimSpace(ipv4)

		// IPv6: prefer non-link-local (fe80::)
		ipv6 := execCommandSafe(ctx,
			"ifconfig "+iface+" | awk '/inet6 / && !/fe80:/{print $2; exit}'")
		if ipv6 == "" {
			// Fall back to link-local
			ipv6 = execCommandSafe(ctx,
				"ifconfig "+iface+" | awk '/inet6 fe80:/{print $2; exit}'")
		}
		// Strip %interface suffix
		if i := strings.Index(ipv6, "%"); i >= 0 {
			ipv6 = ipv6[:i]
		}
		info.IPv6 = strings.TrimSpace(ipv6)
	}
}
