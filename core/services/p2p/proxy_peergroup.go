package p2p

import (
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	creproxy "github.com/smartcontractkit/chainlink-protos/cre/impl/proxy"
)

// NewProxyBackedPeerGroupFactory adapts the p2p proxy client's
// creproxy.PeerGroupFactory to a libocr networking.PeerGroupFactory, so it can
// back Don2DonSharedPeer in place of the local rage peer. The method sets are
// identical; this only bridges the (structurally equal) types.
func NewProxyBackedPeerGroupFactory(inner creproxy.PeerGroupFactory) networking.PeerGroupFactory {
	return proxyPeerGroupFactory{inner: inner}
}

type proxyPeerGroupFactory struct {
	inner creproxy.PeerGroupFactory
}

func (a proxyPeerGroupFactory) NewPeerGroup(configDigest ocr2types.ConfigDigest, peerIDs []string, bootstrappers []commontypes.BootstrapperLocator) (networking.PeerGroup, error) {
	bs := make([]creproxy.BootstrapperInfo, len(bootstrappers))
	for i, b := range bootstrappers {
		bs[i] = creproxy.BootstrapperInfo{PeerID: b.PeerID, Addrs: b.Addrs}
	}
	pg, err := a.inner.NewPeerGroup([32]byte(configDigest), peerIDs, bs)
	if err != nil {
		return nil, err
	}
	return proxyPeerGroup{inner: pg}, nil
}

type proxyPeerGroup struct {
	inner creproxy.PeerGroup
}

func (a proxyPeerGroup) NewStream(remotePeerID string, args networking.NewStreamArgs) (networking.Stream, error) {
	a1, ok := args.(networking.NewStreamArgs1)
	if !ok {
		return nil, fmt.Errorf("unsupported NewStreamArgs type %T", args)
	}
	st, err := a.inner.NewStream(remotePeerID, creproxy.StreamArgs{
		StreamName:         a1.StreamName,
		OutgoingBufferSize: a1.OutgoingBufferSize,
		IncomingBufferSize: a1.IncomingBufferSize,
		MaxMessageLength:   a1.MaxMessageLength,
		MessagesLimit:      creproxy.RateLimit{Rate: a1.MessagesLimit.Rate, Capacity: a1.MessagesLimit.Capacity},
		BytesLimit:         creproxy.RateLimit{Rate: a1.BytesLimit.Rate, Capacity: a1.BytesLimit.Capacity},
	})
	if err != nil {
		return nil, err
	}
	// creproxy.PeerGroupStream's method set matches networking.Stream.
	return st, nil
}

func (a proxyPeerGroup) Close() error {
	return a.inner.Close()
}
