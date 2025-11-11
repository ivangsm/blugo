package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ivangsm/blugo/internal/i18n"
	"github.com/ivangsm/blugo/internal/models"
)

// View renders the user interface.
func (m Model) View() string {
	if m.err != nil {
		return m.renderErrorView()
	}

	if m.manager == nil {
		return m.renderLoadingView()
	}

	if !m.ready {
		return "\n  Initializing..."
	}

	// Return viewport view (content is updated in Update())
	return m.viewport.View()
}

// renderFullContent generates the complete content for the viewport
func (m Model) renderFullContent() string {
	// Calculate effective width with maximum limit for very large terminals
	maxWidth := min(m.width, GetMaxWidth())
	leftPadding := (m.width - maxWidth) / 2

	// Main layout
	sections := []string{}

	// Header
	sections = append(sections, m.renderHeader())

	// Passkey prompt (if exists)
	if m.pairingPasskey != nil {
		sections = append(sections, "", m.renderPasskeyPrompt(), "")
	}

	// Status bar (if exists)
	if m.statusMessage != "" {
		sections = append(sections, "", m.renderStatusBar())
	}

	sections = append(sections, "")

	// Main content - always one column
	sections = append(sections, m.renderSingleColumnLayout())

	sections = append(sections, "")

	// Adapter information table (moved to the bottom)
	sections = append(sections, m.renderAdapterTable())

	sections = append(sections, "")

	// Footer
	sections = append(sections, m.renderFooter())

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Center content in very large terminals
	if m.width > GetMaxWidth() {
		content = lipgloss.NewStyle().
			PaddingLeft(leftPadding).
			Render(content)
	}

	return content
}

// renderErrorView renders the error view.
func (m Model) renderErrorView() string {
	errorText := i18n.T.Error
	if Emoji(EmojiError) != "" {
		errorText = Emoji(EmojiError) + " " + errorText
	}
	title := ErrorStyle.Render("\n" + errorText)
	message := m.err.Error()
	help := HelpStyle.Render("\n" + i18n.T.HelpGeneral + "\n")

	return lipgloss.JoinVertical(lipgloss.Left, title, "", message, help)
}

// renderLoadingView renders the loading view.
func (m Model) renderLoadingView() string {
	loadingText := i18n.T.Initializing
	if Emoji(EmojiLoading) != "" {
		loadingText = Emoji(EmojiLoading) + " " + loadingText
	}
	return TitleStyle.Render(loadingText) + "\n"
}

// renderSingleColumnLayout renders the single column layout.
func (m Model) renderSingleColumnLayout() string {
	sections := []string{}

	// Available devices section with table
	sections = append(sections, m.renderDevicesTable())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderTwoColumnLayout renders the two column layout.
func (m Model) renderTwoColumnLayout() string {
	leftColumn := m.renderFoundDevicesSection()
	rightColumn := m.renderConnectedDevicesSection()

	// Use effective width with maximum limit
	effectiveWidth := min(m.width, 160)

	// Calculate the width of each column with space between them
	columnWidth := (effectiveWidth - 8) / 2

	// Apply width to columns
	leftStyled := lipgloss.NewStyle().Width(columnWidth).Render(leftColumn)
	rightStyled := lipgloss.NewStyle().Width(columnWidth).Render(rightColumn)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyled,
		"    ", // More space between columns
		rightStyled,
	)
}

// renderFoundDevicesSection renders the available devices section.
func (m Model) renderFoundDevicesSection() string {
	foundDevices := m.GetFoundDevices()
	isFocused := m.focusSection == "found"

	sections := []string{}

	// Header
	header := renderSectionHeader(Emoji(EmojiAvailable), i18n.T.AvailableDevices, len(foundDevices), isFocused)
	sections = append(sections, header)
	sections = append(sections, m.renderSeparator())

	// Device list
	if len(foundDevices) == 0 {
		sections = append(sections, renderEmptyState(getEmptyAvailableDevicesMessage()))
	} else {
		deviceList := m.renderFoundDevicesList(foundDevices, isFocused)
		sections = append(sections, deviceList)
	}

	// Apply panel border if focused
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Use effective width consistent with other components
	effectiveWidth := min(m.width, GetMaxWidth())

	if isFocused {
		return FocusedPanelStyle.Width(effectiveWidth - 4).Render(content)
	}
	return PanelStyle.Width(effectiveWidth - 4).Render(content)
}

// renderConnectedDevicesSection renders the connected devices section.
func (m Model) renderConnectedDevicesSection() string {
	connectedDevices := m.GetConnectedDevices()
	isFocused := m.focusSection == "connected"

	sections := []string{}

	// Header
	header := renderSectionHeader(Emoji(EmojiConnected), i18n.T.ConnectedDevices, len(connectedDevices), isFocused)
	sections = append(sections, header)
	sections = append(sections, m.renderSeparator())

	// Device list
	if len(connectedDevices) == 0 {
		sections = append(sections, renderEmptyState(getEmptyConnectedDevicesMessage()))
	} else {
		deviceList := m.renderConnectedDevicesList(connectedDevices, isFocused)
		sections = append(sections, deviceList)
	}

	// Apply panel border if focused
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Use effective width consistent with other components
	effectiveWidth := min(m.width, GetMaxWidth())

	if isFocused {
		return FocusedPanelStyle.Width(effectiveWidth - 4).Render(content)
	}
	return PanelStyle.Width(effectiveWidth - 4).Render(content)
}

// renderFoundDevicesList renders the list of available devices.
func (m Model) renderFoundDevicesList(devices []*models.Device, isFocused bool) string {
	items := []string{}

	for i, dev := range devices {
		isSelected := isFocused && i == m.selectedIndex
		// Show RSSI only for non-connected devices
		showRSSI := !dev.Connected
		item := renderDeviceItem(dev, isSelected, showRSSI)
		items = append(items, item)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderConnectedDevicesList renders the list of connected devices.
func (m Model) renderConnectedDevicesList(devices []*models.Device, isFocused bool) string {
	items := []string{}

	for i, dev := range devices {
		isSelected := isFocused && i == m.selectedIndex
		item := renderDeviceItem(dev, isSelected, false)
		items = append(items, item)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
