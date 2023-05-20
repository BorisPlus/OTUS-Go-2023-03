# Домашнее задание №5 «Параллельное исполнение»

Описание [задания](./README.md) переработано в части графиков.

> **Для формирования данного отчета запустить**
>
> ```shell
> $ cd ../report_templator/
> $ go test templator.go hw05_parallel_execution_test.go
> ```

## Реализации

### Вспомогательное

Для мониторинга процессов разработан класс сбора статистики

<details>
<summary>см. "statistic.go"</summary>

```go
package hw05parallelexecution

import (
    "fmt"
    "sync"
)

// StatisticsMonitor - структура учета статистики выполнения набора задач.
type StatisticsMonitor struct {
    rwMutex *sync.RWMutex
    // Статистика в отношении горутин

    goroutinesCountLimit   uint // Лимит на общее число горутин
    goroutinesCountInit    uint // Всего подготавливалось к запуcку горутин
    startedGoroutinesCount uint // Всего было запущено горутин
    doneGoroutinesCount    uint // Всего горутин исполнилось

    // Статистика в отношении задач

    tasksCount            uint // Общее число задач
    tasksCountInit        uint // Всего подготавливалось к запуcку задач
    startedTasksCount     uint // Всего было запущено задач
    doneTasksCount        uint // Всего задач исполнилось успешно
    errorsTasksCountLimit uint // Лимит на число задач, завершившихся с ошибками
    errorsTasksCount      uint // Всего задач завершилось с ошибками
}

func (statMonitor StatisticsMonitor) String() string {
    return fmt.Sprintf(`СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: %d
        Лимит на общее число горутин (вычилено): %d
        Всего подготавливалось к запуcку горутин: %d
        Всего было запущено горутин: %d
        Всего горутин исполнилось: %d
    ЗАДАЧИ
        Общее число задач: %d
        Всего подготавливалось к запуcку задач: %d
        Всего было запущено задач: %d
        Всего задач исполнилось успешно: %d
        Лимит на число задач, завершившихся с ошибками: %d
        Всего задач завершилось с ошибками: %d`,
        // Статистика в отношении горутин
        statMonitor.GoroutinesCountLimit(false),
        statMonitor.GoroutinesCountLimit(true),
        statMonitor.GoroutinesCountInit(),
        statMonitor.StartedGoroutinesCount(),
        statMonitor.DoneGoroutinesCount(),
        // Статистика в отношении задач
        statMonitor.TasksCount(),
        statMonitor.TasksCountInit(),
        statMonitor.StartedTasksCount(),
        statMonitor.DoneTasksCount(),
        statMonitor.ErrorsTasksCountLimit(),
        statMonitor.ErrorsTasksCount(),
    )
}

// SetErrorsTasksCountLimit() - устаноавливает ограничение на число задач, завершающихся ошибками.
// При этом значение 0 не ограничивает число задач с ошибками.
func (statMonitor *StatisticsMonitor) SetErrorsTasksCountLimit(limit uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.errorsTasksCountLimit = limit
}

// ErrorsTasksCountLimit() - ограничение на число задач, завершающихся ошибками.
// При этом значение 0 не ограничивает число задач с ошибками.
func (statMonitor *StatisticsMonitor) ErrorsTasksCountLimit() uint {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    return statMonitor.errorsTasksCountLimit
}

// cleverGoroutinesCountLimit() - ограничение на число горутин "по-умному".
// И без этого отработает, но если задач меньше, чем воркеров,
// то зачем тогда все воркеры запускать.
// Воркеров столько - сколько и задач, если задач априори меньше воркеров.
func (statMonitor *StatisticsMonitor) cleverGoroutinesCountLimit() uint {
    // TODO: double lock if call from other method
    // defer statMonitor.rwMutex.Unlock()
    // statMonitor.rwMutex.Lock()
    if statMonitor.goroutinesCountLimit <= statMonitor.tasksCount {
        return statMonitor.goroutinesCountLimit
    }
    return statMonitor.tasksCount
}

// DoesErrorsLimitExceeded() - критерий останова.
func (statMonitor *StatisticsMonitor) DoesErrorsLimitExceeded() bool {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.errorsTasksCountLimit > 0 && statMonitor.errorsTasksCount >= statMonitor.errorsTasksCountLimit
}

// goroutinesCountLimit

// SetGoroutinesCountLimit() - устанавливает ограничение на число горутин.
func (statMonitor *StatisticsMonitor) SetGoroutinesCountLimit(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.goroutinesCountLimit = v
}

func (statMonitor *StatisticsMonitor) GoroutinesCountLimit(clever bool) uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    if clever {
        return statMonitor.cleverGoroutinesCountLimit()
    }
    return statMonitor.goroutinesCountLimit
}

func (statMonitor *StatisticsMonitor) IncGoroutinesCountLimit() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.goroutinesCountLimit++
}

// goroutinesCountInit

func (statMonitor *StatisticsMonitor) SetGoroutinesCountInit(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.goroutinesCountInit = v
}

func (statMonitor *StatisticsMonitor) GoroutinesCountInit() uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.goroutinesCountInit
}

func (statMonitor *StatisticsMonitor) IncGoroutinesCountInit() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.goroutinesCountInit++
}

// startedGoroutinesCount

func (statMonitor *StatisticsMonitor) SetStartedGoroutinesCount(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.startedGoroutinesCount = v
}

func (statMonitor *StatisticsMonitor) StartedGoroutinesCount() uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.startedGoroutinesCount
}

func (statMonitor *StatisticsMonitor) IncStartedGoroutinesCount() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.startedGoroutinesCount++
}

// doneGoroutinesCount

func (statMonitor *StatisticsMonitor) SetDoneGoroutinesCount(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.doneGoroutinesCount = v
}

func (statMonitor *StatisticsMonitor) DoneGoroutinesCount() uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.doneGoroutinesCount
}

func (statMonitor *StatisticsMonitor) IncDoneGoroutinesCount() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.doneGoroutinesCount++
}

// tasksCount

func (statMonitor *StatisticsMonitor) SetTasksCount(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.tasksCount = v
}

func (statMonitor *StatisticsMonitor) TasksCount() uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.tasksCount
}

func (statMonitor *StatisticsMonitor) IncTasksCount() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.tasksCount++
}

// tasksCountInit

func (statMonitor *StatisticsMonitor) SetTasksCountInit(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.tasksCountInit = v
}

func (statMonitor *StatisticsMonitor) TasksCountInit() uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.tasksCountInit
}

func (statMonitor *StatisticsMonitor) IncTasksCountInit() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.tasksCountInit++
}

// startedTasksCount

func (statMonitor *StatisticsMonitor) SetStartedTasksCount(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.startedTasksCount = v
}

func (statMonitor *StatisticsMonitor) StartedTasksCount() uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.startedTasksCount
}

func (statMonitor *StatisticsMonitor) IncStartedTasksCount() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.startedTasksCount++
}

// doneTasksCount

func (statMonitor *StatisticsMonitor) SetDoneTasksCount(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.doneTasksCount = v
}

func (statMonitor *StatisticsMonitor) DoneTasksCount() uint {
    defer statMonitor.rwMutex.RUnlock()
    statMonitor.rwMutex.RLock()
    return statMonitor.doneTasksCount
}

func (statMonitor *StatisticsMonitor) IncDoneTasksCount() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.doneTasksCount++
}

// SetErrorsTasksCount() - число задач, завершающихся ошибками.
func (statMonitor *StatisticsMonitor) SetErrorsTasksCount(v uint) {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.errorsTasksCount = v
}

func (statMonitor *StatisticsMonitor) ErrorsTasksCount() uint {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    return statMonitor.errorsTasksCount
}

func (statMonitor *StatisticsMonitor) IncErrorsTasksCount() {
    defer statMonitor.rwMutex.Unlock()
    statMonitor.rwMutex.Lock()
    statMonitor.errorsTasksCount++
}

```

