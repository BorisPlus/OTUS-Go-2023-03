# Дз 14

## Особенности реализации

1. Отправщик рассылщика представлен в виде [интерфейса](internal/interfaces/transmitter.go), его реализация в качестве [перекладчика](internal/rmq/notifier.go) из одной очреди в другую RabbitMQ-одчередь.
2. [Архивируется](cmd/archiver/archiver.go) то, что просрочено по уведомлениям.

```go
if event.StartAt.Before(time.Now())
```

3. [Рассылка](cmd/sender/sender.go) производится, если текущее время попало в период между началом уведомления и началом события, при этом тут же отправляется в архив.

```go
if event.StartAt.Add(-time.Duration(event.Duration)*time.Second).Before(now) && now.Before(event.StartAt)
```

## Процессы

Порядок выполнения ДЗ:

* установить локально очередь сообщений RabbitMQ (или сразу через Docker, если знаете как);
* создать процесс Планировщик (`scheduler`), который периодически сканирует основную базу данных,
выбирая события о которых нужно напомнить:
  * при запуске процесс должен подключаться к RabbitMQ и создавать все необходимые структуры
    (топики и пр.) в ней;
  * процесс должен выбирать сообытия для которых следует отправить уведомление (у события есть соотв. поле),
    создавать для каждого Уведомление (описание сущности см. в [ТЗ](./CALENDAR.MD)),
    сериализовать его (например, в JSON) и складывать в очередь;
  * процесс должен очищать старые (произошедшие более 1 года назад) события.
* создать процесс Рассыльщик (`sender`), который читает сообщения из очереди и шлёт уведомления;
непосредственно отправку делать не нужно - достаточно логировать сообщения / выводить в STDOUT.
* настройки подключения к очереди, периодичность запуска и пр. настройки процессов вынести в конфиг проекта;
* работу с кроликом вынести в отдельный пакет, который будут использовать пакеты, реализующие процессы выше.

Процессы не должны зависеть от конкретной реализации RMQ-клиента.

В результате компиляции проекта (`make build`) должно получаться 3 отдельных исполняемых файла
(по одному на микросервис):

* API (`calendar`);
* Планировщик (`calendar_scheduler`);
* Рассыльщик (`calendar_sender`).

Каждый из сервисов должен принимать путь файлу конфигурации:

```bash
./calendar           --config=/path/to/calendar_config.yaml
./calendar_scheduler --config=/path/to/scheduler_config.yaml
./calendar_sender    --config=/path/to/sender_config.yaml
```

### Кадендарь

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/calendar/ --config ./configs/calendar.yaml
```

### Наполнение тестовыми данными

* 10 кандидатов для архивирования.
* 10 на первоочередную отправку.
* 10 на отложенную отправку (через некоторое время).

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/dataset/ --config ./configs/calendar.yaml
```

### (Пере-)настройка очередей

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/rmqdevops/ --with-drop --config ./configs/rmqdevops.yaml
```

### Планировщик

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/sheduler/ --config ./configs/sheduler.yaml
```

### Архиватор

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/archiver/ --config ./configs/archiver.yaml
```

### Рассылщик

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/sender/ --config ./configs/sender.yaml
```

## Заметки для себя

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
