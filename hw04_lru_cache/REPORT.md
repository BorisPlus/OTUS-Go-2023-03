# Домашнее задание №4. «LRU-кэш».

## Реализации

Интерфейс Lister и структуры для двусвязного списка

<details>
<summary>см. код</summary>

```go
package hw04lrucache

// Lister - интерфейс двусвязного списка.
type Lister interface {
    Len() int
    Front() *ListItem
    Back() *ListItem
    PushFront(v interface{}) *ListItem
    PushBack(v interface{}) *ListItem
    Remove(i *ListItem)
    MoveToFront(i *ListItem)
}

// ListItem - элемент двусвязного списка.
type ListItem struct {
    Data interface{}
    Prev *ListItem
    Next *ListItem
}

// List - структура двусвязного списка.
type List struct {
    len   int
    front *ListItem
    back  *ListItem
}

// Len() - получить длину двусвязного списка.
func (list *List) Len() int {
    return list.len
}

// Front() - получить первый элемент двусвязного списка.
func (list *List) Front() *ListItem {
    return list.front
}

// Back() - получить последний элемент двусвязного списка.
func (list *List) Back() *ListItem {
    return list.back
}

// PushFront() - добавить значение в начало двусвязного списка.
func (list *List) PushFront(data interface{}) *ListItem {
    item := &ListItem{
        Data: data,
        Prev: nil,
        Next: nil,
    }
    if list.len == 0 {
        list.front = item
        list.back = item
    } else {
        item.Next = list.front
        list.front.Prev = item
        list.front = item
    }
    list.len++
    return item
}

// PushBack() - добавить значение в конец двусвязного списка.
func (list *List) PushBack(data interface{}) *ListItem {
    item := &ListItem{
        Data: data,
        Prev: nil,
        Next: nil,
    }
    if list.len == 0 {
        list.front = item
        list.back = item
    } else {
        item.Prev = list.back
        list.back.Next = item
        list.back = item
    }
    list.len++
    return item
}

// Remove() - удалить элемент из двусвязного списка.
func (list *List) Remove(i *ListItem) {
    // TODO:
    // не наглядно, но интересно
    // Ai, Aj = Aj, Ai
    if i.Prev != nil {
        i.Prev.Next = i.Next
    } else {
        list.front = i.Next
    }
    if i.Next != nil {
        i.Next.Prev = i.Prev
    } else {
        list.back = i.Prev
    }
    list.len--
}

// MoveToFront() - переместить элемент в начало двусвязного списка.
func (list *List) MoveToFront(i *ListItem) {
    if i.Prev != nil {
        i.Prev.Next = i.Next
    } else {
        list.front = i.Next
    }
    if i.Next != nil {
        i.Next.Prev = i.Prev
    } else {
        list.back = i.Prev
    }

    i.Prev = nil
    i.Next = list.front
    if list.front != nil {
        list.front.Prev = i
        list.front = i
    } else {
        list.front = i
        list.back = i
    }
}

func NewList() Lister {
    return new(List)
}

```

</details>

Интерфейс Cacher и структуры для кэша

<details>
<summary>см. код</summary>

