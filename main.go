package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

// --- Estilos ---
var (
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Padding(0, 1)
	headerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	deviceStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Padding(0, 2)
	connectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Padding(0, 2)
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("57")).Bold(true).Padding(0, 2)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	separatorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	statusStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
	connectingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	warningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	passkeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true).Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("51"))
)

const (
	bluezService      = "org.bluez"
	bluezAdapterIface = "org.bluez.Adapter1"
	bluezDeviceIface  = "org.bluez.Device1"
	bluezAgentIface   = "org.bluez.Agent1"
)

var agentPath = dbus.ObjectPath("/org/bluez/agent_gob")

// --- Tipos de dispositivos ---
type Device struct {
	Path       dbus.ObjectPath
	Address    string
	Name       string
	Alias      string
	Paired     bool
	Trusted    bool
	Connected  bool
	RSSI       int16
	Icon       string
	Class      uint32
	LastSeen   time.Time
}

// --- Mensajes ---
type initMsg struct {
	conn    *dbus.Conn
	adapter dbus.ObjectPath
	err     error
}

type scanningMsg struct {
	scanning bool
}

type deviceUpdateMsg struct {
	devices map[string]*Device
}

type statusMsg struct {
	message string
	isError bool
}

type connectResultMsg struct {
	address string
	success bool
	err     error
}

type pairResultMsg struct {
	address string
	success bool
	err     error
}

type passkeyDisplayMsg struct {
	passkey uint32
	device  string
}

type passkeyConfirmedMsg struct{}

type tickMsg time.Time

// --- Agente de Bluetooth ---
type Agent struct {
	program *tea.Program
}

const agentIntrospection = `
<node>
	<interface name="org.bluez.Agent1">
		<method name="Release"></method>
		<method name="RequestPinCode">
			<arg type="o" direction="in"/>
			<arg type="s" direction="out"/>
		</method>
		<method name="DisplayPinCode">
			<arg type="o" direction="in"/>
			<arg type="s" direction="in"/>
		</method>
		<method name="RequestPasskey">
			<arg type="o" direction="in"/>
			<arg type="u" direction="out"/>
		</method>
		<method name="DisplayPasskey">
			<arg type="o" direction="in"/>
			<arg type="u" direction="in"/>
			<arg type="q" direction="in"/>
		</method>
		<method name="RequestConfirmation">
			<arg type="o" direction="in"/>
			<arg type="u" direction="in"/>
		</method>
		<method name="RequestAuthorization">
			<arg type="o" direction="in"/>
		</method>
		<method name="AuthorizeService">
			<arg type="o" direction="in"/>
			<arg type="s" direction="in"/>
		</method>
		<method name="Cancel"></method>
	</interface>
	<interface name="org.freedesktop.DBus.Introspectable">
		<method name="Introspect">
			<arg name="data" type="s" direction="out"/>
		</method>
	</interface>
</node>
`

var passkeyChannel chan uint32
var confirmChannel chan bool

func (a *Agent) Release() *dbus.Error {
	return nil
}

func (a *Agent) RequestPinCode(device dbus.ObjectPath) (string, *dbus.Error) {
	// PIN antiguo - devolver un PIN por defecto
	return "0000", nil
}

func (a *Agent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	// Mostrar PIN en la interfaz
	if a.program != nil {
		a.program.Send(statusMsg{message: fmt.Sprintf("PIN: %s", pincode), isError: false})
	}
	return nil
}

func (a *Agent) RequestPasskey(device dbus.ObjectPath) (uint32, *dbus.Error) {
	// Solicitar passkey - no implementado en TUI
	return 0, dbus.MakeFailedError(fmt.Errorf("no se puede solicitar passkey en TUI"))
}

func (a *Agent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	// Mostrar passkey de 6 dígitos - común para teclados
	if a.program != nil {
		// Enviar el passkey a la UI
		passkeyChannel <- passkey
	}

	// Esperar confirmación del usuario
	confirmed := <-confirmChannel
	if !confirmed {
		return dbus.MakeFailedError(fmt.Errorf("pairing cancelado por el usuario"))
	}

	return nil
}

func (a *Agent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	// Confirmar passkey - mostrar y esperar confirmación
	if a.program != nil {
		passkeyChannel <- passkey
	}

	confirmed := <-confirmChannel
	if !confirmed {
		return dbus.MakeFailedError(fmt.Errorf("confirmación rechazada"))
	}

	return nil
}

