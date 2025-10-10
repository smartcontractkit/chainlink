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
		maxAge          time.Duration
		cleanupInterval time.Duration
		wantErr         bool
	}{
		{
			name:            "valid cache with cleanup",
			maxAge:          time.Second,
			cleanupInterval: time.Millisecond * 100,
			wantErr:         false,
		},
		{
			name:            "valid cache without cleanup",
			maxAge:          time.Second,
			cleanupInterval: 0,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(tt.maxAge, tt.cleanupInterval)
			require.NotNil(t, cache)
			assert.Equal(t, tt.maxAge, cache.maxAge)
			assert.Equal(t, tt.cleanupInterval, cache.cleanupInterval)
			assert.NotNil(t, cache.values)
			assert.NotNil(t, cache.closeChan)
		})
	}
}

func TestCache_Add_Get(t *testing.T) {
	tests := []struct {
		name      string
		streamID  llotypes.StreamID
		value     llo.StreamValue
		seqNr     uint64
		maxAge    time.Duration
		wantValue llo.StreamValue
		wantFound bool
		beforeGet func(cache *Cache)
	}{
		{
			name:      "get existing value",
			streamID:  1,
			value:     &mockStreamValue{value: []byte{42}},
			seqNr:     10,
			maxAge:    time.Second,
			wantValue: &mockStreamValue{value: []byte{42}},
			wantFound: true,
		},
		{
			name:      "get non-existent value",
			streamID:  1,
			value:     &mockStreamValue{value: []byte{42}},
			seqNr:     10,
			maxAge:    time.Second,
			wantValue: nil,
			wantFound: false,
		},
		{
			name:      "get expired by sequence number",
			streamID:  1,
			value:     &mockStreamValue{value: []byte{42}},
			seqNr:     5,
			maxAge:    time.Second,
			wantValue: nil,
			wantFound: false,
			beforeGet: func(cache *Cache) {
				cache.SetLastTransmissionSeqNr(10)
			},
		},
		{
			name:      "get expired by age",
			streamID:  1,
			value:     &mockStreamValue{value: []byte{42}},
			seqNr:     10,
			maxAge:    time.Nanosecond * 100,
			wantValue: nil,
			wantFound: false,
			beforeGet: func(_ *Cache) {
				time.Sleep(time.Millisecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(tt.maxAge, 0)

			if tt.wantFound {
				cache.Add(tt.streamID, tt.value, tt.seqNr)
			}

			if tt.beforeGet != nil {
				tt.beforeGet(cache)
			}

			assert.Equal(t, tt.wantValue, cache.Get(tt.streamID))
		})
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(time.Nanosecond*100, time.Millisecond)
	streamID := llotypes.StreamID(1)
	value := &mockStreamValue{value: []byte{42}}

	cache.Add(streamID, value, 10)
	time.Sleep(time.Millisecond * 2)
	assert.Nil(t, cache.Get(streamID))
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(time.Second, 0)
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
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, 1)
			}
		}(i)
	}
	wg.Wait()

	// Verify all values were added correctly
	for i := range uint32(numGoroutines) {
		for j := range numOperations {
			streamID := i*numOperations + j
			assert.Equal(t, &mockStreamValue{value: []byte{byte(i)}}, cache.Get(streamID))
		}
	}
}

func TestCache_ConcurrentReadWrite(t *testing.T) {
	cache := NewCache(time.Second, 0)
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
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, uint64(j))
			}
		}(i)
	}

	// Start read goroutines
	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for j := range numOperations {
				streamID := id*numOperations + j
				cache.Get(streamID)
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_ConcurrentAddGet(t *testing.T) {
	cache := NewCache(time.Second, 0)
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
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, 1)
			}
		}(i)
	}

	// Start Get goroutines
	for i := range uint32(numGoroutines) {
		go func(id uint32) {
			defer wg.Done()
			for j := range numOperations {
				streamID := id*numOperations + j
				cache.Get(streamID)
			}
		}(i)
	}

	wg.Wait()
}
