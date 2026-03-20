package collector

import (
	"context"
	"os/exec"
	"strings"
	"sync"

	"github.com/orion/mystart/internal/config"
)

// Collector orchestrates system information gathering.
type Collector struct{}

// New returns a new Collector.
func New() *Collector {
	return &Collector{}
}

// GatherAll collects all system metrics. System identification is gathered
// first; all other categories run concurrently.
func (c *Collector) GatherAll() (*SystemInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.OverallTimeout)
	defer cancel()

	info := &SystemInfo{}

	// System must complete first; other gatherers may need hostname / user.
	gatherSystem(ctx, info)

	var wg sync.WaitGroup
	wg.Add(6)
	go func() { defer wg.Done(); gatherCPU(ctx, info) }()
	go func() { defer wg.Done(); gatherMemory(ctx, info) }()
	go func() { defer wg.Done(); gatherDisk(ctx, info) }()
	go func() { defer wg.Done(); gatherNetwork(ctx, info) }()
	go func() { defer wg.Done(); gatherUptime(ctx, info) }()
	go func() { defer wg.Done(); gatherUsers(ctx, info) }()
	wg.Wait()

	return info, nil
}

// execCommand runs a shell command via sh -c and returns trimmed stdout.
func execCommand(ctx context.Context, cmd string) (string, error) {
	c := exec.CommandContext(ctx, "sh", "-c", cmd)
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// execCommandSafe runs a command and returns "" on any error.
func execCommandSafe(ctx context.Context, cmd string) string {
	out, _ := execCommand(ctx, cmd)
	return out
}
