package bluetooth

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/gob/internal/models"
)

const (
	bluezService       = "org.bluez"
	bluezAdapterIface  = "org.bluez.Adapter1"
	bluezDeviceIface   = "org.bluez.Device1"
	bluezBatteryIface  = "org.bluez.Battery1"
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

// GetAdapterInfo obtiene la información del adaptador Bluetooth.
func (m *Manager) GetAdapterInfo() (*models.Adapter, error) {
	obj := m.conn.Object(bluezService, m.adapter)

	adapter := &models.Adapter{
		Path: m.adapter,
	}

	// Obtener todas las propiedades del adaptador
	props, err := obj.GetProperty("org.bluez.Adapter1.Address")
	if err == nil {
		if address, ok := props.Value().(string); ok {
			adapter.Address = address
		}
	}

	props, err = obj.GetProperty("org.bluez.Adapter1.Name")
	if err == nil {
		if name, ok := props.Value().(string); ok {
			adapter.Name = name
		}
	}

	props, err = obj.GetProperty("org.bluez.Adapter1.Alias")
	if err == nil {
		if alias, ok := props.Value().(string); ok {
			adapter.Alias = alias
		}
	}

	props, err = obj.GetProperty("org.bluez.Adapter1.Powered")
	if err == nil {
		if powered, ok := props.Value().(bool); ok {
			adapter.Powered = powered
		}
	}

	props, err = obj.GetProperty("org.bluez.Adapter1.Discoverable")
	if err == nil {
		if discoverable, ok := props.Value().(bool); ok {
			adapter.Discoverable = discoverable
		}
	}

	props, err = obj.GetProperty("org.bluez.Adapter1.Pairable")
	if err == nil {
		if pairable, ok := props.Value().(bool); ok {
			adapter.Pairable = pairable
		}
	}

	props, err = obj.GetProperty("org.bluez.Adapter1.Discovering")
	if err == nil {
		if discovering, ok := props.Value().(bool); ok {
			adapter.Discovering = discovering
		}
	}

	return adapter, nil
}

// SetAdapterPowered enciende o apaga el adaptador Bluetooth.
func (m *Manager) SetAdapterPowered(powered bool) error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call("org.freedesktop.DBus.Properties.Set", 0,
		bluezAdapterIface, "Powered", dbus.MakeVariant(powered)).Err
	if err != nil {
		return fmt.Errorf("error al cambiar estado del adaptador: %w", err)
	}
	return nil
}

// SetAdapterDiscoverable activa o desactiva el modo discoverable.
func (m *Manager) SetAdapterDiscoverable(discoverable bool) error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call("org.freedesktop.DBus.Properties.Set", 0,
		bluezAdapterIface, "Discoverable", dbus.MakeVariant(discoverable)).Err
	if err != nil {
		return fmt.Errorf("error al cambiar modo discoverable: %w", err)
	}
	return nil
}

// SetAdapterPairable activa o desactiva el modo pairable.
func (m *Manager) SetAdapterPairable(pairable bool) error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call("org.freedesktop.DBus.Properties.Set", 0,
		bluezAdapterIface, "Pairable", dbus.MakeVariant(pairable)).Err
	if err != nil {
		return fmt.Errorf("error al cambiar modo pairable: %w", err)
	}
	return nil
}

// SetAdapterAlias cambia el alias (nombre visible) del adaptador.
func (m *Manager) SetAdapterAlias(alias string) error {
	obj := m.conn.Object(bluezService, m.adapter)
	err := obj.Call("org.freedesktop.DBus.Properties.Set", 0,
		bluezAdapterIface, "Alias", dbus.MakeVariant(alias)).Err
	if err != nil {
		return fmt.Errorf("error al cambiar alias del adaptador: %w", err)
	}
	return nil
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
