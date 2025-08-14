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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	setSlogHandler()

	store, err := storage.MakeCreatureStore()
	if err != nil {
		log.Fatalf("Failed to load creature data: %v", err)
	}

	app := application{
		config: config{
			addr: env.GetString("ADDR", fallbackListenAddr),
		},
		store: store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

func setSlogHandler() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))
}
