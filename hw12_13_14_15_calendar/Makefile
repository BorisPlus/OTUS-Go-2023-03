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

migrate-sql-up: # make PG_DSN="host=hw15_postgres_prod user=hw15user password=hw15user dbname=hw15 sslmode=disable port=5432" migrate-sql-up
	goose -dir ./migrations/migrations/ postgres "$(PG_DSN)" up

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
	docker network create hw15_network_dev --subnet=172.133.0.0/24 2> /dev/null || true

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
	docker network rm hw15_network_dev  2> /dev/null || true

docker-build-microservice: # Example: make MICROSERVICE=rmqdevops CONFIG_FILE=rmqdevops.compose.yaml docker-build-microservice
	docker build \
		--build-arg=LDFLAGS="${LDFLAGS}" \
		--build-arg=MICROSERVICE="$(MICROSERVICE)" \
		--build-arg=CONFIG_FILE="$(CONFIG_FILE)" \
		-t hw15_$(MICROSERVICE):dev \
		-f docker/Dockerfile .

docker-build:
	MICROSERVICE=calendar 	CONFIG_FILE=calendar.yaml	make docker-build-microservice
	MICROSERVICE=sheduler 	CONFIG_FILE=sheduler.yaml	make docker-build-microservice
	MICROSERVICE=archiver 	CONFIG_FILE=archiver.yaml	make docker-build-microservice
	MICROSERVICE=sender   	CONFIG_FILE=sender.yaml		make docker-build-microservice
	MICROSERVICE=rmqdevops	CONFIG_FILE=rmqdevops.yaml	make docker-build-microservice

docker-run-calendar:
	MICROSERVICE=calendar 	CONFIG_FILE=calendar.yaml	make docker-build-microservice
	docker run -it --name=hw15_calendar_dev --net hw15_network_dev --ip 172.133.0.2 -p 8888:8080   -p 5000:5000 -d hw15_calendar:dev 

docker-run-rabbitmq:
	docker run     --name=hw15_rabbitmq_dev --net hw15_network_dev --ip 172.133.0.3 -p 15672:15672 -p 5672:5672 -d -e RABBITMQ_DEFAULT_USER=hw15user -e RABBITMQ_DEFAULT_PASS=hw15user rabbitmq:3.10.7-management

docker-run-microservice: docker-build-microservice # Example: make MICROSERVICE=rmqdevops CONFIG_FILE=rmqdevops.yaml docker-run-microservice
	docker run -it --name=hw15_$(MICROSERVICE)_dev -d hw15_$(MICROSERVICE):dev 

docker-run: docker-deconfigure \
			docker-configure
							make docker-run-rabbitmq
							sleep 20
	MICROSERVICE=rmqdevops	CONFIG_FILE=rmqdevops.yaml make docker-run-microservice
							make docker-run-calendar
							sleep 20
	MICROSERVICE=sheduler 		CONFIG_FILE=sheduler.yaml make docker-run-microservice
							sleep 20
	MICROSERVICE=archiver		CONFIG_FILE=archiver.yaml make docker-run-microservice
							sleep 20
	MICROSERVICE=sender 		CONFIG_FILE=sender.yaml make docker-run-microservice
	
docker-start-microservice: # Example: make MICROSERVICE=rmqdevops docker-start-microservice
	docker start hw15_$(MICROSERVICE)_dev

docker-start:
	MICROSERVICE=calendar 	make docker-start-microservice
	MICROSERVICE=sheduler	make docker-start-microservice
	MICROSERVICE=archiver	make docker-start-microservice
	MICROSERVICE=sender		make docker-start-microservice
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

docker-compose-deconfigure:
	docker rmi -f hw15_checker_prod		2> /dev/null || true
	docker rmi -f hw15_dataset_prod		2> /dev/null || true
	docker rmi -f hw15_archiver_prod	2> /dev/null || true
	docker rmi -f hw15_sender_prod		2> /dev/null || true
	docker rmi -f hw15_sheduler_prod	2> /dev/null || true
	docker rmi -f hw15_rmqdevops_prod	2> /dev/null || true
	docker rmi -f hw15_rabbitmq_prod	2> /dev/null || true
	docker rmi -f hw15_goose_prod		2> /dev/null || true
	docker rmi -f hw15_pgadmin_prod		2> /dev/null || true
	docker rmi -f hw15_postgres_prod	2> /dev/null || true
	docker rm  -f hw15_checker_prod		2> /dev/null || true
	docker rm  -f hw15_dataset_prod		2> /dev/null || true
	docker rm  -f hw15_archiver_prod	2> /dev/null || true 
	docker rm  -f hw15_sender_prod		2> /dev/null || true
	docker rm  -f hw15_sheduler_prod	2> /dev/null || true
	docker rm  -f hw15_rmqdevops_prod	2> /dev/null || true
	docker rm  -f hw15_rabbitmq_prod	2> /dev/null || true
	docker rm  -f hw15_goose_prod		2> /dev/null || true
	docker rm  -f hw15_pgadmin_prod		2> /dev/null || true
	docker rm  -f hw15_postgres_prod	2> /dev/null || true
	docker rm  -f hw15_calendar_prod	2> /dev/null || true
	docker network rm hw15_public_prod	2> /dev/null || true
	docker network rm hw15_private_prod	2> /dev/null || true

docker-compose-up: docker-compose-deconfigure
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker/docker-compose.yaml up --build -d --remove-orphans

docker-compose-down:
	docker-compose -f ./docker/docker-compose.yaml down

# Internal test

integration_test: docker-compose-deconfigure
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker/docker-compose.yaml -f ./docker/docker-compose.checker.yaml up --exit-code-from checker --abort-on-container-exit

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
