package capabilities

import (
	"context"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type registrySyncer struct {
	peerWrapper p2ptypes.PeerWrapper
	registry    types.CapabilitiesRegistry
	dispatcher  remotetypes.Dispatcher
	lggr        logger.Logger
}

var _ services.Service = &registrySyncer{}

// RegistrySyncer updates local Registry to match its onchain counterpart
func NewRegistrySyncer(peerWrapper p2ptypes.PeerWrapper, registry types.CapabilitiesRegistry, dispatcher remotetypes.Dispatcher, lggr logger.Logger) *registrySyncer {
	return &registrySyncer{
		peerWrapper: peerWrapper,
		registry:    registry,
		dispatcher:  dispatcher,
		lggr:        lggr,
	}
}

func (s *registrySyncer) Start(ctx context.Context) error {
	// NOTE: temporary hard-coded DONs
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
	donInfo := &remotetypes.DON{
		ID: "don1",
	}
	for _, peerID := range peerIDs {
		var p ragetypes.PeerID
		err := p.UnmarshalText([]byte(peerID))
		if err != nil {
			return err
		}
		peers[p] = defaultStreamConfig
		donInfo.Members = append(donInfo.Members, p)
	}
	err := s.peerWrapper.GetPeer().UpdateConnections(peers)
	if err != nil {
		return err
	}
	// NOTE: temporary hard-coded capabilities
	capId := "sample_remote_target"
	targetCap := remote.NewRemoteTargetCaller(commoncap.CapabilityInfo{
		ID:             capId,
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
	}, donInfo, s.dispatcher, s.lggr)
	err = s.registry.Add(ctx, targetCap)
	if err != nil {
		s.lggr.Error("failed to add remote target capability to registry")
		return err
	}
	err = s.dispatcher.SetReceiver(capId, donInfo.ID, targetCap)
	if err != nil {
		s.lggr.Errorw("failed to set receiver", "capabilityId", capId, "donId", donInfo.ID, "error", err)
		return err
	}
	s.lggr.Info("registry syncer started")
	return nil
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
