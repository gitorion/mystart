//go:build darwin
// +build darwin

package collector

import (
	"bufio"
	"context"
	"strconv"
	"strings"
)

func gatherUsers(ctx context.Context, info *SystemInfo) {
	// Active sessions and unique logged-in user count via `w -h`
	wOut := execCommandSafe(ctx, "w -h")
	if wOut != "" {
		lines := strings.Split(strings.TrimSpace(wOut), "\n")
		info.ActiveSessions = len(lines)
		users := make(map[string]bool)
		sc := bufio.NewScanner(strings.NewReader(wOut))
		for sc.Scan() {
			f := strings.Fields(sc.Text())
			if len(f) > 0 {
				users[f[0]] = true
			}
		}
		info.UsersLoggedIn = len(users)
	}

	// Process counts
	psAll := execCommandSafe(ctx, "ps aux | wc -l")
	if n, err := strconv.Atoi(strings.TrimSpace(psAll)); err == nil && n > 1 {
		info.ProcessesTotal = n - 1 // subtract header line
	}
	psUser := execCommandSafe(ctx, "ps -U "+info.User+" | wc -l")
	if n, err := strconv.Atoi(strings.TrimSpace(psUser)); err == nil && n > 1 {
		info.ProcessesUser = n - 1
	}

	// Last login: skip the first entry (current session), use the second
	lastOut := execCommandSafe(ctx, "last -100 | grep -v '^wtmp' | grep -v '^$'")
	if lastOut != "" {
		sc := bufio.NewScanner(strings.NewReader(lastOut))
		var count int
		for sc.Scan() {
			line := sc.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}
			count++
			if count == 2 { // second entry = previous login
				info.LastLogin = parseLastLine(line)
				break
			}
		}
	}
}

// parseLastLine extracts a human-readable string from a `last` output line.
func parseLastLine(line string) string {
	fields := strings.Fields(line)
	// Format: user tty from Mon Day Time Year
	if len(fields) < 4 {
		return strings.TrimSpace(line)
	}
	user := fields[0]
	var from, when string
	// Check if field[2] looks like a hostname/IP (not a day-of-week)
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
