package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

const (
	overfluxResourcePercentage  float64 = 2.5
	overpowerResourcePercentage float64 = 5
	maxDamagePercentage         float64 = 8
	dataURL                             = "https://raw.githubusercontent.com/mathiasbynens/tibia-json/main/data/bestiary.json"
)

type BreakpointSummary struct {
	NeutralElementalDamage       float64
	StrongestElementalDamage     float64
	StrongestElementalPercentage float64
	Overflux                     CharmSummary
	Overpower                    CharmSummary
}

type CharmSummary struct {
	BreakEvenNeutralResourceNeeded   float64
	BreakEvenStrongestResourceNeeded float64
	MaxDamage                        float64
	MaxDamageResourceNeeded          float64
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
	neutral, strongest := cs.GetElementalCharmDamage(creature)
	_, highest := cs.GetResistances(creature)
	maxDamageAllowed := math.Round(getPercentage(creature.Hitpoints, maxDamagePercentage))

	manaNeededNeutral := getResourceNeeded(neutral, overfluxResourcePercentage)
	manaNeededStrongest := getResourceNeeded(strongest, overfluxResourcePercentage)
	manaNeededMax := getResourceNeeded(maxDamageAllowed, overfluxResourcePercentage)

	healthNeededStrongest := getResourceNeeded(strongest, overpowerResourcePercentage)
	healthNeededNeutral := getResourceNeeded(neutral, overpowerResourcePercentage)
	healthNeededMax := getResourceNeeded(maxDamageAllowed, overpowerResourcePercentage)

	return &BreakpointSummary{
		NeutralElementalDamage:       neutral,
		StrongestElementalDamage:     strongest,
		StrongestElementalPercentage: math.Round(highest * 100),
		Overflux: CharmSummary{
			BreakEvenNeutralResourceNeeded:   manaNeededNeutral,
			BreakEvenStrongestResourceNeeded: manaNeededStrongest,
			MaxDamage:                        maxDamageAllowed,
			MaxDamageResourceNeeded:          manaNeededMax,
		},
		Overpower: CharmSummary{
			BreakEvenNeutralResourceNeeded:   healthNeededNeutral,
			BreakEvenStrongestResourceNeeded: healthNeededStrongest,
			MaxDamage:                        maxDamageAllowed,
			MaxDamageResourceNeeded:          healthNeededMax,
		},
	}
}

func (cs *CreatureStore) GetElementalCharmDamage(creature *Creature) (float64, float64) {
	var neutral float64 = -1
	var strongest float64 = -1

	_, highest := cs.GetResistances(creature)

	neutral = math.Round(getPercentage(creature.Hitpoints, 5))
	strongest = math.Round(neutral * highest)

	return neutral, strongest
}

func (cs *CreatureStore) GetResistances(creature *Creature) ([]float64, float64) {
	var highest float64 = -1

	resistances := make([]float64, 0)
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

	return store, nil
}

func getPercentage(from float64, target float64) float64 {
	return from * (target / 100)
}

func getResourceNeeded(target float64, percentage float64) float64 {
	return math.Round(target / (percentage / 100))
}
