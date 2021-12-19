package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
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
		c := NewCache(3)
		for i := 0; i < 4; i++ {
			c.Set(Key(fmt.Sprintf("key%d", i)), i)
		}
		_, ok := c.Get("key0")
		require.False(t, ok)

		for i := 1; i < 4; i++ {
			val, ok := c.Get(Key(fmt.Sprintf("key%d", i)))
			require.True(t, ok)
			require.Equal(t, i, val)
		}

		val, ok := c.Get("key1")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("key2")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("key3")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("key1")
		require.True(t, ok)
		require.Equal(t, 1, val)

		// Add new item
		c.Set("key4", 4)

		// Removed old item
		_, ok = c.Get("key2")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.

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
