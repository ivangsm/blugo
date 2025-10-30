package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ivangsm/blugo/internal/i18n"
	"github.com/ivangsm/blugo/internal/models"
)

// View renderiza la interfaz de usuario.
func (m Model) View() string {
	if m.err != nil {
		return m.renderErrorView()
	}

	if m.manager == nil {
		return m.renderLoadingView()
	}

	// Calcular ancho efectivo con límite máximo para terminales muy grandes
	maxWidth := min(m.width, 140)
	leftPadding := (m.width - maxWidth) / 2

	// Layout principal
	sections := []string{}

	// Header
	sections = append(sections, m.renderHeader())

	// Passkey prompt (si existe)
	if m.pairingPasskey != nil {
		sections = append(sections, "", m.renderPasskeyPrompt(), "")
	}

	// Status bar (si existe)
	if m.statusMessage != "" {
		sections = append(sections, "", m.renderStatusBar())
	}

	sections = append(sections, "")

	// Tabla de información del adaptador (siempre visible)
	sections = append(sections, m.renderAdapterTable())

	sections = append(sections, "")

	// Contenido principal - siempre una columna
	sections = append(sections, m.renderSingleColumnLayout())

	sections = append(sections, "")

	// Footer
	sections = append(sections, m.renderFooter())

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Centrar contenido en terminales muy grandes
	if m.width > 140 {
		return lipgloss.NewStyle().
			PaddingLeft(leftPadding).
			Render(content)
	}

	return content
}

// renderErrorView renderiza la vista de error.
func (m Model) renderErrorView() string {
	title := ErrorStyle.Render("\n❌ " + i18n.T.Error)
	message := m.err.Error()
	help := HelpStyle.Render("\n" + i18n.T.HelpGeneral + "\n")

	return lipgloss.JoinVertical(lipgloss.Left, title, "", message, help)
}

// renderLoadingView renderiza la vista de carga.
func (m Model) renderLoadingView() string {
	return TitleStyle.Render("⚙ " + i18n.T.Initializing) + "\n"
}

// renderSingleColumnLayout renderiza el layout de una columna.
func (m Model) renderSingleColumnLayout() string {
	sections := []string{}

	// Sección de dispositivos disponibles
	sections = append(sections, m.renderFoundDevicesSection())

	sections = append(sections, "")

	// Sección de dispositivos conectados
	sections = append(sections, m.renderConnectedDevicesSection())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderTwoColumnLayout renderiza el layout de dos columnas.
func (m Model) renderTwoColumnLayout() string {
	leftColumn := m.renderFoundDevicesSection()
	rightColumn := m.renderConnectedDevicesSection()

	// Usar ancho efectivo con límite máximo
	effectiveWidth := min(m.width, 160)

	// Calcular el ancho de cada columna con espacio entre ellas
	columnWidth := (effectiveWidth - 8) / 2

	// Aplicar el ancho a las columnas
	leftStyled := lipgloss.NewStyle().Width(columnWidth).Render(leftColumn)
	rightStyled := lipgloss.NewStyle().Width(columnWidth).Render(rightColumn)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyled,
		"    ", // Más espacio entre columnas
		rightStyled,
	)
}

// renderFoundDevicesSection renderiza la sección de dispositivos disponibles.
func (m Model) renderFoundDevicesSection() string {
	foundDevices := m.GetFoundDevices()
	isFocused := m.focusSection == "found"

	sections := []string{}

	// Header
	header := renderSectionHeader("📡", i18n.T.AvailableDevices, len(foundDevices), isFocused)
	sections = append(sections, header)
	sections = append(sections, m.renderSeparator())

	// Lista de dispositivos
	if len(foundDevices) == 0 {
		sections = append(sections, renderEmptyState(getEmptyAvailableDevicesMessage()))
	} else {
		deviceList := m.renderFoundDevicesList(foundDevices, isFocused)
		sections = append(sections, deviceList)
	}

	// Aplicar borde de panel si está enfocado
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Usar ancho efectivo consistente con otros componentes
	effectiveWidth := min(m.width, 140)

	if isFocused {
		return FocusedPanelStyle.Width(effectiveWidth - 4).Render(content)
	}
	return PanelStyle.Width(effectiveWidth - 4).Render(content)
}

// renderConnectedDevicesSection renderiza la sección de dispositivos conectados.
func (m Model) renderConnectedDevicesSection() string {
	connectedDevices := m.GetConnectedDevices()
	isFocused := m.focusSection == "connected"

	sections := []string{}

	// Header
	header := renderSectionHeader("🔗", i18n.T.ConnectedDevices, len(connectedDevices), isFocused)
	sections = append(sections, header)
	sections = append(sections, m.renderSeparator())

	// Lista de dispositivos
	if len(connectedDevices) == 0 {
		sections = append(sections, renderEmptyState(getEmptyConnectedDevicesMessage()))
	} else {
		deviceList := m.renderConnectedDevicesList(connectedDevices, isFocused)
		sections = append(sections, deviceList)
	}

	// Aplicar borde de panel si está enfocado
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Usar ancho efectivo consistente con otros componentes
	effectiveWidth := min(m.width, 140)

	if isFocused {
		return FocusedPanelStyle.Width(effectiveWidth - 4).Render(content)
	}
	return PanelStyle.Width(effectiveWidth - 4).Render(content)
}

// renderFoundDevicesList renderiza la lista de dispositivos disponibles.
func (m Model) renderFoundDevicesList(devices []*models.Device, isFocused bool) string {
	items := []string{}

	for i, dev := range devices {
		isSelected := isFocused && i == m.selectedIndex
		item := renderDeviceItem(dev, isSelected, true)
		items = append(items, item)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderConnectedDevicesList renderiza la lista de dispositivos conectados.
func (m Model) renderConnectedDevicesList(devices []*models.Device, isFocused bool) string {
	items := []string{}

	for i, dev := range devices {
		isSelected := isFocused && i == m.selectedIndex
		item := renderDeviceItem(dev, isSelected, false)
		items = append(items, item)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
