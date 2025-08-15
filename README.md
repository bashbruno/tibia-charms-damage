# Tibia Charms Damage Calculator

A web/cli app for comparing the damage of overflux/overpower charms compared to elemental charms in Tibia.

## Prerequisites

- Go installed on your machine

## Running Locally

1. Run the application:

```bash
make run
```

or

```bash
go run ./cmd/api
```

The application will start on port 8000 by default. You can configure the port by setting the `ADDR` environment variable.

## Build

```bash
go build -o <binary-name>
```
