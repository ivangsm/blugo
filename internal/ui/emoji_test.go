package ui

import (
	"testing"

	"github.com/ivangsm/blugo/internal/config"
)

func TestEmoji(t *testing.T) {
	tests := []struct {
		name       string
		emoji      string
		showEmojis bool
		configNil  bool
		expected   string
	}{
		{
			name:       "returns emoji when ShowEmojis is true",
			emoji:      "🎧",
			showEmojis: true,
			configNil:  false,
			expected:   "🎧",
		},
		{
			name:       "returns empty string when ShowEmojis is false",
			emoji:      "🎧",
			showEmojis: false,
			configNil:  false,
			expected:   "",
		},
		{
			name:       "returns empty string when config is nil",
			emoji:      "🎧",
			showEmojis: false,
			configNil:  true,
			expected:   "",
		},
		{
			name:       "handles empty emoji string",
			emoji:      "",
			showEmojis: true,
			configNil:  false,
			expected:   "",
		},
		{
			name:       "handles multi-character emoji",
			emoji:      "👍🏻",
			showEmojis: true,
			configNil:  false,
			expected:   "👍🏻",
		},
		{
			name:       "handles text when emojis enabled",
			emoji:      "test",
			showEmojis: true,
			configNil:  false,
			expected:   "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config
			originalConfig := config.Global
			defer func() { config.Global = originalConfig }()

			if tt.configNil {
				config.Global = nil
			} else {
				config.Global = &config.Config{ShowEmojis: tt.showEmojis}
			}

			got := Emoji(tt.emoji)
			if got != tt.expected {
				t.Errorf("Emoji(%q) = %q, want %q", tt.emoji, got, tt.expected)
			}
		})
	}
}

// TestEmojiConstants verifies that all emoji constants are defined
func TestEmojiConstants(t *testing.T) {
	// Test that constants are not empty strings
	constants := map[string]string{
		"EmojiError":               EmojiError,
		"EmojiLoading":             EmojiLoading,
		"EmojiAvailable":           EmojiAvailable,
		"EmojiConnected":           EmojiConnected,
		"EmojiAppTitle":            EmojiAppTitle,
		"EmojiScanning":            EmojiScanning,
		"EmojiPaused":              EmojiPaused,
		"EmojiSuccess":             EmojiSuccess,
		"EmojiPairingKey":          EmojiPairingKey,
		"EmojiKeyboard":            EmojiKeyboard,
		"EmojiSelector":            EmojiSelector,
		"EmojiFocusMarker":         EmojiFocusMarker,
		"EmojiHeadphones":          EmojiHeadphones,
		"EmojiPhone":               EmojiPhone,
		"EmojiComputer":            EmojiComputer,
		"EmojiKeyboardDev":         EmojiKeyboardDev,
		"EmojiMouse":               EmojiMouse,
		"EmojiGaming":              EmojiGaming,
		"EmojiCamera":              EmojiCamera,
		"EmojiPrinter":             EmojiPrinter,
		"EmojiGeneric":             EmojiGeneric,
		"EmojiBatteryFull":         EmojiBatteryFull,
		"EmojiBatteryLow":          EmojiBatteryLow,
		"EmojiAdapterOff":          EmojiAdapterOff,
		"EmojiAdapterDiscovering":  EmojiAdapterDiscovering,
		"EmojiAdapterOn":           EmojiAdapterOn,
	}

	for name, value := range constants {
		if value == "" {
			t.Errorf("Constant %s should not be empty", name)
		}
	}
}
