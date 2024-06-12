package chainlink

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	relay "github.com/smartcontractkit/chainlink-common/pkg/loop/adapters/relay"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services"
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
	LegacyCosmosChains() LegacyCosmosContainer
}

type ChainStatuser interface {
	ChainStatus(ctx context.Context, id types.RelayID) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error)
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

// CoreRelayerChainInteroperators implements [RelayerChainInteroperators]
// as needed for the core [chainlink.Application]
type CoreRelayerChainInteroperators struct {
	mu           sync.Mutex
	loopRelayers map[types.RelayID]loop.Relayer
	legacyChains legacyChains

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

// InitEVM is a option for instantiating evm relayers
func InitEVM(ctx context.Context, factory RelayerFactory, config EVMFactoryConfig) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) (err error) {
		adapters, err2 := factory.NewEVM(ctx, config)
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
		op.legacyChains.EVMChains = legacyevm.NewLegacyChains(legacyMap, config.AppConfig.EVMConfigs())
		return nil
	}
}

// InitCosmos is a option for instantiating Cosmos relayers
func InitCosmos(ctx context.Context, factory RelayerFactory, config CosmosFactoryConfig) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) (err error) {
		adapters, err2 := factory.NewCosmos(config)
		if err2 != nil {
			return fmt.Errorf("failed to setup Cosmos relayer: %w", err2)
		}
		legacyMap := make(map[string]cosmos.Chain)

		for id, a := range adapters {
			op.srvs = append(op.srvs, a)
			op.loopRelayers[id] = a
			legacyMap[id.ChainID] = a.Chain()
		}
		op.legacyChains.CosmosChains = NewLegacyCosmos(legacyMap)

		return nil
	}
}

// InitSolana is a option for instantiating Solana relayers
func InitSolana(ctx context.Context, factory RelayerFactory, config SolanaFactoryConfig) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) error {
		solRelayers, err := factory.NewSolana(config.Keystore, config.TOMLConfigs)
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
func InitStarknet(ctx context.Context, factory RelayerFactory, config StarkNetFactoryConfig) CoreRelayerChainInitFunc {
	return func(op *CoreRelayerChainInteroperators) (err error) {
		starkRelayers, err := factory.NewStarkNet(config.Keystore, config.TOMLConfigs)
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

// Get a [loop.Relayer] by id
func (rs *CoreRelayerChainInteroperators) Get(id types.RelayID) (loop.Relayer, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	lr, exist := rs.loopRelayers[id]
	if !exist {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchRelayer, id)
	}
	return lr, nil
}

func (rs *CoreRelayerChainInteroperators) GetIDToRelayerMap() (map[types.RelayID]loop.Relayer, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	result := make(map[types.RelayID]loop.Relayer)
	for id, relayer := range rs.loopRelayers {
		result[id] = relayer
	}

	return result, nil
}

// LegacyEVMChains returns a container with all the evm chains
// TODO BCF-2511
func (rs *CoreRelayerChainInteroperators) LegacyEVMChains() legacyevm.LegacyChainContainer {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.legacyChains.EVMChains
}

// LegacyCosmosChains returns a container with all the cosmos chains
// TODO BCF-2511
func (rs *CoreRelayerChainInteroperators) LegacyCosmosChains() LegacyCosmosContainer {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.legacyChains.CosmosChains
}

// ChainStatus gets [types.ChainStatus]
func (rs *CoreRelayerChainInteroperators) ChainStatus(ctx context.Context, id types.RelayID) (types.ChainStatus, error) {
	lr, err := rs.Get(id)
	if err != nil {
		return types.ChainStatus{}, fmt.Errorf("%w: error getting chain status: %w", chains.ErrNotFound, err)
	}

	return lr.GetChainStatus(ctx)
}

func (rs *CoreRelayerChainInteroperators) ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error) {
	var (
		stats    []types.ChainStatus
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
		stats = append(stats, stat)
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
	if len(result) > limit && limit > 0 {
		return result[offset : offset+limit], count, nil
	}
	return result[offset:], count, nil
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
// A typical usage pattern to use [List] with [FilterByType] to obtain a set of [RelayerChainInteroperators]
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
	EVMChains    legacyevm.LegacyChainContainer
	CosmosChains LegacyCosmosContainer
}

// LegacyCosmosContainer is container interface for Cosmos chains
type LegacyCosmosContainer interface {
	Get(id string) (adapters.Chain, error)
	Len() int
	List(ids ...string) ([]adapters.Chain, error)
	Slice() []adapters.Chain
}

type LegacyCosmos = chains.ChainsKV[adapters.Chain]

var _ LegacyCosmosContainer = &LegacyCosmos{}

func NewLegacyCosmos(m map[string]adapters.Chain) *LegacyCosmos {
	return chains.NewChainsKV[adapters.Chain](m)
}

type CosmosLoopRelayerChainer interface {
	loop.Relayer
	Chain() adapters.Chain
}

type CosmosLoopRelayerChain struct {
	loop.Relayer
	chain adapters.Chain
}

func NewCosmosLoopRelayerChain(r *cosmos.Relayer, s adapters.Chain) *CosmosLoopRelayerChain {
	ra := relay.NewServerAdapter(r, s)
	return &CosmosLoopRelayerChain{
		Relayer: ra,
		chain:   s,
	}
}
func (r *CosmosLoopRelayerChain) Chain() adapters.Chain {
	return r.chain
}

var _ CosmosLoopRelayerChainer = &CosmosLoopRelayerChain{}
