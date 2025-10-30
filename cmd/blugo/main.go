package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivangsm/blugo/internal/config"
	"github.com/ivangsm/blugo/internal/i18n"
	"github.com/ivangsm/blugo/internal/ui"
)

// version is set at build time via -ldflags
var version = "dev"

func main() {
	// Parse command line flags
	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(showVersion, "v", false, "Show version information (shorthand)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Blugo %s\n", version)
		fmt.Println("Bluetooth Manager for Linux")
		os.Exit(0)
	}
	// Initialize configuration
	if err := config.Init(); err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Set language from config
	i18n.InitFromConfig(config.Global.Language)

	// Crear el modelo inicial
	m := ui.NewModel()

	// Crear el programa de Bubble Tea con opciones para limpiar la terminal
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Usar pantalla alternativa
		tea.WithMouseCellMotion(), // Habilitar mouse (opcional)
	)

	// Inicializar en segundo plano
	go func() {
		initCmd := ui.InitializeCmd(p)
		msg := initCmd()
		p.Send(msg)
	}()

	// Ejecutar el programa
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al ejecutar el programa: %v\n", err)
		os.Exit(1)
	}
}
