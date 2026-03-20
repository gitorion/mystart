package main

import (
	"fmt"
	"os"

	"github.com/orion/mystart/internal/collector"
	"github.com/orion/mystart/internal/display"
)

func main() {
	c := collector.New()
	info, err := c.GatherAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	display.Render(info)
}