func (a *Agent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	// Autorizar dispositivo - aceptar automáticamente
	return nil
}

func (a *Agent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	// Autorizar servicio - aceptar automáticamente
	return nil
}

func (a *Agent) Cancel() *dbus.Error {
	// Cancelar pairing
	if confirmChannel != nil {
		select {
		case confirmChannel <- false:
		default:
		}
	}
	return nil
}

func unregisterAgent(conn *dbus.Conn) {
	obj := conn.Object(bluezService, dbus.ObjectPath("/org/bluez"))
	_ = obj.Call("org.bluez.AgentManager1.UnregisterAgent", 0, agentPath).Err
}

func registerAgent(conn *dbus.Conn, program *tea.Program) error {
	agent := &Agent{program: program}

	// Primero intentar desregistrar cualquier agente existente
	unregisterAgent(conn)

	// Exportar el agente en DBus (no necesitamos RequestName)
	// El agente funciona perfectamente exportándose directamente en una ruta
	err := conn.Export(agent, agentPath, bluezAgentIface)
	if err != nil {
		return fmt.Errorf("no se pudo exportar agente: %w", err)
	}

	// Exportar introspección
	err = conn.Export(introspect.Introspectable(agentIntrospection), agentPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return fmt.Errorf("no se pudo exportar introspección: %w", err)
	}

	// Registrar el agente con BlueZ - probar diferentes capacidades
	obj := conn.Object(bluezService, dbus.ObjectPath("/org/bluez"))

	// Intentar con KeyboardDisplay primero
	err = obj.Call("org.bluez.AgentManager1.RegisterAgent", 0, agentPath, "KeyboardDisplay").Err
	if err != nil {
		// Si falla, intentar con NoInputNoOutput (más permisivo)
		err = obj.Call("org.bluez.AgentManager1.RegisterAgent", 0, agentPath, "NoInputNoOutput").Err
		if err != nil {
			return fmt.Errorf("no se pudo registrar agente (probé KeyboardDisplay y NoInputNoOutput): %w", err)
		}
	}

	// Establecer como agente por defecto (no crítico)
	_ = obj.Call("org.bluez.AgentManager1.RequestDefaultAgent", 0, agentPath).Err

	return nil
}

func waitForPasskey() tea.Msg {
	passkey := <-passkeyChannel
	return passkeyDisplayMsg{passkey: passkey}
}

// --- Modelo ---
type model struct {
	conn               *dbus.Conn
	adapter            dbus.ObjectPath
	devices            map[string]*Device
	selectedIndex      int
	focusSection       string // "found" o "connected"
	statusMessage      string
	isError            bool
	scanning           bool
	busy               bool
	err                error
	pairingPasskey     *uint32
	waitingForPasskey  bool
}

// --- Helpers de DBus ---
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

func startDiscovery(conn *dbus.Conn, adapter dbus.ObjectPath) error {
	obj := conn.Object(bluezService, adapter)
	return obj.Call(bluezAdapterIface+".StartDiscovery", 0).Err
}

func stopDiscovery(conn *dbus.Conn, adapter dbus.ObjectPath) error {
	obj := conn.Object(bluezService, adapter)
	return obj.Call(bluezAdapterIface+".StopDiscovery", 0).Err
}

func getDevices(conn *dbus.Conn) (map[string]*Device, error) {
	obj := conn.Object(bluezService, "/")
	var paths map[string]map[string]map[string]dbus.Variant
	err := obj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&paths)
	if err != nil {
		return nil, err
	}

	devices := make(map[string]*Device)
	for path, interfaces := range paths {
		if props, ok := interfaces[bluezDeviceIface]; ok {
			dev := &Device{
				Path:     dbus.ObjectPath(path),
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

			// Usar Alias si no hay Name
			if dev.Name == "" && dev.Alias != "" {
				dev.Name = dev.Alias
			}

			devices[dev.Address] = dev
		}
	}

	return devices, nil
}

func pairDevice(conn *dbus.Conn, devicePath dbus.ObjectPath) error {
	obj := conn.Object(bluezService, devicePath)
	call := obj.Call(bluezDeviceIface+".Pair", 0)
	return call.Err
}

func trustDevice(conn *dbus.Conn, devicePath dbus.ObjectPath) error {
	obj := conn.Object(bluezService, devicePath)
	return obj.Call("org.freedesktop.DBus.Properties.Set", 0,
		bluezDeviceIface, "Trusted", dbus.MakeVariant(true)).Err
}

