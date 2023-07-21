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

// Chains holds a ChainSet for each type of chain.
type legacyChains struct {
	EVMChains *evm.Chains //evmloop.LoopRelayAdapter //evm.ChainSet                      // map[relay.Identifier]evm.ChainSet

	CosmosChains *cosmos.Chains
}

type RelayChainInterchangers interface {
	Services() []services.ServiceCtx

	List(filter FilterFn) RelayChainInterchangers

	LoopRelayerStorer
	LegacyChainGetter
	OperationalStatuser
}

type LoopRelayerStorer interface {
	ocr2.RelayGetter
	Put(id relay.Identifier, r loop.Relayer) error
	PutBatch(b map[relay.Identifier]loop.Relayer) error
	Slice() []loop.Relayer
}
type LegacyChainGetter interface {
	LegacyEVMChains() *evm.Chains
	LegacyCosmosChains() *cosmos.Chains
}

type OperationalStatuser interface {
	chains.ChainStatuser
	chains.NodesStatuser
}

var _ RelayChainInterchangers = &RelayChainInteroperators{}

type RelayChainInteroperators struct {
	mu       sync.Mutex
	relayers map[relay.Identifier]loop.Relayer
	chains   legacyChains
}

func NewRelayers() *RelayChainInteroperators {
	return &RelayChainInteroperators{
		relayers: make(map[relay.Identifier]loop.Relayer),
		chains:   legacyChains{EVMChains: evm.NewLegacyChains(), CosmosChains: new(cosmos.Chains)},
	}
}

var ErrNoSuchRelayer = errors.New("relayer does not exist")

func (rs *RelayChainInteroperators) Get(id relay.Identifier) (loop.Relayer, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	r, exist := rs.relayers[id]
	if !exist {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchRelayer, id)
	}
	return r, nil
}

func (rs *RelayChainInteroperators) putOne(id relay.Identifier, r loop.Relayer) error {

	// backward compatibility. this is bit gross to type cast but it hides the details from products.
	switch id.Network {
	case relay.EVM:
		adapter, ok := r.(evmrelayer.LoopRelayAdapter)
		if !ok {
			return fmt.Errorf("unsupported evm loop relayer implementation. got %t want (evmrelayer.LoopRelayAdapter)", r)
		}

		rs.chains.EVMChains.Put(id.ChainID.String(), adapter.Chain())
		if adapter.Default() {
			dflt, _ := rs.chains.EVMChains.Default()
			if dflt != nil {
				return fmt.Errorf("multiple default evm chains. %s, %s", dflt.ID(), adapter.Chain().ID())
			}
			rs.chains.EVMChains.SetDefault(adapter.Chain())
		}
	case relay.Cosmos:
		adapter, ok := r.(cosmos.LoopRelayAdapter)
		if !ok {
			return fmt.Errorf("unsupported cosmos loop relayer implementation. got %t want (cosmos.LoopRelayAdapter)", r)
		}

		rs.chains.CosmosChains.Put(id.ChainID.String(), adapter.Chain())
	}

	rs.relayers[id] = r
	return nil
}

func (rs *RelayChainInteroperators) Put(id relay.Identifier, r loop.Relayer) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.putOne(id, r)
}

// TODO maybe make Relayer[U,V] where u,v are the chain specific types and then make this generic
// Relayer[U evm.LoopAdapter, Vcosmos.LoopAdapter]
// (rs Relayer[U,V] PutBatch(map[](loop.relayer|U|V))
func (rs *RelayChainInteroperators) PutBatch(b map[relay.Identifier]loop.Relayer) (err error) {
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

func (rs *RelayChainInteroperators) LegacyEVMChains() *evm.Chains {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	return rs.chains.EVMChains
}

func (rs *RelayChainInteroperators) LegacyCosmosChains() *cosmos.Chains {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.chains.CosmosChains
}

func (rs *RelayChainInteroperators) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
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

func (rs *RelayChainInteroperators) ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error) {
	// chain statuses are not dynamic; the call would be better named as ChainConfig or such.
	// lazily create a cache and use that case for the offset and limit to ensure deterministic results

	return nil, 0, nil
}

func (rs *RelayChainInteroperators) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	return nil, 0, nil
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

func (rs *RelayChainInteroperators) List(filter FilterFn) RelayChainInterchangers {

	var matches map[relay.Identifier]loop.Relayer
	rs.mu.Lock()
	for id, relayer := range rs.relayers {
		if filter(id) {
			matches[id] = relayer
		}
	}
	rs.mu.Unlock()
	return &RelayChainInteroperators{
		relayers: matches,
	}
}

func (rs *RelayChainInteroperators) Slice() []loop.Relayer {
	var result []loop.Relayer
	for _, r := range rs.relayers {
		result = append(result, r)
	}
	return result
}
func (rs *RelayChainInteroperators) Services() (s []services.ServiceCtx) {
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
