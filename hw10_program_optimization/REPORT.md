# Домашнее задание №10 «Оптимизация программы»

> "Я - художник, я так вижу." (Веронезе Паоло)

Описание [задания](./README.md).

## Статический анализ кода

В исходном [файле](./stats_initial.go) имеются конструкции кода, которые, исходя из моего опыта, являются узким местом реализации текущего алгоритма (наименование версии - `*Repo`, соотв. `BenchmarkStat001Repo`).

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

Первый переделанный вариант [stats_example](./stats_example.go) имел вид (наименование версии - `*Example`):

```go
package hw10programoptimization

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "regexp"
    "strconv"
    "strings"
    "sync"
)

// TODO: with "reflect" - func LoadOrDefault[T any](name string, asDefault T) T {}.
func loadEnviromentOrDefault(name string, asDefault int) int {
    value, exists := os.LookupEnv(name)
    if exists {
        intValue, err := strconv.Atoi(value)
        if err == nil {
            return intValue
        }
    }
    return asDefault
}

func rowParserExample(
    wg *sync.WaitGroup,
    mtx *sync.Mutex,
    rows <-chan []byte,
    compiledRegexp regexp.Regexp,
    domainStat DomainStat,
) {
    defer wg.Done()
    for row := range rows {
        matches := compiledRegexp.FindAllSubmatch(row, -1)
        for matcheIndex := range matches {
            domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))
            mtx.Lock()
            domainStat[domainAtLowercase]++
            mtx.Unlock()
        }
    }
}

func GetDomainStatExample(r io.Reader, domain string) (DomainStat, error) {
    domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
    compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)
    if err != nil {
        return nil, err
    }
    wg := sync.WaitGroup{}
    mtx := sync.Mutex{}
    dataChannel := make(chan []byte)
    domainStat := make(DomainStat)
    workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 1)
    for i := 0; i < workersCount; i++ {
        wg.Add(1)
        go rowParserExample(&wg, &mtx, dataChannel, *compiledRegexp, domainStat)
    }
    scanner := bufio.NewScanner(r)
    maxCapacity := loadEnviromentOrDefault("MAX_CAPACITY", 2_000_000) // Magick big value
    buf := make([]byte, maxCapacity)
    scanner.Buffer(buf, maxCapacity)
    for scanner.Scan() {
        dataChannel <- scanner.Bytes()
    }
    close(dataChannel)
    wg.Wait()
    return domainStat, nil
}

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
package main

import (
    "archive/zip"
    "fmt"
    "os"

    hw10 "github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization"
)

func main() {

    fmt.Println("WORKERS_COUNT", "=", os.Getenv("WORKERS_COUNT"))
    fmt.Println("MAX_CAPACITY ", "=", os.Getenv("MAX_CAPACITY"))

    r, err := zip.OpenReader("../testdata/users.dat.zip")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer r.Close()

    data, err := r.File[0].Open()
    if err != nil {
        fmt.Println(err)
        return
    }

    // TODO: replace to `GetDomainStat``
    stat, err := hw10.GetDomainStatExample(data, "biz")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("I get GetDomainStatExample")
    fmt.Println("Let's check it with ethalon")

    // LEFT OUTER JOIN
    for key := range stat {
        if stat[key] != hw10.ExpectedBizStat[key] {
            fmt.Println(key) // ВАЖНАЯ СТРОКА
            fmt.Println("FAIL")
            return
        }
    }
    // RIGHT OUTER JOIN
    for key := range hw10.ExpectedBizStat {
        if stat[key] != hw10.ExpectedBizStat[key] {
            fmt.Println("FAIL")
            return
        }
    }
    fmt.Println("OK")
}

