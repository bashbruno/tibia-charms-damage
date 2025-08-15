package main

import (
	"net/http"
)

const indexHTML = "./web/templates/index.html"

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	app.htmlResponse(w, r, indexHTML)
}

func (app *application) staticAssetsHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("web/static"))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}
