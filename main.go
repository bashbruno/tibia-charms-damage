package main

import (
	"fmt"
	"log"

	"github.com/bashbruno/tibia-charms-damage/storage"
)

var creatureStore *storage.CreatureStore

func main() {
	store, err := storage.InitializeCreatureStore()
	if err != nil {
		log.Fatalf("Failed to load creature data: %v", err)
	}
	creatureStore = store

	matches := creatureStore.FuzzyFind("drag")
	fmt.Printf("Found %d creatures matching 'drag':\n", len(matches))
	for _, creature := range matches {
		fmt.Printf("  - %s\n", creature.Name)
	}
}