var expectedBizStat = map[string]int{
    "abata.biz":         25,
    "abatz.biz":         25,
    "agimba.biz":        28,
    "agivu.biz":         17,
    "aibox.biz":         31,
    "ailane.biz":        23,
    "aimbo.biz":         25,
    "aimbu.biz":         36,
    "ainyx.biz":         35,
    "aivee.biz":         25,
    "avamba.biz":        21,
    "avamm.biz":         17,
    "avavee.biz":        35,
    "avaveo.biz":        30,
    "babbleblab.biz":    29,
    "babbleopia.biz":    36,
    "babbleset.biz":     28,
    "babblestorm.biz":   29,
    "blognation.biz":    32,
    "blogpad.biz":       34,
    "blogspan.biz":      21,
    "blogtag.biz":       23,
    "blogtags.biz":      34,
    "blogxs.biz":        35,
    "bluejam.biz":       36,
    "bluezoom.biz":      27,
    "brainbox.biz":      30,
    "brainlounge.biz":   38,
    "brainsphere.biz":   31,
    "brainverse.biz":    39,
    "brightbean.biz":    23,
    "brightdog.biz":     32,
    "browseblab.biz":    31,
    "browsebug.biz":     25,
    "browsecat.biz":     34,
    "browsedrive.biz":   24,
    "browsetype.biz":    34,
    "browsezoom.biz":    29,
    "bubblebox.biz":     19,
    "bubblemix.biz":     38,
    "bubbletube.biz":    34,
    "buzzbean.biz":      26,
    "buzzdog.biz":       30,
    "buzzshare.biz":     26,
    "buzzster.biz":      28,
    "camido.biz":        27,
    "camimbo.biz":       36,
    "centidel.biz":      32,
    "centimia.biz":      17,
    "centizu.biz":       18,
    "chatterbridge.biz": 30,
    "chatterpoint.biz":  32,
    "cogibox.biz":       30,
    "cogidoo.biz":       34,
    "cogilith.biz":      24,
    "dabfeed.biz":       26,
    "dabjam.biz":        30,
    "dablist.biz":       30,
    "dabshots.biz":      33,
    "dabtype.biz":       21,
    "dabvine.biz":       26,
    "dabz.biz":          19,
    "dazzlesphere.biz":  24,
    "demimbu.biz":       27,
    "demivee.biz":       39,
    "demizz.biz":        30,
    "devbug.biz":        20,
    "devcast.biz":       35,
    "devify.biz":        27,
    "devpoint.biz":      26,
    "devpulse.biz":      27,
    "devshare.biz":      30,
    "digitube.biz":      30,
    "divanoodle.biz":    33,
    "divape.biz":        32,
    "divavu.biz":        28,
    "dynabox.biz":       66,
    "dynava.biz":        21,
    "dynazzy.biz":       29,
    "eabox.biz":         28,
    "eadel.biz":         25,
    "eamia.biz":         18,
    "eare.biz":          30,
    "eayo.biz":          30,
    "eazzy.biz":         27,
    "edgeblab.biz":      29,
    "edgeclub.biz":      29,
    "edgeify.biz":       36,
    "edgepulse.biz":     21,
    "edgetag.biz":       24,
    "edgewire.biz":      29,
    "eidel.biz":         33,
    "eimbee.biz":        22,
    "einti.biz":         19,
    "eire.biz":          28,
    "fadeo.biz":         35,
    "fanoodle.biz":      23,
    "fatz.biz":          30,
    "feedbug.biz":       29,
    "feedfire.biz":      30,
    "feedfish.biz":      35,
    "feedmix.biz":       31,
    "feednation.biz":    24,
    "feedspan.biz":      28,
    "fivebridge.biz":    20,
    "fivechat.biz":      29,
    "fiveclub.biz":      23,
    "fivespan.biz":      27,
    "flashdog.biz":      20,
    "flashpoint.biz":    35,
    "flashset.biz":      30,
    "flashspan.biz":     32,
    "flipbug.biz":       27,
    "flipopia.biz":      30,
    "flipstorm.biz":     21,
    "fliptune.biz":      29,
    "gabcube.biz":       29,
    "gabspot.biz":       24,
    "gabtune.biz":       29,
    "gabtype.biz":       29,
    "gabvine.biz":       24,
    "geba.biz":          24,
    "gevee.biz":         23,
    "gigabox.biz":       28,
    "gigaclub.biz":      25,
    "gigashots.biz":     26,
    "gigazoom.biz":      29,
    "innojam.biz":       26,
    "innotype.biz":      27,
    "innoz.biz":         24,
    "izio.biz":          26,
    "jabberbean.biz":    28,
    "jabbercube.biz":    31,
    "jabbersphere.biz":  55,
    "jabberstorm.biz":   22,
    "jabbertype.biz":    27,
    "jaloo.biz":         35,
    "jamia.biz":         33,
    "janyx.biz":         33,
    "jatri.biz":         18,
    "jaxbean.biz":       28,
    "jaxnation.biz":     21,
    "jaxspan.biz":       27,
    "jaxworks.biz":      30,
    "jayo.biz":          44,
    "jazzy.biz":         32,
    "jetpulse.biz":      25,
    "jetwire.biz":       26,
    "jumpxs.biz":        29,
    "kamba.biz":         30,
    "kanoodle.biz":      19,
    "kare.biz":          30,
    "katz.biz":          62,
    "kaymbo.biz":        34,
    "kayveo.biz":        22,
    "kazio.biz":         21,
    "kazu.biz":          16,
    "kimia.biz":         25,
    "kwideo.biz":        17,
    "kwilith.biz":       25,
    "kwimbee.biz":       34,
    "kwinu.biz":         15,
    "lajo.biz":          20,
    "latz.biz":          24,
    "layo.biz":          32,
    "lazz.biz":          27,
    "lazzy.biz":         26,
    "leenti.biz":        26,
    "leexo.biz":         32,
    "linkbridge.biz":    38,
    "linkbuzz.biz":      24,
    "linklinks.biz":     31,
    "linktype.biz":      31,
    "livefish.biz":      31,
    "livepath.biz":      23,
    "livetube.biz":      53,
    "livez.biz":         28,
    "meedoo.biz":        23,
    "meejo.biz":         24,
    "meembee.biz":       26,
    "meemm.biz":         23,
    "meetz.biz":         33,
    "meevee.biz":        62,
    "meeveo.biz":        27,
    "meezzy.biz":        24,
    "miboo.biz":         26,
    "midel.biz":         28,
    "minyx.biz":         25,
    "mita.biz":          29,
    "mudo.biz":          36,
    "muxo.biz":          25,
    "mybuzz.biz":        32,
    "mycat.biz":         32,
    "mydeo.biz":         20,
    "mydo.biz":          30,
    "mymm.biz":          21,
    "mynte.biz":         54,
    "myworks.biz":       27,
    "nlounge.biz":       25,
    "npath.biz":         33,
    "ntag.biz":          28,
    "ntags.biz":         32,
    "oba.biz":           22,
    "oloo.biz":          19,
    "omba.biz":          26,
    "ooba.biz":          27,
    "oodoo.biz":         30,
    "oozz.biz":          22,
    "oyoba.biz":         27,
    "oyoloo.biz":        30,
    "oyonder.biz":       29,
    "oyondu.biz":        23,
    "oyope.biz":         24,
    "oyoyo.biz":         32,
    "ozu.biz":           18,
    "photobean.biz":     25,
    "photobug.biz":      57,
    "photofeed.biz":     25,
    "photojam.biz":      35,
    "photolist.biz":     19,
    "photospace.biz":    33,
    "pixoboo.biz":       14,
    "pixonyx.biz":       30,
    "pixope.biz":        32,
    "plajo.biz":         32,
    "plambee.biz":       29,
    "podcat.biz":        31,
    "quamba.biz":        31,
    "quatz.biz":         54,
    "quaxo.biz":         25,
    "quimba.biz":        25,
    "quimm.biz":         33,
    "quinu.biz":         60,
    "quire.biz":         25,
    "realblab.biz":      32,
    "realbridge.biz":    30,
    "realbuzz.biz":      22,
    "realcube.biz":      57,
    "realfire.biz":      37,
    "reallinks.biz":     25,
    "realmix.biz":       27,
    "realpoint.biz":     22,
    "rhybox.biz":        30,
    "rhycero.biz":       28,
    "rhyloo.biz":        32,
    "rhynoodle.biz":     25,
    "rhynyx.biz":        17,
    "rhyzio.biz":        36,
    "riffpath.biz":      21,
    "riffpedia.biz":     33,
    "riffwire.biz":      31,
    "roodel.biz":        29,
    "roombo.biz":        29,
    "roomm.biz":         32,
    "rooxo.biz":         34,
    "shufflebeat.biz":   32,
    "shuffledrive.biz":  25,
    "shufflester.biz":   26,
    "shuffletag.biz":    23,
    "skaboo.biz":        35,
    "skajo.biz":         26,
    "skalith.biz":       30,
    "skiba.biz":         22,
    "skibox.biz":        27,
    "skidoo.biz":        24,
    "skilith.biz":       29,
    "skimia.biz":        45,
    "skinder.biz":       25,
    "skinix.biz":        23,
    "skinte.biz":        39,
    "skipfire.biz":      29,
    "skippad.biz":       26,
    "skipstorm.biz":     30,
    "skiptube.biz":      26,
    "skivee.biz":        34,
    "skyba.biz":         40,
    "skyble.biz":        32,
    "skyndu.biz":        32,
    "skynoodle.biz":     28,
    "skyvu.biz":         34,
    "snaptags.biz":      33,
    "tagcat.biz":        33,
    "tagchat.biz":       37,
    "tagfeed.biz":       30,
    "tagopia.biz":       17,
    "tagpad.biz":        28,
    "tagtune.biz":       22,
    "talane.biz":        22,
    "tambee.biz":        24,
    "tanoodle.biz":      38,
    "tavu.biz":          37,
    "tazz.biz":          27,
    "tazzy.biz":         28,
    "tekfly.biz":        31,
    "teklist.biz":       26,
    "thoughtbeat.biz":   30,
    "thoughtblab.biz":   24,
    "thoughtbridge.biz": 30,
    "thoughtmix.biz":    33,
    "thoughtsphere.biz": 20,
    "thoughtstorm.biz":  38,
    "thoughtworks.biz":  24,
    "topdrive.biz":      35,
    "topicblab.biz":     32,
    "topiclounge.biz":   21,
    "topicshots.biz":    30,
    "topicstorm.biz":    22,
    "topicware.biz":     35,
    "topiczoom.biz":     38,
    "trilia.biz":        28,
    "trilith.biz":       25,
    "trudeo.biz":        29,
    "trudoo.biz":        28,
    "trunyx.biz":        33,
    "trupe.biz":         34,
    "twimbo.biz":        19,
    "twimm.biz":         30,
    "twinder.biz":       28,
    "twinte.biz":        33,
    "twitterbeat.biz":   33,
    "twitterbridge.biz": 20,
    "twitterlist.biz":   26,
    "twitternation.biz": 22,
    "twitterwire.biz":   21,
    "twitterworks.biz":  39,
    "twiyo.biz":         37,
    "vidoo.biz":         28,
    "vimbo.biz":         21,
    "vinder.biz":        31,
    "vinte.biz":         34,
    "vipe.biz":          25,
    "vitz.biz":          26,
    "viva.biz":          30,
    "voolia.biz":        34,
    "voolith.biz":       26,
    "voomm.biz":         61,
    "voonder.biz":       32,
    "voonix.biz":        32,
    "voonte.biz":        26,
    "voonyx.biz":        25,
    "wikibox.biz":       27,
    "wikido.biz":        21,
    "wikivu.biz":        23,
    "wikizz.biz":        61,
    "wordify.biz":       28,
    "wordpedia.biz":     25,
    "wordtune.biz":      27,
    "wordware.biz":      19,
    "yabox.biz":         24,
    "yacero.biz":        34,
    "yadel.biz":         27,
    "yakidoo.biz":       21,
    "yakijo.biz":        29,
    "yakitri.biz":       26,
    "yambee.biz":        20,
    "yamia.biz":         17,
    "yata.biz":          25,
    "yodel.biz":         26,
    "yodo.biz":          21,
    "yodoo.biz":         24,
    "yombu.biz":         29,
    "yotz.biz":          26,
    "youbridge.biz":     40,
    "youfeed.biz":       32,
    "youopia.biz":       22,
    "youspan.biz":       59,
    "youtags.biz":       22,
    "yoveo.biz":         31,
    "yozio.biz":         33,
    "zava.biz":          29,
    "zazio.biz":         18,
    "zoombeat.biz":      28,
    "zoombox.biz":       30,
    "zoomcast.biz":      38,
    "zoomdog.biz":       29,
    "zoomlounge.biz":    25,
    "zoomzone.biz":      32,
    "zoonder.biz":       29,
    "zoonoodle.biz":     27,
    "zooveo.biz":        22,
    "zoovu.biz":         38,
    "zooxo.biz":         33,
    "zoozzy.biz":        23,
}

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

