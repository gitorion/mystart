# mystart

A cross-platform terminal system information tool written in Go. Displays comprehensive system metrics in a clean, colour-coded dashboard.

Works on any Linux distribution or macOS system for any user — no configuration required.

## Example Output

```
╭──────────────────────────────────────────────────────────────────────────╮
│                                                                          │
│                         ◈  SERVER01 STATUS  ◈                            │
│                           server01  ·  alice                             │
│                                                                          │
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
│  ◆ GPU                                                                   │
│    Model               NVIDIA GeForce RTX 3080                           │
│    VRAM                3456 / 10240 MB                                   │
│    Usage               34%                                               │
│    Temperature         62°C                                              │
│    Driver              535.129.03                                         │
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

Each section header has its own accent colour. Progress bars change colour automatically: green below 60 %, yellow 60–80 %, red above 80 %. The GPU section only appears when a GPU is detected.

## Supported Platforms

| Platform | Architecture |
|----------|-------------|
| macOS | Intel (amd64), Apple Silicon (arm64) |
| Linux | x86_64 (amd64) |

## Features

- **System** — hostname, user, OS, kernel, shell, uptime
- **Processor** — model, cores/threads, frequency, usage %, load averages
- **Memory** — RAM and swap with colour-coded usage bars
- **Storage** — all mounted filesystems with usage bars (auto GB / TB)
- **GPU** — model, VRAM, usage, temperature, driver (when available)
- **Network** — primary IPv4 and IPv6 addresses
- **Sessions** — logged-in users, active sessions, process counts, last login

## Requirements

- **Go 1.21+** — [install from golang.org](https://golang.org/dl/)
- Standard system utilities only — no additional packages needed

---

## Quick Start

Install `mystart` and have it run automatically every time you open a terminal.

### macOS

**1. Install Go** (if not already installed) — easiest via [Homebrew](https://brew.sh):
```bash
brew install go
```

**2. Install mystart:**
```bash
go install github.com/gitorion/mystart/cmd/mystart@latest
```

**3. Add to shell startup:**

zsh (default on macOS):
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.zshrc
echo 'mystart' >> ~/.zshrc
source ~/.zshrc
```

bash:
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bash_profile
echo 'mystart' >> ~/.bash_profile
source ~/.bash_profile
```

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
go install github.com/gitorion/mystart/cmd/mystart@latest
```

**3. Add to shell startup:**

bash (default on most Linux distributions):
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
echo 'mystart' >> ~/.bashrc
source ~/.bashrc
```

zsh:
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.zshrc
echo 'mystart' >> ~/.zshrc
source ~/.zshrc
```

---

## Updating

To update to the latest version, run the same install command again:

```bash
go install github.com/gitorion/mystart/cmd/mystart@latest
```

This overwrites the existing binary automatically.

---

## Other Installation Options

### Build from source

```bash
git clone https://github.com/gitorion/mystart.git
cd mystart
make build
./bin/mystart
```

### Run without installing

```bash
git clone https://github.com/gitorion/mystart.git
cd mystart
go run ./cmd/mystart/
```

## Project Structure

```
mystart/
├── cmd/mystart/
│   └── main.go              # Entry point
├── internal/
│   ├── collector/
│   │   ├── types.go          # SystemInfo and DiskMount types
│   │   ├── collector.go      # Orchestration and exec helpers
│   │   ├── system_{darwin,linux}.go
│   │   ├── cpu_{darwin,linux}.go
│   │   ├── memory_{darwin,linux}.go
│   │   ├── disk_{darwin,linux}.go
│   │   ├── gpu_{darwin,linux}.go
│   │   ├── network_{darwin,linux}.go
│   │   ├── uptime_{darwin,linux}.go
│   │   └── users_{darwin,linux}.go
│   ├── display/
│   │   └── display.go        # Colour-coded box UI
│   └── config/
│       └── config.go         # Layout and timeout constants
├── Makefile
├── go.mod
└── go.sum
```

Platform-specific implementations use Go build tags (`//go:build darwin`, `//go:build linux`) so only the correct code is compiled for each target.

## Dependencies

- [`github.com/fatih/color`](https://github.com/fatih/color) — terminal colour output (only external dependency)

## License

See LICENSE file for details.
