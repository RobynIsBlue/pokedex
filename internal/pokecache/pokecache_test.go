package pokecache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewCache(t *testing.T) {
	cases := []struct {
		input time.Duration
		expected *Cache
	}{
		{
			input: time.Millisecond * 5,
			expected: &Cache{
				Interval: time.Millisecond * 5,
				CacheMap: map[string]cacheEntry{},
			},
		},
	}
	for _, c := range cases {
		t.Run("get cache", func(t *testing.T) {
			actual := NewCache(c.input)
			require.Equal(t, actual.Interval, c.expected.Interval)
			require.Equal(t, len(actual.CacheMap), len(c.expected.CacheMap))
		})
	}
}

func TestAddGet(t *testing.T) {
	cases := []struct {
		input string
		expected []byte
	}{
		{
			input: "https://pergamino.com",
			expected: []byte("coffee"),
			},
			{
				input: "",
				expected: []byte(""),
			},
		}
	for _, c := range cases {
		t.Run("add get", func(t *testing.T) {
			cache := NewCache(time.Millisecond * 5)
			cache.Add(c.input, c.expected)
			<-time.After(time.Millisecond * 4)
			actual, _ := cache.Get(c.input)
			require.Equal(t, c.expected, actual)
		})
	}
}


func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + time.Millisecond * 6
	cache := NewCache(baseTime)
	cache.Add("https://pergamino.com", []byte("coffee"))

	got, _ := cache.Get("https://pergamino.com")
	require.Equal(t, []byte("coffee"), got)

	time.Sleep(waitTime)

	got, _ = cache.Get("https://pergamino.com")
	require.Empty(t, got)

}