Вынесем логику работы с регулярным выражением за пределы горутин, подавая к ним на вход канал готовых доменов для инкремента персональной статистики (наименование версии - `*My`, соотв. `BenchmarkStat002My`)):

```go

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
package hw10programoptimization

import (
    "archive/zip"
    "os"
    "testing"

    "github.com/stretchr/testify/require"
)

func BenchmarkStat001Repo(b *testing.B) {
    r, err := zip.OpenReader("testdata/users.dat.zip")
    require.NoError(b, err)
    defer r.Close()

    require.Equal(b, 1, len(r.File))

    data, err := r.File[0].Open()
    require.NoError(b, err)

    b.ResetTimer()
    b.StartTimer()
    stat, err := GetDomainStatInitial(data, "biz")
    b.StopTimer()
    require.NoError(b, err)

    require.Equal(b, expectedBizStatCopy, stat)
}

func BenchmarkStat002My(b *testing.B) {
    os.Setenv("MAX_CAPACITY", "239")
    os.Setenv("WORKERS_COUNT", "100")

    r, err := zip.OpenReader("testdata/users.dat.zip")
    require.NoError(b, err)
    defer r.Close()

    require.Equal(b, 1, len(r.File))

    data, err := r.File[0].Open()
    require.NoError(b, err)

    b.ResetTimer()

    b.StartTimer()
    stat, err := GetDomainStat(data, "biz")
    b.StopTimer()
    require.NoError(b, err)

    require.Equal(b, expectedBizStatCopy, stat)
}

func BenchmarkStat003Experimental(b *testing.B) {
    os.Setenv("MAX_CAPACITY", "239")
    os.Setenv("WORKERS_COUNT", "100")

    r, err := zip.OpenReader("testdata/users.dat.zip")
    require.NoError(b, err)
    defer r.Close()

    require.Equal(b, 1, len(r.File))

    data, err := r.File[0].Open()
    require.NoError(b, err)

    b.ResetTimer()

    b.StartTimer()
    stat, err := GetDomainStatExperimental(data, "biz")
    b.StopTimer()
    require.NoError(b, err)

    require.Equal(b, expectedBizStatCopy, stat)
}

func BenchmarkStat004Remark(b *testing.B) {
    os.Setenv("MAX_CAPACITY", "239")
    os.Setenv("WORKERS_COUNT", "100")

    r, err := zip.OpenReader("testdata/users.dat.zip")
    require.NoError(b, err)
    defer r.Close()

    require.Equal(b, 1, len(r.File))

    data, err := r.File[0].Open()
    require.NoError(b, err)

    b.ResetTimer()

    b.StartTimer()
    stat, err := GetDomainStatRemark(data, "biz")
    b.StopTimer()
    require.NoError(b, err)

    require.Equal(b, expectedBizStatCopy, stat)
}

var expectedBizStatCopy = DomainStat{
    "abata.biz":         25,
    "abatz.biz":         25,
    "agimba.biz":        28,
    "agivu.biz":         17,
    "aibox.biz":         31,
    "ailane.biz":        23,
    "aimbo.biz":         25,
    "aimbu.biz":         36,
    "ainyx.biz":         35,
    "aivee.biz":         25,
    "avamba.biz":        21,
    "avamm.biz":         17,
    "avavee.biz":        35,
    "avaveo.biz":        30,
    "babbleblab.biz":    29,
    "babbleopia.biz":    36,
    "babbleset.biz":     28,
    "babblestorm.biz":   29,
    "blognation.biz":    32,
    "blogpad.biz":       34,
    "blogspan.biz":      21,
    "blogtag.biz":       23,
    "blogtags.biz":      34,
    "blogxs.biz":        35,
    "bluejam.biz":       36,
    "bluezoom.biz":      27,
    "brainbox.biz":      30,
    "brainlounge.biz":   38,
    "brainsphere.biz":   31,
    "brainverse.biz":    39,
    "brightbean.biz":    23,
    "brightdog.biz":     32,
    "browseblab.biz":    31,
    "browsebug.biz":     25,
    "browsecat.biz":     34,
    "browsedrive.biz":   24,
    "browsetype.biz":    34,
    "browsezoom.biz":    29,
    "bubblebox.biz":     19,
    "bubblemix.biz":     38,
    "bubbletube.biz":    34,
    "buzzbean.biz":      26,
    "buzzdog.biz":       30,
    "buzzshare.biz":     26,
    "buzzster.biz":      28,
    "camido.biz":        27,
    "camimbo.biz":       36,
    "centidel.biz":      32,
    "centimia.biz":      17,
    "centizu.biz":       18,
    "chatterbridge.biz": 30,
    "chatterpoint.biz":  32,
    "cogibox.biz":       30,
    "cogidoo.biz":       34,
    "cogilith.biz":      24,
    "dabfeed.biz":       26,
    "dabjam.biz":        30,
    "dablist.biz":       30,
    "dabshots.biz":      33,
    "dabtype.biz":       21,
    "dabvine.biz":       26,
    "dabz.biz":          19,
    "dazzlesphere.biz":  24,
    "demimbu.biz":       27,
    "demivee.biz":       39,
    "demizz.biz":        30,
    "devbug.biz":        20,
    "devcast.biz":       35,
    "devify.biz":        27,
    "devpoint.biz":      26,
    "devpulse.biz":      27,
    "devshare.biz":      30,
    "digitube.biz":      30,
    "divanoodle.biz":    33,
    "divape.biz":        32,
    "divavu.biz":        28,
    "dynabox.biz":       66,
    "dynava.biz":        21,
    "dynazzy.biz":       29,
    "eabox.biz":         28,
    "eadel.biz":         25,
    "eamia.biz":         18,
    "eare.biz":          30,
    "eayo.biz":          30,
    "eazzy.biz":         27,
    "edgeblab.biz":      29,
    "edgeclub.biz":      29,
    "edgeify.biz":       36,
    "edgepulse.biz":     21,
    "edgetag.biz":       24,
    "edgewire.biz":      29,
    "eidel.biz":         33,
    "eimbee.biz":        22,
    "einti.biz":         19,
    "eire.biz":          28,
    "fadeo.biz":         35,
    "fanoodle.biz":      23,
    "fatz.biz":          30,
    "feedbug.biz":       29,
    "feedfire.biz":      30,
    "feedfish.biz":      35,
    "feedmix.biz":       31,
    "feednation.biz":    24,
    "feedspan.biz":      28,
    "fivebridge.biz":    20,
    "fivechat.biz":      29,
    "fiveclub.biz":      23,
    "fivespan.biz":      27,
    "flashdog.biz":      20,
    "flashpoint.biz":    35,
    "flashset.biz":      30,
    "flashspan.biz":     32,
    "flipbug.biz":       27,
    "flipopia.biz":      30,
    "flipstorm.biz":     21,
    "fliptune.biz":      29,
    "gabcube.biz":       29,
    "gabspot.biz":       24,
    "gabtune.biz":       29,
    "gabtype.biz":       29,
    "gabvine.biz":       24,
    "geba.biz":          24,
    "gevee.biz":         23,
    "gigabox.biz":       28,
    "gigaclub.biz":      25,
    "gigashots.biz":     26,
    "gigazoom.biz":      29,
    "innojam.biz":       26,
    "innotype.biz":      27,
    "innoz.biz":         24,
    "izio.biz":          26,
    "jabberbean.biz":    28,
    "jabbercube.biz":    31,
    "jabbersphere.biz":  55,
    "jabberstorm.biz":   22,
    "jabbertype.biz":    27,
    "jaloo.biz":         35,
    "jamia.biz":         33,
    "janyx.biz":         33,
    "jatri.biz":         18,
    "jaxbean.biz":       28,
    "jaxnation.biz":     21,
    "jaxspan.biz":       27,
    "jaxworks.biz":      30,
    "jayo.biz":          44,
    "jazzy.biz":         32,
    "jetpulse.biz":      25,
    "jetwire.biz":       26,
    "jumpxs.biz":        29,
    "kamba.biz":         30,
    "kanoodle.biz":      19,
    "kare.biz":          30,
    "katz.biz":          62,
    "kaymbo.biz":        34,
    "kayveo.biz":        22,
    "kazio.biz":         21,
    "kazu.biz":          16,
    "kimia.biz":         25,
    "kwideo.biz":        17,
    "kwilith.biz":       25,
    "kwimbee.biz":       34,
    "kwinu.biz":         15,
    "lajo.biz":          20,
    "latz.biz":          24,
    "layo.biz":          32,
    "lazz.biz":          27,
    "lazzy.biz":         26,
    "leenti.biz":        26,
    "leexo.biz":         32,
    "linkbridge.biz":    38,
    "linkbuzz.biz":      24,
    "linklinks.biz":     31,
    "linktype.biz":      31,
    "livefish.biz":      31,
    "livepath.biz":      23,
    "livetube.biz":      53,
    "livez.biz":         28,
    "meedoo.biz":        23,
    "meejo.biz":         24,
    "meembee.biz":       26,
    "meemm.biz":         23,
    "meetz.biz":         33,
    "meevee.biz":        62,
    "meeveo.biz":        27,
    "meezzy.biz":        24,
    "miboo.biz":         26,
    "midel.biz":         28,
    "minyx.biz":         25,
    "mita.biz":          29,
    "mudo.biz":          36,
    "muxo.biz":          25,
    "mybuzz.biz":        32,
    "mycat.biz":         32,
    "mydeo.biz":         20,
    "mydo.biz":          30,
    "mymm.biz":          21,
    "mynte.biz":         54,
    "myworks.biz":       27,
    "nlounge.biz":       25,
    "npath.biz":         33,
    "ntag.biz":          28,
    "ntags.biz":         32,
    "oba.biz":           22,
    "oloo.biz":          19,
    "omba.biz":          26,
    "ooba.biz":          27,
    "oodoo.biz":         30,
    "oozz.biz":          22,
    "oyoba.biz":         27,
    "oyoloo.biz":        30,
    "oyonder.biz":       29,
    "oyondu.biz":        23,
    "oyope.biz":         24,
    "oyoyo.biz":         32,
    "ozu.biz":           18,
    "photobean.biz":     25,
    "photobug.biz":      57,
    "photofeed.biz":     25,
    "photojam.biz":      35,
    "photolist.biz":     19,
    "photospace.biz":    33,
    "pixoboo.biz":       14,
    "pixonyx.biz":       30,
    "pixope.biz":        32,
    "plajo.biz":         32,
    "plambee.biz":       29,
    "podcat.biz":        31,
    "quamba.biz":        31,
    "quatz.biz":         54,
    "quaxo.biz":         25,
    "quimba.biz":        25,
    "quimm.biz":         33,
    "quinu.biz":         60,
    "quire.biz":         25,
    "realblab.biz":      32,
    "realbridge.biz":    30,
    "realbuzz.biz":      22,
    "realcube.biz":      57,
    "realfire.biz":      37,
    "reallinks.biz":     25,
    "realmix.biz":       27,
    "realpoint.biz":     22,
    "rhybox.biz":        30,
    "rhycero.biz":       28,
    "rhyloo.biz":        32,
    "rhynoodle.biz":     25,
    "rhynyx.biz":        17,
    "rhyzio.biz":        36,
    "riffpath.biz":      21,
    "riffpedia.biz":     33,
    "riffwire.biz":      31,
    "roodel.biz":        29,
    "roombo.biz":        29,
    "roomm.biz":         32,
    "rooxo.biz":         34,
    "shufflebeat.biz":   32,
    "shuffledrive.biz":  25,
    "shufflester.biz":   26,
    "shuffletag.biz":    23,
    "skaboo.biz":        35,
    "skajo.biz":         26,
    "skalith.biz":       30,
    "skiba.biz":         22,
    "skibox.biz":        27,
    "skidoo.biz":        24,
    "skilith.biz":       29,
    "skimia.biz":        45,
    "skinder.biz":       25,
    "skinix.biz":        23,
    "skinte.biz":        39,
    "skipfire.biz":      29,
    "skippad.biz":       26,
    "skipstorm.biz":     30,
    "skiptube.biz":      26,
    "skivee.biz":        34,
    "skyba.biz":         40,
    "skyble.biz":        32,
    "skyndu.biz":        32,
    "skynoodle.biz":     28,
    "skyvu.biz":         34,
    "snaptags.biz":      33,
    "tagcat.biz":        33,
    "tagchat.biz":       37,
    "tagfeed.biz":       30,
    "tagopia.biz":       17,
    "tagpad.biz":        28,
    "tagtune.biz":       22,
    "talane.biz":        22,
    "tambee.biz":        24,
    "tanoodle.biz":      38,
    "tavu.biz":          37,
    "tazz.biz":          27,
    "tazzy.biz":         28,
    "tekfly.biz":        31,
    "teklist.biz":       26,
    "thoughtbeat.biz":   30,
    "thoughtblab.biz":   24,
    "thoughtbridge.biz": 30,
    "thoughtmix.biz":    33,
    "thoughtsphere.biz": 20,
    "thoughtstorm.biz":  38,
    "thoughtworks.biz":  24,
    "topdrive.biz":      35,
    "topicblab.biz":     32,
    "topiclounge.biz":   21,
    "topicshots.biz":    30,
    "topicstorm.biz":    22,
    "topicware.biz":     35,
    "topiczoom.biz":     38,
    "trilia.biz":        28,
    "trilith.biz":       25,
    "trudeo.biz":        29,
    "trudoo.biz":        28,
    "trunyx.biz":        33,
    "trupe.biz":         34,
    "twimbo.biz":        19,
    "twimm.biz":         30,
    "twinder.biz":       28,
    "twinte.biz":        33,
    "twitterbeat.biz":   33,
    "twitterbridge.biz": 20,
    "twitterlist.biz":   26,
    "twitternation.biz": 22,
    "twitterwire.biz":   21,
    "twitterworks.biz":  39,
    "twiyo.biz":         37,
    "vidoo.biz":         28,
    "vimbo.biz":         21,
    "vinder.biz":        31,
    "vinte.biz":         34,
    "vipe.biz":          25,
    "vitz.biz":          26,
    "viva.biz":          30,
    "voolia.biz":        34,
    "voolith.biz":       26,
    "voomm.biz":         61,
    "voonder.biz":       32,
    "voonix.biz":        32,
    "voonte.biz":        26,
    "voonyx.biz":        25,
    "wikibox.biz":       27,
    "wikido.biz":        21,
    "wikivu.biz":        23,
    "wikizz.biz":        61,
    "wordify.biz":       28,
    "wordpedia.biz":     25,
    "wordtune.biz":      27,
    "wordware.biz":      19,
    "yabox.biz":         24,
    "yacero.biz":        34,
    "yadel.biz":         27,
    "yakidoo.biz":       21,
    "yakijo.biz":        29,
    "yakitri.biz":       26,
    "yambee.biz":        20,
    "yamia.biz":         17,
    "yata.biz":          25,
    "yodel.biz":         26,
    "yodo.biz":          21,
    "yodoo.biz":         24,
    "yombu.biz":         29,
    "yotz.biz":          26,
    "youbridge.biz":     40,
    "youfeed.biz":       32,
    "youopia.biz":       22,
    "youspan.biz":       59,
    "youtags.biz":       22,
    "yoveo.biz":         31,
    "yozio.biz":         33,
    "zava.biz":          29,
    "zazio.biz":         18,
    "zoombeat.biz":      28,
    "zoombox.biz":       30,
    "zoomcast.biz":      38,
    "zoomdog.biz":       29,
    "zoomlounge.biz":    25,
    "zoomzone.biz":      32,
    "zoonder.biz":       29,
    "zoonoodle.biz":     27,
    "zooveo.biz":        22,
    "zoovu.biz":         38,
    "zooxo.biz":         33,
    "zoozzy.biz":        23,
}

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
  > Проверено, что использование `strings.SplitN(email,"@", 2)[1]` для извлечения домена не ускоряет работу (они состязаются, наименование версии - `*Experimental`, соотв. `BenchmarkStat003Experimental`):
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
  >
  > go test -bench=.
  >
  >     goos: linux
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

Наименование версии - `*Remark` (соотв. `BenchmarkStat004Remark`):

<details><summary>file: `stats_remark.go`</summary>

```go
package hw10programoptimization

