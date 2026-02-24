package observation

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

type mockStreamValue struct {
	value []byte
}

func (m *mockStreamValue) Value() any {
	return m.value
}

func (m *mockStreamValue) MarshalBinary() ([]byte, error) {
	return m.value, nil
}

func (m *mockStreamValue) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}
	m.value = data
	return nil
}

func (m *mockStreamValue) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%d", m.value), nil
}

func (m *mockStreamValue) UnmarshalText(data []byte) error {
	m.value = data
	return nil
}

func (m *mockStreamValue) Type() llo.LLOStreamValue_Type {
	return llo.LLOStreamValue_TimestampedStreamValue
}

func TestNewCache(t *testing.T) {
	tests := []struct {
		name            string
		cleanupInterval time.Duration
		wantErr         bool
	}{
		{
			name:            "valid cache with cleanup",
			cleanupInterval: time.Millisecond * 100,
			wantErr:         false,
		},
		{
			name:            "valid cache without cleanup",
			cleanupInterval: 0,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(tt.cleanupInterval)
			defer cache.Close()
			require.NotNil(t, cache)
			assert.Equal(t, tt.cleanupInterval, cache.cleanupInterval)
			assert.NotNil(t, cache.values)
			assert.NotNil(t, cache.closeChan)
			assert.NotNil(t, cache.metricsCh)
		})
	}
}

func TestCache_AddMany(t *testing.T) {
	t.Run("adds multiple values with same TTL", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		ttl := time.Second
		values := map[llotypes.StreamID]llo.StreamValue{
			1: &mockStreamValue{value: []byte{1}},
			2: &mockStreamValue{value: []byte{2}},
			3: &mockStreamValue{value: []byte{3}},
		}
		cache.AddMany(values, ttl)

		for id, want := range values {
			got, _ := cache.Get(id)
			assert.Equal(t, want, got)
		}
	})

	t.Run("empty map is a no-op", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.AddMany(map[llotypes.StreamID]llo.StreamValue{}, time.Second)
		val, _ := cache.Get(1)
		assert.Nil(t, val)
	})

	t.Run("single entry", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.AddMany(map[llotypes.StreamID]llo.StreamValue{
			42: &mockStreamValue{value: []byte{42}},
		}, time.Minute)
		got, _ := cache.Get(42)
		assert.Equal(t, &mockStreamValue{value: []byte{42}}, got)
	})

	t.Run("overwrites existing entries", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.Add(1, &mockStreamValue{value: []byte{0}}, time.Second)
		cache.AddMany(map[llotypes.StreamID]llo.StreamValue{
			1: &mockStreamValue{value: []byte{100}},
		}, time.Second)
		got, _ := cache.Get(1)
		assert.Equal(t, &mockStreamValue{value: []byte{100}}, got)
	})
}

func TestCache_UpdateStreamValues(t *testing.T) {
	t.Run("fills map with cached values", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.AddMany(map[llotypes.StreamID]llo.StreamValue{
			1: &mockStreamValue{value: []byte{1}},
			2: &mockStreamValue{value: []byte{2}},
			3: &mockStreamValue{value: []byte{3}},
		}, time.Second)

		streamValues := llo.StreamValues{1: nil, 2: nil, 3: nil}
		cache.UpdateStreamValues(streamValues)

		assert.Equal(t, &mockStreamValue{value: []byte{1}}, streamValues[1])
		assert.Equal(t, &mockStreamValue{value: []byte{2}}, streamValues[2])
		assert.Equal(t, &mockStreamValue{value: []byte{3}}, streamValues[3])
	})

	t.Run("misses remain nil", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.Add(1, &mockStreamValue{value: []byte{1}}, time.Second)

		streamValues := llo.StreamValues{1: nil, 2: nil, 99: nil}
		cache.UpdateStreamValues(streamValues)

		assert.Equal(t, &mockStreamValue{value: []byte{1}}, streamValues[1])
		assert.Nil(t, streamValues[2])
		assert.Nil(t, streamValues[99])
	})

	t.Run("empty map is a no-op", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.Add(1, &mockStreamValue{value: []byte{1}}, time.Second)
		streamValues := llo.StreamValues{}
		cache.UpdateStreamValues(streamValues)
		assert.Empty(t, streamValues)
	})

	t.Run("expired entries are filled with nil", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.Add(1, &mockStreamValue{value: []byte{1}}, time.Nanosecond*100)
		time.Sleep(time.Millisecond)

		streamValues := llo.StreamValues{1: nil}
		cache.UpdateStreamValues(streamValues)
		assert.Nil(t, streamValues[1])
	})

	t.Run("overwrites existing values in map", func(t *testing.T) {
		cache := NewCache(0)
		defer cache.Close()
		cache.Add(1, &mockStreamValue{value: []byte{100}}, time.Second)

		streamValues := llo.StreamValues{1: &mockStreamValue{value: []byte{0}}}
		cache.UpdateStreamValues(streamValues)
		assert.Equal(t, &mockStreamValue{value: []byte{100}}, streamValues[1])
	})
}

