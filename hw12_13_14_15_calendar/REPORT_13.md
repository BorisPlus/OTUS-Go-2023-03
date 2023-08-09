# Дз 13

```bash
go run ./hw12_13_14_15_calendar/cmd/calendar/ --config=./hw12_13_14_15_calendar/configs/config.yaml
```

> ```text
> -v - для verbose
> ```

Create

```bash
curl -d "{\"title\":\"title\",\"startat\":\"2021-02-18T21:54:42+02:00\",\"duration\":0,\"description\": \"description\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -X POST "http://localhost:8080/api/events/create"
```

Read

```bash
curl -H "Content-Type: application/json" "http://localhost:8080/api/events/1"
```

Update

```bash
curl -d "{\"PK\":2,\"title\":\"title\",\"startat\":\"2023-02-18T21:54:42+02:00\",\"duration\":0,\"description\": \"description\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -X PATCH "http://localhost:8080/api/events/1/update"
```

Read

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/1"
```

Delete

```bash
curl -v -H "Content-Type: application/json" -X DELETE "http://localhost:8080/api/events/2/delete"
```

Read

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/2"
```

## GRPC

[GRPC](https://grpc.io/docs/languages/go/quickstart/)

```bash
go get google.golang.org/protobuf@v1.30.0
go get google.golang.org/grpc@v1.55.0
go list
protoc --go_out=plugins=grpc:pkg/api ./calendar/calendar.proto
export PATH="$PATH:$(go env GOPATH)/bin"

cd hw12_13_14_15_calendar/internal/protobuf/calendar/

protoc --go_out=plugins:grpc:pkg/api ./internal/protobuf/protofiles/calendar.proto
protoc --go_out=. --go-grpc_out=./api --go_out=plugins=grpc,import_path=hw12_13_14_15_calendar/internal/models:/tmp/dest ./protofiles/calendar.proto

protoc ---go-grpc_out=./api --go_out=plugins=grpc,import_path=hw12_13_14_15_calendar/internal/models:/tmp/dest ./protofiles/calendar.proto
protoc --go_out=./api --go-grpc_out=./api ./protofiles/calendar.proto


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/protobuf/calendar/calendar.proto
```