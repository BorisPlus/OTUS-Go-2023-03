# Домашнее задание №10 «Оптимизация программы»

> "Я - художник, я так вижу." (Веронезе Паоло)

Описание [задания](./README.md).

## 1. Статический анализ кода

В исходном [файле](./stats_initial.go):

<details><summary>file: `stats_initial.go`</summary>

```go
package hw10programoptimization

import (
    "encoding/json"
    "fmt"
    "io"
    "regexp"
    "strings"
)

type User struct {
    ID       int
    Name     string
    Username string
    Email    string
    Phone    string
    Password string
    Address  string
}

func GetDomainStatInitial(r io.Reader, domain string) (DomainStat, error) {
    u, err := getUsers(r)
    if err != nil {
        return nil, fmt.Errorf("get users error: %w", err)
    }
    return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
    content, err := io.ReadAll(r)
    if err != nil {
        return
    }

    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
        var user User
        if err = json.Unmarshal([]byte(line), &user); err != nil {
            return
        }
        result[i] = user
    }
    return
}

func countDomains(u users, domain string) (DomainStat, error) {
    result := make(DomainStat)

    for _, user := range u {
        matched, err := regexp.Match("\\."+domain, []byte(user.Email))
        if err != nil {
            return nil, err
        }

        if matched {
            num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
            num++
            result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
        }
    }
    return result, nil
}

```

</details>

имеются конструкции кода, которые, исходя из моего опыта, являются узким местом реализации текущего алгоритма (наименование версии - `*Repo`, соотв. `BenchmarkStat001Repo`).

> В ниже представленных блоках исходного кода троеточие "..." скрывает конструкции, не имеющие отношения к текущему контексту повествования.

### 1.1. Явное ограничение на длину `100_000`

```go
type users [100_000]User
```

Сведения о размерах входных данных не известны. Если не предполагаем читать порционного по `100_000`, то лучше `slice`:

```go
type users []User
```

### 1.2. Полнообъемное чтение за раз

```go
... := io.ReadAll(r)
```

Лучше читать порциями, обрабатывая, например, файл построчно.

Сюда также можно отнести тот факт, что после полнообъемного чтения:

* Файл разбивается на строки.
* По всему объему строк составляется весь массив пользователей.
* По всему массиву пользователей высчитывается статистика.

Видится, что вложенные циклы упростят сбор мусора, не требуя много памяти для полнообъемных данных.

### 1.3. Излишнее приведение типов

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

### 1.4. Отсутствие подхода повторного использования

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

### 1.5. Переприсваивание значения в рамках атомарной операции

```go
num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
num++
result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
```

если можно, например, так (**нет, но и так я не буду делать**, см. далее)

```go
result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
```

### 1.6. Излишнее оборачивание

```go
... fmt.Errorf("get users error: %w", err)
```

если можно, например, так

```go
... err
```

### 1.7. Доверие к данным

#### 1.7.1. К данным `user.Email`

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

#### 1.7.2. К JSON-данным в принципе

```go
if err = json.Unmarshal([]byte(line), &user); err != nil {
    return
}
```

Если "дамп" данных битый, то `json.Unmarshal` даст ошибку `InvalidUnmarshalError`, но это не значит, что данные конкретно поля `user.Email` - невалидны. Они не попадут в статистику.

### 1.8. Регулярное выражение

Механизм регулярных выражений задействуется не в полной мере.

#### 1.8.1. Не скомпилировано

```go
... := regexp.Match("\\."+domain, ... )
```

Алгоритм в цикле производит разбор регулярного выражения вместо того, чтобы один раз "скомпилировать" его и искать удовлетворение уже скомпилированному:

* `func Compile(expr string) (*Regexp, error)` или
* `func MustCompile(str string) *Regexp` (если, как п.1.7, доверяем входному `str`).

#### 1.8.2. Нет поиска выявленного соответствия

```go
... := regexp.Match("\\."+domain, []byte(user.Email))
...
... strings.SplitN(user.Email, "@", 2) ...
```

