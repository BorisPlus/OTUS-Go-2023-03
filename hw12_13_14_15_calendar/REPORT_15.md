# Дз 15

## Докеризация нативная

```bash
make docker-run

docker ps
```

| CONTAINER ID |           IMAGE            |        COMMAND         |    CREATED     |     STATUS     |                                                                        PORTS                                                                        |       NAMES       |
|:------------:|:--------------------------:|:----------------------:|:--------------:|:--------------:|:---------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------:|
| d286833a2c1f |      hw15_sender:dev       | "/opt/service/servic…" | 38 seconds ago | Up 36  seconds |                                                                                                                                                     |  hw15_sender_dev  |
| 50d6854928ec |     hw15_archiver:dev      | "/opt/service/servic…" | 4 minutes ago  |  Up 4 minutes  |                                                                                                                                                     | hw15_archiver_dev |
| 6ae0a84cb891 |     hw15_sheduler:dev      | "/opt/service/servic…" | 7 minutes ago  |  Up 7 minutes  |                                                                                                                                                     | hw15_sheduler_dev |
| 3a246dacc0ff | rabbitmq:3.10.7-management | "docker-entrypoint.s…" | 14 minutes ago | Up 14 minutes  | 4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, :::5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp, :::15672->15672/tcp | hw15_rabbitmq_dev |
| 02b6eef11aa3 |     hw15_calendar:dev      | "/opt/service/servic…" | 14 minutes ago | Up 14 minutes  |                                0.0.0.0:5000->5000/tcp, :::5000->5000/tcp, 0.0.0.0:8888->8080/tcp, :::8888->8080/tcp                                 | hw15_calendar_dev |

Загрузим тестовые 30 записей:

```bash
go run ./cmd/dataset/ --config ./configs/dataset.yaml
```

Наблюдаем тестовые 10 отправленных:

```bash
docker logs hw15_sender_dev

Notice "PARIATUR IMPEDIT ID QUO SOLUTA" send to "pariatur@Wordpedia.name"
Notice "ASPERIORES CUM" send to "RogerBennett@Minyx.edu"
Notice "RERUM OMNIS QUIDEM MODI" send to "iWilliams@Feedfire.org"
Notice "TEMPORE OPTIO EA RERUM" send to "dolores_autem@Roomm.edu"
Notice "QUO EA FUGA" send to "molestiae@Quinu.name"
Notice "ET ACCUSANTIUM SOLUTA CONSECTETUR VERO" send to "xTucker@Vitz.edu"
Notice "NIHIL RERUM IPSAM" send to "placeat_esse_quasi@Blognation.biz"
Notice "EUM QUOS" send to "xGray@Midel.org"
Notice "EXPLICABO AUTEM IPSAM MODI QUI" send to "LouiseColeman@Centidel.info"
Notice "PLACEAT ILLO RECUSANDAE VOLUPTATE" send to "gRay@Katz.com"
```

Наблюдаем тестовые 10 принудительно заархивированных:

```bash
docker logs hw15_archiver_dev

INFO [2023-10-14 22:26:14] Transmitter.Start()
INFO [2023-10-14 22:26:14] Transmitter step.
INFO [2023-10-14 22:31:09] new candidate.
INFO [2023-10-14 22:31:09] Must be archived: 19
INFO [2023-10-14 22:31:09] new candidate.
INFO [2023-10-14 22:31:19] Transmitter step.
INFO [2023-10-14 22:31:20] Must be archived: 15
INFO [2023-10-14 22:31:20] Must be archived: 11
INFO [2023-10-14 22:31:20] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] Must be archived: 13
INFO [2023-10-14 22:31:23] Must be archived: 16
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] Must be archived: 14
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] Must be archived: 12
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] Must be archived: 20
INFO [2023-10-14 22:31:23] Must be archived: 18
INFO [2023-10-14 22:31:23] new candidate.
INFO [2023-10-14 22:31:23] Must be archived: 17
INFO [2023-10-14 22:31:33] Transmitter step.
INFO [2023-10-14 22:31:37] new candidate.
...
```

## Докеризация нативная