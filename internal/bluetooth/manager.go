package bluetooth

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/gob/internal/models"
)

const (
	bluezService      = "org.bluez"
	bluezAdapterIface = "org.bluez.Adapter1"
	bluezDeviceIface  = "org.bluez.Device1"
)

// Manager gestiona la conexión con BlueZ a través de DBus.
type Manager struct {
	conn    *dbus.Conn
	adapter dbus.ObjectPath
}

// NewManager crea una nueva instancia del manager de Bluetooth.
func NewManager() (*Manager, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a DBus: %w", err)
	}

	adapter, err := getAdapter(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("no se encontró adaptador Bluetooth: %w", err)
	}

	return &Manager{
		conn:    conn,
		adapter: adapter,
	}, nil
}

// Close cierra la conexión con DBus.
func (m *Manager) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

// GetConnection devuelve la conexión DBus.
func (m *Manager) GetConnection() *dbus.Conn {
	return m.conn
}

// GetAdapter devuelve el path del adaptador Bluetooth.
func (m *Manager) GetAdapter() dbus.ObjectPath {
	return m.adapter
}

// GetDevices obtiene todos los dispositivos Bluetooth conocidos.
func (m *Manager) GetDevices() (map[string]*models.Device, error) {
	return getDevices(m.conn)
}

// getAdapter encuentra el primer adaptador Bluetooth disponible.
func getAdapter(conn *dbus.Conn) (dbus.ObjectPath, error) {
	obj := conn.Object(bluezService, "/")
	var paths map[string]map[string]map[string]dbus.Variant
	err := obj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&paths)
	if err != nil {
		return "", err
	}

	for path, interfaces := range paths {
		if _, ok := interfaces[bluezAdapterIface]; ok {
			return dbus.ObjectPath(path), nil
		}
	}
	return "", fmt.Errorf("no se encontró adaptador Bluetooth")
}
