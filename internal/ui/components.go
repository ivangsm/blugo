package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/i18n"
	"github.com/ivangsm/blugo/internal/models"
)

// renderHeader renders the application header.
func (m Model) renderHeader() string {
	titleText := i18n.T.AppTitle
	if Emoji(EmojiAppTitle) != "" {
		titleText = Emoji(EmojiAppTitle) + " " + titleText
	}
	title := TitleStyle.Render(titleText)

	scanStatus := ""
	if m.scanning {
		statusText := i18n.T.Scanning
		if Emoji(EmojiScanning) != "" {
			statusText = Emoji(EmojiScanning) + " " + statusText
		}
		scanStatus = ScanningBadgeStyle.Render(statusText)
	} else {
		statusText := i18n.T.Paused
		if Emoji(EmojiPaused) != "" {
			statusText = Emoji(EmojiPaused) + " " + statusText
		}
		scanStatus = MutedStyle.Render(statusText)
	}

	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())

	// Center the title and align status to the right
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

// renderFooter renders the footer with help.
func (m Model) renderFooter() string {
	var helpText string

	if m.pairingPasskey != nil {
		helpText = HelpStyle.Render(i18n.T.HelpPairing)
	} else if m.showHelp {
		// Show full help when expanded
		helpText = HelpStyle.Render(
			i18n.T.HelpNavigation + " | " + i18n.T.HelpExpanded + "\n" +
			i18n.T.HelpAdapterControl + "\n" +
			i18n.T.HelpScroll,
		)
	} else {
		// Show collapsed help
		helpText = HelpStyle.Render(i18n.T.HelpCollapsed)
	}

	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())

	if effectiveWidth > 0 {
		return FooterBoxStyle.Width(effectiveWidth - 2).Render(helpText)
	}

	return FooterBoxStyle.Render(helpText)
}

// renderStatusBar renders the status bar.
func (m Model) renderStatusBar() string {
	if m.statusMessage == "" {
		return ""
	}

	var styled string
	if m.busy {
		text := m.statusMessage
		if Emoji(EmojiLoading) != "" {
			text = Emoji(EmojiLoading) + " " + text
		}
		styled = ConnectingStyle.Render(text)
	} else if m.isError {
		text := m.statusMessage
		if Emoji(EmojiError) != "" {
			text = Emoji(EmojiError) + " " + text
		}
		styled = ErrorStyle.Render(text)
	} else {
		text := m.statusMessage
		if Emoji(EmojiSuccess) != "" {
			text = Emoji(EmojiSuccess) + " " + text
		}
		styled = SuccessStyle.Render(text)
	}

	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())

	if effectiveWidth > 0 {
		return BoxStyle.Width(effectiveWidth - 4).Align(lipgloss.Center).Render(styled)
	}

	return BoxStyle.Render(styled)
}

// renderPasskeyPrompt renders the passkey prompt.
func (m Model) renderPasskeyPrompt() string {
	passkeyFormat := i18n.T.PairingCode
	if Emoji(EmojiPairingKey) != "" {
		passkeyFormat = Emoji(EmojiPairingKey) + " " + passkeyFormat
	}
	passkeyText := fmt.Sprintf(passkeyFormat, *m.pairingPasskey)

	instructionText := i18n.T.PairingInstruction
	if Emoji(EmojiKeyboard) != "" {
		instructionText = Emoji(EmojiKeyboard) + "  " + instructionText
	}
	instruction := WarningStyle.Render(instructionText)
	confirm := HelpStyle.Render(i18n.T.PairingConfirm)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		passkeyText,
		"",
		instruction,
		confirm,
	)

	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())

	if effectiveWidth > 0 {
		return PasskeyBoxStyle.Width(min(effectiveWidth-4, 70)).Render(content)
	}

	return PasskeyBoxStyle.Render(content)
}

// renderDeviceCount renders the device counter.
func renderDeviceCount(count int) string {
	return MutedStyle.Render(fmt.Sprintf("(%d)", count))
}

// renderSectionHeader renders a section header.
func renderSectionHeader(icon, title string, count int, isFocused bool) string {
	countStr := renderDeviceCount(count)

	focusMarker := ""
	if isFocused && Emoji(EmojiFocusMarker) != "" {
		focusMarker = SelectedStyle.Render(" " + Emoji(EmojiFocusMarker))
	}

	header := fmt.Sprintf("%s %s %s%s", icon, title, countStr, focusMarker)

	return HeaderStyle.Render(header)
}

