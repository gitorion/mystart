# MyStart

A cross-platform system information display tool written in Go, providing comprehensive system metrics and status information.

## Supported Platforms

- ✅ **Linux** (Ubuntu, Debian, Fedora, Arch, etc.)
- ✅ **macOS** (Intel and Apple Silicon)

## Features

- System identification (OS, distribution, hostname, user)
- CPU information (cores, threads, frequency, usage, load average)
- Memory and swap usage
- Disk usage and pool statistics
- Network configuration (IPv4/IPv6)
- User session information
- System uptime breakdown
- Host-specific features (Linux only):
  - Fan speed monitoring (via liquidctl for specific hosts)
  - VPN and Transmission status monitoring

## Installation

### Prerequisites

- Go 1.21 or higher
- Linux or macOS operating system

### Quick Build

```bash
# Build for current platform
make build

# Or manually
go build -o mystart ./cmd/mystart
```

### Cross-Platform Builds

```bash
# Build for Linux
make build-linux

# Build for macOS
make build-darwin

# Build for all platforms
make build-all
```

### Install

```bash
make install

# Or manually
go install ./cmd/mystart
```

## Usage

Simply run the binary:

```bash
# Run directly
./bin/mystart

# Or if installed
mystart

# Using make
make run
```

### Example Output

```
=======================================================================
User: orion	Host: server01 👉 VM - Dockerised Plex Media Server
=======================================================================
[*] Login details		: 13 Feb 22:03 still from ttys000
[*] System details		: ubuntu | Ubuntu 22.04 LTS
[*] System uptime		: 14 days 14 hours 8 minutes 30 seconds
[*] System load		: 2.12 2.78 2.81
[*] CPU info			: 27.20 % in use of 8cores/16threads at 3.60GHz
[*] Memory in use		: 1.60G of 8.00G
[*] Swap memory in use		: 0.27G of 4.00G
[*] Root disk usage		: 45% of 100G
[*] Disk pool size		: 2.5TB
[*] Disk pool used		: 1.8TB
[*] System processes		: orion running 149, total of 241 running on server01
[*] Users			: 2 user(s) currently logged in
[*] Sessions		: 3 current active session(s)
[*] Last system login		: orion on pts/0, 13 Feb 21:25 from 192.168.1.100
=======================================================================
[*] IPv4 address		: 192.168.1.50
[*] IPv6 address		: fe80::1234:5678:90ab:cdef
=======================================================================
```

## Project Structure

```
mystart/
├── cmd/
│   └── mystart/           # Application entry point
│       └── main.go
├── internal/
│   ├── collector/         # System information collection
│   │   ├── collector.go   # Main collector coordinator
│   │   ├── types.go       # Shared types and utilities
│   │   ├── *_linux.go     # Linux-specific implementations
│   │   ├── *_darwin.go    # macOS-specific implementations
│   │   └── host.go        # Host-specific features
│   ├── display/           # Display formatting
│   │   └── display.go
│   └── config/            # Configuration and constants
│       └── config.go
├── Makefile               # Build automation
├── go.mod
├── go.sum
└── README.md
```

### Platform-Specific Code

The project uses Go build tags to provide platform-specific implementations:

- **Linux**: Uses `/proc` filesystem and `ip` command
- **macOS**: Uses `sysctl`, `sw_vers`, and BSD utilities

## Design Principles

This project follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md):

- Error handling with custom error types
- Proper naming conventions (package names, functions, variables)
- Context-based timeout management
- Clear separation of concerns
- Minimal use of global state
- Proper resource cleanup with defer
- Pre-allocated slices where capacity is known
- Interface-based design where appropriate

## Dependencies

- [github.com/fatih/color](https://github.com/fatih/color) - Terminal color output

## License

See LICENSE file for details.
