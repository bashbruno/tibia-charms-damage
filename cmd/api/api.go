package main

import (
	"errors"
	"log/slog"
	"net/http"
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
	mux.HandleFunc("GET /", app.homeHandler)
	mux.HandleFunc("GET /static/", app.staticAssetsHandler)
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
