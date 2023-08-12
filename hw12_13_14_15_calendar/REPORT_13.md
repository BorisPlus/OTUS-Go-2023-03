# Дз 13


## Тестирование

```bash
go test -v ./internal/models/.
=== RUN   TestLogger
Unmarshaled. OK.
--- PASS: TestLogger (0.00s)
PASS
ok      hw12_13_14_15_calendar/internal/models  0.008s




```

## Из задач ДЗ

## HTTP API

Для HTTP API необходимо:

* расширить "hello-world" сервер из [ДЗ №12](./12_README.md) до полноценного API;
* создать отдельный пакет для кода HTTP сервера;
* реализовать хэндлеры, при необходимости выделив структуры запросов и ответов;
* сохранить логирование запросов, реализованное в [ДЗ №12](./12_README.md).

Описание решение:

* Файлы HTTP-сервера размещены в директории [internal/server/http/](internal/server/http/)
* Для RESTful API разработан модуль роутинга, поддерживающий регулярные выражения для путей [regexphandlers](internal/server/http/regexphandlers) с сохранением возможности задействовать ранее реализованный `middleware` для каждого пути отдельно. На примере файла [api/router.go](internal/server/http/api/router.go) тезисно:

```go
*regexped.NewRegexpHandler(
    `/api/version`,
    none,
    middleware.Instance().Listen(commonHandlers.VersionHandler{}),
),
*regexped.NewRegexpHandler(
    `/api/events/`,
    none,
    middleware.Instance().Listen(apiEvents.EventsListHandler{Logger: logger, App: app}),
),
*regexped.NewRegexpHandler(
    `/api/events/create`,
    none,
    apiEvents.EventsCreateHandler{Logger: logger, App: app},
),
*regexped.NewRegexpHandler(
    `/api/events/{numeric}`,
    id,
    apiEvents.EventsGetHandler{Logger: logger, App: app},
),
*regexped.NewRegexpHandler(
    `/api/events/{numeric}/update`,
    id,
    apiEvents.EventsUpdateHandler{Logger: logger, App: app},
),
*regexped.NewRegexpHandler(
    `/api/events/{numeric}/delete`,
    id,
    apiEvents.EventsDeleteHandler{Logger: logger, App: app},
),
```

* дополнительная особенность имплементации - единый формат [ответа](internal/server/http/api/api_response/core.go)
* тестирование сервера [server_test.go](internal/server/http/server_test.go)

```bash
go test -v ./internal/server/http/server.go ./internal/server/http/server_test.go 

=== RUN   TestServerStopNotStarted
INFO [2023-08-11 21:41:28] HTTPServer.Stop()
--- PASS: TestServerStopNotStarted (0.00s)

=== RUN   TestServerStopNormally
INFO [2023-08-11 21:41:28] HTTPServer.Start()
INFO [2023-08-11 21:41:31] HTTPServer.Stop()
--- PASS: TestServerStopNormally (3.00s)

=== RUN   TestServerStopBySignalNoWait
INFO [2023-08-11 21:41:31] HTTPServer.Start()
INFO [2023-08-11 21:41:31] HTTPServer - Graceful Shutdown
INFO [2023-08-11 21:41:31] HTTPServer.Stop()
--- PASS: TestServerStopBySignalNoWait (0.00s)

=== RUN   TestServerStopBySignalWithWait
INFO [2023-08-11 21:41:31] HTTPServer.Start()
INFO [2023-08-11 21:41:34] HTTPServer - Graceful Shutdown
INFO [2023-08-11 21:41:34] HTTPServer.Stop()
--- PASS: TestServerStopBySignalWithWait (3.00s)

=== RUN   TestServerStopByCancelNoWait
INFO [2023-08-11 21:41:34] HTTPServer.Start()
INFO [2023-08-11 21:41:34] HTTPServer - Graceful Shutdown
INFO [2023-08-11 21:41:34] HTTPServer.Stop()
--- PASS: TestServerStopByCancelNoWait (0.00s)

=== RUN   TestServerStopByCancelWithWait
INFO [2023-08-11 21:41:34] HTTPServer.Start()
INFO [2023-08-11 21:41:37] HTTPServer - Graceful Shutdown
INFO [2023-08-11 21:41:37] HTTPServer.Stop()
--- PASS: TestServerStopByCancelWithWait (3.00s)

=== RUN   TestServerCode
INFO [2023-08-11 21:41:40] {StatusCode:418 UserAgent:Go-http-client/1.1 ClientIPAddress:127.0.0.1:47808 HTTPMethod:GET HTTPVersion:HTTP/1.1 URLPath:/ StartAt:2023-08-12 00:41:40.295600571 +0300 MSK m=+12.020231410 Latency:4.7µs}
OK. StatusCode '418'
--- PASS: TestServerCode (3.01s)
PASS
ok      command-line-arguments  12.027s
```

