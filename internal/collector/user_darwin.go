//go:build darwin
// +build darwin

package collector

import (
	"bufio"
	"context"
	"fmt"
	"strings"
)

// GatherUserInfo collects user session and process information on macOS.
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
	// Use 'last' command which is available on macOS
	lastLog := execCommandSafe(ctx, `last | head -n 2 | tail -1 | awk '{print $1 " on " $2 ", " $4 " " $5 " " $6 " " $7}'`)
	if lastLog == "" {
		lastLog = "No login history available"
	}
	s.LastLog = lastLog

	// Get current login info
	thisLog := execCommandSafe(ctx, `last $USER | head -n 1 | awk '{print $4 " " $5 " " $6 " " $7 " from " $3}'`)
	if thisLog == "" {
		thisLog = "Currently logged in"
	}
	s.ThisLog = thisLog

	return nil
}

// gatherProcessInfo retrieves process count information.
func (s *SystemInfo) gatherProcessInfo(ctx context.Context) error {
	processesUser, err := execCommand(ctx, "ps aux | grep -i $USER | wc -l")
	if err != nil {
		return fmt.Errorf("getting user processes: %w", err)
	}
	s.ProcessesUser = strings.TrimSpace(processesUser)

	processesAll, err := execCommand(ctx, "ps aux | wc -l")
	if err != nil {
		return fmt.Errorf("getting all processes: %w", err)
	}
	s.ProcessesAll = strings.TrimSpace(processesAll)

	return nil
}

// gatherSessionInfo retrieves active session information.
func (s *SystemInfo) gatherSessionInfo(ctx context.Context) error {
	activeSessions, err := execCommand(ctx, "w -h | wc -l")
	if err != nil {
		return fmt.Errorf("getting active sessions: %w", err)
	}
	s.ActiveSessions = strings.TrimSpace(activeSessions)

	// Count unique users
	if err := s.countUniqueUsers(ctx); err != nil {
		return err
	}

	return nil
}

// countUniqueUsers counts the number of unique logged-in users.
func (s *SystemInfo) countUniqueUsers(ctx context.Context) error {
	output := execCommandSafe(ctx, "w -h | awk '{print $1}'")
	if output == "" {
		s.Users = 0
		return nil
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
