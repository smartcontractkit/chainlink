package capabilities

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/internal/json"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
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

// CapabilityID is the unique identifier of the capability in the CR.
// It is calculated as keccak256(abi.encode(capabilityType, capabilityVersion)).
type CapabilityID = [32]byte

type Capability struct {
	ID CapabilityID
	// The `Name` is a partially qualified ID for the capability.
	// Validation: ^[a-z0-9_\-:]{1,32}$
	Name string
	// Semver, e.g., "1.2.3"
	Version      string
	ResponseType CapabilityResponseType
	// An address to the capability configuration contract. Having this defined
	// on a capability enforces consistent configuration across DON instances
	// serving the same capability.
	//
	// The main use cases are:
	// 1) Sharing capability configuration across DON instances
	// 2) Inspect and modify on-chain configuration without off-chain
	// capability code.
	ConfigurationContract common.Address
}

// RemoteRegistryState contains a local cache of the CapabilityRegistry deployed
// on-chain. It is updated by the syncer and is otherwise read-only.
type RemoteRegistryState struct {
	Capabilities  map[CapabilityID]Capability
	CapabilityIDs []CapabilityID
}

type registrySyncer struct {
	peerWrapper           p2ptypes.PeerWrapper
	registry              core.CapabilitiesRegistry
	dispatcher            remotetypes.Dispatcher
	subServices           []services.Service
	wg                    sync.WaitGroup
	lggr                  logger.Logger
	remoteRegistryState   RemoteRegistryState
	chainReader           commontypes.ContractReader
	relayer               loop.Relayer
	remoteRegistryAddress string
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

const maxRetryCount = 60

// RegistrySyncer updates local Registry to match its onchain counterpart
func NewRegistrySyncer(
	peerWrapper p2ptypes.PeerWrapper,
	registry core.CapabilitiesRegistry,
	dispatcher remotetypes.Dispatcher,
	lggr logger.Logger,
	relayer loop.Relayer,
	remoteRegistryAddress string,
) *registrySyncer {
	return &registrySyncer{
		peerWrapper:           peerWrapper,
		registry:              registry,
		dispatcher:            dispatcher,
		lggr:                  lggr,
		relayer:               relayer,
		remoteRegistryAddress: remoteRegistryAddress,
	}
}

func (s *registrySyncer) Start(ctx context.Context) error {
	s.wg.Add(1)
	go s.launch(ctx)
	return nil
}

// NOTE: this implementation of the Syncer is temporary and will be replaced by one
// that reads the configuration from chain (KS-117).
func (s *registrySyncer) launch(ctx context.Context) {
	// Creating a JSON blob for chain-agnostic configuration

	contractReaderConfig := evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"capability_registry": {
				ContractABI: keystone_capability_registry.CapabilityRegistryABI,
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"get_capabilities": {
						ChainSpecificName: "getCapabilities",
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"labelledName": "name"}},
						},
					},
				},
			},
		},
	}

	contractReaderConfigEncoded, err := json.Marshal(contractReaderConfig)
	if err != nil {
		s.lggr.Errorw("failed to marshal contract reader config", "error", err)
	}

	fmt.Println("contractReaderConfigEncoded", contractReaderConfigEncoded)

	cr, err := s.relayer.NewContractReader(ctx, contractReaderConfigEncoded)

	if err != nil {
		// TODO: What to do here?
		s.lggr.Errorw("failed to create contract reader", "error", err)
	}

	err = cr.Bind(ctx, []commontypes.BoundContract{
		{
			Name:    "capability_registry",
			Address: s.remoteRegistryAddress,
		}})

	if err != nil {
		// TODO: What to do here?
		s.lggr.Errorw("failed to bind to capability registry", "error", err)
	}

	defer s.wg.Done()
	// NOTE: temporary hard-coded DONs
	workflowDONPeers := []string{
		"12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N",
		"12D3KooWG1AyvwmCpZ93J8pBQUE1SuzrjDXnT4BeouncHR3jWLCG",
		"12D3KooWGeUKZBRMbx27FUTgBwZa9Ap9Ym92mywwpuqkEtz8XWyv",
		"12D3KooW9zYWQv3STmDeNDidyzxsJSTxoCTLicafgfeEz9nhwhC4",
	}
	triggerDONPeers := []string{
		"12D3KooWBaiTbbRwwt2fbNifiL7Ew9tn3vds9AJE3Nf3eaVBX36m",
		"12D3KooWS7JSY9fzSfWgbCE1S3W2LNY6ZVpRuun74moVBkKj6utE",
		"12D3KooWMMTDXcWhpVnwrdAer1jnVARTmnr3RyT3v7Djg8ZuoBh9",
		"12D3KooWGzVXsKxXsF4zLgxSDM8Gzx1ywq2pZef4PrHMKuVg4K3P",
		"12D3KooWSyjmmzjVtCzwN7bXzZQFmWiJRuVcKBerNjVgL7HdLJBW",
		"12D3KooWLGz9gzhrNsvyM6XnXS3JRkZoQdEzuAvysovnSChNK5ZK",
		"12D3KooWAvZnvknFAfSiUYjATyhzEJLTeKvAzpcLELHi4ogM3GET",
	}
	triggerDONSigners := []string{
		"0x9CcE7293a4Cc2621b61193135A95928735e4795F",
		"0x3c775F20bCB2108C1A818741Ce332Bb5fe0dB925",
		"0x50314239e2CF05555ceeD53E7F47eB2A8Eab0dbB",
		"0xd76A4f98898c3b9A72b244476d7337b50D54BCd8",
		"0x656A873f6895b8a03Fb112dE927d43FA54B2c92A",
		"0x5d1e87d87bF2e0cD4Ea64F381a2dbF45e5f0a553",
		"0x91d9b0062265514f012Eb8fABA59372fD9520f56",
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
		s.lggr.Errorw("failed to add peers to workflow DON info", "error", err)
		return
	}
	triggerCapabilityDonInfo := capabilities.DON{ID: "capabilityDon1", F: 1} // NOTE: misconfiguration - should be 2
	if err := addPeersToDONInfo(triggerDONPeers, &triggerCapabilityDonInfo); err != nil {
		s.lggr.Errorw("failed to add peers to trigger DON info", "error", err)
		return
	}
	err := s.peerWrapper.GetPeer().UpdateConnections(allPeers)
	if err != nil {
		s.lggr.Errorw("failed to update connections", "error", err)
		return
	}
	// NOTE: temporary hard-coded capabilities
	capId := "streams-trigger"
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
		RegistrationExpiryMs:    60000,
		MinResponsesToAggregate: uint32(triggerCapabilityDonInfo.F) + 1,
	}
	if slices.Contains(workflowDONPeers, myId) {
		s.lggr.Info("member of a workflow DON - starting remote subscribers")
		codec := streams.NewCodec(s.lggr)
		aggregator := triggers.NewMercuryRemoteAggregator(codec, hexStringsToBytes(triggerDONSigners), int(triggerCapabilityDonInfo.F+1), s.lggr)
		triggerCap := remote.NewTriggerSubscriber(config, triggerInfo, triggerCapabilityDonInfo, workflowDonInfo, s.dispatcher, aggregator, s.lggr)
		err = s.registry.Add(ctx, triggerCap)
		if err != nil {
			s.lggr.Errorw("failed to add remote target capability to registry", "error", err)
			return
		}
		err = s.dispatcher.SetReceiver(capId, triggerCapabilityDonInfo.ID, triggerCap)
		if err != nil {
			s.lggr.Errorw("workflow DON failed to set receiver", "capabilityId", capId, "donId", triggerCapabilityDonInfo.ID, "error", err)
			return
		}
		s.subServices = append(s.subServices, triggerCap)
	}
	if slices.Contains(triggerDONPeers, myId) {
		s.lggr.Info("member of a capability DON - starting remote publishers")

		/*{
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
		}*/

		count := 0
		for {
			count++
			if count > maxRetryCount {
				s.lggr.Error("failed to get Streams Trigger from the Registry")
				return
			}
			underlying, err2 := s.registry.GetTrigger(ctx, capId)
			if err2 != nil {
				// NOTE: it's possible that the jobs are not launched yet at this moment.
				// If not found yet, Syncer won't add to Registry but retry on the next tick.
				s.lggr.Infow("trigger not found yet ...", "capabilityId", capId, "error", err2)
				time.Sleep(1 * time.Second)
				continue
			}
			workflowDONs := map[string]capabilities.DON{
				workflowDonInfo.ID: workflowDonInfo,
			}
			triggerCap := remote.NewTriggerPublisher(config, underlying, triggerInfo, triggerCapabilityDonInfo, workflowDONs, s.dispatcher, s.lggr)
			err = s.dispatcher.SetReceiver(capId, triggerCapabilityDonInfo.ID, triggerCap)
			if err != nil {
				s.lggr.Errorw("capability DON failed to set receiver", "capabilityId", capId, "donId", triggerCapabilityDonInfo.ID, "error", err)
				return
			}
			s.subServices = append(s.subServices, triggerCap)
			break
		}
	}
	// NOTE: temporary service start - should be managed by capability creation
	for _, srv := range s.subServices {
		err = srv.Start(ctx)
		if err != nil {
			s.lggr.Errorw("failed to start remote trigger caller", "error", err)
			return
		}
	}
	s.lggr.Info("registry syncer started")
}

func (s *registrySyncer) Close() error {
	s.wg.Wait()
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

		reports := []datastreams.FeedReport{
			{
				FeedID:               "0x0003fbba4fce42f65d6032b18aee53efdf526cc734ad296cb57565979d883bdd",
				FullReport:           []byte{0x11, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[0].Bytes(),
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x0003c317fec7fad514c67aacc6366bf2f007ce37100e3cddcacd0ccaa1f3746d",
				FullReport:           []byte{0x22, 0xaa, 0xbb, 0xcc},
				BenchmarkPrice:       prices[1].Bytes(),
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               "0x0003da6ab44ea9296674d80fe2b041738189103d6b4ea9a4d34e2f891fa93d12",
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

func hexStringsToBytes(strs []string) (res [][]byte) {
	for _, s := range strs {
		b, _ := hex.DecodeString(s[2:])
		res = append(res, b)
	}
	return res
}
