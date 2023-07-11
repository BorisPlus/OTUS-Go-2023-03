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
=== RUN   TestTelnetClient
=== RUN   TestTelnetClient/basic
...Try connect to 127.0.0.1:45823
...Connected to 127.0.0.1:45823
...Disconnected from 127.0.0.1:45823
--- PASS: TestTelnetClient (0.00s)
    --- PASS: TestTelnetClient/basic (0.00s)
PASS
ok      command-line-arguments    0.006s

```

### Shell

Был вынужден внести правки в `test.sh`, так как в `nc` версии `[v1.10-47]` под `Debian` необходимы дополнительные параметры `-s` и `-p` для запуска локального TCP-сервера: `nc -ls localhost -p 4242`. Наглядно видно на примере:

<details>
<summary>test.native.sh</summary>

```bash
#!/usr/bin/env bash
set -xeuo pipefail

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -ls localhost -p 4242 >/tmp/nc.out &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | nc localhost 4242 > /tmp/telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "FAIL: unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals /tmp/nc.out "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals /tmp/telnet.out "${expected_telnet_out}"

echo "PASS"

```

</details>

> Для сохранения лога запускать как:
>
> ```bash
> ./test.native.sh >/dev/null 2>test.native.sh.out
>

```text
+ echo -e 'Hello\nFrom\nNC\n'
+ cat
+ nc -ls localhost -p 4242
+ NC_PID=204346
+ sleep 1
+ echo -e 'I\nam\nTELNET client\n'
+ TL_PID=204369
+ cat
+ sleep 5
+ nc localhost 4242
+ kill 204369

+ kill 204346

+ expected_nc_out='I
am
TELNET client'
+ fileEquals /tmp/nc.out 'I
am
TELNET client'
+ local fileData
++ cat /tmp/nc.out
+ fileData='I
am
TELNET client'
+ '[' 'I
am
TELNET client' = 'I
am
TELNET client' ']'
+ expected_telnet_out='Hello
From
NC'
+ fileEquals /tmp/telnet.out 'Hello
From
NC'
+ local fileData
++ cat /tmp/telnet.out
+ fileData='Hello
From
NC'
+ '[' 'Hello
From
NC' = 'Hello
From
NC' ']'
+ echo PASS

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
+ go build -o go-telnet.goc
+ echo -e 'Hello\nFrom\nNC\n'
+ NC_PID=204500
+ nc -ls localhost -p 4242
+ sleep 1
+ cat
+ echo -e 'I\nam\nTELNET client\n'
+ ./go-telnet.goc localhost 4242
+ TL_PID=204507
+ cat
+ sleep 5
+ kill 204507
+ kill 204500

+ expected_nc_out='I
am
TELNET client'
+ fileEquals /tmp/nc.out 'I
am
TELNET client'
+ local fileData
++ cat /tmp/nc.out
+ fileData='I
am
TELNET client'
+ '[' 'I
am
TELNET client' = 'I
am
TELNET client' ']'
+ expected_telnet_out='...Try connect to localhost:4242
...Connected to localhost:4242
Hello
From
NC'
+ fileEquals /tmp/telnet.out '...Try connect to localhost:4242
...Connected to localhost:4242
Hello
From
NC'
+ local fileData
++ cat /tmp/telnet.out
+ fileData='...Try connect to localhost:4242
...Connected to localhost:4242
Hello
From
NC'
+ '[' '...Try connect to localhost:4242
...Connected to localhost:4242
Hello
From
NC' = '...Try connect to localhost:4242
...Connected to localhost:4242
Hello
From
NC' ']'
+ rm -f go-telnet.goc
+ echo PASS

```

### Тестирование аргументов

> Частично задействуется чужая кодовая база в отношении проверки соответствия наименования хоста валидному доменному имени. Функциями стандартных модулей пакета `net` (`net.url.Parse`, `net.SplitHostPort`) этого не добился, они уверенно допускают невалидные домены.

```bash
go test -v main_test.go main.go foreign_code_base.go telnet.go > main_test.go.txt
```

```text
=== RUN   TestArgParsePositive
It's Ok. for args [--timeout=10s localhost 4242]
It's Ok. for args [--timeout=5s localhost.com 23]
It's Ok. for args [--timeout=11s telnet.localhost.com 23]
It's Ok. for args [127.0.0.1 23]
It's Ok. for args [--timeout=10s 1.1.1.1 65535]
--- PASS: TestArgParsePositive (0.00s)
=== RUN   TestArgParseNegative
It's Ok. Get expected error HOST (as domain): invalid character '.' at offset 0: label can't begin with a period
         for args [--timeout=5s .40ca1host.com 23]
It's Ok. Get expected error HOST (as domain): top level domain '1' at offset 17 begins with a digit
         for args [--timeout=11s telnet.localhost.1 23]
It's Ok. Get expected error HOST (as domain): invalid character '=' at offset 9
         for args [--timeout=11s telnet.net.]
It's Ok. Get expected error TIMEOUT parsing error
         for args [127.0.0.1 --timeout=11s 23]
It's Ok. Get expected error HOST (as domain): top level domain '257' at offset 6 begins with a digit
         for args [--timeout=10s 1.1.1.257 65535]
It's Ok. Get expected error HOST (as domain): top level domain '1' at offset 6 begins with a digit
         for args [--timeout=10s 1.1.1.1 65537]
It's Ok. Get expected error HOST (ip-address): is gateway
         for args [--timeout=10s 1.1.1.0 1]
--- PASS: TestArgParseNegative (0.00s)
PASS
ok      command-line-arguments    0.006s

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
