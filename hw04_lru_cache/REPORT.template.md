# Домашнее задание №4. «LRU-кэш».

## Реализации

Интерфейс Lister и структуры для двусвязного списка

<details>
<summary>см. код</summary>

```go
{{ list.go }}
```

</details>

Интерфейс Cacher и структуры для кэша

<details>
<summary>см. код</summary>

```go
{{ cache.go }}
```

</details>

## Документация

<details>
<summary>см. "go doc -all ./ > doc.txt"</summary>

```text
{{ doc.txt }}
```

</details>

## Результаты тестирование

### Результаты тестирование двусвязного списка

```shell
go test -cover list.go  list_test.go 
ok      command-line-arguments  0.007s  coverage: 95.2% of statements
```

hw04_lru_cache 0.007s coverage: **95.2%** of statements

```shell
go test -v ./list.go list_stringer.go list_test.go > list_test.txt
```

<details>
<summary>см. "list_test.txt"</summary>

```text

{{ list_test.txt }}

```

</details>

### Результаты тестирование кэша

```shell
go test -v list.go cache.go cache_test_data.go cache_test.go > cache_test.txt
```

<details>
<summary>см. "cache_test.txt"</summary>

```text

{{ cache_test.txt }}

```

</details>

## Benchmark или как я 0(1) сложность предъявлял

### Тестовые данные для Benchmark

```go

{{ cache_test_data.go }}

```

### Подход Benchmark (с использованием Benchmark.ReportMetric)

В начале я не нашел простой штатной возможности через Benchmark-тест продемонстрировать оценку времени исполнения шага цикла - среднее значение времени, затрачиваемое на операции кэша (Set или Get). Поэтому реализовал сбор данных о времени исполнения той или иной операции посредством Benchmark.ReportMetric по:

* среднему времени, затрачиваемому на добавление элемента в кэш (med_t/set);
* дисперсии/отклонению времени, затрачиваемому на добавление элемента в кэш дисперсии (disp_t/set);
* среднему времени на извлечение элемента из кэша (med_t/get);
* дисперсии/отклонению времени на извлечение элемента из кэша (disp_t/get).

```go
{{ cache_benchmark_cli_test.go }}
```

```shell
go test -bench=. -count=5 -benchmem list.go cache.go cache_test_data.go cache_benchmark_cli_test.go > cache_benchmark_cli_test.txt
```

```text
{{ cache_benchmark_cli_test.txt }}
```

**Замечание**: Данные таблицы ниже возможно не соответствуют тексту отчета тестирования выше, они взяты из предыдущего запуска, характер статистики подобный.

Для наглядности измерений (среднего значения и дисперсии) фиксировал время в величине - `duration.Nanoseconds()`.

| Данных | Ср.вр.Set | Дисп.вр.Set | Ср.вр.Get | Дисп.вр.Get |
|:-------------|:-------------:|------------------:|:-------------:|------------------:|
| 10           |     49.80     |              3.76 |     47800     |            259741 |
| 10           |     48.10     |             55.29 |     48052     |            281135 |
| 10           |     48.50     |             100.7 |     47891     |           1090610 |
| 10           |     45.90     |             47.89 |     45412     |          36647782 |
| 10           |     48.60     |               191 |     43505     |           2552472 |
| 100          |     58.26     |             935.4 |     54387     |        1559880554 |
| 100          |     77.94     |             36534 |     53932     |        3146918878 |
| 100          |     62.59     |             163.5 |     78304     |       14245183309 |
| 100          |     71.80     |              1097 |     46249     |          98270456 |
| 100          |     98.01     |             15979 |     70143     |         741504218 |
| 10000        |     46.91     |             64.02 |     44899     |         117657540 |
| 10000        |     54.44     |              1539 |     63521     |        6992726481 |
| 10000        |     49.17     |             263.1 |     45740     |         142750167 |
| 10000        |     68.23     |              5691 |     55182     |        2081767470 |
| 10000        |     74.28     |              9136 |     78429     |       40837710386 |

Как видно, скорость добавления и извлечения мз Кэша не зависит от объема исходных данных (для Set с некоторой дисперсией, **неустановленной причины**). При росте данных на несколько порядков среднее время имеет рост, хотя совсем незначительный (**это тоже загадка пока**), но точно не линейный.

**Вопрос**: чем обусловлены **причины** выше?

### Запуск Benchmark (с использованием BenchmarkResult)

Потом я нашел, как достучаться до [BenchmarkResult](https://www.practical-go-lessons.com/chap-34-benchmarks#run-with-code-without-cli) и провести вышеописанные вычисления (почти все) штатными средствами Benchmark.

```go
{{ cache_benchmark_nocli_test.go }}
```

```shell
go test -v list.go cache.go cache_test_data.go cache_benchmark_nocli_test.go > cache_benchmark_nocli_test.txt
```

В отчете ниже видно, что усредненные значения ресурсов, потребляемых за итеративный шаг, постоянны, то есть независимы от них, то есть - это явно O(1).

```text
{{ cache_benchmark_nocli_test.txt }}
```

**Замечание**: Данные таблицы ниже возможно не соответствуют тексту отчета тестирования выше, они взяты из предыдущего запуска, характер статистики подобный.

| Данных       | Set (AVERAGE nanosecs.) | Get (AVERAGE nanosecs.) |
|:-------------|:-----------------------:|------------------------:|
| 10           |        3101.0000        |               2130.2000 |
| 100          |        2858.7300        |               1573.6400 |
| 10000        |        3374.5771        |               1876.0627 |

### Вывод по Benchmark

* Вариант с BenchmarkResult видится более "штатным", но нет ясности, как его получить непосредственно в Benchmark-тесте, это было бы правильнее, чем запускать его в обычном тесте.
* Вариант с Benchmark.ReportMetric дает больше гибкости в оценке поведения кода, так стало возможным посчитать не просто среднее, но и дисперсию (хотя по ней имеются вопросы выше) времени исполнения тестируемого блока кода.

Было бы интересно передавать в Benchmark функцию подсчета необходимой метрики и получать BenchmarkResult непосредственно в Benchmark-тесте.

**Замечание**: В Go1.20.3 [имеется]("https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/testing/benchmark.go"):

> ```go
> // Elapsed returns the measured elapsed time of the benchmark.
> // The duration reported by Elapsed matches the one measured by
> // StartTimer, StopTimer, and ResetTimer.
> func (b *B) Elapsed() time.Duration {
>    d := b.duration
>    if b.timerOn {
>        d += time.Since(b.start)
>    }
>    return d
> }
> ```

</details>

## Развитие

Сделать кэш горутино-безопасным.
