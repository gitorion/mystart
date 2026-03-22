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
	// valueWidth = innerWidth - 4 (indent) - labelWidth - 2 (sep) = 50
	valueWidth = 50
)

// ─────────────────────────────────────────────────────────────
// Colour palette
// ─────────────────────────────────────────────────────────────

var (
	// Structural
	cBorder = color.New(color.FgCyan)
	cDim    = color.New(color.FgHiBlack)

	// Labels & values
	cLabel = color.New(color.FgHiWhite, color.Bold)
	cValue = color.New(color.FgHiWhite)

	// Title area
	cTitleDeco = color.New(color.FgHiYellow, color.Bold)
	cTitleText = color.New(color.FgHiCyan, color.Bold)

	// Section accent colours — each section gets its own
	cSystem   = color.New(color.FgHiCyan, color.Bold)
	cProc     = color.New(color.FgHiYellow, color.Bold)
	cMem      = color.New(color.FgHiMagenta, color.Bold)
	cStorage  = color.New(color.FgHiBlue, color.Bold)
	cNetwork  = color.New(color.FgHiGreen, color.Bold)
	cSessions = color.New(color.FgHiRed, color.Bold)

	// Progress bar fills (Hi variants for vibrancy)
	cBarGreen  = color.New(color.FgHiGreen)
	cBarYellow = color.New(color.FgHiYellow)
	cBarRed    = color.New(color.FgHiRed)

	// Semantic value colours
	cHostname = color.New(color.FgHiGreen, color.Bold)
	cUser     = color.New(color.FgHiCyan, color.Bold)
	cIP       = color.New(color.FgHiCyan)
	cUptime   = color.New(color.FgHiGreen)
	cOSInfo   = color.New(color.FgHiYellow)
)

// ─────────────────────────────────────────────────────────────
// Public entry point
// ─────────────────────────────────────────────────────────────

