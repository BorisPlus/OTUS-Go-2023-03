# Дз 15

## Докеризация нативная

__Замечание:__ сборка проекта с вариантом взаимодействия с "In-memory database".

```bash
make docker-run
...

docker ps -a

CONTAINER ID   IMAGE                        COMMAND                  CREATED          STATUS                      PORTS                                                                                                         NAMES
dd5578c964a5   hw15_sender_dev:dev          "/opt/service/servic…"   29 minutes ago   Up 29 minutes                                                                                                                             hw15_sender_dev
092fd25f57fc   hw15_archiver_dev:dev        "/opt/service/servic…"   34 minutes ago   Up 34 minutes                                                                                                                             hw15_archiver_dev
6e48e883bd30   hw15_sheduler_dev:dev        "/opt/service/servic…"   37 minutes ago   Up 37 minutes                                                                                                                             hw15_sheduler_dev
d03e622fc3c3   hw15_calendar_dev:dev        "/opt/service/servic…"   40 minutes ago   Up 40 minutes               0.0.0.0:5080->5000/tcp, 0.0.0.0:8080->8000/tcp                                                                hw15_calendar_dev
ef04f1b3b1fd   hw15_rmqdevops_dev:dev       "/opt/service/servic…"   43 minutes ago   Exited (0) 43 minutes ago                                                                                                                 hw15_rmqdevops_dev
2636705fc5fa   rabbitmq:3.10.7-management   "docker-entrypoint.s…"   49 minutes ago   Up 49 minutes               4369/tcp, 5671/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:8672->5672/tcp, 0.0.0.0:18672->15672/tcp   hw15_rabbitmq_dev
```

Загрузим тестовые 30 записей:

```bash
go run ./cmd/dataset/ --config ./configs/dataset.standalone.yaml
```

Наблюдаем тестовые 10 из 30, которые априори должны быть незамедлительно отправлены:

```bash
docker logs -f hw15_sender_dev

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

Наблюдаем тестовые 10 из 30, которые априори должны быть заархивированы (см. "Must be archived"):

```bash
docker logs -f hw15_archiver_dev

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

__Замечание:__ сборка проекта с вариантом взаимодействия с "Postgres" в качестве базы данных.

__Замечание:__ задействуется механизм запуска контейнеров по критериям зависимости от "действительной" готовности к работе внутреконтейнерных сервисов.

```bash
make docker-compose-up
...

docker ps -a

CONTAINER ID   IMAGE                        COMMAND                  CREATED              STATUS                         PORTS                                                                                                         NAMES
1034e9e9eb75   hw15_archiver_prod           "/opt/service/servic…"   58 seconds ago       Up 54 seconds                                                                                                                                hw15_archiver_prod
3ea54591821c   hw15_sheduler_prod           "/opt/service/servic…"   58 seconds ago       Up 55 seconds                                                                                                                                hw15_sheduler_prod
536a97940792   hw15_sender_prod             "/opt/service/servic…"   58 seconds ago       Up 55 seconds                                                                                                                                hw15_sender_prod
3a0b3b72e933   hw15_rmqdevops_prod          "/opt/service/servic…"   About a minute ago   Exited (0) 59 seconds ago                                                                                                                    hw15_rmqdevops_prod
3297433dad16   hw15_calendar_prod           "/opt/service/servic…"   2 minutes ago        Up About a minute              0.0.0.0:5000->5000/tcp, 0.0.0.0:8888->8000/tcp                                                                hw15_calendar_prod
5c774758d5eb   hw15_goose_prod              "/bin/sh -c 'goose -…"   2 minutes ago        Exited (0) 2 minutes ago                                                                                                                     hw15_goose_prod
4a3be29f4cfc   rabbitmq:3.10.7-management   "docker-entrypoint.s…"   2 minutes ago        Up 2 minutes (healthy)                                                                                                                       hw15_rabbitmq_prod
efe90c4f6fb9   postgres:13.3                "docker-entrypoint.s…"   2 minutes ago        Up 2 minutes (healthy)                                                                                                                       hw15_postgres_prod
ac9486e35a62   dpage/pgadmin4               "/entrypoint.sh"         2 minutes ago        Up 2 minutes                   443/tcp, 0.0.0.0:8883->80/tcp                                                                                 hw15_pgadmin_prod
```

