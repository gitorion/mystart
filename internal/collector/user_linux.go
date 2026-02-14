//go:build linux
// +build linux

package collector

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// GatherUserInfo collects user session and process information on Linux.
func (s *SystemInfo) GatherUserInfo(ctx context.Context) error {
	if err := s.gatherLoginInfo(ctx); err != nil {
		return err
	}

	if err := s.gatherProcessInfo(ctx); err != nil {
		return err
	}

	if err := s.gatherSessionInfo(ctx); err != nil {
		return err
	}

	return nil
}

// gatherLoginInfo retrieves last login information.
func (s *SystemInfo) gatherLoginInfo(ctx context.Context) error {
	lastLog, err := execCommand(ctx, `last | head -n 2 | tail -1 | awk '{print $1 " on " $2 ", " $4 " " $5 " " $6 " " $7 " from " $3}'`)
	if err != nil {
		return fmt.Errorf("getting last login: %w", err)
	}
	s.LastLog = lastLog

	// Check which lastlog command is available
	_, err = exec.LookPath("lastlog")
	var thisLogCmd string
	if err == nil {
		thisLogCmd = `lastlog -u $USER | tail -n 1 | awk '{print $4 " " $5 " " $6 " " $7 " from " $3}'`
	} else {
		thisLogCmd = `lastlog2 -u $USER | tail -n 1 | awk '{print $4 " " $5 " " $6 " " $7 " from " $3}'`
	}

	thisLog, err := execCommand(ctx, thisLogCmd)
	if err != nil {
		return fmt.Errorf("getting this login: %w", err)
	}
	s.ThisLog = thisLog

	return nil
}

// gatherProcessInfo retrieves process count information.
func (s *SystemInfo) gatherProcessInfo(ctx context.Context) error {
	processesUser, err := execCommand(ctx, "ps -aux | grep -i $USER | wc -l")
	if err != nil {
		return fmt.Errorf("getting user processes: %w", err)
	}
	s.ProcessesUser = processesUser

	processesAll, err := execCommand(ctx, "ps -aux | wc -l")
	if err != nil {
		return fmt.Errorf("getting all processes: %w", err)
	}
	s.ProcessesAll = processesAll

	return nil
}

// gatherSessionInfo retrieves active session information.
func (s *SystemInfo) gatherSessionInfo(ctx context.Context) error {
	activeSessions, err := execCommand(ctx, "w | awk '{print $1}'| sed 1,2d | wc -l")
	if err != nil {
		return fmt.Errorf("getting active sessions: %w", err)
	}
	s.ActiveSessions = activeSessions

	// Count unique users
	if err := s.countUniqueUsers(ctx); err != nil {
		return err
	}

	return nil
}

// countUniqueUsers counts the number of unique logged-in users.
func (s *SystemInfo) countUniqueUsers(ctx context.Context) error {
	output, err := execCommand(ctx, "w | awk '{print $1}'| sed 1,2d")
	if err != nil {
		return fmt.Errorf("getting user list: %w", err)
	}

	// Use map as set to track unique users
	userSet := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		user := scanner.Text()
		if user != "" {
			userSet[user] = struct{}{}
		}
	}

	s.Users = len(userSet)
	return nil
}
