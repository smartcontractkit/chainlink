package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

type Syncer interface {
	services.Service
	AddListener(h ...registrysyncer.Listener)
}

type ContractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

type RegistrySyncer interface {
	Sync(ctx context.Context, isInitialSync bool) error
	AddListener(listeners ...registrysyncer.Listener)
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
	listeners            []registrysyncer.Listener
	reader               types.ContractReader
	initReader           func(ctx context.Context, lggr logger.Logger, relayer ContractReaderFactory, capabilitiesContract types.BoundContract) (types.ContractReader, error)
	relayer              ContractReaderFactory
	capabilitiesContract types.BoundContract
	getPeerID            func() (p2ptypes.PeerID, error)

	orm registrysyncer.ORM

	updateChan chan *registrysyncer.LocalRegistry

	wg   sync.WaitGroup
	lggr logger.Logger
	mu   sync.RWMutex
}

var _ services.Service = &registrySyncer{}

var defaultTickInterval = 12 * time.Second

// New instantiates a new RegistrySyncer
func New(
	lggr logger.Logger,
	getPeerID func() (p2ptypes.PeerID, error),
	relayer ContractReaderFactory,
	registryAddress string,
	orm registrysyncer.ORM,
) (RegistrySyncer, error) {
	metricLabeler, err := newSyncerMetricLabeler()
	if err != nil {
		return nil, fmt.Errorf("failed to create syncer metric labeler: %w", err)
	}

	return &registrySyncer{
		metrics:    metricLabeler,
		stopCh:     make(services.StopChan),
		updateChan: make(chan *registrysyncer.LocalRegistry),
		lggr:       logger.Named(lggr, "RegistrySyncer"),
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
	contractReaderConfigEncoded, err := json.Marshal(buildV2ContractReaderConfig())
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

// buildV2ContractReaderConfig creates the contract reader configuration for V2 capabilities registry
func buildV2ContractReaderConfig() config.ChainReaderConfig {
	return config.ChainReaderConfig{
		Contracts: map[string]config.ChainContractReader{
			"CapabilitiesRegistry": {
				ContractABI: capabilities_registry_v2.CapabilitiesRegistryABI,
				Configs: map[string]*config.ChainReaderDefinition{
					"getDONs": {
						ChainSpecificName: "getDONs",
					},
					"getCapabilities": {
						ChainSpecificName: "getCapabilities",
					},
					"getNodes": {
						ChainSpecificName: "getNodes",
					},
					"getDONsInFamily": {
						ChainSpecificName: "getDONsInFamily",
					},
					"getHistoricalDONInfo": {
						ChainSpecificName: "getHistoricalDONInfo",
					},
					"getNode": {
						ChainSpecificName: "getNode",
					},
					"getNodeOperator": {
						ChainSpecificName: "getNodeOperator",
					},
					"getNodeOperators": {
						ChainSpecificName: "getNodeOperators",
					},
					"getNodesByP2PIds": {
						ChainSpecificName: "getNodesByP2PIds",
					},
					"isCapabilityDeprecated": {
						ChainSpecificName: "isCapabilityDeprecated",
					},
				},
			},
		},
	}
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

func (s *registrySyncer) importOnchainRegistry(ctx context.Context) (*registrysyncer.LocalRegistry, error) {
	caps := []capabilities_registry_v2.CapabilitiesRegistryCapabilityInfo{}
	// TODO support pagination if needed
	// Using large limit for now to avoid pagination complexity
	// since we don't expect to have that many capabilities
	params := struct {
		Start *big.Int
		Limit *big.Int
	}{Start: big.NewInt(0), Limit: big.NewInt(1024)}
	err := s.reader.GetLatestValue(ctx, s.capabilitiesContract.ReadIdentifier("getCapabilities"), primitives.Unconfirmed, params, &caps)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest value for getCapabilities: %w", err)
	}

	idsToCapabilities := map[string]registrysyncer.Capability{}
	for _, c := range caps {
		capabilityType, _, parseErr := parseCapabilityMetadata(c.Metadata)
		if parseErr != nil {
			s.lggr.Warnw("failed to parse capability metadata, skipping", "capabilityID", c.CapabilityId, "error", parseErr)
			continue
		}
		idsToCapabilities[c.CapabilityId] = registrysyncer.Capability{
			ID:             c.CapabilityId,
			CapabilityType: toCapabilityType(capabilityType),
		}
	}

	dons := []capabilities_registry_v2.CapabilitiesRegistryDONInfo{}

	err = s.reader.GetLatestValue(ctx, s.capabilitiesContract.ReadIdentifier("getDONs"), primitives.Unconfirmed, params, &dons)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest value for getDONs: %w", err)
	}

	idsToDONs := map[registrysyncer.DonID]registrysyncer.DON{}
	for _, d := range dons {
		cc := map[string]registrysyncer.CapabilityConfiguration{}
		for _, dc := range d.CapabilityConfigurations {
			cc[dc.CapabilityId] = registrysyncer.CapabilityConfiguration{
				Config: dc.Config,
			}
		}

		idsToDONs[registrysyncer.DonID(d.Id)] = registrysyncer.DON{
			DON:                      *toDONInfo(d),
			CapabilityConfigurations: cc,
		}
	}

	nodes := []capabilities_registry_v2.INodeInfoProviderNodeInfo{}

	err = s.reader.GetLatestValue(ctx, s.capabilitiesContract.ReadIdentifier("getNodes"), primitives.Unconfirmed, params, &nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest value for getNodes: %w", err)
	}

	idsToNodes := map[p2ptypes.PeerID]registrysyncer.NodeInfo{}
	for _, node := range nodes {
		nodeInfo := registrysyncer.NodeInfo{
			NodeOperatorID:      node.NodeOperatorId,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			Signer:              node.Signer,
			P2pID:               node.P2pId,
			EncryptionPublicKey: node.EncryptionPublicKey,
			CapabilitiesDONIds:  make([]*big.Int, 0, len(node.CapabilitiesDONIds)),
			HashedCapabilityIDs: make([][32]byte, 0, len(node.CapabilityIds)),
			CsaKey:              node.CsaKey,
			CapabilityIDs:       node.CapabilityIds,
		}
		copy(nodeInfo.CapabilitiesDONIds, node.CapabilitiesDONIds)

		// Backfill hashed capability IDs
		for _, capID := range node.CapabilityIds {
			hashedCapID, err := HashCapabilityID(capID)
			if err != nil {
				s.lggr.Warnw("failed to hash capability ID, skipping", "capabilityID", capID, "error", err)
				continue
			}
			nodeInfo.HashedCapabilityIDs = append(nodeInfo.HashedCapabilityIDs, hashedCapID)
		}

		idsToNodes[node.P2pId] = nodeInfo
	}

	return &registrysyncer.LocalRegistry{
		Logger:            s.lggr,
		GetPeerID:         s.getPeerID,
		IDsToDONs:         idsToDONs,
		IDsToCapabilities: idsToCapabilities,
		IDsToNodes:        idsToNodes,
	}, nil
}

func (s *registrySyncer) Sync(ctx context.Context, isInitialSync bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.listeners) == 0 {
		s.lggr.Warn("sync called, but no listeners are registered; nooping")
		return nil
	}

	if s.reader == nil {
		reader, err := s.initReader(ctx, s.lggr, s.relayer, s.capabilitiesContract)
		if err != nil {
			return err
		}

		s.reader = reader
	}

	var latestRegistry *registrysyncer.LocalRegistry
	var err error

	if isInitialSync {
		s.lggr.Debug("syncing with local registry")
		latestRegistry, err = s.orm.LatestLocalRegistry(ctx)
		if err != nil {
			s.lggr.Warnw("failed to sync with local registry, using remote registry instead", "error", err)
		} else {
			latestRegistry.Logger = s.lggr
			latestRegistry.GetPeerID = s.getPeerID
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
			return nil
		case s.updateChan <- latestRegistry:
			// Successfully sent state
			s.lggr.Debug("remote registry update triggered successfully")
		default:
			// No one is ready to receive the state, handle accordingly
			s.lggr.Debug("no listeners on update channel, remote registry update skipped")
		}
	}

	for _, listener := range s.listeners {
		lrCopy := registrysyncer.DeepCopyLocalRegistry(latestRegistry)
		if err := listener.OnNewRegistry(ctx, &lrCopy); err != nil {
			s.lggr.Errorf("error calling launcher: %s", err)
			s.metrics.incrementLauncherFailureCounter(ctx)
		}
	}

	return nil
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

func toDONInfo(don capabilities_registry_v2.CapabilitiesRegistryDONInfo) *capabilities.DON {
	peerIDs := []p2ptypes.PeerID{}
	for _, p := range don.NodeP2PIds {
		peerIDs = append(peerIDs, p)
	}

	return &capabilities.DON{
		Name:             don.Name,
		ID:               don.Id,
		Families:         don.DonFamilies,
		ConfigVersion:    don.ConfigCount,
		Members:          peerIDs,
		F:                don.F,
		IsPublic:         don.IsPublic,
		AcceptsWorkflows: don.AcceptsWorkflows,
		Config:           don.Config,
	}
}

func (s *registrySyncer) AddListener(listeners ...registrysyncer.Listener) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners = append(s.listeners, listeners...)
}