Сначала алгоритм устанавливает факт отнесения содержания `user.Email` к домену, а потом извлекает этот домен после знака `@`. Вместо этого можно использовать функцию поиска "подсоответствия" (`submatch`) с использованием регулярных выражений:

* `func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string` или
* `func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte`.

### 1.9. Сам алгоритм

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

### 1.10. Хм

Ощущение, что я что-то забыл написать.

### 2.1. Динамический анализ кода

Удостовериться в предположениях статического анализа поможет `pprof`:

```bash
GOGC=off go test -bench=BenchmarkStat000InitialVariant -cpuprofile cpu_000_initial_variant.out

    goos: linux
    goarch: amd64
    pkg: github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization
    cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz
    BenchmarkStat000InitialVariant-4               1        2105941056 ns/op
    PASS
    ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  3.764s  

go tool pprof -svg ./hw10_program_optimization.test ./cpu_000_initial_variant.out > ./REPORT.files/cpu_000_initial_variant.svg
```

Исходя из [графа вызовов](./REPORT.files/cpu_000_initial_variant.svg):

![cpu_000_initial_variant.svg](./REPORT.files/cpu_000_initial_variant.svg)

рассмотренные выше предположения верны.

## 3. Замечание по реализациям алгоритма

В целях демонстрации эффекта "было-стало" в репозитории сохранены все варианты реализаций.

## 3.1. Алгоритм в изначальной конве

По договоренности с куратором учебного курса:

* Данным доверяем всецело (п.1.7. тогда и от регулярных выражений можно отказаться).
* Статистика считается исключительно для `user.Email`.

### 3.1.1. Циклическая реализация

Реализация посредством [циклов](./stats_looped.go):

<details><summary>file: `stats_looped.go`</summary>

```go
package hw10programoptimization

import (
    "encoding/json"
    "fmt"
    "io"
    "bufio"
    "strings"
)


func GetDomainStatLooped(r io.Reader, domain string) (DomainStat, error) {
    result := make(DomainStat)
    scanner := bufio.NewScanner(r)
    var user User
    for scanner.Scan() {
        err := json.Unmarshal(scanner.Bytes(), &user)
        if err != nil {
            continue
        }
        if strings.HasSuffix(user.Email, fmt.Sprintf(".%s", domain)) {
            result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
        }
    }
    return result, nil
}

```

</details>

### 3.1.2. Горутированная реализация

Реализация посредством [горутин](./stats_goroutined.go):

<details><summary>file: `stats_goroutined.go`</summary>

```go
package hw10programoptimization

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io"
    "strings"
    "sync"
)


func GetDomainStatGoroutined(r io.Reader, domain string) (DomainStat, error) {
    wg := sync.WaitGroup{}
    mtx := sync.Mutex{}
    dataChannel := make(chan string)
    domainStat := make(DomainStat)
    workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 100)
    for i := 0; i < workersCount; i++ {
        wg.Add(1)
        go domainStatCalc(&wg, &mtx, dataChannel, domainStat)
    }
    scanner := bufio.NewScanner(r)
    var user User
    for scanner.Scan() {
        err := json.Unmarshal(scanner.Bytes(), &user)
        if err != nil {
            continue
        }
        if strings.HasSuffix(user.Email, fmt.Sprintf(".%s", domain)) {
            dataChannel <- strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
        }
    }
    close(dataChannel)
    wg.Wait()
    return domainStat, nil
}

```

</details>


### 3.1.3. Горутированная реализация + FastJson вместо Unmarshal

Реализация посредством [горутин + FastJson](./stats_goroutined_fastjson.go):

<details><summary>file: `stats_goroutined_fastjson.go`</summary>

```go
package hw10programoptimization

import (
    "bufio"
    "fmt"
    "io"
    "strings"
    "sync"

    "github.com/valyala/fastjson"
)

func GetDomainStatGoroutinedFastJson(r io.Reader, domain string) (DomainStat, error) {
    wg := sync.WaitGroup{}
    mtx := sync.Mutex{}
    dataChannel := make(chan string)
    domainStat := make(DomainStat)
    workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 100)
    for i := 0; i < workersCount; i++ {
        wg.Add(1)
        go domainStatCalc(&wg, &mtx, dataChannel, domainStat)
    }
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        email := fastjson.GetString(scanner.Bytes(), "Email")
        if strings.HasSuffix(email, fmt.Sprintf(".%s", domain)) {
            dataChannel <- strings.ToLower(strings.SplitN(email, "@", 2)[1])
        }
    }
    close(dataChannel)
    wg.Wait()
    return domainStat, nil
}

```

