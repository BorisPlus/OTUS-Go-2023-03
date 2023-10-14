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

run-microservice: build # make MICROSERVICE=calendar run-microservice
	"./bin/$(MICROSERVICE)" --config ./configs/$(MICROSERVICE).yaml 

version: build  # make MICROSERVICE=calendar build-microservice
	"./bin/$(MICROSERVICE)" version

# migrate-sql-up: 
# 	goose -dir ./migrations/migrations/ postgres "$(PG_DSN)" up

# migrate-sql-down: 
# 	goose -dir ./migrations/migrations/ postgres "$(PG_DSN)" down

# migrate-rmq-up: 
# 	go run ./cmd/management/rmqdevops/ --config ./configs/rmqdevops.yaml

# migrate-rmq-down: 
# 	go run ./cmd/management/rmqdevops/ --drop --config ./configs/rmqdevops.yaml

# migrate-rmq-up-with-drop: 
# 	go run ./cmd/management/rmqdevops/ --with-drop --config ./configs/rmqdevops.yaml

# dataset-example: 
# 	go run ./cmd/dataset/ --config ./configs/calendar.yaml

# Dockerize

docker-configure:
	docker network create --subnet=172.18.0.0/24 calendar 2> /dev/null || true

docker-deconfigure:
	docker rmi -f hw15_calendar         2> /dev/null || true
	docker rmi -f hw15_rabbitmq         2> /dev/null || true
	docker rmi -f hw15_sheduler         2> /dev/null || true
	docker rmi -f hw15_archiver         2> /dev/null || true
	docker rmi -f hw15_sender           2> /dev/null || true
	docker rmi -f hw15_rmqdevops        2> /dev/null || true
	docker rm  -f hw15_calendar_dev     2> /dev/null || true
	docker rm  -f hw15_rabbitmq_dev     2> /dev/null || true
	docker rm  -f hw15_sheduler_dev     2> /dev/null || true
	docker rm  -f hw15_archiver_dev     2> /dev/null || true
	docker rm  -f hw15_sender_dev       2> /dev/null || true
	docker rm  -f hw15_rmqdevops_dev    2> /dev/null || true
	docker network rm calendar          2> /dev/null || true

docker-build-microservice: # Example: make MICROSERVICE=calendar docker-build-microservice
	docker build \
		--build-arg=LDFLAGS="${LDFLAGS}" \
		--build-arg=MICROSERVICE="$(MICROSERVICE)" \
		-t hw15_$(MICROSERVICE):dev \
		-f docker/Dockerfile .

docker-build:
	MICROSERVICE=calendar 	make docker-build-microservice
	MICROSERVICE=sheduler 	make docker-build-microservice
	MICROSERVICE=archiver 	make docker-build-microservice
	MICROSERVICE=sender   	make docker-build-microservice
	MICROSERVICE=rmqdevops	make docker-build-microservice

docker-run-calendar:
	MICROSERVICE=calendar 	make docker-build-microservice
	docker run -it --name=hw15_calendar_dev --net calendar --ip 172.18.0.2 -p 8888:8080   -p 5000:5000 -d hw15_calendar:dev 

docker-run-rabbitmq:
	docker run     --name=hw15_rabbitmq_dev --net calendar --ip 172.18.0.3 -p 15672:15672 -p 5672:5672 -d rabbitmq:3.10.7-management

docker-run-microservice: docker-build-microservice # Example: make MICROSERVICE=rmqdevops docker-run-microservice
	docker run -it --name=hw15_$(MICROSERVICE)_dev -d hw15_$(MICROSERVICE):dev 

docker-run: docker-deconfigure docker-configure
							make docker-run-calendar
							make docker-run-rabbitmq
	MICROSERVICE=rmqdevops	make docker-run-microservice
	MICROSERVICE=sheduler 	make docker-run-microservice
	MICROSERVICE=archiver 	make docker-run-microservice
	MICROSERVICE=sender   	make docker-run-microservice
	
docker-start-microservice: # Example: make MICROSERVICE=calendar docker-start-microservice
	docker start hw15_$(MICROSERVICE)_dev

docker-start:
	MICROSERVICE=calendar 	make docker-start-microservice
	MICROSERVICE=sheduler 	make docker-start-microservice
	MICROSERVICE=archiver 	make docker-start-microservice
	MICROSERVICE=sender   	make docker-start-microservice
	MICROSERVICE=rmqdevops	make docker-start-microservice

docker-stop-microservice: # Example: make MICROSERVICE=calendar docker-stop-microservice
	docker stop hw15_$(MICROSERVICE)_dev

docker-stop:
	MICROSERVICE=calendar 	make docker-stop-microservice
	MICROSERVICE=sheduler 	make docker-stop-microservice
	MICROSERVICE=archiver 	make docker-stop-microservice
	MICROSERVICE=sender   	make docker-stop-microservice
	MICROSERVICE=rmqdevops	make docker-stop-microservice
	
# Composing

docker-compose-up:
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker/docker-compose.yaml up -d

docker-compose-down:
	docker-compose -f ./docker/docker-compose.yaml down

# Be careful

prune:
	docker stop -f $$(docker ps -a -q) || true
	docker rm -f $$(docker ps -a -q) || true
	docker rmi -f $$(docker images -q) || true
	docker image prune -f
	docker network prune -f
	docker container prune -f

#

.PHONY: build run build-img run-img version test lint
