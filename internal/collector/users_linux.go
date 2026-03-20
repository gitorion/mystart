//go:build linux
// +build linux

package collector

import (
	"bufio"
	"context"
	"strconv"
	"strings"
)

func gatherUsers(ctx context.Context, info *SystemInfo) {
	// Active sessions and unique logged-in user count via `w`
	wOut := execCommandSafe(ctx, "w -h")
	if wOut != "" {
		trimmed := strings.TrimSpace(wOut)
		if trimmed != "" {
			lines := strings.Split(trimmed, "\n")
			info.ActiveSessions = len(lines)
			users := make(map[string]bool)
			for _, line := range lines {
				f := strings.Fields(line)
				if len(f) > 0 {
					users[f[0]] = true
				}
			}
			info.UsersLoggedIn = len(users)
		}
	}

	// Process counts
	psAll := execCommandSafe(ctx, "ps -e --no-headers | wc -l")
	if n, err := strconv.Atoi(strings.TrimSpace(psAll)); err == nil {
		info.ProcessesTotal = n
	}
	if info.User != "" {
		psUser := execCommandSafe(ctx, "ps -U "+info.User+" --no-headers | wc -l")
		if n, err := strconv.Atoi(strings.TrimSpace(psUser)); err == nil {
			info.ProcessesUser = n
		}
	}

	// Last login: skip the first entry (current session), use the second
	lastOut := execCommandSafe(ctx, "last -100 | grep -v '^wtmp' | grep -v '^$' | grep -v '^reboot'")
	if lastOut != "" {
		sc := bufio.NewScanner(strings.NewReader(lastOut))
		var count int
		for sc.Scan() {
			line := sc.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}
			count++
			if count == 2 {
				info.LastLogin = parseLastLine(line)
				break
			}
		}
	}
}

// parseLastLine extracts a human-readable string from a `last` output line.
func parseLastLine(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return strings.TrimSpace(line)
	}
	user := fields[0]
	var from, when string
	days := map[string]bool{"Mon": true, "Tue": true, "Wed": true, "Thu": true,
		"Fri": true, "Sat": true, "Sun": true}
	if !days[fields[2]] && len(fields) >= 6 {
		from = fields[2]
		when = strings.Join(fields[3:7], " ")
	} else if len(fields) >= 5 {
		when = strings.Join(fields[2:6], " ")
	}
	if from != "" {
		return user + " from " + from + "  " + when
	}
	return user + "  " + when
}
