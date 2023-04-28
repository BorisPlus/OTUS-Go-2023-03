package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

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
		_, ok_zero := zeroCache.Get("0")
		require.False(t, ok_zero)
		ok_zero = zeroCache.Set("0", "zero")

		cache := NewCache(2)

		ok_one := cache.Set("1", "one")
		require.False(t, ok_one)

		ok_first := cache.Set("1", "first")
		require.True(t, ok_first)

		_, ok_one = cache.Get("1")
		require.True(t, ok_one)

		ok_second := cache.Set("2", "second")
		require.False(t, ok_second)

		ok_third := cache.Set("3", "third")
		require.False(t, ok_third)
		
		_, ok_one = cache.Get("1")
		require.False(t, ok_one)
	})

	t.Run("clear", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("1", "first")
		cache.Set("2", "second")
		cache.Clear()

		_, ok_one := cache.Get("1")
		require.False(t, ok_one)
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
