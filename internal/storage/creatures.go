package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bashbruno/tibia-charms-damage/internal/env"
)

type Creature struct {
	BestiaryClass  string  `json:"bestiaryClass"`
	Name           string  `json:"name"`
	BestiaryLevel  string  `json:"bestiaryLevel"`
	Occurrence     string  `json:"occurrence"`
	CharmPoints    float64 `json:"charmPoints"`
	Experience     float64 `json:"experience"`
	Hitpoints      float64 `json:"hitpoints"`
	Armor          float64 `json:"armor"`
	Mitigation     float64 `json:"mitigation"`
	PhysicalDmgMod float64 `json:"physicalDmgMod"`
	EarthDmgMod    float64 `json:"earthDmgMod"`
	FireDmgMod     float64 `json:"fireDmgMod"`
	DeathDmgMod    float64 `json:"deathDmgDod"`
	EnergyDmgMod   float64 `json:"energyDmgMod"`
	HolyDmgMod     float64 `json:"holyDmgMod"`
	IceDmgMod      float64 `json:"iceDmgMod"`
	HealDmgMod     float64 `json:"healDmgMod"`
}

type CreatureStore struct {
	byName map[string]*Creature
	all    []Creature
}

func LoadCreatures() (*CreatureStore, error) {
	dataURL := env.GetString("DATA_URL", "")
	if dataURL == "" {
		return nil, fmt.Errorf("invalid DATA_URL")
	}

	resp, err := http.Get(dataURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data file: %w", err)
	}

	var creatures []Creature
	if err := json.Unmarshal(data, &creatures); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	store := &CreatureStore{
		byName: make(map[string]*Creature, len(creatures)),
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
	lowerName := strings.ToLower(name)
	creature, exists := cs.byName[lowerName]
	return creature, exists
}

func (cs *CreatureStore) FuzzyFind(searchTerm string) []*Creature {
	var matches []*Creature
	lowerSearch := strings.ToLower(searchTerm)

	for _, creature := range cs.GetAll() {
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
	return len(cs.GetAll())
}

func MakeCreatureStore() (*CreatureStore, error) {
	store, err := LoadCreatures()
	if err != nil {
		return nil, err
	}

	slog.Info("Successfully loaded creatures into memory", "count", store.Count())
	return store, nil
}