// Render prints the full system status dashboard.
func Render(info *collector.SystemInfo) {
	p := fmt.Println

	// ── Title block ─────────────────────────────────────────
	p(topBorder())
	p(emptyRow())
	p(titleLine())
	p(subtitleLine(info.Hostname, info.User))
	p(emptyRow())
	p(divider())

	// ── SYSTEM ──────────────────────────────────────────────
	p(sectionHeader("SYSTEM", cSystem))
	p(rowC("Hostname", info.Hostname, cHostname))
	p(rowC("User", info.User, cUser))
	p(rowC("OS", info.OS, cOSInfo))
	p(rowC("Kernel", info.Kernel, cOSInfo))
	if info.Shell != "" {
		p(row("Shell", info.Shell))
	}
	p(rowC("Uptime", formatUptime(info), cUptime))
	p(divider())

	// ── PROCESSOR ───────────────────────────────────────────
	p(sectionHeader("PROCESSOR", cProc))
	if info.CPUModel != "" {
		p(row("Model", info.CPUModel))
	}
	p(row("Cores / Threads", formatCores(info)))
	p(barRowPct("Usage", info.CPUUsage))
	p(row("Load Average", formatLoad(info)))
	p(divider())

	// ── MEMORY ──────────────────────────────────────────────
	p(sectionHeader("MEMORY", cMem))
	p(barRowGB("RAM", info.MemUsedGB, info.MemTotalGB))
	if info.SwapTotalGB > 0 {
		p(barRowGB("Swap", info.SwapUsedGB, info.SwapTotalGB))
	}
	p(divider())

	// ── STORAGE ─────────────────────────────────────────────
	p(sectionHeader("STORAGE", cStorage))
	for _, m := range info.DiskMounts {
		p(barRowDisk(m))
	}
	p(divider())

	// ── NETWORK ─────────────────────────────────────────────
	p(sectionHeader("NETWORK", cNetwork))
	if info.IPv4 != "" {
		p(rowC("IPv4", info.IPv4, cIP))
	} else {
		p(rowC("IPv4", "unavailable", cDim))
	}
	if info.IPv6 != "" {
		p(rowC("IPv6", info.IPv6, cIP))
	} else {
		p(rowC("IPv6", "unavailable", cDim))
	}
	p(divider())

	// ── SESSIONS ────────────────────────────────────────────
	p(sectionHeader("SESSIONS", cSessions))
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

// row renders a label → value line in default colours.
func row(label, value string) string {
	return rowC(label, value, cValue)
}

// rowC renders a label → value line with a custom value colour.
func rowC(label, value string, vc *color.Color) string {
	vr := []rune(value)
	if len(vr) > valueWidth {
		vr = vr[:valueWidth-3]
		value = string(vr) + "..."
	}
	padding := valueWidth - len([]rune(value))
	return fmt.Sprintf("%s    %s  %s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		vc.Sprint(value),
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// barRowPct renders a progress bar for a percentage value (e.g. CPU usage).
func barRowPct(label string, pct float64) string {
	bar, barColor := makeBar(pct)
	pctStr := fmt.Sprintf("%.1f%%", pct)
	text := "  " + barColor.Sprint(pctStr)
	textVisual := 2 + len(pctStr)
	padding := valueWidth - barWidth - textVisual
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s    %s  %s%s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		bar,
		text,
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
	bar, barColor := makeBar(pct)
	sizeStr := fmt.Sprintf("  %.1f / %.1f GB  ", used, total)
	pctStr := fmt.Sprintf("(%.1f%%)", pct)
	textVisual := len(sizeStr) + len(pctStr)
	padding := valueWidth - barWidth - textVisual
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s    %s  %s%s%s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		bar,
		cValue.Sprint(sizeStr),
		barColor.Sprint(pctStr),
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// barRowDisk renders a progress bar for a DiskMount, auto-selecting GB or TB.
func barRowDisk(m collector.DiskMount) string {
	bar, barColor := makeBar(m.UsedPercent)
	var sizeStr, pctStr string
	if m.TotalGB >= 1000 {
		sizeStr = fmt.Sprintf("  %.1f / %.1f TB  ",
			m.UsedGB/1024, m.TotalGB/1024)
	} else {
		sizeStr = fmt.Sprintf("  %.1f / %.1f GB  ",
			m.UsedGB, m.TotalGB)
	}
	pctStr = fmt.Sprintf("(%.1f%%)", m.UsedPercent)

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

	textVisual := len(sizeStr) + len(pctStr)
	padding := valueWidth - barWidth - textVisual
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s    %s  %s%s%s%s%s",
		cBorder.Sprint("│"),
		cLabel.Sprintf("%-18s", label),
		bar,
		cValue.Sprint(sizeStr),
		barColor.Sprint(pctStr),
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

func emptyRow() string {
	return cBorder.Sprint("│") + strings.Repeat(" ", innerWidth) + cBorder.Sprint("│")
}

func titleLine() string {
	// Visual: "◈  SYSTEM STATUS  ◈" = 19 chars
	const visualLen = 19
	total := innerWidth - visualLen
	left := total / 2
	right := total - left
	return fmt.Sprintf("%s%s%s  %s  %s%s%s",
		cBorder.Sprint("│"),
		strings.Repeat(" ", left),
		cTitleDeco.Sprint("◈"),
		cTitleText.Sprint("SYSTEM STATUS"),
		cTitleDeco.Sprint("◈"),
		strings.Repeat(" ", right),
		cBorder.Sprint("│"),
	)
}

func subtitleLine(hostname, user string) string {
	sep := "  " + cDim.Sprint("·") + "  "
	hostC := cHostname.Sprint(hostname)
	userC := cUser.Sprint(user)
	visualLen := len([]rune(hostname)) + 5 + len([]rune(user)) // "  ·  " = 5
	total := innerWidth - visualLen
	left := total / 2
	right := total - left
	if left < 0 {
		left, right = 0, 0
	}
	return fmt.Sprintf("%s%s%s%s%s%s%s",
		cBorder.Sprint("│"),
		strings.Repeat(" ", left),
		hostC, sep, userC,
		strings.Repeat(" ", right),
		cBorder.Sprint("│"),
	)
}

func sectionHeader(title string, accent *color.Color) string {
	icon := accent.Sprint("◆")
	titleC := accent.Sprint(title)
	visualLen := 2 + 2 + len([]rune(title)) // "  " + "◆ " + title
	padding := innerWidth - visualLen
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s  %s %s%s%s",
		cBorder.Sprint("│"),
		icon, titleC,
		strings.Repeat(" ", padding),
		cBorder.Sprint("│"),
	)
}

// ─────────────────────────────────────────────────────────────
// Progress bar
// ─────────────────────────────────────────────────────────────

// makeBar returns a coloured progress bar and the fill colour used.
func makeBar(pct float64) (string, *color.Color) {
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
		fill = cBarRed
	case pct >= 60:
		fill = cBarYellow
	default:
		fill = cBarGreen
	}
	bar := fill.Sprint(strings.Repeat("█", filled)) +
		cDim.Sprint(strings.Repeat("░", barWidth-filled))
	return bar, fill
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
