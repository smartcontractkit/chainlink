package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
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

type contractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
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

	orm *syncerORM

	updateChan     chan *LocalRegistry
	testUpdateChan chan bool

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
	ds sqlutil.DataSource,
) (*registrySyncer, error) {
	orm := newORM(ds, lggr)
	return &registrySyncer{
		stopCh:          make(services.StopChan),
		updateChan:      make(chan *LocalRegistry),
		lggr:            lggr.Named("RegistrySyncer"),
		relayer:         relayer,
		registryAddress: registryAddress,
		initReader:      newReader,
		orm:             &orm,
		peerWrapper:     peerWrapper,
	}, nil
}

func newTestSyncer(
	lggr logger.Logger,
	peerWrapper p2ptypes.PeerWrapper,
	relayer contractReaderFactory,
	registryAddress string,
	ds sqlutil.DataSource,
) (*registrySyncer, error) {
	rs, err := New(lggr, peerWrapper, relayer, registryAddress, ds)
	if err != nil {
		return nil, err
	}
	rs.testUpdateChan = make(chan bool, 100)
	return rs, nil
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
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.updateStateLoop()
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
	err := s.sync(ctx, true)
	if err != nil {
		s.lggr.Errorw("failed to sync with remote registry", "error", err)
	}

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.lggr.Debug("starting regular sync with the remote registry")
			err := s.sync(ctx, false)
			if err != nil {
				s.lggr.Errorw("failed to sync with remote registry", "error", err)
			}
		}
	}
}

