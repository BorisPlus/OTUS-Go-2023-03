
# Fat image with sources
FROM golang:1.19 as development

ARG MICROSERVICE
ARG LDFLAGS
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
RUN CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o /opt/service/service ./hw12_13_14_15_calendar/cmd/${MICROSERVICE}/

# Slim image
FROM alpine:3.9 as production

ARG MICROSERVICE
ARG CONFIG_FILE

LABEL ORGANIZATION="OTUS Online Education"
LABEL MODULE="${MICROSERVICE}"
LABEL MAINTAINERS="admin@localhost"

COPY --from=development /opt/service/service /opt/service/service

COPY configs/${CONFIG_FILE} "/etc/service/config.yaml"
RUN echo "CONFIG_FILE=${CONFIG_FILE}"
# One imege - one process
ENTRYPOINT [ "/opt/service/service", "--config", "/etc/service/config.yaml" ]
