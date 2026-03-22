package collector

// DiskMount represents a single mounted filesystem.
type DiskMount struct {
	Path        string
	UsedGB      float64
	TotalGB     float64
	UsedPercent float64
}

// SystemInfo holds all collected system metrics.
type SystemInfo struct {
	// System identification
	Hostname string
	User     string
	OS       string
	Kernel   string
	Shell    string

	// Uptime
	UptimeDays    int
	UptimeHours   int
	UptimeMinutes int
	UptimeSeconds int

	// CPU
	CPUModel   string
	CPUCores   int
	CPUThreads int
	CPUHz      float64 // GHz; 0 if unknown
	CPUUsage   float64 // percentage 0–100

	// Load averages
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64

	// Memory (GB)
	MemTotalGB  float64
	MemUsedGB   float64
	SwapTotalGB float64
	SwapUsedGB  float64

	// Disk – one entry per visible mount point
	DiskMounts []DiskMount

	// GPU (empty if no GPU detected)
	GPUModel   string
	GPUVram    string // human-readable, e.g. "8 GB"
	GPUUsage   string // percentage string, e.g. "34%", empty if unavailable
	GPUTemp    string // e.g. "62°C", empty if unavailable
	GPUDriver  string // driver/metal version, empty if unavailable

	// Network
	IPv4 string
	IPv6 string

	// Sessions / processes
	ProcessesUser  int
	ProcessesTotal int
	ActiveSessions int
	UsersLoggedIn  int
	LastLogin      string
}
