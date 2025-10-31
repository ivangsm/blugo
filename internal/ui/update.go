package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/i18n"
)

// updateViewportContent updates the viewport with current content
func (m *Model) updateViewportContent() {
	if m.ready {
		content := m.renderFullContent()
		m.viewport.SetContent(content)
	}
}

// Update handles model updates.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// First time - initialize viewport
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			// Update viewport size
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

		// Update viewport content
		m.updateViewportContent()

		// tea.ClearScreen clears the screen during resize
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

	case ForgetDeviceMsg:
		return m.handleForgetDevice(msg)

	case AdapterUpdateMsg:
		return m.handleAdapterUpdate(msg)

	case AdapterPropertyChangedMsg:
		return m.handleAdapterPropertyChanged(msg)

	case TickMsg:
		return m.handleTick()

	case tea.MouseMsg:
		// Handle mouse wheel scrolling
		switch msg.Type {
		case tea.MouseWheelUp:
			m.viewport.LineUp(3)
		case tea.MouseWheelDown:
			m.viewport.LineDown(3)
		}
		return m, nil
	}

	return m, nil
}

// handleKeyPress handles pressed keys.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If we are waiting for passkey confirmation
	if m.pairingPasskey != nil {
		return m.handlePasskeyConfirmation(msg)
	}

	// If we are busy, only allow exit
	if m.busy {
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m.quit()
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m.quit()

	// Viewport scrolling
	case "pgup":
		m.viewport.ViewUp()
		return m, nil

	case "pgdown":
		m.viewport.ViewDown()
		return m, nil

	case "ctrl+up":
		m.viewport.LineUp(3)
		return m, nil

	case "ctrl+down":
		m.viewport.LineDown(3)
		return m, nil

	case "home":
		m.viewport.GotoTop()
		return m, nil

	case "end":
		m.viewport.GotoBottom()
		return m, nil

	// Device navigation
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
			m.updateViewportContent()
		}

	case "down", "j":
		maxIndex := len(m.GetFoundDevices()) - 1
		if m.selectedIndex < maxIndex {
			m.selectedIndex++
			m.updateViewportContent()
		}

	// Tab navigation removed since we only have one section now

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
		// Toggle Powered (turn Bluetooth on/off)
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

		// Save language preference if configured
		if config.Global != nil && config.Global.RememberLanguage {
			// Update config with new language
			newLang := string(i18n.GetCurrentLanguage())
			config.Global.Language = newLang

			// Save config to disk
			_ = config.Global.Save()
		}

		m.updateViewportContent()
		return m, nil
	}

	return m, nil
}

// handlePasskeyConfirmation handles passkey confirmation.
func (m Model) handlePasskeyConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "y":
		// Confirm pairing
		if m.agent != nil {
			m.agent.GetConfirmChannel() <- true
		}
		m.pairingPasskey = nil
		m.statusMessage = i18n.T.StatusConfirmingPairing
		return m, nil

	case "n", "esc":
		// Cancel pairing
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

// handleEnter handles the Enter key.
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}

	dev := m.GetSelectedDevice()
	if dev == nil {
		return m, nil
	}

	if dev.Connected {
		// Disconnect device
		m.busy = true
		m.statusMessage = fmt.Sprintf(i18n.T.Disconnecting, dev.GetDisplayName())
		return m, disconnectFromDeviceCmd(m.manager, dev)
	} else {
		// Connect device
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
	}
}

// handleForget handles the forget device action.
func (m Model) handleForget() (tea.Model, tea.Cmd) {
	if m.manager == nil {
		return m, nil
	}

	dev := m.GetSelectedDevice()
	if dev == nil {
		return m, nil
	}

	if dev.Paired {
		m.busy = true
		m.statusMessage = fmt.Sprintf(i18n.T.Forgetting, dev.GetDisplayName())
		return m, forgetDeviceCmd(m.manager, dev)
	}

	return m, nil
}

// handleInit handles the initialization message.
func (m Model) handleInit(msg InitMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.err = msg.Err
		return m, nil
	}

	m.manager = msg.Manager
	m.agent = msg.Agent
	m.scanning = msg.Scanning // Use the actual scanning state from init
	if msg.Scanning {
		m.statusMessage = i18n.T.ScanEnabled
	} else {
		m.statusMessage = i18n.T.ScanPaused
	}
	m.updateViewportContent()
	return m, tea.Batch(
		updateDevicesCmd(m.manager),
		updateAdapterInfoCmd(m.manager),
	)
}