* тестирование RESTful API ветки сервера [server_api_test.go](internal/server/http/server_api_test.go)

```bash
go test -v ./internal/server/http/server.go ./internal/server/http/server_api_test.go 

=== RUN   TestServerAPICreatePKSequence
INFO [2023-08-12 06:23:51] HTTPServer.Start()
OK: get event PK 1
OK: get event PK 2
OK: get event PK 3
INFO [2023-08-12 06:23:51] HTTPServer.Stop()
--- PASS: TestServerAPICreatePKSequence (0.01s)

=== RUN   TestServerAPIVersion
INFO [2023-08-12 06:23:51] HTTPServer.Start()
OK: {"method":"api.version","error":"","data":{"Version":"1.0.0"}}
INFO [2023-08-12 06:23:51] HTTPServer.Stop()
--- PASS: TestServerAPIVersion (0.00s)

PASS
ok      command-line-arguments  0.024s
```

### GRPC API

Для GRPC API необходимо:

* создать отдельную директорию для Protobuf спецификаций;
* создать Protobuf файлы с описанием всех методов API, объектов запросов и ответов (
т.к. объект Event будет использоваться во многих ответах разумно выделить его в отдельный message);
* создать отдельный пакет для кода GRPC сервера;
* добавить в Makefile команду `generate`; `make generate` - вызывает `go generate`, которая в свою очередь
генерирует код GRPC сервера на основе Protobuf спецификаций;
* написать код, связывающий GRPC сервер с методами доменной области (бизнес логикой);
* логировать каждый запрос по аналогии с HTTP API.

Описание решение:

* Файлы GRPC-сервера размещены в директории [internal/server/rpc/](internal/server/rpc/);
* Файлы Protobub-спецификации размещены в директории [internal/server/rpc/protofiles](internal/server/rpc/protofiles);
* `make protoc` генерирует код GRPC-сущностей на Go на основе Protobuf спецификаций [internal/server/rpc/api](internal/server/rpc/api);
* Имплементация клиента [internal/server/rpc/client/client.go](internal/server/rpc/client/client.go);
* Имплементация сервера [internal/server/rpc/server/server.go](internal/server/rpc/server/server.go), где реализация логгирования протокола для Unary и Stream задействует интерфейс ранее разработанного логгера;
* Тестирование имплементаций [internal/server/rpc/client/client.go](internal/server/rpc/integration/integration_test.go):

