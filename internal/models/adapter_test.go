package models

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/config"
)

func TestAdapter_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		adapter  Adapter
		expected string
	}{
		{
			name: "returns Alias when all fields are present",
			adapter: Adapter{
				Alias:   "My Adapter",
				Name:    "hci0",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: "My Adapter",
		},
		{
			name: "returns Name when Alias is empty",
			adapter: Adapter{
				Alias:   "",
				Name:    "hci0",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: "hci0",
		},
		{
			name: "returns Address when Alias and Name are empty",
			adapter: Adapter{
				Alias:   "",
				Name:    "",
				Address: "AA:BB:CC:DD:EE:FF",
			},
			expected: "AA:BB:CC:DD:EE:FF",
		},
		{
			name: "returns empty Address when all fields are empty",
			adapter: Adapter{
				Alias:   "",
				Name:    "",
				Address: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.adapter.GetDisplayName()
			if got != tt.expected {
				t.Errorf("GetDisplayName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAdapter_GetStatusIcon(t *testing.T) {
	// Setup: Enable emojis for testing
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: true}

	tests := []struct {
		name        string
		powered     bool
		discovering bool
		expected    string
	}{
		{
			name:        "returns black circle when powered off",
			powered:     false,
			discovering: false,
			expected:    "⚫",
		},
		{
			name:        "returns black circle when powered off even if discovering",
			powered:     false,
			discovering: true,
			expected:    "⚫",
		},
		{
			name:        "returns magnifying glass when powered on and discovering",
			powered:     true,
			discovering: true,
			expected:    "🔍",
		},
		{
			name:        "returns blue circle when powered on and not discovering",
			powered:     true,
			discovering: false,
			expected:    "🔵",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := Adapter{
				Powered:     tt.powered,
				Discovering: tt.discovering,
			}
			got := adapter.GetStatusIcon()
			if got != tt.expected {
				t.Errorf("GetStatusIcon() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAdapter_GetStatusIcon_WithEmojisDisabled(t *testing.T) {
	// Setup: Disable emojis
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: false}

	tests := []struct {
		name        string
		powered     bool
		discovering bool
	}{
		{
			name:        "returns empty string when powered off",
			powered:     false,
			discovering: false,
		},
		{
			name:        "returns empty string when powered on and discovering",
			powered:     true,
			discovering: true,
		},
		{
			name:        "returns empty string when powered on and not discovering",
			powered:     true,
			discovering: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := Adapter{
				Powered:     tt.powered,
				Discovering: tt.discovering,
			}
			got := adapter.GetStatusIcon()
			if got != "" {
				t.Errorf("GetStatusIcon() with emojis disabled = %v, want empty string", got)
			}
		})
	}
}

// TestAdapter_Integration tests the Adapter struct with realistic data
func TestAdapter_Integration(t *testing.T) {
	adapter := Adapter{
		Path:         dbus.ObjectPath("/org/bluez/hci0"),
		Address:      "AA:BB:CC:DD:EE:FF",
		Name:         "hci0",
		Alias:        "My Bluetooth Adapter",
		Powered:      true,
		Discoverable: true,
		Pairable:     true,
		Discovering:  true,
	}

	// Enable emojis
	originalConfig := config.Global
	defer func() { config.Global = originalConfig }()
	config.Global = &config.Config{ShowEmojis: true}

	// Test all methods together
	if adapter.GetDisplayName() != "My Bluetooth Adapter" {
		t.Errorf("GetDisplayName() failed for integrated test")
	}

	if adapter.GetStatusIcon() != "🔍" {
		t.Errorf("GetStatusIcon() failed for integrated test, got %v, want 🔍", adapter.GetStatusIcon())
	}

	// Test with powered off
	adapter.Powered = false
	if adapter.GetStatusIcon() != "⚫" {
		t.Errorf("GetStatusIcon() failed for powered off test, got %v, want ⚫", adapter.GetStatusIcon())
	}

	// Test with powered on but not discovering
	adapter.Powered = true
	adapter.Discovering = false
	if adapter.GetStatusIcon() != "🔵" {
		t.Errorf("GetStatusIcon() failed for powered on test, got %v, want 🔵", adapter.GetStatusIcon())
	}
}