</details>

Тестирование

<details>
<summary>см. "statistic_test.go"</summary>

```go
package hw05parallelexecution

import (
    "sync"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestStatisticsMonitorBasic(t *testing.T) {
    statisticsMonitor := StatisticsMonitor{rwMutex: &sync.RWMutex{}}

    var i uint
    // Статистика в отношении горутин

    statisticsMonitor.SetGoroutinesCountLimit(i)
    require.Equal(t, uint(0), statisticsMonitor.GoroutinesCountLimit(false), "GoroutinesCountLimit set/get. OK.")
    statisticsMonitor.IncGoroutinesCountLimit()
    i++
    require.Equal(t, i, statisticsMonitor.GoroutinesCountLimit(false), "GoroutinesCountLimit increment. OK.")

    i++
    statisticsMonitor.SetGoroutinesCountInit(i)
    require.Equal(t, uint(2), statisticsMonitor.GoroutinesCountInit(), "GoroutinesCountInit set/get. OK.")
    statisticsMonitor.IncGoroutinesCountInit()
    i++
    require.Equal(t, i, statisticsMonitor.GoroutinesCountInit(), "GoroutinesCountInit increment. OK.")

    i++
    statisticsMonitor.SetStartedGoroutinesCount(i)
    require.Equal(t, uint(4), statisticsMonitor.StartedGoroutinesCount(), "StartedGoroutinesCount set/get. OK.")
    statisticsMonitor.IncStartedGoroutinesCount()
    i++
    require.Equal(t, i, statisticsMonitor.StartedGoroutinesCount(), "StartedGoroutinesCount increment. OK.")

    i++
    statisticsMonitor.SetDoneGoroutinesCount(i)
    require.Equal(t, uint(6), statisticsMonitor.DoneGoroutinesCount(), "DoneGoroutinesCount set/get. OK.")
    statisticsMonitor.IncDoneGoroutinesCount()
    i++
    require.Equal(t, i, statisticsMonitor.DoneGoroutinesCount(), "DoneGoroutinesCount increment. OK.")

    // Статистика в отношении задач

    i++
    statisticsMonitor.SetTasksCount(i)
    require.Equal(t, uint(8), statisticsMonitor.TasksCount(), "TasksCount set/get. OK.")
    statisticsMonitor.IncTasksCount()
    i++
    require.Equal(t, i, statisticsMonitor.TasksCount(), "TasksCount increment. OK.")

    i++
    statisticsMonitor.SetTasksCountInit(i)
    require.Equal(t, uint(10), statisticsMonitor.TasksCountInit(), "TasksCountInit set/get. OK.")
    statisticsMonitor.IncTasksCountInit()
    i++
    require.Equal(t, i, statisticsMonitor.TasksCountInit(), "TasksCountInit increment. OK.")

    i++
    statisticsMonitor.SetStartedTasksCount(i)
    require.Equal(t, uint(12), statisticsMonitor.StartedTasksCount(), "StartedTasksCount set/get. OK.")
    statisticsMonitor.IncStartedTasksCount()
    i++
    require.Equal(t, i, statisticsMonitor.StartedTasksCount(), "StartedTasksCount increment. OK.")

    i++
    statisticsMonitor.SetDoneTasksCount(i)
    require.Equal(t, uint(14), statisticsMonitor.DoneTasksCount(), "DoneTasksCount set/get. OK.")
    statisticsMonitor.IncDoneTasksCount()
    i++
    require.Equal(t, i, statisticsMonitor.DoneTasksCount(), "DoneTasksCount increment. OK.")

    i++
    statisticsMonitor.SetErrorsTasksCount(i)
    require.Equal(t, uint(16), statisticsMonitor.ErrorsTasksCount(), "ErrorsTasksCount set/get. OK.")
    statisticsMonitor.IncErrorsTasksCount()
    i++
    require.Equal(t, i, statisticsMonitor.ErrorsTasksCount(), "ErrorsTasksCount increment. OK.")

    i++
    statisticsMonitor.SetErrorsTasksCountLimit(i)
    require.Equal(t, uint(18), statisticsMonitor.ErrorsTasksCountLimit(), "ErrorsTasksCountLimit set/get. OK.")

    // Вычисляемое значение числа горутин, если задач априори меньше

    statisticsMonitor.SetGoroutinesCountLimit(100)
    statisticsMonitor.SetTasksCount(10)
    require.Equal(t, uint(10), statisticsMonitor.GoroutinesCountLimit(true), "GoroutinesCountLimit(true) M<N. OK.")

    statisticsMonitor.SetGoroutinesCountLimit(10)
    statisticsMonitor.SetTasksCount(100)
    require.Equal(t, uint(10), statisticsMonitor.GoroutinesCountLimit(true), "GoroutinesCountLimit(true) M<N. OK.")

    // Критерий останова
    statisticsMonitor.SetErrorsTasksCountLimit(10)
    statisticsMonitor.SetErrorsTasksCount(10)
    require.True(t, statisticsMonitor.DoesErrorsLimitExceeded(), "10 of 10 = ErrorsLimitExceeded. OK.")

    statisticsMonitor.SetErrorsTasksCountLimit(10)
    statisticsMonitor.SetErrorsTasksCount(100)
    require.True(t, statisticsMonitor.DoesErrorsLimitExceeded(), "100 of 10 = ErrorsLimitExceeded. OK.")

    statisticsMonitor.SetErrorsTasksCountLimit(100)
    statisticsMonitor.SetErrorsTasksCount(10)
    require.False(t, statisticsMonitor.DoesErrorsLimitExceeded(), "10 of 100 = so Not ErrorsLimitExceeded. OK.")

    statisticsMonitor.SetErrorsTasksCountLimit(0)
    statisticsMonitor.SetErrorsTasksCount(10)
    require.False(t, statisticsMonitor.DoesErrorsLimitExceeded(), "10 of 0 = so Not ErrorsLimitExceeded. OK.")
}

```

