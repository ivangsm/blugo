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

// InitializeCmd inicializa el manager de Bluetooth y el agente.
func InitializeCmd(program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		manager, err := bluetooth.NewManager()
		if err != nil {
			return InitMsg{Err: err}
		}

		// Crear y registrar el agente
		btAgent := agent.NewAgent(program)
		err = btAgent.Register(manager.GetConnection())
		if err != nil {
			// No es crítico, la app funcionará pero puede requerir pairing manual
			fmt.Fprintf(os.Stderr, "⚠️  Advertencia: No se pudo registrar agente de pairing: %v\n", err)
			fmt.Fprintf(os.Stderr, "   La app funcionará pero algunos dispositivos pueden requerir pairing manual.\n")
		}

		// Iniciar descubrimiento (if enabled in config)
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

// toggleScanningCmd alterna el estado de escaneo.
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

// updateDevicesCmd actualiza la lista de dispositivos.
func updateDevicesCmd(manager *bluetooth.Manager) tea.Cmd {
	return func() tea.Msg {
		devices, err := manager.GetDevices()
		if err != nil {
			return StatusMsg{Message: fmt.Sprintf("Error al obtener dispositivos: %s", err), IsError: true}
		}
		return DeviceUpdateMsg{Devices: devices}
	}
}

// connectToDeviceCmd conecta a un dispositivo.
func connectToDeviceCmd(manager *bluetooth.Manager, dev *models.Device) tea.Cmd {
	return func() tea.Msg {
		// Si no está pareado, intentar pairing
		if !dev.Paired {
			err := manager.PairDevice(dev.Path)
			if err != nil {
				return ConnectResultMsg{Address: dev.Address, Success: false, Err: fmt.Errorf("error al parear: %w", err)}
			}

			// Confiar en el dispositivo (if enabled in config)
			if config.Global != nil && config.Global.AutoTrustOnPair {
				_ = manager.TrustDevice(dev.Path)
			}

			// Esperar un momento después del pairing (configurable)
			delay := 1000 // Default 1 second
			if config.Global != nil {
				delay = config.Global.PairingDelay
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		// Small delay before attempting connection (helps with reconnection after disconnect)
		time.Sleep(500 * time.Millisecond)

		// Conectar
		err := manager.ConnectDevice(dev.Path)
		if err != nil {
			return ConnectResultMsg{Address: dev.Address, Success: false, Err: fmt.Errorf("error al conectar: %w", err)}
		}

		return ConnectResultMsg{Address: dev.Address, Success: true}
	}
}

// disconnectFromDeviceCmd desconecta de un dispositivo.
func disconnectFromDeviceCmd(manager *bluetooth.Manager, dev *models.Device) tea.Cmd {
	return func() tea.Msg {
		err := manager.DisconnectDevice(dev.Path)
		if err != nil {
			return ConnectResultMsg{Address: dev.Address, Success: false, Err: fmt.Errorf("error al desconectar: %w", err)}
		}
		return ConnectResultMsg{Address: dev.Address, Success: true}
	}
}

// forgetDeviceCmd olvida (elimina) un dispositivo.
func forgetDeviceCmd(manager *bluetooth.Manager, dev *models.Device) tea.Cmd {
	return func() tea.Msg {
		// Desconectar primero si está conectado
		if dev.Connected {
			_ = manager.DisconnectDevice(dev.Path)
			// Wait before removing (configurable)
			delay := 500 // Default 500ms
			if config.Global != nil {
				delay = config.Global.DisconnectDelay
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		// Eliminar el dispositivo
		err := manager.RemoveDevice(dev.Path)
		if err != nil {
			return StatusMsg{Message: fmt.Sprintf("Error al olvidar: %s", err), IsError: true}
		}

		return ForgetDeviceMsg{Address: dev.Address, Message: fmt.Sprintf("Dispositivo %s olvidado", dev.GetDisplayName())}
	}
}

// waitForPasskeyCmd espera a que se reciba un passkey.
func waitForPasskeyCmd(agent *agent.Agent) tea.Cmd {
	if agent == nil {
		return nil
	}

	return func() tea.Msg {
		passkey := <-agent.GetPasskeyChannel()
		return PasskeyDisplayMsg{Passkey: passkey}
	}
}

// tickCmd genera un tick periódico.
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

// updateAdapterInfoCmd actualiza la información del adaptador.
func updateAdapterInfoCmd(manager *bluetooth.Manager) tea.Cmd {
	return func() tea.Msg {
		adapter, err := manager.GetAdapterInfo()
		if err != nil {
			return StatusMsg{Message: fmt.Sprintf("Error al obtener info del adaptador: %s", err), IsError: true}
		}
		return AdapterUpdateMsg{Adapter: adapter}
	}
}

// toggleAdapterPoweredCmd enciende o apaga el adaptador.
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

// toggleAdapterDiscoverableCmd activa o desactiva el modo discoverable.
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

// toggleAdapterPairableCmd activa o desactiva el modo pairable.
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
