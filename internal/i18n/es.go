package i18n

// spanishTranslations contains all Spanish translations
var spanishTranslations = Translations{
	// App
	AppTitle:     "🔵 GOB - Gestor Bluetooth",
	Scanning:     "🔍 Escaneando",
	Paused:       "⏸ Pausado",
	Initializing: "⚙ Inicializando Bluetooth...",

	// Sections
	AvailableDevices: "📡 DISPOSITIVOS DISPONIBLES",
	ConnectedDevices: "🔗 DISPOSITIVOS CONECTADOS",
	AdapterInfo:      "Adaptador Bluetooth",

	// Device info
	NoDevicesAvailable: "No hay dispositivos disponibles",
	NoDevicesConnected: "No hay dispositivos conectados",

	// Actions
	Connecting:    "Conectando a %s...",
	Disconnecting: "Desconectando de %s (manteniendo pairing)...",
	Pairing:       "Pareando %s...",
	Forgetting:    "Olvidando %s...",

	// Adapter
	AdapterPoweringOn:        "Encendiendo adaptador Bluetooth...",
	AdapterPoweringOff:       "Apagando adaptador Bluetooth...",
	AdapterPoweredOn:         "✓ Adaptador Bluetooth encendido",
	AdapterPoweredOff:        "✓ Adaptador Bluetooth apagado",
	DiscoverableActivating:   "Activando modo discoverable...",
	DiscoverableDeactivating: "Desactivando modo discoverable...",
	DiscoverableOn:           "✓ Modo discoverable activado",
	DiscoverableOff:          "✓ Modo discoverable desactivado",
	PairableActivating:       "Activando modo pairable...",
	PairableDeactivating:     "Desactivando modo pairable...",
	PairableOn:               "✓ Modo pairable activado",
	PairableOff:              "✓ Modo pairable desactivado",

	// Status messages
	ScanEnabled:        "Escaneo activado",
	ScanPaused:         "Escaneo pausado",
	Connected:          "✓ Conectado a %s",
	Disconnected:       "✓ Desconectado de %s",
	DisconnectedPaired: "✓ Desconectado de %s (aún pareado)",
	Forgotten:          "✓ Dispositivo olvidado",

	// Errors
	Error: "❌ Error",

	// Pairing
	PairingCode:        "🔑 CÓDIGO DE PAIRING: %06d",
	PairingInstruction: "⌨️  Escribe este código en tu teclado y presiona Enter",
	PairingConfirm:     "Luego presiona Enter aquí para confirmar, o Esc/N para cancelar",
	PairingCancelled:   "Pairing cancelado",

	// Help
	HelpNavigation:     "↑/↓: navegar | Tab: cambiar | Enter: conectar | D/X: olvidar",
	HelpActions:        "↑/↓: navegar | Tab: cambiar | Enter: desconectar | D/X: olvidar",
	HelpAdapterControl: "S: escaneo | P: powered | V: discoverable | B: pairable | L: idioma | R: refrescar | Q: salir",
	HelpGeneral:        "Q: salir",
	HelpPairing:        "Enter: confirmar | N/Esc: cancelar | Q: salir",

	// Adapter table
	AdapterName:         "Nombre",
	AdapterAlias:        "Alias",
	AdapterPower:        "Energía",
	AdapterPairable:     "Pairable",
	AdapterDiscoverable: "Descubrible",

	// Status
	StatusOn:  "ON",
	StatusOff: "OFF",

	// Badges
	BadgePaired:    "PAREADO",
	BadgeConnected: "CONECTADO",
	BadgeTrusted:   "Confiable",
}
