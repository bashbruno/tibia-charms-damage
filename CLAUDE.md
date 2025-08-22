# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

Run the web application locally:
```bash
make api
```

Build for production:
```bash
make build
```

Run with Docker:
```bash
make docker
```

Run directly with Go (requires setting PORT and DATA_URL environment variables):
```bash
go run ./cmd/api
```

Default environment variables are defined in the Makefile:
- PORT=8000
- DATA_URL=https://raw.githubusercontent.com/mathiasbynens/tibia-json/main/data/bestiary.json

## Architecture

This is a Go web application that calculates damage comparisons for Tibia game charms. The application has a simple architecture:

### Core Components

- **cmd/api/**: Main application entry point and HTTP handlers
  - `main.go`: Application bootstrap and initialization
  - `api.go`: HTTP server setup and routing with ServeMux
  - `handlers.go`: Request handlers for home page and static assets
  - `html.go`: HTML template rendering utilities
  - `errors.go`: Error handling utilities

- **internal/storage/**: Data layer for creature information
  - `creatures.go`: Core business logic for damage calculations, creature data loading from external JSON API, and breakpoint analysis for charm comparisons

- **internal/env/**: Environment variable utilities
  - `env.go`: Helper functions for reading environment variables with fallbacks

- **web/**: Frontend assets
  - `templates/`: HTML templates (index.html, results.html, card.html)
  - `static/`: CSS and images

### Key Business Logic

The application calculates damage breakpoints for Tibia charms:
- **Overflux charms**: Use 2.5% of current mana per damage point
- **Overpower charms**: Use 5% of current health per damage point  
- **Elemental charms**: Deal 5% of creature's hitpoints as base damage
- **Maximum damage cap**: 8% of creature's hitpoints

The core calculation logic is in `internal/storage/creatures.go` with methods like:
- `GetBreakpoints()`: Calculates resource requirements for different charm types
- `GetElementalCharmDamage()`: Determines neutral and strongest elemental damage
- `GetResistances()`: Finds creature's highest elemental resistance

### Data Flow

1. Application fetches creature data from external JSON API on startup
2. User searches for creatures via web interface
3. Application performs fuzzy search and calculates charm damage breakpoints
4. Results are rendered using HTML templates with damage comparisons

The application is stateless and loads all creature data into memory at startup for fast lookups.