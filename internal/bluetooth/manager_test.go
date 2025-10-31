package bluetooth

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

// TestManager_GetConnection tests the GetConnection method
func TestManager_GetConnection(t *testing.T) {
	// Create a manager with a nil connection
	m := &Manager{conn: nil}

	conn := m.GetConnection()
	if conn != nil {
		t.Errorf("GetConnection() should return nil when conn is nil")
	}

	// We can't create a real connection without DBus running,
	// but we can verify the method works
}

// TestManager_GetAdapter tests the GetAdapter method
func TestManager_GetAdapter(t *testing.T) {
	testPath := dbus.ObjectPath("/org/bluez/hci0")
	m := &Manager{adapter: testPath}

	adapter := m.GetAdapter()
	if adapter != testPath {
		t.Errorf("GetAdapter() = %v, want %v", adapter, testPath)
	}
}

// TestManager_Close tests the Close method
func TestManager_Close(t *testing.T) {
	t.Run("returns nil when conn is nil", func(t *testing.T) {
		m := &Manager{conn: nil}
		err := m.Close()
		if err != nil {
			t.Errorf("Close() with nil conn should return nil, got %v", err)
		}
	})

	// We can't test Close with a real connection without DBus running
}

// TestManager_GetDevices tests that GetDevices method exists
func TestManager_GetDevices(t *testing.T) {
	// We can't test this without a real DBus connection
	// Verify the method exists
	_ = (*Manager).GetDevices
	t.Log("GetDevices method exists on Manager type")
}

// TestConstants verifies the BlueZ constants are defined
func TestConstants(t *testing.T) {
	if bluezService != "org.bluez" {
		t.Errorf("bluezService = %v, want org.bluez", bluezService)
	}
	if bluezAdapterIface != "org.bluez.Adapter1" {
		t.Errorf("bluezAdapterIface = %v, want org.bluez.Adapter1", bluezAdapterIface)
	}
	if bluezDeviceIface != "org.bluez.Device1" {
		t.Errorf("bluezDeviceIface = %v, want org.bluez.Device1", bluezDeviceIface)
	}
	if bluezBatteryIface != "org.bluez.Battery1" {
		t.Errorf("bluezBatteryIface = %v, want org.bluez.Battery1", bluezBatteryIface)
	}
}

// TestManager_Struct verifies the Manager structure
func TestManager_Struct(t *testing.T) {
	m := &Manager{
		conn:    nil,
		adapter: dbus.ObjectPath("/org/bluez/hci0"),
	}

	if m.conn != nil {
		t.Errorf("conn should be nil")
	}
	if m.adapter != "/org/bluez/hci0" {
		t.Errorf("adapter = %v, want /org/bluez/hci0", m.adapter)
	}
}
