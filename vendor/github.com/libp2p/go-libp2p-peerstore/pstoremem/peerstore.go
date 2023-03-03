package pstoremem

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"io"
)

type pstoremem struct {
	peerstore.Metrics

	*memoryKeyBook
	*memoryAddrBook
	*memoryProtoBook
	*memoryPeerMetadata
}

// NewPeerstore creates an in-memory threadsafe collection of peers.
func NewPeerstore() *pstoremem {
	return &pstoremem{
		Metrics:            pstore.NewMetrics(),
		memoryKeyBook:      NewKeyBook(),
		memoryAddrBook:     NewAddrBook(),
		memoryProtoBook:    NewProtoBook(),
		memoryPeerMetadata: NewPeerMetadata(),
	}
}

func (ps *pstoremem) Close() (err error) {
	var errs []error
	weakClose := func(name string, c interface{}) {
		if cl, ok := c.(io.Closer); ok {
			if err = cl.Close(); err != nil {
				errs = append(errs, fmt.Errorf("%s error: %s", name, err))
			}
		}
	}

	weakClose("keybook", ps.memoryKeyBook)
	weakClose("addressbook", ps.memoryAddrBook)
	weakClose("protobook", ps.memoryProtoBook)
	weakClose("peermetadata", ps.memoryPeerMetadata)

	if len(errs) > 0 {
		return fmt.Errorf("failed while closing peerstore; err(s): %q", errs)
	}
	return nil
}

func (ps *pstoremem) Peers() peer.IDSlice {
	set := map[peer.ID]struct{}{}
	for _, p := range ps.PeersWithKeys() {
		set[p] = struct{}{}
	}
	for _, p := range ps.PeersWithAddrs() {
		set[p] = struct{}{}
	}

	pps := make(peer.IDSlice, 0, len(set))
	for p := range set {
		pps = append(pps, p)
	}
	return pps
}

func (ps *pstoremem) PeerInfo(p peer.ID) peer.AddrInfo {
	return peer.AddrInfo{
		ID:    p,
		Addrs: ps.memoryAddrBook.Addrs(p),
	}
}
