package bluetooth

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

// StartDiscovery inicia el escaneo de dispositivos Bluetooth.
func (m *Manager) StartDiscovery() error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call(bluezAdapterIface+".StartDiscovery", 0).Err
	if err != nil {
		return fmt.Errorf("no se pudo iniciar descubrimiento: %w", err)
	}
	return nil
}

// StopDiscovery detiene el escaneo de dispositivos Bluetooth.
func (m *Manager) StopDiscovery() error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call(bluezAdapterIface+".StopDiscovery", 0).Err
	if err != nil {
		return fmt.Errorf("no se pudo detener descubrimiento: %w", err)
	}
	return nil
}

// RemoveDevice elimina un dispositivo del adaptador.
func (m *Manager) RemoveDevice(devicePath dbus.ObjectPath) error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call(bluezAdapterIface+".RemoveDevice", 0, devicePath).Err
	if err != nil {
		return fmt.Errorf("no se pudo eliminar dispositivo: %w", err)
	}
	return nil
}
