package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/bashbruno/tibia-charms-damage/internal/env"
	"github.com/bashbruno/tibia-charms-damage/internal/storage"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
	store, err := storage.MakeCreatureStore()
	if err != nil {
		log.Fatalf("Failed to load creature data: %v", err)
	}

	return &application{
		config: config{
			addr: env.GetString("ADDR", fallbackListenAddr),
		},
		store: store,
	}
}
