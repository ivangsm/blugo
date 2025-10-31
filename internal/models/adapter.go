package models

import "github.com/godbus/dbus/v5"

// Adapter represents a Bluetooth adapter in the system.
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

// GetDisplayName returns the display name of the adapter.
func (a *Adapter) GetDisplayName() string {
	if a.Alias != "" {
		return a.Alias
	}
	if a.Name != "" {
		return a.Name
	}
	return a.Address
}

// GetStatusIcon returns the icon based on adapter status.
func (a *Adapter) GetStatusIcon() string {
	if !a.Powered {
		return emoji("⚫") // Off
	}
	if a.Discovering {
		return emoji("🔍") // Scanning
	}
	return emoji("🔵") // On
}
