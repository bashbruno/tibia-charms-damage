PORT = 8000

api:
	@PORT=$(PORT) go run ./cmd/api

cli:
	@go run ./cmd/cli

build-api:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/api ./cmd/api

build-cli:
	@go build -ldflags="-s -w" -o bin/tibia ./cmd/cli

test:
	@go test ./... -v

docker:
	@docker build -t tibia-charms .
	@docker run -e PORT=$(PORT) -p $(PORT):$(PORT) tibia-charms
