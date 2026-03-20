package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gitorion/mystart/internal/collector"
)

// Layout constants (all values are visual character widths).
const (
	boxWidth   = 76
	innerWidth = boxWidth - 2 // 74, between the vertical borders
	labelWidth = 18
	barWidth   = 20
	// valueWidth = innerWidth - 4 (indent) - labelWidth - 2 (sep) - 0 (no trailing space needed before │)
	// 74 - 4 - 18 - 2 = 50
	valueWidth = 50
)

// Colour palette
var (
	cBorder   = color.New(color.FgCyan)
	cSection  = color.New(color.FgYellow, color.Bold)
	cLabel    = color.New(color.FgMagenta)
	cValue    = color.New(color.FgHiWhite)
	cDim      = color.New(color.FgHiBlack)
	cTitle    = color.New(color.FgCyan, color.Bold)
	cSubtitle = color.New(color.FgHiWhite)
	cGreen    = color.New(color.FgGreen)
	cYellow   = color.New(color.FgYellow)
	cRed      = color.New(color.FgRed)
)

// ─────────────────────────────────────────────────────────────
// Public entry point
// ─────────────────────────────────────────────────────────────

// Render prints the full system status dashboard.
func Render(info *collector.SystemInfo) {
	p := fmt.Println

	p(topBorder())
	p(titleLine("◈  SYSTEM STATUS  ◈"))
	p(subtitleLine(info.Hostname + "  ·  " + info.User))
	p(divider())

	// ── SYSTEM ──────────────────────────────────────────────
	p(sectionHeader("SYSTEM"))
	p(row("Hostname", info.Hostname))
	p(row("User", info.User))
	p(row("OS", info.OS))
	p(row("Kernel", info.Kernel))
	if info.Shell != "" {
		p(row("Shell", info.Shell))
	}
	p(row("Uptime", formatUptime(info)))
	p(divider())

	// ── PROCESSOR ───────────────────────────────────────────
	p(sectionHeader("PROCESSOR"))
	if info.CPUModel != "" {
		p(row("Model", info.CPUModel))
	}
	p(row("Cores / Threads", formatCores(info)))
	p(barRowPct("Usage", info.CPUUsage))
	p(row("Load Average", formatLoad(info)))
	p(divider())

	// ── MEMORY ──────────────────────────────────────────────
	p(sectionHeader("MEMORY"))
	p(barRowGB("RAM", info.MemUsedGB, info.MemTotalGB))
	if info.SwapTotalGB > 0 {
		p(barRowGB("Swap", info.SwapUsedGB, info.SwapTotalGB))
	}
	p(divider())

	// ── STORAGE ─────────────────────────────────────────────
	p(sectionHeader("STORAGE"))
	for _, m := range info.DiskMounts {
		p(barRowDisk(m))
	}
	p(divider())

	// ── NETWORK ─────────────────────────────────────────────
	p(sectionHeader("NETWORK"))
	if info.IPv4 != "" {
		p(row("IPv4", info.IPv4))
	} else {
		p(row("IPv4", "unavailable"))
	}
	if info.IPv6 != "" {
		p(row("IPv6", info.IPv6))
	} else {
		p(row("IPv6", "unavailable"))
	}
	p(divider())

	// ── SESSIONS ────────────────────────────────────────────
	p(sectionHeader("SESSIONS"))
	p(row("Users", formatUsers(info)))
	p(row("Processes", formatProcesses(info)))
	if info.LastLogin != "" {
		p(row("Last Login", info.LastLogin))
	}
	p(bottomBorder())
}

// ─────────────────────────────────────────────────────────────
// Row builders
// ─────────────────────────────────────────────────────────────