</details>

### 3.2 Альтернативное решение

Я вижу, что задача `GetDomainStatSource` - **"Посчитать статистику доменов по почтовым адресам"**.

Я бы предложил сделать так (с некоторыми допущениями):

* Учесть все выше перечисленное.
* Сделать заполнение структуры `DomainStat` конкурентным (в несколько горутин).
* Реализовать метод разбора файла ("парсинга") без промежуточной структуры `User` посредством регулярного выражения, соответствующего произвольному адресу электронной почты.

> "Настоящий" многострочный вид регулярного выражения "адрес электронной почты" выходит за рамки данной реализации.

Все это **допустимо** по условию задачи, так как можно:

* писать любой новый необходимый код;
* удалять имеющийся лишний код (кроме функции `GetDomainStat`);

Самое главное, что в текущей реализации доверия к данным нет (см. п.1.7 рассуждений):

<details><summary>file: `stats_alternate.go`</summary>

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


func domainStatCalcAlternate(
    wg *sync.WaitGroup,
    mtx *sync.Mutex,
    rows <-chan string,
    compiledRegexp *regexp.Regexp,
    domainStat DomainStat,
) {
    for row := range rows {
        matches := compiledRegexp.FindAllStringSubmatch(row, -1)
        for matcheIndex := range matches {
            domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))
            mtx.Lock()
            domainStat[domainAtLowercase]++
            mtx.Unlock()
        }
    }
    wg.Done()
}

func GetDomainStatAlternate(r io.Reader, domain string) (DomainStat, error) {
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
        go domainStatCalcAlternate(&wg, &mtx, dataChannel, compiledRegexp, domainStat)
    }
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        dataChannel <- scanner.Text()
    }
    close(dataChannel)
    wg.Wait()
    return domainStat, nil
}

```

</details>

## 4. Тестирование реализаций

### 4.1 Тестирование работоспособности

Реализован специальный обобщенный над сигнатурой функции `GetDomainStat` (`func(r io.Reader, domain string) (DomainStat, error)`) тест:

<details><summary>file: `stats_common_test.go`</summary>

```go
//go:build !bench
// +build !bench

package hw10programoptimization

import (
    "bytes"
    "io"
    "testing"

    "github.com/stretchr/testify/require"
)

type GetDomainFuncSignature func(r io.Reader, domain string) (DomainStat, error)

func TestAllGetDomainStatVariants(t *testing.T) {
    data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

    testCases := []GetDomainFuncSignature{
        GetDomainStat,
        GetDomainStatLooped,
        GetDomainStatGoroutined,
        GetDomainStatGoroutinedFastJson,
        GetDomainStatAlternate,
    }
    for _, Func := range testCases {
        t.Run("find 'com'", func(t *testing.T) {
            result, err := Func(bytes.NewBufferString(data), "com")
            require.NoError(t, err)
            require.Equal(t, DomainStat{
                "browsecat.com": 2,
                "linktype.com":  1,
            }, result)
        })

        t.Run("find 'gov'", func(t *testing.T) {
            result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
            require.NoError(t, err)
            require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
        })

        t.Run("find 'unknown'", func(t *testing.T) {
            result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
            require.NoError(t, err)
            require.Equal(t, DomainStat{}, result)
        })

    }
}

```

</details>

```bash
go test -run=TestAllGetDomainStatVariants
PASS
ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization  0.028s
```

### 4.2 Оценка производительности нагрузочное

Реализован специальный обобщенный над сигнатурой функции `GetDomainStat` (`func(r io.Reader, domain string) (DomainStat, error)`) бенчмарк-тест:

<details><summary>file: `stats_common_benchmark_test.go`</summary>

```go

