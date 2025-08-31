# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
- **Web API**: `make api` (runs on port 8000 by default)
- **CLI Interface**: `make cli` 
- **Docker**: `docker compose up -d`
- **Direct Go**: `go run ./cmd/api` (requires setting PORT environment variable)

### Building and Testing
- **Build**: `make build` (creates optimized Linux binary in `bin/api`)
- **Test**: `make test` (runs all tests with verbose output)
- **Docker Build**: `make docker` (builds and runs container)

## Architecture Overview

This is a Go application for calculating Tibia charm damage comparisons with dual interfaces:

### Core Components
- **cmd/api/**: HTTP web server with HTML templates and static assets
  - Uses standard Go `net/http` with custom routing
  - Serves static files from `web/static/`
  - HTML templates in `web/templates/`
  
- **cmd/cli/**: Terminal user interface built with Charm's Bubble Tea framework
  - Interactive creature search and damage calculations
  - Uses lipgloss for styling

- **internal/storage/**: Core business logic and data management
  - `CreatureStore`: In-memory creature database loaded from external JSON API
  - Implements charm damage calculations (overflux vs overpower vs elemental)
  - Fetches data from `tibia-json` GitHub repository

- **internal/env/**: Environment variable utilities with fallback support

### Data Flow
1. Application startup loads creature data from external API into memory
2. Both interfaces use the same `CreatureStore` for creature lookup and calculations
3. Calculations compare elemental charm damage against overflux/overpower charms
4. Results show breakeven resource requirements and maximum damage potential

### Key Calculations
- Elemental charms: 5% base damage, modified by creature resistances
- Overflux: 2.5% resource cost for variable damage
- Overpower: 5% resource cost for variable damage  
- Maximum damage cap: 8% of creature hitpoints

## Environment Variables
- `PORT`: Server port (default: 8000)