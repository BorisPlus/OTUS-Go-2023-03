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
{{ statistic.go }}
```

</details>

Тестирование

<details>
<summary>см. "statistic_test.go"</summary>

```go
{{ statistic_test.go }}
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
{{ run.go }}
```

</details>

#### Тестирование

* первые M-задач с ошибками

```shell
go test -v -run TestRunFirstMTasksErrors ./ > TestRunFirstMTasksErrors.txt
```

```text
{{ TestRunFirstMTasksErrors.txt }}
```

* все N-задач без ошибок

```shell
go test -v -run TestRunAllTasksWithoutAnyError ./ > TestRunAllTasksWithoutAnyError.txt
```

Видно, что почти в 5 (0.21...) раз ускорили

```text
{{ TestRunAllTasksWithoutAnyError.txt }}
```

* игнорирование ошибок в принципе

Видно, что число ошибочных задач заведомо случайно

> Кстати, если запускать тест несколько раз, то окажется, что распределение между 1 или 2 в rand.Intn(2) не совсем случайно.

```shell
go test -v -run TestRunWithUnlimitedErrorsCount ./ > TestRunWithUnlimitedErrorsCount.txt
```

```text
{{ TestRunWithUnlimitedErrorsCount.txt }}
```

* Снижение числа горутин, необходимых для обработки меньшего числа задач

```shell
go test -v -run TestRun4TaskWith5Gorutine ./ > TestRun4TaskWith5Gorutine.txt
```

Ключевое в

```text
{{ TestRun4TaskWith5Gorutine.txt }}
```

это **Всего подготавливалось к запуcку горутин: 4**, а не **5**
