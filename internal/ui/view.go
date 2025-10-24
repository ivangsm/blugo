package ui

import (
	"fmt"

	"github.com/ivangsm/gob/internal/models"
)

// View renderiza la interfaz de usuario.
func (m Model) View() string {
	if m.err != nil {
		return ErrorStyle.Render(fmt.Sprintf("\n❌ Error: %s\n\nPresiona 'q' para salir\n", m.err))
	}

	if m.manager == nil {
		return TitleStyle.Render("⚙ Inicializando Bluetooth...") + "\n"
	}

	s := "\n"
	s += TitleStyle.Render("🔵 Gestor Bluetooth (BlueZ)") + "\n\n"

	// Mostrar passkey si está activo
	if m.pairingPasskey != nil {
		s += m.renderPasskeyPrompt()
	}

	// Mensaje de estado
	if m.statusMessage != "" {
		s += m.renderStatusMessage()
	}

	// Indicador de escaneo
	s += m.renderScanIndicator()

	// Sección de dispositivos encontrados
	s += m.renderFoundDevices()

	s += "\n"

	// Sección de dispositivos conectados
	s += m.renderConnectedDevices()

	s += "\n"

	// Ayuda
	s += m.renderHelp()

	return s
}

// renderPasskeyPrompt renderiza el prompt de passkey.
func (m Model) renderPasskeyPrompt() string {
	s := PasskeyStyle.Render(fmt.Sprintf("🔑 CÓDIGO DE PAIRING: %06d", *m.pairingPasskey)) + "\n\n"
	s += WarningStyle.Render("⌨️  Escribe este código en tu teclado y presiona Enter") + "\n"
	s += HelpStyle.Render("Luego presiona Enter aquí para confirmar, o Esc/N para cancelar") + "\n\n"
	return s
}

// renderStatusMessage renderiza el mensaje de estado.
func (m Model) renderStatusMessage() string {
	if m.busy {
		return ConnectingStyle.Render("⚙ "+m.statusMessage) + "\n\n"
	} else if m.isError {
		return ErrorStyle.Render(m.statusMessage) + "\n\n"
	} else {
		return StatusStyle.Render(m.statusMessage) + "\n\n"
	}
}

// renderScanIndicator renderiza el indicador de escaneo.
func (m Model) renderScanIndicator() string {
	scanIndicator := "⏸ Pausado"
	if m.scanning {
		scanIndicator = "🔍 Escaneando"
	}
	return WarningStyle.Render(scanIndicator) + "\n\n"
}

// renderFoundDevices renderiza la sección de dispositivos disponibles.
func (m Model) renderFoundDevices() string {
	foundDevices := m.GetFoundDevices()
	focusMarker := ""
	if m.focusSection == "found" {
		focusMarker = " ◀"
	}

	s := HeaderStyle.Render("📡 DISPOSITIVOS DISPONIBLES"+focusMarker) + " "
	s += fmt.Sprintf("(%d)", len(foundDevices))
	s += "\n"
	s += SeparatorStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n"

	if len(foundDevices) == 0 {
		s += DeviceStyle.Render("  No hay dispositivos disponibles") + "\n"
	} else {
		for i, dev := range foundDevices {
			s += m.renderDevice(dev, i, m.focusSection == "found")
		}
	}

	return s
}

// renderConnectedDevices renderiza la sección de dispositivos conectados.
func (m Model) renderConnectedDevices() string {
	connectedDevices := m.GetConnectedDevices()
	focusMarker := ""
	if m.focusSection == "connected" {
		focusMarker = " ◀"
	}

	s := HeaderStyle.Render("🔗 DISPOSITIVOS CONECTADOS"+focusMarker) + " "
	s += fmt.Sprintf("(%d)", len(connectedDevices))
	s += "\n"
	s += SeparatorStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n"

	if len(connectedDevices) == 0 {
		s += ConnectedStyle.Render("  No hay dispositivos conectados") + "\n"
	} else {
		for i, dev := range connectedDevices {
			s += m.renderConnectedDevice(dev, i, m.focusSection == "connected")
		}
	}

	return s
}

// renderDevice renderiza un dispositivo disponible.
func (m Model) renderDevice(dev *models.Device, index int, isFocused bool) string {
	icon := dev.GetIcon()
	name := dev.GetDisplayName()

	pairedMarker := ""
	if dev.Paired {
		pairedMarker = " [PAREADO]"
	}

	rssiInfo := ""
	if dev.RSSI != 0 {
		rssiInfo = fmt.Sprintf(" | %d dBm", dev.RSSI)
	}

	line := fmt.Sprintf("  %s %s (%s)%s%s", icon, name, dev.Address, rssiInfo, pairedMarker)

	if isFocused && index == m.selectedIndex {
		return SelectedStyle.Render("▶ "+line) + "\n"
	}
	return DeviceStyle.Render(line) + "\n"
}

// renderConnectedDevice renderiza un dispositivo conectado.
func (m Model) renderConnectedDevice(dev *models.Device, index int, isFocused bool) string {
	icon := dev.GetIcon()
	name := dev.GetDisplayName()

	trustedMarker := ""
	if dev.Trusted {
		trustedMarker = " | Confiable"
	}

	line := fmt.Sprintf("  %s %s (%s)%s", icon, name, dev.Address, trustedMarker)

	if isFocused && index == m.selectedIndex {
		return SelectedStyle.Render("▶ "+line) + "\n"
	}
	return ConnectedStyle.Render(line) + "\n"
}

// renderHelp renderiza la ayuda.
func (m Model) renderHelp() string {
	if m.pairingPasskey != nil {
		return HelpStyle.Render("Enter: confirmar pairing | N/Esc: cancelar | Q: salir") + "\n"
	}

	var helpText string
	if m.focusSection == "connected" {
		helpText = "↑/↓: navegar | Tab: cambiar sección | Enter: desconectar | D/X: desconectar y olvidar\nS: pausar escaneo | R: refrescar | Q: salir"
	} else {
		helpText = "↑/↓: navegar | Tab: cambiar sección | Enter: conectar | D/X: olvidar pareado\nS: pausar escaneo | R: refrescar | Q: salir"
	}
	return HelpStyle.Render(helpText) + "\n"
}
