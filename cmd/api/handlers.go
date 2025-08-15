package main

import (
	"net/http"
)

const layoutHTML = "./web/templates/index.html"

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	creatures := app.store.GetAll()
	data := map[string]any{
		"Query":     query,
		"Creatures": creatures,
	}

	app.htmlResponse(w, r, data, layoutHTML)
}

func (app *application) assetsHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("web/static"))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}
