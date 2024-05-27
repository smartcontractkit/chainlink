package capabilities

import (
	"context"
	"math/big"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
)

// CapabilityResponseType indicates whether remote response requires
// aggregation or is an already aggregated report. There are multiple
// possible ways to aggregate.
type CapabilityResponseType int

// CapabilityResponseType enum values.
const (
	// No additional aggregation is needed on the remote response.
	Report CapabilityResponseType = 0
	// A number of identical observations need to be aggregated.
	ObservationIdentical CapabilityResponseType = 1
)

type registrySyncer struct {
	peerWrapper    p2ptypes.PeerWrapper
	registry       core.CapabilitiesRegistry
	remoteRegistry *remoteRegistry
	dispatcher     remotetypes.Dispatcher
	subServices    []services.Service
	lggr           logger.Logger
	client         evmclient.Client
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

var CALLER_ADDRESS = types.MustEIP55Address("0x0000000000000000000000000000000000000001").Address()

// Should extract chainReader from relayer and pass it to the syncer
// RegistrySyncer updates local Registry to match its onchain counterpart
func NewRegistrySyncer(
	peerWrapper p2ptypes.PeerWrapper,
	registry core.CapabilitiesRegistry,
	dispatcher remotetypes.Dispatcher,
	lggr logger.Logger,
	remoteRegistry *remoteRegistry,
	client evmclient.Client,
) *registrySyncer {
	// db := pgtest.NewSqlxDB(t)
	// lpOpts := logpoller.Opts{
	// 	PollPeriod:               time.Millisecond,
	// 	FinalityDepth:            4,
	// 	BackfillBatchSize:        1,
	// 	RpcBatchSize:             1,
	// 	KeepFinalizedBlocksDepth: 10000,
	// }
	// lp := logpoller.NewLogPoller(
	// 	logpoller.NewORM(testutils.SimulatedChainID, db, lggr),
	// 	simulatedBackendClient,
	// 	lggr,
	// 	lpOpts,
	// )

	// chainConfig := types.ChainReaderConfig{
	// 	Contracts: map[string]types.ChainContractReader{
	// 		"capability_registry": {
	// 			ContractABI: keystone_capability_registry.CapabilityRegistryABI,
	// 			Configs: map[string]*types.ChainReaderDefinition{
	// 				"get_capabilities": {
	// 					ChainSpecificName: "getCapabilities",
	// 					OutputModifications: codec.ModifiersConfig{
	// 						&codec.RenameModifierConfig{Fields: map[string]string{"labelledName": "name"}},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// cr, err := evm.NewChainReaderService(ctx, lggr, lp, simulatedBackendClient, chainConfig)

	// cr.Bind(ctx, []commontypes.BoundContract{
	// 	{
	// 		Name:    "capability_registry",
	// 		Address: capabilityRegistry.Address().String(),
	// 	}})

	return &registrySyncer{
		peerWrapper:    peerWrapper,
		registry:       registry,
		dispatcher:     dispatcher,
		lggr:           lggr,
		remoteRegistry: remoteRegistry,
		client:         client,
	}
}

func (s *registrySyncer) Start(ctx context.Context) error {
	// INITIALIZE SYNCER

	// type Cap struct {
	// 	Name                  string
	// 	Version               string
	// 	ResponseType          int
	// 	ConfigurationContract []byte
	// }

	// var returnedCapabilities []Cap

	// err = cr.GetLatestValue(ctx, "capability_registry", "get_capabilities", nil, &returnedCapabilities)

	// fmt.Println("Returned capabilities:", returnedCapabilities)

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
	capId := "mercury-trigger"
	triggerInfo := capabilities.CapabilityInfo{
		ID:             capId,
		CapabilityType: capabilities.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
		Version:        "0.0.1",
		DON:            &triggerCapabilityDonInfo,
	}
	myId := s.peerWrapper.GetPeer().ID().String()
	config := remotetypes.RemoteTriggerConfig{
		RegistrationRefreshMs:   20000,
		MinResponsesToAggregate: uint32(triggerCapabilityDonInfo.F) + 1,
	}
	if slices.Contains(workflowDONPeers, myId) {
		s.lggr.Info("member of a workflow DON - starting remote subscribers")
		codec := streams.NewCodec()
		aggregator := triggers.NewMercuryRemoteAggregator(codec, s.lggr)
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

		{
			// ---- This is for local tests only, until a full-blown Syncer is implemented
			// ---- Normally this is set up asynchronously (by the Relayer + job specs in Mercury's case)
			localTrigger := triggers.NewMercuryTriggerService(1000, s.lggr)
			mockMercuryDataProducer := NewMockMercuryDataProducer(localTrigger, s.lggr)
			err = s.registry.Add(ctx, localTrigger)
			if err != nil {
				s.lggr.Errorw("failed to add local trigger capability to registry", "error", err)
				return err
			}
			s.subServices = append(s.subServices, localTrigger)
			s.subServices = append(s.subServices, mockMercuryDataProducer)
			// ----
		}

		underlying, err2 := s.registry.GetTrigger(ctx, capId)
		if err2 != nil {
			// NOTE: it's possible that the jobs are not launched yet at this moment.
			// If not found yet, Syncer won't add to Registry but retry on the next tick.
			return err2
		}
		workflowDONs := map[string]capabilities.DON{
			workflowDonInfo.ID: workflowDonInfo,
		}
		triggerCap := remote.NewTriggerPublisher(config, underlying, triggerInfo, triggerCapabilityDonInfo, workflowDONs, s.dispatcher, s.lggr)
		err = s.dispatcher.SetReceiver(capId, triggerCapabilityDonInfo.ID, triggerCap)
		if err != nil {
			s.lggr.Errorw("capability DON failed to set receiver", "capabilityId", capId, "donId", triggerCapabilityDonInfo.ID, "error", err)
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

	capabilityRegistry, err := kcr.NewCapabilityRegistry(s.remoteRegistry.address, s.client)
	if err != nil {
		s.lggr.Errorw("failed to create capability registry", "error", err)
		return err
	}

	capabilities, err := capabilityRegistry.GetCapabilities(&bind.CallOpts{})
	if err != nil {
		s.lggr.Errorw("failed to get capabilities from on-chain registry", "error", err)
		return err
	}

	for _, capability := range capabilities {
		capabilityID, err := capabilityRegistry.GetHashedCapabilityId(&bind.CallOpts{}, capability.LabelledName, capability.Version)
		if err != nil {
			s.lggr.Errorw("failed to get capability ID", "error", err)
			return err
		}

		capabilityStruct := Capability{
			ID:                    capabilityID,
			Name:                  capability.LabelledName,
			Version:               capability.Version,
			ResponseType:          CapabilityResponseType(capability.ResponseType),
			ConfigurationContract: capability.ConfigurationContract,
		}
		s.lggr.Infof("capability struct %v", capabilityStruct)

		s.remoteRegistry.capabilities = append(s.remoteRegistry.capabilities, capabilityStruct)
	}

	s.lggr.Infof("capabilities, %v", capabilities)
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

	prices := []*big.Int{big.NewInt(300000), big.NewInt(40000), big.NewInt(5000000)}

	for range ticker.C {
		for i := range prices {
			prices[i].Add(prices[i], big.NewInt(1))
		}

		reports := []mercury.FeedReport{
			{
				FeedID:               "0x1111111111111111111100000000000000000000000000000000000000000000",
				FullReport:           []byte{0x11, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[0].Bytes(),
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x2222222222222222222200000000000000000000000000000000000000000000",
				FullReport:           []byte{0x22, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[1].Bytes(),
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x3333333333333333333300000000000000000000000000000000000000000000",
				FullReport:           []byte{0x33, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[2].Bytes(),
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
