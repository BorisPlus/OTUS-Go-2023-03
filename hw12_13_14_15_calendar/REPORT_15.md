# Дз 15

## Докеризация нативная

```bash
make docker-run

docker ps
```

CONTAINER ID IMAGE COMMANDCREATEDSTATUSPORTS NAMES
f630e6aadf25 hw15_sender:dev "/opt/service/servic…" 2 hours agoUp 2 hourshw15_sender_dev
01dd1e67541f hw15_archiver:dev "/opt/service/servic…" 2 hours agoUp 2 hourshw15_archiver_dev
3f8b88fa0298 hw15_sheduler:dev "/opt/service/servic…" 2 hours agoUp 2 hourshw15_sheduler_dev
63e4e01f8efe hw15_calendar:dev "/opt/service/servic…" 2 hours agoUp 2 hours0.0.0.0:5000->5000/tcp, 0.0.0.0:8888->8080/tcphw15_calendar_dev
1b7a388dbfdd rabbitmq:3.10.7-management"docker-entrypoint.s…" 2 hours agoUp 2 hours4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp hw15_rabbitmq_dev
3a51d5545067 phpmyadmin"/docker-entrypoint.…" 5 weeks agoUp 4 hours0.0.0.0:8070->80/tcp, :::8070->80/tcp phpmyadmin-phpmyadmin-1
35ab29c88c94 mariadb:10.6"docker-entrypoint.s…" 5 weeks agoUp 4 hours3306/tcpphpmyadmin-db-1
c64363606b2b mariadb:latest"docker-entrypoint.s…" 5 weeks agoUp 4 hours0.0.0.0:32769->3306/tcp, :::32769->3306/tcp charming_matsumoto
032da9337391 mongo:latest"docker-entrypoint.s…" 7 weeks agoUp 10 minutes 0.0.0.0:32916->27017/tcp, :::32916->27017/tcp opencellid
7b65587539ae postgres:latest "docker-entrypoint.s…" 2 months ago Up 4 hours0.0.0.0:15432->5432/tcp, :::15432->5432/tcp postgres
fa97d8c9fe52 dpage/pgadmin4:latest "/entrypoint.sh" 2 months ago Up 4 hours443/tcp, 0.0.0.0:18080->80/tcp, :::18080->80/tcppgadmin4
46bc55affcb0 portainer/portainer-ce:latest "/portainer" 2 months ago Up 4 hours0.0.0.0:8000->8000/tcp, :::8000->8000/tcp, 0.0.0.0:9443->9443/tcp, :::9443->9443/tcp, 9000/tcpportainer
Загрузим тестовые 30 записей:

```bash
go run ./cmd/dataset/ --config ./configs/dataset.standalone.yaml
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

Наблюдаем тестовые 10, которые априори должны были быть заархивированы:

```bash
docker logs hw15_archiver_dev
INFO [2023-10-18 12:26:58] Transmitter.Start()
INFO [2023-10-18 12:26:58] Transmitter step.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] Must be archived: 14
INFO [2023-10-18 14:39:42] Must be archived: 18
INFO [2023-10-18 14:39:42] Must be archived: 11
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] Must be archived: 12
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] Must be archived: 15
INFO [2023-10-18 14:39:42] Must be archived: 16
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] Must be archived: 17
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:42] Must be archived: 13
INFO [2023-10-18 14:39:42] new candidate.
INFO [2023-10-18 14:39:46] new candidate.
INFO [2023-10-18 14:39:56] Transmitter step.
INFO [2023-10-18 14:39:58] new candidate.
INFO [2023-10-18 14:40:08] Transmitter step.
INFO [2023-10-18 14:40:08] new candidate.
INFO [2023-10-18 14:40:08] Must be archived: 20
INFO [2023-10-18 14:40:08] new candidate.
INFO [2023-10-18 14:40:18] Transmitter step.
INFO [2023-10-18 14:40:18] new candidate.
INFO [2023-10-18 14:40:18] new candidate.
INFO [2023-10-18 14:40:28] Transmitter step.
INFO [2023-10-18 14:40:28] new candidate.
INFO [2023-10-18 14:40:28] new candidate.
INFO [2023-10-18 14:40:28] new candidate.
INFO [2023-10-18 14:40:28] new candidate.
INFO [2023-10-18 14:40:28] new candidate.
INFO [2023-10-18 14:40:38] Transmitter step.
INFO [2023-10-18 14:40:38] new candidate.
INFO [2023-10-18 14:40:38] new candidate.
INFO [2023-10-18 14:40:38] new candidate.
INFO [2023-10-18 14:40:38] Must be archived: 19
INFO [2023-10-18 14:40:38] new candidate.
INFO [2023-10-18 14:40:41] new candidate.
INFO [2023-10-18 14:40:41] new candidate.
INFO [2023-10-18 14:40:41] new candidate.
INFO [2023-10-18 14:40:41] new candidate.
INFO [2023-10-18 14:40:41] new candidate.
INFO [2023-10-18 14:40:51] Transmitter step.
INFO [2023-10-18 14:40:51] new candidate.
INFO [2023-10-18 14:40:51] new candidate.
INFO [2023-10-18 14:40:51] new candidate.
INFO [2023-10-18 14:40:51] new candidate.
...

```

## Докеризация компоуз

```bash
make docker-compose-up
```
