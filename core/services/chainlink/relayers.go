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
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	evmrelayer "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

// Chains holds a ChainSet for each type of chain.
type chains struct {
	EVMChains *evm.Chains //evmloop.LoopRelayAdapter //evm.ChainSet                      // map[relay.Identifier]evm.ChainSet

	CosmosChains *cosmos.Chains
}

type Relayers struct {
	mu       sync.Mutex
	relayers map[relay.Identifier]loop.Relayer
	chains   chains
	// have to treat evm as special because of the Default func. maybe there is a cleaner way...
	//defaultEvmID relay.Identifier
}

var ErrNoSuchRelayer = errors.New("relayer does not exist")

// TODO generics to simplify
/*
func (rs *Relayers) EVM(ids ...relay.Identifier) (map[relay.Identifier]evmloop.LoopRelayAdapter, error) {
	var (
		result map[relay.Identifier]evmloop.LoopRelayAdapter
		err    error
	)
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if len(ids) == 0 {
		return rs.chains.EVM, nil
	}
	for _, id := range ids {
		r, ok := rs.chains.EVM[id]
		if !ok {
			err = errors.Join(err, fmt.Errorf("no such id %s", id))
			continue
		}
		result[id] = r
	}
	return result, err
}

func (rs *Relayers) Cosmos(ids ...relay.Identifier) (map[relay.Identifier]cosmos.LoopRelayAdapter, error) {
	var (
		result map[relay.Identifier]cosmos.LoopRelayAdapter
		err    error
	)
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if len(ids) == 0 {
		return rs.chains.Cosmos, nil
	}
	for _, id := range ids {
		r, ok := rs.chains.Cosmos[id]
		if !ok {
			err = errors.Join(err, fmt.Errorf("no such id %s", id))
			continue
		}
		result[id] = r
	}
	return result, err
}

func (rs *Relayers) PutEVM(id relay.Identifier, l evmloop.LoopRelayAdapter) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.chains.EVM[id] = l
}

func (rs *Relayers) PutCosmos(id relay.Identifier, l cosmos.LoopRelayAdapter) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.chains.Cosmos[id] = l
}
*/

func (rs *Relayers) Get(id relay.Identifier) (loop.Relayer, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	r, exist := rs.relayers[id]
	if !exist {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchRelayer, id)
	}
	return r, nil
}

func (rs *Relayers) putOne(id relay.Identifier, r loop.Relayer) error {

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

func (rs *Relayers) Put(id relay.Identifier, r loop.Relayer) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.putOne(id, r)
}

// TODO maybe make Relayer[U,V] where u,v are the chain specific types and then make this generic
// Relayer[U evm.LoopAdapter, Vcosmos.LoopAdapter]
// (rs Relayer[U,V] PutBatch(map[](loop.relayer|U|V))
func (rs *Relayers) PutBatch(b map[relay.Identifier]loop.Relayer) (err error) {
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

func (rs *Relayers) LegacyEVMChains() *evm.Chains {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	// for backward compatibility we return an empty, non-nil value here
	// all other chains/relayers can be nil...
	if rs.chains.EVMChains == nil {
		rs.chains.EVMChains = evm.NewLegacyChains()
	}
	return rs.chains.EVMChains
}

func (rs *Relayers) LegacyCosmosChains() *cosmos.Chains {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.chains.CosmosChains
}

func (rs *Relayers) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
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

func (rs *Relayers) ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error) {
	// chain statuses are not dynamic; the call would be better named as ChainConfig or such.
	// lazily create a cache and use that case for the offset and limit to ensure deterministic results

	return nil, 0, nil
}

func (rs *Relayers) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	return nil, 0, nil
}

// backwards compatibility with Default func for evm chainset
/*
func (rs *Relayers) DefaultEVM() (loop.Relayer, error) {
	return rs.Get(rs.defaultEvmID)
}

func (rs *Relayers) SetDefaultEVM(id relay.Identifier) error {
	_, err := rs.Get(id)
	if err != nil {
		return fmt.Errorf("failed to set default evm relayer. has it been put?: %w", err)
	}
	rs.defaultEvmID = id
	return nil
}
*/
type FilterFn func(id relay.Identifier) bool

var AllRelayers = func(id relay.Identifier) bool {
	return true
}

func FilterByType(network relay.Network) func(id relay.Identifier) bool {
	return func(id relay.Identifier) bool {
		return id.Network == network
	}
}

func (rs *Relayers) List(filter FilterFn) *Relayers {

	var matches map[relay.Identifier]loop.Relayer
	rs.mu.Lock()
	for id, relayer := range rs.relayers {
		if filter(id) {
			matches[id] = relayer
		}
	}
	rs.mu.Unlock()
	return &Relayers{
		relayers: matches,
		// note filtering may cause the default evm to no longer be in the returned set.
		//defaultEvmID: rs.defaultEvmID,
	}
}

func (rs *Relayers) Slice() []loop.Relayer {
	var result []loop.Relayer
	for _, r := range rs.relayers {
		result = append(result, r)
	}
	return result
}
func (rs *Relayers) services() (s []services.ServiceCtx) {
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
