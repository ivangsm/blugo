package ui

import (
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/models"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.devices == nil {
		t.Errorf("NewModel() should initialize devices map")
	}
	if m.deviceOrder == nil {
		t.Errorf("NewModel() should initialize deviceOrder slice")
	}
	if m.focusSection != "found" {
		t.Errorf("NewModel() focusSection = %v, want 'found'", m.focusSection)
	}
	if m.selectedIndex != 0 {
		t.Errorf("NewModel() selectedIndex = %v, want 0", m.selectedIndex)
	}
	if len(m.devices) != 0 {
		t.Errorf("NewModel() should have empty devices map")
	}
	if len(m.deviceOrder) != 0 {
		t.Errorf("NewModel() should have empty deviceOrder slice")
	}
}

func TestModel_GetFoundDevices(t *testing.T) {
	tests := []struct {
		name            string
		devices         map[string]*models.Device
		deviceOrder     []string
		config          *config.Config
		expectedCount   int
		validateDevices func(*testing.T, []*models.Device)
	}{
		{
			name: "returns only non-connected devices",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Device 1",
					Connected: false,
					RSSI:      -50,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Device 2",
					Connected: true,
					RSSI:      -60,
				},
				"77:88:99:AA:BB:CC": {
					Address:   "77:88:99:AA:BB:CC",
					Name:      "Device 3",
					Connected: false,
					RSSI:      -70,
				},
			},
			deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66", "77:88:99:AA:BB:CC"},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -100},
			expectedCount: 2,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				for _, dev := range devices {
					if dev.Connected {
						t.Errorf("GetFoundDevices() should not return connected devices")
					}
				}
			},
		},
		{
			name: "maintains stable insertion order",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "First Device",
					Connected: false,
					Paired:    false,
					RSSI:      -30,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Second Device",
					Connected: false,
					Paired:    true,
					RSSI:      -80,
				},
			},
			deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66"},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -100},
			expectedCount: 2,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 2 {
					t.Fatalf("Expected 2 devices, got %d", len(devices))
				}
				// First device should be the one first in deviceOrder (not sorted by paired status)
				if devices[0].Address != "AA:BB:CC:DD:EE:FF" {
					t.Errorf("First device should maintain order, got %v", devices[0].Name)
				}
				// Second device should be the second in deviceOrder
				if devices[1].Address != "11:22:33:44:55:66" {
					t.Errorf("Second device should maintain order, got %v", devices[1].Name)
				}
			},
		},
		{
			name: "maintains order regardless of RSSI changes",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Weak Signal",
					Connected: false,
					Paired:    false,
					RSSI:      -80,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Strong Signal",
					Connected: false,
					Paired:    false,
					RSSI:      -30,
				},
				"77:88:99:AA:BB:CC": {
					Address:   "77:88:99:AA:BB:CC",
					Name:      "Medium Signal",
					Connected: false,
					Paired:    false,
					RSSI:      -50,
				},
			},
			deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66", "77:88:99:AA:BB:CC"},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -100},
			expectedCount: 3,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 3 {
					t.Fatalf("Expected 3 devices, got %d", len(devices))
				}
				// Should maintain insertion order, not sorted by RSSI
				if devices[0].Address != "AA:BB:CC:DD:EE:FF" {
					t.Errorf("First device should maintain order")
				}
				if devices[1].Address != "11:22:33:44:55:66" {
					t.Errorf("Second device should maintain order")
				}
				if devices[2].Address != "77:88:99:AA:BB:CC" {
					t.Errorf("Third device should maintain order")
				}
			},
		},
		{
			name: "filters unnamed devices when configured",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Named Device",
					Connected: false,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "",
					Alias:     "",
					Connected: false,
				},
				"77:88:99:AA:BB:CC": {
					Address:   "77:88:99:AA:BB:CC",
					Name:      "",
					Alias:     "Device Alias",
					Connected: false,
				},
			},
			deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66", "77:88:99:AA:BB:CC"},
			config:        &config.Config{HideUnnamedDevices: true, MinRSSIThreshold: -100},
			expectedCount: 2,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				for _, dev := range devices {
					if dev.Name == "" && dev.Alias == "" {
						t.Errorf("GetFoundDevices() should filter unnamed devices when configured")
					}
				}
			},
		},
		{
			name: "does not filter unnamed devices when not configured",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Named Device",
					Connected: false,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "",
					Alias:     "",
					Connected: false,
				},
			},
			deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66"},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -100},
			expectedCount: 2,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 2 {
					t.Errorf("GetFoundDevices() should not filter unnamed devices when HideUnnamedDevices is false")
				}
			},
		},
		{
			name: "filters by RSSI threshold",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Strong Signal",
					Connected: false,
					RSSI:      -40,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Weak Signal",
					Connected: false,
					RSSI:      -90,
				},
				"77:88:99:AA:BB:CC": {
					Address:   "77:88:99:AA:BB:CC",
					Name:      "Medium Signal",
					Connected: false,
					RSSI:      -70,
				},
			},
			deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66", "77:88:99:AA:BB:CC"},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -80},
			expectedCount: 2,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				for _, dev := range devices {
					if dev.RSSI != 0 && dev.RSSI < -80 {
						t.Errorf("GetFoundDevices() should filter devices below RSSI threshold, got device with RSSI %d", dev.RSSI)
					}
				}
			},
		},
		{
			name: "includes devices with RSSI 0 (no signal data)",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "No Signal Data",
					Connected: false,
					RSSI:      0,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Weak Signal",
					Connected: false,
					RSSI:      -90,
				},
			},
			deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66"},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -50},
			expectedCount: 1,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				// Should include device with RSSI 0
				hasZeroRSSI := false
				for _, dev := range devices {
					if dev.RSSI == 0 {
						hasZeroRSSI = true
					}
				}
				if !hasZeroRSSI {
					t.Errorf("GetFoundDevices() should include devices with RSSI 0")
				}
			},
		},
		{
			name:          "returns empty list when no devices",
			devices:       map[string]*models.Device{},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -100},
			expectedCount: 0,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 0 {
					t.Errorf("GetFoundDevices() should return empty list when no devices")
				}
			},
		},
		{
			name: "returns empty list when all devices are connected",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Device 1",
					Connected: true,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Device 2",
					Connected: true,
				},
			},
			config:        &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -100},
			expectedCount: 0,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 0 {
					t.Errorf("GetFoundDevices() should return empty list when all devices are connected")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config for test
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()
			config.Global = tt.config

			// Create model with test devices
			m := Model{devices: tt.devices, deviceOrder: tt.deviceOrder}

			// Get found devices
			foundDevices := m.GetFoundDevices()

			// Validate count
			if len(foundDevices) != tt.expectedCount {
				t.Errorf("GetFoundDevices() returned %d devices, want %d", len(foundDevices), tt.expectedCount)
			}

			// Run custom validation
			if tt.validateDevices != nil {
				tt.validateDevices(t, foundDevices)
			}
		})
	}
}

