# docker build --no-cache --build-arg MICROSERVICE=calendar --progress plain -f ./docker/Dockerfile --tag=hw15/calendar:v1 .
# docker build --no-cache --build-arg MICROSERVICE=sheduler --progress plain -f ./docker/Dockerfile --tag=hw15/sheduler:v1 .
# docker build --no-cache --build-arg MICROSERVICE=archiver --progress plain -f ./docker/Dockerfile --tag=hw15/archiver:v1 .
# docker build --no-cache --build-arg MICROSERVICE=sender   --progress plain -f ./docker/Dockerfile --tag=hw15/sender:v1 .
# docker builder prune
# docker build --build-arg MICROSERVICE=calendar --progress plain -f ./docker/Dockerfile --tag=hw15/calendar:v2 .
# docker build --build-arg MICROSERVICE=sheduler --progress plain -f ./docker/Dockerfile --tag=hw15/sheduler:v2 .
# docker build --build-arg MICROSERVICE=archiver --progress plain -f ./docker/Dockerfile --tag=hw15/archiver:v2 .
# docker build --build-arg MICROSERVICE=sender   --progress plain -f ./docker/Dockerfile --tag=hw15/sender:v2 .
# Fat image with sources
FROM golang:1.19 as build

ARG MICROSERVICE
ENV BIN_FILE /opt/${MICROSERVICE}/${MICROSERVICE}
ENV SRC_DIR /go/src
ENV CODE_DIR ${SRC_DIR}/hw12_13_14_15_calendar

# Feauture project file
WORKDIR ${SRC_DIR}
COPY docker/additional/${MICROSERVICE}.go.work go.work

# Go requirements
WORKDIR ${CODE_DIR}
COPY go.mod .
COPY go.sum .

WORKDIR ${CODE_DIR}/cmd/${MICROSERVICE}
COPY /cmd/${MICROSERVICE}/go.mod .
COPY /cmd/${MICROSERVICE}/go.sum .

RUN go mod download -x

# Source-code files
WORKDIR ${CODE_DIR}/cmd/${MICROSERVICE}
COPY cmd/${MICROSERVICE}/main.go ./ 

WORKDIR ${CODE_DIR}/internal
COPY ./internal .

# Special compile with out additional info
WORKDIR ${SRC_DIR}
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o ${BIN_FILE} ./hw12_13_14_15_calendar/cmd/${MICROSERVICE}/

# Slim image
FROM alpine:3.9 as alpine

LABEL ORGANIZATION="OTUS Online Education"
LABEL MODULE="${MICROSERVICE}"
LABEL MAINTAINERS="admin@localhost"

ARG MICROSERVICE
ENV BIN_FILE "/opt/${MICROSERVICE}/${MICROSERVICE}"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/${MICROSERVICE}/config.yaml
COPY configs/${MICROSERVICE}.yaml ${CONFIG_FILE}

# One imege - one process
ENV COMMAND "/opt/${MICROSERVICE}/${MICROSERVICE} --config /etc/${MICROSERVICE}/config.yaml"
RUN echo ${COMMAND}

ENTRYPOINT ${COMMAND}