```go
package hw04lrucache

type Key string

// KeyValue - в хранилище будет учтена пара.
//
// Пара пригодится при извлесении элемента из списка и
// необходимостью поиска в карте, в частности, при
// очистке абсолютно заполненного кэша.
type KeyValue struct {
    key   Key
    value interface{}
}

// Cacher - интерфейс хранения кэша.
type Cacher interface {
    Set(key Key, value interface{}) bool
    Get(key Key) (interface{}, bool)
    Clear()
}

// LruCache - структура кэша.
type LruCache struct {
    capacity int
    list     Lister
    items    map[Key]*ListItem // Примерно как key-value database, доступ быстрый
}

// Set - уставновка элемента в кэш.
func (cache *LruCache) Set(key Key, value interface{}) bool {
    if cache.capacity == 0 {
        return false
    }
    item, exists := cache.items[key]
    if exists {
        item.Data = KeyValue{key, value}
        cache.list.MoveToFront(item)
        return true
    }
    if cache.capacity == cache.list.Len() {
        back := cache.list.Back()
        delete(cache.items, back.Data.(KeyValue).key)
        cache.list.Remove(back)
        back = nil
    }
    // newItem := cache.list.PushFront(KeyValue{key, value})
    cache.items[key] = cache.list.PushFront(KeyValue{key, value})
    return false
}

// Get - получение элемента из кэша.
func (cache *LruCache) Get(key Key) (interface{}, bool) {
    item, exists := cache.items[key]
    if exists {
        cache.list.MoveToFront(item)
        return item.Data.(KeyValue).value, true
    }
    return nil, false
}

// Clear - "очистка" кэша.
func (cache *LruCache) Clear() {
    // TODO:
    // а так эффективно?
    // cache.queue = nil
    // cache.items = nil
    *cache = *(NewCache(cache.capacity).(*LruCache))
}

// NewCache - функция-конструктор кэша.
func NewCache(capacity int) Cacher {
    return &LruCache{
        capacity: capacity,
        list:     NewList(),
        items:    make(map[Key]*ListItem, capacity),
    }
}

```

</details>

## Документация

<details>
<summary>см. "go doc -all ./ > doc.txt"</summary>

