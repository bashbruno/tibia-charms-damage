package main

import (
	"net/http"

	"github.com/bashbruno/tibia-charms-damage/internal/storage"
)

const (
	layoutHTML  = "./web/templates/index.html"
	resultsHTML = "./web/templates/results.html"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var matches []*storage.Creature
	if query != "" {
		matches = app.store.FuzzyFind(query)
	}

	data := map[string]any{
		"Query":   query,
		"Matches": matches,
	}

	app.htmlResponse(w, r, data, layoutHTML, resultsHTML)
}

func (app *application) assetsHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("web/static"))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}