</details>

```shell
$ go test -cover ./statistic.go ./statistic_test.go 
    ok      command-line-arguments  0.006s  coverage: 95.8% of statements
```

### Воркеры

Сохраняя заданный для реализации интерфейс функции:

<details>
<summary>см. "run.go"</summary>

```go
package hw05parallelexecution

import (
    "errors"
    "fmt"
    "sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(wtg *sync.WaitGroup, tasksChan <-chan Task, stat *StatisticsMonitor) {
    defer func(stat *StatisticsMonitor) {
        stat.IncDoneGoroutinesCount()
        wtg.Done()
    }(stat)

    stat.IncStartedGoroutinesCount()

    for task := range tasksChan {
        stat.IncTasksCountInit()
        if !stat.DoesErrorsLimitExceeded() {
            stat.IncStartedTasksCount()
            taskReturnError := task() != nil
            if taskReturnError {
                stat.IncErrorsTasksCount()
            } else {
                stat.IncDoneTasksCount()
            }
        } else {
            break
        }
    }
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workTogetherTasksCountLimit, errorsCountLimit int) error {
    mtx := sync.RWMutex{}

    stat := StatisticsMonitor{rwMutex: &mtx}

    stat.SetErrorsTasksCountLimit(uint(errorsCountLimit))
    stat.SetTasksCount(uint(len(tasks)))
    stat.SetGoroutinesCountLimit(uint(workTogetherTasksCountLimit))

    fmt.Printf("\nИСХОДНАЯ\n%s\n", stat)

    defer func() {
        fmt.Printf("\nИТОГОВАЯ\n%s\n", stat)
    }()

    tasksChan := make(chan Task, len(tasks))
    for _, task := range tasks {
        tasksChan <- task
    }
    close(tasksChan)

    var workerIndex uint
    wtGr := sync.WaitGroup{}
    for workerIndex = 1; workerIndex <= stat.GoroutinesCountLimit(true); workerIndex++ {
        stat.IncGoroutinesCountInit()

        wtGr.Add(1)
        go worker(&wtGr, tasksChan, &stat)
    }

    wtGr.Wait()

    if stat.DoesErrorsLimitExceeded() {
        fmt.Println("Errors was limit!!!")
        return ErrErrorsLimitExceeded
    }
    return nil
}

```

