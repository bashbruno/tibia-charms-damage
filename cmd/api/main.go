package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/bashbruno/tibia-charms-damage/internal/env"
	"github.com/bashbruno/tibia-charms-damage/internal/storage"
)

const fallbackListenAddr = ":8000"

func init() {
}

func main() {
	app := makeApp()
	mux := app.mount()
	log.Fatal(app.run(mux))
}

func makeApp() *application {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))

	store, err := storage.MakeCreatureStore()
	if err != nil {
		log.Fatalf("Failed to load creature data: %v", err)
	}

	return &application{
		config: config{
			addr: env.GetString("PORT", fallbackListenAddr),
		},
		store: store,
	}
}
