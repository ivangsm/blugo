package models

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/config"
)

// Device representa un dispositivo Bluetooth.
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
	Battery   *uint8 // Nivel de batería (0-100), nil si no disponible
	LastSeen  time.Time
}

// emoji returns the emoji if ShowEmojis is enabled, otherwise empty string
func emoji(e string) string {
	if config.Global != nil && config.Global.ShowEmojis {
		return e
	}
	return ""
}

// GetDisplayName devuelve el nombre a mostrar del dispositivo.
// Prioriza: Name > Alias > Address.
func (d *Device) GetDisplayName() string {
	if d.Name != "" {
		return d.Name
	}
	if d.Alias != "" {
		return d.Alias
	}
	return d.Address
}

// IsAvailable determina si el dispositivo está disponible pero no conectado.
func (d *Device) IsAvailable() bool {
	return !d.Connected
}

// GetIcon devuelve el icono apropiado según el tipo de dispositivo.
func (d *Device) GetIcon() string {
	// Iconos según el tipo de dispositivo
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

	// Fallback basado en clase de dispositivo
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

// GetBatteryInfo devuelve el icono y texto de la batería.
// Retorna ("", "") si no hay información de batería disponible.
func (d *Device) GetBatteryInfo() (icon string, text string) {
	if d.Battery == nil {
		return "", ""
	}

	level := *d.Battery

	// Elegir icono según el nivel
	switch {
	case level >= 90:
		icon = emoji("🔋") // Batería llena
	case level >= 60:
		icon = emoji("🔋") // Batería alta
	case level >= 30:
		icon = emoji("🔋") // Batería media
	case level >= 15:
		icon = emoji("🪫") // Batería baja
	default:
		icon = emoji("🪫") // Batería muy baja/crítica
	}

	// Formato del texto
	text = fmt.Sprintf("%d%%", level)

	return icon, text
}

// HasBattery indica si el dispositivo reporta nivel de batería.
func (d *Device) HasBattery() bool {
	return d.Battery != nil
}