```text
package hw04lrucache // import "github.com/BorisPlus/OTUS-Go-2023-03/tree/master/hw04_lru_cache"


VARIABLES

var TestDatasets = []struct {
    data []KeyValue
}{
    {data: generateData(10)},
    {data: generateData(100)},
    {data: generateData(10000)},
}
    Тестовые данные.


TYPES

type Cacher interface {
    Set(key Key, value interface{}) bool
    Get(key Key) (interface{}, bool)
    Clear()
}
    Cacher - интерфейс хранения кэша.

func NewCache(capacity int) Cacher
    NewCache - функция-конструктор кэша.

type Key string

type KeyValue struct {
    // Has unexported fields.
}
    KeyValue - в хранилище будет учтена пара.

    Пара пригодится при извлесении элемента из списка и необходимостью поиска в
    карте, в частности, при очистке абсолютно заполненного кэша.

func (kv KeyValue) String() string
    String - наглядное представление значения KeyValue-структуры.

type List struct {
    // Has unexported fields.
}
    List - структура двусвязного списка.

func (list *List) Back() *ListItem
    Back() - получить последний элемент двусвязного списка.

func (list *List) Front() *ListItem
    Front() - получить первый элемент двусвязного списка.

func (list *List) Len() int
    Len() - получить длину двусвязного списка.

func (list *List) MoveToFront(i *ListItem)
    MoveToFront() - переместить элемент в начало двусвязного списка.

func (list *List) PushBack(data interface{}) *ListItem
    PushBack() - добавить значение в конец двусвязного списка.

func (list *List) PushFront(data interface{}) *ListItem
    PushFront() - добавить значение в начало двусвязного списка.

func (list *List) Remove(i *ListItem)
    Remove() - удалить элемент из двусвязного списка.

func (list *List) String() string
    String - наглядное представление всего двусвязного списка.

    Например,

    - пустой список:

        (nil:0x0)
            |
            V
        (nil:0x0)

    - список из двух элементов:

            (nil:0x0)
                |
                V
        -------------------
        Item: 0xc00002e3a0 <--------┐
        -------------------         |
        Data: 2                     |
        Prev: 0x0                   |
        Next: 0xc00002e380  >>>-----|---┐ Next 0xc00002e380
        -------------------         |   | ссылается на
                |                   |   | блок 0xc00002e380
                V                   |   |
        -------------------         |   |
        Item: 0xc00002e380  <-----------┘
        -------------------         | Prev 0xc00002e3a0
        Data: 1                     | ссылается на
        Prev: 0xc00002e3a0  >>>-----┘ блок 0xc00002e3a0
        Next: 0x0
        -------------------
                |
                V
            (nil:0x0)

type ListItem struct {
    Data interface{}
    Prev *ListItem
    Next *ListItem
}
    ListItem - элемент двусвязного списка.

func (listItem *ListItem) String() string
    String - наглядное представление значения элемента двусвязного списка.

    Например,

        -------------------             -------------------
        Item: 0xc00002e400              Item: 0xc00002e400
        -------------------             -------------------
        Data: 30                или     Data: 30
        Prev: 0xc00002e3c0              Prev: 0x0
        Next: 0xc00002e440              Next: 0x0
        -------------------             -------------------

type Lister interface {
    Len() int
    Front() *ListItem
    Back() *ListItem
    PushFront(v interface{}) *ListItem
    PushBack(v interface{}) *ListItem
    Remove(i *ListItem)
    MoveToFront(i *ListItem)
}
    Lister - интерфейс двусвязного списка.

func NewList() Lister

type LruCache struct {
    // Has unexported fields.
}
    LruCache - структура кэша.

func (cache *LruCache) Clear()
    Clear - "очистка" кэша.

func (cache *LruCache) Get(key Key) (interface{}, bool)
    Get - получение элемента из кэша.

func (cache *LruCache) Set(key Key, value interface{}) bool
    Set - уставновка элемента в кэш.


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

=== RUN   TestList
=== RUN   TestList/zero-value_list-item_test

[zero-value] is:
 
-------------------
Item: 0xc00002e2e0
-------------------
Data: <nil>
Prev: 0x0
Next: 0x0
-------------------
=== RUN   TestList/list-items_referencies_test

[1] <--> [2] <--> [3]

[1] is:
 
-------------------
Item: 0xc00002e300
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

[1] become:
 
-------------------
Item: 0xc00002e320
-------------------
Data: 2
Prev: 0xc00002e300
Next: 0x0
-------------------

[2] is:
 
-------------------
Item: 0xc00002e320
-------------------
Data: 2
Prev: 0xc00002e300
Next: 0x0
-------------------

[2] become:
 
-------------------
Item: 0xc00002e320
-------------------
Data: 2
Prev: 0xc00002e300
Next: 0xc00002e340
-------------------

[3] is:
 
-------------------
Item: 0xc00002e340
-------------------
Data: 3
Prev: 0xc00002e320
Next: 0x0
-------------------
first.Next.Next.Next is nil
first.Next.Next is third
third.Prev.Prev.Prev is nil
third.Prev.Prev is first
=== RUN   TestList/empty_list_test
=== RUN   TestList/little_list_test_#1

List was:
 
    (nil:0x0)
    |
    V
    (nil:0x0)

Item was:
 
-------------------
Item: 0xc00002e360
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item become:
 
-------------------
Item: 0xc00002e360
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

List become:
 
    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e360
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------
    |
    V
    (nil:0x0)

Back become:
 
-------------------
Item: 0xc00002e360
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Front become:
 
-------------------
Item: 0xc00002e360
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------
=== RUN   TestList/little_list_test_#2

List was:
 
    (nil:0x0)
    |
    V
    (nil:0x0)

Item [1] become:
 
-------------------
Item: 0xc00002e380
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item [back] become:
 
-------------------
Item: 0xc00002e380
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item [front] become:
 
-------------------
Item: 0xc00002e380
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item [2] become:
 
-------------------
Item: 0xc00002e3a0
-------------------
Data: 2
Prev: 0x0
Next: 0xc00002e380
-------------------

Item [back] become:
 
-------------------
Item: 0xc00002e380
-------------------
Data: 1
Prev: 0xc00002e3a0
Next: 0x0
-------------------

Item [front] become:
 
-------------------
Item: 0xc00002e3a0
-------------------
Data: 2
Prev: 0x0
Next: 0xc00002e380
-------------------

List become:
 
    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e3a0
-------------------
Data: 2
Prev: 0x0
Next: 0xc00002e380
-------------------
    |
    V
-------------------
Item: 0xc00002e380
-------------------
Data: 1
Prev: 0xc00002e3a0
Next: 0x0
-------------------
    |
    V
    (nil:0x0)

Item [back] removed:
 
-------------------
Item: 0xc00002e380
-------------------
Data: 1
Prev: 0xc00002e3a0
Next: 0x0
-------------------

List become:
 
    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e3a0
-------------------
Data: 2
Prev: 0x0
Next: 0x0
-------------------
    |
    V
    (nil:0x0)

Item [back] become 2 step:
 
-------------------
Item: 0xc00002e3a0
-------------------
Data: 2
Prev: 0x0
Next: 0x0
-------------------

Item [back] removed  2 step:
 
-------------------
Item: 0xc00002e3a0
-------------------
Data: 2
Prev: 0x0
Next: 0x0
-------------------

List become:
 
    (nil:0x0)
    |
    V
    (nil:0x0)
--- PASS: TestList (0.00s)
    --- PASS: TestList/zero-value_list-item_test (0.00s)
    --- PASS: TestList/list-items_referencies_test (0.00s)
    --- PASS: TestList/empty_list_test (0.00s)
    --- PASS: TestList/little_list_test_#1 (0.00s)
    --- PASS: TestList/little_list_test_#2 (0.00s)
=== RUN   TestListComplex
=== RUN   TestListComplex/complex_processing_test

[10] become:
 
    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e3c0
-------------------
Data: 10
Prev: 0x0
Next: 0x0
-------------------
    |
    V
    (nil:0x0)

[10, 20] become:
 
    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e3c0
-------------------
Data: 10
Prev: 0x0
Next: 0xc00002e3e0
-------------------
    |
    V
-------------------
Item: 0xc00002e3e0
-------------------
Data: 20
Prev: 0xc00002e3c0
Next: 0x0
-------------------
    |
    V
    (nil:0x0)

[10, 20, 30] become:
 
    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e3c0
-------------------
Data: 10
Prev: 0x0
Next: 0xc00002e3e0
-------------------
    |
    V
-------------------
Item: 0xc00002e3e0
-------------------
Data: 20
Prev: 0xc00002e3c0
Next: 0xc00002e400
-------------------
    |
    V
-------------------
Item: 0xc00002e400
-------------------
Data: 30
Prev: 0xc00002e3e0
Next: 0x0
-------------------
    |
    V
    (nil:0x0)
middle.Value is 20. OK
[10, 30] become:
 
    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e3c0
-------------------
Data: 10
Prev: 0x0
Next: 0xc00002e400
-------------------
    |
    V
-------------------
Item: 0xc00002e400
-------------------
Data: 30
Prev: 0xc00002e3c0
Next: 0x0
-------------------
    |
    V
    (nil:0x0)
[10, 30] mix [40, 50, 60, 70, 80] with mod(2, idx)
list.Len() is 7. OK
list.Front().Value is 80. OK
list.Back().Value is 70. OK
Forward stroke check for [80, 60, 40, 10, 30, 50, 70]. OK
Reverse stroke check for [80, 60, 40, 10, 30, 50, 70]. OK
Move front to front check for [80, 60, 40, 10, 30, 50, 70]. OK
Remove and PushBack last - check for [80, 60, 40, 10, 30, 50, 70]. OK
Check for list.Front().Prev and list.Back().Next is nils. OK
Move back to front check for [80, 60, 40, 10, 30, 50, 70]. OK
list become

    (nil:0x0)
    |
    V
-------------------
Item: 0xc00002e4c0
-------------------
Data: 70
Prev: 0x0
Next: 0xc00002e4a0
-------------------
    |
    V
-------------------
Item: 0xc00002e4a0
-------------------
Data: 80
Prev: 0xc00002e4c0
Next: 0xc00002e460
-------------------
    |
    V
-------------------
Item: 0xc00002e460
-------------------
Data: 60
Prev: 0xc00002e4a0
Next: 0xc00002e420
-------------------
    |
    V
-------------------
Item: 0xc00002e420
-------------------
Data: 40
Prev: 0xc00002e460
Next: 0xc00002e3c0
-------------------
    |
    V
-------------------
Item: 0xc00002e3c0
-------------------
Data: 10
Prev: 0xc00002e420
Next: 0xc00002e400
-------------------
    |
    V
-------------------
Item: 0xc00002e400
-------------------
Data: 30
Prev: 0xc00002e3c0
Next: 0xc00002e440
-------------------
    |
    V
-------------------
Item: 0xc00002e440
-------------------
Data: 50
Prev: 0xc00002e400
Next: 0x0
-------------------
    |
    V
    (nil:0x0)
--- PASS: TestListComplex (0.00s)
    --- PASS: TestListComplex/complex_processing_test (0.00s)
PASS
ok      command-line-arguments    0.005s


```

