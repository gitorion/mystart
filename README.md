# mystart

A cross-platform terminal system information tool written in Go. Displays comprehensive system metrics in a clean, colour-coded dashboard.

Works on any Linux distribution or macOS system for any user — no configuration required.

## Supported Platforms

| Platform | Architecture |
|----------|-------------|
| macOS | Intel (amd64), Apple Silicon (arm64) |
| Linux | x86_64 (amd64) |

## Features

- **System** — hostname, user, OS name & version, kernel, shell, uptime
- **Processor** — model name, physical/logical cores, frequency, usage %, load averages
- **Memory** — RAM and swap with usage bars
- **Storage** — all mounted filesystems with usage bars (auto GB / TB)
- **Network** — primary IPv4 and IPv6 addresses
- **Sessions** — logged-in users, active sessions, process counts, last login

Progress bars change colour automatically: green below 60 %, yellow 60–80 %, red above 80 %.

## Requirements

- **Go 1.21+** — [install from golang.org](https://golang.org/dl/)
- Standard system utilities only — no additional packages needed

---

## Quick Start

The fastest way to get `mystart` running automatically every time you open a terminal.

### macOS

**1. Install Go** (if not already installed) — easiest via [Homebrew](https://brew.sh):
```bash
brew install go
```

**2. Install mystart:**
```bash
go install github.com/orion/mystart/cmd/mystart@latest
```

**3. Add to shell startup** — zsh is the default shell on macOS:
```bash
echo 'mystart' >> ~/.zshrc
```

Open a new terminal and mystart will run automatically.

---

### Linux

**1. Install Go** (if not already installed):
```bash
# Ubuntu / Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go

# Or download directly from golang.org/dl
```

**2. Install mystart:**
```bash
go install github.com/orion/mystart/cmd/mystart@latest
```

**3. Add to shell startup** — bash is the common default on Linux:
```bash
echo 'mystart' >> ~/.bashrc
```

Open a new terminal and mystart will run automatically.

---

> **`mystart` not found after install?** Go's bin directory may not be in your `$PATH`. Add this to your rc file before the `mystart` line, then open a new terminal:
> ```bash
> export PATH="$PATH:$HOME/go/bin"
> ```

---

## Updating

To update to the latest version, run the same install command again:

```bash
go install github.com/orion/mystart/cmd/mystart@latest
```

This overwrites the existing binary automatically.

---

## Other Installation Options

### Build from source

```bash
git clone https://github.com/orion/mystart.git
cd mystart
make build
./bin/mystart
```

### Run without installing

```bash
git clone https://github.com/orion/mystart.git
cd mystart
go run ./cmd/mystart/
```

## Build targets

```bash
make build          # build for current platform  →  bin/mystart
make build-linux    # cross-compile for Linux     →  bin/mystart-linux
make build-darwin   # cross-compile for macOS     →  bin/mystart-darwin-arm64 / amd64
make build-all      # all of the above
make install        # install to $GOPATH/bin
make run            # build and run immediately
make clean          # remove bin/
```

## Example Output

```
╭──────────────────────────────────────────────────────────────────────────╮
│                           ◈  SYSTEM STATUS  ◈                            │
│                           server01  ·  alice                             │
╠══════════════════════════════════════════════════════════════════════════╣
│  ◆ SYSTEM                                                                │
│    Hostname            server01                                          │
│    User                alice                                             │
│    OS                  Ubuntu 22.04.3 LTS                                │
│    Kernel              Linux 5.15.0-91-generic x86_64                   │
│    Shell               bash                                              │
│    Uptime              14 days, 8 hours, 30 minutes, 12 seconds          │
╠══════════════════════════════════════════════════════════════════════════╣
│  ◆ PROCESSOR                                                             │
│    Model               Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz        │
│    Cores / Threads     8 physical · 16 logical · 3.80 GHz               │
│    Usage               ████░░░░░░░░░░░░░░░░  21.3%                       │
│    Load Average        2.12  2.78  2.81   (1m / 5m / 15m)                │
╠══════════════════════════════════════════════════════════════════════════╣
│  ◆ MEMORY                                                                │
│    RAM                 ████░░░░░░░░░░░░░░░░  1.6 / 8.0 GB  (20.0%)       │
│    Swap                █░░░░░░░░░░░░░░░░░░░  0.3 / 4.0 GB  (6.8%)        │
╠══════════════════════════════════════════════════════════════════════════╣
│  ◆ STORAGE                                                               │
│    /                   █████████░░░░░░░░░░░  45.0 / 100.0 GB  (45.0%)   │
│    /home               ██████████████░░░░░░  420.0 / 600.0 GB  (70.0%)  │
╠══════════════════════════════════════════════════════════════════════════╣
│  ◆ NETWORK                                                               │
│    IPv4                192.168.1.50                                      │
│    IPv6                fd00::1:2:3:4                                     │
╠══════════════════════════════════════════════════════════════════════════╣
│  ◆ SESSIONS                                                              │
│    Users               2 logged in · 3 active sessions                   │
│    Processes           149 user · 241 total                              │
│    Last Login          alice from 192.168.1.100  Mon 13 Feb 21:25        │
╰──────────────────────────────────────────────────────────────────────────╯
```

## Project Structure

```
mystart/
├── cmd/mystart/           # Binary entry point
│   └── main.go
├── internal/
│   ├── collector/         # System metric collection
│   │   ├── types.go       # SystemInfo and DiskMount types
│   │   ├── collector.go   # Orchestration and exec helpers
│   │   ├── system_{darwin,linux}.go
│   │   ├── cpu_{darwin,linux}.go
│   │   ├── memory_{darwin,linux}.go
│   │   ├── disk_{darwin,linux}.go
│   │   ├── network_{darwin,linux}.go
│   │   ├── uptime_{darwin,linux}.go
│   │   └── users_{darwin,linux}.go
│   ├── display/
│   │   └── display.go     # Coloured box UI
│   └── config/
│       └── config.go      # Layout and timeout constants
├── Makefile
├── go.mod
└── go.sum
```

Platform-specific implementations use Go build tags (`//go:build darwin`, `//go:build linux`) so only the correct code is compiled for each target.

## Dependencies

- [`github.com/fatih/color`](https://github.com/fatih/color) — terminal colour output (only external dependency)

## License

See LICENSE file for details.
