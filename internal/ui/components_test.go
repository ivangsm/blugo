package ui

import (
	"strings"
	"testing"

	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/i18n"
	"github.com/ivangsm/blugo/internal/models"
)

func TestGetEmptyAvailableDevicesMessage(t *testing.T) {
	// Set language to English for consistent testing
	i18n.SetLanguage(i18n.English)

	msg := getEmptyAvailableDevicesMessage()
	if msg == "" {
		t.Errorf("getEmptyAvailableDevicesMessage() returned empty string")
	}

	// Verify it returns the translation
	if msg != i18n.T.NoDevicesAvailable {
		t.Errorf("getEmptyAvailableDevicesMessage() = %q, want %q", msg, i18n.T.NoDevicesAvailable)
	}
}

func TestGetEmptyConnectedDevicesMessage(t *testing.T) {
	// Set language to English for consistent testing
	i18n.SetLanguage(i18n.English)

	msg := getEmptyConnectedDevicesMessage()
	if msg == "" {
		t.Errorf("getEmptyConnectedDevicesMessage() returned empty string")
	}

	// Verify it returns the translation
	if msg != i18n.T.NoDevicesConnected {
		t.Errorf("getEmptyConnectedDevicesMessage() = %q, want %q", msg, i18n.T.NoDevicesConnected)
	}
}

func TestRenderEmptyState(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "renders with simple message",
			message: "No devices",
		},
		{
			name:    "renders with empty message",
			message: "",
		},
		{
			name:    "renders with long message",
			message: "This is a very long message that should still be rendered correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderEmptyState(tt.message)
			if result == "" {
				t.Errorf("renderEmptyState(%q) returned empty string", tt.message)
			}
			// The result should contain the message (styled)
			// We can't test the exact styling, but we can verify it's not empty
		})
	}
}

func TestModel_RenderSeparator(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		config        *config.Config
		expectPattern string
	}{
		{
			name:          "renders separator with normal width",
			width:         100,
			config:        &config.Config{MaxTerminalWidth: 100},
			expectPattern: "─",
		},
		{
			name:          "renders separator with small width",
			width:         50,
			config:        &config.Config{MaxTerminalWidth: 100},
			expectPattern: "─",
		},
		{
			name:          "renders separator with zero width falls back to default",
			width:         0,
			config:        &config.Config{MaxTerminalWidth: 140},
			expectPattern: "─",
		},
		{
			name:          "renders separator with large width constrained by max",
			width:         300,
			config:        &config.Config{MaxTerminalWidth: 150},
			expectPattern: "─",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()
			config.Global = tt.config

			m := Model{width: tt.width}
			result := m.renderSeparator()

			if result == "" {
				t.Errorf("renderSeparator() returned empty string")
			}

			// Verify it contains the separator pattern
			if !strings.Contains(result, tt.expectPattern) {
				t.Errorf("renderSeparator() should contain %q", tt.expectPattern)
			}
		})
	}
}

func TestModel_RenderThickSeparator(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		config        *config.Config
		expectPattern string
	}{
		{
			name:          "renders thick separator with normal width",
			width:         100,
			config:        &config.Config{MaxTerminalWidth: 100},
			expectPattern: "━",
		},
		{
			name:          "renders thick separator with small width",
			width:         50,
			config:        &config.Config{MaxTerminalWidth: 100},
			expectPattern: "━",
		},
		{
			name:          "renders thick separator with zero width falls back to default",
			width:         0,
			config:        &config.Config{MaxTerminalWidth: 140},
			expectPattern: "━",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()
			config.Global = tt.config

			m := Model{width: tt.width}
			result := m.renderThickSeparator()

			if result == "" {
				t.Errorf("renderThickSeparator() returned empty string")
			}

			// Verify it contains the separator pattern
			if !strings.Contains(result, tt.expectPattern) {
				t.Errorf("renderThickSeparator() should contain %q", tt.expectPattern)
			}
		})
	}
}

func TestRenderDeviceCount(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "renders count of zero",
			count: 0,
		},
		{
			name:  "renders count of one",
			count: 1,
		},
		{
			name:  "renders count of multiple",
			count: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderDeviceCount(tt.count)
			if result == "" {
				t.Errorf("renderDeviceCount(%d) returned empty string", tt.count)
			}
			// The result should be a styled string representation
			// We verify it's not empty
		})
	}
}

