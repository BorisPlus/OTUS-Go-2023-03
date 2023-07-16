# Домашнее задание №10 «Оптимизация программы»

> "Я - художник, я так вижу." (Веронезе Паоло)

Описание [задания](./README.md).

## Статический анализ кода

В исходном [файле](./stats_initial.go) имеются конструкции кода, которые, исходя из моего опыта, являются узким местом реализации текущего алгоритма.

> В ниже представленных блоках исходного кода троеточие "..." скрывает конструкции, не имеющие отношения к текущему контексту повествования.

### 1. Явное ограничение на длину `100_000`

```go
type users [100_000]User
```

Сведения о размерах входных данных не известны. Если не предполагаем читать порционного по `100_000`, то лучше `slice`:

```go
type users []User
```

### 2. Полнообъемное чтение за раз

```go
... := io.ReadAll(r)
```

Лучше читать порциями, обрабатывая, например, файл построчно.

### 3. Излишнее приведение типов

```go
lines := strings.Split(string(content), "\n")
for i, line := range lines {
    ...
    if err = json.Unmarshal([]byte(line), &user); err != nil {
        return
    }
    ...
}
```

На моментах `string(content)` и `[]byte(line)` происходит последовательное приведение типов `[]byte` -> `string` и обратно `string` -> `[]byte` одних и тех же данных. Можно подобрать методы без последовательного перевода.

### 4. Отсутствие подхода повторного использования

Так, в:

```go
for i, line := range lines {
    var user User
    if err = json.Unmarshal([]byte(line), &user); err != nil {
        return
    }
    ...
}
```

`var user User` необходимо вынести на уровень выше

и в:

```go
num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
num++
result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
```

`strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])` заменить, например, на (**но и так я не буду делать**, см. далее):

```go
key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])

num := result[key]
num++
result[key] = num
```

### 5. Переприсваивание значения в рамках атомарной операции

```go
num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
num++
result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
```

если можно, например, так (**нет, но и так я не буду делать**, см. далее)

```go
result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
```

### 6. Излишнее оборачивание

```go
... fmt.Errorf("get users error: %w", err)
```

если можно, например, так

```go
... err
```

### 7. Доверие к данным

#### 7.1. К данным `user.Email`

```go
matched, err := regexp.Match("\\."+domain, []byte(user.Email))
...
if matched {
    num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
    ...
}
```

Точно ли `user.Email` содержит `@`, если оканчивается на `.com`, `.biz` или пр.?
Ведь если `user.Email` не содержит, то второй элемент (вот этот `strings.SplitN(user.Email, "@", 2)[1]`) даст ошибку:

```text
panic: runtime error: index out of range [1] with length 1
```

#### 7.2. К JSON-данным в принципе

```go
if err = json.Unmarshal([]byte(line), &user); err != nil {
    return
}
```

Если "дамп" данных битый, то `json.Unmarshal` даст ошибку `InvalidUnmarshalError`, но это не значит, что данные конкретно поля `user.Email` - невалидны. Они не попадут в статистику.

### 8. Регулярное выражение

Механизм регулярных выражений задействуется не в полной мере.

#### 8.1. Не скомпилировано

```go
... := regexp.Match("\\."+domain, ... )
```

Алгоритм в цикле производит разбор регулярного выражения вместо того, чтобы один раз "скомпилировать" его и искать удовлетворение уже скомпилированному:

* `func Compile(expr string) (*Regexp, error)` или
* `func MustCompile(str string) *Regexp` (если, как п.7, доверяем входному `str`).

#### 8.2. Нет поиска выявленного соответствия

```go
... := regexp.Match("\\."+domain, []byte(user.Email))
...
... strings.SplitN(user.Email, "@", 2) ...
```

Сначала алгоритм устанавливает факт отнесения содержания `user.Email` к домену, а потом извлекает этот домен после знака `@`. Вместо этого можно использовать функцию поиска "подсоответствия" (`submatch`) с использованием регулярных выражений:

* `func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string` или
* `func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte`.

### 9. Сам алгоритм

Если присмотреться ко всему пакету, то видно, что публичны только:

* `DomainStat`
* `GetDomainStatSource`
* `User`

При этом между собой явно связаны только `DomainStat` и `GetDomainStatSource`, а `User` как-то - "сбоку".

Скорее всего в проекте может задействоваться структура `User`, тогда ее необходимо вынести в отдельный файл.

Сложно сделать какие-то выводы о возможных источниках данных, на которых планируется применять исследуемую функцию `GetDomainStatSource`:

* Это исключительно строго типизированные JSON-данные в формате как в тестах и в `users.dat.zip`? Может ли присутствовать адрес электронной почты в каком-то другом строковом поле (например, в `Username`)?
* А может ли быть подан на вход функции `GetDomainStatSource`:
  * MongoDB-дамп с электронными почтами пользователей?
  * CSV-файл с электронными почтами пользователей?
  * Лог SMTP-сервера?
  * Плоский файл с текстовыми материалами и контактными реквизитами их авторов (адресами электронной почты)?

### 10. Хм

Ощущение, что я что-то забыл написать.

## Предложение по реализации

Я вижу, что задача `GetDomainStatSource` - **"Посчитать статистику доменов по почтовым адресам"**.

Я бы предложил сделать так (с некоторыми допущениями):

* Учесть все выше перечисленное.
* Сделать заполнение структуры `DomainStat` конкуретным (в несколько горутин).
* Реализовать метод разбора файла ("парсинга") без промежуточной структуры `User` посредством регулярного выражения, соответствующего произвольному адресу электронной почты.

> "Настоящий" многострочный вид регулярного выражения "адрес электронной почты" выходит за рамки данной реализации.

Все это **допустимо** по условию задачи, так как можно:

* писать любой новый необходимый код;
* удалять имеющийся лишний код (кроме функции `GetDomainStat`);

## Тестирование реализации

Первый переделанный вариант [stats_example](./stats_example.go) имел вид:

```go
{{ stats_example.go }}
```

> Наименованиям дописан суффикс "Example". Это сделано для отсутствия конфликтов по коду в итоговом проекте. Однокурсники, если вы захотите поэкспериментировать вместе со мной, то просто скопируйте себе этот файл, удалив суффиксы "Example". Но можете просто дочитать до конца.

### Проверка на работоспособность

```text
 go test -v -run=TestGetDomainStat ./
=== RUN   TestGetDomainStat
=== RUN   TestGetDomainStat/find_'com'
=== RUN   TestGetDomainStat/find_'gov'
=== RUN   TestGetDomainStat/find_'unknown'
--- PASS: TestGetDomainStat (0.62s)
    --- PASS: TestGetDomainStat/find_'com' (0.40s)
    --- PASS: TestGetDomainStat/find_'gov' (0.11s)
    --- PASS: TestGetDomainStat/find_'unknown' (0.10s)
PASS
ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  0.637s
```

### Нагрузочное тестирование

Функция проверки нагрузки имелась в изначальном репозитории:

```bash
go test -v -count=1 -timeout=30s -tags bench .
```

Нагрузка на изначальный вариант, представленный в учебном репозитории:

```go
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:46: time used: 1.58459501s / 300ms
    stats_optimization_test.go:47: memory used: 308Mb / 30Mb
    stats_optimization_test.go:49: 
                Error Trace:    stats_optimization_test.go:49
                Error:          "1584595010" is not less than "300000000"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       the program is too slow
--- FAIL: TestGetDomainStat_Time_And_Memory (1.59s)
FAIL
FAIL    github.com/fixme_my_friend/hw10_program_optimization    1.599s
FAIL
```

Нагрузка на реализованный мною вариант:

```go
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:46: time used: 330.227538ms / 300ms
    stats_optimization_test.go:47: memory used: 8Mb / 30Mb
    stats_optimization_test.go:49: 
                Error Trace:    stats_optimization_test.go:49
                Error:          "330227538" is not less than "300000000"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       the program is too slow
--- FAIL: TestGetDomainStat_Time_And_Memory (6.10s)
FAIL
FAIL    github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization    6.101s
FAIL
```

Видно, что нагрузка снизилась:

* на память - в 38.5 раз (`8Mb` против `308Mb`)
* по скорости - в 4.798 раз (`330.227538ms` против `1.58459501s`)

## Замечания

### Замечание о скорости алгоритма

Реализация не достигла требуемых `300ms`, что предположительно связано с локальными ограничениями на вычислительную мощность (железо). Так как, исходя из опроса аудитории курса, исполнение даже изначальной реализации не превышало по длительности 1 секунду, когда как у меня оно составило 1.58 секунды.

### Замечание о вариантах конкурентного доступа к Map-структуре

В рамках конкурентного доступа к объекту `DomainStat` для инкрементирования его значения я успешно опробовал:

* `sync.Mutex` (при этом с выводом в отдельную структуру).
* `sync/atomic` (особенность в том, что `DomainStat` по условию менять нельзя, а он содержит `int`, а `sync/atomic` имеет их конкретные архитектурные реализации `int32`/`int64` - `AddInt32`/`AddInt64`).
* `sync.Map` (описание в Приложении).

### Замечание о числе горутин-обработчиков

Число горутин-обработчиков задаются через переменную окружения `WORKERS_COUNT`:

```go
workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 1)
```

```bash
WORKERS_COUNT=1000 go test -v -count=1 -timeout=30s -tags bench .
```

### Замечание о большом объеме буфера чтения

Объем буфера чтения задается через переменную окружения `MAX_CAPACITY`:

```go
...
scanner := bufio.NewScanner(r)
maxCapacity := loadEnviromentOrDefault("MAX_CAPACITY", 2_000_000) // Magick default big value
fmt.Println("maxCapacity =", maxCapacity)
buf := make([]byte, maxCapacity)
scanner.Buffer(buf, maxCapacity)
...
```

и имеет большое значение по умолчанию - `2_000_000` (для текущего варианта алгоритма).

Вызвано это частично ограничением `bufio.NewScanner`:

```go
const (
    // MaxScanTokenSize is the maximum size used to buffer a token
    // unless the user provides an explicit buffer with Scanner.Buffer.
    // The actual maximum token size may be smaller as the buffer
    // may need to include, for instance, a newline.
    MaxScanTokenSize = 64 * 1024
)
```

и неустановленной причиной, которую я и рассматриваю далее. 

Тестируем с этим **ОБЪЕМНЫМ** буфером (`2_000_000`) и **БОЛЬШИМ** числом горутин (`1000`), тогда будет все **ОК**:

```bash
MAX_CAPACITY=2000000 WORKERS_COUNT=1000 go test -v -run=TestBigDataGetDomainStat ./

    === RUN   TestBigDataGetDomainStat
    --- PASS: TestBigDataGetDomainStat (0.43s)
    PASS
    ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  (cached)
```

Тестируем с **ОБЪЯСНИМЫМ ОБЪЕМОМ** буфера (`100_000` - число строк в тестовом файле) и **ОДНОЙ** горутиной, тогда будет все **ОК**:

```bash
MAX_CAPACITY=100000 WORKERS_COUNT=1 go test -v -run=TestBigDataGetDomainStat ./

    === RUN   TestBigDataGetDomainStat
    --- PASS: TestBigDataGetDomainStat (0.55s)
    PASS
    ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  0.556s
```

а вот, например, **ТОТ ЖЕ** буфер (`100_000`) и уже **ДВЕ** горутины, тогда будет **FAIL**:

