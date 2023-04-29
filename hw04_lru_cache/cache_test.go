package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// go test -v list.go list_stringer.go cache.go cache_stringer.go  cache_test.go

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("small-cache test", func(t *testing.T) {
		zeroCache := NewCache(0)
		_, okZero := zeroCache.Get("0")
		require.False(t, okZero)
		okZero = zeroCache.Set("0", "zero")
		require.False(t, okZero)

		cache := NewCache(2)

		okOne := cache.Set("1", "one")
		require.False(t, okOne)

		okFirst := cache.Set("1", "first")
		require.True(t, okFirst)

		_, okOne = cache.Get("1")
		require.True(t, okOne)

		okSecond := cache.Set("2", "second")
		require.False(t, okSecond)

		okThird := cache.Set("3", "third")
		require.False(t, okThird)

		_, okOne = cache.Get("1")
		require.False(t, okOne)
	})

	t.Run("clear", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("1", "first")
		cache.Set("2", "second")
		cache.Clear()

		_, okOne := cache.Get("1")
		require.False(t, okOne)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Write me
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

// go test -bench=. -benchmem list.go cache.go cache_stringer.go  cache_test.go
// go test -bench=Prime list.go cache.go cache_stringer.go  cache_test.go
// Операции, необходимые для подсчета статистики времени добавления

// Степень двойки
func x2(x float64) float64 {
	return x * x
}

// Сумма
func sum(arr []float64) float64 {
	var sum float64 = 0
	for _, value := range arr {
		sum += value
	}
	return sum
}

// Среднее значение (времени проведения операции Set/Get)
func mediana(data []float64) float64 {
	return float64(sum(data)) / float64(len(data))
}

// Дисперсия (времени проведения операции Set/Get)
func dispersion(data []float64) float64 {
	dataMediana := mediana(data)
	var x2up = []float64{}
	for _, value := range data {
		x2up = append(x2up, x2(value-dataMediana))
	}
	return sum(x2up) / float64(len(data))
}

// Тестирование на объемах
func generateData(count int) []KeyValue {
	keyValues := []KeyValue{}

	for key, value := range rand.Perm(count) {
		keyValues = append(keyValues, KeyValue{key: Key(fmt.Sprint(key)), value: value})
	}
	return keyValues
}

// тестовые данные
var TestCases = []struct {
	data []KeyValue
}{
	{data: generateData(1)},
	{data: generateData(100)},
	{data: generateData(10000)},
}

var keyValues = []KeyValue{}

func BenchmarkSet(b *testing.B) {
	for _, testCase := range TestCases {
		b.Run(fmt.Sprintf("%d", len(testCase.data)), func(b *testing.B) {
			// чтоб не было операций вымещения capasity=valuesCount
			cache := NewCache(5)
			var durationsSet, durationsGet []float64 = []float64{}, []float64{}
			b.ResetTimer()
			// Собираем данные для статистики в отношении метода Set
			for _, keyValue := range testCase.data {

				start := time.Now()
				b.StartTimer()

				cache.Set(keyValue.key, keyValue.value)

				b.StopTimer()
				duration := time.Since(start)

				durationsSet = append(durationsSet, float64(duration.Seconds()))
			}

			// Собираем данные для статистики в отношении метода Get
			for _, keyValue := range testCase.data {
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