// renderDeviceItem renders a device item.
func renderDeviceItem(dev *models.Device, isSelected bool, showRSSI bool) string {
	icon := DeviceIconStyle.Render(dev.GetIcon())
	name := DeviceNameStyle.Render(dev.GetDisplayName())

	// Address (conditional based on config)
	parts := []string{icon, name}
	showAddress := true
	if config.Global != nil {
		showAddress = config.Global.ShowDeviceAddress
	}
	if showAddress {
		address := DeviceAddressStyle.Render(fmt.Sprintf("(%s)", dev.Address))
		parts = append(parts, address)
	}

	// Additional information
	var info []string

	// RSSI (conditional based on config)
	showRSSIConfig := true
	if config.Global != nil {
		showRSSIConfig = config.Global.ShowRSSI
	}
	if showRSSI && showRSSIConfig && dev.RSSI != 0 {
		info = append(info, DeviceInfoStyle.Render(fmt.Sprintf("%d dBm", dev.RSSI)))
	}

	// Battery (conditional based on config)
	showBatteryConfig := true
	if config.Global != nil {
		showBatteryConfig = config.Global.ShowBattery
	}
	if showBatteryConfig && dev.HasBattery() {
		battIcon, battText := dev.GetBatteryInfo()
		batteryStr := fmt.Sprintf("%s %s", battIcon, battText)
		batteryStyled := GetBatteryStyle(*dev.Battery).Render(batteryStr)
		info = append(info, batteryStyled)
	}

	// Badges
	var badges []string
	if dev.Connected {
		badges = append(badges, ConnectedBadgeStyle.Render(i18n.T.BadgeConnected))
	} else if dev.Paired {
		badges = append(badges, PairedBadgeStyle.Render(i18n.T.BadgePaired))
	}
	if dev.Trusted {
		badges = append(badges, SuccessStyle.Render(i18n.T.BadgeTrusted))
	}

	// Build the line (parts already initialized above)

	if len(info) > 0 {
		parts = append(parts, MutedStyle.Render("|"), strings.Join(info, " "+MutedStyle.Render("|")+" "))
	}

	if len(badges) > 0 {
		parts = append(parts, strings.Join(badges, " "))
	}

	line := strings.Join(parts, " ")

	if isSelected {
		prefix := ""
		if Emoji(EmojiSelector) != "" {
			prefix = Emoji(EmojiSelector) + " "
		}
		return SelectedDeviceItemStyle.Render(prefix + line)
	}

	return DeviceItemStyle.Render(line)
}

// renderEmptyState renders an empty state.
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

// renderSeparator renders a separator.
func (m Model) renderSeparator() string {
	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())

	if effectiveWidth > 0 {
		return SeparatorStyle.Render(strings.Repeat("─", effectiveWidth-4))
	}
	return SeparatorStyle.Render(strings.Repeat("─", 80))
}

// renderThickSeparator renders a thick separator.
func (m Model) renderThickSeparator() string {
	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())

	if effectiveWidth > 0 {
		return ThickSeparatorStyle.Render(strings.Repeat("━", effectiveWidth-4))
	}
	return ThickSeparatorStyle.Render(strings.Repeat("━", 80))
}

// renderAdapterTable renders the adapter information table.
func (m Model) renderAdapterTable() string {
	if m.adapter == nil {
		return BoxStyle.Render(MutedStyle.Render(i18n.T.StatusLoadingAdapterInfo))
	}

	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())
	if effectiveWidth <= 0 {
		effectiveWidth = 80
	}

	// BoxStyle has Padding(0, 1) and borders
	// Border = 2 chars, Padding = 2 chars (1 on each side) = 4 total
	// Additional adjustment for alignment
	availableWidth := effectiveWidth - 6

	// Calculate column width (5 columns, distribute evenly)
	colWidth := availableWidth / 5

	// Header style
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(secondaryColor).
		Width(colWidth).
		Align(lipgloss.Center)

	// Cell style
	cellStyle := lipgloss.NewStyle().
		Width(colWidth).
		Align(lipgloss.Center)

	// Build header row
	headers := []string{
		headerStyle.Render(i18n.T.AdapterName),
		headerStyle.Render(i18n.T.AdapterAlias),
		headerStyle.Render(i18n.T.AdapterPower),
		headerStyle.Render(i18n.T.AdapterPairable),
		headerStyle.Render(i18n.T.AdapterDiscoverable),
	}
	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headers...)

	// Build separator
	separator := SeparatorStyle.Render(strings.Repeat("─", availableWidth))

	// Build data row with color coding
	powerText := ErrorStyle.Render(i18n.T.StatusOff)
	if m.adapter.Powered {
		powerText = SuccessStyle.Render(i18n.T.StatusOn)
	}

	pairableText := MutedStyle.Render(i18n.T.StatusOff)
	if m.adapter.Pairable {
		pairableText = SuccessStyle.Render(i18n.T.StatusOn)
	}

	discoverableText := MutedStyle.Render(i18n.T.StatusOff)
	if m.adapter.Discoverable {
		discoverableText = SuccessStyle.Render(i18n.T.StatusOn)
	}

	cells := []string{
		cellStyle.Render(m.adapter.Name),
		cellStyle.Render(m.adapter.Alias),
		cellStyle.Render(powerText),
		cellStyle.Render(pairableText),
		cellStyle.Render(discoverableText),
	}
	dataRow := lipgloss.JoinHorizontal(lipgloss.Top, cells...)

	// Combine table content
	tableContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerRow,
		separator,
		dataRow,
	)

	// Add title header with icon
	title := HeaderStyle.Render(
		fmt.Sprintf("%s %s", m.adapter.GetStatusIcon(), i18n.T.AdapterInfo),
	)
	titleSep := SeparatorStyle.Render(strings.Repeat("─", availableWidth))

	table := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		titleSep,
		tableContent,
	)

	return BoxStyle.Width(effectiveWidth - 4).Render(table)
}
