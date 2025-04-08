package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
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

type ContractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

type RegistrySyncer interface {
	Sync(ctx context.Context, isInitialSync bool) error
	AddLauncher(launchers ...Launcher)
	Start(ctx context.Context) error
	Close() error
	Ready() error
	HealthReport() map[string]error
	Name() string
}

type registrySyncer struct {
	services.StateMachine
	metrics              *syncerMetricLabeler
	stopCh               services.StopChan
	launchers            []Launcher
	reader               types.ContractReader
	initReader           func(ctx context.Context, lggr logger.Logger, relayer ContractReaderFactory, capabilitiesContract types.BoundContract) (types.ContractReader, error)
	relayer              ContractReaderFactory
	capabilitiesContract types.BoundContract
	getPeerID            func() (p2ptypes.PeerID, error)

	orm ORM

	updateChan chan *LocalRegistry

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
	getPeerID func() (p2ptypes.PeerID, error),
	relayer ContractReaderFactory,
	registryAddress string,
	orm ORM,
) (RegistrySyncer, error) {
	metricLabeler, err := newSyncerMetricLabeler()
	if err != nil {
		return nil, fmt.Errorf("failed to create syncer metric labeler: %w", err)
	}

	return &registrySyncer{
		metrics:    metricLabeler,
		stopCh:     make(services.StopChan),
		updateChan: make(chan *LocalRegistry),
		lggr:       lggr.Named("RegistrySyncer"),
		relayer:    relayer,
		capabilitiesContract: types.BoundContract{
			Address: registryAddress,
			Name:    "CapabilitiesRegistry",
		},
		initReader: newReader,
		orm:        orm,
		getPeerID:  getPeerID,
	}, nil
}

// NOTE: this can't be called while initializing the syncer and needs to be called in the sync loop.
// This is because Bind() makes an onchain call to verify that the contract address exists, and if
// called during initialization, this results in a "no live nodes" error.
func newReader(ctx context.Context, lggr logger.Logger, relayer ContractReaderFactory, capabilitiesContract types.BoundContract) (types.ContractReader, error) {
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

	err = cr.Bind(ctx, []types.BoundContract{capabilitiesContract})

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
	// fire for the first time at T+N, where N is the interval. We do not
	// increment RemoteRegistryFailureCounter the first time
	s.lggr.Debug("starting initial sync with remote registry")
	err := s.Sync(ctx, true)
	if err != nil {
		s.lggr.Errorw("failed to sync with remote registry", "error", err)
	}

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.lggr.Debug("starting regular sync with the remote registry")
			err := s.Sync(ctx, false)
			if err != nil {
				s.lggr.Errorw("failed to sync with remote registry", "error", err)
				s.metrics.incrementRemoteRegistryFailureCounter(ctx)
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
			if err := s.orm.AddLocalRegistry(ctx, *localRegistry); err != nil {
				s.lggr.Errorw("failed to save state to local registry", "error", err)
			}
		}
	}
}

func (s *registrySyncer) importOnchainRegistry(ctx context.Context) (*LocalRegistry, error) {
	caps := []kcr.CapabilitiesRegistryCapabilityInfo{}

	err := s.reader.GetLatestValue(ctx, s.capabilitiesContract.ReadIdentifier("getCapabilities"), primitives.Unconfirmed, nil, &caps)
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

	err = s.reader.GetLatestValue(ctx, s.capabilitiesContract.ReadIdentifier("getDONs"), primitives.Unconfirmed, nil, &dons)
	if err != nil {
		return nil, err
	}

	idsToDONs := map[DonID]DON{}
	for _, d := range dons {
		cc := map[string]CapabilityConfiguration{}
		for _, dc := range d.CapabilityConfigurations {
			cid, ok := hashedIDsToCapabilityIDs[dc.CapabilityId]
			if !ok {
				return nil, fmt.Errorf("invariant violation: could not find full ID for hashed ID %s", dc.CapabilityId)
			}

			cc[cid] = CapabilityConfiguration{
				Config: dc.Config,
			}
		}

		idsToDONs[DonID(d.Id)] = DON{
			DON:                      *toDONInfo(d),
			CapabilityConfigurations: cc,
		}
	}

	nodes := []kcr.INodeInfoProviderNodeInfo{}

	err = s.reader.GetLatestValue(ctx, s.capabilitiesContract.ReadIdentifier("getNodes"), primitives.Unconfirmed, nil, &nodes)
	if err != nil {
		return nil, err
	}

	idsToNodes := map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{}
	for _, node := range nodes {
		idsToNodes[node.P2pId] = node
	}

	return &LocalRegistry{
		lggr:              s.lggr,
		getPeerID:         s.getPeerID,
		IDsToDONs:         idsToDONs,
		IDsToCapabilities: idsToCapabilities,
		IDsToNodes:        idsToNodes,
	}, nil
}

