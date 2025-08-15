package main

import (
	"net/http"
	"text/template"
)

func (app *application) htmlResponse(w http.ResponseWriter, r *http.Request, filePath string) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		app.logError(r, err)
		http.Error(w, "Unable to load HTML template", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title": "Toma",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		app.logError(r, err)
		http.Error(w, "Unable to render HTML template", http.StatusInternalServerError)
	}
}
