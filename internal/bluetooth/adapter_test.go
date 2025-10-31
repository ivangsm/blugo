package bluetooth

import (
	"testing"
)

// Note: Most bluetooth functions require a real DBus connection and cannot be
// easily tested without integration tests or comprehensive mocking.
// The device parsing logic is tested in device_test.go

// TestAdapterMethodsExist verifies the adapter methods are defined
func TestAdapterMethodsExist(t *testing.T) {
	// This test verifies that the Manager type has the expected methods
	// We can't call them without a real DBus connection

	// Just verify the package compiles and has the expected structure
	_ = (*Manager).StartDiscovery
	_ = (*Manager).StopDiscovery
	_ = (*Manager).RemoveDevice

	t.Log("Adapter methods exist on Manager type")
}
