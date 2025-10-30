package models

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
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
			return "🎧"
		case "phone", "smartphone":
			return "📱"
		case "computer", "laptop":
			return "💻"
		case "input-keyboard":
			return "⌨️"
		case "input-mouse":
			return "🖱️"
		case "input-gaming":
			return "🎮"
		case "camera":
			return "📷"
		case "printer":
			return "🖨️"
		}
	}

	// Fallback basado en clase de dispositivo
	majorClass := (d.Class >> 8) & 0x1F
	switch majorClass {
	case 1: // Computer
		return "💻"
	case 2: // Phone
		return "📱"
	case 4: // Audio/Video
		return "🎧"
	case 5: // Peripheral (keyboard, mouse, etc)
		return "⌨️"
	case 6: // Imaging (printer, camera)
		return "📷"
	}

	return "📶"
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
		icon = "🔋" // Batería llena
	case level >= 60:
		icon = "🔋" // Batería alta
	case level >= 30:
		icon = "🔋" // Batería media
	case level >= 15:
		icon = "🪫" // Batería baja
	default:
		icon = "🪫" // Batería muy baja/crítica
	}

	// Formato del texto
	text = fmt.Sprintf("%d%%", level)

	return icon, text
}

// HasBattery indica si el dispositivo reporta nivel de batería.
func (d *Device) HasBattery() bool {
	return d.Battery != nil
}