```bash
MAX_CAPACITY=100000 WORKERS_COUNT=2 go test -v -run=TestBigDataGetDomainStat ./
...
...
--- FAIL: TestBigDataGetDomainStat (0.52s)
FAIL
```

Почему именно с **ОБЪЕМНЫМ** буфером (`2_000_000`) работает сразу? Вед это значение не соответствует ни одному из:

* Размер архива `users.dat.zip` - `5300876` байт.
* Размер файла `users.dat` - `17375349` байт после распаковки `users.dat.zip`.
* В файле `users.dat` - `100000` строк.
* Максимальная длина считываемой строки тестового файла - `239` байт.
* benchmark-тест запускается "произвольное" число раз.

При этом при периодическом запуске benchmark-теста в случае его "случайного" успеха это значение в `2_000_000` можно постепенно снижать и снижать (видимо влияют "горячие" данные, будет продемонстрировано).

Запустим специально подготовленный код [`main.go`](experimantal/main.go), с помощью которого "вручную" сверим результат посчитанной статистики с ожидаемыми эталонными данными (не хотелось менять исходный тест, так как по условию задачи - это запрещено, можно, конечно, просто его продублировать, но решено так):

<details><summary>file: `experimantal/main.go`</summary>

```go
{{ experimantal/main.go }}
```

</details>

```bash
cd ./experimantal/
```

Смотрите, вот уже установленное ограничение `MaxScanTokenSize` в `bufio`.

С `100_000` работает стабильно успешно:

```bash
MAX_CAPACITY=100000 WORKERS_COUNT=1 go run ./main.go 
    
    WORKERS_COUNT = 1
    MAX_CAPACITY  = 64000
    I get GetDomainStat
    Let's check it with ethalon
    OK
```

С `64_000` не работает (периодически выпадает `FAIL`):

```bash
MAX_CAPACITY=64000 WORKERS_COUNT=1 go run ./main.go 

    WORKERS_COUNT = 1
    MAX_CAPACITY  = 63999
    I get GetDomainStat
    Let's check it with ethalon
    centimia.biz <---- та "ВАЖНАЯ СТРОКА" (см. комментарий выше), ну ОК, допустим что-то там не посчиталось
    FAIL

MAX_CAPACITY=64000 WORKERS_COUNT=1 go run ./main.go 

    WORKERS_COUNT = 1
    MAX_CAPACITY  = 63999
    I get GetDomainStat
    Let's check it with ethalon
    OK
```

Если пристально рассмотреть расхождение, то обнаружится, что в результате статистики `type DomainStat map[string]int` имеется:

* В string-ключах отсутствуют некоторые домены.
* В int-значениях некоторых доменов показатели больше или меньше эталонных.
* **Самое удивительное** то, иногда в string-ключах присутствуют **произвольные** короткие подстроки из тестового файла (**вообще не удовлетворяющие регулярному выражению электронной почты**, например, середина json-строки со знаком двоеточия ":"), при этом и в самой json-строке почта тоже **никак не удовлетворяет** регулярному выражению.

Примеры нестабильного поведения:

```bash
$ MAX_CAPACITY=64000 WORKERS_COUNT=2 go run ./main.go 
    WORKERS_COUNT = 2
    MAX_CAPACITY  = 100000
    I get GetDomainStat
    Let's check it with ethalon
    voolia.biz <---- та "ВАЖНАЯ СТРОКА" (см. комментарий выше), ну ОК, допустим что-то там не посчиталось
    FAIL

$ MAX_CAPACITY=128000 WORKERS_COUNT=2 go run ./main.go 
    WORKERS_COUNT = 2
    MAX_CAPACITY  = 128000
    I get GetDomainStat
    Let's check it with ethalon
    ,"username" <---- та "ВАЖНАЯ СТРОКА" (см. комментарий выше), как она вообще с запятой сюда попала ?!
    FAIL
```

Не буду томить.

## Переформатируем логику и всё заработает