</details>

#### Тестирование

* первые M-задач с ошибками

```shell
go test -v -run TestRunFirstMTasksErrors ./ > TestRunFirstMTasksErrors.txt
```

```text
=== RUN   TestRunFirstMTasksErrors
=== RUN   TestRunFirstMTasksErrors/If_were_errors_in_first_M_tasks,_than_finished_not_more_N+M_tasks

ИСХОДНАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 10
        Лимит на общее число горутин (вычилено): 10
        Всего подготавливалось к запуcку горутин: 0
        Всего было запущено горутин: 0
        Всего горутин исполнилось: 0
    ЗАДАЧИ
        Общее число задач: 50
        Всего подготавливалось к запуcку задач: 0
        Всего было запущено задач: 0
        Всего задач исполнилось успешно: 0
        Лимит на число задач, завершившихся с ошибками: 23
        Всего задач завершилось с ошибками: 0
Errors was limit!!!

ИТОГОВАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 10
        Лимит на общее число горутин (вычилено): 10
        Всего подготавливалось к запуcку горутин: 10
        Всего было запущено горутин: 10
        Всего горутин исполнилось: 10
    ЗАДАЧИ
        Общее число задач: 50
        Всего подготавливалось к запуcку задач: 42
        Всего было запущено задач: 32
        Всего задач исполнилось успешно: 0
        Лимит на число задач, завершившихся с ошибками: 23
        Всего задач завершилось с ошибками: 32
--- PASS: TestRunFirstMTasksErrors (0.20s)
    --- PASS: TestRunFirstMTasksErrors/If_were_errors_in_first_M_tasks,_than_finished_not_more_N+M_tasks (0.20s)
PASS
ok      hw05parallelexecution    (cached)

```

