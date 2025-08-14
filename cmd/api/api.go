package main

import (
	"errors"
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/bashbruno/tibia-charms-damage/internal/storage"
)

type application struct {
	store  *storage.CreatureStore
	config config
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		app.htmlResponse(w, r, "./web/templates/index.html")
	})

	return mux
}

func (app *application) run(mux http.Handler) error {
	slog.Info("Starting http server", "port", app.config.addr)

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	slog.Info("Server has stopped", "addr", app.config.addr)

	return nil
}

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

func (app *application) logError(r *http.Request, err error) {
	slog.Error("Err", "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "error", err.Error())
}