## Интеграционное тестирование

__Замечание:__  Поскольку контейнеры с миграциями (для `RabbitMQ` и `Postgres`) успешно завершаются (то есть "абортятся"), таким образом передавая очередь запуска другим контейнерам, вариант использования в `docker-compose` флага `--exit-code-from` (а значит и `--abort-on-container-exit`) не подходит.

__Замечание:__  В рамках переопределений файла `docker-compose.integration_test.yaml` происходит запуск ["чекера"](configs/checker.integration_test.yaml), имеющего полный доступ ко всем микросервисам, в том числе из внутреннего контура.

```bash
make integration_test
...

Creating hw15_rabbitmq_integration_test  ... done
Creating hw15_postgres_integration_test  ... done
Creating hw15_pgadmin_integration_test   ... done
Creating hw15_goose_integration_test     ... done
Creating hw15_calendar_integration_test  ... done
Creating hw15_rmqdevops_integration_test ... done
Creating hw15_sender_integration_test    ... done
Creating hw15_archiver_integration_test  ... done
Creating hw15_sheduler_integration_test  ... done
Creating hw15_checker_integration_test   ... done

docker logs -f hw15_checker_integration_test 
INFO [2023-10-22 17:39:12] OK. Get all notices of sended events
INFO [2023-10-22 17:40:13] OK. Get all notices of archived events (also after send).
INFO [2023-10-22 17:40:13] OK. All selected titles of created events are correct.
INFO [2023-10-22 17:40:13] OK. Get expected count of events: 30.
INFO [2023-10-22 17:40:13] Everything all right.
```

```bash
docker ps -a

CONTAINER ID   IMAGE                             COMMAND                  CREATED          STATUS                      PORTS                                                                                     NAMES
224b09718808   hw15_checker_integration_test     "/opt/service/servic…"   14 minutes ago   Exited (0) 13 minutes ago                                                                                             hw15_checker_integration_test
e51b05997d06   hw15_sheduler_integration_test    "/opt/service/servic…"   15 minutes ago   Up 15 minutes                                                                                                         hw15_sheduler_integration_test
acc858c82ce6   hw15_archiver_integration_test    "/opt/service/servic…"   15 minutes ago   Up 15 minutes                                                                                                         hw15_archiver_integration_test
01c96b7c5768   hw15_sender_integration_test      "/opt/service/servic…"   15 minutes ago   Up 15 minutes                                                                                                         hw15_sender_integration_test
0c50b4de47f5   hw15_rmqdevops_integration_test   "/opt/service/servic…"   15 minutes ago   Exited (0) 15 minutes ago                                                                                             hw15_rmqdevops_integration_test
6a593289e2f8   hw15_calendar_integration_test    "/opt/service/servic…"   15 minutes ago   Up 15 minutes (healthy)     0.0.0.0:5000->5000/tcp, 0.0.0.0:8888->8080/tcp                                            hw15_calendar_integration_test
a3b494527ef4   hw15_goose_integration_test       "/bin/sh -c 'goose -…"   15 minutes ago   Exited (0) 15 minutes ago                                                                                             hw15_goose_integration_test
8c5eaa955f88   postgres:13.3                     "docker-entrypoint.s…"   15 minutes ago   Up 15 minutes (healthy)     0.0.0.0:5432->5432/tcp                                                                    hw15_postgres_integration_test
8be13bc79822   dpage/pgadmin4                    "/entrypoint.sh"         15 minutes ago   Up 15 minutes               443/tcp, 0.0.0.0:8883->80/tcp                                                             hw15_pgadmin_integration_test
496ecd245908   rabbitmq:3.10.7-management        "docker-entrypoint.s…"   15 minutes ago   Up 15 minutes (healthy)     4369/tcp, 5671/tcp, 15671-15672/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:5672->5672/tcp   hw15_rabbitmq_integration_test
```

Видно, что `hw15_checker_integration_test` завершился с `Exited (0)`

## Для себя

* как посредством `docker-compose.integration_test.yaml` не добавлять порты, а переопределять.
