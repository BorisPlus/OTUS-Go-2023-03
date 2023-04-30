# Домашнее задание №4 «LRU-кэш»

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

```

{{ doc.txt }}

```

</details>

## Результаты тестирование

```shell
go test -cover ./
```

hw04_lru_cache 0.007s coverage: **96.1%** of statements

### Результаты тестирование двусвязного списка

```shell
go test -v ./list.go list_stringer.go list_test.go > list_testing.txt
```

<details>
<summary>см. "list_testing.txt"</summary>

```text

{{ list_testing.txt }}

```

</details>

### Результаты тестирование кэша

```shell
go test -v list.go list_stringer.go cache.go cache_stringer.go cache_test.go > cache_testing.txt
```

<details>
<summary>см. "cache_testing.txt"</summary>

```text

{{ cache_testing.txt }}

```

</details>

## 0(1) сложность?

Откровенно скажу, я не нашел штатной возможности через Benchmark-тест продемонстрировать оценку времени испонения шага цикла - среднее значение времени, затрачиваемое на операции кэша - Set или Get. Поэтому сам реализовал сбор данных о времени исполнения той или иной операции, а в тоге в Benchmark-тест добавил метрики по:

* среднему времени, затрачиваемому на добавление элемента в кэш (med_t/set);
* дисперсии/отклонению времени, затрачиваемому на  добавление эдемента в кэш
дисперсии (disp_t/set);
* среднему времени на извлечение элемента из кэша (med_t/get);
* дисперсии/отклонению времен на извлечение элемента из кэша (disp_t/get).

**Замечание**: В Go1.20.3 [имеется]("https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/testing/benchmark.go"):

```go
// Elapsed returns the measured elapsed time of the benchmark.
// The duration reported by Elapsed matches the one measured by
// StartTimer, StopTimer, and ResetTimer.
func (b *B) Elapsed() time.Duration {
    d := b.duration
    if b.timerOn {
        d += time.Since(b.start)
    }
    return d
}
```

```shell
go test -bench=. -count=5 -benchmem list.go cache.go cache_test.go > O?.txt
```

<details>
<summary>см. "cache_testing.txt"</summary>

```text
{{ O?.txt }}
```

</details>

**Замечание**: Данные таблицы ниже не соотвествуют отчету тестирования выше. Сведения предыдущего отчета, но характер статистики подобный.

| Данных    | Ср.вр.Set |Дисп.вр.Set| Ср.вр.Get |Дисп.вр.Get|
|:----------|:---------:|:---------:|:---------:|:---------:|
|         1 | 0.0000195 | 0         |  41836    | 0         |
|         1 | 0.0000298 | 0         |  51545    | 0         |
|         1 | 0.0000585 | 0         | 123273    | 0         |
|         1 | 0.0000242 | 0         |  49169    | 0         |
|         1 | 0.0000194 | 0         |  43093    | 0         |
|       100 | 0.0000499 | 0         |  48718    | 11188478  |
|       100 | 0.0000590 | 0         |  46080    | 45719711  |
|       100 | 0.0000555 | 0         |  45047    | 20888291  |
|       100 | 0.0000552 | 0         |  49476    | 25700855  |
|       100 | 0.0000524 | 0         |  43033    | 28605871  |
|     10000 | 0.0000990 | 0         |  71010    | 110489575006|
|     10000 | 0.0000648 | 0         |  74181    | 63383030103|
|     10000 | 0.0000572 | 0         |  54419    | 1443327760|
|     10000 | 0.0000527 | 0         |  56246    | 7718550376|
|     10000 | 0.0000541 | 0         |  54864    | 1806996526|

Как видно, скорость добавления и извлечения мз Кэша не зависит от объема исходных данных (для Set с некоторой дисперсией, неустановленной причины).
При росте данных на несколько порядков среднее время имеет рост, хотя совсем незначительный (это тоже загадка пока), но точно не линейный.

## Развитие

Cделать кэш горутино-безопасным.
