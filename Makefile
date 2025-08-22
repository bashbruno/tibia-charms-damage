PORT = 8000
DATA_URL = https://raw.githubusercontent.com/mathiasbynens/tibia-json/main/data/bestiary.json

api:
	@PORT=:$(PORT) DATA_URL=$(DATA_URL) go run ./cmd/api

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/api ./cmd/api

docker:
	@docker build --build-arg PORT_ARG=$(PORT) --build-arg DATA_URL_ARG="$(DATA_URL)" -t tibia-charms .
	@docker run -p $(PORT):$(PORT) tibia-charms