func (s *registrySyncer) Sync(ctx context.Context, isInitialSync bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.launchers) == 0 {
		s.lggr.Warn("sync called, but no launchers are registered; nooping")
		return nil
	}

	if s.reader == nil {
		reader, err := s.initReader(ctx, s.lggr, s.relayer, s.capabilitiesContract)
		if err != nil {
			return err
		}

		s.reader = reader
	}

	var latestRegistry *LocalRegistry
	var err error

	if isInitialSync {
		s.lggr.Debug("syncing with local registry")
		latestRegistry, err = s.orm.LatestLocalRegistry(ctx)
		if err != nil {
			s.lggr.Warnw("failed to sync with local registry, using remote registry instead", "error", err)
		} else {
			latestRegistry.lggr = s.lggr
			latestRegistry.getPeerID = s.getPeerID
		}
	}

	if latestRegistry == nil {
		s.lggr.Debug("syncing with remote registry")
		importedRegistry, err := s.importOnchainRegistry(ctx)
		if err != nil {
			return fmt.Errorf("failed to sync with remote registry: %w", err)
		}
		latestRegistry = importedRegistry
		// Attempt to send local registry to the update channel without blocking
		// This is to prevent the tests from hanging if they are not calling `Start()` on the syncer
		select {
		case <-s.stopCh:
			s.lggr.Debug("sync cancelled, stopping")
		case s.updateChan <- latestRegistry:
			// Successfully sent state
			s.lggr.Debug("remote registry update triggered successfully")
		default:
			// No one is ready to receive the state, handle accordingly
			s.lggr.Debug("no listeners on update channel, remote registry update skipped")
		}
	}

	for _, h := range s.launchers {
		lrCopy := deepCopyLocalRegistry(latestRegistry)
		if err := h.Launch(ctx, &lrCopy); err != nil {
			s.lggr.Errorf("error calling launcher: %s", err)
			s.metrics.incrementLauncherFailureCounter(ctx)
		}
	}

	return nil
}

func deepCopyLocalRegistry(lr *LocalRegistry) LocalRegistry {
	var lrCopy LocalRegistry
	lrCopy.lggr = lr.lggr
	lrCopy.getPeerID = lr.getPeerID
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
		capCfgs := make(map[string]CapabilityConfiguration, len(don.CapabilityConfigurations))
		for capID, capCfg := range don.CapabilityConfigurations {
			capCfgs[capID] = CapabilityConfiguration{
				Config: capCfg.Config,
			}
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

	lrCopy.IDsToNodes = make(map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo, len(lr.IDsToNodes))
	for id, node := range lr.IDsToNodes {
		nodeInfo := kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      node.NodeOperatorId,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			Signer:              node.Signer,
			P2pId:               node.P2pId,
			EncryptionPublicKey: node.EncryptionPublicKey,
			HashedCapabilityIds: make([][32]byte, len(node.HashedCapabilityIds)),
			CapabilitiesDONIds:  make([]*big.Int, len(node.CapabilitiesDONIds)),
		}
		copy(nodeInfo.HashedCapabilityIds, node.HashedCapabilityIds)
		copy(nodeInfo.CapabilitiesDONIds, node.CapabilitiesDONIds)
		lrCopy.IDsToNodes[id] = nodeInfo
	}

	return lrCopy
}

type ContractCapabilityType uint8

const (
	ContractCapabilityTypeTrigger ContractCapabilityType = iota
	ContractCapabilityTypeAction
	ContractCapabilityTypeConsensus
	ContractCapabilityTypeTarget
)

func toCapabilityType(capabilityType uint8) capabilities.CapabilityType {
	switch ContractCapabilityType(capabilityType) {
	case ContractCapabilityTypeTrigger:
		return capabilities.CapabilityTypeTrigger
	case ContractCapabilityTypeAction:
		return capabilities.CapabilityTypeAction
	case ContractCapabilityTypeConsensus:
		return capabilities.CapabilityTypeConsensus
	case ContractCapabilityTypeTarget:
		return capabilities.CapabilityTypeTarget
	default:
		return capabilities.CapabilityTypeUnknown
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
		s.mu.Lock()
		defer s.mu.Unlock()
		close(s.updateChan)
		s.wg.Wait()
		return nil
	})
}

func (s *registrySyncer) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *registrySyncer) Name() string {
	return s.lggr.Name()
}
