package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strings"

	"github.com/bashbruno/tibia-charms-damage/internal/env"
)

const (
	overfluxResourcePercentage  float64 = 2.5
	overpowerResourcePercentage float64 = 5
)

type BreakpointSummary struct {
	NeutralElementalDamage float64
	WeakestElementalDamage float64
	Overflux               CharmSummary
	Overpower              CharmSummary
}

type CharmSummary struct {
	BreakEvenNeutralResourceNeeded float64
	BreakEvenWeakestResourceNeeded float64
	MaxDamage                      float64
	MaxDamageResourceNeeded        float64
}

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
			matches = append(matches, &creature)
		}
	}

	return matches
}

func (cs *CreatureStore) GetBreakpoints(creature *Creature) *BreakpointSummary {
	neutral, weakest := cs.GetElementalCharmDamage(creature)
	maxDamageAllowed := getPercentage(creature.Hitpoints, 8)

	manaNeededNeutral := getResourceNeeded(neutral, overfluxResourcePercentage)
	manaNeededWeakest := getResourceNeeded(weakest, overfluxResourcePercentage)
	manaNeededMax := getResourceNeeded(maxDamageAllowed, overfluxResourcePercentage)

	healthNeededWeakest := getResourceNeeded(weakest, overpowerResourcePercentage)
	healthNeededNeutral := getResourceNeeded(neutral, overpowerResourcePercentage)
	healthNeededMax := getResourceNeeded(maxDamageAllowed, overpowerResourcePercentage)

	return &BreakpointSummary{
		NeutralElementalDamage: neutral,
		WeakestElementalDamage: weakest,
		Overflux: CharmSummary{
			BreakEvenNeutralResourceNeeded: manaNeededNeutral,
			BreakEvenWeakestResourceNeeded: manaNeededWeakest,
			MaxDamage:                      maxDamageAllowed,
			MaxDamageResourceNeeded:        manaNeededMax,
		},
		Overpower: CharmSummary{
			BreakEvenNeutralResourceNeeded: healthNeededNeutral,
			BreakEvenWeakestResourceNeeded: healthNeededWeakest,
			MaxDamage:                      maxDamageAllowed,
			MaxDamageResourceNeeded:        healthNeededMax,
		},
	}
}

func (cs *CreatureStore) GetElementalCharmDamage(creature *Creature) (float64, float64) {
	var neutral float64 = -1
	var weakest float64 = -1

	_, highest := cs.GetResistances(creature)

	neutral = getPercentage(creature.Hitpoints, 5)
	weakest = neutral * highest

	return neutral, weakest
}

func (cs *CreatureStore) GetResistances(creature *Creature) ([]float64, float64) {
	var highest float64 = -1

	resistances := make([]float64, 7)
	resistances = append(resistances,
		creature.FireDmgMod,
		creature.DeathDmgMod,
		creature.EarthDmgMod,
		creature.EnergyDmgMod,
		creature.HolyDmgMod,
		creature.IceDmgMod,
		creature.PhysicalDmgMod)

	for _, r := range resistances {
		highest = math.Max(highest, r)
	}

	return resistances, highest
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

func getPercentage(from float64, target float64) float64 {
	return from * (target / 100)
}

func getResourceNeeded(target float64, percentage float64) float64 {
	return target / (percentage / 100)
}
