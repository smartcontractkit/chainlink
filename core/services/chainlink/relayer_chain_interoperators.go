package chainlink

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
)

var ErrNoSuchRelayer = errors.New("relayer does not exist")

// RelayerChainInteroperators
// encapsulates relayers and chains and is the primary entry point for
// the node to access relayers, get legacy chains associated to a relayer
// and get status about the chains and nodes
type RelayerChainInteroperators interface {
	Services() []services.ServiceCtx

	List(filter FilterFn) RelayerChainInteroperators

	LoopRelayerStorer
	LegacyChainer
	ChainsNodesStatuser
}

// LoopRelayerStorer is key-value like interface for storing and
// retrieving [loop.Relayer]
type LoopRelayerStorer interface {
	ocr2.RelayGetter
	Slice() []loop.Relayer
}

// LegacyChainer is an interface for getting legacy chains
// This will be deprecated/removed when products depend only
// on the relayer interface.
type LegacyChainer interface {
	LegacyEVMChains() legacyevm.LegacyChainContainer
}

// NetworkChainStatus is a ChainStatus from a particlar Network.
type NetworkChainStatus struct {
	Network string
	types.ChainStatus
}

type ChainStatuser interface {
	ChainStatus(ctx context.Context, id types.RelayID) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]NetworkChainStatus, int, error)
}

// NodesStatuser is an interface for node configuration and state.
// TODO BCF-2440, BCF-2511 may need Node(ctx,name) to get a node status by name
type NodesStatuser interface {
	NodeStatuses(ctx context.Context, offset, limit int, relayIDs ...types.RelayID) (nodes []types.NodeStatus, count int, err error)
}

// ChainsNodesStatuser report statuses about chains and nodes
type ChainsNodesStatuser interface {
	ChainStatuser
	NodesStatuser
}

var _ RelayerChainInteroperators = &CoreRelayerChainInteroperators{}

type DummyFactory interface {
	NewDummy(config DummyFactoryConfig) (loop.Relayer, error)
}

// CoreRelayerChainInteroperators implements [RelayerChainInteroperators]
// as needed for the core [chainlink.Application]
type CoreRelayerChainInteroperators struct {
	mu           sync.Mutex
	loopRelayers map[types.RelayID]loop.Relayer
	legacyChains legacyChains

	dummyFactory DummyFactory

	// we keep an explicit list of services because the legacy implementations have more than
	// just the relayer service
	srvs []services.ServiceCtx
}

func NewCoreRelayerChainInteroperators(initFuncs ...CoreRelayerChainInitFunc) (*CoreRelayerChainInteroperators, error) {
	cr := &CoreRelayerChainInteroperators{
		loopRelayers: make(map[types.RelayID]loop.Relayer),
		srvs:         make([]services.ServiceCtx, 0),
	}
	for _, initFn := range initFuncs {
		err := initFn(cr)
		if err != nil {
			return nil, err
		}
	}
	return cr, nil
}

// CoreRelayerChainInitFunc is a hook in the constructor to create relayers from a factory.
type CoreRelayerChainInitFunc func(op *CoreRelayerChainInteroperators) error

// InitDummy instantiates a dummy relayer
func InitDummy(factory RelayerFactory) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) error {
		op.dummyFactory = &factory
		return nil
	}
}

// InitEVM is a option for instantiating evm relayers
func InitEVM(factory RelayerFactory, config EVMFactoryConfig) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) (err error) {
		adapters, err2 := factory.NewEVM(config)
		if err2 != nil {
			return fmt.Errorf("failed to setup EVM relayer: %w", err2)
		}

		legacyMap := make(map[string]legacyevm.Chain)
		for id, a := range adapters {
			// adapter is a service
			op.srvs = append(op.srvs, a)
			op.loopRelayers[id] = a
			legacyMap[id.ChainID] = a.Chain()
		}
		op.legacyChains.EVMChains = legacyevm.NewLegacyChains(legacyMap, config.ChainConfigs)
		return nil
	}
}

// InitCosmos is a option for instantiating Cosmos relayers
func InitCosmos(factory RelayerFactory, ks keystore.Cosmos, chainCfgs RawConfigs) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) (err error) {
		loopKs := &keystore.CosmosLoopSigner{Cosmos: ks}
		relayers, err := factory.NewCosmos(loopKs, chainCfgs)
		if err != nil {
			return fmt.Errorf("failed to setup Cosmos relayer: %w", err)
		}
		for id, relayer := range relayers {
			op.srvs = append(op.srvs, relayer)
			op.loopRelayers[id] = relayer
		}

		return nil
	}
}