Вынесем логику работы с регулярным выражением за пределы горутин, подавая к ним на вход канал готовых доменов для инкремента персональной статистики:

```go
{{ stat.go }}
```

И теперь все заработает даже с меньшим объемом `MAX_CAPACITY`, который не должен быть меньше `239`, так как именно такова максимальная длина входной строки репозиторного тестового файла из архива `users.dat.zip`, и таким образом, именно такой максимальный объем необходим для буфера чтения.

```bash
MAX_CAPACITY=239 WORKERS_COUNT=1 go test -v -count=1 -timeout=30s -tags bench 

    === RUN   TestGetDomainStat_Time_And_Memory
        stats_optimization_test.go:46: time used: 722.672834ms / 300ms
        stats_optimization_test.go:47: memory used: 3Mb / 30Mb
        stats_optimization_test.go:49: 
                    Error Trace:    stats_optimization_test.go:49
                    Error:          "722672834" is not less than "300000000"
                    Test:           TestGetDomainStat_Time_And_Memory
                    Messages:       the program is too slow
    --- FAIL: TestGetDomainStat_Time_And_Memory (20.06s)
    FAIL
    exit status 1
    FAIL    github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  20.066s

MAX_CAPACITY=239 WORKERS_COUNT=10 go test -v -count=1 -timeout=30s -tags bench 

    === RUN   TestGetDomainStat_Time_And_Memory
        stats_optimization_test.go:46: time used: 407.30801ms / 300ms
        stats_optimization_test.go:47: memory used: 3Mb / 30Mb
        stats_optimization_test.go:49: 
                    Error Trace:    stats_optimization_test.go:49
                    Error:          "407308010" is not less than "300000000"
                    Test:           TestGetDomainStat_Time_And_Memory
                    Messages:       the program is too slow
    --- FAIL: TestGetDomainStat_Time_And_Memory (10.51s)
    FAIL
    exit status 1
    FAIL    github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  10.519s

MAX_CAPACITY=239 WORKERS_COUNT=100 go test -v -count=1 -timeout=30s -tags bench 
    === RUN   TestGetDomainStat_Time_And_Memory
        stats_optimization_test.go:46: time used: 414.144843ms / 300ms
        stats_optimization_test.go:47: memory used: 3Mb / 30Mb
        stats_optimization_test.go:49: 
                    Error Trace:    stats_optimization_test.go:49
                    Error:          "414144843" is not less than "300000000"
                    Test:           TestGetDomainStat_Time_And_Memory
                    Messages:       the program is too slow
    --- FAIL: TestGetDomainStat_Time_And_Memory (10.73s)
    FAIL
    exit status 1
    FAIL    github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  10.740s
```

Однако увеличение числа горутин с `1` до `100` все же не позволяет приблизиться к успеху в `300` миллисекунд. Как указано выше, скорее всего это "железо".

## Сравнение производительности

Для сравнения был реализован benchmark-тест:

<details><summary>file: `stats_benchmark_test.go`</summary>

```go
{{ stats_benchmark_test.go }}
```

</details>

наглядно показавший увеличение скорости исполнения:

```bash
go test -bench=.

    goos: linux
    goarch: amd64
    pkg: github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization
    cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz

    BenchmarkStat001Repo-4                1      1566185476 ns/op
    BenchmarkStat002My-4         1000000000      0.4648 ns/op

    PASS
    ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  13.380s
```

</details>

```bash
GOGC=off go test -bench=BenchmarkStat002My -cpuprofile cpu.out

    goos: linux
    goarch: amd64
    pkg: github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization
    cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz
    BenchmarkStat002My-4    1000000000               0.4082 ns/op
    PASS
    ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  10.851s

go tool pprof -svg ./hw10_program_optimization.test ./cpu.out > ./REPORT.files/cpu.svg
```

Исходя из [графа вызовов](./REPORT.files/cpu.svg):

