package models

import "testing"

func BenchmarkNormalizeMAC(b *testing.B) {
	mac := "AA:BB:CC:DD:EE:FF"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NormalizeMAC(mac)
	}
}

func BenchmarkIsAliasMACAddress(b *testing.B) {
	alias := "AA-BB-CC-DD-EE-FF"
	address := "AA:BB:CC:DD:EE:FF"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsAliasMACAddress(alias, address)
	}
}

func BenchmarkIsAliasMACAddress_NotMAC(b *testing.B) {
	alias := "My Bluetooth Device"
	address := "AA:BB:CC:DD:EE:FF"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsAliasMACAddress(alias, address)
	}
}

func BenchmarkHasRealName(b *testing.B) {
	dev := &Device{
		Name:    "",
		Alias:   "AA-BB-CC-DD-EE-FF",
		Address: "AA:BB:CC:DD:EE:FF",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dev.HasRealName()
	}
}

func BenchmarkHasRealName_WithName(b *testing.B) {
	dev := &Device{
		Name:    "Real Device",
		Alias:   "AA-BB-CC-DD-EE-FF",
		Address: "AA:BB:CC:DD:EE:FF",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dev.HasRealName()
	}
}