func connectDevice(conn *dbus.Conn, devicePath dbus.ObjectPath) error {
	obj := conn.Object(bluezService, devicePath)
	call := obj.Call(bluezDeviceIface+".Connect", 0)
	return call.Err
}

func disconnectDevice(conn *dbus.Conn, devicePath dbus.ObjectPath) error {
	obj := conn.Object(bluezService, devicePath)
	call := obj.Call(bluezDeviceIface+".Disconnect", 0)
	return call.Err
}

func removeDevice(conn *dbus.Conn, adapter dbus.ObjectPath, devicePath dbus.ObjectPath) error {
	obj := conn.Object(bluezService, adapter)
	call := obj.Call(bluezAdapterIface+".RemoveDevice", 0, devicePath)
	return call.Err
}

func getDeviceIcon(icon string, class uint32) string {
	// Íconos según el tipo de dispositivo
	if icon != "" {
		switch icon {
		case "audio-card", "audio-headset", "audio-headphones":
			return "🎧"
		case "phone", "smartphone":
			return "📱"
		case "computer", "laptop":
			return "💻"
		case "input-keyboard":
			return "⌨️"
		case "input-mouse":
			return "🖱️"
		case "input-gaming":
			return "🎮"
		case "camera":
			return "📷"
		case "printer":
			return "🖨️"
		}
	}

	// Fallback basado en clase de dispositivo
	majorClass := (class >> 8) & 0x1F
	switch majorClass {
	case 1: // Computer
		return "💻"
	case 2: // Phone
		return "📱"
	case 4: // Audio/Video
		return "🎧"
	case 5: // Peripheral (keyboard, mouse, etc)
		return "⌨️"
	case 6: // Imaging (printer, camera)
		return "📷"
	}

	return "📶"
}

// --- Comandos ---
func initialize(program *tea.Program) tea.Msg {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return initMsg{err: fmt.Errorf("no se pudo conectar a DBus: %w", err)}
	}

	adapter, err := getAdapter(conn)
	if err != nil {
		return initMsg{err: fmt.Errorf("no se encontró adaptador Bluetooth: %w", err)}
	}

	// Intentar registrar agente para pairing (opcional)
	err = registerAgent(conn, program)
	if err != nil {
		// No es crítico - la app funcionará pero puede necesitar confirmación manual para algunos dispositivos
		fmt.Fprintf(os.Stderr, "⚠️  Advertencia: No se pudo registrar agente de pairing: %v\n", err)
		fmt.Fprintf(os.Stderr, "   La app funcionará pero algunos dispositivos pueden requerir pairing manual.\n")
	} else {
		fmt.Fprintf(os.Stderr, "✓ Agente de pairing registrado exitosamente\n")
	}

	// Iniciar descubrimiento
	err = startDiscovery(conn, adapter)
	if err != nil {
		return initMsg{err: fmt.Errorf("no se pudo iniciar descubrimiento: %w", err)}
	}

	return initMsg{conn: conn, adapter: adapter}
}

func toggleScanning(conn *dbus.Conn, adapter dbus.ObjectPath, currentlyScanning bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		if currentlyScanning {
			err = stopDiscovery(conn, adapter)
		} else {
			err = startDiscovery(conn, adapter)
		}

		if err != nil {
			return statusMsg{message: fmt.Sprintf("Error al cambiar escaneo: %s", err), isError: true}
		}

		return scanningMsg{scanning: !currentlyScanning}
	}
}

func updateDevices(conn *dbus.Conn) tea.Cmd {
	return func() tea.Msg {
		devices, err := getDevices(conn)
		if err != nil {
			return statusMsg{message: fmt.Sprintf("Error al obtener dispositivos: %s", err), isError: true}
		}
		return deviceUpdateMsg{devices: devices}
	}
}

func connectToDevice(conn *dbus.Conn, dev *Device) tea.Cmd {
	return func() tea.Msg {
		// Si no está pareado, intentar pairing
		if !dev.Paired {
			err := pairDevice(conn, dev.Path)
			if err != nil {
				return connectResultMsg{address: dev.Address, success: false, err: fmt.Errorf("error al parear: %w", err)}
			}

			// Confiar en el dispositivo
			_ = trustDevice(conn, dev.Path)

			// Esperar un momento después del pairing
			time.Sleep(1 * time.Second)
		}

		// Conectar
		err := connectDevice(conn, dev.Path)
		if err != nil {
			return connectResultMsg{address: dev.Address, success: false, err: fmt.Errorf("error al conectar: %w", err)}
		}

		return connectResultMsg{address: dev.Address, success: true}
	}
}

