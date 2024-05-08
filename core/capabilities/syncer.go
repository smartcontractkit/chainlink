package capabilities

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type registrySyncer struct {
	peerWrapper p2ptypes.PeerWrapper
	registry    core.CapabilitiesRegistry
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
		Rate:     100.0,
		Capacity: 1000,
	},
	BytesRateLimiter: ragep2p.TokenBucketParams{
		Rate:     100000.0,
		Capacity: 1000000,
	},
}

// RegistrySyncer updates local Registry to match its onchain counterpart
func NewRegistrySyncer(peerWrapper p2ptypes.PeerWrapper, registry core.CapabilitiesRegistry, dispatcher remotetypes.Dispatcher, lggr logger.Logger) *registrySyncer {
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
		"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N",
		"12D3KooWG1AyvwmCpZ93J8pBQUE1SuzrjDXnT4BeouncHR3jWLCG",
		"12D3KooWGeUKZBRMbx27FUTgBwZa9Ap9Ym92mywwpuqkEtz8XWyv",
		"12D3KooW9zYWQv3STmDeNDidyzxsJSTxoCTLicafgfeEz9nhwhC4",
	}
	triggerDONPeers := []string{
		"12D3KooWJrthXtnPHw7xyHFAxo6NxifYTvc8igKYaA6wRRRqtsMb",
		"12D3KooWFQekP9sGex4XhqEJav5EScjTpDVtDqJFg1JvrePBCEGJ",
		"12D3KooWFLEq4hYtdyKWwe47dXGEbSiHMZhmr5xLSJNhpfiEz8NF",
		"12D3KooWN2hztiXNNS1jMQTTvvPRYcarK1C7T3Mdqk4x4gwyo5WS",
	}
	allPeers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)
	addPeersToDONInfo := func(peers []string, donInfo *capabilities.DON) error {
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
	workflowDonInfo := capabilities.DON{ID: "workflowDon1", F: 1}
	if err := addPeersToDONInfo(workflowDONPeers, &workflowDonInfo); err != nil {
		return err
	}
	triggerCapabilityDonInfo := capabilities.DON{ID: "capabilityDon1", F: 1}
	if err := addPeersToDONInfo(triggerDONPeers, &triggerCapabilityDonInfo); err != nil {
		return err
	}
	err := s.peerWrapper.GetPeer().UpdateConnections(allPeers)
	if err != nil {
		return err
	}
	// NOTE: temporary hard-coded capabilities
	capId := "mercury-trigger@1"
	triggerInfo := commoncap.CapabilityInfo{
		ID:             capId,
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
		DON:            &triggerCapabilityDonInfo,
	}
	myId := s.peerWrapper.GetPeer().ID().String()
	config := remotetypes.RemoteTriggerConfig{
		RegistrationRefreshMs:   20000,
		MinResponsesToAggregate: uint32(triggerCapabilityDonInfo.F) + 1,
	}
	if slices.Contains(workflowDONPeers, myId) {
		s.lggr.Info("member of a workflow DON - starting remote subscribers")
		aggregator := triggers.NewMercuryRemoteAggregator(s.lggr)
		triggerCap := remote.NewTriggerSubscriber(config, triggerInfo, triggerCapabilityDonInfo, workflowDonInfo, s.dispatcher, aggregator, s.lggr)
		err = s.registry.Add(ctx, triggerCap)
		if err != nil {
			s.lggr.Errorw("failed to add remote target capability to registry", "error", err)
			return err
		}
		err = s.dispatcher.SetReceiver(capId, triggerCapabilityDonInfo.ID, triggerCap)
		if err != nil {
			s.lggr.Errorw("workflow DON failed to set receiver", "capabilityId", capId, "donId", triggerCapabilityDonInfo.ID, "error", err)
			return err
		}
		s.subServices = append(s.subServices, triggerCap)
	}
	if slices.Contains(triggerDONPeers, myId) {
		s.lggr.Info("member of a capability DON - starting remote publishers")
		workflowDONs := map[string]capabilities.DON{
			workflowDonInfo.ID: workflowDonInfo,
		}
		underlying := triggers.NewMercuryTriggerService(1000, s.lggr)
		triggerCap := remote.NewTriggerPublisher(config, underlying, triggerInfo, triggerCapabilityDonInfo, workflowDONs, s.dispatcher, s.lggr)
		err = s.dispatcher.SetReceiver(capId, triggerCapabilityDonInfo.ID, triggerCap)
		if err != nil {
			s.lggr.Errorw("capability DON failed to set receiver", "capabilityId", capId, "donId", triggerCapabilityDonInfo.ID, "error", err)
			return err
		}
		s.subServices = append(s.subServices, underlying)
		s.subServices = append(s.subServices, triggerCap)
		// NOTE: temporary mock Mercury data producer
		mockMercuryDataProducer := NewMockMercuryDataProducer(underlying, s.lggr)
		s.subServices = append(s.subServices, mockMercuryDataProducer)
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

type mockMercuryDataProducer struct {
	trigger *triggers.MercuryTriggerService
	wg      sync.WaitGroup
	closeCh chan struct{}
	lggr    logger.Logger
}

var _ services.Service = &mockMercuryDataProducer{}

func NewMockMercuryDataProducer(trigger *triggers.MercuryTriggerService, lggr logger.Logger) *mockMercuryDataProducer {
	return &mockMercuryDataProducer{
		trigger: trigger,
		closeCh: make(chan struct{}),
		lggr:    lggr,
	}
}

func (m *mockMercuryDataProducer) Start(ctx context.Context) error {
	m.wg.Add(1)
	go m.loop()
	return nil
}

func (m *mockMercuryDataProducer) loop() {
	defer m.wg.Done()

	sleepSec := 60
	ticker := time.NewTicker(time.Duration(sleepSec) * time.Second)
	defer ticker.Stop()

	prices := []int64{300000, 40000, 5000000}

	for range ticker.C {
		for i := range prices {
			prices[i] = prices[i] + 1
		}

		reports := []mercury.FeedReport{
			{
				FeedID:               "0x1111111111111111111100000000000000000000000000000000000000000000",
				FullReport:           []byte{0x11, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[0],
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x2222222222222222222200000000000000000000000000000000000000000000",
				FullReport:           []byte{0x22, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[1],
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x3333333333333333333300000000000000000000000000000000000000000000",
				FullReport:           []byte{0x33, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[2],
				ObservationTimestamp: time.Now().Unix(),
			},
		}

		m.lggr.Infow("New set of Mercury reports", "timestamp", time.Now().Unix(), "payload", reports)
		err := m.trigger.ProcessReport(reports)
		if err != nil {
			m.lggr.Errorw("failed to process Mercury reports", "err", err, "timestamp", time.Now().Unix(), "payload", reports)
		}
	}
}

func (m *mockMercuryDataProducer) Close() error {
	close(m.closeCh)
	m.wg.Wait()
	return nil
}

func (m *mockMercuryDataProducer) HealthReport() map[string]error {
	return nil
}

func (m *mockMercuryDataProducer) Ready() error {
	return nil
}

func (m *mockMercuryDataProducer) Name() string {
	return "mockMercuryDataProducer"
}
