package tokendata

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
)

type CachedReader struct {
	Reader

	cache      map[uint64][]byte
	cacheMutex sync.RWMutex
}

func NewCachedReader(reader Reader) *CachedReader {
	return &CachedReader{
		Reader: reader,
		cache:  make(map[uint64][]byte),
	}
}

// ReadTokenData tries to get the tokenData from cache, if not found then calls the underlying reader
// and updates the cache.
func (r *CachedReader) ReadTokenData(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([]byte, error) {
	r.cacheMutex.RLock()
	data, ok := r.cache[msg.SequenceNumber]
	r.cacheMutex.RUnlock()

	if ok {
		return data, nil
	}

	tokenData, err := r.Reader.ReadTokenData(ctx, msg)
	if err != nil {
		return []byte{}, err
	}

	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()

	// Update the cache
	r.cache[msg.SequenceNumber] = tokenData

	return tokenData, nil
}
