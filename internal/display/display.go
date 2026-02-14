package display

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/orion/mystart/internal/collector"
	"github.com/orion/mystart/internal/config"
)

// Formatter handles the formatting and display of system information.
type Formatter struct {
	green   func(a ...interface{}) string
	magenta func(a ...interface{}) string
	cyan    func(a ...interface{}) string
	yellow  func(a ...interface{}) string
	blue    func(a ...interface{}) string
	red     func(a ...interface{}) string
	reset   func(a ...interface{}) string
}

// NewFormatter creates a new display formatter.
func NewFormatter() *Formatter {
	return &Formatter{
		green:   color.New(color.FgGreen).SprintFunc(),
		magenta: color.New(color.FgMagenta).SprintFunc(),
		cyan:    color.New(color.FgCyan).SprintFunc(),
		yellow:  color.New(color.FgYellow).SprintFunc(),
		blue:    color.New(color.FgBlue).SprintFunc(),
		red:     color.New(color.FgRed).SprintFunc(),
		reset:   color.New(color.Reset).SprintFunc(),
	}
}

// Display prints the formatted system information.
func (f *Formatter) Display(info *collector.SystemInfo) {
	messages := f.buildMessages(info)

	for _, msg := range messages {
		fmt.Println(msg)
	}
}

// buildMessages constructs the display messages from SystemInfo.
func (f *Formatter) buildMessages(info *collector.SystemInfo) []string {
	lineBorder := f.magenta("=======================================================================")

	messages := make([]string, 0, 35)

	// Header
	messages = append(messages, lineBorder)
	messages = append(messages, fmt.Sprintf("%sUser: %s%s\tHost: %s%s %s %s%s",
		f.reset(), f.green(info.User), f.reset(), f.green(info.Host),
		config.PointRight, info.HostTask, f.reset(), f.reset()))
	messages = append(messages, lineBorder)

	// Login information
	messages = append(messages, f.formatLineWithColor("Login details", info.ThisLog, f.cyan))

	// System information
	systemDetails := fmt.Sprintf("%s | %s", info.Distro, info.Uname)
	messages = append(messages, f.formatLineWithColor("System details", systemDetails, f.yellow))

	uptimeStr := fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		info.UptimeDays, info.UptimeHours, info.UptimeMinutes, info.UptimeSeconds)
	messages = append(messages, f.formatLineWithColor("System uptime", uptimeStr, f.cyan))
	messages = append(messages, f.formatLineWithColor("System load", info.LoadAvg, f.yellow))

	// CPU information
	var cpuInfoValue string
	if info.CPUHz > 0 {
		cpuInfoValue = fmt.Sprintf("%s in use of %d cores/%d threads at %.2fGHz",
			info.CPUUsage, info.CPUCores, info.CPUThreads, info.CPUHz)
	} else {
		cpuInfoValue = fmt.Sprintf("%s in use of %d cores/%d threads",
			info.CPUUsage, info.CPUCores, info.CPUThreads)
	}
	messages = append(messages, f.formatLineWithColor("CPU info", cpuInfoValue, f.magenta))

	// Memory information
	messages = append(messages, f.formatLineWithColor("Memory in use",
		fmt.Sprintf("%s of %s", info.MemUsed, info.MemTotal), f.cyan))
	messages = append(messages, f.formatLineWithColor("Swap memory in use",
		fmt.Sprintf("%s of %s", info.SwapUsed, info.SwapTotal), f.cyan))

	// Disk information
	messages = append(messages, f.formatLineWithColor("Root disk usage",
		fmt.Sprintf("%s of %s", info.DiskUse, info.DiskSize), f.yellow))
	messages = append(messages, f.formatLineWithColor("Disk pool size", info.DiskPoolSize, f.yellow))
	messages = append(messages, f.formatLineWithColor("Disk pool used", info.DiskPoolUsed, f.yellow))

	// Process and user information
	processInfo := fmt.Sprintf("%s running %s, total of %s running on %s",
		info.User, info.ProcessesUser, info.ProcessesAll, info.Host)
	messages = append(messages, f.formatLineWithColor("System processes", processInfo, f.blue))

	userInfo := fmt.Sprintf("%d user(s) currently logged in", info.Users)
	messages = append(messages, f.formatLineWithColor("Users", userInfo, f.blue))

	sessionInfo := fmt.Sprintf("%s current active session(s)", info.ActiveSessions)
	messages = append(messages, f.formatLineWithColor("Sessions", sessionInfo, f.blue))

	messages = append(messages, f.formatLineWithColor("Last system login", info.LastLog, f.cyan))

	// Fan information (for specific hosts)
	if info.User == "root" && info.Host == "saturn" {
		messages = append(messages, lineBorder)
		messages = append(messages, f.formatLineWithColor("Fans 1 & 2",
			fmt.Sprintf("%s/rpm & %s/rpm", info.Fan1, info.Fan2), f.blue))
		messages = append(messages, f.formatLineWithColor("Fans 3 & 4",
			fmt.Sprintf("%s/rpm & %s/rpm", info.Fan3, info.Fan4), f.blue))
		messages = append(messages, f.formatLineWithColor("Fans 5 & 6",
			fmt.Sprintf("%s/rpm & %s/rpm", info.Fan5, info.Fan6), f.blue))
	}

	// Network information
	messages = append(messages, lineBorder)
	messages = append(messages, f.formatLineWithColor("IPv4 address", info.IPv4, f.cyan))
	messages = append(messages, f.formatLineWithColor("IPv6 address", info.IPv6, f.cyan))
	messages = append(messages, lineBorder)

	// VPN information (for specific hosts)
	if info.User == "orion" && info.Host == "titan" {
		messages = append(messages, f.formatLineWithColor("VPN address", info.NordAddr, f.yellow))
		messages = append(messages, f.formatLineWithColor("Transmission address", info.TransAddr, f.yellow))
		messages = append(messages, lineBorder)
		messages = append(messages, f.formatLineWithColor("Transmission status", info.VPNCheck, f.green))
		messages = append(messages, f.formatLineWithColor("Transkick status", info.TranskickStatus, f.green))
		messages = append(messages, lineBorder)
	}

	return messages
}

// formatLine formats a single information line with consistent spacing.
func (f *Formatter) formatLine(label, value string) string {
	return f.formatLineWithColor(label, value, f.magenta)
}

// formatLineWithColor formats a single information line with consistent spacing and custom color.
func (f *Formatter) formatLineWithColor(label, value string, colorFunc func(a ...interface{}) string) string {
	const labelWidth = 25 // Fixed width for label column
	padding := labelWidth - len(label)
	if padding < 1 {
		padding = 1
	}
	spaces := ""
	for i := 0; i < padding; i++ {
		spaces += " "
	}
	return fmt.Sprintf("%s[*]%s %s%s: %s%s",
		f.green(""), f.reset(), label, spaces, colorFunc(""), value)
}
