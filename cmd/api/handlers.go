package main

import (
	"net/http"

	"github.com/bashbruno/tibia-charms-damage/internal/storage"
)

const (
	layoutHTML  = "./web/templates/index.html"
	resultsHTML = "./web/templates/results.html"
	cardHTML    = "./web/templates/card.html"
)

type searchResult struct {
	Creature *storage.Creature
	Summary  *storage.BreakpointSummary
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	selectedName := r.URL.Query().Get("selected")

	var result []searchResult
	var selectedCreature *searchResult

	if query != "" {
		matches := app.store.FuzzyFind(query)
		for _, m := range matches {
			summary := app.store.GetBreakpoints(m)
			searchRes := searchResult{
				Creature: m,
				Summary:  summary,
			}
			result = append(result, searchRes)

			if selectedName != "" && m.Name == selectedName {
				selectedCreature = &searchRes
			}
		}
	}

	data := map[string]any{
		"Query":            query,
		"Result":           result,
		"SelectedCreature": selectedCreature,
	}

	app.htmlResponse(w, r, data, layoutHTML, resultsHTML, cardHTML)
}

func (app *application) assetsHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("web/static"))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}
