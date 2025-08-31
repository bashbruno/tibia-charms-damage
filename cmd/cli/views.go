package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/list"
)

func (m model) View() string {
	s := inputView(m)

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
	tmpl += subtleStyle.Render("ctrl+c or esc to quit")

	return fmt.Sprintf(tmpl, title, m.query.View())
}

func choicesView(m model) string {
	c := m.choicesCursor

	tmpl := headerStyle.Render("Which of these possible results did you mean?")
	tmpl += "\n\n"
	tmpl += "%s\n\n"
	tmpl += subtleStyle.Render("j/k, tab/shift+tab or up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("ctrl+c, q or esc to quit")

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

	general := list.New(
		fmt.Sprintf("Neutral Elemental Damage (100%%): %.0f", t.NeutralElementalDamage),
		fmt.Sprintf("Strongest%s Elemental Damage (%.0f%%): %.0f\n", asteriskStyle, t.StrongestElementalPercentage, t.StrongestElementalDamage),
	)

	overflux := list.New(
		fmt.Sprintf("Mana needed to outdamage Neutral: %.0f", t.Overflux.BreakEvenNeutralResourceNeeded),
		fmt.Sprintf("Mana needed to outdamage Strongest: %.0f", t.Overflux.BreakEvenStrongestResourceNeeded),
		fmt.Sprintf("Mana needed to hit damage cap%s (%.0f): %.0f\n", asteriskStyle, t.Overflux.MaxDamage, t.Overflux.MaxDamageResourceNeeded),
	)

	overpower := list.New(
		fmt.Sprintf("Health needed to outdamage Neutral: %.0f", t.Overpower.BreakEvenNeutralResourceNeeded),
		fmt.Sprintf("Health needed to outdamage Strongest: %.0f", t.Overpower.BreakEvenStrongestResourceNeeded),
		fmt.Sprintf("Health needed to hit damage cap%s (%.0f): %.0f", asteriskStyle, t.Overpower.MaxDamage, t.Overpower.MaxDamageResourceNeeded),
	)

	l := list.New(
		"General", general,
		"Overflux", overflux,
		"Overpower", overpower,
	)
	l.Enumerator(list.Dash)
	l.EnumeratorStyle(enumeratorStyle)
	l.ItemStyle(itemStyle)

	tmpl := headerStyle.Render("%s (%.0f HP)")
	tmpl += "\n\n"
	tmpl += "%s\n\n"
	tmpl += fmt.Sprintf("%s Strongest refers to the highest elemental vulnerability of the creature", asteriskStyle)
	tmpl += "\n"
	tmpl += fmt.Sprintf("%s Overpower and Overflux are capped at 8%% of the creature's health", asteriskStyle)
	tmpl += "\n\n"
	tmpl += subtleStyle.Render("ctrl+c, q or esc to quit") + dotStyle +
		subtleStyle.Render("enter: search again")

	return fmt.Sprintf(
		tmpl,
		c.Name,
		c.Hitpoints,
		l,
	)
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
