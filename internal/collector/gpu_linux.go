//go:build linux
// +build linux

package collector

import (
	"context"
	"strings"
)

func gatherGPU(ctx context.Context, info *SystemInfo) {
	// Try NVIDIA first (most common discrete GPU)
	if gatherNvidiaGPU(ctx, info) {
		return
	}

	// Fall back to lspci for any GPU
	gatherLspciGPU(ctx, info)
}

func gatherNvidiaGPU(ctx context.Context, info *SystemInfo) bool {
	// nvidia-smi gives model, temp, usage, memory in one call
	out := execCommandSafe(ctx,
		`nvidia-smi --query-gpu=name,temperature.gpu,utilization.gpu,memory.total,memory.used,driver_version --format=csv,noheader,nounits 2>/dev/null`)
	if out == "" {
		return false
	}

	// Format: "NVIDIA GeForce RTX 3080, 62, 34, 10240, 3456, 535.129.03"
	fields := strings.SplitN(out, ", ", 6)
	if len(fields) < 1 {
		return false
	}

	info.GPUModel = strings.TrimSpace(fields[0])

	if len(fields) >= 2 && fields[1] != "" {
		info.GPUTemp = strings.TrimSpace(fields[1]) + "°C"
	}
	if len(fields) >= 3 && fields[2] != "" {
		info.GPUUsage = strings.TrimSpace(fields[2]) + "%"
	}
	if len(fields) >= 5 {
		total := strings.TrimSpace(fields[3])
		used := strings.TrimSpace(fields[4])
		if total != "" {
			info.GPUVram = used + " / " + total + " MB"
		}
	}
	if len(fields) >= 6 && fields[5] != "" {
		info.GPUDriver = strings.TrimSpace(fields[5])
	}

	return true
}

func gatherLspciGPU(ctx context.Context, info *SystemInfo) {
	// Try VGA compatible controller
	out := execCommandSafe(ctx, `lspci 2>/dev/null | grep -i 'vga\|3d\|display' | head -1`)
	if out == "" {
		return
	}

	// Format: "01:00.0 VGA compatible controller: NVIDIA Corporation ... [GeForce RTX 3080]"
	if idx := strings.Index(out, ": "); idx >= 0 {
		info.GPUModel = strings.TrimSpace(out[idx+2:])
	}
}
