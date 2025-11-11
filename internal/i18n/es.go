package i18n

// spanishTranslations contains all Spanish translations
var spanishTranslations = Translations{
	// App
	AppTitle:     "BLUGO - Gestor Bluetooth",
	Scanning:     "Escaneando",
	Paused:       "Pausado",
	Initializing: "Inicializando Bluetooth...",

	// Sections
	AvailableDevices: "DISPOSITIVOS DISPONIBLES",
	ConnectedDevices: "DISPOSITIVOS CONECTADOS",
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
	AdapterPoweredOn:         "Adaptador Bluetooth encendido",
	AdapterPoweredOff:        "Adaptador Bluetooth apagado",
	DiscoverableActivating:   "Activando modo discoverable...",
	DiscoverableDeactivating: "Desactivando modo discoverable...",
	DiscoverableOn:           "Modo discoverable activado",
	DiscoverableOff:          "Modo discoverable desactivado",
	PairableActivating:       "Activando modo pairable...",
	PairableDeactivating:     "Desactivando modo pairable...",
	PairableOn:               "Modo pairable activado",
	PairableOff:              "Modo pairable desactivado",

	// Status messages
	ScanEnabled:        "Escaneo activado",
	ScanPaused:         "Escaneo pausado",
	Connected:          "Conectado a %s",
	Disconnected:       "Desconectado de %s",
	DisconnectedPaired: "Desconectado de %s (aún pareado)",
	Forgotten:          "Dispositivo olvidado",

	// Errors
	Error:           "Error",
	ErrorScanToggle: "Error al cambiar escaneo: %s",

	// Pairing
	PairingCode:        "CÓDIGO DE PAIRING: %06d",
	PairingInstruction: "Escribe este código en tu teclado y presiona Enter",
	PairingConfirm:     "Luego presiona Enter aquí para confirmar, o Esc/N para cancelar",
	PairingCancelled:   "Pairing cancelado",

	// Help
	HelpNavigation:     "↑↓, kj: navegar | enter: conectar/desconectar | d/x: olvidar | q: salir",
	HelpActions:        "↑↓, kj: navegar | enter: desconectar | d/x: olvidar",
	HelpAdapterControl: "s: escaneo | p: encendido | v: descubrible | b: pairable | l: idioma | r: refrescar",
	HelpScroll:         "RePág/AvPág: página | Ctrl+↑↓, kj: scroll | Inicio/Fin: arriba/abajo | Rueda ratón: scroll",
	HelpGeneral:        "q: salir",
	HelpPairing:        "enter: confirmar | n/esc: cancelar | q: salir",
	HelpCollapsed:      "?: mostrar ayuda | q: salir",
	HelpExpanded:       "?: ocultar ayuda",

	// Adapter table
	AdapterName:         "Nombre",
	AdapterAlias:        "Alias",
	AdapterPower:        "Energía",
	AdapterPairable:     "Pairable",
	AdapterDiscoverable: "Descubrible",

	// Device table
	DeviceIcon:    "Icono",
	DeviceName:    "Nombre",
	DeviceAddress: "Dirección",
	DeviceRSSI:    "Señal",
	DeviceBattery: "Batería",
	DeviceStatus:  "Estado",

	// Status
	StatusOn:  "ON",
	StatusOff: "OFF",

	// Badges
	BadgePaired:    "PAREADO",
	BadgeConnected: "CONECTADO",
	BadgeTrusted:   "Confiable",

	// Error messages
	ErrorDBusConnection:         "No se pudo conectar a DBus",
	ErrorAdapterNotFound:        "No se encontró adaptador Bluetooth",
	ErrorStartDiscovery:         "No se pudo iniciar descubrimiento",
	ErrorStopDiscovery:          "No se pudo detener descubrimiento",
	ErrorRemoveDevice:           "No se pudo eliminar dispositivo",
	ErrorPairDevice:             "Error al parear dispositivo",
	ErrorTrustDevice:            "Error al confiar en dispositivo",
	ErrorConnectDevice:          "Error al conectar dispositivo",
	ErrorDisconnectDevice:       "Error al desconectar dispositivo",
	ErrorGetDevices:             "Error al obtener dispositivos",
	ErrorGetAdapterInfo:         "Error al obtener info del adaptador",
	ErrorSetAdapterPowered:      "Error al cambiar estado del adaptador",
	ErrorSetAdapterDiscoverable: "Error al cambiar modo discoverable",
	ErrorSetAdapterPairable:     "Error al cambiar modo pairable",
	ErrorSetAdapterAlias:        "Error al cambiar alias del adaptador",
	ErrorForgetDevice:           "Error al olvidar dispositivo",
	ErrorChangeProperty:         "Error al cambiar",

	// Status messages
	StatusConfirmingPairing:  "Confirmando pairing...",
	StatusLoadingAdapterInfo: "Cargando información del adaptador...",

	// Warnings
	WarningAgentRegistration:       "⚠️  Advertencia: No se pudo registrar agente de pairing",
	WarningAgentRegistrationDetail: "   La app funcionará pero algunos dispositivos pueden requerir pairing manual.",

	// Agent errors (internal)
	ErrorRequestPasskey:   "No se puede solicitar passkey en TUI",
	ErrorPairingCancelled: "Pairing cancelado por el usuario",
	ErrorConfirmRejected:  "Confirmación rechazada",
	ErrorExportAgent:      "No se pudo exportar agente",
	ErrorExportIntrospect: "No se pudo exportar introspección",
	ErrorRegisterAgent:    "No se pudo registrar agente",
}
