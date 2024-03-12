package p2p

import (
	"context"

	ocrnetworking "github.com/smartcontractkit/libocr/networking/types"
)

var _ ocrnetworking.DiscovererDatabase = &InMemoryDiscovererDatabase{}

type InMemoryDiscovererDatabase struct {
	store map[string][]byte
}

func NewInMemoryDiscovererDatabase() *InMemoryDiscovererDatabase {
	return &InMemoryDiscovererDatabase{make(map[string][]byte)}
}

func (d *InMemoryDiscovererDatabase) StoreAnnouncement(ctx context.Context, peerID string, ann []byte) error {
	d.store[peerID] = ann
	return nil
}

func (d *InMemoryDiscovererDatabase) ReadAnnouncements(ctx context.Context, peerIDs []string) (map[string][]byte, error) {
	ret := make(map[string][]byte)
	for _, peerID := range peerIDs {
		if ann, ok := d.store[peerID]; ok {
			ret[peerID] = ann
		}
	}
	return ret, nil
}
