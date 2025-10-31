package ui

import (
	"sort"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivangsm/blugo/internal/agent"
	"github.com/ivangsm/blugo/internal/bluetooth"
	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/models"
)

// Model represents the state of the TUI application.
type Model struct {
	manager           *bluetooth.Manager
	agent             *agent.Agent
	adapter           *models.Adapter
	devices           map[string]*models.Device
	deviceOrder       []string // Track insertion order of device addresses
	selectedIndex     int
	focusSection      string // "found" o "connected"
	statusMessage     string
	isError           bool
	scanning          bool
	busy              bool
	err               error
	pairingPasskey    *uint32
	waitingForPasskey bool
	width             int // Terminal width
	height            int // Terminal height
	viewport          viewport.Model
	ready             bool // Indicates if the viewport is ready
}

// NewModel creates a new UI model.
func NewModel() Model {
	return Model{
		devices:       make(map[string]*models.Device),
		deviceOrder:   make([]string, 0),
		focusSection:  "found",
		selectedIndex: 0,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// GetFoundDevices returns all devices (available and connected).
// Paired devices are sorted to the top, then maintains stable insertion order.
func (m Model) GetFoundDevices() []*models.Device {
	devices := make([]*models.Device, 0)

	// Get config values
	hideUnnamed := false
	minRSSI := -100
	if config.Global != nil {
		hideUnnamed = config.Global.HideUnnamedDevices
		minRSSI = config.Global.MinRSSIThreshold
	}

	// Use deviceOrder to maintain stable ordering
	for _, addr := range m.deviceOrder {
		dev, exists := m.devices[addr]
		if !exists {
			continue
		}

		// Filter unnamed devices if configured
		if hideUnnamed && dev.Name == "" && dev.Alias == "" {
			continue
		}

		// Filter by RSSI threshold if configured (only for non-connected devices)
		if !dev.Connected && dev.RSSI != 0 && dev.RSSI < int16(minRSSI) {
			continue
		}

		devices = append(devices, dev)
	}

	// Sort: paired devices first, then by insertion order
	sort.Slice(devices, func(i, j int) bool {
		// Paired devices come first
		if devices[i].Paired && !devices[j].Paired {
			return true
		}
		if !devices[i].Paired && devices[j].Paired {
			return false
		}
		// For devices with same paired status, maintain insertion order
		return false // Keep original order
	})

	return devices
}

// GetConnectedDevices returns the connected devices.
func (m Model) GetConnectedDevices() []*models.Device {
	devices := make([]*models.Device, 0)
	for _, dev := range m.devices {
		if dev.Connected {
			devices = append(devices, dev)
		}
	}
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].LastSeen.Before(devices[j].LastSeen)
	})
	return devices
}

// GetSelectedDevice returns the currently selected device.
func (m Model) GetSelectedDevice() *models.Device {
	var devices []*models.Device
	if m.focusSection == "connected" {
		devices = m.GetConnectedDevices()
	} else {
		devices = m.GetFoundDevices()
	}

	if m.selectedIndex < len(devices) {
		return devices[m.selectedIndex]
	}
	return nil
}
