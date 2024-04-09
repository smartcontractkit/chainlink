package capabilities

import (
	"context"
	"slices"

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
	subServices []services.Service
	lggr        logger.Logger
}

var _ services.Service = &registrySyncer{}

var defaultStreamConfig = p2ptypes.StreamConfig{
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
	workflowDONPeers := []string{
		"12D3KooWF3dVeJ6YoT5HFnYhmwQWWMoEwVFzJQ5kKCMX3ZityxMC",
		"12D3KooWQsmok6aD8PZqt3RnJhQRrNzKHLficq7zYFRp7kZ1hHP8",
		"12D3KooWJbZLiMuGeKw78s3LM5TNgBTJHcF39DraxLu14bucG9RN",
		"12D3KooWGqfSPhHKmQycfhRjgUDE2vg9YWZN27Eue8idb2ZUk6EH",
	}
	capabilityDONPeers := []string{
		"12D3KooWHCcyTPmYFB1ydNvNcXw5WyAomRzGSFu1B7hpB4yi8Smf",
		"12D3KooWPv6eqJvYz7TcQWk4Y4XjZ1uQ7mUKahdDXj65ht95zH6a",
	}
	allPeers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)
	addPeersToDONInfo := func(peers []string, donInfo *remotetypes.DON) error {
		for _, peerID := range peers {
			var p ragetypes.PeerID
			err := p.UnmarshalText([]byte(peerID))
			if err != nil {
				return err
			}
			allPeers[p] = defaultStreamConfig
			donInfo.Members = append(donInfo.Members, p)
		}
		return nil
	}
	workflowDonInfo := remotetypes.DON{ID: "workflowDon1"}
	if err := addPeersToDONInfo(workflowDONPeers, &workflowDonInfo); err != nil {
		return err
	}
	capabilityDonInfo := remotetypes.DON{ID: "capabilityDon1"}
	if err := addPeersToDONInfo(capabilityDONPeers, &capabilityDonInfo); err != nil {
		return err
	}
	err := s.peerWrapper.GetPeer().UpdateConnections(allPeers)
	if err != nil {
		return err
	}
	// NOTE: temporary hard-coded capabilities
	capId := "sample_remote_trigger"
	triggerInfo := commoncap.CapabilityInfo{
		ID:             capId,
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
		Version:        "0.0.1",
	}
	myId := s.peerWrapper.GetPeer().ID().String()
	config := remotetypes.RemoteTriggerConfig{
		RegistrationRefreshMs: 20000,
	}
	if slices.Contains(workflowDONPeers, myId) {
		s.lggr.Info("member of a workflow DON - starting remote subscribers")
		triggerCap := remote.NewTriggerSubscriber(config, triggerInfo, capabilityDonInfo, workflowDonInfo, s.dispatcher, nil, s.lggr)
		err = s.registry.Add(ctx, triggerCap)
		if err != nil {
			s.lggr.Errorw("failed to add remote target capability to registry", "error", err)
			return err
		}
		err = s.dispatcher.SetReceiver(capId, capabilityDonInfo.ID, triggerCap)
		if err != nil {
			s.lggr.Errorw("failed to set receiver", "capabilityId", capId, "donId", capabilityDonInfo.ID, "error", err)
			return err
		}
		s.subServices = append(s.subServices, triggerCap)
	}
	if slices.Contains(capabilityDONPeers, myId) {
		s.lggr.Info("member of a capability DON - starting remote publishers")
		workflowDONs := map[string]remotetypes.DON{
			workflowDonInfo.ID: workflowDonInfo,
		}
		underlying := &noOpTrigger{info: triggerInfo, lggr: s.lggr}
		triggerCap := remote.NewTriggerPublisher(config, underlying, triggerInfo, capabilityDonInfo, workflowDONs, s.dispatcher, s.lggr)
		err = s.dispatcher.SetReceiver(capId, capabilityDonInfo.ID, triggerCap)
		if err != nil {
			s.lggr.Errorw("failed to set receiver", "capabilityId", capId, "donId", capabilityDonInfo.ID, "error", err)
			return err
		}
		s.subServices = append(s.subServices, triggerCap)
	}
	// NOTE: temporary service start - should be managed by capability creation
	for _, srv := range s.subServices {
		err = srv.Start(ctx)
		if err != nil {
			s.lggr.Errorw("failed to start remote trigger caller", "error", err)
			return err
		}
	}
	s.lggr.Info("registry syncer started")
	return nil
}

func (s *registrySyncer) Close() error {
	for _, subService := range s.subServices {
		err := subService.Close()
		if err != nil {
			s.lggr.Errorw("failed to close a sub-service", "name", subService.Name(), "error", err)
		}
	}
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

type noOpTrigger struct {
	info commoncap.CapabilityInfo
	lggr logger.Logger
}

func (t *noOpTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return t.info, nil
}

func (t *noOpTrigger) RegisterTrigger(_ context.Context, _ chan<- commoncap.CapabilityResponse, request commoncap.CapabilityRequest) error {
	t.lggr.Infow("no-op trigger RegisterTrigger", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (t *noOpTrigger) UnregisterTrigger(_ context.Context, request commoncap.CapabilityRequest) error {
	t.lggr.Infow("no-op trigger RegisterTrigger", "workflowID", request.Metadata.WorkflowID)
	return nil
}
