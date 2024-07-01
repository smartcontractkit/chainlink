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
	stopCh      services.StopChan
	launchers   []Launcher
	reader      types.ContractReader
	peerWrapper p2ptypes.PeerWrapper

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
	ctx, _ := stopCh.NewCtx()
	reader, err := newReader(ctx, lggr, relayer, registryAddress)
	if err != nil {
		return nil, err
	}

	return newSyncer(
		stopCh,
		lggr.Named("RegistrySyncer"),
		reader,
	), nil
}

type contractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

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

func newSyncer(
	stopCh services.StopChan,
	lggr logger.Logger,
	reader types.ContractReader,
) *registrySyncer {
	return &registrySyncer{
		stopCh: stopCh,
		lggr:   lggr,
		reader: reader,
	}
}

func (s *registrySyncer) Start(ctx context.Context) error {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.syncLoop()
	}()
	return nil
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
		rtc.RegistrationRefreshMs = prtc.RegistrationRefreshMs
		rtc.RegistrationExpiryMs = prtc.RegistrationExpiryMs
		rtc.MinResponsesToAggregate = prtc.MinResponsesToAggregate
		rtc.MessageExpiryMs = prtc.MessageExpiryMs
	}

	return capabilities.CapabilityConfiguration{
		ExecuteConfig:       values.FromMapValueProto(cconf.ExecuteConfig),
		RemoteTriggerConfig: rtc,
	}, nil
}

func (s *registrySyncer) state(ctx context.Context) (*LocalRegistry, error) {
	caps := []kcr.CapabilitiesRegistryCapabilityInfo{}
	err := s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getCapabilities", nil, &caps)
	if err != nil {
		return nil, err
	}

	idsToCapabilities := map[CapabilityID]Capability{}
	hashedIDsToCapabilityIDs := map[[32]byte]CapabilityID{}
	for _, c := range caps {
		cid := CapabilityID(fmt.Sprintf("%s@%s", c.LabelledName, c.Version))
		idsToCapabilities[cid] = Capability{
			ID:             cid,
			CapabilityType: toCapabilityType(c.CapabilityType),
		}

		hashedIDsToCapabilityIDs[c.HashedId] = cid
	}

	dons := []kcr.CapabilitiesRegistryDONInfo{}
	err = s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getDONs", nil, &dons)
	if err != nil {
		return nil, err
	}

	idsToDONs := map[DonID]DON{}
	for _, d := range dons {
		cc := map[CapabilityID]capabilities.CapabilityConfiguration{}
		for _, dc := range d.CapabilityConfigurations {
			cid, ok := hashedIDsToCapabilityIDs[dc.CapabilityId]
			if !ok {
				return nil, fmt.Errorf("invariant violation: could not find full ID for hashed ID %s", dc.CapabilityId)
			}

			cconf, innerErr := unmarshalCapabilityConfig(dc.Config)
			if err != nil {
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
	err = s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getNodes", nil, &nodes)
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

	state, err := s.state(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync with remote registry: %w", err)
	}

	for _, h := range s.launchers {
		if err := h.Launch(ctx, state); err != nil {
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
		ID:            don.Id,
		ConfigVersion: don.ConfigCount,
		Members:       peerIDs,
		F:             don.F,
	}
}

func (s *registrySyncer) AddLauncher(launchers ...Launcher) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.launchers = append(s.launchers, launchers...)
}

func (s *registrySyncer) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	return nil
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
