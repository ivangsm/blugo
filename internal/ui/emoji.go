package ui

import "github.com/ivangsm/blugo/internal/config"

// Emoji returns the emoji if ShowEmojis is enabled in config, otherwise empty string
func Emoji(emoji string) string {
	if config.Global != nil && config.Global.ShowEmojis {
		return emoji
	}
	return ""
}

// Common UI emojis
const (
	EmojiError       = "❌"
	EmojiLoading     = "⚙"
	EmojiAvailable   = "📡"
	EmojiConnected   = "🔗"
	EmojiAppTitle    = "🔵"
	EmojiScanning    = "🔍"
	EmojiPaused      = "⏸"
	EmojiSuccess     = "✓"
	EmojiPairingKey  = "🔑"
	EmojiKeyboard    = "⌨️"
	EmojiSelector    = "▶"
	EmojiFocusMarker = "◀"
)

// Device type emojis
const (
	EmojiHeadphones = "🎧"
	EmojiPhone      = "📱"
	EmojiComputer   = "💻"
	EmojiKeyboardDev = "⌨️"
	EmojiMouse      = "🖱️"
	EmojiGaming     = "🎮"
	EmojiCamera     = "📷"
	EmojiPrinter    = "🖨️"
	EmojiGeneric    = "📶"
)

// Battery emojis
const (
	EmojiBatteryFull = "🔋"
	EmojiBatteryLow  = "🪫"
)

// Adapter emojis
const (
	EmojiAdapterOff        = "⚫"
	EmojiAdapterDiscovering = "🔍"
	EmojiAdapterOn         = "🔵"
)