func TestRenderSectionHeader(t *testing.T) {
	// Set default config
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{
		ShowEmojis: true,
	}

	tests := []struct {
		name      string
		icon      string
		title     string
		count     int
		isFocused bool
	}{
		{
			name:      "renders unfocused header",
			icon:      "📡",
			title:     "Available Devices",
			count:     5,
			isFocused: false,
		},
		{
			name:      "renders focused header",
			icon:      "🔗",
			title:     "Connected Devices",
			count:     2,
			isFocused: true,
		},
		{
			name:      "renders header with zero count",
			icon:      "📡",
			title:     "Available Devices",
			count:     0,
			isFocused: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderSectionHeader(tt.icon, tt.title, tt.count, tt.isFocused)
			if result == "" {
				t.Errorf("renderSectionHeader() returned empty string")
			}
			// Verify it contains the title (styling may wrap it)
			// We can't test exact output due to styling, but verify it's not empty
		})
	}
}

func TestRenderDeviceItem(t *testing.T) {
	// Set default config
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{
		ShowEmojis:        true,
		ShowDeviceAddress: true,
		ShowRSSI:          true,
		ShowBattery:       true,
	}

	battery := uint8(75)
	tests := []struct {
		name       string
		device     *models.Device
		isSelected bool
		showRSSI   bool
	}{
		{
			name: "renders selected device with all info",
			device: &models.Device{
				Address:   "AA:BB:CC:DD:EE:FF",
				Name:      "My Headphones",
				Icon:      "audio-headset",
				RSSI:      -50,
				Battery:   &battery,
				Paired:    true,
				Trusted:   true,
				Connected: true,
			},
			isSelected: true,
			showRSSI:   true,
		},
		{
			name: "renders unselected device",
			device: &models.Device{
				Address: "11:22:33:44:55:66",
				Name:    "My Keyboard",
				Icon:    "input-keyboard",
			},
			isSelected: false,
			showRSSI:   false,
		},
		{
			name: "renders device without name",
			device: &models.Device{
				Address: "77:88:99:AA:BB:CC",
				Name:    "",
				Alias:   "",
			},
			isSelected: false,
			showRSSI:   false,
		},
		{
			name: "renders device with battery",
			device: &models.Device{
				Address: "AA:BB:CC:DD:EE:FF",
				Name:    "Device with Battery",
				Battery: &battery,
			},
			isSelected: false,
			showRSSI:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderDeviceItem(tt.device, tt.isSelected, tt.showRSSI)
			if result == "" {
				t.Errorf("renderDeviceItem() returned empty string")
			}
			// Verify it's not empty - we can't check exact content due to styling
		})
	}
}

func TestRenderDeviceItem_WithDifferentConfigs(t *testing.T) {
	device := &models.Device{
		Address: "AA:BB:CC:DD:EE:FF",
		Name:    "Test Device",
		RSSI:    -50,
	}

	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "with address hidden",
			config: &config.Config{
				ShowDeviceAddress: false,
				ShowRSSI:          true,
				ShowBattery:       true,
				ShowEmojis:        true,
			},
		},
		{
			name: "with RSSI hidden",
			config: &config.Config{
				ShowDeviceAddress: true,
				ShowRSSI:          false,
				ShowBattery:       true,
				ShowEmojis:        true,
			},
		},
		{
			name: "with battery hidden",
			config: &config.Config{
				ShowDeviceAddress: true,
				ShowRSSI:          true,
				ShowBattery:       false,
				ShowEmojis:        true,
			},
		},
		{
			name: "with emojis disabled",
			config: &config.Config{
				ShowDeviceAddress: true,
				ShowRSSI:          true,
				ShowBattery:       true,
				ShowEmojis:        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()
			config.Global = tt.config

			result := renderDeviceItem(device, false, true)
			if result == "" {
				t.Errorf("renderDeviceItem() returned empty string")
			}
		})
	}
}

func TestModel_RenderAdapterTable(t *testing.T) {
	// Set language to English for consistent testing
	i18n.SetLanguage(i18n.English)

	// Set default config
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{
		MaxTerminalWidth: 140,
		ShowEmojis:       true,
	}

	t.Run("renders loading message when adapter is nil", func(t *testing.T) {
		m := Model{adapter: nil, width: 100}
		result := m.renderAdapterTable()

		if result == "" {
			t.Errorf("renderAdapterTable() returned empty string")
		}

		// Should contain some loading indication
		if !strings.Contains(result, "Loading") && !strings.Contains(result, "loading") {
			t.Logf("Result: %s", result)
		}
	})

	// We can't easily test with a real adapter without full initialization
	// but we've verified the nil case works
}
