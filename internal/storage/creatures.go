package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"
)

const (
	overfluxResourcePercentage  float64 = 2.5
	overpowerResourcePercentage float64 = 5
	maxDamagePercentage         float64 = 8
	dataURL                             = "https://raw.githubusercontent.com/mathiasbynens/tibia-json/main/data/bestiary.json"
	startingLevel                       = 8
	startingHealth                      = 185
	startingMana                        = 90
)

type ResourceType int

const (
	ResourceMana ResourceType = iota
	ResourceHealth
)

type Class struct {
	Name               string
	HealthGainPerLevel int
	ManaGainPerLevel   int
}

var classes = map[string]Class{
	"Knight": {
		Name:               "Knight",
		HealthGainPerLevel: 15,
		ManaGainPerLevel:   5,
	},
	"Mage": {
		Name:               "Mage",
		HealthGainPerLevel: 5,
		ManaGainPerLevel:   30,
	},
	"Paladin": {
		Name:               "Paladin",
		HealthGainPerLevel: 10,
		ManaGainPerLevel:   15,
	},
	"Monk": {
		Name:               "Monk",
		HealthGainPerLevel: 10,
		ManaGainPerLevel:   10,
	},
}

// ClassLevels holds the required level per class for a single resource breakpoint.
type ClassLevels struct {
	Knight  int
	Mage    int
	Paladin int
	Monk    int
}

// ElementBreakpoint is a per-element row: its charm damage and what Overflux/Overpower need to match it.
type ElementBreakpoint struct {
	Element               string
	ResistancePercent     float64
	CharmDamage           float64
	ExceedsCap            bool
	OverfluxManaNeeded    float64
	OverfluxLevels        ClassLevels
	OverpowerHealthNeeded float64
	OverpowerLevels       ClassLevels
}

// DamageCap holds the absolute cap breakpoint (8% of HP).
type DamageCap struct {
	MaxDamage             float64
	OverfluxManaNeeded    float64
	OverfluxLevels        ClassLevels
	OverpowerHealthNeeded float64
	OverpowerLevels       ClassLevels
}

// BreakpointSummary is the top-level result returned by GetBreakpoints().
type BreakpointSummary struct {
	Elements []ElementBreakpoint
	Cap      DamageCap
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
	DeathDmgMod    float64 `json:"deathDmgMod"`
	EnergyDmgMod   float64 `json:"energyDmgMod"`
	HolyDmgMod     float64 `json:"holyDmgMod"`
	IceDmgMod      float64 `json:"iceDmgMod"`
	HealDmgMod     float64 `json:"healMod"`
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
	hp := creature.Hitpoints
	cap := math.Round(getPercentage(hp, maxDamagePercentage))

	elementMods := []struct {
		name string
		mod  float64
	}{
		{"🔥 Fire", creature.FireDmgMod},
		{"💀 Death", creature.DeathDmgMod},
		{"🌿 Earth", creature.EarthDmgMod},
		{"⚡ Energy", creature.EnergyDmgMod},
		{"✨ Holy", creature.HolyDmgMod},
		{"❄️ Ice", creature.IceDmgMod},
		{"⚔️ Physical", creature.PhysicalDmgMod},
	}

	elements := make([]ElementBreakpoint, 0, len(elementMods))
	for _, em := range elementMods {
		charmDmg := math.Round(getPercentage(hp, 5) * em.mod)
		exceedsCap := charmDmg > cap

		eb := ElementBreakpoint{
			Element:           em.name,
			ResistancePercent: math.Round(em.mod * 100),
			CharmDamage:       charmDmg,
			ExceedsCap:        exceedsCap,
		}

		if !exceedsCap {
			eb.OverfluxManaNeeded = getResourceNeeded(charmDmg, overfluxResourcePercentage)
			eb.OverfluxLevels = calculateClassLevels(eb.OverfluxManaNeeded, ResourceMana)
			eb.OverpowerHealthNeeded = getResourceNeeded(charmDmg, overpowerResourcePercentage)
			eb.OverpowerLevels = calculateClassLevels(eb.OverpowerHealthNeeded, ResourceHealth)
		}

		elements = append(elements, eb)
	}

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].CharmDamage > elements[j].CharmDamage
	})

	manaNeededCap := getResourceNeeded(cap, overfluxResourcePercentage)
	healthNeededCap := getResourceNeeded(cap, overpowerResourcePercentage)

	return &BreakpointSummary{
		Elements: elements,
		Cap: DamageCap{
			MaxDamage:             cap,
			OverfluxManaNeeded:    manaNeededCap,
			OverfluxLevels:        calculateClassLevels(manaNeededCap, ResourceMana),
			OverpowerHealthNeeded: healthNeededCap,
			OverpowerLevels:       calculateClassLevels(healthNeededCap, ResourceHealth),
		},
	}
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

func calculateLevelForResource(requiredAmount float64, class Class, resourceType ResourceType) int {
	var startingAmount float64
	var gainPerLevel int

	switch resourceType {
	case ResourceMana:
		startingAmount = startingMana
		gainPerLevel = class.ManaGainPerLevel
	case ResourceHealth:
		startingAmount = startingHealth
		gainPerLevel = class.HealthGainPerLevel
	}

	if requiredAmount <= startingAmount {
		return startingLevel
	}

	additionalResourceNeeded := requiredAmount - startingAmount
	levelsNeeded := math.Ceil(additionalResourceNeeded / float64(gainPerLevel))
	return startingLevel + int(levelsNeeded)
}

func calculateClassLevels(resource float64, rt ResourceType) ClassLevels {
	return ClassLevels{
		Knight:  calculateLevelForResource(resource, classes["Knight"], rt),
		Mage:    calculateLevelForResource(resource, classes["Mage"], rt),
		Paladin: calculateLevelForResource(resource, classes["Paladin"], rt),
		Monk:    calculateLevelForResource(resource, classes["Monk"], rt),
	}
}
