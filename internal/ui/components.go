package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ivangsm/gob/internal/models"
)

// renderHeader renderiza el encabezado de la aplicación.
func (m Model) renderHeader() string {
	title := TitleStyle.Render("🔵 GOB - Gestor Bluetooth")

	scanStatus := ""
	if m.scanning {
		scanStatus = ScanningBadgeStyle.Render("🔍 Escaneando")
	} else {
		scanStatus = MutedStyle.Render("⏸ Pausado")
	}

	// Centrar el título y alinear status a la derecha
	if m.width > 0 {
		titleWidth := lipgloss.Width(title)
		statusWidth := lipgloss.Width(scanStatus)
		spacing := max(1, m.width-titleWidth-statusWidth-4)

		header := lipgloss.JoinHorizontal(
			lipgloss.Top,
			title,
			strings.Repeat(" ", spacing),
			scanStatus,
		)

		return HeaderBoxStyle.Width(m.width - 2).Render(header)
	}

	return HeaderBoxStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", scanStatus))
}

// renderFooter renderiza el pie de página con ayuda.
func (m Model) renderFooter() string {
	var helpText string

	if m.pairingPasskey != nil {
		helpText = HelpStyle.Render("Enter: confirmar | N/Esc: cancelar | Q: salir")
	} else if m.focusSection == "connected" {
		helpText = HelpStyle.Render("↑/↓: navegar | Tab: cambiar | Enter: desconectar | D/X: olvidar\nS: escaneo | P: powered | V: discoverable | B: pairable | R: refrescar | Q: salir")
	} else {
		helpText = HelpStyle.Render("↑/↓: navegar | Tab: cambiar | Enter: conectar | D/X: olvidar\nS: escaneo | P: powered | V: discoverable | B: pairable | R: refrescar | Q: salir")
	}

	if m.width > 0 {
		return FooterBoxStyle.Width(m.width - 2).Render(helpText)
	}

	return FooterBoxStyle.Render(helpText)
}

// renderStatusBar renderiza la barra de estado.
func (m Model) renderStatusBar() string {
	if m.statusMessage == "" {
		return ""
	}

	var styled string
	if m.busy {
		styled = ConnectingStyle.Render("⚙ " + m.statusMessage)
	} else if m.isError {
		styled = ErrorStyle.Render("❌ " + m.statusMessage)
	} else {
		styled = SuccessStyle.Render("✓ " + m.statusMessage)
	}

	if m.width > 0 {
		return BoxStyle.Width(m.width - 4).Align(lipgloss.Center).Render(styled)
	}

	return BoxStyle.Render(styled)
}

// renderPasskeyPrompt renderiza el prompt de passkey.
func (m Model) renderPasskeyPrompt() string {
	passkeyText := fmt.Sprintf("🔑 CÓDIGO DE PAIRING: %06d", *m.pairingPasskey)
	instruction := WarningStyle.Render("⌨️  Escribe este código en tu teclado y presiona Enter")
	confirm := HelpStyle.Render("Luego presiona Enter aquí para confirmar, o Esc/N para cancelar")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		passkeyText,
		"",
		instruction,
		confirm,
	)

	if m.width > 0 {
		return PasskeyBoxStyle.Width(min(m.width-4, 70)).Render(content)
	}

	return PasskeyBoxStyle.Render(content)
}

// renderDeviceCount renderiza el contador de dispositivos.
func renderDeviceCount(count int) string {
	return MutedStyle.Render(fmt.Sprintf("(%d)", count))
}

// renderSectionHeader renderiza el encabezado de una sección.
func renderSectionHeader(icon, title string, count int, isFocused bool) string {
	countStr := renderDeviceCount(count)

	focusMarker := ""
	if isFocused {
		focusMarker = SelectedStyle.Render(" ◀")
	}

	header := fmt.Sprintf("%s %s %s%s", icon, title, countStr, focusMarker)

	return HeaderStyle.Render(header)
}

