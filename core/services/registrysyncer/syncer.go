package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type Launcher interface {
	Launch(ctx context.Context, registry *LocalRegistry) error
}

type Syncer interface {
	services.Service
	AddLauncher(h ...Launcher)
}

type registrySyncer struct {
	services.StateMachine
	stopCh          services.StopChan
	launchers       []Launcher
	reader          types.ContractReader
	initReader      func(ctx context.Context, lggr logger.Logger, relayer contractReaderFactory, registryAddress string) (types.ContractReader, error)
	relayer         contractReaderFactory
	registryAddress string
	peerWrapper     p2ptypes.PeerWrapper

	wg   sync.WaitGroup
	lggr logger.Logger
	mu   sync.RWMutex
}

var _ services.Service = &registrySyncer{}

var (
	defaultTickInterval = 12 * time.Second
)

// New instantiates a new RegistrySyncer
func New(
	lggr logger.Logger,
	peerWrapper p2ptypes.PeerWrapper,
	relayer contractReaderFactory,
	registryAddress string,
) (*registrySyncer, error) {
	stopCh := make(services.StopChan)
	return &registrySyncer{
		stopCh:          stopCh,
		lggr:            lggr.Named("RegistrySyncer"),
		relayer:         relayer,
		registryAddress: registryAddress,
		initReader:      newReader,
		peerWrapper:     peerWrapper,
	}, nil
}

type contractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

// NOTE: this can't be called while initializing the syncer and needs to be called in the sync loop.
// This is because Bind() makes an onchain call to verify that the contract address exists, and if
// called during initialization, this results in a "no live nodes" error.
func newReader(ctx context.Context, lggr logger.Logger, relayer contractReaderFactory, remoteRegistryAddress string) (types.ContractReader, error) {
	contractReaderConfig := evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"CapabilitiesRegistry": {
				ContractABI: kcr.CapabilitiesRegistryABI,
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"getDONs": {
						ChainSpecificName: "getDONs",
					},
					"getCapabilities": {
						ChainSpecificName: "getCapabilities",
					},
					"getNodes": {
						ChainSpecificName: "getNodes",
					},
				},
			},
		},
	}

	contractReaderConfigEncoded, err := json.Marshal(contractReaderConfig)
	if err != nil {
		return nil, err
	}

	cr, err := relayer.NewContractReader(ctx, contractReaderConfigEncoded)
	if err != nil {
		return nil, err
	}

	err = cr.Bind(ctx, []types.BoundContract{
		{
			Address: remoteRegistryAddress,
			Name:    "CapabilitiesRegistry",
		},
	})

	return cr, err
}

func (s *registrySyncer) Start(ctx context.Context) error {
	return s.StartOnce("RegistrySyncer", func() error {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.syncLoop()
		}()
		return nil
	})
}

func (s *registrySyncer) syncLoop() {
	ctx, cancel := s.stopCh.NewCtx()
	defer cancel()

	ticker := time.NewTicker(defaultTickInterval)
	defer ticker.Stop()

	// Sync for a first time outside the loop; this means we'll start a remote
	// sync immediately once spinning up syncLoop, as by default a ticker will
	// fire for the first time at T+N, where N is the interval.
	s.lggr.Debug("starting initial sync with remote registry")
	err := s.sync(ctx)
	if err != nil {
		s.lggr.Errorw("failed to sync with remote registry", "error", err)
	}

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.lggr.Debug("starting regular sync with the remote registry")
			err := s.sync(ctx)
			if err != nil {
				s.lggr.Errorw("failed to sync with remote registry", "error", err)
			}
		}
	}
}

func unmarshalCapabilityConfig(data []byte) (capabilities.CapabilityConfiguration, error) {
	cconf := &capabilitiespb.CapabilityConfig{}
	err := proto.Unmarshal(data, cconf)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}

	var rtc capabilities.RemoteTriggerConfig
	if prtc := cconf.GetRemoteTriggerConfig(); prtc != nil {
		rtc.RegistrationRefresh = prtc.RegistrationRefresh.AsDuration()
		rtc.RegistrationExpiry = prtc.RegistrationExpiry.AsDuration()
		rtc.MinResponsesToAggregate = prtc.MinResponsesToAggregate
		rtc.MessageExpiry = prtc.MessageExpiry.AsDuration()
	}

	return capabilities.CapabilityConfiguration{
		DefaultConfig:       values.FromMapValueProto(cconf.DefaultConfig),
		RemoteTriggerConfig: rtc,
	}, nil
}

