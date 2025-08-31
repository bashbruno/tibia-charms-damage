package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/bashbruno/tibia-charms-damage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type application struct {
	store *storage.CreatureStore
}

const (
	dotChar        string = " • "
	asteriskChar   string = "*"
	greenClr              = lipgloss.Color("#00D787")
	redClr                = lipgloss.Color("9")
	purpleClr             = lipgloss.Color("99")
	darkGrayClr           = lipgloss.Color("241")
	ligtherGrayClr        = lipgloss.Color("236")
)

var (
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(redClr))
	headerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(purpleClr))
	subtleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(darkGrayClr))
	dotStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(ligtherGrayClr)).Render(dotChar)
	checkboxStyle   = lipgloss.NewStyle().Foreground(greenClr)
	enumeratorStyle = lipgloss.NewStyle().Foreground(greenClr).MarginRight(1)
	itemStyle       = lipgloss.NewStyle().Foreground(greenClr).MarginRight(1)
	asteriskStyle   = lipgloss.NewStyle().Foreground(greenClr).Render(asteriskChar)
)

func main() {
	app := makeApp()
	p := tea.NewProgram(makeInitialModel(app.store))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func makeApp() *application {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))

	store, err := storage.MakeCreatureStore()
	if err != nil {
		log.Fatalf("Failed to load creature data: %v", err)
	}

	return &application{
		store: store,
	}
}
