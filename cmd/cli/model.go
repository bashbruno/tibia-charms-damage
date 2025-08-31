package main

import (
	"github.com/bashbruno/tibia-charms-damage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	query         textinput.Model
	choicesCursor int
	store         *storage.CreatureStore
	results       []*storage.Creature
	target        *storage.Creature
	state         state
}

type state struct {
	showInputView   bool
	showChoicesView bool
	showResultView  bool
	hasQueried      bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func reset(m model) tea.Model {
	m.choicesCursor = 0
	m.target = nil
	m.results = nil
	m.query.SetValue("")
	m.state.showInputView = true
	m.state.showChoicesView = false
	m.state.showResultView = false
	m.state.hasQueried = false
	return m
}

func makeInitialModel(store *storage.CreatureStore) model {
	ti := textinput.New()
	ti.Placeholder = "dragon"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.PromptStyle = lipgloss.NewStyle().Foreground(greenClr)

	return model{
		query:         ti,
		choicesCursor: 0,
		results:       nil,
		target:        nil,
		store:         store,
		state: state{
			hasQueried:      false,
			showInputView:   true,
			showChoicesView: false,
			showResultView:  false,
		},
	}
}
