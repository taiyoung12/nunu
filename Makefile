.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest

.PHONY: db-up
db-up:
	cd ./deploy/docker-compose && docker compose up -d

.PHONY: db-down
db-down:
	cd ./deploy/docker-compose && docker compose down

.PHONY: run
run:
	go run ./cmd/server -conf config/local.yml

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

.PHONY: wire
wire:
	cd cmd/server/wire && wire

.PHONY: test
test:
	go test ./... -v

.PHONY: lint
lint:
	go vet ./...
