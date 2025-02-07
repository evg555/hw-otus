CALENDAR := "./bin/calendar"
SCHEDULER = "./bin/calendar_scheduler"
SENDER = "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

GOOSE=goose
MIGRATIONS_DIR=migrations

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

build:
	go build -v -o $(CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(CALENDAR) -config ./configs/calendar_config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(CALENDAR) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

install-protoc-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate: install-protoc-deps
	go generate ./...

.PHONY: build run build-img run-img version test lint migrate-up migrate-down generate