func TestCache_Add_Get(t *testing.T) {
	tests := []struct {
		name      string
		streamID  llotypes.StreamID
		value     llo.StreamValue
		ttl       time.Duration
		wantValue llo.StreamValue
		beforeGet func(cache *Cache)
	}{
		{
			name:      "get existing value",
			streamID:  1,
			value:     &mockStreamValue{value: []byte{42}},
			ttl:       time.Second,
			wantValue: &mockStreamValue{value: []byte{42}},
		},
		{
			name:      "get non-existent value",
			streamID:  1,
			ttl:       time.Second,
			wantValue: nil,
		},
		{
			name:      "get expired by age",
			streamID:  1,
			value:     &mockStreamValue{value: []byte{42}},
			ttl:       time.Nanosecond * 100,
			wantValue: nil,
			beforeGet: func(_ *Cache) {
				time.Sleep(time.Millisecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(0)
			defer cache.Close()

			if tt.value != nil {
				cache.Add(tt.streamID, tt.value, tt.ttl)
			}

			if tt.beforeGet != nil {
				tt.beforeGet(cache)
			}

			val, _ := cache.Get(tt.streamID)
			assert.Equal(t, tt.wantValue, val)
		})
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(time.Millisecond)
	defer cache.Close()
	streamID := llotypes.StreamID(1)
	value := &mockStreamValue{value: []byte{42}}

	cache.Add(streamID, value, time.Nanosecond*100)
	time.Sleep(time.Millisecond * 2)

	gotValue, _ := cache.Get(streamID)
	assert.Nil(t, gotValue)
}

func TestCache_Close(t *testing.T) {
	t.Run("double close does not panic", func(t *testing.T) {
		cache := NewCache(time.Millisecond)
		cache.Add(1, &mockStreamValue{value: []byte{1}}, time.Second)

		require.NoError(t, cache.Close())
		require.NotPanics(t, func() {
			require.NoError(t, cache.Close())
		})
	})

	t.Run("nils values on close", func(t *testing.T) {
		cache := NewCache(0)
		cache.Add(1, &mockStreamValue{value: []byte{1}}, time.Second)
		require.NoError(t, cache.Close())
		assert.Nil(t, cache.values)
	})
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(0)
	defer cache.Close()
	const numGoroutines = 10
	const numOperations = uint32(1000)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Test concurrent Add operations
	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for j := range numOperations {
				streamID := id*numOperations + j
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, time.Second)
			}
		}(i)
	}
	wg.Wait()

	// Verify all values were added correctly
	for i := range uint32(numGoroutines) {
		for j := range numOperations {
			streamID := i*numOperations + j
			val, _ := cache.Get(streamID)
			assert.Equal(t, &mockStreamValue{value: []byte{byte(i)}}, val)
		}
	}
}

func TestCache_ConcurrentReadWrite(t *testing.T) {
	cache := NewCache(0)
	defer cache.Close()
	const numGoroutines = 10
	const numOperations = uint32(1000)

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // Double for read and write goroutines

	// Start write goroutines
	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for j := range numOperations {
				streamID := id*numOperations + j
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, time.Second)
			}
		}(i)
	}

	// Start read goroutines
	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for j := range numOperations {
				streamID := id*numOperations + j
				_, _ = cache.Get(streamID)
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_ConcurrentAddGet(t *testing.T) {
	cache := NewCache(0)
	defer cache.Close()
	const numGoroutines = 10
	const numOperations = uint32(1000)

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // Double for Add and Get goroutines

	// Start Add goroutines
	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for j := range numOperations {
				streamID := id*numOperations + j
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, time.Second)
			}
		}(i)
	}

	// Start Get goroutines
	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for j := range numOperations {
				streamID := id*numOperations + j
				_, _ = cache.Get(streamID)
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_ConcurrentAddMany(t *testing.T) {
	cache := NewCache(0)
	defer cache.Close()
	const numGoroutines = 10
	const batchSize = uint32(100)
	const numBatches = uint32(10)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for b := range numBatches {
				batch := make(map[llotypes.StreamID]llo.StreamValue, batchSize)
				for j := range batchSize {
					streamID := id*numBatches*batchSize + b*batchSize + j
					batch[streamID] = &mockStreamValue{value: []byte{byte(id)}}
				}
				cache.AddMany(batch, time.Second)
			}
		}(i)
	}
	wg.Wait()

	for i := range uint32(numGoroutines) {
		for b := range numBatches {
			for j := range batchSize {
				streamID := i*numBatches*batchSize + b*batchSize + j
				val, _ := cache.Get(streamID)
				assert.Equal(t, &mockStreamValue{value: []byte{byte(i)}}, val)
			}
		}
	}
}

func TestCache_ConcurrentAddManyUpdateStreamValues(t *testing.T) {
	cache := NewCache(0)
	defer cache.Close()
	const numWriters = 5
	const numReaders = 5
	const batchSize = uint32(50)
	const numIterations = uint32(100)

	var wg sync.WaitGroup
	wg.Add(numWriters + numReaders)

	for i := range uint32(numWriters) {
		go func(id uint32) {
			defer wg.Done()
			for iter := range numIterations {
				batch := make(map[llotypes.StreamID]llo.StreamValue, batchSize)
				for j := range batchSize {
					streamID := id*batchSize + j
					batch[streamID] = &mockStreamValue{value: []byte{byte(iter)}}
				}
				cache.AddMany(batch, time.Second)
			}
		}(i)
	}

	for i := range uint32(numReaders) {
		go func(id uint32) {
			defer wg.Done()
			for range numIterations {
				sv := make(llo.StreamValues, batchSize)
				for j := range batchSize {
					sv[id*batchSize+j] = nil
				}
				cache.UpdateStreamValues(sv)
			}
		}(i)
	}

	wg.Wait()
}