// InitSolana is a option for instantiating Solana relayers
func InitSolana(factory RelayerFactory, ks keystore.Solana, config SolanaFactoryConfig) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) error {
		loopKs := &keystore.SolanaLooppSigner{Solana: ks}
		solRelayers, err := factory.NewSolana(loopKs, config)
		if err != nil {
			return fmt.Errorf("failed to setup Solana relayer: %w", err)
		}

		for id, relayer := range solRelayers {
			op.srvs = append(op.srvs, relayer)
			op.loopRelayers[id] = relayer
		}

		return nil
	}
}

// InitStarknet is a option for instantiating Starknet relayers
func InitStarknet(factory RelayerFactory, ks keystore.StarkNet, chainCfgs RawConfigs) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) (err error) {
		loopKs := &keystore.StarknetLooppSigner{StarkNet: ks}
		starkRelayers, err := factory.NewStarkNet(loopKs, chainCfgs)
		if err != nil {
			return fmt.Errorf("failed to setup StarkNet relayer: %w", err)
		}

		for id, relayer := range starkRelayers {
			op.srvs = append(op.srvs, relayer)
			op.loopRelayers[id] = relayer
		}

		return nil
	}
}

// InitAptos is a option for instantiating Aptos relayers
func InitAptos(factory RelayerFactory, ks keystore.Aptos, chainCfgs RawConfigs) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) (err error) {
		loopKs := &keystore.AptosLooppSigner{Aptos: ks}
		relayers, err := factory.NewAptos(loopKs, chainCfgs)
		if err != nil {
			return fmt.Errorf("failed to setup aptos relayer: %w", err)
		}

		for id, relayer := range relayers {
			op.srvs = append(op.srvs, relayer)
			op.loopRelayers[id] = relayer
		}

		return nil
	}
}

// InitTron is a option for instantiating Tron relayers
func InitTron(factory RelayerFactory, ks keystore.Tron, chainCfgs RawConfigs) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) error {
		loopKs := &keystore.TronLOOPSigner{Tron: ks}
		tronRelayers, err := factory.NewTron(loopKs, chainCfgs)
		if err != nil {
			return fmt.Errorf("failed to setup Tron relayer: %w", err)
		}

		for id, relayer := range tronRelayers {
			op.srvs = append(op.srvs, relayer)
			op.loopRelayers[id] = relayer
		}

		return nil
	}
}

// Get a [loop.Relayer] by id
func (rs *CoreRelayerChainInteroperators) Get(id types.RelayID) (loop.Relayer, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	lr, exist := rs.loopRelayers[id]
	if !exist {
		// lazily create dummy relayers
		if id.Network == "dummy" {
			var err error
			lr, err = rs.dummyFactory.NewDummy(DummyFactoryConfig{id.ChainID})
			if err != nil {
				return nil, err
			}
			rs.loopRelayers[id] = lr
			return lr, nil
		}
		return nil, fmt.Errorf("%w: %s", ErrNoSuchRelayer, id)
	}
	return lr, nil
}

func (rs *CoreRelayerChainInteroperators) GetIDToRelayerMap() map[types.RelayID]loop.Relayer {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	result := make(map[types.RelayID]loop.Relayer)
	for id, relayer := range rs.loopRelayers {
		result[id] = relayer
	}

	return result
}

// LegacyEVMChains returns a container with all the evm chains
// TODO BCF-2511
func (rs *CoreRelayerChainInteroperators) LegacyEVMChains() legacyevm.LegacyChainContainer {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.legacyChains.EVMChains
}

// ChainStatus gets [types.ChainStatus]
func (rs *CoreRelayerChainInteroperators) ChainStatus(ctx context.Context, id types.RelayID) (types.ChainStatus, error) {
	lr, err := rs.Get(id)
	if err != nil {
		return types.ChainStatus{}, fmt.Errorf("%w: error getting chain status: %w", chains.ErrNotFound, err)
	}

	return lr.GetChainStatus(ctx)
}

