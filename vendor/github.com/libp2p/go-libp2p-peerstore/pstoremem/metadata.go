package pstoremem

import (
	"sync"

	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
)

var internKeys = map[string]bool{
	"AgentVersion":    true,
	"ProtocolVersion": true,
}

type metakey struct {
	id  peer.ID
	key string
}

type memoryPeerMetadata struct {
	// store other data, like versions
	//ds ds.ThreadSafeDatastore
	ds       map[metakey]interface{}
	dslock   sync.RWMutex
	interned map[string]interface{}
}

var _ pstore.PeerMetadata = (*memoryPeerMetadata)(nil)

func NewPeerMetadata() *memoryPeerMetadata {
	return &memoryPeerMetadata{
		ds:       make(map[metakey]interface{}),
		interned: make(map[string]interface{}),
	}
}

func (ps *memoryPeerMetadata) Put(p peer.ID, key string, val interface{}) error {
	if err := p.Validate(); err != nil {
		return err
	}
	ps.dslock.Lock()
	defer ps.dslock.Unlock()
	if vals, ok := val.(string); ok && internKeys[key] {
		if interned, ok := ps.interned[vals]; ok {
			val = interned
		} else {
			ps.interned[vals] = val
		}
	}
	ps.ds[metakey{p, key}] = val
	return nil
}

func (ps *memoryPeerMetadata) Get(p peer.ID, key string) (interface{}, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	ps.dslock.RLock()
	defer ps.dslock.RUnlock()
	i, ok := ps.ds[metakey{p, key}]
	if !ok {
		return nil, pstore.ErrNotFound
	}
	return i, nil
}