![REPORT.files/cpu.svg](./REPORT.files/cpu.svg)

* На скорость декомпрессии мне не повлиять никак, как и ...
* На скорость работы скомпилированного регулярного выражения
  > Проверено, что использование `strings.SplitN(email,"@", 2)[1]` для извлечения домена не ускоряет работу (они состязаются):
  >
  > ```bash
  > go test -bench=.
  > 
  >     goos: linux
  >     goarch: amd64
  >     pkg: github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization
  >     cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz
  >     BenchmarkStat001Repo-4                           1        2019547597 ns/op
  >     BenchmarkStat002My-4                    1000000000        0.4622     ns/op
  >     BenchmarkStat003Experimental-4          1000000000        0.3936     ns/op
  >     PASS
  >     ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  18.818s
  >  go test -bench=.
  >   goos: linux
  >     goarch: amd64
  >     pkg: github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization
  >     cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz
  >     BenchmarkStat001Repo-4                           1        1620842078 ns/op
  >     BenchmarkStat002My-4                    1000000000        0.4703     ns/op
  >     BenchmarkStat003Experimental-4          1000000000        0.5113     ns/op
  >     PASS
  >     ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  25.584s
  > ```

## Вывод. Вот в чем проблема
  
Так как изменения промежуточного варианта к исходному коснулись только двух вещей в горутинах, то предполагаю, что узкое место в промежуточном варианте заключалось:

* Либо в передаче объекта скомпилированного регулярного выражения `regexp.Regexp`, которое в итоге не корректно извлекало подстроки. Это вариант маловероятен, так как переменная скомпилированного регулярного выражения передавалась по значению (в итоге отказался от передачи горутинам скомпилированного регулярного выражения).
* Либо в канале слайса байт `[]byte` (он был изменен на канал строк `string`). Я не поверил, но именно в этом варианте дело. Я изменил промежуточный вариант только в части, касающейся канала строк, все также передавая горутинам скомпилированное регулярное выражение (можно и по ссылке). Ошибок соответствия эталону на малых значениях `MAX_CAPACITY` (`239`) больше **НЕТ**.

<details><summary>file: `stats_remark.go`</summary>

```go
{{ stats_remark.go }}
```

</details>

```text
go test -v -bench=.

    goos: linux
    goarch: amd64
    pkg: github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization
    cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz
    BenchmarkStat001Repo-4                         1        1810331956 ns/op
    BenchmarkStat002My-4                    1000000000               0.4144 ns/op
    BenchmarkStat003Experimental-4          1000000000               0.4073 ns/op
    BenchmarkStat004Remark-4                1000000000               0.4036 ns/op
    PASS
    ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  34.815s
```

Задача решена.

А по достижении скорости исполнения `300` миллисекунд, однокурсники, пожалуйста, замерьте у себя мою реализацию [stats.go](./stats.go).

## Приложение

### sync.Map

Конкурентное заполнение структуры `sync.Map`:

```go
v, ok := syncMap.LoadOrStore(domainAtLowercase, 1)
if ok {
    syncMap.Store(domainAtLowercase, v.(int)+1)
}
```

Заполнение сигнатурного `DomainStat`:

```go
syncMap.Range(func(key, value interface{}) bool {
    domainStat[key.(string)] = value.(int)
    return true
})
```

### Чтение чанками

```go
b := make([]byte, 1)
chunk := make([]byte, 0)
for {
    _, err := r.Read(b)
    if err != nil {
        if errors.Is(err, io.EOF) {
            rowsChannel <- chunk
            break
        }
        fmt.Println(err)
        break
    }
    if b[0] == '\n' {
        rowsChannel <- chunk
        chunk = make([]byte, 0)
        continue
    }
    chunk = append(chunk, b...)
}
```

### bufio.NewReaderSize

```go
reader := bytes.NewReader(r)
scanners := bufio.NewReaderSize(reader, 10000000)
```
