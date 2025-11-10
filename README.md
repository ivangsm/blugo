# BLUGO - Bluetooth Manager for Linux

> A minimalist Bluetooth manager with a modern Terminal User Interface (TUI) for Linux, built with Go and BlueZ.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![BlueZ](https://img.shields.io/badge/BlueZ-5.0+-blue.svg)](http://www.bluez.org/)

[Español](README.es.md) | English

---

## Features

### Device Management
- **Automatic scanning** of nearby Bluetooth devices
- **Automatic pairing** with passkey authentication support
- **Connect/disconnect** devices easily
- **Forget devices** to remove pairing from system
- **Detailed information**: name, MAC address, signal strength (RSSI), and device type
- **Battery indicator** with dynamic colors for compatible devices

### Adapter Control
- **Power control**: Turn Bluetooth adapter on/off (key `P`)
- **Discoverable mode**: Make adapter visible to other devices (key `V`)
- **Pairable mode**: Allow pairing with new devices (key `B`)
- **Adapter information**: View detailed adapter status and configuration

### Modern Interface
- **Responsive design**: Adapts to any terminal size
- **Clean layout**: Single-column design with proper spacing
- **Visual feedback**: Icons, badges, and color-coded status
- **Real-time updates**: Automatic refresh every 2 seconds
- **Scan control**: Pause/resume scanning with visual indicator
- **Alt-screen mode**: Clean terminal on exit and resize
- **Internationalization**: Full i18n support (English/Spanish), easily extensible for more languages

---

## Screenshots

```
╭────────────────────────────────────────────────────────────────╮
│                  🔵 BLUGO - Bluetooth Manager                  │
│                         🔍 Scanning                            │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│                                                                │
│ 📡 Available Devices (3)                                       │
│ ────────────────────────────────────────────────────────────   │
│                                                                │
│ ▶ 🎧 Sony WH-1000XM4 (AA:BB:CC:DD:EE:FF)                      │
│   | -45 dBm | 🔋 85%                                           │
│   [CONNECTED]                                                  │
│                                                                │
│   ⌨️  Keychron K3 (11:22:33:44:55:66)                          │
│   | -38 dBm | 🔋 60%                                           │
│   [PAIRED]                                                     │
│                                                                │
│   🖱️  Logitech MX Master 3 (FF:EE:DD:CC:BB:AA)                │
│   | -52 dBm | 🪫 12%                                           │
│                                                                │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│                                                                │
│ Adapter Info                                                   │
│ ────────────────────────────────────────────────────────────   │
│ Name:          hci0                                            │
│ Alias:         My Laptop                                       │
│ Power:         ✓ On                                            │
│ Pairable:      ✓ On                                            │
│ Discoverable:  ✗ Off                                           │
│                                                                │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│ ↑/k: Select  Enter: Connect/Disconnect  D/X: Forget           │
│ S: Toggle Scan  P: Power  V: Discoverable  B: Pairable        │
│ L: Language  Q: Quit                                           │
╰────────────────────────────────────────────────────────────────╯
```

---

## Requirements

- **Linux** with BlueZ installed
- **Go 1.25** or higher
- **Bluetooth adapter** compatible with BlueZ
- **DBus** system access

### Installing BlueZ

**Arch Linux / Manjaro:**
```bash
sudo pacman -S bluez bluez-utils
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```

**Ubuntu / Debian:**
```bash
sudo apt install bluez bluetooth
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```

**Fedora:**
```bash
sudo dnf install bluez bluez-tools
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```

---

## Installation

### Option 1: Using Make (Recommended)

```bash
# Clone the repository
git clone https://github.com/ivangsm/blugo.git
cd blugo

# Build
make build

# (Optional) Install system-wide
make install
```

### Option 2: Manual with Go

```bash
# Clone the repository
git clone https://github.com/ivangsm/blugo.git
cd blugo

# Download dependencies
go mod download

# Build from cmd/blugo
go build -o blugo ./cmd/blugo

# (Optional) Install system-wide
sudo mv blugo /usr/local/bin/
```

---

## Usage

Simply run the binary:
```bash
./blugo
```

Or if installed system-wide:
```bash
blugo
```

### Keyboard Controls

**Navigation:**
- `↑/↓` or `k/j`: Navigate between devices
- `PgUp/PgDn`: Scroll viewport by page
- `Ctrl+↑/↓`: Scroll viewport by line
- `Home/End`: Jump to top/bottom of list
- `r`: Manually refresh device list

**Device Actions:**
- `Enter`: Connect to available device / Disconnect from connected device
- `d` or `x`: Forget device (disconnect and remove pairing)
- `s`: Pause/resume device scanning

**Adapter Control:**
- `p`: Turn Bluetooth adapter on/off
- `v`: Toggle Discoverable mode
- `b`: Toggle Pairable mode
- `l`: Switch language (English/Spanish)

**General:**
- `q` or `Ctrl+C`: Exit application

**During Pairing:**
- `Enter` or `y`: Confirm pairing code
- `n` or `Esc`: Cancel pairing

---

### Project Structure

```
blugo/
├── cmd/blugo/              # Application entry point
├── internal/
│   ├── models/           # Data models
│   ├── agent/            # Bluetooth pairing agent
│   ├── bluetooth/        # Bluetooth/DBus management
│   └── ui/               # Terminal User Interface
│       ├── styles.go     # Lipgloss styles
│       ├── components.go # Reusable UI components
│       ├── model.go      # Application state
│       ├── update.go     # Update logic (TEA)
│       ├── view.go       # Rendering logic (TEA)
│       ├── messages.go   # Message types
│       └── commands.go   # Async commands
├── Makefile
├── Dockerfile
├── go.mod
└── README.md
```

### Design Principles

- **Separation of Concerns**: Each package has a single, clear responsibility
- **Single Responsibility**: Each file handles a specific aspect
- **Dependency Inversion**: High-level modules don't depend on low-level implementation details
- **The Elm Architecture**: Reactive UI with Model-Update-View pattern

### Technologies

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: TUI framework based on The Elm Architecture
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)**: Terminal styling library
- **[godbus](https://github.com/godbus/dbus)**: DBus client for Go

---

## Development

### Available Make Commands

```bash
make build        # Build the application
make run          # Build and run
make install      # Install to /usr/local/bin
make clean        # Clean build artifacts
make test         # Run tests
make fmt          # Format code
make help         # View all commands
```

### Adding New Features

**New Bluetooth functionality:**
1. Add method in `bluetooth/adapter.go` or `bluetooth/device.go`
2. Create command in `ui/commands.go`
3. Add handler in `ui/update.go`
4. Update view in `ui/view.go` if needed

**New UI section:**
1. Add state to model in `ui/model.go`
2. Create message in `ui/messages.go`
3. Implement handler in `ui/update.go`
4. Create rendering function in `ui/view.go`

---

## Docker

The application can run in a Docker container, though it requires privileged access to DBus and Bluetooth hardware:

**Build the image:**
```bash
docker build -t gob .
```

**Run:**
```bash
docker run --rm -it --privileged --net=host \
  -v /var/run/dbus:/var/run/dbus \
  gob
```

**Note:** Docker usage is limited due to Bluetooth hardware access requirements. Native installation is recommended.

---

## Roadmap

### Current Version
- ✅ Battery indicator for compatible devices
- ✅ Complete adapter control (power, discoverable, pairable)
- ✅ Modern responsive TUI with proper layouts
- ✅ Color-coded badges and status indicators
- ✅ Real-time scanning with pause/resume
- ✅ Clean terminal handling (alt-screen mode)
- ✅ Persistent configuration
- ✅ Unit and integration tests
- ✅ Configuration file support (TOML/YAML)

### Planned Features
- [ ] Support for multiple Bluetooth adapters
- [ ] Enhanced logging and debugging
- [ ] Customizable color themes

---

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

MIT

---

## Author

Ivan - [@ivangsm](https://github.com/ivangsm)

---

## Acknowledgments

- BlueZ project for the Linux Bluetooth stack
- Charm.sh for excellent TUI tools
- The Go community for amazing libraries and support
