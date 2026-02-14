package collector

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/orion/mystart/internal/config"
)

// GatherHostSpecificInfo collects information specific to certain hosts.
func (s *SystemInfo) GatherHostSpecificInfo(ctx context.Context) {
	// Fan information for saturn host as root
	if s.User == "root" && s.Host == "saturn" {
		s.gatherFanInfo(ctx)
	} else {
		s.setAllFansUnavailable()
	}

	// VPN/Transmission info for titan host as orion user
	if s.User == "orion" && s.Host == "titan" {
		s.gatherVPNInfo(ctx)
	}
}

// gatherFanInfo collects fan speed information using liquidctl.
func (s *SystemInfo) gatherFanInfo(ctx context.Context) {
	// Create context with timeout for liquidctl command
	cmdCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "liquidctl", "status")
	if err := cmd.Run(); err != nil {
		s.setAllFansUnavailable()
		return
	}

	// Gather individual fan speeds
	s.Fan1 = s.getFanSpeed(ctx, 1)
	s.Fan2 = s.getFanSpeed(ctx, 2)
	s.Fan3 = s.getFanSpeed(ctx, 3)
	s.Fan4 = s.getFanSpeed(ctx, 4)
	s.Fan5 = s.getFanSpeed(ctx, 5)
	s.Fan6 = s.getFanSpeed(ctx, 6)
}

// getFanSpeed retrieves the speed for a specific fan number.
func (s *SystemInfo) getFanSpeed(ctx context.Context, fanNum int) string {
	cmd := fmt.Sprintf(`liquidctl status | grep "Fan %d" | awk '{print $5}'`, fanNum)
	speed := execCommandSafe(ctx, cmd)
	if speed == "" {
		return config.CrossMark
	}
	return speed
}

// setAllFansUnavailable sets all fan readings to unavailable.
func (s *SystemInfo) setAllFansUnavailable() {
	s.Fan1 = config.CrossMark
	s.Fan2 = config.CrossMark
	s.Fan3 = config.CrossMark
	s.Fan4 = config.CrossMark
	s.Fan5 = config.CrossMark
	s.Fan6 = config.CrossMark
}

// gatherVPNInfo collects VPN and Transmission information.
func (s *SystemInfo) gatherVPNInfo(ctx context.Context) {
	s.checkTranskickService(ctx)
	s.checkVPNAddress(ctx)
	s.checkTransmissionAddress(ctx)
	s.verifyVPNProtection()
}

// checkTranskickService checks the status of the transkick service.
func (s *SystemInfo) checkTranskickService(ctx context.Context) {
	status := execCommandSafe(ctx, `systemctl status transkick.service | grep Active | awk '{ print $2, $3 }'`)

	switch status {
	case "active (running)":
		s.TranskickStatus = fmt.Sprintf("%s %s", status, config.ThumbUp)
	case "Unit transkick.service could not be found.":
		s.TranskickStatus = config.CrossMark
	case "":
		s.TranskickStatus = config.CrossMark
	default:
		s.TranskickStatus = fmt.Sprintf("%s %s", status, config.StopEmoji)
	}
}

// checkVPNAddress retrieves the VPN container's external IP address.
func (s *SystemInfo) checkVPNAddress(ctx context.Context) {
	nordCtx, nordCancel := context.WithTimeout(ctx, 3*time.Second)
	defer nordCancel()

	cmd := exec.CommandContext(nordCtx, "docker", "exec", "nord", "curl", "ifconfig.io")
	output, err := cmd.Output()
	if err != nil {
		s.NordAddr = config.CrossMark
	} else {
		s.NordAddr = strings.TrimSpace(string(output))
	}
}

// checkTransmissionAddress retrieves the Transmission container's external IP address.
func (s *SystemInfo) checkTransmissionAddress(ctx context.Context) {
	transCtx, transCancel := context.WithTimeout(ctx, 3*time.Second)
	defer transCancel()

	cmd := exec.CommandContext(transCtx, "docker", "exec", "transmission", "curl", "ifconfig.io")
	output, err := cmd.Output()
	if err != nil {
		s.TransAddr = config.CrossMark
	} else {
		s.TransAddr = strings.TrimSpace(string(output))
	}
}

// verifyVPNProtection checks if Transmission is protected by VPN.
func (s *SystemInfo) verifyVPNProtection() {
	if s.TranskickStatus == config.CrossMark {
		s.VPNCheck = config.CrossMark
		return
	}

	if s.NordAddr == s.TransAddr && s.NordAddr != config.CrossMark {
		s.VPNCheck = fmt.Sprintf("Protected %s", config.ThumbUp)
	} else {
		s.VPNCheck = fmt.Sprintf("Unprotected %s", config.StopEmoji)
	}
}
