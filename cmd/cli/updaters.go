package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// make sure these always quit the application
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	if m.state.showInputView {
		return updateInput(msg, m)
	}

	if m.state.showChoicesView {
		return updateChoices(msg, m)
	}

	if m.state.showResultView {
		return updateResult(msg, m)
	}

	return updateLoading(msg, m)
}

func updateLoading(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case downloadMsg:
		m.store = msg.store
		m.state.isLoadingData = false
		m.state.showInputView = true
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func updateInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			q := strings.TrimSpace(m.query.Value())
			r := []rune(q)
			if len(r) > 0 {
				m.state.hasQueried = true
				m.results = m.store.FuzzyFind(q)
				len := len(m.results)

				switch len {
				case 0:
					m.query.SetValue("")
					m.state.showChoicesView = false
					m.state.showResultView = false
					m.state.showInputView = true
				case 1:
					m.target = m.results[0]
					m.state.showResultView = true
					m.state.showChoicesView = false
					m.state.showInputView = false
				default:
					m.state.showChoicesView = true
					m.state.showResultView = false
					m.state.showInputView = false
				}

			}
		default:
			m.query, cmd = m.query.Update(msg)
		}
	}

	return m, cmd
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j", "down", "tab":
			m.choicesCursor++
			len := len(m.results) - 1
			if m.choicesCursor > len {
				m.choicesCursor = len
			}
		case "k", "up", "shift+tab":
			m.choicesCursor--
			if m.choicesCursor < 0 {
				m.choicesCursor = 0
			}
		case "enter":
			m.target = m.results[m.choicesCursor]
			m.state.showChoicesView = false
			m.state.showInputView = false
			m.state.showResultView = true
			return m, nil
		}
	}

	return m, nil
}

func updateResult(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			m := reset(m)
			return m, nil
		}
	}

	return m, nil
}
