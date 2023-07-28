package chainlink

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/hashicorp/go-multierror"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	evmrelayer "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

// RelayerChainInteroperators
// encapulates relayers and chains and is the primary entry point for
// the node to store relayers, get legacy chains associated to a relayer
// and get status about the chains and nodes
//
//go:generate mockery --quiet --name RelayerChainInteroperators --output ./mocks/ --case=underscore
type RelayerChainInteroperators interface {
	Services() []services.ServiceCtx

	List(filter FilterFn) RelayerChainInteroperators

	LoopRelayerStorer
	LegacyChainer
	ChainsNodesStatuser
}

// LoopRelayerStorer is key-value like interface for storing and
// retrieving [loop.Relayer] by [relay.Identifier]
type LoopRelayerStorer interface {
	ocr2.RelayGetter
	Put(id relay.Identifier, r loop.Relayer) error
	PutBatch(b map[relay.Identifier]loop.Relayer) error
	Slice() []loop.Relayer
}

// LegacyChainer is an interface for getting legacy chains
// This will be deprecated/removed when products depend only
// on the relayer interface.
type LegacyChainer interface {
	LegacyEVMChains() evm.LegacyChainContainer //evm.LegacyChainContainer
	LegacyCosmosChains() *cosmos.Chains
}

// ChainsNodesStatuser report statuses about chains and nodes
type ChainsNodesStatuser interface {
	chains.ChainStatuser
	chains.NodesStatuser
}

var _ RelayerChainInteroperators = &CoreRelayerChainInteroperators{}

// CoreRelayerChainInteroperators implements [RelayerChainInteroperators]
// as needed for the core [chainlink.Application]
type CoreRelayerChainInteroperators struct {
	mu           sync.Mutex
	relayers     map[relay.Identifier]loop.Relayer
	legacyChains legacyChains
}

func NewCoreRelayerChainInteroperators() *CoreRelayerChainInteroperators {
	return &CoreRelayerChainInteroperators{
		relayers:     make(map[relay.Identifier]loop.Relayer),
		legacyChains: legacyChains{EVMChains: evm.NewLegacyChains(), CosmosChains: new(cosmos.Chains)},
	}
}

var ErrNoSuchRelayer = errors.New("relayer does not exist")

func (rs *CoreRelayerChainInteroperators) Get(id relay.Identifier) (loop.Relayer, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	r, exist := rs.relayers[id]
	if !exist {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchRelayer, id)
	}
	return r, nil
}

func (rs *CoreRelayerChainInteroperators) putOne(id relay.Identifier, r loop.Relayer) error {

	// backward compatibility. this is bit gross to type cast but it hides the details from products.
	switch id.Network {
	case relay.EVM:
		adapter, ok := r.(evmrelayer.LoopRelayAdapter)
		if !ok {
			return fmt.Errorf("unsupported evm loop relayer implementation. got %t want (evmrelayer.LoopRelayAdapter)", r)
		}

		rs.legacyChains.EVMChains.Put(id.ChainID.String(), adapter.Chain())
		if adapter.Default() {
			dflt, _ := rs.legacyChains.EVMChains.Default()
			if dflt != nil {
				return fmt.Errorf("multiple default evm chains. %s, %s", dflt.ID(), adapter.Chain().ID())
			}
			rs.legacyChains.EVMChains.SetDefault(adapter.Chain())
		}
	case relay.Cosmos:
		adapter, ok := r.(cosmos.LoopRelayAdapter)
		if !ok {
			return fmt.Errorf("unsupported cosmos loop relayer implementation. got %t want (cosmos.LoopRelayAdapter)", r)
		}

		rs.legacyChains.CosmosChains.Put(id.ChainID.String(), adapter.Chain())
	}

	rs.relayers[id] = r
	return nil
}

func (rs *CoreRelayerChainInteroperators) Put(id relay.Identifier, r loop.Relayer) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.putOne(id, r)
}

