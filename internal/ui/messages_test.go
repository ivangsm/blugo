package ui

import (
	"errors"
	"testing"
	"time"

	"github.com/ivangsm/blugo/internal/models"
)

func TestInitMsg(t *testing.T) {
	msg := InitMsg{
		Manager:  nil,
		Agent:    nil,
		Scanning: true,
		Err:      errors.New("test error"),
	}

	if msg.Scanning != true {
		t.Errorf("InitMsg.Scanning = %v, want true", msg.Scanning)
	}
	if msg.Err == nil {
		t.Errorf("InitMsg.Err should not be nil")
	}
	if msg.Manager != nil {
		t.Errorf("InitMsg.Manager should be nil")
	}
	if msg.Agent != nil {
		t.Errorf("InitMsg.Agent should be nil")
	}
}

func TestScanningMsg(t *testing.T) {
	tests := []struct {
		name     string
		scanning bool
	}{
		{
			name:     "scanning enabled",
			scanning: true,
		},
		{
			name:     "scanning disabled",
			scanning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ScanningMsg{Scanning: tt.scanning}
			if msg.Scanning != tt.scanning {
				t.Errorf("ScanningMsg.Scanning = %v, want %v", msg.Scanning, tt.scanning)
			}
		})
	}
}

func TestDeviceUpdateMsg(t *testing.T) {
	devices := map[string]*models.Device{
		"AA:BB:CC:DD:EE:FF": {
			Address: "AA:BB:CC:DD:EE:FF",
			Name:    "Test Device",
		},
	}

	msg := DeviceUpdateMsg{Devices: devices}

	if msg.Devices == nil {
		t.Errorf("DeviceUpdateMsg.Devices should not be nil")
	}
	if len(msg.Devices) != 1 {
		t.Errorf("DeviceUpdateMsg.Devices length = %d, want 1", len(msg.Devices))
	}
}

func TestStatusMsg(t *testing.T) {
	tests := []struct {
		name    string
		message string
		isError bool
	}{
		{
			name:    "success status",
			message: "Connected successfully",
			isError: false,
		},
		{
			name:    "error status",
			message: "Connection failed",
			isError: true,
		},
		{
			name:    "empty message",
			message: "",
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := StatusMsg{
				Message: tt.message,
				IsError: tt.isError,
			}

			if msg.Message != tt.message {
				t.Errorf("StatusMsg.Message = %v, want %v", msg.Message, tt.message)
			}
			if msg.IsError != tt.isError {
				t.Errorf("StatusMsg.IsError = %v, want %v", msg.IsError, tt.isError)
			}
		})
	}
}

func TestConnectResultMsg(t *testing.T) {
	tests := []struct {
		name    string
		address string
		success bool
		err     error
	}{
		{
			name:    "successful connection",
			address: "AA:BB:CC:DD:EE:FF",
			success: true,
			err:     nil,
		},
		{
			name:    "failed connection",
			address: "11:22:33:44:55:66",
			success: false,
			err:     errors.New("connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ConnectResultMsg{
				Address: tt.address,
				Success: tt.success,
				Err:     tt.err,
			}

			if msg.Address != tt.address {
				t.Errorf("ConnectResultMsg.Address = %v, want %v", msg.Address, tt.address)
			}
			if msg.Success != tt.success {
				t.Errorf("ConnectResultMsg.Success = %v, want %v", msg.Success, tt.success)
			}
			if (msg.Err != nil) != (tt.err != nil) {
				t.Errorf("ConnectResultMsg.Err presence mismatch")
			}
		})
	}
}

func TestPairResultMsg(t *testing.T) {
	tests := []struct {
		name    string
		address string
		success bool
		err     error
	}{
		{
			name:    "successful pairing",
			address: "AA:BB:CC:DD:EE:FF",
			success: true,
			err:     nil,
		},
		{
			name:    "failed pairing",
			address: "11:22:33:44:55:66",
			success: false,
			err:     errors.New("pairing failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := PairResultMsg{
				Address: tt.address,
				Success: tt.success,
				Err:     tt.err,
			}

			if msg.Address != tt.address {
				t.Errorf("PairResultMsg.Address = %v, want %v", msg.Address, tt.address)
			}
			if msg.Success != tt.success {
				t.Errorf("PairResultMsg.Success = %v, want %v", msg.Success, tt.success)
			}
			if (msg.Err != nil) != (tt.err != nil) {
				t.Errorf("PairResultMsg.Err presence mismatch")
			}
		})
	}
}

func TestPasskeyDisplayMsg(t *testing.T) {
	tests := []struct {
		name    string
		passkey uint32
		device  string
	}{
		{
			name:    "valid passkey",
			passkey: 123456,
			device:  "My Device",
		},
		{
			name:    "zero passkey",
			passkey: 0,
			device:  "Another Device",
		},
		{
			name:    "max passkey",
			passkey: 999999,
			device:  "Test Device",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := PasskeyDisplayMsg{
				Passkey: tt.passkey,
				Device:  tt.device,
			}

			if msg.Passkey != tt.passkey {
				t.Errorf("PasskeyDisplayMsg.Passkey = %v, want %v", msg.Passkey, tt.passkey)
			}
			if msg.Device != tt.device {
				t.Errorf("PasskeyDisplayMsg.Device = %v, want %v", msg.Device, tt.device)
			}
		})
	}
}

func TestPasskeyConfirmedMsg(t *testing.T) {
	// Empty struct, just verify it can be created
	msg := PasskeyConfirmedMsg{}
	_ = msg
	t.Log("PasskeyConfirmedMsg created successfully")
}

func TestAdapterUpdateMsg(t *testing.T) {
	adapter := &models.Adapter{
		Address: "AA:BB:CC:DD:EE:FF",
		Name:    "hci0",
		Powered: true,
	}

	msg := AdapterUpdateMsg{Adapter: adapter}

	if msg.Adapter == nil {
		t.Errorf("AdapterUpdateMsg.Adapter should not be nil")
	}
	if msg.Adapter.Address != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("AdapterUpdateMsg.Adapter.Address = %v, want AA:BB:CC:DD:EE:FF", msg.Adapter.Address)
	}
}

func TestAdapterPropertyChangedMsg(t *testing.T) {
	tests := []struct {
		name     string
		property string
		success  bool
		err      error
	}{
		{
			name:     "powered property changed successfully",
			property: "Powered",
			success:  true,
			err:      nil,
		},
		{
			name:     "discoverable property change failed",
			property: "Discoverable",
			success:  false,
			err:      errors.New("failed to change property"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := AdapterPropertyChangedMsg{
				Property: tt.property,
				Success:  tt.success,
				Err:      tt.err,
			}

			if msg.Property != tt.property {
				t.Errorf("AdapterPropertyChangedMsg.Property = %v, want %v", msg.Property, tt.property)
			}
			if msg.Success != tt.success {
				t.Errorf("AdapterPropertyChangedMsg.Success = %v, want %v", msg.Success, tt.success)
			}
			if (msg.Err != nil) != (tt.err != nil) {
				t.Errorf("AdapterPropertyChangedMsg.Err presence mismatch")
			}
		})
	}
}

func TestTickMsg(t *testing.T) {
	now := time.Now()
	msg := TickMsg(now)

	tickTime := time.Time(msg)
	if !tickTime.Equal(now) {
		t.Errorf("TickMsg time = %v, want %v", tickTime, now)
	}
}
