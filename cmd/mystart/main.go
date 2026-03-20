package main

import (
	"fmt"
	"os"

	"github.com/gitorion/mystart/internal/collector"
	"github.com/gitorion/mystart/internal/display"
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
