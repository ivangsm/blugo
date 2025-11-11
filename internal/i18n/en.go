package i18n

// englishTranslations contains all English translations
var englishTranslations = Translations{
	// App
	AppTitle:     "BLUGO - Bluetooth Manager",
	Scanning:     "Scanning",
	Paused:       "Paused",
	Initializing: "Initializing Bluetooth...",

	// Sections
	AvailableDevices: "AVAILABLE DEVICES",
	ConnectedDevices: "CONNECTED DEVICES",
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
	AdapterPoweredOn:         "Bluetooth adapter powered on",
	AdapterPoweredOff:        "Bluetooth adapter powered off",
	DiscoverableActivating:   "Activating discoverable mode...",
	DiscoverableDeactivating: "Deactivating discoverable mode...",
	DiscoverableOn:           "Discoverable mode activated",
	DiscoverableOff:          "Discoverable mode deactivated",
	PairableActivating:       "Activating pairable mode...",
	PairableDeactivating:     "Deactivating pairable mode...",
	PairableOn:               "Pairable mode activated",
	PairableOff:              "Pairable mode deactivated",

	// Status messages
	ScanEnabled:        "Scanning enabled",
	ScanPaused:         "Scanning paused",
	Connected:          "Connected to %s",
	Disconnected:       "Disconnected from %s",
	DisconnectedPaired: "Disconnected from %s (still paired)",
	Forgotten:          "Device forgotten",

	// Errors
	Error:           "Error",
	ErrorScanToggle: "Error toggling scan: %s",

	// Pairing
	PairingCode:        "PAIRING CODE: %06d",
	PairingInstruction: "Type this code on your keyboard and press Enter",
	PairingConfirm:     "Then press Enter here to confirm, or Esc/N to cancel",
	PairingCancelled:   "Pairing cancelled",

	// Help
	HelpNavigation:     "↑↓, kj: navigate | enter: connect/disconnect | d/x: forget | q: quit",
	HelpActions:        "↑↓, kj: navigate | enter: disconnect | d/x: forget",
	HelpAdapterControl: "s: scan | p: power | v: discoverable | b: pairable | l: language | r: refresh",
	HelpScroll:         "PgUp/PgDn: scroll page | Ctrl+↑↓, kj: scroll | Home/End: top/bottom | Mouse wheel: scroll",
	HelpGeneral:        "q: quit",
	HelpPairing:        "enter: confirm | n/esc: cancel | q: quit",
	HelpCollapsed:      "?: toggle help | q: quit",
	HelpExpanded:       "?: hide help",

	// Adapter table
	AdapterName:         "Name",
	AdapterAlias:        "Alias",
	AdapterPower:        "Power",
	AdapterPairable:     "Pairable",
	AdapterDiscoverable: "Discoverable",

	// Device table
	DeviceIcon:    "Icon",
	DeviceName:    "Name",
	DeviceAddress: "Address",
	DeviceRSSI:    "Signal",
	DeviceBattery: "Battery",
	DeviceStatus:  "Status",

	// Status
	StatusOn:  "ON",
	StatusOff: "OFF",

	// Badges
	BadgePaired:    "PAIRED",
	BadgeConnected: "CONNECTED",
	BadgeTrusted:   "Trusted",

	// Error messages
	ErrorDBusConnection:         "Could not connect to DBus",
	ErrorAdapterNotFound:        "No Bluetooth adapter found",
	ErrorStartDiscovery:         "Could not start discovery",
	ErrorStopDiscovery:          "Could not stop discovery",
	ErrorRemoveDevice:           "Could not remove device",
	ErrorPairDevice:             "Error pairing device",
	ErrorTrustDevice:            "Error trusting device",
	ErrorConnectDevice:          "Error connecting device",
	ErrorDisconnectDevice:       "Error disconnecting device",
	ErrorGetDevices:             "Error getting devices",
	ErrorGetAdapterInfo:         "Error getting adapter info",
	ErrorSetAdapterPowered:      "Error changing adapter state",
	ErrorSetAdapterDiscoverable: "Error changing discoverable mode",
	ErrorSetAdapterPairable:     "Error changing pairable mode",
	ErrorSetAdapterAlias:        "Error changing adapter alias",
	ErrorForgetDevice:           "Error forgetting device",
	ErrorChangeProperty:         "Error changing",

	// Status messages
	StatusConfirmingPairing:  "Confirming pairing...",
	StatusLoadingAdapterInfo: "Loading adapter information...",

	// Warnings
	WarningAgentRegistration:       "⚠️  Warning: Could not register pairing agent",
	WarningAgentRegistrationDetail: "   The app will work but some devices may require manual pairing.",

	// Agent errors (internal)
	ErrorRequestPasskey:   "Cannot request passkey in TUI",
	ErrorPairingCancelled: "Pairing cancelled by user",
	ErrorConfirmRejected:  "Confirmation rejected",
	ErrorExportAgent:      "Could not export agent",
	ErrorExportIntrospect: "Could not export introspection",
	ErrorRegisterAgent:    "Could not register agent",
}
