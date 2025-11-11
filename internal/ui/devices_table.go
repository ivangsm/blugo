package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/i18n"
)

// initDevicesTable initializes the devices table with current devices
func (m *Model) initDevicesTable() {
	devices := m.GetFoundDevices()

	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())
	if effectiveWidth <= 0 {
		effectiveWidth = 80
	}

	// Get config values to know which columns to show
	showAddress := true
	showRSSI := true
	showBattery := true
	if config.Global != nil {
		showAddress = config.Global.ShowDeviceAddress
		showRSSI = config.Global.ShowRSSI
		showBattery = config.Global.ShowBattery
	}

	// The panel will use (effectiveWidth - 4) for its content width
	// We need to account for table's internal spacing (columns have separators)
	// Each column adds ~3 chars for spacing/borders
	panelContentWidth := effectiveWidth - 4

	// Count visible columns
	visibleColumns := 2 // icon + name (always visible)
	if showAddress {
		visibleColumns++
	}
	if showRSSI {
		visibleColumns++
	}
	if showBattery {
		visibleColumns++
	}
	visibleColumns++ // status (always visible)

	// Account for column separators and padding
	tableOverhead := visibleColumns * 3
	availableWidth := panelContentWidth - tableOverhead

	if availableWidth < 30 {
		availableWidth = 30 // Minimum width
	}

	// Define minimum column widths
	iconWidth := 2
	minNameWidth := 15
	minStatusWidth := 12

	// Optional columns
	addressWidth := 0
	rssiWidth := 0
	batteryWidth := 0

	if showAddress {
		addressWidth = 17
	}
	if showRSSI {
		rssiWidth = 7
	}
	if showBattery {
		batteryWidth = 9
	}

	// Calculate how much space is left after fixed columns
	fixedWidth := iconWidth + addressWidth + rssiWidth + batteryWidth
	remainingWidth := availableWidth - fixedWidth

	// Distribute remaining width between name and status proportionally
	// Name gets 60%, Status gets 40%
	nameWidth := (remainingWidth * 60) / 100
	statusWidth := remainingWidth - nameWidth

	// Ensure minimums
	if nameWidth < minNameWidth {
		nameWidth = minNameWidth
	}
	if statusWidth < minStatusWidth {
		statusWidth = minStatusWidth
	}

	// Build columns list (only include columns with width > 0)
	columns := []table.Column{
		{Title: "", Width: iconWidth}, // Empty title for icon column
		{Title: i18n.T.DeviceName, Width: nameWidth},
	}

	if addressWidth > 0 {
		columns = append(columns, table.Column{Title: i18n.T.DeviceAddress, Width: addressWidth})
	}
	if rssiWidth > 0 {
		columns = append(columns, table.Column{Title: i18n.T.DeviceRSSI, Width: rssiWidth})
	}
	if batteryWidth > 0 {
		columns = append(columns, table.Column{Title: i18n.T.DeviceBattery, Width: batteryWidth})
	}
	columns = append(columns, table.Column{Title: i18n.T.DeviceStatus, Width: statusWidth})

	// Build rows
	rows := []table.Row{}
	for _, dev := range devices {
		// Icon
		icon := dev.GetIcon()

		// Name
		name := dev.GetDisplayName()

		// Status (badges)
		status := ""
		if dev.Connected {
			status = i18n.T.BadgeConnected
		} else if dev.Paired {
			status = i18n.T.BadgePaired
		}
		if dev.Trusted {
			if status != "" {
				status += " "
			}
			status += i18n.T.BadgeTrusted
		}

		// Build row dynamically based on which columns are shown
		row := table.Row{icon, name}

		if addressWidth > 0 {
			row = append(row, dev.Address)
		}

		if rssiWidth > 0 {
			rssi := ""
			if dev.RSSI != 0 && !dev.Connected {
				rssi = fmt.Sprintf("%d dBm", dev.RSSI)
			}
			row = append(row, rssi)
		}

		if batteryWidth > 0 {
			battery := ""
			if dev.HasBattery() {
				battIcon, battText := dev.GetBatteryInfo()
				battery = fmt.Sprintf("%s %s", battIcon, battText)
			}
			row = append(row, battery)
		}

		row = append(row, status)

		rows = append(rows, row)
	}

	// Preserve cursor position before recreating table
	currentCursor := 0
	if m.devicesTable.Cursor() >= 0 {
		currentCursor = m.devicesTable.Cursor()
	}

	// Create table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(min(len(devices)+1, 15)), // +1 for header, max 15 visible rows
	)

	// Apply custom styles using theme colors
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		BorderBottom(true).
		Bold(true).
		Foreground(secondaryColor)

	s.Selected = s.Selected.
		Foreground(highlightColor).
		Background(accentColor).
		Bold(false)

	t.SetStyles(s)

	// Restore cursor position (ensure it's within bounds)
	if currentCursor < len(rows) {
		t.SetCursor(currentCursor)
	} else if len(rows) > 0 {
		t.SetCursor(len(rows) - 1)
	}

	m.devicesTable = t
}

// renderDevicesTable renders the devices table
func (m Model) renderDevicesTable() string {
	devices := m.GetFoundDevices()

	// Use effective width
	effectiveWidth := min(m.width, GetMaxWidth())

	// FocusedPanelStyle has Padding(1, 2) and borders
	// Border = 2 chars, Padding = 4 chars (2 on each side) = 6 total
	// Additional adjustment for table alignment
	panelWidth := effectiveWidth - 10

	// The separator should match the actual content width inside the panel
	separatorWidth := panelWidth

	// Header
	header := renderSectionHeader(
		Emoji(EmojiAvailable),
		i18n.T.AvailableDevices,
		len(devices),
		true,
	)

	// Separator that fits the content width
	separator := SeparatorStyle.Render(strings.Repeat("─", separatorWidth))

	sections := []string{header}

	if len(devices) == 0 {
		sections = append(sections, separator)
		sections = append(sections, renderEmptyState(getEmptyAvailableDevicesMessage()))
	} else {
		sections = append(sections, separator)
		sections = append(sections, m.devicesTable.View())
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return FocusedPanelStyle.Width(effectiveWidth - 4).Render(content)
}
