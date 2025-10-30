package ui

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivangsm/blugo/internal/agent"
	"github.com/ivangsm/blugo/internal/bluetooth"
	"github.com/ivangsm/blugo/internal/models"
)

// Model representa el estado de la aplicación TUI.
type Model struct {
	manager           *bluetooth.Manager
	agent             *agent.Agent
	adapter           *models.Adapter
	devices           map[string]*models.Device
	selectedIndex     int
	focusSection      string // "found" o "connected"
	statusMessage     string
	isError           bool
	scanning          bool
	busy              bool
	err               error
	pairingPasskey    *uint32
	waitingForPasskey bool
	width             int // Ancho de la terminal
	height            int // Alto de la terminal
}

// NewModel crea un nuevo modelo de UI.
func NewModel() Model {
	return Model{
		devices:       make(map[string]*models.Device),
		focusSection:  "found",
		selectedIndex: 0,
	}
}

// Init inicializa el modelo.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// GetFoundDevices devuelve los dispositivos disponibles (no conectados).
func (m Model) GetFoundDevices() []*models.Device {
	devices := make([]*models.Device, 0)
	for _, dev := range m.devices {
		if !dev.Connected {
			devices = append(devices, dev)
		}
	}
	sort.Slice(devices, func(i, j int) bool {
		// Pareados primero
		if devices[i].Paired != devices[j].Paired {
			return devices[i].Paired
		}
		// Luego por RSSI
		return devices[i].RSSI > devices[j].RSSI
	})
	return devices
}

// GetConnectedDevices devuelve los dispositivos conectados.
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

// GetSelectedDevice devuelve el dispositivo actualmente seleccionado.
func (m Model) GetSelectedDevice() *models.Device {
	var devices []*models.Device
	if m.focusSection == "found" {
		devices = m.GetFoundDevices()
	} else {
		devices = m.GetConnectedDevices()
	}

	if m.selectedIndex < len(devices) {
		return devices[m.selectedIndex]
	}
	return nil
}