```bash
go test -v ./internal/server/rpc/integration_test/integration_test.go 

=== RUN   TestIntegration
INFO [2023-08-12 08:55:18] server listening at localhost:5000
INFO [2023-08-12 08:55:18] GRPCServer.Start()
INFO [2023-08-12 08:55:23] grpcClient.Client{}
INFO [2023-08-12 08:55:23] localhost:5000
INFO [2023-08-12 08:55:23] UnaryInterceptor: "/calendar.Application/CreateEvent" <-- OBJECT{title:"Title 1"}
INFO [2023-08-12 08:55:23] createdEvent1.PK = 1
INFO [2023-08-12 08:55:23] UnaryInterceptor: "/calendar.Application/CreateEvent" <-- OBJECT{title:"Title 2"}
INFO [2023-08-12 08:55:23] createdEvent2.PK = 2
INFO [2023-08-12 08:55:23] UnaryInterceptor: "/calendar.Application/DeleteEvent" <-- OBJECT{p_k:2  title:"Title 2"  start_at:{}}
INFO [2023-08-12 08:55:23] deletedEvent2.PK = 2
INFO [2023-08-12 08:55:23] UnaryInterceptor: "/calendar.Application/CreateEvent" <-- OBJECT{title:"Title 3"}
INFO [2023-08-12 08:55:23] createdEvent3.PK = 3
INFO [2023-08-12 08:55:23] UnaryInterceptor: "/calendar.Application/ReadEvent" <-- OBJECT{pk:3}
INFO [2023-08-12 08:55:23] pbEvent3Copy.Title = Title 3
INFO [2023-08-12 08:55:23] StreamInterceptor: "/calendar.Application/ListEvents" %!s(MISSING)
INFO [2023-08-12 08:55:23] [p_k:1  title:"Title 1"  start_at:{} p_k:3  title:"Title 3"  start_at:{}]
INFO [2023-08-12 08:55:23] RPCServer - Graceful Shutdown
INFO [2023-08-12 08:55:23] GRPCServer.GracefulStop()
INFO [2023-08-12 08:55:23] GracefulStop
INFO [2023-08-12 08:55:23] GRPCServer.GracefulStop()
--- PASS: TestIntegration (5.02s)

=== RUN   TestInterceptorLogging
It's OK.
--- PASS: TestInterceptorLogging (0.01s)
PASS
ok      command-line-arguments  5.053s
```

## Запуск конечной реализации

```bash
go run ./cmd/calendar/ --config ./configs/config.yaml 
2023/08/12 12:12:26 HTTP Config - {Host:localhost Port:8080 ReadTimeout:5s ReadHeaderTimeout:5s WriteTimeout:5s MaxHeaderBytes:1048576}
2023/08/12 12:12:26 RPC Config - {Host:localhost Port:5000}
2023/08/12 12:12:26 Println Calendar is running...
INFO [2023-08-12 09:12:26] calendar is running...
INFO [2023-08-12 09:12:26] HTTPServer.Start()
INFO [2023-08-12 09:12:26] GRPCServer.Start()
^C
INFO [2023-08-12 09:12:37] GRPCServer.GracefulStop()
INFO [2023-08-12 09:12:37] RPCServer - Graceful Shutdown
INFO [2023-08-12 09:12:37] GRPCServer.GracefulStop()
INFO [2023-08-12 09:12:37] HTTPServer.Stop()
INFO [2023-08-12 09:12:37] HTTPServer - Graceful Shutdown
INFO [2023-08-12 09:12:37] HTTPServer.Stop()
INFO [2023-08-12 09:12:37] Complex Shutting down was done gracefully by signal.
```

> Замечание: Видно, что методы `HTTPServer.Stop()` и `GRPCServer.GracefulStop()` запускаются дважды, что связано с особенностями автоматической остановки процессов при завершении контекста. Это можно победить, внеся дополнительную логику с `sync.Once` внутри объектов имплементаций, но в целях наглядности для данного конечного решения данная особенность оставлена.

## Заметки

### Тест HTTP API вручную cURL

* Create - HTTP POST "/api/events/create"
* Read - HTTP GET "/api/events/{NUMBER}" и "/api/events/"
* Update - HTTP PATCH "/api/events/{NUMBER}/update"
* Delete - HTTP DELETE "/api/events/{NUMBER}/delete"
* Вызовы с несоотвествующим HTTP методом дают ошибку

```bash
go run ./hw12_13_14_15_calendar/cmd/calendar/ --config=./hw12_13_14_15_calendar/configs/config.yaml
```

Create

```bash
curl -v -d "{\"title\":\"title\",\"startat\":\"2023-08-11T21:54:42+02:00\",\"duration\":0,\"description\": \"description\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -POST "http://localhost:8080/api/events/create"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /api/events/create HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> Content-Length: 129
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 09:56:13 GMT
< Content-Length: 193
< 
* Connection #0 to host localhost left intact
{"method":"api.events.create","error":"","data":{"item":{"PK":1,"title":"title","startat":"2023-08-11T21:54:42+02:00","duration":0,"description":"description","owner":"owner","notifyearly":0}}}
```

