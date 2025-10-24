package ui

import "github.com/charmbracelet/lipgloss"

// Estilos de la aplicación.
var (
	TitleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Padding(0, 1)
	HeaderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	DeviceStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Padding(0, 2)
	ConnectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Padding(0, 2)
	SelectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("57")).Bold(true).Padding(0, 2)
	ErrorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	HelpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	SeparatorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	StatusStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
	ConnectingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	WarningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	PasskeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true).Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("51"))
)
