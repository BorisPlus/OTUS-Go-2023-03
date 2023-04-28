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
	} else {
		if cache.capacity == cache.list.Len() {
			back := cache.list.Back()
			delete(cache.items, back.Data.(KeyValue).key)
			cache.list.Remove(back)
			back = nil
		}
		newItem := cache.list.PushFront(KeyValue{key, value})
		cache.items[key] = newItem
		return false
	}
}

// Get - получение элемента из кэша.
func (cache *LruCache) Get(key Key) (interface{}, bool) {
	item, exists := cache.items[key]
	if exists {
		cache.list.MoveToFront(item)
		return item.Data.(KeyValue).value, true
	} else {
		return nil, false
	}
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
