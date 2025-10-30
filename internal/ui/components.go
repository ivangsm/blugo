package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ivangsm/gob/internal/i18n"
	"github.com/ivangsm/gob/internal/models"
)

// renderHeader renderiza el encabezado de la aplicación.
func (m Model) renderHeader() string {
	title := TitleStyle.Render(i18n.T.AppTitle)

	scanStatus := ""
	if m.scanning {
		scanStatus = ScanningBadgeStyle.Render(i18n.T.Scanning)
	} else {
		scanStatus = MutedStyle.Render(i18n.T.Paused)
	}

	// Usar ancho efectivo
	effectiveWidth := min(m.width, 140)

	// Centrar el título y alinear status a la derecha
	if effectiveWidth > 0 {
		titleWidth := lipgloss.Width(title)
		statusWidth := lipgloss.Width(scanStatus)
		spacing := max(1, effectiveWidth-titleWidth-statusWidth-4)

		header := lipgloss.JoinHorizontal(
			lipgloss.Top,
			title,
			strings.Repeat(" ", spacing),
			scanStatus,
		)

		return HeaderBoxStyle.Width(effectiveWidth - 2).Render(header)
	}

	return HeaderBoxStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", scanStatus))
}

// renderFooter renderiza el pie de página con ayuda.
func (m Model) renderFooter() string {
	var helpText string

	if m.pairingPasskey != nil {
		helpText = HelpStyle.Render(i18n.T.HelpPairing)
	} else if m.focusSection == "connected" {
		helpText = HelpStyle.Render(i18n.T.HelpActions + "\n" + i18n.T.HelpAdapterControl)
	} else {
		helpText = HelpStyle.Render(i18n.T.HelpNavigation + "\n" + i18n.T.HelpAdapterControl)
	}

	// Usar ancho efectivo
	effectiveWidth := min(m.width, 140)

	if effectiveWidth > 0 {
		return FooterBoxStyle.Width(effectiveWidth - 2).Render(helpText)
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

	// Usar ancho efectivo
	effectiveWidth := min(m.width, 140)

	if effectiveWidth > 0 {
		return BoxStyle.Width(effectiveWidth - 4).Align(lipgloss.Center).Render(styled)
	}

	return BoxStyle.Render(styled)
}

// renderPasskeyPrompt renderiza el prompt de passkey.
func (m Model) renderPasskeyPrompt() string {
	passkeyText := fmt.Sprintf(i18n.T.PairingCode, *m.pairingPasskey)
	instruction := WarningStyle.Render(i18n.T.PairingInstruction)
	confirm := HelpStyle.Render(i18n.T.PairingConfirm)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		passkeyText,
		"",
		instruction,
		confirm,
	)

	// Usar ancho efectivo
	effectiveWidth := min(m.width, 140)

	if effectiveWidth > 0 {
		return PasskeyBoxStyle.Width(min(effectiveWidth-4, 70)).Render(content)
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
		badges = append(badges, PairedBadgeStyle.Render(i18n.T.BadgePaired))
	}
	if dev.Trusted {
		badges = append(badges, SuccessStyle.Render(i18n.T.BadgeTrusted))
	}
	if dev.Connected {
		badges = append(badges, ConnectedBadgeStyle.Render(i18n.T.BadgeConnected))
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

// getEmptyDevicesMessage returns the appropriate empty state message
func getEmptyAvailableDevicesMessage() string {
	return i18n.T.NoDevicesAvailable
}

// getEmptyConnectedDevicesMessage returns the appropriate empty state message
func getEmptyConnectedDevicesMessage() string {
	return i18n.T.NoDevicesConnected
}

// renderSeparator renderiza un separador.
func (m Model) renderSeparator() string {
	// Usar ancho efectivo
	effectiveWidth := min(m.width, 140)

	if effectiveWidth > 0 {
		return SeparatorStyle.Render(strings.Repeat("─", effectiveWidth-4))
	}
	return SeparatorStyle.Render(strings.Repeat("─", 80))
}

// renderThickSeparator renderiza un separador grueso.
func (m Model) renderThickSeparator() string {
	// Usar ancho efectivo
	effectiveWidth := min(m.width, 140)

	if effectiveWidth > 0 {
		return ThickSeparatorStyle.Render(strings.Repeat("━", effectiveWidth-4))
	}
	return ThickSeparatorStyle.Render(strings.Repeat("━", 80))
}

// renderAdapterTable renderiza la tabla de información del adaptador.
func (m Model) renderAdapterTable() string {
	if m.adapter == nil {
		return BoxStyle.Render(MutedStyle.Render("Loading adapter information..."))
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
		headerCellStyle.Render(i18n.T.AdapterName),
		headerCellStyle.Render(i18n.T.AdapterAlias),
		headerCellStyle.Render(i18n.T.AdapterPower),
		headerCellStyle.Render(i18n.T.AdapterPairable),
		headerCellStyle.Render(i18n.T.AdapterDiscoverable),
	}

	// Values
	nameVal := valueCellStyle.Render(m.adapter.Name)
	aliasVal := valueCellStyle.Render(m.adapter.Alias)

	// Power con color
	var powerVal string
	if m.adapter.Powered {
		powerVal = valueCellStyle.Render(SuccessStyle.Render(i18n.T.StatusOn))
	} else {
		powerVal = valueCellStyle.Render(ErrorStyle.Render(i18n.T.StatusOff))
	}

	// Pairable con color
	var pairableVal string
	if m.adapter.Pairable {
		pairableVal = valueCellStyle.Render(SuccessStyle.Render(i18n.T.StatusOn))
	} else {
		pairableVal = valueCellStyle.Render(MutedStyle.Render(i18n.T.StatusOff))
	}

	// Discoverable con color
	var discoverableVal string
	if m.adapter.Discoverable {
		discoverableVal = valueCellStyle.Render(SuccessStyle.Render(i18n.T.StatusOn))
	} else {
		discoverableVal = valueCellStyle.Render(MutedStyle.Render(i18n.T.StatusOff))
	}

	values := []string{nameVal, aliasVal, powerVal, pairableVal, discoverableVal}

	// Crear filas
	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headers...)
	valueRow := lipgloss.JoinHorizontal(lipgloss.Top, values...)
	separator := SeparatorStyle.Render(strings.Repeat("─", lipgloss.Width(headerRow)))

	table := lipgloss.JoinVertical(
		lipgloss.Left,
		HeaderStyle.Render(fmt.Sprintf("%s %s", m.adapter.GetStatusIcon(), i18n.T.AdapterInfo)),
		separator,
		headerRow,
		valueRow,
	)

	// Usar ancho efectivo
	effectiveWidth := min(m.width, 140)

	if effectiveWidth > 0 {
		return BoxStyle.Width(effectiveWidth - 4).Render(table)
	}

	return BoxStyle.Render(table)
}
