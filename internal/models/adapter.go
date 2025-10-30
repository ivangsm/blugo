package models

import "github.com/godbus/dbus/v5"

// Adapter representa un adaptador Bluetooth del sistema.
type Adapter struct {
	Path         dbus.ObjectPath
	Address      string
	Name         string
	Alias        string
	Powered      bool
	Discoverable bool
	Pairable     bool
	Discovering  bool
}

// GetDisplayName retorna el nombre a mostrar del adaptador.
func (a *Adapter) GetDisplayName() string {
	if a.Alias != "" {
		return a.Alias
	}
	if a.Name != "" {
		return a.Name
	}
	return a.Address
}

// GetStatusIcon retorna el icono según el estado del adaptador.
func (a *Adapter) GetStatusIcon() string {
	if !a.Powered {
		return "⚫" // Apagado
	}
	if a.Discovering {
		return "🔍" // Escaneando
	}
	return "🔵" // Encendido
}
