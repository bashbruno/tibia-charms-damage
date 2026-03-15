package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/list"
)

func (m model) View() string {
	s := spinnerView(m)

	if m.state.showInputView {
		s = inputView(m)
	}

	if m.state.showResultView {
		s = resultView(m)
	}

	if m.state.showChoicesView {
		s = choicesView(m)
	}

	return s
}

func inputView(m model) string {
	noResult := m.state.hasQueried && len(m.results) == 0
	title := headerStyle.Render("What creature are you looking for?")
	if noResult {
		title = errorStyle.Render("Couldn't find a creature that matches that name! Try again:")
	}

	tmpl := "%s\n\n%s\n\n"
	tmpl += subtleStyle.Render("ctrl+c or esc: quit")

	return fmt.Sprintf(tmpl, title, m.query.View())
}

func choicesView(m model) string {
	c := m.choicesCursor

	tmpl := headerStyle.Render("Which of these possible results did you mean?")
	tmpl += "\n\n"
	tmpl += "%s\n\n"
	tmpl += subtleStyle.Render("j/k, tab/shift+tab or up/down: select") + ("\n") +
		subtleStyle.Render("enter: choose") + ("\n") +
		subtleStyle.Render("ctrl+c, q or esc: quit")

	var checkboxes []string
	for i, r := range m.results {
		checkboxes = append(checkboxes, checkbox(r.Name, c == i))
	}
	choicesStr := strings.Join(checkboxes, "\n")

	return fmt.Sprintf(tmpl, choicesStr)
}

func resultView(m model) string {
	c := m.target
	t := m.store.GetBreakpoints(c)

	var overfluxItems []any
	for _, el := range t.Elements {
		if el.ExceedsCap {
			overfluxItems = append(overfluxItems, fmt.Sprintf("%s (%.0f%%): %.0f dmg — exceeds cap", el.Element, el.ResistancePercent, el.CharmDamage))
		} else {
			overfluxItems = append(overfluxItems, fmt.Sprintf("%s (%.0f%%): %.0f dmg → %.0f MP (EK %d | MS %d | RP %d | EM %d)",
				el.Element, el.ResistancePercent, el.CharmDamage, el.OverfluxManaNeeded,
				el.OverfluxLevels.Knight, el.OverfluxLevels.Mage, el.OverfluxLevels.Paladin, el.OverfluxLevels.Monk))
		}
	}
	overfluxItems = append(overfluxItems, fmt.Sprintf("%sDamage cap: %.0f dmg → %.0f MP (EK %d | MS %d | RP %d | EM %d)", asteriskStyle,
		t.Cap.MaxDamage, t.Cap.OverfluxManaNeeded,
		t.Cap.OverfluxLevels.Knight, t.Cap.OverfluxLevels.Mage, t.Cap.OverfluxLevels.Paladin, t.Cap.OverfluxLevels.Monk))

	var overpowerItems []any
	for _, el := range t.Elements {
		if el.ExceedsCap {
			overpowerItems = append(overpowerItems, fmt.Sprintf("%s (%.0f%%): %.0f dmg — exceeds cap", el.Element, el.ResistancePercent, el.CharmDamage))
		} else {
			overpowerItems = append(overpowerItems, fmt.Sprintf("%s (%.0f%%): %.0f dmg → %.0f HP (EK %d | MS %d | RP %d | EM %d)",
				el.Element, el.ResistancePercent, el.CharmDamage, el.OverpowerHealthNeeded,
				el.OverpowerLevels.Knight, el.OverpowerLevels.Mage, el.OverpowerLevels.Paladin, el.OverpowerLevels.Monk))
		}
	}
	overpowerItems = append(overpowerItems, fmt.Sprintf("%sDamage cap: %.0f dmg → %.0f HP (EK %d | MS %d | RP %d | EM %d)", asteriskStyle,
		t.Cap.MaxDamage, t.Cap.OverpowerHealthNeeded,
		t.Cap.OverpowerLevels.Knight, t.Cap.OverpowerLevels.Mage, t.Cap.OverpowerLevels.Paladin, t.Cap.OverpowerLevels.Monk))

	overflux := list.New(overfluxItems...)
	overpower := list.New(overpowerItems...)

	l := list.New(
		"Overflux", overflux,
		"Overpower", overpower,
	)
	l.Enumerator(list.Dash)
	l.EnumeratorStyle(enumeratorStyle)
	l.ItemStyle(itemStyle)

	tmpl := headerStyle.Render("%s (%.0f HP)")
	tmpl += "\n\n"
	tmpl += "%s\n\n"
	tmpl += asteriskStyle + " Overpower and Overflux are capped at 8%% of the creature's health"
	tmpl += "\n\n"
	tmpl += subtleStyle.Render("ctrl+c, q or esc: quit") + ("\n") +
		subtleStyle.Render("enter: new search")

	return fmt.Sprintf(
		tmpl,
		c.Name,
		c.Hitpoints,
		l,
	)
}

func spinnerView(m model) string {
	str := fmt.Sprintf("\n\n   %s Loading data...\n\n", m.spinner.View())
	return str
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
