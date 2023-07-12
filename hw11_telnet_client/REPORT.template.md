# Домашнее задание №11 «Клиент TELNET»

Описание [задания](./README.md). Я переработал представление диалогов взаимодействия клиента и сервера в виде таблицы ([вот](https://github.com/BorisPlus/OTUS-Go-2023-03/blob/80c7b6f09f5b8c79f554c9652fb397937f4142bc/hw11_telnet_client/README.md#L29) и [вот](https://github.com/BorisPlus/OTUS-Go-2023-03/blob/80c7b6f09f5b8c79f554c9652fb397937f4142bc/hw11_telnet_client/README.md#L80)). Кажется, стало немного нагляднее.

## Реализация

Боль.

## Тестирование

### Тестирование интерфейса

```bash
go test -v telnet.go telnet_test.go > telnet_test.go.txt
```

```text
{{ telnet_test.go.txt }}
```

### Shell

Был вынужден внести правки в `test.sh`, так как в `nc` версии `[v1.10-47]` под `Debian` необходимы дополнительные параметры `-s` и `-p` для запуска локального TCP-сервера: `nc -ls localhost -p 4242`. Наглядно видно на примере:

<details>
<summary>test.native.sh</summary>

```bash
{{ test.native.sh }}
```

</details>

> Для сохранения лога запускать как:
>
> ```bash
> ./test.native.sh >/dev/null 2>test.native.sh.out
>

```text
{{ test.native.sh.out }}
```

С учетом правок:

```bash
./test.sh
```

> Для сохранения лога запускать как:
>
> ```bash
> ./test.sh >/dev/null 2>test.sh.out
>

```text
{{ test.sh.out }}
```

### Тестирование аргументов

> Частично задействуется чужая кодовая база в отношении проверки соответствия наименования хоста валидному доменному имени. Функциями стандартных модулей пакета `net` (`net.url.Parse`, `net.SplitHostPort`) этого не добился, они уверенно допускают невалидные домены.

```bash
go test -v main_test.go main.go foreign_code_base.go telnet.go > main_test.go.txt
```

```text
{{ main_test.go.txt }}
```

### Интерактивное тестирование

Рассматриваются варианты, описанные в задании.

#### 1) Сервер обрывает соединение

Затруднительным оказалось реализовать вариант `1) Сервер обрывает соединение`. Сервер закрывает соединение без возможности отправки сообщения в уже закрытый канал связи клиенту, уведомлений `...Connection was closed by peer` клиент не получает. Но продемонстрирую вариант с иным способом прерывания диалога со стороны сервера, используя

```bash
nc -vvl -s localhost -p 4242 -c '
    set -x
    sleep 2
    echo "Hello from NC."
    sleep 2
    echo "I am Artificial intelligence."
    sleep 2
    echo "It''s joke. How are you?"
    sleep 10
    echo "Oh, Goodbye, client!"
    sleep 2
'
```

![./REPORT.files/1.gif]()

> Во этом варианте есть какая-то особенность. Если не вводить сообщения во время диалога (например "Fine"), то клиент "ругается" только на второе свое сообщение (необходим второй перенос). Можно просто два Enter подряд поставить после прощания сервера, и только тогда клиент просигнализирует о потере связи с сервером.

#### 2) Клиент завершает ввод

Вариант `2) Клиент завершает ввод` продемонстрирован ниже посредством ручного ввода сообщений попеременно от имени сервера и клиента

![./REPORT.files/2.gif]()

### Эффект горутин

Приходящие сообщения сервера "вставляются" в сообщения клиента. Перевод каретки не поможет.

Пример печатающего сервера

```bash
nc -vvl -s localhost -p 4242 -c '
    set -x
    sleep 5
    echo "This is server message"
    sleep 5
    echo "This is server message"
    sleep 5
    echo "This is server message"
    sleep 5
    echo "This is server message"
    sleep 5
    echo "This is server message"
    sleep 5
'
```

Клиент печатает от 0 до 9 бесконечно без нажатия на Enter.

```bash
./go-telnet localhost 4242
```

![./REPORT.files/3.gif]()

## Заметки для себя

### Поменял `.golangci.yml`

```bash
golangci-lint run --out-format=github-actions ./ 
::error file=telnet_test.go,line=11,col=2::import 'github.com/stretchr/testify/require' is not allowed from list 'Main' (depguard)
```

### Ничего себе, они еще живы и их немало

https://www.telnetbbsguide.com/

### Для отчета

```bash
cd ../hw11_telnet_client/
go mod tidy
rm -f ./.sync
golangci-lint run --out-format=github-actions ./ 

./test.native.sh >/dev/null 2>test.native.sh.out

./test.sh >/dev/null 2>test.sh.out
grep "FAIL" ./test.sh.out

go test -v telnet.go telnet_test.go > telnet_test.go.txt
go test -v main_test.go main.go foreign_code_base.go telnet.go  > main_test.go.txt
cd ../report_templator/
go test templator.go hw11_telnet_client_test.go
cd ../hw11_telnet_client/

grep ".txt }}" ./REPORT.md
grep ".out }}" ./REPORT.md
grep ".go }}"  ./REPORT.md
grep ".sh }}"  ./REPORT.md
```