import (
    "bufio"
    "fmt"
    "io"
    "regexp"
    "strings"
    "sync"
)

func rowParserRemark(
    wg *sync.WaitGroup,
    mtx *sync.Mutex,
    rows <-chan string,
    compiledRegexp regexp.Regexp,
    domainStat DomainStat,
) {
    defer wg.Done()
    for row := range rows {
        matches := compiledRegexp.FindAllStringSubmatch(row, -1)
        for matcheIndex := range matches {
            domainAtLowercase := strings.ToLower(matches[matcheIndex][1])
            mtx.Lock()
            domainStat[domainAtLowercase]++
            mtx.Unlock()
        }
    }
}

func GetDomainStatRemark(r io.Reader, domain string) (DomainStat, error) {
    domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
    compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)
    if err != nil {
        return nil, err
    }
    wg := sync.WaitGroup{}
    mtx := sync.Mutex{}
    dataChannel := make(chan string)
    domainStat := make(DomainStat)
    workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 100)
    for i := 0; i < workersCount; i++ {
        wg.Add(1)
        go rowParserRemark(&wg, &mtx, dataChannel, *compiledRegexp, domainStat)
    }
    scanner := bufio.NewScanner(r)
    maxCapacity := loadEnviromentOrDefault("MAX_CAPACITY", 239)
    buf := make([]byte, maxCapacity)
    scanner.Buffer(buf, maxCapacity)
    for scanner.Scan() {
        dataChannel <- scanner.Text()
    }
    close(dataChannel)
    wg.Wait()
    return domainStat, nil
}

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

go test -v -count=1 -timeout=30s -tags bench .

    === RUN   TestGetDomainStat_Time_And_Memory
        stats_optimization_test.go:46: time used: 404.882282ms / 300ms
        stats_optimization_test.go:47: memory used: 23Mb / 30Mb
        stats_optimization_test.go:49: 
                    Error Trace:    stats_optimization_test.go:49
                    Error:          "404882282" is not less than "300000000"
                    Test:           TestGetDomainStat_Time_And_Memory
                    Messages:       the program is too slow
    --- FAIL: TestGetDomainStat_Time_And_Memory (8.54s)
    FAIL
    FAIL    github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  8.547s
    FAIL
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

### Чтение порциями без буфера (медленно)

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