func (s *registrySyncer) localRegistry(ctx context.Context) (*LocalRegistry, error) {
	caps := []kcr.CapabilitiesRegistryCapabilityInfo{}
	err := s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getCapabilities", primitives.Unconfirmed, nil, &caps)
	if err != nil {
		return nil, err
	}

	idsToCapabilities := map[string]Capability{}
	hashedIDsToCapabilityIDs := map[[32]byte]string{}
	for _, c := range caps {
		cid := fmt.Sprintf("%s@%s", c.LabelledName, c.Version)
		idsToCapabilities[cid] = Capability{
			ID:             cid,
			CapabilityType: toCapabilityType(c.CapabilityType),
		}

		hashedIDsToCapabilityIDs[c.HashedId] = cid
	}

	dons := []kcr.CapabilitiesRegistryDONInfo{}
	err = s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getDONs", primitives.Unconfirmed, nil, &dons)
	if err != nil {
		return nil, err
	}

	idsToDONs := map[DonID]DON{}
	for _, d := range dons {
		cc := map[string]capabilities.CapabilityConfiguration{}
		for _, dc := range d.CapabilityConfigurations {
			cid, ok := hashedIDsToCapabilityIDs[dc.CapabilityId]
			if !ok {
				return nil, fmt.Errorf("invariant violation: could not find full ID for hashed ID %s", dc.CapabilityId)
			}

			cconf, innerErr := unmarshalCapabilityConfig(dc.Config)
			if innerErr != nil {
				return nil, innerErr
			}

			cconf.RemoteTriggerConfig.ApplyDefaults()

			cc[cid] = cconf
		}

		idsToDONs[DonID(d.Id)] = DON{
			DON:                      *toDONInfo(d),
			CapabilityConfigurations: cc,
		}
	}

	nodes := []kcr.CapabilitiesRegistryNodeInfo{}
	err = s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getNodes", primitives.Unconfirmed, nil, &nodes)
	if err != nil {
		return nil, err
	}

	idsToNodes := map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{}
	for _, node := range nodes {
		idsToNodes[node.P2pId] = node
	}

	return &LocalRegistry{
		lggr:              s.lggr,
		peerWrapper:       s.peerWrapper,
		IDsToDONs:         idsToDONs,
		IDsToCapabilities: idsToCapabilities,
		IDsToNodes:        idsToNodes,
	}, nil
}

func (s *registrySyncer) sync(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.launchers) == 0 {
		s.lggr.Warn("sync called, but no launchers are registered; nooping")
		return nil
	}

	if s.reader == nil {
		reader, err := s.initReader(ctx, s.lggr, s.relayer, s.registryAddress)
		if err != nil {
			return err
		}

		s.reader = reader
	}

	lr, err := s.localRegistry(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync with remote registry: %w", err)
	}

	for _, h := range s.launchers {
		if err := h.Launch(ctx, lr); err != nil {
			s.lggr.Errorf("error calling launcher: %s", err)
		}
	}

	return nil
}

func toCapabilityType(capabilityType uint8) capabilities.CapabilityType {
	switch capabilityType {
	case 0:
		return capabilities.CapabilityTypeTrigger
	case 1:
		return capabilities.CapabilityTypeAction
	case 2:
		return capabilities.CapabilityTypeConsensus
	case 3:
		return capabilities.CapabilityTypeTarget
	default:
		// Not found
		return capabilities.CapabilityType(-1)
	}
}

func toDONInfo(don kcr.CapabilitiesRegistryDONInfo) *capabilities.DON {
	peerIDs := []p2ptypes.PeerID{}
	for _, p := range don.NodeP2PIds {
		peerIDs = append(peerIDs, p)
	}

	return &capabilities.DON{
		ID:               don.Id,
		ConfigVersion:    don.ConfigCount,
		Members:          peerIDs,
		F:                don.F,
		IsPublic:         don.IsPublic,
		AcceptsWorkflows: don.AcceptsWorkflows,
	}
}

func (s *registrySyncer) AddLauncher(launchers ...Launcher) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.launchers = append(s.launchers, launchers...)
}

func (s *registrySyncer) Close() error {
	return s.StopOnce("RegistrySyncer", func() error {
		close(s.stopCh)
		s.wg.Wait()
		return nil
	})
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
