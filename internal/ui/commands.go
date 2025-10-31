package ui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivangsm/blugo/internal/agent"
	"github.com/ivangsm/blugo/internal/bluetooth"
	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/i18n"
	"github.com/ivangsm/blugo/internal/models"
)

// InitializeCmd initializes the Bluetooth manager and agent.
func InitializeCmd(program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		manager, err := bluetooth.NewManager()
		if err != nil {
			return InitMsg{Err: err}
		}

		// Create and register the agent
		btAgent := agent.NewAgent(program)
		err = btAgent.Register(manager.GetConnection())
		if err != nil {
			// Not critical, the app will work but may require manual pairing
			fmt.Fprintf(os.Stderr, "%s: %v\n", i18n.T.WarningAgentRegistration, err)
			fmt.Fprintf(os.Stderr, "%s\n", i18n.T.WarningAgentRegistrationDetail)
		}

		// Start discovery (if enabled in config)
		autoStart := true // Default
		if config.Global != nil {
			autoStart = config.Global.AutoStartScanning
		}

		scanningStarted := false
		if autoStart {
			err = manager.StartDiscovery()
			if err != nil {
				return InitMsg{Err: fmt.Errorf("no se pudo iniciar descubrimiento: %w", err)}
			}
			scanningStarted = true
		}

		return InitMsg{Manager: manager, Agent: btAgent, Scanning: scanningStarted}
	}
}

// toggleScanningCmd toggles scanning state.
func toggleScanningCmd(manager *bluetooth.Manager, currentlyScanning bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		if currentlyScanning {
			err = manager.StopDiscovery()
		} else {
			err = manager.StartDiscovery()
		}

		if err != nil {
			return StatusMsg{Message: fmt.Sprintf(i18n.T.ErrorScanToggle, err), IsError: true}
		}

		return ScanningMsg{Scanning: !currentlyScanning}
	}
}

// updateDevicesCmd updates the device list.
func updateDevicesCmd(manager *bluetooth.Manager) tea.Cmd {
	return func() tea.Msg {
		devices, err := manager.GetDevices()
		if err != nil {
			return StatusMsg{Message: fmt.Sprintf(i18n.T.ErrorGetDevices+": %s", err), IsError: true}
		}
		return DeviceUpdateMsg{Devices: devices}
	}
}

// connectToDeviceCmd connects to a device.
func connectToDeviceCmd(manager *bluetooth.Manager, dev *models.Device) tea.Cmd {
	return func() tea.Msg {
		// If not paired, try pairing
		if !dev.Paired {
			err := manager.PairDevice(dev.Path)
			if err != nil {
				return ConnectResultMsg{Address: dev.Address, Success: false, Err: fmt.Errorf("error al parear: %w", err)}
			}

			// Trust the device (if enabled in config)
			if config.Global != nil && config.Global.AutoTrustOnPair {
				_ = manager.TrustDevice(dev.Path)
			}

			// Wait a moment after pairing (configurable)
			delay := 1000 // Default 1 second
			if config.Global != nil {
				delay = config.Global.PairingDelay
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		// Small delay before attempting connection (helps with reconnection after disconnect)
		time.Sleep(500 * time.Millisecond)

		// Connect
		err := manager.ConnectDevice(dev.Path)
		if err != nil {
			return ConnectResultMsg{Address: dev.Address, Success: false, Err: fmt.Errorf("error al conectar: %w", err)}
		}

		return ConnectResultMsg{Address: dev.Address, Success: true}
	}
}

// disconnectFromDeviceCmd disconnects from a device.
func disconnectFromDeviceCmd(manager *bluetooth.Manager, dev *models.Device) tea.Cmd {
	return func() tea.Msg {
		err := manager.DisconnectDevice(dev.Path)
		if err != nil {
			return ConnectResultMsg{Address: dev.Address, Success: false, Err: fmt.Errorf("error al desconectar: %w", err)}
		}
		return ConnectResultMsg{Address: dev.Address, Success: true}
	}
}

// forgetDeviceCmd forgets (removes) a device.
func forgetDeviceCmd(manager *bluetooth.Manager, dev *models.Device) tea.Cmd {
	return func() tea.Msg {
		// Disconnect first if connected
		if dev.Connected {
			_ = manager.DisconnectDevice(dev.Path)
			// Wait before removing (configurable)
			delay := 500 // Default 500ms
			if config.Global != nil {
				delay = config.Global.DisconnectDelay
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		// Remove the device
		err := manager.RemoveDevice(dev.Path)
		if err != nil {
			return StatusMsg{Message: fmt.Sprintf(i18n.T.ErrorForgetDevice+": %s", err), IsError: true}
		}

		return ForgetDeviceMsg{Address: dev.Address, Message: fmt.Sprintf(i18n.T.Forgotten+": %s", dev.GetDisplayName())}
	}
}

// waitForPasskeyCmd waits for a passkey to be received.
func waitForPasskeyCmd(agent *agent.Agent) tea.Cmd {
	if agent == nil {
		return nil
	}

	return func() tea.Msg {
		passkey := <-agent.GetPasskeyChannel()
		return PasskeyDisplayMsg{Passkey: passkey}
	}
}

// tickCmd generates a periodic tick.
func tickCmd() tea.Cmd {
	interval := 2 // Default fallback
	if config.Global != nil {
		interval = config.Global.RefreshInterval
		// Validate range
		if interval < 1 {
			interval = 1
		} else if interval > 10 {
			interval = 10
		}
	}
	return tea.Tick(time.Duration(interval)*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// updateAdapterInfoCmd updates the adapter information.
func updateAdapterInfoCmd(manager *bluetooth.Manager) tea.Cmd {
	return func() tea.Msg {
		adapter, err := manager.GetAdapterInfo()
		if err != nil {
			return StatusMsg{Message: fmt.Sprintf(i18n.T.ErrorGetAdapterInfo+": %s", err), IsError: true}
		}
		return AdapterUpdateMsg{Adapter: adapter}
	}
}

// toggleAdapterPoweredCmd turns the adapter on or off.
func toggleAdapterPoweredCmd(manager *bluetooth.Manager, currentState bool) tea.Cmd {
	return func() tea.Msg {
		newState := !currentState
		err := manager.SetAdapterPowered(newState)
		if err != nil {
			return AdapterPropertyChangedMsg{Property: "Powered", Success: false, Err: err}
		}
		return AdapterPropertyChangedMsg{Property: "Powered", Success: true}
	}
}

// toggleAdapterDiscoverableCmd enables or disables discoverable mode.
func toggleAdapterDiscoverableCmd(manager *bluetooth.Manager, currentState bool) tea.Cmd {
	return func() tea.Msg {
		newState := !currentState
		err := manager.SetAdapterDiscoverable(newState)
		if err != nil {
			return AdapterPropertyChangedMsg{Property: "Discoverable", Success: false, Err: err}
		}
		return AdapterPropertyChangedMsg{Property: "Discoverable", Success: true}
	}
}

// toggleAdapterPairableCmd enables or disables pairable mode.
func toggleAdapterPairableCmd(manager *bluetooth.Manager, currentState bool) tea.Cmd {
	return func() tea.Msg {
		newState := !currentState
		err := manager.SetAdapterPairable(newState)
		if err != nil {
			return AdapterPropertyChangedMsg{Property: "Pairable", Success: false, Err: err}
		}
		return AdapterPropertyChangedMsg{Property: "Pairable", Success: true}
	}
}
