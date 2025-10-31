package ui

import (
	"time"

	"github.com/ivangsm/blugo/internal/agent"
	"github.com/ivangsm/blugo/internal/bluetooth"
	"github.com/ivangsm/blugo/internal/models"
)

// InitMsg indicates that initialization has completed.
type InitMsg struct {
	Manager  *bluetooth.Manager
	Agent    *agent.Agent
	Scanning bool // Indicates if scanning was started
	Err      error
}

// ScanningMsg indicates a change in scanning state.
type ScanningMsg struct {
	Scanning bool
}

// DeviceUpdateMsg contains updated devices.
type DeviceUpdateMsg struct {
	Devices map[string]*models.Device
}

// StatusMsg shows a status message.
type StatusMsg struct {
	Message string
	IsError bool
}

// ConnectResultMsg indicates the result of a connection operation.
type ConnectResultMsg struct {
	Address string
	Success bool
	Err     error
}

// PairResultMsg indicates the result of a pairing operation.
type PairResultMsg struct {
	Address string
	Success bool
	Err     error
}

// PasskeyDisplayMsg requests to display a passkey.
type PasskeyDisplayMsg struct {
	Passkey uint32
	Device  string
}

// PasskeyConfirmedMsg indicates that the passkey was confirmed.
type PasskeyConfirmedMsg struct{}

// AdapterUpdateMsg contains updated adapter information.
type AdapterUpdateMsg struct {
	Adapter *models.Adapter
}

// AdapterPropertyChangedMsg indicates that an adapter property changed.
type AdapterPropertyChangedMsg struct {
	Property string
	Success  bool
	Err      error
}

// ForgetDeviceMsg indicates that a device was forgotten.
type ForgetDeviceMsg struct {
	Address string
	Message string
}

// TickMsg is a clock tick for periodic updates.
type TickMsg time.Time
