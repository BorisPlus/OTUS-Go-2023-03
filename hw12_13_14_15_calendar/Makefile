BIN := "./bin/calendar.goc"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o "./bin/calendar.goc" -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o "./bin/sheduler.goc" -ldflags "$(LDFLAGS)" ./cmd/sheduler
	go build -v -o "./bin/archiver.goc" -ldflags "$(LDFLAGS)" ./cmd/archiver
	go build -v -o "./bin/sender.goc" -ldflags "$(LDFLAGS)" ./cmd/sender
	go build -v -o "./bin/rmqdevops.goc" -ldflags "$(LDFLAGS)" ./cmd/management/rmqdevops
	go build -v -o "./bin/dataset.goc" -ldflags "$(LDFLAGS)" ./cmd/management/dataset

run: build
	$(BIN) --config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

compose-up:
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker/docker-compose.yaml up -d

compose-down:
	docker-compose -f ./docker/docker-compose.yaml down

protoc:
	protoc --go-grpc_out=./internal/server/rpc/api --go_out=./internal/server/rpc/api ./internal/server/rpc/protofiles/calendar.proto

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

migrate-sql-up: 
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
	go run ./cmd/management/dataset/ --config ./configs/calendar.yaml

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run --out-format=github-actions ./...

fixlint:
	golangci-lint run --out-format=github-actions --fix ./...

prune:
	docker image prune -f
	docker network prune -f
	docker container prune -f
	docker rm -f "$(docker ps -a -q)" || true
	docker rmi -f "$(docker images -q)" || true

.PHONY: build run build-img run-img version test lint
