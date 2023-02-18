package coordinator

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	b := NewBlockCache[int](time.Second)

	assert.Equal(t, time.Second, time.Duration(b.evictionWindow), "must set correct blockEvictionWindow")
}

func TestCache(t *testing.T) {
	t.Run("Happy path, no overwrites.", func(t *testing.T) {

		now := time.Now().UTC()

		tests := []struct {
			Key   common.Hash
			Value int
		}{
			{Key: common.HexToHash("0x0"), Value: 1},
			{Key: common.HexToHash("0x1"), Value: 2},
			{Key: common.HexToHash("0x2"), Value: 3},
			{Key: common.HexToHash("0x3"), Value: 4},
			{Key: common.HexToHash("0x4"), Value: 5},
		}

		c := NewBlockCache[int](time.Second * 100)

		// Populate cache with ordered items.
		for i, test := range tests {
			c.CacheItem(test.Value, test.Key, getSecondsAfterNow(now, i))
			item := c.GetItem(test.Key)
			assert.Equal(t, test.Value, *item)
		}

		// Ensure cache has 5 items, with the newest and oldest pointers correct.
		assert.Equal(t, 5, len(c.cache), "cache should contain 5 keys")

		// Evict all items.
		evictionTime := getSecondsAfterNow(now, 105)
		c.EvictExpiredItems(evictionTime)
		assert.Equal(t, 0, len(c.cache), "cache should contain 0 keys")

		// Cache a new item.
		c.CacheItem(tests[0].Value, tests[0].Key, getSecondsAfterNow(now, 10))
		item := c.GetItem(tests[0].Key)
		assert.Equal(t, tests[0].Value, *item)

		// Attempting a new eviction should have no effect.
		c.EvictExpiredItems(evictionTime)
		assert.Equal(t, 1, len(c.cache), "cache should contain 1 key")

		// Reduce eviction window.
		c.SetEvictonWindow(time.Second * 50)

		// Attempting a new eviction will remove the added item.
		c.EvictExpiredItems(evictionTime)
		assert.Equal(t, 0, len(c.cache), "cache should contain 0 keys")
	})

	t.Run("Happy path, override middle item.", func(t *testing.T) {

		now := time.Now().UTC()

		tests := []struct {
			Key   common.Hash
			Value int
		}{
			{Key: common.HexToHash("0x0"), Value: 1},
			{Key: common.HexToHash("0x1"), Value: 2},
			{Key: common.HexToHash("0x2"), Value: 3},
			{Key: common.HexToHash("0x3"), Value: 4},
			{Key: common.HexToHash("0x1"), Value: 5},
		}

		c := NewBlockCache[int](time.Duration(time.Second * 100))

		// Populate cache with items.
		for i, test := range tests {
			c.CacheItem(test.Value, test.Key, getSecondsAfterNow(now, i))
			item := c.GetItem(test.Key)
			assert.Equal(t, test.Value, *item)
		}

		// Ensure cache has 4 items, with the newest and oldest pointers correct.
		assert.Equal(t, 4, len(c.cache), "cache should contain 4 keys")

		// Evict all but two items.
		c.EvictExpiredItems(getSecondsAfterNow(now, 103))
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")

		// Evict all but one items.
		c.EvictExpiredItems(getSecondsAfterNow(now, 104))
		assert.Equal(t, 1, len(c.cache), "cache should contain 1 keys")

		// Evict remaining item.
		c.EvictExpiredItems(getSecondsAfterNow(now, 105))
		assert.Equal(t, 0, len(c.cache), "cache should contain 0 keys")
	})

	t.Run("Happy path, override last item.", func(t *testing.T) {

		now := time.Now().UTC()

		tests := []struct {
			Key   common.Hash
			Value int
		}{
			{Key: common.HexToHash("0x0"), Value: 1},
			{Key: common.HexToHash("0x1"), Value: 2},
			{Key: common.HexToHash("0x2"), Value: 3},
			{Key: common.HexToHash("0x3"), Value: 4},
			{Key: common.HexToHash("0x0"), Value: 5},
		}

		c := NewBlockCache[int](time.Duration(time.Second * 100))

		// Populate cache with items.
		for i, test := range tests {
			c.CacheItem(test.Value, test.Key, getSecondsAfterNow(now, i))
			item := c.GetItem(test.Key)
			assert.Equal(t, test.Value, *item)
		}

		// Ensure cache has 4 items, with the newest and oldest pointers correct.
		assert.Equal(t, 4, len(c.cache), "cache should contain 4 keys")

		// Evict all but one item.
		c.EvictExpiredItems(getSecondsAfterNow(now, 104))
		assert.Equal(t, 1, len(c.cache), "cache should contain 1 keys")

		// Cache a new item.
		c.CacheItem(tests[1].Value, tests[1].Key, getSecondsAfterNow(now, 110))
		item := c.GetItem(tests[1].Key)
		assert.Equal(t, tests[1].Value, *item)

		// Assert correct length.
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")

		// Replace the oldest item.
		c.CacheItem(tests[0].Value, tests[0].Key, getSecondsAfterNow(now, 111))
		item = c.GetItem(tests[0].Key)
		assert.Equal(t, tests[0].Value, *item)

		// Assert correct length.
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")

		// Replace the newest item.
		c.CacheItem(tests[0].Value, tests[0].Key, getSecondsAfterNow(now, 112))
		item = c.GetItem(tests[0].Key)
		assert.Equal(t, tests[0].Value, *item)

		// Assert correct length.
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")
	})
}

func getSecondsAfterNow(now time.Time, i int) time.Time {
	return now.Add(time.Duration(i) * time.Second)
}