// row renders a plain label → value line.
func row(label, value string) string {
	vr := []rune(value)
	if len(vr) > valueWidth {
		vr = vr[:valueWidth-3]
		value = string(vr) + "..."
	}
	padding := valueWidth - len([]rune(value))
	return fmt.Sprintf("%s    %s  %s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		cValue.Sprint(value),
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// barRowPct renders a progress bar for a plain percentage value (e.g. CPU usage).
func barRowPct(label string, pct float64) string {
	bar := makeBar(pct)
	text := fmt.Sprintf("  %.1f%%", pct)
	padding := valueWidth - barWidth - len(text)
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s    %s  %s%s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		bar,
		cValue.Sprint(text),
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// barRowGB renders a progress bar for a used/total value in GB.
func barRowGB(label string, used, total float64) string {
	pct := 0.0
	if total > 0 {
		pct = used / total * 100
	}
	bar := makeBar(pct)
	text := fmt.Sprintf("  %.1f / %.1f GB  (%.1f%%)", used, total, pct)
	padding := valueWidth - barWidth - len(text)
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s    %s  %s%s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		bar,
		cValue.Sprint(text),
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// barRowDisk renders a progress bar for a DiskMount, auto-selecting GB or TB.
func barRowDisk(m collector.DiskMount) string {
	bar := makeBar(m.UsedPercent)
	var text string
	if m.TotalGB >= 1000 {
		text = fmt.Sprintf("  %.1f / %.1f TB  (%.1f%%)",
			m.UsedGB/1024, m.TotalGB/1024, m.UsedPercent)
	} else {
		text = fmt.Sprintf("  %.1f / %.1f GB  (%.1f%%)",
			m.UsedGB, m.TotalGB, m.UsedPercent)
	}

	// Shorten well-known long paths, then truncate if still too long
	label := m.Path
	switch label {
	case "/System/Volumes/Data":
		label = "/data (APFS)"
	}
	lr := []rune(label)
	if len(lr) > labelWidth {
		label = "…" + string(lr[len(lr)-(labelWidth-1):])
	}

	padding := valueWidth - barWidth - len(text)
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s    %s  %s%s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		bar,
		cValue.Sprint(text),
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// ─────────────────────────────────────────────────────────────
// Structural elements
// ─────────────────────────────────────────────────────────────

func topBorder() string {
	return cBorder.Sprint("╭" + strings.Repeat("─", innerWidth) + "╮")
}

func bottomBorder() string {
	return cBorder.Sprint("╰" + strings.Repeat("─", innerWidth) + "╯")
}

func divider() string {
	return cBorder.Sprint("╠" + strings.Repeat("═", innerWidth) + "╣")
}

func titleLine(text string) string {
	return centeredLine(text, cTitle)
}

func subtitleLine(text string) string {
	return centeredLine(text, cSubtitle)
}

func centeredLine(text string, c *color.Color) string {
	textLen := len([]rune(text))
	total := innerWidth - textLen
	left := total / 2
	right := total - left
	if left < 0 {
		left, right = 0, 0
	}
	return fmt.Sprintf("%s%s%s%s%s",
		cBorder.Sprint("│"),
		strings.Repeat(" ", left),
		c.Sprint(text),
		strings.Repeat(" ", right),
		cBorder.Sprint("│"),
	)
}

func sectionHeader(title string) string {
	content := "◆ " + title
	padding := innerWidth - 2 - len([]rune(content)) // 2 = leading spaces "  "
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s  %s%s%s",
		cBorder.Sprint("│"),
		cSection.Sprint(content),
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// ─────────────────────────────────────────────────────────────
// Progress bar
// ─────────────────────────────────────────────────────────────

func makeBar(pct float64) string {
	filled := int(pct / 100.0 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}
	var fill *color.Color
	switch {
	case pct >= 80:
		fill = cRed
	case pct >= 60:
		fill = cYellow
	default:
		fill = cGreen
	}
	return fill.Sprint(strings.Repeat("█", filled)) +
		cDim.Sprint(strings.Repeat("░", barWidth-filled))
}

// ─────────────────────────────────────────────────────────────
// Value formatters
// ─────────────────────────────────────────────────────────────

func formatUptime(info *collector.SystemInfo) string {
	parts := []string{}
	if info.UptimeDays > 0 {
		parts = append(parts, fmt.Sprintf("%d day%s", info.UptimeDays, plural(info.UptimeDays)))
	}
	if info.UptimeHours > 0 {
		parts = append(parts, fmt.Sprintf("%d hour%s", info.UptimeHours, plural(info.UptimeHours)))
	}
	if info.UptimeMinutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minute%s", info.UptimeMinutes, plural(info.UptimeMinutes)))
	}
	parts = append(parts, fmt.Sprintf("%d second%s", info.UptimeSeconds, plural(info.UptimeSeconds)))
	return strings.Join(parts, ", ")
}

func formatCores(info *collector.SystemInfo) string {
	if info.CPUHz > 0 {
		return fmt.Sprintf("%d physical · %d logical · %.2f GHz",
			info.CPUCores, info.CPUThreads, info.CPUHz)
	}
	return fmt.Sprintf("%d physical · %d logical", info.CPUCores, info.CPUThreads)
}

func formatLoad(info *collector.SystemInfo) string {
	return fmt.Sprintf("%.2f  %.2f  %.2f   (1m / 5m / 15m)",
		info.LoadAvg1, info.LoadAvg5, info.LoadAvg15)
}

func formatUsers(info *collector.SystemInfo) string {
	return fmt.Sprintf("%d logged in · %d active session%s",
		info.UsersLoggedIn, info.ActiveSessions, plural(info.ActiveSessions))
}

func formatProcesses(info *collector.SystemInfo) string {
	return fmt.Sprintf("%d user · %d total", info.ProcessesUser, info.ProcessesTotal)
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
