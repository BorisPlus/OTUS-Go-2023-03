# Домашнее задание №4 «LRU-кэш»

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

```

package hw04lrucache // import "github.com/BorisPlus/OTUS-Go-2023-03/tree/master/hw04_lru_cache"


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
        Next: 0xc00002e380  >>>-----|---┐
        -------------------         |   | ссылается на...
                |                   |   |
                V                   |   |
        -------------------         |   |
        Item: 0xc00002e380  <-----------┘
        -------------------         |
        Data: 1                     | ссылается на...
        Prev: 0xc00002e3a0  >>>-----┘
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

=== RUN   TestList
=== RUN   TestList/zero-value_list-item_test

[zero-value] is:
 
-------------------
Item: 0xc0000e42a0
-------------------
Data: <nil>
Prev: 0x0
Next: 0x0
-------------------
=== RUN   TestList/list-items_referencies_test

[1] <--> [2] <--> [3]

[1] is:
 
-------------------
Item: 0xc0000e42c0
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

[1] become:
 
-------------------
Item: 0xc0000e42e0
-------------------
Data: 2
Prev: 0xc0000e42c0
Next: 0x0
-------------------

[2] is:
 
-------------------
Item: 0xc0000e42e0
-------------------
Data: 2
Prev: 0xc0000e42c0
Next: 0x0
-------------------

[2] become:
 
-------------------
Item: 0xc0000e42e0
-------------------
Data: 2
Prev: 0xc0000e42c0
Next: 0xc0000e4300
-------------------

[3] is:
 
-------------------
Item: 0xc0000e4300
-------------------
Data: 3
Prev: 0xc0000e42e0
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
Item: 0xc0000e4320
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item become:
 
-------------------
Item: 0xc0000e4320
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
Item: 0xc0000e4320
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
Item: 0xc0000e4320
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Front become:
 
-------------------
Item: 0xc0000e4320
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
Item: 0xc0000e4340
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item [back] become:
 
-------------------
Item: 0xc0000e4340
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item [front] become:
 
-------------------
Item: 0xc0000e4340
-------------------
Data: 1
Prev: 0x0
Next: 0x0
-------------------

Item [2] become:
 
-------------------
Item: 0xc0000e4360
-------------------
Data: 2
Prev: 0x0
Next: 0xc0000e4340
-------------------

Item [back] become:
 
-------------------
Item: 0xc0000e4340
-------------------
Data: 1
Prev: 0xc0000e4360
Next: 0x0
-------------------

Item [front] become:
 
-------------------
Item: 0xc0000e4360
-------------------
Data: 2
Prev: 0x0
Next: 0xc0000e4340
-------------------

List become:
 
    (nil:0x0)
	|
	V
-------------------
Item: 0xc0000e4360
-------------------
Data: 2
Prev: 0x0
Next: 0xc0000e4340
-------------------
	|
	V
-------------------
Item: 0xc0000e4340
-------------------
Data: 1
Prev: 0xc0000e4360
Next: 0x0
-------------------
	|
	V
    (nil:0x0)

Item [back] removed:
 
-------------------
Item: 0xc0000e4340
-------------------
Data: 1
Prev: 0xc0000e4360
Next: 0x0
-------------------

List become:
 
    (nil:0x0)
	|
	V
-------------------
Item: 0xc0000e4360
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
Item: 0xc0000e4360
-------------------
Data: 2
Prev: 0x0
Next: 0x0
-------------------

Item [back] removed  2 step:
 
-------------------
Item: 0xc0000e4360
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
Item: 0xc0000e4380
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
Item: 0xc0000e4380
-------------------
Data: 10
Prev: 0x0
Next: 0xc0000e43a0
-------------------
	|
	V
-------------------
Item: 0xc0000e43a0
-------------------
Data: 20
Prev: 0xc0000e4380
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
Item: 0xc0000e4380
-------------------
Data: 10
Prev: 0x0
Next: 0xc0000e43a0
-------------------
	|
	V
-------------------
Item: 0xc0000e43a0
-------------------
Data: 20
Prev: 0xc0000e4380
Next: 0xc0000e43c0
-------------------
	|
	V
-------------------
Item: 0xc0000e43c0
-------------------
Data: 30
Prev: 0xc0000e43a0
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
Item: 0xc0000e4380
-------------------
Data: 10
Prev: 0x0
Next: 0xc0000e43c0
-------------------
	|
	V
-------------------
Item: 0xc0000e43c0
-------------------
Data: 30
Prev: 0xc0000e4380
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
Item: 0xc0000e4480
-------------------
Data: 70
Prev: 0x0
Next: 0xc0000e4460
-------------------
	|
	V
-------------------
Item: 0xc0000e4460
-------------------
Data: 80
Prev: 0xc0000e4480
Next: 0xc0000e4420
-------------------
	|
	V
-------------------
Item: 0xc0000e4420
-------------------
Data: 60
Prev: 0xc0000e4460
Next: 0xc0000e43e0
-------------------
	|
	V
-------------------
Item: 0xc0000e43e0
-------------------
Data: 40
Prev: 0xc0000e4420
Next: 0xc0000e4380
-------------------
	|
	V
-------------------
Item: 0xc0000e4380
-------------------
Data: 10
Prev: 0xc0000e43e0
Next: 0xc0000e43c0
-------------------
	|
	V
-------------------
Item: 0xc0000e43c0
-------------------
Data: 30
Prev: 0xc0000e4380
Next: 0xc0000e4400
-------------------
	|
	V
-------------------
Item: 0xc0000e4400
-------------------
Data: 50
Prev: 0xc0000e43c0
Next: 0x0
-------------------
	|
	V
    (nil:0x0)
--- PASS: TestListComplex (0.00s)
    --- PASS: TestListComplex/complex_processing_test (0.00s)
PASS
ok  	command-line-arguments	0.009s


```

</details>

### Результаты тестирование кэша

```shell
go test -v list.go list_stringer.go cache.go cache_stringer.go cache_test.go > cache_testing.txt
```

<details>
<summary>см. "cache_testing.txt"</summary>

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
    cache_test.go:98: 
--- SKIP: TestCacheMultithreading (0.00s)
PASS
ok  	command-line-arguments	0.007s


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
FAIL	command-line-arguments [build failed]
FAIL

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
