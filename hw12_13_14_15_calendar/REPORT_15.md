# Дз 15

## Докеризация нативная

```bash
make docker-run

docker ps
```

```text
CONTAINER ID   IMAGE                        COMMAND                  CREATED          STATUS          PORTS                                                                                                         NAMES
046df462da34   hw15_sender:dev              "/opt/service/servic…"   4 minutes ago    Up 4 minutes                                                                                                                  hw15_sender_dev
c69dc78f093f   hw15_archiver:dev            "/opt/service/servic…"   8 minutes ago    Up 8 minutes                                                                                                                  hw15_archiver_dev
55e2e05aac39   hw15_sheduler:dev            "/opt/service/servic…"   12 minutes ago   Up 12 minutes                                                                                                                 hw15_sheduler_dev
de2b7a5c7397   hw15_calendar:dev            "/opt/service/servic…"   17 minutes ago   Up 17 minutes   0.0.0.0:5000->5000/tcp, 0.0.0.0:8888->8080/tcp                                                                hw15_calendar_dev
f685a7a07d01   rabbitmq:3.10.7-management   "docker-entrypoint.s…"   24 minutes ago   Up 24 minutes   4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp   hw15_rabbitmq_dev
```

Загрузим тестовые 30 записей:

```bash
go run ./cmd/dataset/ --config ./configs/dataset.standalone.yaml
```

Наблюдаем тестовые 10 из 30, которые априори должны быть незамедлительно отправлены:

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

Наблюдаем тестовые 10 из 30, которые априори должны быть заархивированы:

```bash
docker logs hw15_archiver_dev

INFO [2023-10-18 13:26:58] Transmitter.Start()
INFO [2023-10-18 13:26:58] Transmitter step.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] Must be archived: 14
INFO [2023-10-18 13:39:42] Must be archived: 18
INFO [2023-10-18 13:39:42] Must be archived: 11
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] Must be archived: 12
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] Must be archived: 15
INFO [2023-10-18 13:39:42] Must be archived: 16
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] Must be archived: 17
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:42] Must be archived: 13
INFO [2023-10-18 13:39:42] new candidate.
INFO [2023-10-18 13:39:46] new candidate.
INFO [2023-10-18 13:39:56] Transmitter step.
INFO [2023-10-18 13:39:58] new candidate.
INFO [2023-10-18 13:40:08] Transmitter step.
INFO [2023-10-18 13:40:08] new candidate.
INFO [2023-10-18 13:40:08] Must be archived: 20
INFO [2023-10-18 13:40:08] new candidate.
INFO [2023-10-18 13:40:18] Transmitter step.
INFO [2023-10-18 13:40:18] new candidate.
INFO [2023-10-18 13:40:18] new candidate.
INFO [2023-10-18 13:40:28] Transmitter step.
INFO [2023-10-18 13:40:28] new candidate.
INFO [2023-10-18 13:40:28] new candidate.
INFO [2023-10-18 13:40:28] new candidate.
INFO [2023-10-18 13:40:28] new candidate.
INFO [2023-10-18 13:40:28] new candidate.
INFO [2023-10-18 13:40:38] Transmitter step.
INFO [2023-10-18 13:40:38] new candidate.
INFO [2023-10-18 13:40:38] new candidate.
INFO [2023-10-18 13:40:38] new candidate.
INFO [2023-10-18 13:40:38] Must be archived: 19
INFO [2023-10-18 13:40:38] new candidate.
INFO [2023-10-18 13:40:41] new candidate.
INFO [2023-10-18 13:40:41] new candidate.
INFO [2023-10-18 13:40:41] new candidate.
INFO [2023-10-18 13:40:41] new candidate.
INFO [2023-10-18 13:40:41] new candidate.
INFO [2023-10-18 13:40:51] Transmitter step.
INFO [2023-10-18 13:40:51] new candidate.
INFO [2023-10-18 13:40:51] new candidate.
INFO [2023-10-18 13:40:51] new candidate.
INFO [2023-10-18 13:40:51] new candidate.
...
```

## Докеризация компоуз

```bash
make docker-compose-up
```

```bash
docker ps 
```

```text
CONTAINER ID   IMAGE                        COMMAND                  CREATED         STATUS                          PORTS                                                                                                         NAMES
7cc7533e420c   docker_sheduler              "/opt/service/servic…"   2 minutes ago   Restarting (0) 27 seconds ago                                                                                                                 hw15_sheduler_prod
557f35cd8597   hw15_calendar_prod           "/opt/service/servic…"   2 minutes ago   Up 2 minutes (healthy)          0.0.0.0:5000->5000/tcp, 0.0.0.0:8888->8080/tcp                                                                hw15_calendar_prod
a71652b1ebd4   docker_rmq_migration         "/opt/service/servic…"   2 minutes ago   Restarting (0) 56 seconds ago                                                                                                                 hw15_rmqdevops_prod
72bd605a195a   rabbitmq:3.10.7-management   "docker-entrypoint.s…"   2 minutes ago   Up 2 minutes (healthy)          4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp   hw15_rabbitmq_prod
a1719aa5a0c8   postgres:13.3                "docker-entrypoint.s…"   2 minutes ago   Up 2 minutes (healthy)                                                                                                                        hw15_postgres_prod
8052150783cd   dpage/pgadmin4               "/entrypoint.sh"         2 minutes ago   Up 2 minutes                    443/tcp, 0.0.0.0:8883->80/tcp                                                                                 hw15_pgadmin_prod

```