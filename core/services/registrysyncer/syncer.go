package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type HashedCapabilityID [32]byte
type DonID uint32

type State struct {
	IDsToDONs         map[DonID]kcr.CapabilitiesRegistryDONInfo
	IDsToNodes        map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo
	IDsToCapabilities map[HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo
}

type Launcher interface {
	Launch(ctx context.Context, state State) error
}

type Syncer interface {
	services.Service
	AddLauncher(h ...Launcher)
}

type registrySyncer struct {
	stopCh    services.StopChan
	launchers []Launcher
	reader    types.ContractReader

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

func (s *registrySyncer) state(ctx context.Context) (State, error) {
	dons := []kcr.CapabilitiesRegistryDONInfo{}
	err := s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getDONs", nil, &dons)
	if err != nil {
		return State{}, err
	}

	idsToDONs := map[DonID]kcr.CapabilitiesRegistryDONInfo{}
	for _, d := range dons {
		idsToDONs[DonID(d.Id)] = d
	}

	caps := []kcr.CapabilitiesRegistryCapabilityInfo{}
	err = s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getCapabilities", nil, &caps)
	if err != nil {
		return State{}, err
	}

	idsToCapabilities := map[HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{}
	for _, c := range caps {
		idsToCapabilities[c.HashedId] = c
	}

	nodes := []kcr.CapabilitiesRegistryNodeInfo{}
	err = s.reader.GetLatestValue(ctx, "CapabilitiesRegistry", "getNodes", nil, &nodes)
	if err != nil {
		return State{}, err
	}

	idsToNodes := map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{}
	for _, node := range nodes {
		idsToNodes[node.P2pId] = node
	}

	return State{IDsToDONs: idsToDONs, IDsToCapabilities: idsToCapabilities, IDsToNodes: idsToNodes}, nil
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