* все N-задач без ошибок

```shell
go test -v -run TestRunAllTasksWithoutAnyError ./ > TestRunAllTasksWithoutAnyError.txt
```

Видно, что почти в 5 (0.21...) раз ускорили

```text
=== RUN   TestRunAllTasksWithoutAnyError
=== RUN   TestRunAllTasksWithoutAnyError/Tasks_without_errors

ИСХОДНАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 5
        Всего подготавливалось к запуcку горутин: 0
        Всего было запущено горутин: 0
        Всего горутин исполнилось: 0
    ЗАДАЧИ
        Общее число задач: 50
        Всего подготавливалось к запуcку задач: 0
        Всего было запущено задач: 0
        Всего задач исполнилось успешно: 0
        Лимит на число задач, завершившихся с ошибками: 1
        Всего задач завершилось с ошибками: 0

ИТОГОВАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 5
        Всего подготавливалось к запуcку горутин: 5
        Всего было запущено горутин: 5
        Всего горутин исполнилось: 5
    ЗАДАЧИ
        Общее число задач: 50
        Всего подготавливалось к запуcку задач: 50
        Всего было запущено задач: 50
        Всего задач исполнилось успешно: 50
        Лимит на число задач, завершившихся с ошибками: 1
        Всего задач завершилось с ошибками: 0
One-thread time 2557000000
Multi-thread time 542108113
[One-thread time]/[Multi-thread time] 0.21200943
--- PASS: TestRunAllTasksWithoutAnyError (0.54s)
    --- PASS: TestRunAllTasksWithoutAnyError/Tasks_without_errors (0.54s)
PASS
ok      hw05parallelexecution    0.549s

```

* игнорирование ошибок в принципе

Видно, что число ошибочных задач заведомо случайно

> Кстати, если запускать тест несколько раз, то окажется, что распределение между 1 или 2 в rand.Intn(2) не совсем случайно.

```shell
go test -v -run TestRunWithUnlimitedErrorsCount ./ > TestRunWithUnlimitedErrorsCount.txt
```