Read

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/1"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/events/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 09:57:21 GMT
< Content-Length: 190
< 
* Connection #0 to host localhost left intact
{"method":"api.events.get","error":"","data":{"item":{"PK":1,"title":"title","startat":"2023-08-11T21:54:42+02:00","duration":0,"description":"description","owner":"owner","notifyearly":0}}}
```

Update (на примере даты, c противоречием параметров из Url `api/events/1/update` и Body `PK:2`), важен Url параметр в REST API:

```bash
curl -v -d "{\"PK\":2,\"title\":\"title\",\"startat\":\"2024-02-18T21:54:42+02:00\",\"duration\":0,\"description\": \"description\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -X PATCH "http://localhost:8080/api/events/1/update"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> PATCH /api/events/1/update HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> Content-Length: 136
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 09:57:59 GMT
< Content-Length: 193
< 
* Connection #0 to host localhost left intact
{"method":"api.events.update","error":"","data":{"item":{"PK":1,"title":"title","startat":"2024-02-18T21:54:42+02:00","duration":0,"description":"description","owner":"owner","notifyearly":0}}}
```

Read

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/1"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/events/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 09:58:34 GMT
< Content-Length: 190
< 
* Connection #0 to host localhost left intact
{"method":"api.events.get","error":"","data":{"item":{"PK":1,"title":"title","startat":"2024-02-18T21:54:42+02:00","duration":0,"description":"description","owner":"owner","notifyearly":0}}}
```

Create дополнительных объектов

```bash
curl -d "{\"title\":\"title 2\",\"startat\":\"2023-08-11T21:54:42+02:00\",\"duration\":0,\"description\": \"description 2\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -POST "http://localhost:8080/api/events/create"
```

```text
{"method":"api.events.create","error":"","data":{"item":{"PK":2,"title":"title 2","startat":"2023-08-11T21:54:42+02:00","duration":0,"description":"description 2","owner":"owner","notifyearly":0}}}
```

```bash
curl -d "{\"title\":\"title 3\",\"startat\":\"2023-08-11T21:54:42+02:00\",\"duration\":0,\"description\": \"description 3\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -POST "http://localhost:8080/api/events/create"
```

```text
{"method":"api.events.create","error":"","data":{"item":{"PK":3,"title":"title 3","startat":"2023-08-11T21:54:42+02:00","duration":0,"description":"description 3","owner":"owner","notifyearly":0}}}
```

Delete

```bash
curl -v -H "Content-Type: application/json" -X DELETE "http://localhost:8080/api/events/3/delete"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> DELETE /api/events/3/delete HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 10:01:43 GMT
< Content-Length: 197
< 
* Connection #0 to host localhost left intact
{"method":"api.events.delete","error":"","data":{"item":{"PK":3,"title":"title 3","startat":"2023-08-11T21:54:42+02:00","duration":0,"description":"description 3","owner":"owner","notifyearly":0}}}
```

Read не существующего

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/3"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/events/3 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 10:02:08 GMT
< Content-Length: 59
< 
* Connection #0 to host localhost left intact
{"method":"api.events.get","error":"","data":{"item":null}}
```

Попытка Delete через HTTP POST даст `"error":"invalid HTTP method"`:

```bash
curl -v -H "Content-Type: application/json" -X POST "http://localhost:8080/api/events/3/delete"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /api/events/3/delete HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> 
< HTTP/1.1 405 Method Not Allowed
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 10:03:01 GMT
< Content-Length: 72
< 
* Connection #0 to host localhost left intact
{"method":"api.events.delete","error":"invalid HTTP method","data":null}
```

Read списка

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/"
```

```text
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/events/ HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Content-Type: application/json
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sat, 12 Aug 2023 10:03:43 GMT
< Content-Length: 334
< 
* Connection #0 to host localhost left intact
{"method":"api.events.list","error":"","data":{"items":[{"PK":1,"title":"title","startat":"2024-02-18T21:54:42+02:00","duration":0,"description":"description","owner":"owner","notifyearly":0},{"PK":2,"title":"title 2","startat":"2023-08-11T21:54:42+02:00","duration":0,"description":"description 2","owner":"owner","notifyearly":0}]}}
```

### заметки GRPC

[GRPC](https://grpc.io/docs/languages/go/quickstart/)

> TODO: Надо в .profile:
>
> ```bash
> export PATH="$PATH:$(go env GOPATH)/bin"
> ```

### Иное

Важны названия веток для автотестов.

Конфлик линтеров.
