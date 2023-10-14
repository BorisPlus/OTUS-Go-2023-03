
# Golang
FROM golang:1.19

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Feauture project file
WORKDIR /goose
COPY migrations/migrations/ migrations/

# One imege - one process
ENTRYPOINT goose -dir /goose/migrations $NOTATION "$DSN" $ACTION 
