package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivangsm/blugo/internal/i18n"
)

// Update maneja las actualizaciones del modelo.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// tea.ClearScreen limpia la pantalla durante el resize
		return m, tea.ClearScreen

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case InitMsg:
		return m.handleInit(msg)

	case ScanningMsg:
		return m.handleScanning(msg)

	case DeviceUpdateMsg:
		return m.handleDeviceUpdate(msg)

	case PasskeyDisplayMsg:
		return m.handlePasskeyDisplay(msg)

	case ConnectResultMsg:
		return m.handleConnectResult(msg)

	case StatusMsg:
		return m.handleStatus(msg)

	case AdapterUpdateMsg:
		return m.handleAdapterUpdate(msg)

	case AdapterPropertyChangedMsg:
		return m.handleAdapterPropertyChanged(msg)

	case TickMsg:
		return m.handleTick()
	}

	return m, nil
}

// handleKeyPress maneja las teclas presionadas.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Si estamos esperando confirmación de passkey
	if m.pairingPasskey != nil {
		return m.handlePasskeyConfirmation(msg)
	}

	// Si estamos ocupados, solo permitir salir
	if m.busy {
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m.quit()
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m.quit()

	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}

	case "down", "j":
		maxIndex := 0
		if m.focusSection == "found" {
			maxIndex = len(m.GetFoundDevices()) - 1
		} else {
			maxIndex = len(m.GetConnectedDevices()) - 1
		}
		if m.selectedIndex < maxIndex {
			m.selectedIndex++
		}

	case "tab":
		if m.focusSection == "found" && len(m.GetConnectedDevices()) > 0 {
			m.focusSection = "connected"
			m.selectedIndex = 0
		} else if m.focusSection == "connected" && len(m.GetFoundDevices()) > 0 {
			m.focusSection = "found"
			m.selectedIndex = 0
		}

	case "s":
		if m.manager != nil {
			return m, toggleScanningCmd(m.manager, m.scanning)
		}

	case "enter":
		return m.handleEnter()

	case "d", "x":
		return m.handleForget()

	case "r":
		if m.manager != nil {
			return m, updateDevicesCmd(m.manager)
		}

	case "p":
		// Toggle Powered (encender/apagar Bluetooth)
		if m.manager != nil && m.adapter != nil {
			m.busy = true
			if m.adapter.Powered {
				m.statusMessage = i18n.T.AdapterPoweringOff
			} else {
				m.statusMessage = i18n.T.AdapterPoweringOn
			}
			return m, toggleAdapterPoweredCmd(m.manager, m.adapter.Powered)
		}

	case "v":
		// Toggle Discoverable
		if m.manager != nil && m.adapter != nil {
			m.busy = true
			if m.adapter.Discoverable {
				m.statusMessage = i18n.T.DiscoverableDeactivating
			} else {
				m.statusMessage = i18n.T.DiscoverableActivating
			}
			return m, toggleAdapterDiscoverableCmd(m.manager, m.adapter.Discoverable)
		}

	case "b":
		// Toggle Pairable
		if m.manager != nil && m.adapter != nil {
			m.busy = true
			if m.adapter.Pairable {
				m.statusMessage = i18n.T.PairableDeactivating
			} else {
				m.statusMessage = i18n.T.PairableActivating
			}
			return m, toggleAdapterPairableCmd(m.manager, m.adapter.Pairable)
		}

	case "l":
		// Toggle Language
		i18n.ToggleLanguage()
		return m, nil
	}

	return m, nil
}

// handlePasskeyConfirmation maneja la confirmación del passkey.
func (m Model) handlePasskeyConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "y":
		// Confirmar pairing
		if m.agent != nil {
			m.agent.GetConfirmChannel() <- true
		}
		m.pairingPasskey = nil
		m.statusMessage = "Confirming pairing..."
		return m, nil

	case "n", "esc":
		// Cancelar pairing
		if m.agent != nil {
			m.agent.GetConfirmChannel() <- false
		}
		m.pairingPasskey = nil
		m.busy = false
		m.statusMessage = i18n.T.PairingCancelled
		return m, nil

	case "ctrl+c", "q":
		if m.agent != nil {
			m.agent.GetConfirmChannel() <- false
		}
		return m.quit()
	}

	return m, nil
}

// handleEnter maneja la tecla Enter.
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}

	dev := m.GetSelectedDevice()
	if dev == nil {
		return m, nil
	}

	if m.focusSection == "found" {
		// Conectar dispositivo
		m.busy = true
		if dev.Paired {
			m.statusMessage = fmt.Sprintf(i18n.T.Connecting, dev.GetDisplayName())
		} else {
			m.statusMessage = fmt.Sprintf(i18n.T.Pairing, dev.GetDisplayName())
			m.waitingForPasskey = true
		}
		return m, tea.Batch(
			connectToDeviceCmd(m.manager, dev),
			waitForPasskeyCmd(m.agent),
		)
	} else {
		// Desconectar dispositivo
		m.busy = true
		m.statusMessage = fmt.Sprintf(i18n.T.Disconnecting, dev.GetDisplayName())
		return m, disconnectFromDeviceCmd(m.manager, dev)
	}
}