</details>

### Результаты тестирование кэша

```shell
go test -v list.go cache.go cache_test_data.go cache_test.go > cache_test.txt
```

<details>
<summary>см. "cache_test.txt"</summary>

```text

=== RUN   TestCache
=== RUN   TestCache/empty_cache
=== RUN   TestCache/small-cache_test
=== RUN   TestCache/clear
=== RUN   TestCache/simple
=== RUN   TestCache/purge_logic
--- PASS: TestCache (0.00s)
    --- PASS: TestCache/empty_cache (0.00s)
    --- PASS: TestCache/small-cache_test (0.00s)
    --- PASS: TestCache/clear (0.00s)
    --- PASS: TestCache/simple (0.00s)
    --- PASS: TestCache/purge_logic (0.00s)
=== RUN   TestCacheMultithreading
    cache_test.go:100: 
--- SKIP: TestCacheMultithreading (0.00s)
PASS
ok      command-line-arguments    0.035s


```

</details>

## Benchmark или как я 0(1) сложность предъявлял

### Тестовые данные для Benchmark

```go

package hw04lrucache

import (
    "fmt"
    "math/rand"
)

func generateData(count int) []KeyValue {
    keyValues := []KeyValue{}

    for key, value := range rand.Perm(count) {
        keyValues = append(keyValues, KeyValue{key: Key(fmt.Sprint(key)), value: value})
    }
    return keyValues
}

var TestDatasets = []struct {
    data []KeyValue
}{
    {data: generateData(10)},
    {data: generateData(100)},
    {data: generateData(10000)},
}


```

