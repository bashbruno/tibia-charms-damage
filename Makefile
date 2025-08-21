api:
	@go run ./cmd/api

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api 