// handleForget maneja la acción de olvidar un dispositivo.
func (m Model) handleForget() (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}

	dev := m.GetSelectedDevice()
	if dev == nil {
		return m, nil
	}

	if m.focusSection == "connected" || (m.focusSection == "found" && dev.Paired) {
		m.busy = true
		m.statusMessage = fmt.Sprintf(i18n.T.Forgetting, dev.GetDisplayName())
		return m, forgetDeviceCmd(m.manager, dev)
	}

	return m, nil
}

// handleInit maneja el mensaje de inicialización.
func (m Model) handleInit(msg InitMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.err = msg.Err
		return m, nil
	}

	m.manager = msg.Manager
	m.agent = msg.Agent
	m.scanning = true
	m.statusMessage = "Scanning Bluetooth devices..."
	return m, tea.Batch(
		updateDevicesCmd(m.manager),
		updateAdapterInfoCmd(m.manager),
	)
}

// handleScanning maneja el cambio de estado de escaneo.
func (m Model) handleScanning(msg ScanningMsg) (tea.Model, tea.Cmd) {
	m.scanning = msg.Scanning
	if msg.Scanning {
		m.statusMessage = i18n.T.ScanEnabled
	} else {
		m.statusMessage = i18n.T.ScanPaused
	}
	return m, nil
}

// handleDeviceUpdate maneja la actualización de dispositivos.
func (m Model) handleDeviceUpdate(msg DeviceUpdateMsg) (tea.Model, tea.Cmd) {
	// Actualizar solo dispositivos nuevos o modificados
	for addr, newDev := range msg.Devices {
		if oldDev, exists := m.devices[addr]; exists {
			// Mantener LastSeen si el dispositivo ya existía
			if !oldDev.Connected && !newDev.Connected {
				newDev.LastSeen = oldDev.LastSeen
			}
		}
		m.devices[addr] = newDev
	}
	return m, nil
}

// handlePasskeyDisplay maneja la visualización del passkey.
func (m Model) handlePasskeyDisplay(msg PasskeyDisplayMsg) (tea.Model, tea.Cmd) {
	if m.waitingForPasskey {
		m.pairingPasskey = &msg.Passkey
		m.waitingForPasskey = false
	}
	return m, nil
}

// handleConnectResult maneja el resultado de una conexión.
func (m Model) handleConnectResult(msg ConnectResultMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	m.waitingForPasskey = false

	if msg.Err != nil {
		m.statusMessage = fmt.Sprintf("❌ %s", msg.Err.Error())
		m.isError = true
	} else {
		if dev, ok := m.devices[msg.Address]; ok {
			if dev.Connected {
				m.statusMessage = fmt.Sprintf(i18n.T.Connected, dev.GetDisplayName())
			} else {
				if dev.Paired {
					m.statusMessage = fmt.Sprintf(i18n.T.DisconnectedPaired, dev.GetDisplayName())
				} else {
					m.statusMessage = fmt.Sprintf(i18n.T.Disconnected, dev.GetDisplayName())
				}
			}
		}
		m.isError = false
	}

	return m, updateDevicesCmd(m.manager)
}

// handleStatus maneja mensajes de estado.
func (m Model) handleStatus(msg StatusMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	m.statusMessage = msg.Message
	m.isError = msg.IsError
	return m, updateDevicesCmd(m.manager)
}

// handleAdapterUpdate maneja la actualización de información del adaptador.
func (m Model) handleAdapterUpdate(msg AdapterUpdateMsg) (tea.Model, tea.Cmd) {
	m.adapter = msg.Adapter
	return m, nil
}

// handleAdapterPropertyChanged maneja el cambio de una propiedad del adaptador.
func (m Model) handleAdapterPropertyChanged(msg AdapterPropertyChangedMsg) (tea.Model, tea.Cmd) {
	m.busy = false

	if msg.Err != nil {
		m.statusMessage = fmt.Sprintf("❌ Error al cambiar %s: %s", msg.Property, msg.Err.Error())
		m.isError = true
		return m, nil
	}

	// Mensajes de éxito según la propiedad
	switch msg.Property {
	case "Powered":
		if m.adapter != nil && m.adapter.Powered {
			m.statusMessage = i18n.T.AdapterPoweredOff
		} else {
			m.statusMessage = i18n.T.AdapterPoweredOn
		}
	case "Discoverable":
		if m.adapter != nil && m.adapter.Discoverable {
			m.statusMessage = i18n.T.DiscoverableOff
		} else {
			m.statusMessage = i18n.T.DiscoverableOn
		}
	case "Pairable":
		if m.adapter != nil && m.adapter.Pairable {
			m.statusMessage = i18n.T.PairableOff
		} else {
			m.statusMessage = i18n.T.PairableOn
		}
	}

	m.isError = false

	// Actualizar información del adaptador
	return m, updateAdapterInfoCmd(m.manager)
}

// handleTick maneja el tick periódico.
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	cmds = append(cmds, tickCmd())
	if m.manager != nil {
		cmds = append(cmds, updateDevicesCmd(m.manager))
		cmds = append(cmds, updateAdapterInfoCmd(m.manager))
	}
	return m, tea.Batch(cmds...)
}

// quit maneja la salida de la aplicación.
func (m Model) quit() (tea.Model, tea.Cmd) {
	if m.manager != nil && m.scanning {
		_ = m.manager.StopDiscovery()
	}
	if m.manager != nil {
		_ = m.manager.Close()
	}
	return m, tea.Quit
}