### Подход Benchmark (с использованием Benchmark.ReportMetric)

В начале я не нашел простой штатной возможности через Benchmark-тест продемонстрировать оценку времени исполнения шага цикла - среднее значение времени, затрачиваемое на операции кэша (Set или Get). Поэтому реализовал сбор данных о времени исполнения той или иной операции посредством Benchmark.ReportMetric по:

* среднему времени, затрачиваемому на добавление элемента в кэш (med_t/set);
* дисперсии/отклонению времени, затрачиваемому на добавление элемента в кэш дисперсии (disp_t/set);
* среднему времени на извлечение элемента из кэша (med_t/get);
* дисперсии/отклонению времени на извлечение элемента из кэша (disp_t/get).

```go
package hw04lrucache

import (
    "fmt"
    "testing"
    "time"
)

// go test -bench=. -benchmem list.go cache.go cache_test_data.go cache_benchmark_cli_test.go \
// > cache_benchmark_cli_test.txt

// Операции, необходимые для подсчета статистики времени добавления
// x2 - возведение в степень 2.
func x2(x float64) float64 {
    return x * x
}

// sum - сумма.
func sum(arr []float64) float64 {
    var sum float64
    for _, value := range arr {
        sum += value
    }
    return sum
}

// mediana - среднее значение (времени проведения операции Set/Get).
func mediana(data []float64) float64 {
    return sum(data) / float64(len(data))
}

// dispersion - дисперсия (времени проведения операции Set/Get).
func dispersion(data []float64) float64 {
    dataMediana := mediana(data)
    x2up := make([]float64, 0, len(data))
    for _, value := range data {
        x2up = append(x2up, x2(value-dataMediana))
    }
    return sum(x2up) / float64(len(data))
}

func BenchmarkSet(b *testing.B) {
    for _, testDataset := range TestDatasets {
        b.Run(fmt.Sprintf("%d", len(testDataset.data)), func(b *testing.B) {
            // чтоб не было операций вымещения capasity=valuesCount
            cache := NewCache(len(testDataset.data))
            durationsSet, durationsGet := []float64{}, []float64{}
            b.ResetTimer()
            // Собираем данные для статистики в отношении метода Set
            for _, keyValue := range testDataset.data {
                start := time.Now()
                b.StartTimer()

                cache.Set(keyValue.key, keyValue.value)

                b.StopTimer()
                duration := time.Since(start)

                durationsSet = append(durationsSet, float64(duration.Microseconds()))
            }

            // Собираем данные для статистики в отношении метода Get
            for _, keyValue := range testDataset.data {
                start := time.Now()
                b.StartTimer()

                cache.Get(keyValue.key)

                b.StopTimer()
                duration := time.Since(start)

                // durationsGet = append(durationsGet, float64(duration.Microseconds()))
                durationsGet = append(durationsGet, float64(duration.Nanoseconds()))
            }

            // Среднее времени добавления в LRU
            b.ReportMetric(mediana(durationsSet), "med_t/set")
            // Дисперсия времени добавления в LRU
            b.ReportMetric(dispersion(durationsSet), "disp_t/set")
            // Среднее времени взятия из LRU
            b.ReportMetric(mediana(durationsGet), "med_t/get")
            // Дисперсия времени взятия из LRU
            b.ReportMetric(dispersion(durationsGet), "disp_t/get")
        })
    }
}

```

