package config

import "time"

const (
	// Display dimensions
	BoxWidth   = 76
	BarWidth   = 20
	LabelWidth = 18

	// Timeout settings
	CommandTimeout = 5 * time.Second
	OverallTimeout = 30 * time.Second
)
