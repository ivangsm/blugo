package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/config"
)

// Device represents a Bluetooth device.
type Device struct {
	Path      dbus.ObjectPath
	Address   string
	Name      string
	Alias     string
	Paired    bool
	Trusted   bool
	Connected bool
	RSSI      int16
	Icon      string
	Class     uint32
	Battery   *uint8 // Battery level (0-100), nil if not available
	LastSeen  time.Time
}

// emoji returns the emoji if ShowEmojis is enabled, otherwise empty string
func emoji(e string) string {
	if config.Global != nil && config.Global.ShowEmojis {
		return e
	}
	return ""
}

// NormalizeMAC removes separators (: and -) from MAC addresses for comparison
func NormalizeMAC(mac string) string {
	// Use strings.Builder for efficient string concatenation
	var builder strings.Builder
	builder.Grow(12) // MAC address without separators is always 12 characters

	for i := 0; i < len(mac); i++ {
		c := mac[i]
		// Skip separators
		if c == ':' || c == '-' {
			continue
		}
		// Convert to lowercase
		if c >= 'A' && c <= 'F' {
			c += 'a' - 'A'
		}
		builder.WriteByte(c)
	}

	return builder.String()
}

// IsAliasMACAddress checks if the alias is just the MAC address with different separator
func IsAliasMACAddress(alias, address string) bool {
	// Quick length check: normalized MAC should be 12 chars, with separators 17 chars
	if len(alias) != 17 && len(alias) != 12 {
		return false
	}
	if len(address) != 17 && len(address) != 12 {
		return false
	}

	return NormalizeMAC(alias) == NormalizeMAC(address)
}

// GetDisplayName returns the display name of the device.
// Prioritizes: Name > Alias > Address.
func (d *Device) GetDisplayName() string {
	if d.Name != "" {
		return d.Name
	}
	if d.Alias != "" {
		return d.Alias
	}
	return d.Address
}

// IsAvailable determines if the device is available but not connected.
func (d *Device) IsAvailable() bool {
	return !d.Connected
}

// HasRealName checks if the device has a real name (not just a MAC address).
// Returns true if Name is set, or if Alias is set and is not the MAC address.
func (d *Device) HasRealName() bool {
	// If Name is set, it's a real name
	if d.Name != "" {
		return true
	}
	// If Alias is empty, no real name
	if d.Alias == "" {
		return false
	}
	// Check if Alias is not the MAC address
	return !IsAliasMACAddress(d.Alias, d.Address)
}

// GetIcon returns the appropriate icon based on device type.
func (d *Device) GetIcon() string {
	// Icons based on device type
	if d.Icon != "" {
		switch d.Icon {
		case "audio-card", "audio-headset", "audio-headphones":
			return emoji("🎧")
		case "phone", "smartphone":
			return emoji("📱")
		case "computer", "laptop":
			return emoji("💻")
		case "input-keyboard":
			return emoji("⌨️")
		case "input-mouse":
			return emoji("🖱️")
		case "input-gaming":
			return emoji("🎮")
		case "camera":
			return emoji("📷")
		case "printer":
			return emoji("🖨️")
		}
	}

	// Fallback based on device class
	majorClass := (d.Class >> 8) & 0x1F
	switch majorClass {
	case 1: // Computer
		return emoji("💻")
	case 2: // Phone
		return emoji("📱")
	case 4: // Audio/Video
		return emoji("🎧")
	case 5: // Peripheral (keyboard, mouse, etc)
		return emoji("⌨️")
	case 6: // Imaging (printer, camera)
		return emoji("📷")
	}

	return emoji("📶")
}

// GetBatteryInfo returns the battery icon and text.
// Returns ("", "") if no battery information is available.
func (d *Device) GetBatteryInfo() (icon string, text string) {
	if d.Battery == nil {
		return "", ""
	}

	level := *d.Battery

	// Choose icon based on level
	switch {
	case level >= 90:
		icon = emoji("🔋") // Full battery
	case level >= 60:
		icon = emoji("🔋") // High battery
	case level >= 30:
		icon = emoji("🔋") // Medium battery
	case level >= 15:
		icon = emoji("🪫") // Low battery
	default:
		icon = emoji("🪫") // Very low/critical battery
	}

	// Text format
	text = fmt.Sprintf("%d%%", level)

	return icon, text
}

// HasBattery indicates if the device reports battery level.
func (d *Device) HasBattery() bool {
	return d.Battery != nil
}