```shell
go test -bench=. -count=5 -benchmem list.go cache.go cache_test_data.go cache_benchmark_cli_test.go > cache_benchmark_cli_test.txt
```

```text
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i3-2310M CPU @ 2.10GHz
BenchmarkSet/10-4         	1000000000	         0.0000559 ns/op	    259741 disp_t/get	         3.760 disp_t/set	     47800 med_t/get	        49.80 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10-4         	1000000000	         0.0000604 ns/op	    281135 disp_t/get	        55.29 disp_t/set	     48052 med_t/get	        48.10 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10-4         	1000000000	         0.0000574 ns/op	   1090610 disp_t/get	       100.7 disp_t/set	     47891 med_t/get	        48.50 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10-4         	1000000000	         0.0000487 ns/op	  36647782 disp_t/get	        47.89 disp_t/set	     45412 med_t/get	        45.90 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10-4         	1000000000	         0.0000538 ns/op	   2552472 disp_t/get	       191.0 disp_t/set	     43505 med_t/get	        48.60 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/100-4        	1000000000	         0.0005783 ns/op	1559880554 disp_t/get	       935.4 disp_t/set	     54387 med_t/get	        58.26 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/100-4        	1000000000	         0.0005797 ns/op	3146918878 disp_t/get	     36534 disp_t/set	     53932 med_t/get	        77.94 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/100-4        	1000000000	         0.0006964 ns/op	14245183309 disp_t/get	       163.5 disp_t/set	     78304 med_t/get	        62.59 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/100-4        	1000000000	         0.0006041 ns/op	  98270456 disp_t/get	      1097 disp_t/set	     46249 med_t/get	        71.80 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/100-4        	1000000000	         0.0007832 ns/op	 741504218 disp_t/get	     15979 disp_t/set	     70143 med_t/get	        98.01 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10000-4      	1000000000	         0.04917 ns/op	 117657540 disp_t/get	        64.02 disp_t/set	     44899 med_t/get	        46.91 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10000-4      	1000000000	         0.06417 ns/op	6992726481 disp_t/get	      1539 disp_t/set	     63521 med_t/get	        54.44 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10000-4      	1000000000	         0.05047 ns/op	 142750167 disp_t/get	       263.1 disp_t/set	     45740 med_t/get	        49.17 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10000-4      	1000000000	         0.06231 ns/op	2081767470 disp_t/get	      5691 disp_t/set	     55182 med_t/get	        68.23 med_t/set	       0 B/op	       0 allocs/op
BenchmarkSet/10000-4      	1000000000	         0.07280 ns/op	40837710386 disp_t/get	      9136 disp_t/set	     78429 med_t/get	        74.28 med_t/set	       0 B/op	       0 allocs/op
PASS
ok  	command-line-arguments	55.231s

```

