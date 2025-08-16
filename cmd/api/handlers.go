package main

import (
	"net/http"

	"github.com/bashbruno/tibia-charms-damage/internal/storage"
)

const (
	layoutHTML  = "./web/templates/index.html"
	resultsHTML = "./web/templates/results.html"
)

type searchResult struct {
	Creature *storage.Creature
	Summary  *storage.BreakpointSummary
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var result []searchResult

	if query != "" {
		matches := app.store.FuzzyFind(query)
		for _, m := range matches {
			summary := app.store.GetBreakpoints(m)
			result = append(result, searchResult{
				Creature: m,
				Summary:  summary,
			})
		}
	}

	data := map[string]any{
		"Query":  query,
		"Result": result,
	}

	app.htmlResponse(w, r, data, layoutHTML, resultsHTML)
}

func (app *application) assetsHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("web/static"))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}
