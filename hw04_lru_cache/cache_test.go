package hw04lrucache

import (
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

		first := c.Set("first", 1)
		require.False(t, first)

		second := c.Set("second", 2)
		require.False(t, second)

		third := c.Set("third", 3)
		require.False(t, third)

		fourth := c.Set("fourth", 4)
		require.False(t, fourth)

		_, ok := c.Get("first")
		require.False(t, ok)

		val, ok := c.Get("second")
		require.True(t, ok)
		require.Equal(t, 2, val)

		second = c.Set("second", 22)
		require.True(t, second)

		third = c.Set("third", 33)
		require.True(t, third)

		third = c.Set("third", 30)
		require.True(t, third)

		fourth = c.Set("fourth", 42)
		require.True(t, fourth)

		fivth := c.Set("fivth", 55)
		require.False(t, fivth)

		_, ok = c.Get("second")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