**Замечание**: Данные таблицы ниже возможно не соответствуют тексту отчета тестирования выше, они взяты из предыдущего запуска, характер статистики подобный.

Для наглядности измерений фиксировал в величине - `duration.Nanoseconds()`.

| Данных    | Ср.вр.Set |Дисп.вр.Set| Ср.вр.Get |Дисп.вр.Get|
|:----------|:---------:|:---------:|:---------:|----------:|
|        10 | 49.80     | 3.760     |  47800    | 259741    |
|        10 | 48.10     | 55.29     |  48052    | 281135    |
|        10 | 48.50     | 100.7     |  47891    | 1090610   |
|        10 | 45.90     | 47.89     |  45412    | 36647782  |
|        10 | 48.60     | 191.0     |  43505    | 2552472   |
|       100 | 58.26     | 935.4     |  54387    | 1559880554|
|       100 | 77.94     | 36534     |  53932    | 3146918878|
|       100 | 62.59     | 163.5     |  78304    |14245183309|
|       100 | 71.80     | 1097      |  46249    | 98270456  |
|       100 | 98.01     | 15979     |  70143    | 741504218 |
|     10000 | 46.91     | 64.02     |  44899    | 117657540 |
|     10000 | 54.44     | 1539      |  63521    | 6992726481|
|     10000 | 49.17     | 263.1     |  45740    | 142750167 |
|     10000 | 68.23     | 5691      |  55182    | 2081767470|
|     10000 | 74.28     | 9136      |  78429    |40837710386|

Как видно, скорость добавления и извлечения мз Кэша не зависит от объема исходных данных (для Set с некоторой дисперсией, **неустановленной причины**). При росте данных на несколько порядков среднее время имеет рост, хотя совсем незначительный (**это тоже загадка пока**), но точно не линейный.

**Вопрос**: чем обусловлены **причины** выше?

### Запуск Benchmark (с использованием BenchmarkResult)

