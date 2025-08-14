package storage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const DATA = "./assets/data.json"

type Resistance struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Creature struct {
	Name        string       `json:"name"`
	Difficulty  string       `json:"difficulty"`
	Health      float64      `json:"health"`
	Resistances []Resistance `json:"resistances"`
}

type CreatureStore struct {
	byName map[string]*Creature
	all    []Creature
}

func LoadCreatures(filename string) (*CreatureStore, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", filename, err)
	}

	var creatures []Creature
	if err := json.Unmarshal(data, &creatures); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	store := &CreatureStore{
		byName: make(map[string]*Creature),
		all:    creatures,
	}

	for i := range creatures {
		creature := &creatures[i]
		lowerName := strings.ToLower(creature.Name)
		store.byName[lowerName] = creature
	}

	return store, nil
}

func (cs *CreatureStore) GetByName(name string) (*Creature, bool) {
	creature, exists := cs.byName[name]
	return creature, exists
}

func (cs *CreatureStore) FuzzyFind(searchTerm string) []*Creature {
	var matches []*Creature
	lowerSearch := strings.ToLower(searchTerm)

	for _, creature := range cs.all {
		if strings.Contains(strings.ToLower(creature.Name), lowerSearch) {
			creatureCopy := creature
			matches = append(matches, &creatureCopy)
		}
	}

	return matches
}

func (cs *CreatureStore) GetAll() []Creature {
	return cs.all
}

func (cs *CreatureStore) Count() int {
	return len(cs.all)
}

func MakeCreatureStore() (*CreatureStore, error) {
	store, err := LoadCreatures(DATA)
	if err != nil {
		return nil, err
	}

	slog.Info("Successfully loaded creatures into memory", "count", store.Count())
	return store, nil
}
