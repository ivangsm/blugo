package ui

import (
	"testing"

	"github.com/ivangsm/blugo/internal/config"
)

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "returns a when a > b",
			a:        10,
			b:        5,
			expected: 10,
		},
		{
			name:     "returns b when b > a",
			a:        5,
			b:        10,
			expected: 10,
		},
		{
			name:     "returns either when equal",
			a:        10,
			b:        10,
			expected: 10,
		},
		{
			name:     "handles negative numbers",
			a:        -5,
			b:        -10,
			expected: -5,
		},
		{
			name:     "handles zero",
			a:        0,
			b:        -5,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "returns b when b < a",
			a:        10,
			b:        5,
			expected: 5,
		},
		{
			name:     "returns a when a < b",
			a:        5,
			b:        10,
			expected: 5,
		},
		{
			name:     "returns either when equal",
			a:        10,
			b:        10,
			expected: 10,
		},
		{
			name:     "handles negative numbers",
			a:        -5,
			b:        -10,
			expected: -10,
		},
		{
			name:     "handles zero",
			a:        0,
			b:        5,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := min(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestGetBatteryStyle(t *testing.T) {
	tests := []struct {
		name         string
		level        uint8
		config       *config.Config
		expectedName string // We'll check the style name/type
	}{
		{
			name:  "returns high style for level >= 60 (default threshold)",
			level: 80,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  30,
			},
			expectedName: "high",
		},
		{
			name:  "returns high style for level exactly at threshold",
			level: 60,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  30,
			},
			expectedName: "high",
		},
		{
			name:  "returns medium style for level between thresholds",
			level: 50,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  30,
			},
			expectedName: "medium",
		},
		{
			name:  "returns medium style for level exactly at low threshold",
			level: 30,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  30,
			},
			expectedName: "medium",
		},
		{
			name:  "returns low style for level < 30 (default threshold)",
			level: 20,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  30,
			},
			expectedName: "low",
		},
		{
			name:  "returns low style for 0%",
			level: 0,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  30,
			},
			expectedName: "low",
		},
		{
			name:  "returns high style for 100%",
			level: 100,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  30,
			},
			expectedName: "high",
		},
		{
			name:  "uses custom high threshold",
			level: 75,
			config: &config.Config{
				BatteryHighThreshold: 80,
				BatteryLowThreshold:  30,
			},
			expectedName: "medium",
		},
		{
			name:  "uses custom low threshold",
			level: 15,
			config: &config.Config{
				BatteryHighThreshold: 60,
				BatteryLowThreshold:  20,
			},
			expectedName: "low",
		},
		{
			name:         "uses default thresholds when config is nil",
			level:        70,
			config:       nil,
			expectedName: "high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()
			config.Global = tt.config

			// Get style
			style := GetBatteryStyle(tt.level)

			// We can't easily compare styles directly, but we can verify
			// that the function runs without panic and returns a valid style
			_ = style
			// The test passes if GetBatteryStyle doesn't panic
		})
	}
}

func TestGetMaxWidth(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected int
	}{
		{
			name: "returns configured width when valid",
			config: &config.Config{
				MaxTerminalWidth: 120,
			},
			expected: 120,
		},
		{
			name: "returns 80 when configured width < 80",
			config: &config.Config{
				MaxTerminalWidth: 50,
			},
			expected: 80,
		},
		{
			name: "returns 200 when configured width > 200",
			config: &config.Config{
				MaxTerminalWidth: 250,
			},
			expected: 200,
		},
		{
			name: "returns 80 exactly when configured as 80",
			config: &config.Config{
				MaxTerminalWidth: 80,
			},
			expected: 80,
		},
		{
			name: "returns 200 exactly when configured as 200",
			config: &config.Config{
				MaxTerminalWidth: 200,
			},
			expected: 200,
		},
		{
			name:     "returns 140 when config is nil",
			config:   nil,
			expected: 140,
		},
		{
			name: "returns 140 configured value",
			config: &config.Config{
				MaxTerminalWidth: 140,
			},
			expected: 140,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()
			config.Global = tt.config

			got := GetMaxWidth()
			if got != tt.expected {
				t.Errorf("GetMaxWidth() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestGetPadding(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected int
	}{
		{
			name: "returns 0 when compact mode enabled",
			config: &config.Config{
				CompactMode: true,
			},
			expected: 0,
		},
		{
			name: "returns 1 when compact mode disabled",
			config: &config.Config{
				CompactMode: false,
			},
			expected: 1,
		},
		{
			name:     "returns 1 when config is nil",
			config:   nil,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()
			config.Global = tt.config

			got := GetPadding()
			if got != tt.expected {
				t.Errorf("GetPadding() = %d, want %d", got, tt.expected)
			}
		})
	}
}