func (rs *CoreRelayerChainInteroperators) ChainStatuses(ctx context.Context, offset, limit int) ([]NetworkChainStatus, int, error) {
	var (
		stats    []NetworkChainStatus
		totalErr error
	)
	rs.mu.Lock()
	defer rs.mu.Unlock()

	relayerIds := make([]types.RelayID, 0)
	for rid := range rs.loopRelayers {
		relayerIds = append(relayerIds, rid)
	}
	sort.Slice(relayerIds, func(i, j int) bool {
		return relayerIds[i].String() < relayerIds[j].String()
	})
	for _, rid := range relayerIds {
		lr := rs.loopRelayers[rid]
		stat, err := lr.GetChainStatus(ctx)
		if err != nil {
			totalErr = errors.Join(totalErr, err)
			continue
		}
		stats = append(stats, NetworkChainStatus{ChainStatus: stat, Network: rid.Network})
	}

	if totalErr != nil {
		return nil, 0, totalErr
	}
	cnt := len(stats)
	if len(stats) > limit+offset && limit > 0 {
		return stats[offset : offset+limit], cnt, nil
	}
	return stats[offset:], cnt, nil
}

func (rs *CoreRelayerChainInteroperators) Node(ctx context.Context, name string) (types.NodeStatus, error) {
	// This implementation is round-about
	// TODO BFC-2511, may be better in the loop.Relayer interface itself
	stats, _, err := rs.NodeStatuses(ctx, 0, -1)
	if err != nil {
		return types.NodeStatus{}, err
	}
	for _, stat := range stats {
		if stat.Name == name {
			return stat, nil
		}
	}
	return types.NodeStatus{}, fmt.Errorf("node %s: %w", name, chains.ErrNotFound)
}

// ids must be a string representation of relay.Identifier
// ids are a filter; if none are specified, all are returned.
func (rs *CoreRelayerChainInteroperators) NodeStatuses(ctx context.Context, offset, limit int, relayerIDs ...types.RelayID) (nodes []types.NodeStatus, count int, err error) {
	var (
		totalErr error
		result   []types.NodeStatus
	)
	if len(relayerIDs) == 0 {
		for _, lr := range rs.loopRelayers {
			stats, _, total, err := lr.ListNodeStatuses(ctx, int32(limit), "")
			if err != nil {
				totalErr = errors.Join(totalErr, err)
				continue
			}
			result = append(result, stats...)
			count += total
		}
	} else {
		for _, rid := range relayerIDs {
			lr, exist := rs.loopRelayers[rid]
			if !exist {
				totalErr = errors.Join(totalErr, fmt.Errorf("relayer %s does not exist", rid.Name()))
				continue
			}
			nodeStats, _, total, err := lr.ListNodeStatuses(ctx, int32(limit), "")

			if err != nil {
				totalErr = errors.Join(totalErr, err)
				continue
			}
			result = append(result, nodeStats...)
			count += total
		}
	}
	if totalErr != nil {
		return nil, 0, totalErr
	}
	if len(result) < offset {
		return nil, 0, nil
	}
	result = result[offset:]
	if len(result) > limit && limit > 0 {
		return result[:limit], count, nil
	}
	return result, count, nil
}

type FilterFn func(id types.RelayID) bool

var AllRelayers = func(id types.RelayID) bool {
	return true
}

// Returns true if the given network matches id.Network
func FilterRelayersByType(network string) func(id types.RelayID) bool {
	return func(id types.RelayID) bool {
		return id.Network == network
	}
}

// List returns all the [RelayerChainInteroperators] that match the [FilterFn].
// A typical usage pattern to use [List] with [FilterRelayersByType] to obtain a set of [RelayerChainInteroperators]
// for a given chain
func (rs *CoreRelayerChainInteroperators) List(filter FilterFn) RelayerChainInteroperators {
	matches := make(map[types.RelayID]loop.Relayer)
	rs.mu.Lock()
	for id, relayer := range rs.loopRelayers {
		if filter(id) {
			matches[id] = relayer
		}
	}
	rs.mu.Unlock()
	return &CoreRelayerChainInteroperators{
		loopRelayers: matches,
	}
}

// Returns a slice of [loop.Relayer]. A typically usage pattern to is
// use [List(criteria)].Slice() for range based operations
func (rs *CoreRelayerChainInteroperators) Slice() []loop.Relayer {
	var result []loop.Relayer
	for _, r := range rs.loopRelayers {
		result = append(result, r)
	}
	return result
}
func (rs *CoreRelayerChainInteroperators) Services() (s []services.ServiceCtx) {
	return rs.srvs
}

// legacyChains encapsulates the chain-specific dependencies. Will be
// deprecated when chain-specific logic is removed from products.
type legacyChains struct {
	EVMChains legacyevm.LegacyChainContainer
}
