package ui

import "github.com/charmbracelet/lipgloss"

// Colores del tema
var (
	primaryColor   = lipgloss.Color("205")
	secondaryColor = lipgloss.Color("86")
	accentColor    = lipgloss.Color("57")
	successColor   = lipgloss.Color("82")
	warningColor   = lipgloss.Color("214")
	errorColor     = lipgloss.Color("196")
	mutedColor     = lipgloss.Color("240")
	borderColor    = lipgloss.Color("238")
	highlightColor = lipgloss.Color("230")
	infoColor      = lipgloss.Color("51")
)

// Estilos de texto
var (
	TitleStyle      = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Padding(0, 1)
	SubtitleStyle   = lipgloss.NewStyle().Foreground(secondaryColor).Italic(true)
	HeaderStyle     = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Padding(0, 1)
	DeviceStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	ConnectedStyle  = lipgloss.NewStyle().Foreground(successColor)
	SelectedStyle   = lipgloss.NewStyle().Foreground(highlightColor).Background(accentColor).Bold(true)
	ErrorStyle      = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	HelpStyle       = lipgloss.NewStyle().Foreground(mutedColor)
	StatusStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
	ConnectingStyle = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
	WarningStyle    = lipgloss.NewStyle().Foreground(warningColor)
	SuccessStyle    = lipgloss.NewStyle().Foreground(successColor)
	InfoStyle       = lipgloss.NewStyle().Foreground(infoColor)
	MutedStyle      = lipgloss.NewStyle().Foreground(mutedColor)
)

// Estilos de bordes y contenedores
var (
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2)

	FocusedPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(1, 2)

	PasskeyBoxStyle = lipgloss.NewStyle().
			Foreground(infoColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(infoColor).
			Align(lipgloss.Center)

	HeaderBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(borderColor).
			Padding(0, 1)

	FooterBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(borderColor).
			Padding(0, 1)
)

// Estilos de lista
var (
	DeviceItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedDeviceItemStyle = lipgloss.NewStyle().
				Foreground(highlightColor).
				Background(accentColor).
				Bold(true).
				PaddingLeft(1)

	DeviceIconStyle = lipgloss.NewStyle().
			Width(3).
			Align(lipgloss.Center)

	DeviceNameStyle = lipgloss.NewStyle().
			Bold(true)

	DeviceAddressStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	DeviceInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	BatteryHighStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	BatteryMediumStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Bold(true)

	BatteryLowStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)
)

// Estilos de separadores
var (
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(borderColor)

	ThickSeparatorStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true)
)

// Estilos de badges
var (
	PairedBadgeStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Background(lipgloss.Color("22")).
				Padding(0, 1).
				Bold(true)

	ScanningBadgeStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Background(lipgloss.Color("58")).
				Padding(0, 1).
				Bold(true)

	ConnectedBadgeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(successColor).
				Padding(0, 1).
				Bold(true)
)

// Funciones de ayuda para layouts
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetBatteryStyle retorna el estilo apropiado según el nivel de batería.
func GetBatteryStyle(level uint8) lipgloss.Style {
	switch {
	case level >= 60:
		return BatteryHighStyle
	case level >= 30:
		return BatteryMediumStyle
	default:
		return BatteryLowStyle
	}
}
