# Дз 14

## Особенности реализации

> Реализация _немного_ отличается от требуемой, но конструктивно, на мой взгляд, она нагляднее. Так, например, архивирование - это отдельный процесс, что логично с точки зрения одной ответственности за одним компонентом системы. А ахивируются именно просроченные уведомления, а не события в исходной БД (я вообще не люблю удалять данные из БД, заменяя в своих проектах это действие понятием "скрыть с глаз"). Дополнительно, в базу введен признак обработки события планировщиком, что позволяет останавливать исходное приложение и запускать его заново, без ущерба возникновения двойных уведомлений об одном и том же событии.

Процессы `рассылки`, `архивирования` и `отправки` уведомлений сведены к одному общему абстрактному "знаменателю" - `Перекладчик` (`transmitter`):

```go
package transmitter

type Item interface {
   any
}

...

type Streamer interface {
   Connect(context.Context) error
   Disconnect(context.Context) error
}

...

type Source[FROM Item] interface {
   Streamer
   DataChannel(context.Context) (<-chan FROM, error)
   Confirm(context.Context, *FROM) error
   Getback(context.Context, *FROM) error
}

...

type Target[TO Item] interface {
   Streamer
   Put(context.Context, *TO) error
}

...

type Transmitter[FROM Item, TO Item] struct {
   Source         Source[FROM]
   Target         Target[TO]
   Transmit       func(ctx context.Context, candidate FROM) (bool, error)
   Logger         interfaces.Logger
   LoopTimeoutSec int64
}

...

```

> Также стало необходимым разработать структуру дженерика математического множества. Она помогает отслеживать появление цикла обработки. Когда необработанный элемент поступает опять на вход обработчику внутри одного цикла. Но это имплементировано и требует реализации только функции хеша от сохраняемого во множество элемента.

Так:

* Планировщик `sheduler` "перекладывает" события из хранилища (через приложение по RPC) в очередь сообщений "RabbitMQ", переконвертируя объекты `событие` в объекты `уведомление`. При этом фиксируется факт того, что событие было забрано, и при слудующем цикле выборки повтороного его помещения в очередь не происходит.
* Архиватор `archive` архивирует просроченные по времени (`event.StartAt.Before(time.Now())`) рассылки уведомления, "перекладывая" их из одной очереди "q_noticies_sheduled" в другую "q_noticies_archived".
* Рассылщик `sender` отбирает из очереди сообщений актуальные для отправки уведомления (_если текущее время попало в период между началом уведомления о событии и началом самого события_ - `event.StartAt.Add(-time.Duration(event.Duration)*time.Second).Before(now) && now.Before(event.StartAt)`) и отправляет их, перекладывая клонированные данные в две очереди: "q_noticies_send" и "q_noticies_archived" (все происходит через один exchange, но бинд-тег `send` ведет в две очереди). Он содержит отсылку на интерфейс уведомляющего/рассылщика `type Notifier interface { Notify(models.Notice) error }`. В конечной реализации `Notify` - это вывод в STDOUT.

Остается реализовать в терминах:

* источника `Source[FROM Item]` - метод `DataChannel(context.Context) (<-chan FROM, error)` получения данных (из базы данных - это список событий в канал, а для RabbitMQ - это делает консьюмер);
* точки назначения `Target[TO Item]` - метод `Put(context.Context, *TO) error`;
* объявить функцию метода перекладчика `Transmit func(ctx context.Context, candidate FROM) (bool, error)` требуемой бизнес-логики в зависимости от источника и точки назначения.

Для чего именно сделан перекладчик:

* Для главного посыла ДЗ: "Процессы не должны зависеть от конкретной реализации RMQ-клиента";
* Тренировки работы с RabbitMQ на Golang: взяв событие на принятие решения "отправлять"/"архивировать" или "оставить в очереди", его нужно не просто вернуть, а вернуть **в начало** очереди.
* Тренировка работы с интерфейсами, каналами и дженериками в Golang.

## Процессы

```bash
make build
```

### (Пере-)настройка RabbitMQ очередей

