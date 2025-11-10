package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// ThemeMode defines how colors are sourced
type ThemeMode string

const (
	ThemeModeANSI      ThemeMode = "ansi"      // Use terminal ANSI colors
	ThemeModeTrueColor ThemeMode = "truecolor" // Use custom hex colors
)

// ColorScheme holds the color definitions for the UI
type ColorScheme struct {
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Muted      lipgloss.Color
	Border     lipgloss.Color
	Highlight  lipgloss.Color
	Info       lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
}

// CurrentColorScheme is the active color scheme
var CurrentColorScheme *ColorScheme

// InitializeTheme sets up the color scheme based on the theme mode
func InitializeTheme(mode ThemeMode) error {
	switch mode {
	case ThemeModeANSI:
		CurrentColorScheme = getANSIColorScheme()
	case ThemeModeTrueColor:
		CurrentColorScheme = getDefaultTrueColorScheme()
	default:
		CurrentColorScheme = getANSIColorScheme()
	}

	// Update all style variables
	updateStyles()
	return nil
}

// getANSIColorScheme returns a color scheme using terminal ANSI colors
// This respects the terminal's configured color scheme
func getANSIColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary:    lipgloss.Color("5"),  // Magenta
		Secondary:  lipgloss.Color("6"),  // Cyan
		Accent:     lipgloss.Color("4"),  // Blue
		Success:    lipgloss.Color("2"),  // Green
		Warning:    lipgloss.Color("3"),  // Yellow
		Error:      lipgloss.Color("1"),  // Red
		Muted:      lipgloss.Color("8"),  // Bright Black (Gray)
		Border:     lipgloss.Color("8"),  // Bright Black (Gray)
		Highlight:  lipgloss.Color("15"), // Bright White
		Info:       lipgloss.Color("6"),  // Cyan
		Background: lipgloss.Color("0"),  // Black
		Foreground: lipgloss.Color("7"),  // White
	}
}

// getDefaultTrueColorScheme returns the original hardcoded color scheme
func getDefaultTrueColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary:    lipgloss.Color("205"),
		Secondary:  lipgloss.Color("86"),
		Accent:     lipgloss.Color("57"),
		Success:    lipgloss.Color("82"),
		Warning:    lipgloss.Color("214"),
		Error:      lipgloss.Color("196"),
		Muted:      lipgloss.Color("240"),
		Border:     lipgloss.Color("238"),
		Highlight:  lipgloss.Color("230"),
		Info:       lipgloss.Color("51"),
		Background: lipgloss.Color("0"),
		Foreground: lipgloss.Color("252"),
	}
}

// updateStyles updates all style variables with the current color scheme
func updateStyles() {
	if CurrentColorScheme == nil {
		CurrentColorScheme = getANSIColorScheme()
	}

	cs := CurrentColorScheme

	// Update color variables (for backwards compatibility)
	primaryColor = cs.Primary
	secondaryColor = cs.Secondary
	accentColor = cs.Accent
	successColor = cs.Success
	warningColor = cs.Warning
	errorColor = cs.Error
	mutedColor = cs.Muted
	borderColor = cs.Border
	highlightColor = cs.Highlight
	infoColor = cs.Info

	// Update text styles
	TitleStyle = lipgloss.NewStyle().Foreground(cs.Primary).Bold(true).Padding(0, 1)
	SubtitleStyle = lipgloss.NewStyle().Foreground(cs.Secondary).Italic(true)
	HeaderStyle = lipgloss.NewStyle().Foreground(cs.Secondary).Bold(true).Padding(0, 1)
	DeviceStyle = lipgloss.NewStyle().Foreground(cs.Foreground)
	ConnectedStyle = lipgloss.NewStyle().Foreground(cs.Success)
	SelectedStyle = lipgloss.NewStyle().Foreground(cs.Highlight).Background(cs.Accent).Bold(true)
	ErrorStyle = lipgloss.NewStyle().Foreground(cs.Error).Bold(true)
	HelpStyle = lipgloss.NewStyle().Foreground(cs.Muted)
	StatusStyle = lipgloss.NewStyle().Foreground(cs.Muted).Italic(true)
	ConnectingStyle = lipgloss.NewStyle().Foreground(cs.Warning).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(cs.Warning)
	SuccessStyle = lipgloss.NewStyle().Foreground(cs.Success)
	InfoStyle = lipgloss.NewStyle().Foreground(cs.Info)
	MutedStyle = lipgloss.NewStyle().Foreground(cs.Muted)

	// Update border and container styles
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cs.Border).
		Padding(0, 1)

	PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cs.Secondary).
		Padding(1, 2)

	FocusedPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cs.Accent).
		Padding(1, 2)

	PasskeyBoxStyle = lipgloss.NewStyle().
		Foreground(cs.Info).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(cs.Info).
		Align(lipgloss.Center)

	HeaderBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(cs.Border).
		Padding(0, 1)

	FooterBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(cs.Border).
		Padding(0, 1)

	// Update list styles
	DeviceItemStyle = lipgloss.NewStyle().PaddingLeft(2)

	SelectedDeviceItemStyle = lipgloss.NewStyle().
		Foreground(cs.Highlight).
		Background(cs.Accent).
		Bold(true).
		PaddingLeft(1)

	DeviceIconStyle = lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Center)

	DeviceNameStyle = lipgloss.NewStyle().Bold(true)

	DeviceAddressStyle = lipgloss.NewStyle().Foreground(cs.Muted)

	DeviceInfoStyle = lipgloss.NewStyle().Foreground(cs.Muted)

	BatteryHighStyle = lipgloss.NewStyle().
		Foreground(cs.Success).
		Bold(true)

	BatteryMediumStyle = lipgloss.NewStyle().
		Foreground(cs.Warning).
		Bold(true)

	BatteryLowStyle = lipgloss.NewStyle().
		Foreground(cs.Error).
		Bold(true)

	// Update separator styles
	SeparatorStyle = lipgloss.NewStyle().Foreground(cs.Border)

	ThickSeparatorStyle = lipgloss.NewStyle().
		Foreground(cs.Secondary).
		Bold(true)

	// Update badge styles
	PairedBadgeStyle = lipgloss.NewStyle().
		Foreground(cs.Success).
		Background(cs.Background).
		Padding(0, 1).
		Bold(true)

	ScanningBadgeStyle = lipgloss.NewStyle().
		Foreground(cs.Warning).
		Background(cs.Background).
		Padding(0, 1).
		Bold(true)

	ConnectedBadgeStyle = lipgloss.NewStyle().
		Foreground(cs.Background).
		Background(cs.Success).
		Padding(0, 1).
		Bold(true)
}
