package main

import (
	"log"
	"time"

	"github.com/bashbruno/tibia-charms-damage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	asteriskChar string = "*"
	greenClr            = lipgloss.Color("#00D787")
	redClr              = lipgloss.Color("9")
	purpleClr           = lipgloss.Color("99")
	darkGrayClr         = lipgloss.Color("241")
)

var (
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(redClr))
	headerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(purpleClr))
	subtleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(darkGrayClr))
	checkboxStyle   = lipgloss.NewStyle().Foreground(greenClr)
	enumeratorStyle = lipgloss.NewStyle().Foreground(greenClr).MarginRight(1)
	itemStyle       = lipgloss.NewStyle().Foreground(greenClr).MarginRight(1)
	asteriskStyle   = lipgloss.NewStyle().Foreground(greenClr).Render(asteriskChar)
)

type downloadMsg struct {
	store *storage.CreatureStore
}

func main() {
	p := tea.NewProgram(makeInitialModel(), tea.WithAltScreen())

	go func() {
		time.Sleep(300 * time.Millisecond)
		store, err := storage.MakeCreatureStore()
		if err != nil {
			log.Fatal(err)
		}
		p.Send(downloadMsg{
			store: store,
		})
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