Вынесена в отдельный модуль - аля DevOps. Ответственность на подъем вынтреннего устройства RabbitMQ в итоге сосредоточна в атомарном модуле:

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/management/rmqdevops/ --with-drop --config ./configs/rmqdevops.yaml
```

или после `make build`

```bash
cd hw12_13_14_15_calendar 
./bin/mgm_rmqdevops.goc --with-drop --config ./configs/rmqdevops.yaml
```

Лог:

```text
2023/08/29 23:54:47 Connect to "amqp://guest:guest@localhost:5672/".
2023/08/29 23:54:47 Channel open.
2023/08/29 23:54:47 withDrop true
2023/08/29 23:54:47 OK. Drop exchange: "notice_exchange".
2023/08/29 23:54:47 OK. Drop queue: "q_unspecified".
2023/08/29 23:54:47 OK. Drop queue: "q_noticies_sheduled".
2023/08/29 23:54:47 OK. Drop queue: "q_noticies_archived".
2023/08/29 23:54:47 OK. Drop queue: "q_noticies_archived".
2023/08/29 23:54:47 OK. Drop queue: "q_noticies_send".
2023/08/29 23:54:47 Declare exchange: "notice_exchange".
2023/08/29 23:54:47 Declare queue: "q_unspecified".
2023/08/29 23:54:47 Bind queue: "q_unspecified" to exchange "notice_exchange" with routing key "".
2023/08/29 23:54:47 Declare queue: "q_noticies_sheduled".
2023/08/29 23:54:47 Bind queue: "q_noticies_sheduled" to exchange "notice_exchange" with routing key "shedule".
2023/08/29 23:54:47 Declare queue: "q_noticies_archived".
2023/08/29 23:54:47 Bind queue: "q_noticies_archived" to exchange "notice_exchange" with routing key "archive".
2023/08/29 23:54:47 Declare queue: "q_noticies_archived".
2023/08/29 23:54:47 Bind queue: "q_noticies_archived" to exchange "notice_exchange" with routing key "send".
2023/08/29 23:54:47 Declare queue: "q_noticies_send".
2023/08/29 23:54:47 Bind queue: "q_noticies_send" to exchange "notice_exchange" with routing key "send".
2023/08/29 23:54:47 RabbitMQ (re-)configure done.
2023/08/29 23:54:47 Channel close.
2023/08/29 23:54:47 Connection close.
```

![exchange.png](REPORT_14.files/exchange.png)

![empty.png](REPORT_14.files/empty.png)

### Кадендарь

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/calendar/ --config ./configs/calendar.yaml
```

или после `make build`

```bash
cd hw12_13_14_15_calendar 
./bin/calendar/ --config ./configs/calendar.yaml
```

Лог

```text
INFO [2023-08-29 20:59:09] HTTP Config - {Host:localhost Port:8080 ReadTimeout:5s ReadHeaderTimeout:5s WriteTimeout:5s MaxHeaderBytes:1048576}
INFO [2023-08-29 20:59:09] RPC Config - {Host:localhost Port:5000}
INFO [2023-08-29 20:59:09] calendar is running...
INFO [2023-08-29 20:59:09] HTTPServer.Start()
INFO [2023-08-29 20:59:09] GRPCServer.Start()
^C                                                  <----------- В случае прерывания
INFO [2023-08-29 20:59:11] Complex Shutting down was done gracefully by signal.
INFO [2023-08-29 20:59:11] HTTPServer.Stop()
INFO [2023-08-29 20:59:11] GRPCServer.GracefulStop()
```

### Наполнение тестовыми данными

Реализует наполнение хранилища исходного приложения тестовыми данными:

