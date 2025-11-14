package models

import (
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/config"
)

func TestDevice_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		device   Device
		expected string
	}{
		{
			name: "returns Name when all fields are present",
			device: Device{
				Name:    "My Device",
				Alias:   "Device Alias",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: "My Device",
		},
		{
			name: "returns Alias when Name is empty",
			device: Device{
				Name:    "",
				Alias:   "Device Alias",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: "Device Alias",
		},
		{
			name: "returns Address when Name and Alias are empty",
			device: Device{
				Name:    "",
				Alias:   "",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: "AA:BB:CC:DD:EE:FF",
		},
		{
			name: "returns empty Address when all are empty",
			device: Device{
				Name:    "",
				Alias:   "",
				Address: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.device.GetDisplayName()
			if got != tt.expected {
				t.Errorf("GetDisplayName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDevice_IsAvailable(t *testing.T) {
	tests := []struct {
		name      string
		connected bool
		expected  bool
	}{
		{
			name:      "returns true when device is not connected",
			connected: false,
			expected:  true,
		},
		{
			name:      "returns false when device is connected",
			connected: true,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := Device{Connected: tt.connected}
			got := device.IsAvailable()
			if got != tt.expected {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDevice_GetIcon(t *testing.T) {
	// Setup: Enable emojis for testing
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: true}

	tests := []struct {
		name     string
		device   Device
		expected string
	}{
		// Icon-based detection
		{
			name:     "returns headphones emoji for audio-card",
			device:   Device{Icon: "audio-card"},
			expected: "🎧",
		},
		{
			name:     "returns headphones emoji for audio-headset",
			device:   Device{Icon: "audio-headset"},
			expected: "🎧",
		},
		{
			name:     "returns headphones emoji for audio-headphones",
			device:   Device{Icon: "audio-headphones"},
			expected: "🎧",
		},
		{
			name:     "returns phone emoji for phone",
			device:   Device{Icon: "phone"},
			expected: "📱",
		},
		{
			name:     "returns phone emoji for smartphone",
			device:   Device{Icon: "smartphone"},
			expected: "📱",
		},
		{
			name:     "returns computer emoji for computer",
			device:   Device{Icon: "computer"},
			expected: "💻",
		},
		{
			name:     "returns computer emoji for laptop",
			device:   Device{Icon: "laptop"},
			expected: "💻",
		},
		{
			name:     "returns keyboard emoji for input-keyboard",
			device:   Device{Icon: "input-keyboard"},
			expected: "⌨️",
		},
		{
			name:     "returns mouse emoji for input-mouse",
			device:   Device{Icon: "input-mouse"},
			expected: "🖱️",
		},
		{
			name:     "returns gamepad emoji for input-gaming",
			device:   Device{Icon: "input-gaming"},
			expected: "🎮",
		},
		{
			name:     "returns camera emoji for camera",
			device:   Device{Icon: "camera"},
			expected: "📷",
		},
		{
			name:     "returns printer emoji for printer",
			device:   Device{Icon: "printer"},
			expected: "🖨️",
		},
		// Class-based detection (majorClass in bits 12-8)
		{
			name:     "returns computer emoji for class 1 (Computer)",
			device:   Device{Icon: "", Class: 0x0100}, // majorClass = 1
			expected: "💻",
		},
		{
			name:     "returns phone emoji for class 2 (Phone)",
			device:   Device{Icon: "", Class: 0x0200}, // majorClass = 2
			expected: "📱",
		},
		{
			name:     "returns headphones emoji for class 4 (Audio/Video)",
			device:   Device{Icon: "", Class: 0x0400}, // majorClass = 4
			expected: "🎧",
		},
		{
			name:     "returns keyboard emoji for class 5 (Peripheral)",
			device:   Device{Icon: "", Class: 0x0500}, // majorClass = 5
			expected: "⌨️",
		},
		{
			name:     "returns camera emoji for class 6 (Imaging)",
			device:   Device{Icon: "", Class: 0x0600}, // majorClass = 6
			expected: "📷",
		},
		{
			name:     "returns default signal emoji for unknown class",
			device:   Device{Icon: "", Class: 0x0700}, // majorClass = 7 (unknown)
			expected: "📶",
		},
		{
			name:     "returns default signal emoji for no icon and no class",
			device:   Device{Icon: "", Class: 0},
			expected: "📶",
		},
		{
			name:     "prioritizes icon over class",
			device:   Device{Icon: "phone", Class: 0x0100}, // icon says phone, class says computer
			expected: "📱",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.device.GetIcon()
			if got != tt.expected {
				t.Errorf("GetIcon() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDevice_GetIcon_WithEmojisDisabled(t *testing.T) {
	// Setup: Disable emojis
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: false}

	tests := []struct {
		name   string
		device Device
	}{
		{
			name:   "returns empty string when emojis disabled with icon",
			device: Device{Icon: "audio-headset"},
		},
		{
			name:   "returns empty string when emojis disabled with class",
			device: Device{Class: 0x0400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.device.GetIcon()
			if got != "" {
				t.Errorf("GetIcon() with emojis disabled = %v, want empty string", got)
			}
		})
	}
}

func TestDevice_GetBatteryInfo(t *testing.T) {
	// Setup: Enable emojis for testing
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: true}

	// Helper function to create uint8 pointer
	uint8Ptr := func(v uint8) *uint8 {
		return &v
	}

	tests := []struct {
		name         string
		battery      *uint8
		expectedIcon string
		expectedText string
	}{
		{
			name:         "returns empty strings when battery is nil",
			battery:      nil,
			expectedIcon: "",
			expectedText: "",
		},
		{
			name:         "returns full battery icon for 100%",
			battery:      uint8Ptr(100),
			expectedIcon: "🔋",
			expectedText: "100%",
		},
		{
			name:         "returns full battery icon for 90%",
			battery:      uint8Ptr(90),
			expectedIcon: "🔋",
			expectedText: "90%",
		},
		{
			name:         "returns high battery icon for 89%",
			battery:      uint8Ptr(89),
			expectedIcon: "🔋",
			expectedText: "89%",
		},
		{
			name:         "returns high battery icon for 60%",
			battery:      uint8Ptr(60),
			expectedIcon: "🔋",
			expectedText: "60%",
		},
		{
			name:         "returns medium battery icon for 59%",
			battery:      uint8Ptr(59),
			expectedIcon: "🔋",
			expectedText: "59%",
		},
		{
			name:         "returns medium battery icon for 30%",
			battery:      uint8Ptr(30),
			expectedIcon: "🔋",
			expectedText: "30%",
		},
		{
			name:         "returns low battery icon for 29%",
			battery:      uint8Ptr(29),
			expectedIcon: "🪫",
			expectedText: "29%",
		},
		{
			name:         "returns low battery icon for 15%",
			battery:      uint8Ptr(15),
			expectedIcon: "🪫",
			expectedText: "15%",
		},
		{
			name:         "returns critical battery icon for 14%",
			battery:      uint8Ptr(14),
			expectedIcon: "🪫",
			expectedText: "14%",
		},
		{
			name:         "returns critical battery icon for 0%",
			battery:      uint8Ptr(0),
			expectedIcon: "🪫",
			expectedText: "0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := Device{Battery: tt.battery}
			gotIcon, gotText := device.GetBatteryInfo()
			if gotIcon != tt.expectedIcon {
				t.Errorf("GetBatteryInfo() icon = %v, want %v", gotIcon, tt.expectedIcon)
			}
			if gotText != tt.expectedText {
				t.Errorf("GetBatteryInfo() text = %v, want %v", gotText, tt.expectedText)
			}
		})
	}
}

func TestDevice_GetBatteryInfo_WithEmojisDisabled(t *testing.T) {
	// Setup: Disable emojis
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: false}

	battery := uint8(50)
	device := Device{Battery: &battery}
	gotIcon, gotText := device.GetBatteryInfo()

	if gotIcon != "" {
		t.Errorf("GetBatteryInfo() icon with emojis disabled = %v, want empty string", gotIcon)
	}
	if gotText != "50%" {
		t.Errorf("GetBatteryInfo() text = %v, want 50%%", gotText)
	}
}

func TestDevice_HasBattery(t *testing.T) {
	tests := []struct {
		name     string
		battery  *uint8
		expected bool
	}{
		{
			name:     "returns false when battery is nil",
			battery:  nil,
			expected: false,
		},
		{
			name: "returns true when battery is not nil",
			battery: func() *uint8 {
				v := uint8(50)
				return &v
			}(),
			expected: true,
		},
		{
			name: "returns true when battery is 0",
			battery: func() *uint8 {
				v := uint8(0)
				return &v
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := Device{Battery: tt.battery}
			got := device.HasBattery()
			if got != tt.expected {
				t.Errorf("HasBattery() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEmoji(t *testing.T) {
	tests := []struct {
		name       string
		showEmojis bool
		input      string
		expected   string
	}{
		{
			name:       "returns emoji when ShowEmojis is true",
			showEmojis: true,
			input:      "🎧",
			expected:   "🎧",
		},
		{
			name:       "returns empty string when ShowEmojis is false",
			showEmojis: false,
			input:      "🎧",
			expected:   "",
		},
		{
			name:       "returns empty string when config is nil",
			showEmojis: false, // will set config to nil
			input:      "🎧",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()

			// Set config based on test case
			if tt.name == "returns empty string when config is nil" {
				config.Global = nil
			} else {
				config.Global = &config.Config{ShowEmojis: tt.showEmojis}
			}

			got := emoji(tt.input)
			if got != tt.expected {
				t.Errorf("emoji() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestDevice_Integration tests the Device struct with realistic data
func TestDevice_Integration(t *testing.T) {
	battery := uint8(75)
	device := Device{
		Path:      dbus.ObjectPath("/org/bluez/hci0/dev_AA_BB_CC_DD_EE_FF"),
		Address:   "AA:BB:CC:DD:EE:FF",
		Name:      "My Headphones",
		Alias:     "Custom Name",
		Paired:    true,
		Trusted:   true,
		Connected: true,
		RSSI:      -50,
		Icon:      "audio-headset",
		Class:     0x0400,
		Battery:   &battery,
		LastSeen:  time.Now(),
	}

	// Enable emojis
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: true}

	// Test all methods together
	if device.GetDisplayName() != "My Headphones" {
		t.Errorf("GetDisplayName() failed for integrated test")
	}

	if device.IsAvailable() != false {
		t.Errorf("IsAvailable() should return false for connected device")
	}

	if device.GetIcon() != "🎧" {
		t.Errorf("GetIcon() failed for integrated test")
	}

	icon, text := device.GetBatteryInfo()
	if icon != "🔋" || text != "75%" {
		t.Errorf("GetBatteryInfo() failed for integrated test: got (%v, %v), want (🔋, 75%%)", icon, text)
	}

	if !device.HasBattery() {
		t.Errorf("HasBattery() should return true for device with battery")
	}
}

func TestNormalizeMAC(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normalizes MAC with colons",
			input:    "AA:BB:CC:DD:EE:FF",
			expected: "aabbccddeeff",
		},
		{
			name:     "normalizes MAC with dashes",
			input:    "AA-BB-CC-DD-EE-FF",
			expected: "aabbccddeeff",
		},
		{
			name:     "normalizes MAC with mixed case",
			input:    "Aa:Bb:Cc:Dd:Ee:Ff",
			expected: "aabbccddeeff",
		},
		{
			name:     "normalizes already normalized MAC",
			input:    "aabbccddeeff",
			expected: "aabbccddeeff",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeMAC(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeMAC(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsAliasMACAddress(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		address  string
		expected bool
	}{
		{
			name:     "detects MAC with dashes as same as colons",
			alias:    "AA-BB-CC-DD-EE-FF",
			address:  "AA:BB:CC:DD:EE:FF",
			expected: true,
		},
		{
			name:     "detects lowercase MAC with dashes",
			alias:    "aa-bb-cc-dd-ee-ff",
			address:  "AA:BB:CC:DD:EE:FF",
			expected: true,
		},
		{
			name:     "detects different MACs",
			alias:    "11-22-33-44-55-66",
			address:  "AA:BB:CC:DD:EE:FF",
			expected: false,
		},
		{
			name:     "detects real device name vs MAC",
			alias:    "My Bluetooth Device",
			address:  "AA:BB:CC:DD:EE:FF",
			expected: false,
		},
		{
			name:     "handles identical format",
			alias:    "AA:BB:CC:DD:EE:FF",
			address:  "AA:BB:CC:DD:EE:FF",
			expected: true,
		},
		{
			name:     "handles empty alias",
			alias:    "",
			address:  "AA:BB:CC:DD:EE:FF",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAliasMACAddress(tt.alias, tt.address)
			if result != tt.expected {
				t.Errorf("IsAliasMACAddress(%q, %q) = %v, want %v", tt.alias, tt.address, result, tt.expected)
			}
		})
	}
}

func TestDevice_HasRealName(t *testing.T) {
	tests := []struct {
		name     string
		device   Device
		expected bool
	}{
		{
			name: "returns true when Name is set",
			device: Device{
				Name:    "My Device",
				Alias:   "",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: true,
		},
		{
			name: "returns true when Alias is real name",
			device: Device{
				Name:    "",
				Alias:   "Real Device Name",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: true,
		},
		{
			name: "returns false when Alias is MAC with dashes",
			device: Device{
				Name:    "",
				Alias:   "AA-BB-CC-DD-EE-FF",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: false,
		},
		{
			name: "returns false when Alias is MAC with colons",
			device: Device{
				Name:    "",
				Alias:   "AA:BB:CC:DD:EE:FF",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: false,
		},
		{
			name: "returns false when both Name and Alias are empty",
			device: Device{
				Name:    "",
				Alias:   "",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: false,
		},
		{
			name: "returns true when Name is set even if Alias is MAC",
			device: Device{
				Name:    "Real Name",
				Alias:   "AA-BB-CC-DD-EE-FF",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.device.HasRealName()
			if result != tt.expected {
				t.Errorf("HasRealName() = %v, want %v", result, tt.expected)
			}
		})
	}
}
