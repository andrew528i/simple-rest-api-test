.PHONY: install build run test migrate_init migrate migrate_down

API_NAME?=backend
DB_PORT?=5555
DB_HOST?=localhost
DB_CONNECTION_URL?=postgres://$(API_NAME):$(API_NAME)@$(DB_HOST):$(DB_PORT)/$(API_NAME)?sslmode=disable

install:
	go mod tidy
	go mod download

build: install
	go build -o ./bin/$(API_NAME) ./api/
	go build -o ./bin/migrations ./migrations/

run: build
	DB_CONNECTION_URL=$(DB_CONNECTION_URL) \
	./bin/$(API_NAME)

test:
	go test ./api/...

migrate: build
	DB_CONNECTION_URL=$(DB_CONNECTION_URL) \
	./bin/migrations migrate