```text
=== RUN   TestRunWithUnlimitedErrorsCount
=== RUN   TestRunWithUnlimitedErrorsCount/Unlimited_errors_count

ИСХОДНАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 5
        Всего подготавливалось к запуcку горутин: 0
        Всего было запущено горутин: 0
        Всего горутин исполнилось: 0
    ЗАДАЧИ
        Общее число задач: 10
        Всего подготавливалось к запуcку задач: 0
        Всего было запущено задач: 0
        Всего задач исполнилось успешно: 0
        Лимит на число задач, завершившихся с ошибками: 0
        Всего задач завершилось с ошибками: 0

ИТОГОВАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 5
        Всего подготавливалось к запуcку горутин: 5
        Всего было запущено горутин: 5
        Всего горутин исполнилось: 5
    ЗАДАЧИ
        Общее число задач: 10
        Всего подготавливалось к запуcку задач: 10
        Всего было запущено задач: 10
        Всего задач исполнилось успешно: 6
        Лимит на число задач, завершившихся с ошибками: 0
        Всего задач завершилось с ошибками: 4

ИСХОДНАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 5
        Всего подготавливалось к запуcку горутин: 0
        Всего было запущено горутин: 0
        Всего горутин исполнилось: 0
    ЗАДАЧИ
        Общее число задач: 10
        Всего подготавливалось к запуcку задач: 0
        Всего было запущено задач: 0
        Всего задач исполнилось успешно: 0
        Лимит на число задач, завершившихся с ошибками: 18446744073709551615
        Всего задач завершилось с ошибками: 0

ИТОГОВАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 5
        Всего подготавливалось к запуcку горутин: 5
        Всего было запущено горутин: 5
        Всего горутин исполнилось: 5
    ЗАДАЧИ
        Общее число задач: 10
        Всего подготавливалось к запуcку задач: 10
        Всего было запущено задач: 10
        Всего задач исполнилось успешно: 6
        Лимит на число задач, завершившихся с ошибками: 18446744073709551615
        Всего задач завершилось с ошибками: 4
--- PASS: TestRunWithUnlimitedErrorsCount (0.26s)
    --- PASS: TestRunWithUnlimitedErrorsCount/Unlimited_errors_count (0.26s)
PASS
ok      hw05parallelexecution    0.266s

```

* Снижение числа горутин, необходимых для обработки меньшего числа задач

```shell
go test -v -run TestRun4TaskWith5Gorutine ./ > TestRun4TaskWith5Gorutine.txt
```

Ключевое в

```text
=== RUN   TestRun4TaskWith5Gorutine
=== RUN   TestRun4TaskWith5Gorutine/5_goroutines_for_4_tasks

ИСХОДНАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 4
        Всего подготавливалось к запуcку горутин: 0
        Всего было запущено горутин: 0
        Всего горутин исполнилось: 0
    ЗАДАЧИ
        Общее число задач: 4
        Всего подготавливалось к запуcку задач: 0
        Всего было запущено задач: 0
        Всего задач исполнилось успешно: 0
        Лимит на число задач, завершившихся с ошибками: 0
        Всего задач завершилось с ошибками: 0

ИТОГОВАЯ
СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: 5
        Лимит на общее число горутин (вычилено): 4
        Всего подготавливалось к запуcку горутин: 4
        Всего было запущено горутин: 4
        Всего горутин исполнилось: 4
    ЗАДАЧИ
        Общее число задач: 4
        Всего подготавливалось к запуcку задач: 4
        Всего было запущено задач: 4
        Всего задач исполнилось успешно: 4
        Лимит на число задач, завершившихся с ошибками: 0
        Всего задач завершилось с ошибками: 0
--- PASS: TestRun4TaskWith5Gorutine (0.09s)
    --- PASS: TestRun4TaskWith5Gorutine/5_goroutines_for_4_tasks (0.09s)
PASS
ok      hw05parallelexecution    0.095s

```

это **Всего подготавливалось к запуcку горутин: 4**, а не **5**
