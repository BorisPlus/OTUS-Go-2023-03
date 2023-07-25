# Решение

## 0) «Точка входа», запускающая сервис

```bash

$ cd ./hw12_13_14_15_calendar/cmd/calendar/

/hw12_13_14_15_calendar/cmd/calendar$ go run ./

    Please set: '--config=<Path to configuration file>'

/hw12_13_14_15_calendar/cmd/calendar$ go run ./ --config="../../configs/config.yaml"
    
    INFO [2023-07-24 23:52:46] calendar is running...
    INFO [2023-07-24 23:52:46] Server.Start()
    ^C
    INFO [2023-07-24 23:52:47] Shutting down gracefully by signal.
    INFO [2023-07-24 23:52:47] Server.Stop()

/hw12_13_14_15_calendar/cmd/calendar$ go run ./ --config="../../configs/config.yaml"
    
    INFO [2023-07-24 23:52:55] calendar is running...
    INFO [2023-07-24 23:52:55] Server.Start()
    ^Z
    INFO [2023-07-24 23:52:57] Shutting down gracefully by signal.
    INFO [2023-07-24 23:52:57] Server.Stop()
    [1]+  Остановлен    go run ./ --config="../../configs/config.yaml"
```

## 1) Конфигурирование сервиса

Задействовал:

* "github.com/spf13/pflag"
* "github.com/spf13/viper"

Реализация:

* [config.go](internal/config/config.go)
* [main.go](cmd/calendar/main.go)

## 2) Логирование в сервисе

Свой логгер:

* [logger.go](internal/logger/logger.go)
* [logger_test.go](internal/logger/logger_test.go)

## 3) Работа с хранилищем

Реализации

* [In memory database by Golang](internal/storage/gomemory/storage.go)
* [Postgres](internal/storage/pgsqldtb/storage.go)
* [и их обобщенный тест](internal/storage/storage_test.go)

### Предварительно для Postgres

```sql
CREATE USER hw12user WITH PASSWORD 'hw12user';
CREATE DATABASE hw12calendar OWNER hw12user;
```

```text
# pg_hba.conf
host    all    hw12user    127.0.0.1/32    trust
```

Проверка корректности настройки:

```bash
psql -h localhost -p 5432 -d hw12calendar -U hw12user -W

goose -dir ./migrations/migrations/ postgres "user=hw12user password=hw12user host='127.0.0.1' database=hw12calendar" status

    2023/07/25 03:18:59     Applied At                  Migration
    2023/07/25 03:18:59     =======================================
    2023/07/25 03:18:59     Sun Jul 23 00:02:33 2023 -- 00000000000001_create_schema.sql
    2023/07/25 03:18:59     Pending                  -- 20230722000001_create_events.sql
```

## 4) Запуск простого HTTP-сервера

* [Сервер](internal/server/http/server.go)
* [Его мидлваре для "latency"](internal/server/http/middleware.go)
* [Тест](internal/server/http/server_test.go)

## 5) Юнит-тесты

Представлены в подпакетах модуля.

## 6) Makefile

* [Makefile](Makefile)

```bash
make build

    go build -v -o "./bin/calendar.goc" -ldflags "-X main.release="develop" -X main.buildDate=2023-07-25T00:11:15 -X main.gitHash=8ff86cc" ./cmd/calendar
    hw12_13_14_15_calendar/cmd/calendar

make run

    go build -v -o "./bin/calendar.goc" -ldflags "-X main.release="develop" -X main.buildDate=2023-07-25T00:12:22 -X main.gitHash=8ff86cc" ./cmd/calendar
    hw12_13_14_15_calendar/cmd/calendar
    "./bin/calendar.goc" --config ./configs/config.yaml
    INFO [2023-07-25 00:12:23] calendar is running...
    INFO [2023-07-25 00:12:23] Server.Start()
    ^C
    INFO [2023-07-25 00:12:31] Shutting down gracefully by signal.
    INFO [2023-07-25 00:12:31] Server.Stop()
    make: *** [Makefile:11: run] Прерывание

make test 

    go test -race ./internal/...
    ?       hw12_13_14_15_calendar/internal/app     [no test files]
    ?       hw12_13_14_15_calendar/internal/config  [no test files]
    ?       hw12_13_14_15_calendar/internal/interfaces      [no test files]
    ok      hw12_13_14_15_calendar/internal/logger  (cached)
    ?       hw12_13_14_15_calendar/internal/models  [no test files]
    ok      hw12_13_14_15_calendar/internal/server/http     (cached)
    ok      hw12_13_14_15_calendar/internal/storage (cached)
    ok      hw12_13_14_15_calendar/internal/storage/gomemory        0.031s
    ok      hw12_13_14_15_calendar/internal/storage/pgsqldtb        0.036s

make lint

    (which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /home/b/go/bin v1.50.1
    golangci-lint run --out-format=github-actions ./...

make migrate PG_DSN="user=hw12user password=hw12user host='127.0.0.1' database=hw12calendar"

    goose -dir ./migrations/migrations/ postgres "user=hw12user password=hw12user host='127.0.0.1' database=hw12calendar" up
    2023/07/25 03:13:51 OK   20230722000001_create_events.sql (17.42ms)
    2023/07/25 03:13:51 goose: no migrations to run. current version: 20230722000001
```
