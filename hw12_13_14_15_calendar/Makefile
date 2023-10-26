BIN_DIR="./bin"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# General

protoc:
	protoc --go-grpc_out=./internal/server/rpc/rpcapi --go_out=./internal/server/rpc/rpcapi ./internal/server/rpc/protofiles/calendar.proto

test:
	go clean -testcache
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run --out-format=github-actions ./...

fixlint:
	golangci-lint run --out-format=github-actions --fix ./...

# Build

build:
	make MICROSERVICE=calendar 	build-microservice
	make MICROSERVICE=sheduler 	build-microservice
	make MICROSERVICE=archiver 	build-microservice
	make MICROSERVICE=sender 	build-microservice
	make MICROSERVICE=rmqdevops build-microservice
	make MICROSERVICE=dataset 	build-microservice

build-microservice: # make MICROSERVICE=calendar build-microservice
	go build -v -o "$(BIN_DIR)/$(MICROSERVICE)"  -ldflags "$(LDFLAGS)" ./cmd/$(MICROSERVICE)

go-run: # make MICROSERVICE=checker go-run
	go run ./cmd/$(MICROSERVICE)/ --config ./configs/$(MICROSERVICE).yaml

run-microservice: build # make MICROSERVICE=calendar run-microservice
	"./bin/$(MICROSERVICE)" --config ./configs/$(MICROSERVICE).yaml

version: build  # make MICROSERVICE=calendar build-microservice
	"./bin/$(MICROSERVICE)" version

migrate-sql-up: # make PG_DSN="host=hw15_postgres_prod user=hw15user password=hw15user dbname=hw15 sslmode=disable port=5432" migrate-sql-up
	goose -dir ./migrations/migrations/ postgres "$(PG_DSN)" up

migrate-sql-down:
	goose -dir ./migrations/migrations/ postgres "$(PG_DSN)" down

migrate-rmq-up:
	go run ./cmd/management/rmqdevops/ --config ./configs/rmqdevops.yaml

migrate-rmq-down:
	go run ./cmd/management/rmqdevops/ --drop --config ./configs/rmqdevops.yaml

migrate-rmq-up-with-drop:
	go run ./cmd/management/rmqdevops/ --with-drop --config ./configs/rmqdevops.yaml

dataset-example:
	go run ./cmd/dataset/ --config ./configs/calendar.yaml

# 

.PHONY: build run build-img run-img version test lint