// TODO maybe make Relayer[U,V] where u,v are the chain specific types and then make this generic
// Relayer[U evm.LoopAdapter, Vcosmos.LoopAdapter]
// (rs Relayer[U,V] PutBatch(map[](loop.relayer|U|V))
func (rs *CoreRelayerChainInteroperators) PutBatch(b map[relay.Identifier]loop.Relayer) (err error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	for id, r := range b {
		err2 := rs.putOne(id, r)
		if err2 != nil {
			multierror.Append(err, err2)
		}
	}
	return err
}

func (rs *CoreRelayerChainInteroperators) LegacyEVMChains() evm.LegacyChainContainer {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	return rs.legacyChains.EVMChains
}

func (rs *CoreRelayerChainInteroperators) LegacyCosmosChains() *cosmos.Chains {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.legacyChains.CosmosChains
}

func (rs *CoreRelayerChainInteroperators) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
	relayID := new(relay.Identifier)
	err := relayID.UnmarshalString(id)
	if err != nil {
		return types.ChainStatus{}, fmt.Errorf("error getting chainstatus: %w", err)
	}
	relayer, err := rs.Get(*relayID)
	if err != nil {
		return types.ChainStatus{}, fmt.Errorf("error getting chainstatus: %w", err)
	}
	// this call is weird because the [loop.Relayer] interface still requires id
	// but in this context the `relayer` should only have only id
	return relayer.ChainStatus(ctx, id)
}

func (rs *CoreRelayerChainInteroperators) ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error) {
	// chain statuses are not dynamic; the call would be better named as ChainConfig or such.
	// lazily create a cache and use that case for the offset and limit to ensure deterministic results

	return nil, 0, nil
}

// ids must be a string representation of relay.Identifier
func (rs *CoreRelayerChainInteroperators) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	var (
		totalErr error
		result   []types.NodeStatus
	)
	for _, idStr := range chainIDs {
		rid := new(relay.Identifier)
		err := rid.UnmarshalString(idStr)
		if err != nil {
			totalErr = errors.Join(totalErr, err)
			continue
		}
		relayer := rs.relayers[*rid]
		nodeStats, _, err := relayer.NodeStatuses(ctx, offset, limit, rid.ChainID.String())

		if err != nil {
			totalErr = errors.Join(totalErr, err)
			continue
		}
		result = append(result, nodeStats...)

	}
	return result, len(result), totalErr
}

type FilterFn func(id relay.Identifier) bool

var AllRelayers = func(id relay.Identifier) bool {
	return true
}

func FilterByType(network relay.Network) func(id relay.Identifier) bool {
	return func(id relay.Identifier) bool {
		return id.Network == network
	}
}

func (rs *CoreRelayerChainInteroperators) List(filter FilterFn) RelayerChainInteroperators {

	matches := make(map[relay.Identifier]loop.Relayer)
	rs.mu.Lock()
	for id, relayer := range rs.relayers {
		if filter(id) {
			matches[id] = relayer
		}
	}
	rs.mu.Unlock()
	return &CoreRelayerChainInteroperators{
		relayers: matches,
	}
}

func (rs *CoreRelayerChainInteroperators) Slice() []loop.Relayer {
	var result []loop.Relayer
	for _, r := range rs.relayers {
		result = append(result, r)
	}
	return result
}
func (rs *CoreRelayerChainInteroperators) Services() (s []services.ServiceCtx) {
	// TODO. ensure that the services are not duplicated between the chain and relayers...
	s = append(s, sortByChainID(rs.relayers)...)
	return
}

func sortByChainID[V services.ServiceCtx](m map[relay.Identifier]V) []services.ServiceCtx {
	sorted := make([]services.ServiceCtx, len(m))
	ids := make([]relay.Identifier, 0)
	for id := range m {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].ChainID.String() < ids[j].ChainID.String()
	})
	for i := 0; i < len(m); i += 1 {
		sorted[i] = m[ids[i]]
	}
	return sorted
}

// legacyChains encapsulates the chain-specific dependencies. Will be
// deprecated when chain-specific logic is removed from products.
type legacyChains struct {
	EVMChains    evm.LegacyChainContainer
	CosmosChains *cosmos.Chains
}
