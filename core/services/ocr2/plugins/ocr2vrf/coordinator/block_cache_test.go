package coordinator

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	b := NewBlockCache[int](int64(time.Second))

	assert.Equal(t, time.Second, time.Duration(b.blockEvictionWindow), "must set correct blockEvictionWindow")
}

func TestCache(t *testing.T) {
	t.Run("Happy path, no overwrites.", func(t *testing.T) {
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

		c := NewBlockCache[int](int64(100))

		// Populate cache with ordered items.
		for i, test := range tests {
			err := c.CacheItem(test.Value, test.Key, int64(i))
			assert.NoError(t, err)
			item, err := c.GetItem(test.Key)
			assert.NoError(t, err)
			assert.Equal(t, test.Value, *item)
		}

		// Ensure cache has 5 items, with the newest and oldest pointers correct.
		assert.Equal(t, 5, len(c.cache), "cache should contain 5 keys")
		assert.Equal(t, tests[0].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[4].Key, c.newestItem.itemKey)

		// Evict all items.
		c.EvictExpiredItems(int64(105))
		assert.Equal(t, 0, len(c.cache), "cache should contain 0 keys")
		assert.Nil(t, c.oldestItem)
		assert.Nil(t, c.newestItem)

		// Cache a new item.
		err := c.CacheItem(tests[0].Value, tests[0].Key, 10)
		assert.NoError(t, err)
		item, err := c.GetItem(tests[0].Key)
		assert.NoError(t, err)
		assert.Equal(t, tests[0].Value, *item)

		// Ensure cache has 1 item, with the newest and oldest pointers correct.
		assert.Equal(t, 1, len(c.cache), "cache should contain 1 key")
		assert.Equal(t, tests[0].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[0].Key, c.newestItem.itemKey)
	})

	t.Run("Happy path, override middle item.", func(t *testing.T) {
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

		c := NewBlockCache[int](int64(100))

		// Populate cache with items.
		for i, test := range tests {
			err := c.CacheItem(test.Value, test.Key, int64(i))
			assert.NoError(t, err)
			item, err := c.GetItem(test.Key)
			assert.NoError(t, err)
			assert.Equal(t, test.Value, *item)
		}

		// Ensure cache has 4 items, with the newest and oldest pointers correct.
		assert.Equal(t, 4, len(c.cache), "cache should contain 4 keys")
		assert.Equal(t, tests[0].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[1].Key, c.newestItem.itemKey)

		// Evict all but two items.
		c.EvictExpiredItems(int64(103))
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")
		assert.Equal(t, tests[3].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[1].Key, c.newestItem.itemKey)

		// Evict all but one items.
		c.EvictExpiredItems(int64(104))
		assert.Equal(t, 1, len(c.cache), "cache should contain 1 keys")
		assert.Equal(t, tests[1].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[1].Key, c.newestItem.itemKey)

		// Evict remaining item.
		c.EvictExpiredItems(int64(105))
		assert.Equal(t, 0, len(c.cache), "cache should contain 0 keys")
		assert.Nil(t, c.oldestItem)
		assert.Nil(t, c.newestItem)
	})

	t.Run("Happy path, override last item.", func(t *testing.T) {
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

		c := NewBlockCache[int](int64(100))

		// Populate cache with items.
		for i, test := range tests {
			err := c.CacheItem(test.Value, test.Key, int64(i))
			assert.NoError(t, err)
			item, err := c.GetItem(test.Key)
			assert.NoError(t, err)
			assert.Equal(t, test.Value, *item)
		}

		// Ensure cache has 4 items, with the newest and oldest pointers correct.
		assert.Equal(t, 4, len(c.cache), "cache should contain 4 keys")
		assert.Equal(t, tests[1].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[0].Key, c.newestItem.itemKey)

		// Evict all but one item.
		c.EvictExpiredItems(int64(104))
		assert.Equal(t, 1, len(c.cache), "cache should contain 1 keys")
		assert.Equal(t, tests[0].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[0].Key, c.newestItem.itemKey)

		// Cache a new item.
		err := c.CacheItem(tests[1].Value, tests[1].Key, 10)
		assert.NoError(t, err)
		item, err := c.GetItem(tests[1].Key)
		assert.NoError(t, err)
		assert.Equal(t, tests[1].Value, *item)

		// Assert correct pointers.
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")
		assert.Equal(t, tests[0].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[1].Key, c.newestItem.itemKey)

		// Replace the oldest item.
		err = c.CacheItem(tests[0].Value, tests[0].Key, 11)
		assert.NoError(t, err)
		item, err = c.GetItem(tests[0].Key)
		assert.NoError(t, err)
		assert.Equal(t, tests[0].Value, *item)

		// Assert correct pointers.
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")
		assert.Equal(t, tests[1].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[0].Key, c.newestItem.itemKey)

		// Replace the newest item.
		err = c.CacheItem(tests[0].Value, tests[0].Key, 12)
		assert.NoError(t, err)
		item, err = c.GetItem(tests[0].Key)
		assert.NoError(t, err)
		assert.Equal(t, tests[0].Value, *item)

		// Assert correct pointers.
		assert.Equal(t, 2, len(c.cache), "cache should contain 2 keys")
		assert.Equal(t, tests[1].Key, c.oldestItem.itemKey)
		assert.Equal(t, tests[0].Key, c.newestItem.itemKey)
	})

	t.Run("Error, out of order item.", func(t *testing.T) {
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

		c := NewBlockCache[int](int64(100))

		// Attempt to populate cache with out of order items.
		for i, test := range tests {
			err := c.CacheItem(test.Value, test.Key, 5-int64(i))
			if i > 0 {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			item, err := c.GetItem(test.Key)
			if i > 0 {
				assert.Error(t, err)
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.Value, *item)
			}
		}
	})
}
