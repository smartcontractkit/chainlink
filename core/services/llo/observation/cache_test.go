package observation

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

type mockStreamValue struct {
	value []byte
}

func (m *mockStreamValue) Value() interface{} {
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
	return []byte(fmt.Sprintf("%d", m.value)), nil
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
			cache := NewCache(ocr2types.ConfigDigest{}, tt.maxAge, tt.cleanupInterval)
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
			cache := NewCache(ocr2types.ConfigDigest{}, tt.maxAge, 0)

			if tt.wantFound {
				cache.Add(tt.streamID, tt.value, tt.seqNr)
			}

			if tt.beforeGet != nil {
				tt.beforeGet(cache)
			}

			gotValue, gotFound := cache.Get(tt.streamID)
			assert.Equal(t, tt.wantFound, gotFound)
			if tt.wantFound {
				assert.Equal(t, tt.wantValue, gotValue)
			}
		})
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(ocr2types.ConfigDigest{}, time.Nanosecond*100, time.Millisecond)
	streamID := llotypes.StreamID(1)
	value := &mockStreamValue{value: []byte{42}}

	cache.Add(streamID, value, 10)
	time.Sleep(time.Millisecond * 2)

	gotValue, gotFound := cache.Get(streamID)
	assert.False(t, gotFound)
	assert.Nil(t, gotValue)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(ocr2types.ConfigDigest{}, time.Second, 0)
	const numGoroutines = 10
	const numOperations = uint32(1000)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Test concurrent Add operations
	for i := uint32(0); i < numGoroutines; i++ {
		go func(id uint32) {
			defer wg.Done()
			for j := uint32(0); j < numOperations; j++ {
				streamID := id*numOperations + j
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, 1)
			}
		}(i)
	}
	wg.Wait()

	// Verify all values were added correctly
	for i := uint32(0); i < numGoroutines; i++ {
		for j := uint32(0); j < numOperations; j++ {
			streamID := i*numOperations + j
			value, found := cache.Get(streamID)
			assert.True(t, found)
			assert.Equal(t, &mockStreamValue{value: []byte{byte(i)}}, value)
		}
	}
}

func TestCache_ConcurrentReadWrite(t *testing.T) {
	cache := NewCache(ocr2types.ConfigDigest{}, time.Second, 0)
	const numGoroutines = 10
	const numOperations = uint32(1000)

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // Double for read and write goroutines

	// Start write goroutines
	for i := uint32(0); i < numGoroutines; i++ {
		go func(id uint32) {
			defer wg.Done()
			for j := uint32(0); j < numOperations; j++ {
				streamID := id*numOperations + j
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, uint64(j))
			}
		}(i)
	}

	// Start read goroutines
	for i := uint32(0); i < numGoroutines; i++ {
		go func(id uint32) {
			defer wg.Done()
			for j := uint32(0); j < numOperations; j++ {
				streamID := id*numOperations + j
				cache.Get(streamID)
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_ConcurrentAddGet(t *testing.T) {
	cache := NewCache(ocr2types.ConfigDigest{}, time.Second, 0)
	const numGoroutines = 10
	const numOperations = uint32(1000)

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // Double for Add and Get goroutines

	// Start Add goroutines
	for i := uint32(0); i < numGoroutines; i++ {
		go func(id uint32) {
			defer wg.Done()
			for j := uint32(0); j < numOperations; j++ {
				streamID := id*numOperations + j
				cache.Add(streamID, &mockStreamValue{value: []byte{byte(id)}}, 1)
			}
		}(i)
	}

	// Start Get goroutines
	for i := uint32(0); i < numGoroutines; i++ {
		go func(id uint32) {
			defer wg.Done()
			for j := uint32(0); j < numOperations; j++ {
				streamID := id*numOperations + j
				cache.Get(streamID)
			}
		}(i)
	}

	wg.Wait()
}
