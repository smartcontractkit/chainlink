package capabilities

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type registrySyncer struct {
	peerWrapper p2ptypes.PeerWrapper
	registry    types.CapabilitiesRegistry
	lggr        logger.Logger
}

var _ services.Service = &registrySyncer{}

// RegistrySyncer updates local Registry to match its onchain counterpart
func NewRegistrySyncer(peerWrapper p2ptypes.PeerWrapper, registry types.CapabilitiesRegistry, lggr logger.Logger) *registrySyncer {
	return &registrySyncer{
		peerWrapper: peerWrapper,
		registry:    registry,
		lggr:        lggr,
	}
}

func (s *registrySyncer) Start(ctx context.Context) error {
	// NOTE: temporary hard-coded values
	defaultStreamConfig := p2ptypes.StreamConfig{
		IncomingMessageBufferSize: 1000000,
		OutgoingMessageBufferSize: 1000000,
		MaxMessageLenBytes:        100000,
		MessageRateLimiter: ragep2p.TokenBucketParams{
			Rate:     10.0,
			Capacity: 1000,
		},
		BytesRateLimiter: ragep2p.TokenBucketParams{
			Rate:     10.0,
			Capacity: 1000,
		},
	}
	peerIDs := []string{
		"12D3KooWF3dVeJ6YoT5HFnYhmwQWWMoEwVFzJQ5kKCMX3ZityxMC",
		"12D3KooWQsmok6aD8PZqt3RnJhQRrNzKHLficq7zYFRp7kZ1hHP8",
		"12D3KooWJbZLiMuGeKw78s3LM5TNgBTJHcF39DraxLu14bucG9RN",
		"12D3KooWGqfSPhHKmQycfhRjgUDE2vg9YWZN27Eue8idb2ZUk6EH",
	}
	peers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)
	for _, peerID := range peerIDs {
		var p ragetypes.PeerID
		err := p.UnmarshalText([]byte(peerID))
		if err != nil {
			return err
		}
		peers[p] = defaultStreamConfig
	}
	return s.peerWrapper.GetPeer().UpdateConnections(peers)
}

func (s *registrySyncer) Close() error {
	return s.peerWrapper.GetPeer().UpdateConnections(map[ragetypes.PeerID]p2ptypes.StreamConfig{})
}

func (s *registrySyncer) Ready() error {
	return nil
}

func (s *registrySyncer) HealthReport() map[string]error {
	return nil
}

func (s *registrySyncer) Name() string {
	return "RegistrySyncer"
}