func (s *registrySyncer) updateStateLoop() {
	ctx, cancel := s.stopCh.NewCtx()
	defer cancel()

	for {
		select {
		case <-s.stopCh:
			return
		case localRegistry, ok := <-s.updateChan:
			if !ok {
				// channel has been closed, terminating.
				return
			}
			if err := s.orm.addLocalRegistry(ctx, *localRegistry); err != nil {
				s.lggr.Errorw("failed to save state to local registry", "error", err)
			}
			select {
			case s.testUpdateChan <- true:
			default:
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

	var remoteTriggerConfig *capabilities.RemoteTriggerConfig
	var remoteTargetConfig *capabilities.RemoteTargetConfig

	switch cconf.GetRemoteConfig().(type) {
	case *capabilitiespb.CapabilityConfig_RemoteTriggerConfig:
		prtc := cconf.GetRemoteTriggerConfig()
		remoteTriggerConfig = &capabilities.RemoteTriggerConfig{}
		remoteTriggerConfig.RegistrationRefresh = prtc.RegistrationRefresh.AsDuration()
		remoteTriggerConfig.RegistrationExpiry = prtc.RegistrationExpiry.AsDuration()
		remoteTriggerConfig.MinResponsesToAggregate = prtc.MinResponsesToAggregate
		remoteTriggerConfig.MessageExpiry = prtc.MessageExpiry.AsDuration()
	case *capabilitiespb.CapabilityConfig_RemoteTargetConfig:
		prtc := cconf.GetRemoteTargetConfig()
		remoteTargetConfig = &capabilities.RemoteTargetConfig{}
		remoteTargetConfig.RequestHashExcludedAttributes = prtc.RequestHashExcludedAttributes
	}

	dc, err := values.FromMapValueProto(cconf.DefaultConfig)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}

	return capabilities.CapabilityConfiguration{
		DefaultConfig:       dc,
		RemoteTriggerConfig: remoteTriggerConfig,
		RemoteTargetConfig:  remoteTargetConfig,
	}, nil
}

func (s *registrySyncer) localRegistry(ctx context.Context) (*LocalRegistry, error) {
	var caps []kcr.CapabilitiesRegistryCapabilityInfo
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

	var dons []kcr.CapabilitiesRegistryDONInfo
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

			cc[cid] = cconf
		}

		idsToDONs[DonID(d.Id)] = DON{
			DON:                      *toDONInfo(d),
			CapabilityConfigurations: cc,
		}
	}

	var nodes []kcr.CapabilitiesRegistryNodeInfo
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

func (s *registrySyncer) sync(ctx context.Context, isInitialSync bool) error {
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

	var lr *LocalRegistry
	var err error

	if isInitialSync {
		s.lggr.Debug("syncing with local registry")
		lr, err = s.orm.latestLocalRegistry(ctx)
		if err != nil {
			s.lggr.Errorw("failed to sync with local registry, using remote registry instead", "error", err)
		} else {
			lr.lggr = s.lggr
			lr.peerWrapper = s.peerWrapper
		}
	}

	if lr == nil {
		s.lggr.Debug("syncing with remote registry")
		localRegistry, err := s.localRegistry(ctx)
		if err != nil {
			return fmt.Errorf("failed to sync with remote registry: %w", err)
		}
		lr = localRegistry
		// Attempt to send local registry to the update channel without blocking
		// This is to prevent the tests from hanging if they are not calling `Start()` on the syncer
		select {
		case s.updateChan <- lr:
			// Successfully sent state
			s.lggr.Debug("remote registry update triggered successfully")
		default:
			// No one is ready to receive the state, handle accordingly
			s.lggr.Debug("no listeners on update channel, remote registry update skipped")
		}
	}

	for _, h := range s.launchers {
		lrCopy := deepCopyLocalRegistry(lr)
		if err := h.Launch(ctx, &lrCopy); err != nil {
			s.lggr.Errorf("error calling launcher: %s", err)
		}
	}

	return nil
}

func deepCopyLocalRegistry(lr *LocalRegistry) LocalRegistry {
	var lrCopy LocalRegistry
	lrCopy.lggr = lr.lggr
	lrCopy.peerWrapper = lr.peerWrapper
	lrCopy.IDsToDONs = make(map[DonID]DON, len(lr.IDsToDONs))
	for id, don := range lr.IDsToDONs {
		d := capabilities.DON{
			ID:               don.ID,
			ConfigVersion:    don.ConfigVersion,
			Members:          make([]p2ptypes.PeerID, len(don.Members)),
			F:                don.F,
			IsPublic:         don.IsPublic,
			AcceptsWorkflows: don.AcceptsWorkflows,
		}
		copy(d.Members, don.Members)
		capCfgs := make(map[string]capabilities.CapabilityConfiguration, len(don.CapabilityConfigurations))
		for capID, capCfg := range don.CapabilityConfigurations {
			cfg := capabilities.CapabilityConfiguration{
				DefaultConfig:       nil,
				RemoteTriggerConfig: nil,
				RemoteTargetConfig:  nil,
			}
			if capCfg.DefaultConfig != nil {
				defaultConfig := *capCfg.DefaultConfig
				cfg.DefaultConfig = &defaultConfig
			}
			if capCfg.RemoteTriggerConfig != nil {
				remoteTriggerConfig := *capCfg.RemoteTriggerConfig
				cfg.RemoteTriggerConfig = &remoteTriggerConfig
			}
			if capCfg.RemoteTargetConfig != nil {
				remoteTargetConfig := *capCfg.RemoteTargetConfig
				cfg.RemoteTargetConfig = &remoteTargetConfig
			}
			capCfgs[capID] = cfg
		}
		lrCopy.IDsToDONs[id] = DON{
			DON:                      d,
			CapabilityConfigurations: capCfgs,
		}
	}

	lrCopy.IDsToCapabilities = make(map[string]Capability, len(lr.IDsToCapabilities))
	for id, capability := range lr.IDsToCapabilities {
		cp := capability
		lrCopy.IDsToCapabilities[id] = cp
	}

	lrCopy.IDsToNodes = make(map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo, len(lr.IDsToNodes))
	for id, node := range lr.IDsToNodes {
		nodeInfo := kcr.CapabilitiesRegistryNodeInfo{
			NodeOperatorId:      node.NodeOperatorId,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			Signer:              node.Signer,
			P2pId:               node.P2pId,
			HashedCapabilityIds: make([][32]byte, len(node.HashedCapabilityIds)),
			CapabilitiesDONIds:  make([]*big.Int, len(node.CapabilitiesDONIds)),
		}
		copy(nodeInfo.HashedCapabilityIds, node.HashedCapabilityIds)
		copy(nodeInfo.CapabilitiesDONIds, node.CapabilitiesDONIds)
		lrCopy.IDsToNodes[id] = nodeInfo
	}

	return lrCopy
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
		close(s.updateChan)
		if s.testUpdateChan != nil {
			close(s.testUpdateChan)
		}
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