func disconnectFromDevice(conn *dbus.Conn, dev *Device) tea.Cmd {
	return func() tea.Msg {
		err := disconnectDevice(conn, dev.Path)
		if err != nil {
			return connectResultMsg{address: dev.Address, success: false, err: fmt.Errorf("error al desconectar: %w", err)}
		}
		return connectResultMsg{address: dev.Address, success: true}
	}
}

func forgetDevice(conn *dbus.Conn, adapter dbus.ObjectPath, dev *Device) tea.Cmd {
	return func() tea.Msg {
		// Desconectar primero si está conectado
		if dev.Connected {
			_ = disconnectDevice(conn, dev.Path)
			time.Sleep(500 * time.Millisecond)
		}

		// Eliminar el dispositivo
		err := removeDevice(conn, adapter, dev.Path)
		if err != nil {
			return statusMsg{message: fmt.Sprintf("Error al olvidar: %s", err), isError: true}
		}

		return statusMsg{message: fmt.Sprintf("Dispositivo %s olvidado", dev.Name), isError: false}
	}
}

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// --- Modelo inicial ---
func initialModel() model {
	return model{
		devices:       make(map[string]*Device),
		focusSection:  "found",
		selectedIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
	)
}

// --- Helpers ---
func (m model) getFoundDevices() []*Device {
	devices := make([]*Device, 0)
	for _, dev := range m.devices {
		if !dev.Connected {
			devices = append(devices, dev)
		}
	}
	sort.Slice(devices, func(i, j int) bool {
		// Pareados primero
		if devices[i].Paired != devices[j].Paired {
			return devices[i].Paired
		}
		// Luego por RSSI
		return devices[i].RSSI > devices[j].RSSI
	})
	return devices
}

func (m model) getConnectedDevices() []*Device {
	devices := make([]*Device, 0)
	for _, dev := range m.devices {
		if dev.Connected {
			devices = append(devices, dev)
		}
	}
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].LastSeen.Before(devices[j].LastSeen)
	})
	return devices
}

