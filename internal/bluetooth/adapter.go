package bluetooth

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/i18n"
)

// StartDiscovery starts scanning for Bluetooth devices.
func (m *Manager) StartDiscovery() error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call(bluezAdapterIface+".StartDiscovery", 0).Err
	if err != nil {
		return fmt.Errorf(i18n.T.ErrorStartDiscovery+": %w", err)
	}
	return nil
}

// StopDiscovery stops scanning for Bluetooth devices.
func (m *Manager) StopDiscovery() error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call(bluezAdapterIface+".StopDiscovery", 0).Err
	if err != nil {
		return fmt.Errorf(i18n.T.ErrorStopDiscovery+": %w", err)
	}
	return nil
}

// RemoveDevice removes a device from the adapter.
func (m *Manager) RemoveDevice(devicePath dbus.ObjectPath) error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call(bluezAdapterIface+".RemoveDevice", 0, devicePath).Err
	if err != nil {
		return fmt.Errorf(i18n.T.ErrorRemoveDevice+": %w", err)
	}
	return nil
}
