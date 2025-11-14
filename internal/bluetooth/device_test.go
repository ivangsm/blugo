package bluetooth

import (
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/models"
)

func TestParseDevice(t *testing.T) {
	// Helper to create dbus.Variant
	makeVariant := func(v interface{}) dbus.Variant {
		return dbus.MakeVariant(v)
	}

	tests := []struct {
		name       string
		path       dbus.ObjectPath
		interfaces map[string]map[string]dbus.Variant
		props      map[string]dbus.Variant
		validate   func(*testing.T, *models.Device)
	}{
		{
			name: "parses all basic device properties",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address":   makeVariant("AA:BB:CC:DD:EE:FF"),
					"Name":      makeVariant("Test Device"),
					"Alias":     makeVariant("My Device"),
					"Paired":    makeVariant(true),
					"Trusted":   makeVariant(true),
					"Connected": makeVariant(true),
					"RSSI":      makeVariant(int16(-50)),
					"Icon":      makeVariant("audio-headset"),
					"Class":     makeVariant(uint32(0x0400)),
				},
			},
			props: map[string]dbus.Variant{
				"Address":   makeVariant("AA:BB:CC:DD:EE:FF"),
				"Name":      makeVariant("Test Device"),
				"Alias":     makeVariant("My Device"),
				"Paired":    makeVariant(true),
				"Trusted":   makeVariant(true),
				"Connected": makeVariant(true),
				"RSSI":      makeVariant(int16(-50)),
				"Icon":      makeVariant("audio-headset"),
				"Class":     makeVariant(uint32(0x0400)),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Path != "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF" {
					t.Errorf("Path = %v, want /org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF", dev.Path)
				}
				if dev.Address != "AA:BB:CC:DD:EE:FF" {
					t.Errorf("Address = %v, want AA:BB:CC:DD:EE:FF", dev.Address)
				}
				if dev.Name != "Test Device" {
					t.Errorf("Name = %v, want Test Device", dev.Name)
				}
				if dev.Alias != "My Device" {
					t.Errorf("Alias = %v, want My Device", dev.Alias)
				}
				if !dev.Paired {
					t.Errorf("Paired = false, want true")
				}
				if !dev.Trusted {
					t.Errorf("Trusted = false, want true")
				}
				if !dev.Connected {
					t.Errorf("Connected = false, want true")
				}
				if dev.RSSI != -50 {
					t.Errorf("RSSI = %v, want -50", dev.RSSI)
				}
				if dev.Icon != "audio-headset" {
					t.Errorf("Icon = %v, want audio-headset", dev.Icon)
				}
				if dev.Class != 0x0400 {
					t.Errorf("Class = %v, want 0x0400", dev.Class)
				}
			},
		},
		{
			name: "uses Alias as Name when Name is empty",
			path: "/org/bluez/hci0/dev_11_22_33_44_55_66",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("11:22:33:44:55:66"),
					"Name":    makeVariant(""),
					"Alias":   makeVariant("Device Alias"),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("11:22:33:44:55:66"),
				"Name":    makeVariant(""),
				"Alias":   makeVariant("Device Alias"),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Name != "Device Alias" {
					t.Errorf("Name = %v, want Device Alias (should use Alias when Name is empty)", dev.Name)
				}
				if dev.Alias != "Device Alias" {
					t.Errorf("Alias = %v, want Device Alias", dev.Alias)
				}
			},
		},
		{
			name: "handles device with battery information",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
					"Name":    makeVariant("Headphones"),
				},
				bluezBatteryIface: {
					"Percentage": makeVariant(byte(75)),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				"Name":    makeVariant("Headphones"),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Battery == nil {
					t.Errorf("Battery should not be nil")
				} else if *dev.Battery != 75 {
					t.Errorf("Battery = %v, want 75", *dev.Battery)
				}
			},
		},
		{
			name: "handles device without battery information",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
					"Name":    makeVariant("Mouse"),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				"Name":    makeVariant("Mouse"),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Battery != nil {
					t.Errorf("Battery should be nil for device without battery interface")
				}
			},
		},
		{
			name: "handles battery with 0% charge",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				},
				bluezBatteryIface: {
					"Percentage": makeVariant(byte(0)),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Battery == nil {
					t.Errorf("Battery should not be nil even with 0%%")
				} else if *dev.Battery != 0 {
					t.Errorf("Battery = %v, want 0", *dev.Battery)
				}
			},
		},
		{
			name: "handles battery with 100% charge",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				},
				bluezBatteryIface: {
					"Percentage": makeVariant(byte(100)),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Battery == nil {
					t.Errorf("Battery should not be nil")
				} else if *dev.Battery != 100 {
					t.Errorf("Battery = %v, want 100", *dev.Battery)
				}
			},
		},
		{
			name: "handles missing optional properties",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Address != "AA:BB:CC:DD:EE:FF" {
					t.Errorf("Address = %v, want AA:BB:CC:DD:EE:FF", dev.Address)
				}
				if dev.Name != "" {
					t.Errorf("Name should be empty when not provided")
				}
				if dev.Alias != "" {
					t.Errorf("Alias should be empty when not provided")
				}
				if dev.Paired {
					t.Errorf("Paired should be false when not provided")
				}
				if dev.Trusted {
					t.Errorf("Trusted should be false when not provided")
				}
				if dev.Connected {
					t.Errorf("Connected should be false when not provided")
				}
				if dev.RSSI != 0 {
					t.Errorf("RSSI should be 0 when not provided")
				}
				if dev.Icon != "" {
					t.Errorf("Icon should be empty when not provided")
				}
				if dev.Class != 0 {
					t.Errorf("Class should be 0 when not provided")
				}
			},
		},
		{
			name: "handles negative RSSI values",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
					"RSSI":    makeVariant(int16(-100)),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				"RSSI":    makeVariant(int16(-100)),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.RSSI != -100 {
					t.Errorf("RSSI = %v, want -100", dev.RSSI)
				}
			},
		},
		{
			name: "handles various device classes",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
					"Class":   makeVariant(uint32(0x240404)), // Audio device class
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				"Class":   makeVariant(uint32(0x240404)),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.Class != 0x240404 {
					t.Errorf("Class = %v, want 0x240404", dev.Class)
				}
			},
		},
		{
			name: "sets LastSeen timestamp",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
				},
			},
			props: map[string]dbus.Variant{
				"Address": makeVariant("AA:BB:CC:DD:EE:FF"),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if dev.LastSeen.IsZero() {
					t.Errorf("LastSeen should be set to current time")
				}
				// Check that LastSeen is recent (within last second)
				if time.Since(dev.LastSeen) > time.Second {
					t.Errorf("LastSeen should be recent, got %v", dev.LastSeen)
				}
			},
		},
		{
			name: "handles disconnected paired device",
			path: "/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			interfaces: map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Address":   makeVariant("AA:BB:CC:DD:EE:FF"),
					"Name":      makeVariant("My Headphones"),
					"Paired":    makeVariant(true),
					"Connected": makeVariant(false),
				},
			},
			props: map[string]dbus.Variant{
				"Address":   makeVariant("AA:BB:CC:DD:EE:FF"),
				"Name":      makeVariant("My Headphones"),
				"Paired":    makeVariant(true),
				"Connected": makeVariant(false),
			},
			validate: func(t *testing.T, dev *models.Device) {
				if !dev.Paired {
					t.Errorf("Paired should be true")
				}
				if dev.Connected {
					t.Errorf("Connected should be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := parseDevice(tt.path, tt.interfaces, tt.props)
			if dev == nil {
				t.Fatal("parseDevice returned nil")
			}
			tt.validate(t, dev)
		})
	}
}

