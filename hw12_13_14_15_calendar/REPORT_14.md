# Дз 14

## Особенности реализации

1. Отправщик рассылщика представлен в виде [интерфейса](internal/interfaces/transmitter.go), его реализация в качестве [перекладчика](internal/rmq/notifier.go) в RabbitMQ-одчередь.
2. [Архивируется](cmd/archiver/archiver.go) то, что не успело по уведомлениям.

```go
if event.StartAt.Before(time.Now())
```

3. [Рассылка](cmd/sender/sender.go) производится, если текущее время попало в период между началом уведомления и началом события, при этом тут же отправляется в архив.

```go
if event.StartAt.Add(-time.Duration(event.Duration)*time.Second).Before(now) && now.Before(event.StartAt)
```

## Процессы

### Кадендарь

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/calendar/ --config ./configs/calendar.yaml
```

### Наполнение тестовыми данными

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/dataset/ --config ./configs/calendar.yaml
```

### (Пере-)настройка очередей

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/rmqdevops/ drop --config ./configs/rmqdevops.yaml
```

```bash
cd hw12_13_14_15_calendar 
go run ./cmd/rmqdevops/ --config ./configs/rmqdevops.yaml
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