func TestModel_GetConnectedDevices(t *testing.T) {
	// Create test times with known ordering
	time1 := time.Now().Add(-3 * time.Hour)
	time2 := time.Now().Add(-2 * time.Hour)
	time3 := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name            string
		devices         map[string]*models.Device
		expectedCount   int
		validateDevices func(*testing.T, []*models.Device)
	}{
		{
			name: "returns only connected devices",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Device 1",
					Connected: true,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Device 2",
					Connected: false,
				},
				"77:88:99:AA:BB:CC": {
					Address:   "77:88:99:AA:BB:CC",
					Name:      "Device 3",
					Connected: true,
				},
			},
			expectedCount: 2,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				for _, dev := range devices {
					if !dev.Connected {
						t.Errorf("GetConnectedDevices() should only return connected devices")
					}
				}
			},
		},
		{
			name: "sorts by LastSeen (oldest first)",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Device 1",
					Connected: true,
					LastSeen:  time2,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Device 2",
					Connected: true,
					LastSeen:  time1,
				},
				"77:88:99:AA:BB:CC": {
					Address:   "77:88:99:AA:BB:CC",
					Name:      "Device 3",
					Connected: true,
					LastSeen:  time3,
				},
			},
			expectedCount: 3,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 3 {
					t.Fatalf("Expected 3 devices, got %d", len(devices))
				}
				// Should be sorted by LastSeen (oldest first)
				if !devices[0].LastSeen.Equal(time1) {
					t.Errorf("First device should have oldest LastSeen")
				}
				if !devices[1].LastSeen.Equal(time2) {
					t.Errorf("Second device should have middle LastSeen")
				}
				if !devices[2].LastSeen.Equal(time3) {
					t.Errorf("Third device should have newest LastSeen")
				}
			},
		},
		{
			name:          "returns empty list when no devices",
			devices:       map[string]*models.Device{},
			expectedCount: 0,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 0 {
					t.Errorf("GetConnectedDevices() should return empty list when no devices")
				}
			},
		},
		{
			name: "returns empty list when no connected devices",
			devices: map[string]*models.Device{
				"AA:BB:CC:DD:EE:FF": {
					Address:   "AA:BB:CC:DD:EE:FF",
					Name:      "Device 1",
					Connected: false,
				},
				"11:22:33:44:55:66": {
					Address:   "11:22:33:44:55:66",
					Name:      "Device 2",
					Connected: false,
				},
			},
			expectedCount: 0,
			validateDevices: func(t *testing.T, devices []*models.Device) {
				if len(devices) != 0 {
					t.Errorf("GetConnectedDevices() should return empty list when no connected devices")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{devices: tt.devices}
			connectedDevices := m.GetConnectedDevices()

			if len(connectedDevices) != tt.expectedCount {
				t.Errorf("GetConnectedDevices() returned %d devices, want %d", len(connectedDevices), tt.expectedCount)
			}

			if tt.validateDevices != nil {
				tt.validateDevices(t, connectedDevices)
			}
		})
	}
}

