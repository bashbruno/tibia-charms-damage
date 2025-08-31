# Tibia Charms Damage Calculator

A web/cli app for comparing the damage of overflux/overpower charms compared to elemental charms in Tibia.

## Prerequisites

- Go installed on your machine
- Internet connection

## Web Version

You can run it locally using one of these options:

1. Via Makefile:

```bash
make api
```

2. Via Docker:

```bash
docker compose up -d
```

3. Via Go:

```bash
go run ./cmd/api
```

The application will then start on port 8000 by default, which can be changed in the Makefile or by setting the `PORT` variable in your shell.

## CLI Version

You can run the CLI locally using one of these options:

1. Via Makefile:

```bash
make cli
```

2. Via Go:

```bash
go run ./cmd/cli
```

### Building the CLI Binary

#### Linux/MacOS:

Build the binary:

```bash
go build -o bin/tibia ./cmd/cli
```

You can then run the binary with `./bin/tibia`

To be able to run it from any terminal window without specifying the full path to the binary, you can set `GOBIN` and add it to your `PATH`:

In your `.bashrc`, `.zshrc` or equivalent:

```bash
export PATH=$PATH:$HOME/go/bin
```

Then move the created binary to your Go bin directory (run this from the project's folder):

```bash
mv ./bin/tibia ~/go/bin/
```

Then run from anywhere:

```bash
tibia
```

#### Windows:

Build the binary:

```bash
go build -o bin/tibia.exe ./cmd/cli
```

You can then run the binary with `.\bin\tibia.exe`

To be able to run it from any terminal window without specifying the full path to the binary, ensure `%USERPROFILE%\go\bin` is on `PATH`:

1. Open Start Menu -> Edit environment variables
2. Under User variables → find Path → Edit → Add:

```bash
%USERPROFILE%\go\bin
```

Then move the created binary (run this from the project's folder):

```bash
move .\bin\tibia.exe %USERPROFILE%\go\bin\
```

Then run from anywhere:

```bash
tibia
```