// TestParseDevice_EdgeCases tests edge cases and error conditions
func TestParseDevice_EdgeCases(t *testing.T) {
	makeVariant := func(v interface{}) dbus.Variant {
		return dbus.MakeVariant(v)
	}

	t.Run("handles empty props map", func(t *testing.T) {
		dev := parseDevice(
			"/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			map[string]map[string]dbus.Variant{},
			map[string]dbus.Variant{},
		)
		if dev == nil {
			t.Fatal("parseDevice should not return nil even with empty props")
		}
		if dev.Address != "" {
			t.Errorf("Address should be empty")
		}
	})

	t.Run("handles Name without Alias", func(t *testing.T) {
		dev := parseDevice(
			"/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Name": makeVariant("Device Name"),
				},
			},
			map[string]dbus.Variant{
				"Name": makeVariant("Device Name"),
			},
		)
		if dev.Name != "Device Name" {
			t.Errorf("Name = %v, want Device Name", dev.Name)
		}
		if dev.Alias != "" {
			t.Errorf("Alias should be empty")
		}
	})

	t.Run("preserves Name when both Name and Alias exist", func(t *testing.T) {
		dev := parseDevice(
			"/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
			map[string]map[string]dbus.Variant{
				bluezDeviceIface: {
					"Name":  makeVariant("Original Name"),
					"Alias": makeVariant("Device Alias"),
				},
			},
			map[string]dbus.Variant{
				"Name":  makeVariant("Original Name"),
				"Alias": makeVariant("Device Alias"),
			},
		)
		if dev.Name != "Original Name" {
			t.Errorf("Name = %v, want Original Name (should not be overwritten by Alias)", dev.Name)
		}
	})
}


func TestParseDevice_DoesNotUseAliasWhenItIsMAC(t *testing.T) {
	makeVariant := func(v interface{}) dbus.Variant {
		return dbus.MakeVariant(v)
	}

	tests := []struct {
		name    string
		address string
		alias   string
	}{
		{
			name:    "does not use Alias when it matches Address with dashes",
			address: "AA:BB:CC:DD:EE:FF",
			alias:   "AA-BB-CC-DD-EE-FF",
		},
		{
			name:    "does not use Alias when it matches Address exactly",
			address: "AA:BB:CC:DD:EE:FF",
			alias:   "AA:BB:CC:DD:EE:FF",
		},
		{
			name:    "does not use lowercase Alias matching Address",
			address: "AA:BB:CC:DD:EE:FF",
			alias:   "aa-bb-cc-dd-ee-ff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := parseDevice(
				"/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF",
				map[string]map[string]dbus.Variant{
					bluezDeviceIface: {
						"Address": makeVariant(tt.address),
						"Name":    makeVariant(""),
						"Alias":   makeVariant(tt.alias),
					},
				},
				map[string]dbus.Variant{
					"Address": makeVariant(tt.address),
					"Name":    makeVariant(""),
					"Alias":   makeVariant(tt.alias),
				},
			)

			if dev.Name != "" {
				t.Errorf("Name should be empty when Alias is MAC address, got %q", dev.Name)
			}
			if dev.Alias != tt.alias {
				t.Errorf("Alias = %v, want %v", dev.Alias, tt.alias)
			}
		})
	}
}