* 10 кандидатов для архивирования.
* 10 на первоочередную отправку.
* 10 на отложенную отправку (через некоторое время).

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/management/dataset/ --config ./configs/calendar.yaml
```

или после `make build`

```bash
cd hw12_13_14_15_calendar 
./bin/mgm_dataset.goc --with-drop --config ./configs/calendar.yaml
```

Лог

```text
2023/08/30 01:01:51 Put event: {PK:0 Title:AD ISTE VELIT ARCHITECTO NESCIUNT StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:id accusantium ut placeat voluptas quia consectetur aliquam. Owner:accusamus@Eimbee.mil NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:ENIM EIUS SUSCIPIT ALIQUID VOLUPTATEM StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:nihil perspiciatis quia ea ut eos! Owner:DonaldHawkins@Npath.com NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:QUAS EST NATUS UNDE StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:ut fugiat est dolorem id harum est pariatur doloremque dolores non aut. Owner:BeverlyDuncan@Realmix.net NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:PERFERENDIS ET CUM StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:consequatur placeat officiis totam. Owner:aStanley@Yacero.gov NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:QUO SED LABORUM ADIPISCI StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:eum adipisci dolor qui recusandae ratione id. Owner:earum_rem_quod@Wikizz.gov NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:UT RERUM ELIGENDI VOLUPTATEM StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:libero et quia magni ducimus velit. Owner:omnis@Skipstorm.net NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:RECUSANDAE ATQUE CORRUPTI StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:recusandae quasi porro libero. Owner:TimothyCook@Meeveo.org NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:SED MAGNAM DOLORIBUS StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:odio omnis dolorum eveniet quis placeat soluta sapiente dolor modi. Owner:impedit@Gigabox.biz NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:SINT CORRUPTI StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:enim et magnam et. Owner:minima@Thoughtbridge.biz NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:NOBIS IMPEDIT StartAt:2023-08-30 01:51:51.414568391 +0300 MSK m=+3000.002230691 Duration:1800 Description:provident voluptates qui officiis molestias aspernatur fuga dolore eum. Owner:PhyllisTorres@Twinder.net NotifyEarly:3600 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:NEMO ET DOLOR StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:sed maxime ut facere ducimus et. Owner:odio_voluptatem@Skinder.info NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:OFFICIIS FUGIT EIUS QUI MAGNAM StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:cum et ut odio ut quis cumque. Owner:oStanley@Wordware.org NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:LABORUM EOS AUT StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:veniam consequatur ipsum nisi delectus incidunt. Owner:delectus@Twiyo.info NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:EST NEMO VOLUPTAS SIMILIQUE VENIAM StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:qui repudiandae numquam sint amet saepe laudantium quia molestias. Owner:corrupti@Tekfly.info NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:QUIBUSDAM ACCUSANTIUM StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:maxime perspiciatis quasi ipsam nulla voluptas sit ullam. Owner:voluptatem@Rhycero.mil NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:ALIQUID AUT ADIPISCI StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:qui architecto ea. Owner:repudiandae_unde_est@Edgeclub.org NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:NON EST SUNT StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:aperiam ea ad amet aut praesentium aspernatur in non. Owner:pJacobs@Devpoint.info NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:FACERE FUGA RECUSANDAE QUIA StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:reprehenderit aliquam ut nobis consequuntur molestias dignissimos ut natus. Owner:necessitatibus@Twitterworks.name NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:REPELLAT MAXIME AUT StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:et aut aperiam voluptatibus et veniam et. Owner:RyanAllen@Ainyx.edu NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:REPREHENDERIT QUI VOLUPTATE StartAt:2023-08-17 13:01:51.414568391 +0300 MSK m=-1079999.997769309 Duration:1800 Description:nostrum fugit inventore dolor quis sed velit molestias. Owner:MarthaLawrence@Jaloo.biz NotifyEarly:60 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:QUI NEQUE APERIAM UT SAPIENTE StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:voluptate quo consequatur quaerat rem rerum. Owner:dolor_quisquam@Babblestorm.com NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:ASPERIORES AUT PROVIDENT ESSE SIMILIQUE StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:et aliquid odio nam dolores modi. Owner:numquam_aut_et@Lazz.net NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:REPELLENDUS NIHIL DOLORUM StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:quo architecto aut quo molestias quia consequatur nihil enim omnis esse itaque. Owner:gDixon@Edgeify.edu NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:QUI APERIAM ET StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:aut debitis pariatur. Owner:JoseWells@Aimbo.biz NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:DISTINCTIO SED SUNT VOLUPTATEM StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:officia omnis autem dicta. Owner:iure_tempora@Mynte.mil NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:QUASI EUM StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:molestias et quia a vel. Owner:quaerat_qui_est@Trudeo.edu NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:EVENIET PORRO LIBERO VOLUPTATEM ALIQUAM StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:non id fugit repellendus velit. Owner:impedit_cum_neque@Fiveclub.mil NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:IPSA ENIM CORPORIS NUMQUAM StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:modi commodi et. Owner:StephanieHarrison@Skinte.name NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:ET MINIMA REPREHENDERIT TEMPORIBUS StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:iure unde eveniet quo laborum aliquam omnis in doloribus. Owner:MelissaDay@Mydo.info NotifyEarly:1 Sheduled:false}
2023/08/30 01:01:51 Put event: {PK:0 Title:EST DOLOREM MOLESTIAE VEL ET StartAt:2023-08-30 11:01:51.414568391 +0300 MSK m=+36000.002230691 Duration:1800 Description:nulla quis id deserunt veniam possimus quaerat. Owner:pAndrews@Dablist.org NotifyEarly:1 Sheduled:false}
```

### Планировщик

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/sheduler/ --config ./configs/sheduler.yaml
```