func (s *registrySyncer) Close() error {
	return s.StopOnce("RegistrySyncer", func() error {
		close(s.stopCh)
		s.wg.Wait()
		close(s.updateChan)
		return nil
	})
}

func (s *registrySyncer) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *registrySyncer) Name() string {
	return s.lggr.Name()
}

// CapabilityMetadata represents the metadata structure for V2 capabilities
type CapabilityMetadata struct {
	CapabilityType uint8 `json:"capabilityType"`
	ResponseType   uint8 `json:"responseType"`
}

// parseCapabilityMetadata extracts capability type and response type from V2 metadata
func parseCapabilityMetadata(metadata []byte) (capabilityType, responseType uint8, err error) {
	if len(metadata) == 0 {
		return 0, 0, errors.New("metadata is empty")
	}

	var meta CapabilityMetadata
	if err := json.Unmarshal(metadata, &meta); err != nil {
		return 0, 0, fmt.Errorf("invalid metadata: %w", err)
	}

	return meta.CapabilityType, meta.ResponseType, nil
}

// parseCapabilityID parses a V2 capability ID (e.g., "write-chain@1.0.1") into name and version parts
func parseCapabilityID(capabilityID string) (name, version string, err error) {
	parts := strings.Split(capabilityID, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid capability ID format: %s (expected format: name@version)", capabilityID)
	}
	return parts[0], parts[1], nil
}

// hashCapabilityID creates a hashed capability ID from a V2 capability ID string
func HashCapabilityID(capabilityID string) ([32]byte, error) {
	name, version, err := parseCapabilityID(capabilityID)
	if err != nil {
		return [32]byte{}, err
	}

	return common.HashedCapabilityID(name, version)
}