func TestModel_GetSelectedDevice(t *testing.T) {
	devices := map[string]*models.Device{
		"AA:BB:CC:DD:EE:FF": {
			Address:   "AA:BB:CC:DD:EE:FF",
			Name:      "Available Device 1",
			Connected: false,
			RSSI:      -50,
		},
		"11:22:33:44:55:66": {
			Address:   "11:22:33:44:55:66",
			Name:      "Available Device 2",
			Connected: false,
			RSSI:      -60,
		},
		"77:88:99:AA:BB:CC": {
			Address:   "77:88:99:AA:BB:CC",
			Name:      "Connected Device",
			Connected: true,
			LastSeen:  time.Now(),
		},
	}

	tests := []struct {
		name          string
		focusSection  string
		selectedIndex int
		expectedAddr  string // Empty string means expecting nil
	}{
		{
			name:          "returns first available device when focus on found",
			focusSection:  "found",
			selectedIndex: 0,
			expectedAddr:  "AA:BB:CC:DD:EE:FF",
		},
		{
			name:          "returns second available device when focus on found",
			focusSection:  "found",
			selectedIndex: 1,
			expectedAddr:  "11:22:33:44:55:66",
		},
		{
			name:          "returns connected device when focus on connected",
			focusSection:  "connected",
			selectedIndex: 0,
			expectedAddr:  "77:88:99:AA:BB:CC",
		},
		{
			name:          "returns nil when index out of bounds for found",
			focusSection:  "found",
			selectedIndex: 10,
			expectedAddr:  "",
		},
		{
			name:          "returns nil when index out of bounds for connected",
			focusSection:  "connected",
			selectedIndex: 10,
			expectedAddr:  "",
		},
	}

	// Set default config
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{HideUnnamedDevices: false, MinRSSIThreshold: -100}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				devices:       devices,
				deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66", "77:88:99:AA:BB:CC"},
				focusSection:  tt.focusSection,
				selectedIndex: tt.selectedIndex,
			}

			selected := m.GetSelectedDevice()

			if tt.expectedAddr == "" {
				if selected != nil {
					t.Errorf("GetSelectedDevice() = %v, want nil", selected.Address)
				}
			} else {
				if selected == nil {
					t.Fatalf("GetSelectedDevice() = nil, want device with address %s", tt.expectedAddr)
				}
				if selected.Address != tt.expectedAddr {
					t.Errorf("GetSelectedDevice().Address = %v, want %v", selected.Address, tt.expectedAddr)
				}
			}
		})
	}
}

func TestModel_GetSelectedDevice_EmptyLists(t *testing.T) {
	m := Model{
		devices:       map[string]*models.Device{},
		deviceOrder:   []string{},
		focusSection:  "found",
		selectedIndex: 0,
	}

	selected := m.GetSelectedDevice()
	if selected != nil {
		t.Errorf("GetSelectedDevice() should return nil when device list is empty")
	}
}

// TestModel_Integration tests the full workflow of device filtering and selection
func TestModel_Integration(t *testing.T) {
	// Setup devices
	devices := map[string]*models.Device{
		"AA:BB:CC:DD:EE:FF": {
			Path:      dbus.ObjectPath("/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF"),
			Address:   "AA:BB:CC:DD:EE:FF",
			Name:      "My Headphones",
			Paired:    true,
			Connected: false,
			RSSI:      -40,
			Icon:      "audio-headset",
		},
		"11:22:33:44:55:66": {
			Path:      dbus.ObjectPath("/org/bluez/hci0/dev_11_22_33_44_55_66"),
			Address:   "11:22:33:44:55:66",
			Name:      "Unknown Device",
			Paired:    false,
			Connected: false,
			RSSI:      -70,
			Icon:      "",
		},
		"77:88:99:AA:BB:CC": {
			Path:      dbus.ObjectPath("/org/bluez/hci0/dev_77_88_99_AA_BB_CC"),
			Address:   "77:88:99:AA:BB:CC",
			Name:      "My Keyboard",
			Paired:    true,
			Connected: true,
			RSSI:      -50,
			Icon:      "input-keyboard",
			LastSeen:  time.Now(),
		},
	}

	// Set config
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{
		HideUnnamedDevices: false,
		MinRSSIThreshold:   -100,
	}

	m := Model{
		devices:       devices,
		deviceOrder:   []string{"AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66", "77:88:99:AA:BB:CC"},
		focusSection:  "found",
		selectedIndex: 0,
	}

	// Test found devices
	foundDevices := m.GetFoundDevices()
	if len(foundDevices) != 2 {
		t.Errorf("Should have 2 available devices, got %d", len(foundDevices))
	}

	// First should be paired device (sorted by paired status first)
	if !foundDevices[0].Paired {
		t.Errorf("First available device should be paired")
	}

	// Test connected devices
	connectedDevices := m.GetConnectedDevices()
	if len(connectedDevices) != 1 {
		t.Errorf("Should have 1 connected device, got %d", len(connectedDevices))
	}
	if connectedDevices[0].Address != "77:88:99:AA:BB:CC" {
		t.Errorf("Connected device should be keyboard")
	}

	// Test selection
	selected := m.GetSelectedDevice()
	if selected == nil {
		t.Fatal("GetSelectedDevice() should not return nil")
	}
	if !selected.Paired {
		t.Errorf("Selected device should be the paired one")
	}
}
