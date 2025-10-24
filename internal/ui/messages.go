package ui

import (
	"time"

	"github.com/ivangsm/gob/internal/agent"
	"github.com/ivangsm/gob/internal/bluetooth"
	"github.com/ivangsm/gob/internal/models"
)

// InitMsg indica que la inicialización ha completado.
type InitMsg struct {
	Manager *bluetooth.Manager
	Agent   *agent.Agent
	Err     error
}

// ScanningMsg indica un cambio en el estado de escaneo.
type ScanningMsg struct {
	Scanning bool
}

// DeviceUpdateMsg contiene dispositivos actualizados.
type DeviceUpdateMsg struct {
	Devices map[string]*models.Device
}

// StatusMsg muestra un mensaje de estado.
type StatusMsg struct {
	Message string
	IsError bool
}

// ConnectResultMsg indica el resultado de una operación de conexión.
type ConnectResultMsg struct {
	Address string
	Success bool
	Err     error
}

// PairResultMsg indica el resultado de una operación de pairing.
type PairResultMsg struct {
	Address string
	Success bool
	Err     error
}

// PasskeyDisplayMsg solicita mostrar un passkey.
type PasskeyDisplayMsg struct {
	Passkey uint32
	Device  string
}

// PasskeyConfirmedMsg indica que el passkey fue confirmado.
type PasskeyConfirmedMsg struct{}

// TickMsg es un tick del reloj para actualizaciones periódicas.
type TickMsg time.Time
