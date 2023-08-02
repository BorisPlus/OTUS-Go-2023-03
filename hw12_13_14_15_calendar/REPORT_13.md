# Дз 13

```bash
go run ./hw12_13_14_15_calendar/cmd/calendar/ --config=./hw12_13_14_15_calendar/configs/config.yaml

curl -d "{\"title\":\"title\",\"duration\":0,\"description\":\"description\":\"owner\":\"owner\",\"notifyearly\":\"notifyearly\",\"startat\":\"2023-01-01 00:00:00\"}" -H "Content-Type: application/json" -X POST "http://localhost:8080/api/events/"

curl -v -d "{\"title\":\"title\",\"startat\":\"2021-02-18T21:54:42+02:00\",\"duration\":0,\"description\": \"description\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -X POST "http://localhost:8080/api/events/create/"

curl -X POST -H 'Content-Type: application/json' -d "{\"title\": \"title\"}" "http://localhost:8080/api/events/create/"
curl -X POST -H 'Content-Type: application/json' -d "{}" "http://localhost:8080/api/events/create/"
```

```text
-v - для verbose
```

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
curl -d "{\"PK\":2,\"title\":\"title\",\"startat\":\"2023-02-18T21:54:42+02:00\",\"duration\":0,\"description\": \"description\",\"owner\":\"owner\",\"notifyearly\":0}" -H "Content-Type: application/json" -X PUTCH "http://localhost:8080/api/events/1/update"
```

Read

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/1"
```

Delete

```bash
curl -v -H "Content-Type: application/json" -X PUT "http://localhost:8080/api/events/2/delete"
```

Read

```bash
curl -v -H "Content-Type: application/json" "http://localhost:8080/api/events/1"
```


```go

import (
	"encoding/json"
	"fmt"
)

type E struct {
	Error     error  `json:"error"`
}

fmt.Println(json.Marshal(E{fmt.Errorf("qweqw")}))
