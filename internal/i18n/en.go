package i18n

// englishTranslations contains all English translations
var englishTranslations = Translations{
	// App
	AppTitle:     "🔵 GOB - Bluetooth Manager",
	Scanning:     "🔍 Scanning",
	Paused:       "⏸ Paused",
	Initializing: "⚙ Initializing Bluetooth...",

	// Sections
	AvailableDevices: "📡 AVAILABLE DEVICES",
	ConnectedDevices: "🔗 CONNECTED DEVICES",
	AdapterInfo:      "Bluetooth Adapter",

	// Device info
	NoDevicesAvailable: "No devices available",
	NoDevicesConnected: "No connected devices",

	// Actions
	Connecting:    "Connecting to %s...",
	Disconnecting: "Disconnecting from %s (keeping pairing)...",
	Pairing:       "Pairing with %s...",
	Forgetting:    "Forgetting %s...",

	// Adapter
	AdapterPoweringOn:        "Turning Bluetooth adapter on...",
	AdapterPoweringOff:       "Turning Bluetooth adapter off...",
	AdapterPoweredOn:         "✓ Bluetooth adapter powered on",
	AdapterPoweredOff:        "✓ Bluetooth adapter powered off",
	DiscoverableActivating:   "Activating discoverable mode...",
	DiscoverableDeactivating: "Deactivating discoverable mode...",
	DiscoverableOn:           "✓ Discoverable mode activated",
	DiscoverableOff:          "✓ Discoverable mode deactivated",
	PairableActivating:       "Activating pairable mode...",
	PairableDeactivating:     "Deactivating pairable mode...",
	PairableOn:               "✓ Pairable mode activated",
	PairableOff:              "✓ Pairable mode deactivated",

	// Status messages
	ScanEnabled:        "Scanning enabled",
	ScanPaused:         "Scanning paused",
	Connected:          "✓ Connected to %s",
	Disconnected:       "✓ Disconnected from %s",
	DisconnectedPaired: "✓ Disconnected from %s (still paired)",
	Forgotten:          "✓ Device forgotten",

	// Errors
	Error: "❌ Error",

	// Pairing
	PairingCode:        "🔑 PAIRING CODE: %06d",
	PairingInstruction: "⌨️  Type this code on your keyboard and press Enter",
	PairingConfirm:     "Then press Enter here to confirm, or Esc/N to cancel",
	PairingCancelled:   "Pairing cancelled",

	// Help
	HelpNavigation:     "↑/↓: navigate | Tab: switch | Enter: connect | D/X: forget",
	HelpActions:        "↑/↓: navigate | Tab: switch | Enter: disconnect | D/X: forget",
	HelpAdapterControl: "S: scan | P: power | V: discoverable | B: pairable | L: language | R: refresh | Q: quit",
	HelpGeneral:        "Q: quit",
	HelpPairing:        "Enter: confirm | N/Esc: cancel | Q: quit",

	// Adapter table
	AdapterName:         "Name",
	AdapterAlias:        "Alias",
	AdapterPower:        "Power",
	AdapterPairable:     "Pairable",
	AdapterDiscoverable: "Discoverable",

	// Status
	StatusOn:  "ON",
	StatusOff: "OFF",

	// Badges
	BadgePaired:    "PAIRED",
	BadgeConnected: "CONNECTED",
	BadgeTrusted:   "Trusted",
}