или после `make build`

```bash
cd hw12_13_14_15_calendar 
./bin/sheduler.goc --config ./configs/sheduler.yaml
```

```text
INFO [2023-08-29 21:49:23] Sheduler start.
INFO [2023-08-29 21:49:23] Transmitter.Start()
^C                                                  <----------- В случае прерывания
INFO [2023-08-29 21:49:28] Transmitter.Stop()
INFO [2023-08-29 21:49:28] Sheduler end.
```

![sheduled.png](REPORT_14.files/sheduled.png)

### Архиватор

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/archiver/ --config ./configs/archiver.yaml
```

или после `make build`

```bash
cd hw12_13_14_15_calendar 
./bin/archiver.goc --config ./configs/archiver.yaml
```

Лог:

```text
INFO [2023-08-29 21:30:47] Transmitter.Start()
INFO [2023-08-29 21:30:47] Must be archived: 15
INFO [2023-08-29 21:30:47] Must be archived: 20
INFO [2023-08-29 21:30:47] Must be archived: 12
INFO [2023-08-29 21:30:47] Must be archived: 11
INFO [2023-08-29 21:30:47] Must be archived: 17
INFO [2023-08-29 21:30:47] Must be archived: 14
INFO [2023-08-29 21:30:47] Must be archived: 16
INFO [2023-08-29 21:30:47] Must be archived: 13
INFO [2023-08-29 21:30:47] Must be archived: 19
INFO [2023-08-29 21:30:47] Must be archived: 18
^C                                                  <----------- В случае прерывания
INFO [2023-08-29 21:35:55] Transmitter.Stop()
```

![archived.png](REPORT_14.files/archived.png)

### Рассылщик

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/sender/ --config ./configs/sender.yaml
```

или после `make build`

```bash
cd hw12_13_14_15_calendar 
./bin/sender.goc --config ./configs/sender.yaml
```

Лог (по мере добавления уведомлений, например, посредством `mgm_dataset.goc`, будут видны новые отправляемые сообщения):

```text
Notice "AD ISTE VELIT ARCHITECTO NESCIUNT" send to "accusamus@Eimbee.mil"
Notice "UT RERUM ELIGENDI VOLUPTATEM" send to "omnis@Skipstorm.net"
Notice "ENIM EIUS SUSCIPIT ALIQUID VOLUPTATEM" send to "DonaldHawkins@Npath.com"
Notice "QUO SED LABORUM ADIPISCI" send to "earum_rem_quod@Wikizz.gov"
Notice "RECUSANDAE ATQUE CORRUPTI" send to "TimothyCook@Meeveo.org"
Notice "SINT CORRUPTI" send to "minima@Thoughtbridge.biz"
Notice "SED MAGNAM DOLORIBUS" send to "impedit@Gigabox.biz"
Notice "NOBIS IMPEDIT" send to "PhyllisTorres@Twinder.net"
Notice "PERFERENDIS ET CUM" send to "aStanley@Yacero.gov"
Notice "QUAS EST NATUS UNDE" send to "BeverlyDuncan@Realmix.net"
```

![send.png](REPORT_14.files/send.png)

## Заметки (для себя)

[go get github.com/streadway/amqp](https://nuancesprog.ru/p/4907/)

* 5671:5671
* 5672:5672
* 15671:15671
* 15672:15672
* 15691:15691
* 15692:15692
* 25672:25672
* 4369:4369

Тем, кто использует Portainer, в rabbitmq-образе надо дополнительно:

* прокинуть `15672:15672`
* включить

```bash
rabbitmq-plugins enable rabbitmq_management
```

и

```bash
cd  /etc/rabbitmq/conf.d/
echo management_agent.disable_metrics_collector = false > management_agent.disable_metrics_collector.conf
```

* To move messages, the shovel plugin must be enabled, try:

```bash
rabbitmq-plugins enable rabbitmq_shovel rabbitmq_shovel_management

wget -P /opt/rabbitmq/plugins/ https://github.com/noxdafox/rabbitmq-message-deduplication/releases/download/0.6.1/elixir-1.13.4.ez

wget -P /opt/rabbitmq/plugins/ https://github.com/noxdafox/rabbitmq-message-deduplication/releases/download/0.6.1/rabbitmq_message_deduplication-0.6.1.ez

rabbitmq-plugins enable rabbitmq_message_deduplication
```