Потом я нашел, как достучаться до [BenchmarkResult](https://www.practical-go-lessons.com/chap-34-benchmarks#run-with-code-without-cli) и провести вышеописанные вычисления (почти все) штатными средствами Benchmark.

```go
package hw04lrucache

import (
    "fmt"
    "testing"
)

// go test -v list.go cache.go cache_test_data.go cache_benchmark_nocli_test.go > cache_benchmark_nocli_test.txt

func TestBenchmark(t *testing.T) {
    _ = t
    for _, testDataset := range TestDatasets {
        dataset := testDataset.data
        cache := NewCache(len(dataset))

        resForSet := testing.Benchmark(
            func(b *testing.B) {
                b.Helper()
                // Вынести в отдельную функцию это никак из-за необходимости сигнатуры f(b *B).
                // Вынужденное анонимное замыкание на dataset.
                b.ResetTimer()
                for _, keyValue := range dataset {
                    b.StartTimer()
                    cache.Set(keyValue.key, keyValue.value)
                    b.StopTimer()
                }
            },
        )
        fmt.Printf("--------------------------------------------------------\n")
        fmt.Printf("Operation - Set - with dataset %d values\n", len(dataset))
        fmt.Printf("--------------------------------------------------------\n")
        fmt.Printf("Number of run: %d\n", resForSet.N)
        fmt.Printf("Memory allocations: %d\n", resForSet.MemAllocs)
        fmt.Printf("Memory allocations (AVERAGE): %f\n", float64(resForSet.MemAllocs)/float64(len(dataset)))
        fmt.Printf("Number of bytes allocated: %d\n", resForSet.Bytes)
        fmt.Printf("Number of bytes allocated (AVERAGE): %f\n", float64(resForSet.Bytes)/float64(len(dataset)))
        fmt.Printf("Time taken: %s\n", resForSet.T)
        fmt.Printf("Time taken (AVERAGE, nanosecs.): %f  \n", float64(resForSet.T.Nanoseconds())/float64(len(dataset)))
        fmt.Printf("\n\n")
        res := testing.Benchmark(
            func(b *testing.B) {
                b.Helper()
                // Вынести в отдельную функцию это никак из-за необходимости сигнатуры f(b *B).
                // Вынужденное анонимное замыкание на dataset.
                b.ResetTimer()
                for _, keyValue := range dataset {
                    b.StartTimer()
                    cache.Get(keyValue.key)
                    b.StopTimer()
                }
            },
        )
        fmt.Printf("--------------------------------------------------------\n")
        fmt.Printf("Operation - Get - with dataset %d values\n", len(dataset))
        fmt.Printf("--------------------------------------------------------\n")
        fmt.Printf("Number of run: %d\n", res.N)
        fmt.Printf("Memory allocations: %d\n", res.MemAllocs)
        fmt.Printf("Memory allocations (AVERAGE): %f\n", float64(res.MemAllocs)/float64(len(dataset)))
        fmt.Printf("Number of bytes allocated: %d\n", res.Bytes)
        fmt.Printf("Number of bytes allocated (AVERAGE): %f\n", float64(res.Bytes)/float64(len(dataset)))
        fmt.Printf("Time taken: %s\n", res.T)
        fmt.Printf("Time taken (AVERAGE, nanosecs.): %f  \n", float64(res.T.Nanoseconds())/float64(len(dataset)))
        fmt.Printf("\n\n")
    }
}

```

```shell
go test -v list.go cache.go cache_test_data.go cache_benchmark_nocli_test.go > cache_benchmark_nocli_test.txt
```

В отчете ниже видно, что усредненные значения ресурсов, потребляемых за итеративный шаг, постоянны, то есть независимы от них, то есть - это явно O(1).

```text
=== RUN   TestBenchmark
--------------------------------------------------------
Operation - Set - with dataset 10 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 10
Memory allocations (AVERAGE): 1.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 31.01µs
Time taken (AVERAGE, nanosecs.): 3101.000000  


--------------------------------------------------------
Operation - Get - with dataset 10 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 0
Memory allocations (AVERAGE): 0.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 21.302µs
Time taken (AVERAGE, nanosecs.): 2130.200000  


--------------------------------------------------------
Operation - Set - with dataset 100 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 100
Memory allocations (AVERAGE): 1.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 285.873µs
Time taken (AVERAGE, nanosecs.): 2858.730000  


--------------------------------------------------------
Operation - Get - with dataset 100 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 0
Memory allocations (AVERAGE): 0.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 157.364µs
Time taken (AVERAGE, nanosecs.): 1573.640000  


--------------------------------------------------------
Operation - Set - with dataset 10000 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 10000
Memory allocations (AVERAGE): 1.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 33.745771ms
Time taken (AVERAGE, nanosecs.): 3374.577100  


--------------------------------------------------------
Operation - Get - with dataset 10000 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 0
Memory allocations (AVERAGE): 0.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 18.760627ms
Time taken (AVERAGE, nanosecs.): 1876.062700  


--- PASS: TestBenchmark (8.29s)
PASS
ok      command-line-arguments    8.302s

```

**Замечание**: Данные таблицы ниже возможно не соответствуют тексту отчета тестирования выше, они взяты из предыдущего запуска, характер статистики подобный.

| Данных    | Set (AVERAGE nanosecs.) |Get (AVERAGE nanosecs.)| 
|:----------|:---------:|----------:|
|        10 | 3101.0000 | 2130.2000 |
|       100 | 2858.7300 | 1573.6400 |
|     10000 | 3374.5771 | 1876.0627 |

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