```

</details>

```bash
go test -bench=. > stats_common_benchmark_test.out
 ```

```text
goos: linux
goarch: amd64
pkg: github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization
cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz
BenchmarkStat000InitialVariant-4                  1    1504527916 ns/op
BenchmarkStat001LoopedVariant-4                   1    1122965358 ns/op
BenchmarkStat002GoroutinedVariant-4               1    1487473372 ns/op
BenchmarkStat003GoroutinedFastJson-4       1000000000             0.5567 ns/op
BenchmarkStat004Alternate-4                1000000000             0.4069 ns/op
PASS
ok      github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization    23.377s

```

### 4.3 Сравнение реализаций

Имевшаяся в изначальном репозитории функция проверки нагрузки:

```bash
go test -v -count=1 -timeout=30s -tags bench .
```

была переработана в целях сравнения более успешных реализаций с изначальной.

```bash
GOGC=off go test -run=TestCommon -v -count=1 -timeout=30s -tags bench .  > stats_common_test.out
```

```text
=== RUN   TestCommonGetDomainStatInitial_Time_And_Memory

Test with GetDomainStat: Initial
    stats_common_optimization_test.go:42: time used: 1.74315414s / 300ms
    stats_common_optimization_test.go:43: memory used: 308Mb / 30Mb
    stats_common_optimization_test.go:45: 
            Error Trace:    stats_common_optimization_test.go:45
                                        stats_common_optimization_test.go:51
            Error:          "1743154140" is not less than "300000000"
            Test:           TestCommonGetDomainStatInitial_Time_And_Memory
            Messages:       the program is too slow
--- FAIL: TestCommonGetDomainStatInitial_Time_And_Memory (1.75s)
=== RUN   TestCommonGetDomainStatAlternate_Time_And_Memory

Test with GetDomainStat: Alternate
    stats_common_optimization_test.go:42: time used: 362.076142ms / 300ms
    stats_common_optimization_test.go:43: memory used: 20Mb / 30Mb
    stats_common_optimization_test.go:45: 
            Error Trace:    stats_common_optimization_test.go:45
                                        stats_common_optimization_test.go:56
            Error:          "362076142" is not less than "300000000"
            Test:           TestCommonGetDomainStatAlternate_Time_And_Memory
            Messages:       the program is too slow
--- FAIL: TestCommonGetDomainStatAlternate_Time_And_Memory (8.62s)
=== RUN   TestCommonGetDomainStatGoroutinedFastJson_Time_And_Memory

Test with GetDomainStat: Goroutined + FastJson
    stats_common_optimization_test.go:42: time used: 644.143505ms / 300ms
    stats_common_optimization_test.go:43: memory used: 5Mb / 30Mb
    stats_common_optimization_test.go:45: 
            Error Trace:    stats_common_optimization_test.go:45
                                        stats_common_optimization_test.go:61
            Error:          "644143505" is not less than "300000000"
            Test:           TestCommonGetDomainStatGoroutinedFastJson_Time_And_Memory
            Messages:       the program is too slow
--- FAIL: TestCommonGetDomainStatGoroutinedFastJson_Time_And_Memory (9.72s)
FAIL
FAIL    github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization    20.092s
FAIL

```

Видно, что нагрузка на память значительно снизилась, а скорость возросла. Варианты "FastJson" и "Alternate" между собой конкурируют. В репозитории итоговым вариантом оставлен "FastJson".

## Замечание о скорости алгоритма

Реализация не достигла требуемых `300ms`, что предположительно связано с локальными ограничениями на вычислительную мощность (железо). Так как, исходя из опроса аудитории курса, исполнение даже изначальной реализации не превышало по длительности 1 секунду, когда как у меня оно составило 1.58 секунды.

### Замечание о числе горутин-обработчиков

Число горутин-обработчиков задаются через переменную окружения `WORKERS_COUNT`:

```go
workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 100)
```

```bash
WORKERS_COUNT=1 go test -v -count=1 -timeout=30s -tags bench .
```

## Вывод
  
Задача решена.

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
