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
	dotChar      string = " • "
	asteriskChar string = "*"
	greenClr            = lipgloss.Color("#00D787")
)

var (
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	headerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	subtleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dotStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	checkboxStyle   = lipgloss.NewStyle().Foreground(greenClr)
	enumeratorStyle = lipgloss.NewStyle().Foreground(greenClr).MarginRight(1)
	itemStyle       = lipgloss.NewStyle().Foreground(greenClr).MarginRight(1)
	asteriskStyle   = lipgloss.NewStyle().Foreground(greenClr)
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