// --- Update ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Si estamos esperando confirmación de passkey
		if m.pairingPasskey != nil {
			switch msg.String() {
			case "enter", "y":
				// Confirmar pairing
				confirmChannel <- true
				m.pairingPasskey = nil
				m.statusMessage = "Confirmando pairing..."
				return m, nil
			case "n", "esc":
				// Cancelar pairing
				confirmChannel <- false
				m.pairingPasskey = nil
				m.busy = false
				m.statusMessage = "Pairing cancelado"
				return m, nil
			case "ctrl+c", "q":
				if confirmChannel != nil {
					confirmChannel <- false
				}
				if m.conn != nil && m.scanning {
					stopDiscovery(m.conn, m.adapter)
				}
				return m, tea.Quit
			}
			return m, nil
		}

		if m.busy {
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				if m.conn != nil && m.scanning {
					stopDiscovery(m.conn, m.adapter)
				}
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if m.conn != nil && m.scanning {
				stopDiscovery(m.conn, m.adapter)
			}
			return m, tea.Quit

		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}

		case "down", "j":
			maxIndex := 0
			if m.focusSection == "found" {
				maxIndex = len(m.getFoundDevices()) - 1
			} else {
				maxIndex = len(m.getConnectedDevices()) - 1
			}
			if m.selectedIndex < maxIndex {
				m.selectedIndex++
			}

		case "tab":
			if m.focusSection == "found" && len(m.getConnectedDevices()) > 0 {
				m.focusSection = "connected"
				m.selectedIndex = 0
			} else if m.focusSection == "connected" && len(m.getFoundDevices()) > 0 {
				m.focusSection = "found"
				m.selectedIndex = 0
			}

		case "s":
			if m.conn != nil {
				return m, toggleScanning(m.conn, m.adapter, m.scanning)
			}

		case "enter":
			if m.conn == nil {
				return m, nil
			}

			if m.focusSection == "found" {
				devices := m.getFoundDevices()
				if m.selectedIndex < len(devices) {
					dev := devices[m.selectedIndex]
					m.busy = true
					if dev.Paired {
						m.statusMessage = fmt.Sprintf("Conectando a %s...", dev.Name)
					} else {
						m.statusMessage = fmt.Sprintf("Pareando %s...", dev.Name)
						m.waitingForPasskey = true
					}
					return m, tea.Batch(
						connectToDevice(m.conn, dev),
						waitForPasskey,
					)
				}
			} else if m.focusSection == "connected" {
				// Enter solo desconecta pero mantiene el pairing
				devices := m.getConnectedDevices()
				if m.selectedIndex < len(devices) {
					dev := devices[m.selectedIndex]
					m.busy = true
					m.statusMessage = fmt.Sprintf("Desconectando de %s (manteniendo pairing)...", dev.Name)
					return m, disconnectFromDevice(m.conn, dev)
				}
			}

		case "d", "x":
			// D/X desconecta (si está conectado) y olvida el dispositivo
			if m.conn == nil {
				return m, nil
			}

			if m.focusSection == "connected" {
				devices := m.getConnectedDevices()
				if m.selectedIndex < len(devices) {
					dev := devices[m.selectedIndex]
					m.busy = true
					m.statusMessage = fmt.Sprintf("Desconectando y olvidando %s...", dev.Name)
					return m, forgetDevice(m.conn, m.adapter, dev)
				}
			} else if m.focusSection == "found" {
				devices := m.getFoundDevices()
				if m.selectedIndex < len(devices) {
					dev := devices[m.selectedIndex]
					if dev.Paired {
						m.busy = true
						m.statusMessage = fmt.Sprintf("Olvidando %s...", dev.Name)
						return m, forgetDevice(m.conn, m.adapter, dev)
					}
				}
			}

		case "r":
			if m.conn != nil {
				return m, updateDevices(m.conn)
			}
		}

	case initMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.conn = msg.conn
		m.adapter = msg.adapter
		m.scanning = true
		m.statusMessage = "Escaneando dispositivos Bluetooth..."
		return m, updateDevices(m.conn)

	case scanningMsg:
		m.scanning = msg.scanning
		if msg.scanning {
			m.statusMessage = "Escaneo activado"
		} else {
			m.statusMessage = "Escaneo pausado"
		}

	case deviceUpdateMsg:
		// Actualizar solo dispositivos nuevos o modificados
		for addr, newDev := range msg.devices {
			if oldDev, exists := m.devices[addr]; exists {
				// Mantener LastSeen si el dispositivo ya existía
				if !oldDev.Connected && !newDev.Connected {
					newDev.LastSeen = oldDev.LastSeen
				}
			}
			m.devices[addr] = newDev
		}

	case passkeyDisplayMsg:
		if m.waitingForPasskey {
			m.pairingPasskey = &msg.passkey
			m.waitingForPasskey = false
		}

	case connectResultMsg:
		m.busy = false
		m.waitingForPasskey = false
		if msg.err != nil {
			m.statusMessage = fmt.Sprintf("❌ %s", msg.err.Error())
			m.isError = true
		} else {
			if dev, ok := m.devices[msg.address]; ok {
				if dev.Connected {
					m.statusMessage = fmt.Sprintf("✓ Conectado a %s", dev.Name)
				} else {
					// Verificar si sigue pareado
					if dev.Paired {
						m.statusMessage = fmt.Sprintf("✓ Desconectado de %s (aún pareado)", dev.Name)
					} else {
						m.statusMessage = fmt.Sprintf("✓ Desconectado de %s", dev.Name)
					}
				}
			}
			m.isError = false
		}
		// Actualizar dispositivos después de conectar/desconectar
		return m, updateDevices(m.conn)

	case statusMsg:
		m.busy = false
		m.statusMessage = msg.message
		m.isError = msg.isError
		// Actualizar dispositivos después de olvidar
		return m, updateDevices(m.conn)

	case tickMsg:
		var cmds []tea.Cmd
		cmds = append(cmds, tick())
		if m.conn != nil {
			cmds = append(cmds, updateDevices(m.conn))
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// --- View ---
func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("\n❌ Error: %s\n\nPresiona 'q' para salir\n", m.err))
	}

	if m.conn == nil {
		return titleStyle.Render("⚙ Inicializando Bluetooth...") + "\n"
	}

	s := "\n"
	s += titleStyle.Render("🔵 Gestor Bluetooth (BlueZ)") + "\n\n"

	// Mostrar passkey si está activo
	if m.pairingPasskey != nil {
		s += passkeyStyle.Render(fmt.Sprintf("🔑 CÓDIGO DE PAIRING: %06d", *m.pairingPasskey)) + "\n\n"
		s += warningStyle.Render("⌨️  Escribe este código en tu teclado y presiona Enter") + "\n"
		s += helpStyle.Render("Luego presiona Enter aquí para confirmar, o Esc/N para cancelar") + "\n\n"
	}

	// Mensaje de estado
	if m.statusMessage != "" {
		if m.busy {
			s += connectingStyle.Render("⚙ "+m.statusMessage) + "\n\n"
		} else if m.isError {
			s += errorStyle.Render(m.statusMessage) + "\n\n"
		} else {
			s += statusStyle.Render(m.statusMessage) + "\n\n"
		}
	}

	// Indicador de escaneo
	scanIndicator := "⏸ Pausado"
	if m.scanning {
		scanIndicator = "🔍 Escaneando"
	}
	s += warningStyle.Render(scanIndicator) + "\n\n"

	// Sección de dispositivos encontrados
	foundDevices := m.getFoundDevices()
	focusMarker := ""
	if m.focusSection == "found" {
		focusMarker = " ◀"
	}
	s += headerStyle.Render("📡 DISPOSITIVOS DISPONIBLES"+focusMarker) + " "
	s += fmt.Sprintf("(%d)", len(foundDevices))
	s += "\n"
	s += separatorStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n"

	if len(foundDevices) == 0 {
		s += deviceStyle.Render("  No hay dispositivos disponibles") + "\n"
	} else {
		for i, dev := range foundDevices {
			icon := getDeviceIcon(dev.Icon, dev.Class)
			name := dev.Name
			if name == "" {
				name = dev.Address
			}

			pairedMarker := ""
			if dev.Paired {
				pairedMarker = " [PAREADO]"
			}

			rssiInfo := ""
			if dev.RSSI != 0 {
				rssiInfo = fmt.Sprintf(" | %d dBm", dev.RSSI)
			}

			line := fmt.Sprintf("  %s %s (%s)%s%s", icon, name, dev.Address, rssiInfo, pairedMarker)

			if m.focusSection == "found" && i == m.selectedIndex {
				s += selectedStyle.Render("▶ "+line) + "\n"
			} else {
				s += deviceStyle.Render(line) + "\n"
			}
		}
	}

	s += "\n"

	// Sección de dispositivos conectados
	connectedDevices := m.getConnectedDevices()
	focusMarker = ""
	if m.focusSection == "connected" {
		focusMarker = " ◀"
	}
	s += headerStyle.Render("🔗 DISPOSITIVOS CONECTADOS"+focusMarker) + " "
	s += fmt.Sprintf("(%d)", len(connectedDevices))
	s += "\n"
	s += separatorStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n"

	if len(connectedDevices) == 0 {
		s += connectedStyle.Render("  No hay dispositivos conectados") + "\n"
	} else {
		for i, dev := range connectedDevices {
			icon := getDeviceIcon(dev.Icon, dev.Class)
			name := dev.Name
			if name == "" {
				name = dev.Address
			}

			trustedMarker := ""
			if dev.Trusted {
				trustedMarker = " | Confiable"
			}

			line := fmt.Sprintf("  %s %s (%s)%s", icon, name, dev.Address, trustedMarker)

			if m.focusSection == "connected" && i == m.selectedIndex {
				s += selectedStyle.Render("▶ "+line) + "\n"
			} else {
				s += connectedStyle.Render(line) + "\n"
			}
		}
	}

	s += "\n"

	if m.pairingPasskey != nil {
		s += helpStyle.Render("Enter: confirmar pairing | N/Esc: cancelar | Q: salir") + "\n"
	} else {
		var helpText string
		if m.focusSection == "connected" {
			// Ayuda específica para sección de conectados
			helpText = "↑/↓: navegar | Tab: cambiar sección | Enter: desconectar | D/X: desconectar y olvidar\nS: pausar escaneo | R: refrescar | Q: salir"
		} else {
			// Ayuda para sección de disponibles
			helpText = "↑/↓: navegar | Tab: cambiar sección | Enter: conectar | D/X: olvidar pareado\nS: pausar escaneo | R: refrescar | Q: salir"
		}
		s += helpStyle.Render(helpText) + "\n"
	}

	return s
}

func main() {
	// Inicializar canales globales
	passkeyChannel = make(chan uint32, 1)
	confirmChannel = make(chan bool, 1)

	m := initialModel()
	p := tea.NewProgram(m)

	// Inicializar después de crear el programa
	go func() {
		p.Send(initialize(p))
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al ejecutar el programa: %v\n", err)
		os.Exit(1)
	}
}

