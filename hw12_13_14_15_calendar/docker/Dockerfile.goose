
# Golang
FROM golang:1.19

ENV NOTATION postgres
ENV DSN "host=hw15_postgres_prod port=9432 dbname=hw15 user=hw15user password=hw15user sslmode=disable"
ENV ACTION up

# Prepare
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Feauture project file
WORKDIR /goose
COPY migrations/migrations/ migrations/

# One imege - one process
ENTRYPOINT goose -dir /goose/migrations $NOTATION "$DSN" $ACTION 
# docker build --no-cache --progress plain -f ./docker/Dockerfile.goose -d -e NOTATION=postgres -e ACTION=up -e DSN="host=hw15_postgres_prod user=hw15user password=hw15user dbname=hw15 sslmode=disable port=5432" .