// renderDeviceItem renderiza un item de dispositivo.
func renderDeviceItem(dev *models.Device, isSelected bool, showRSSI bool) string {
	icon := DeviceIconStyle.Render(dev.GetIcon())
	name := DeviceNameStyle.Render(dev.GetDisplayName())
	address := DeviceAddressStyle.Render(fmt.Sprintf("(%s)", dev.Address))

	// Información adicional
	var info []string

	// RSSI
	if showRSSI && dev.RSSI != 0 {
		info = append(info, DeviceInfoStyle.Render(fmt.Sprintf("%d dBm", dev.RSSI)))
	}

	// Batería
	if dev.HasBattery() {
		battIcon, battText := dev.GetBatteryInfo()
		batteryStr := fmt.Sprintf("%s %s", battIcon, battText)
		batteryStyled := GetBatteryStyle(*dev.Battery).Render(batteryStr)
		info = append(info, batteryStyled)
	}

	// Badges
	var badges []string
	if dev.Paired {
		badges = append(badges, PairedBadgeStyle.Render("PAREADO"))
	}
	if dev.Trusted {
		badges = append(badges, SuccessStyle.Render("Confiable"))
	}
	if dev.Connected {
		badges = append(badges, ConnectedBadgeStyle.Render("CONECTADO"))
	}

	// Construir la línea
	parts := []string{icon, name, address}

	if len(info) > 0 {
		parts = append(parts, MutedStyle.Render("|"), strings.Join(info, " "+MutedStyle.Render("|")+" "))
	}

	if len(badges) > 0 {
		parts = append(parts, strings.Join(badges, " "))
	}

	line := strings.Join(parts, " ")

	if isSelected {
		return SelectedDeviceItemStyle.Render("▶ " + line)
	}

	return DeviceItemStyle.Render(line)
}

// renderEmptyState renderiza un estado vacío.
func renderEmptyState(message string) string {
	return MutedStyle.Padding(2, 4).Render(message)
}

// renderSeparator renderiza un separador.
func (m Model) renderSeparator() string {
	if m.width > 0 {
		return SeparatorStyle.Render(strings.Repeat("─", m.width-4))
	}
	return SeparatorStyle.Render(strings.Repeat("─", 80))
}

// renderThickSeparator renderiza un separador grueso.
func (m Model) renderThickSeparator() string {
	if m.width > 0 {
		return ThickSeparatorStyle.Render(strings.Repeat("━", m.width-4))
	}
	return ThickSeparatorStyle.Render(strings.Repeat("━", 80))
}

// renderAdapterTable renderiza la tabla de información del adaptador.
func (m Model) renderAdapterTable() string {
	if m.adapter == nil {
		return BoxStyle.Render(MutedStyle.Render("Cargando información del adaptador..."))
	}

	// Crear estilos para la tabla
	headerCellStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(secondaryColor).
		Width(15).
		Align(lipgloss.Left)

	valueCellStyle := lipgloss.NewStyle().
		Width(20).
		Align(lipgloss.Left)

	// Headers
	headers := []string{
		headerCellStyle.Render("Name"),
		headerCellStyle.Render("Alias"),
		headerCellStyle.Render("Power"),
		headerCellStyle.Render("Pairable"),
		headerCellStyle.Render("Discoverable"),
	}

	// Values
	nameVal := valueCellStyle.Render(m.adapter.Name)
	aliasVal := valueCellStyle.Render(m.adapter.Alias)

	// Power con color
	var powerVal string
	if m.adapter.Powered {
		powerVal = valueCellStyle.Render(SuccessStyle.Render("ON"))
	} else {
		powerVal = valueCellStyle.Render(ErrorStyle.Render("OFF"))
	}

	// Pairable con color
	var pairableVal string
	if m.adapter.Pairable {
		pairableVal = valueCellStyle.Render(SuccessStyle.Render("ON"))
	} else {
		pairableVal = valueCellStyle.Render(MutedStyle.Render("OFF"))
	}

	// Discoverable con color
	var discoverableVal string
	if m.adapter.Discoverable {
		discoverableVal = valueCellStyle.Render(SuccessStyle.Render("ON"))
	} else {
		discoverableVal = valueCellStyle.Render(MutedStyle.Render("OFF"))
	}

	values := []string{nameVal, aliasVal, powerVal, pairableVal, discoverableVal}

	// Crear filas
	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headers...)
	valueRow := lipgloss.JoinHorizontal(lipgloss.Top, values...)
	separator := SeparatorStyle.Render(strings.Repeat("─", lipgloss.Width(headerRow)))

	table := lipgloss.JoinVertical(
		lipgloss.Left,
		HeaderStyle.Render(fmt.Sprintf("%s Adaptador Bluetooth", m.adapter.GetStatusIcon())),
		separator,
		headerRow,
		valueRow,
	)

	if m.width > 0 {
		return BoxStyle.Width(m.width - 4).Render(table)
	}

	return BoxStyle.Render(table)
}
