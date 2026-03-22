//go:build darwin
// +build darwin

package collector

import (
	"context"
	"strings"
)

func gatherGPU(ctx context.Context, info *SystemInfo) {
	// system_profiler SPDisplaysDataType gives all GPU info on macOS
	out := execCommandSafe(ctx, "system_profiler SPDisplaysDataType 2>/dev/null")
	if out == "" {
		return
	}

	// Parse key fields from the output
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Chipset Model:") {
			info.GPUModel = strings.TrimSpace(strings.TrimPrefix(line, "Chipset Model:"))
		} else if strings.HasPrefix(line, "VRAM") && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				info.GPUVram = strings.TrimSpace(parts[1])
			}
		} else if strings.HasPrefix(line, "Metal Family:") || strings.HasPrefix(line, "Metal Support:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				info.GPUDriver = strings.TrimSpace(parts[1])
			}
		}
	}
}
