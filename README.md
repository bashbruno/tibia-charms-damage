# Tibia Charms Damage Calculator

A web/cli app for comparing the damage of overflux/overpower charms compared to elemental charms in Tibia.

## Prerequisites

- Go installed on your machine
- Internet connection

## Running Web locally

1. Run the application:

Running it with Go via the Makefile:

```bash
make api
```

or

Running it via the Dockerfile:

```bash
make docker
```

or

Running it directly with Go:

```bash
go run ./cmd/api
```

If you decide to go with the latter option, don't forget to set the `PORT` and `DATA_URL` environment variables in your shell - you can grab the default values from the Makefile.

The application will start on port 8000 by default. You can configure the port by setting the `PORT` variable in the Makefile/your shell.

## Build

```bash
go build -o <binary-name>
```
