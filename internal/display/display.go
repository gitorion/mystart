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
	reset   func(a ...interface{}) string
}

// NewFormatter creates a new display formatter.
func NewFormatter() *Formatter {
	return &Formatter{
		green:   color.New(color.FgGreen).SprintFunc(),
		magenta: color.New(color.FgMagenta).SprintFunc(),
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
	messages = append(messages, f.formatLine("Login details", info.ThisLog))

	// System information
	messages = append(messages, fmt.Sprintf("%s[*]%s System details\t\t:%s %s %s| %s %s",
		f.green(""), f.reset(), f.green(""), info.Distro, f.reset(), f.magenta(""), info.Uname))
	messages = append(messages, fmt.Sprintf("%s[*]%s System uptime\t\t:%s %d days %d hours %d minutes %d seconds",
		f.green(""), f.reset(), f.magenta(""), info.UptimeDays, info.UptimeHours,
		info.UptimeMinutes, info.UptimeSeconds))
	messages = append(messages, f.formatLine("System load", info.LoadAvg))

	// CPU information
	var cpuInfo string
	if info.CPUHz > 0 {
		cpuInfo = fmt.Sprintf("%s[*]%s CPU info\t\t\t:%s %s in use of %dcores/%dthreads at %.2fGHz",
			f.green(""), f.reset(), f.magenta(""), info.CPUUsage, info.CPUCores,
			info.CPUThreads, info.CPUHz)
	} else {
		cpuInfo = fmt.Sprintf("%s[*]%s CPU info\t\t\t:%s %s in use of %dcores/%dthreads",
			f.green(""), f.reset(), f.magenta(""), info.CPUUsage, info.CPUCores,
			info.CPUThreads)
	}
	messages = append(messages, cpuInfo)

	// Memory information
	messages = append(messages, f.formatLine("Memory in use",
		fmt.Sprintf("%s of %s", info.MemUsed, info.MemTotal)))
	messages = append(messages, f.formatLine("Swap memory in use",
		fmt.Sprintf("%s of %s", info.SwapUsed, info.SwapTotal)))

	// Disk information
	messages = append(messages, f.formatLine("Root disk usage",
		fmt.Sprintf("%s of %s", info.DiskUse, info.DiskSize)))
	messages = append(messages, f.formatLine("Disk pool size", info.DiskPoolSize))
	messages = append(messages, f.formatLine("Disk pool used", info.DiskPoolUsed))

	// Process and user information
	messages = append(messages, fmt.Sprintf("%s[*]%s System processes\t\t:%s %s running %s, total of %s running on %s",
		f.green(""), f.reset(), f.magenta(""), info.User, info.ProcessesUser,
		info.ProcessesAll, info.Host))
	messages = append(messages, fmt.Sprintf("%s[*]%s Users\t\t\t:%s %d user(s) currently logged in",
		f.green(""), f.reset(), f.magenta(""), info.Users))
	messages = append(messages, f.formatLine("Sessions",
		fmt.Sprintf("%s current active session(s)", info.ActiveSessions)))
	messages = append(messages, f.formatLine("Last system login", info.LastLog))

	// Fan information (for specific hosts)
	if info.User == "root" && info.Host == "saturn" {
		messages = append(messages, lineBorder)
		messages = append(messages, f.formatLine("Fans 1 & 2",
			fmt.Sprintf("%s/rpm & %s/rpm", info.Fan1, info.Fan2)))
		messages = append(messages, f.formatLine("Fans 3 & 4",
			fmt.Sprintf("%s/rpm & %s/rpm", info.Fan3, info.Fan4)))
		messages = append(messages, f.formatLine("Fans 5 & 6",
			fmt.Sprintf("%s/rpm & %s/rpm", info.Fan5, info.Fan6)))
	}

	// Network information
	messages = append(messages, lineBorder)
	messages = append(messages, f.formatLine("IPv4 address", info.IPv4))
	messages = append(messages, f.formatLine("IPv6 address", info.IPv6))
	messages = append(messages, lineBorder)

	// VPN information (for specific hosts)
	if info.User == "orion" && info.Host == "titan" {
		messages = append(messages, f.formatLine("VPN address", info.NordAddr))
		messages = append(messages, f.formatLine("Transmission address", info.TransAddr))
		messages = append(messages, lineBorder)
		messages = append(messages, f.formatLine("Transmission status", info.VPNCheck))
		messages = append(messages, f.formatLine("Transkick status", info.TranskickStatus))
		messages = append(messages, lineBorder)
	}

	return messages
}

// formatLine formats a single information line with consistent spacing.
func (f *Formatter) formatLine(label, value string) string {
	return fmt.Sprintf("%s[*]%s %s\t\t:%s %s",
		f.green(""), f.reset(), label, f.magenta(""), value)
}