// handleScanning handles scanning state change.
func (m Model) handleScanning(msg ScanningMsg) (tea.Model, tea.Cmd) {
	m.scanning = msg.Scanning
	if msg.Scanning {
		m.statusMessage = i18n.T.ScanEnabled
	} else {
		m.statusMessage = i18n.T.ScanPaused
	}
	m.updateViewportContent()
	return m, nil
}

// handleDeviceUpdate handles device updates.
func (m Model) handleDeviceUpdate(msg DeviceUpdateMsg) (tea.Model, tea.Cmd) {
	// Update only new or modified devices
	for addr, newDev := range msg.Devices {
		if oldDev, exists := m.devices[addr]; exists {
			// Keep LastSeen if device already existed
			if !oldDev.Connected && !newDev.Connected {
				newDev.LastSeen = oldDev.LastSeen
			}
		} else {
			// New device - add to deviceOrder to maintain stable ordering
			m.deviceOrder = append(m.deviceOrder, addr)
		}
		m.devices[addr] = newDev
	}
	m.updateViewportContent()
	return m, nil
}

// handlePasskeyDisplay handles passkey display.
func (m Model) handlePasskeyDisplay(msg PasskeyDisplayMsg) (tea.Model, tea.Cmd) {
	if m.waitingForPasskey {
		m.pairingPasskey = &msg.Passkey
		m.waitingForPasskey = false
	}
	m.updateViewportContent()
	return m, nil
}

// handleConnectResult handles connection result.
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

	// Auto-close passkey prompt immediately on any connection result (success or failure)
	// This provides better responsiveness
	m.pairingPasskey = nil

	m.updateViewportContent()
	return m, updateDevicesCmd(m.manager)
}

// handleStatus handles status messages.
func (m Model) handleStatus(msg StatusMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	m.statusMessage = msg.Message
	m.isError = msg.IsError
	m.updateViewportContent()
	return m, updateDevicesCmd(m.manager)
}

// handleForgetDevice handles device forgetting.
func (m Model) handleForgetDevice(msg ForgetDeviceMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	m.statusMessage = msg.Message
	m.isError = false

	// Remove the device from our local cache
	delete(m.devices, msg.Address)

	// Remove from deviceOrder if present
	for i, addr := range m.deviceOrder {
		if addr == msg.Address {
			m.deviceOrder = append(m.deviceOrder[:i], m.deviceOrder[i+1:]...)
			break
		}
	}

	// Adjust selectedIndex if necessary
	maxIndex := len(m.GetFoundDevices()) - 1
	if m.selectedIndex > maxIndex {
		m.selectedIndex = maxIndex
	}
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}

	m.updateViewportContent()
	return m, updateDevicesCmd(m.manager)
}

// handleAdapterUpdate handles adapter information update.
func (m Model) handleAdapterUpdate(msg AdapterUpdateMsg) (tea.Model, tea.Cmd) {
	m.adapter = msg.Adapter
	m.updateViewportContent()
	return m, nil
}

// handleAdapterPropertyChanged handles adapter property change.
func (m Model) handleAdapterPropertyChanged(msg AdapterPropertyChangedMsg) (tea.Model, tea.Cmd) {
	m.busy = false

	if msg.Err != nil {
		m.statusMessage = fmt.Sprintf("❌ Error al cambiar %s: %s", msg.Property, msg.Err.Error())
		m.isError = true
		m.updateViewportContent()
		return m, nil
	}

	// Success messages based on property
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

	m.updateViewportContent()
	// Update adapter information
	return m, updateAdapterInfoCmd(m.manager)
}

// handleTick handles periodic tick.
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	cmds = append(cmds, tickCmd())
	if m.manager != nil {
		cmds = append(cmds, updateDevicesCmd(m.manager))
		cmds = append(cmds, updateAdapterInfoCmd(m.manager))
	}
	return m, tea.Batch(cmds...)
}

// quit handles application exit.
func (m Model) quit() (tea.Model, tea.Cmd) {
	if m.manager != nil && m.scanning {
		_ = m.manager.StopDiscovery()
	}
	if m.manager != nil {
		_ = m.manager.Close()
	}
	return m, tea.Quit
}
