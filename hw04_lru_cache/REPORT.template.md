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

```text

{{ doc.txt }}

```

</details>

## Результаты тестирование

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

## Развитие

Cделать кэш горутино-безопасным.