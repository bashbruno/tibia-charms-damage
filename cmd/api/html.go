package main

import (
	"html/template"
	"net/http"
)

func (app *application) htmlResponse(w http.ResponseWriter, r *http.Request, data any, filenames ...string) {
	tmpl, err := template.ParseFiles(filenames...)
	if err != nil {
		app.logError(r, err)
		http.Error(w, "Unable to load HTML template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		app.logError(r, err)
		http.Error(w, "Unable to render HTML template", http.StatusInternalServerError)
	}
}
