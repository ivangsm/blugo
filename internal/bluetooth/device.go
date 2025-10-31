package bluetooth

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/ivangsm/blugo/internal/i18n"
	"github.com/ivangsm/blugo/internal/models"
)

// getDevices gets all Bluetooth devices from the system.
func getDevices(conn *dbus.Conn) (map[string]*models.Device, error) {
	obj := conn.Object(bluezService, "/")
	var paths map[string]map[string]map[string]dbus.Variant
	err := obj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&paths)
	if err != nil {
		return nil, err
	}

	devices := make(map[string]*models.Device)
	for path, interfaces := range paths {
		if props, ok := interfaces[bluezDeviceIface]; ok {
			dev := parseDevice(dbus.ObjectPath(path), interfaces, props)
			devices[dev.Address] = dev
		}
	}

	return devices, nil
}

// parseDevice converts DBus properties into a Device model.
func parseDevice(path dbus.ObjectPath, interfaces map[string]map[string]dbus.Variant, props map[string]dbus.Variant) *models.Device {
	dev := &models.Device{
		Path:     path,
		LastSeen: time.Now(),
	}

	if variant, ok := props["Address"]; ok {
		if v, ok := variant.Value().(string); ok {
			dev.Address = v
		}
	}
	if variant, ok := props["Name"]; ok {
		if v, ok := variant.Value().(string); ok {
			dev.Name = v
		}
	}
	if variant, ok := props["Alias"]; ok {
		if v, ok := variant.Value().(string); ok {
			dev.Alias = v
		}
	}
	if variant, ok := props["Paired"]; ok {
		if v, ok := variant.Value().(bool); ok {
			dev.Paired = v
		}
	}
	if variant, ok := props["Trusted"]; ok {
		if v, ok := variant.Value().(bool); ok {
			dev.Trusted = v
		}
	}
	if variant, ok := props["Connected"]; ok {
		if v, ok := variant.Value().(bool); ok {
			dev.Connected = v
		}
	}
	if variant, ok := props["RSSI"]; ok {
		if v, ok := variant.Value().(int16); ok {
			dev.RSSI = v
		}
	}
	if variant, ok := props["Icon"]; ok {
		if v, ok := variant.Value().(string); ok {
			dev.Icon = v
		}
	}
	if variant, ok := props["Class"]; ok {
		if v, ok := variant.Value().(uint32); ok {
			dev.Class = v
		}
	}

	// Use Alias if no Name
	if dev.Name == "" && dev.Alias != "" {
		dev.Name = dev.Alias
	}

	// Get battery information if available
	if batteryProps, ok := interfaces[bluezBatteryIface]; ok {
		if variant, ok := batteryProps["Percentage"]; ok {
			if percentage, ok := variant.Value().(byte); ok {
				dev.Battery = &percentage
			}
		}
	}

	return dev
}

// PairDevice pairs a device.
func (m *Manager) PairDevice(devicePath dbus.ObjectPath) error {
	obj := m.conn.Object(bluezService, devicePath)
	err := obj.Call(bluezDeviceIface+".Pair", 0).Err
	if err != nil {
		return fmt.Errorf(i18n.T.ErrorPairDevice+": %w", err)
	}
	return nil
}

// TrustDevice marks a device as trusted.
func (m *Manager) TrustDevice(devicePath dbus.ObjectPath) error {
	obj := m.conn.Object(bluezService, devicePath)
	err := obj.Call("org.freedesktop.DBus.Properties.Set", 0,
		bluezDeviceIface, "Trusted", dbus.MakeVariant(true)).Err
	if err != nil {
		return fmt.Errorf(i18n.T.ErrorTrustDevice+": %w", err)
	}
	return nil
}

// ConnectDevice connects to a device.
func (m *Manager) ConnectDevice(devicePath dbus.ObjectPath) error {
	obj := m.conn.Object(bluezService, devicePath)
	err := obj.Call(bluezDeviceIface+".Connect", 0).Err
	if err != nil {
		return fmt.Errorf(i18n.T.ErrorConnectDevice+": %w", err)
	}
	return nil
}

// DisconnectDevice disconnects a device.
func (m *Manager) DisconnectDevice(devicePath dbus.ObjectPath) error {
	obj := m.conn.Object(bluezService, devicePath)
	err := obj.Call(bluezDeviceIface+".Disconnect", 0).Err
	if err != nil {
		return fmt.Errorf(i18n.T.ErrorDisconnectDevice+": %w", err)
	}
	return nil
}
