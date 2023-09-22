package relay

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type Network = string
type ChainID = string

var (
	EVM             Network = "evm"
	Cosmos          Network = "cosmos"
	Solana          Network = "solana"
	StarkNet        Network = "starknet"
	SupportedRelays         = map[Network]struct{}{
		EVM:      {},
		Cosmos:   {},
		Solana:   {},
		StarkNet: {},
	}
)

// ID uniquely identifies a relayer by network and chain id
type ID struct {
	Network Network
	ChainID ChainID
}

func (i *ID) Name() string {
	return fmt.Sprintf("%s.%s", i.Network, i.ChainID)
}

func (i *ID) String() string {
	return i.Name()
}
func NewID(n Network, c ChainID) ID {
	return ID{Network: n, ChainID: c}
}

var idRegex = regexp.MustCompile(
	fmt.Sprintf("^((%s)|(%s)|(%s)|(%s))\\.", EVM, Cosmos, Solana, StarkNet),
)

func (i *ID) UnmarshalString(s string) error {
	idxs := idRegex.FindStringIndex(s)
	if idxs == nil {
		return fmt.Errorf("error unmarshaling Identifier. %q does not match expected pattern", s)
	}
	// ignore the `.` in the match by dropping last rune
	network := s[idxs[0] : idxs[1]-1]
	chainID := s[idxs[1]:]
	newID := &ID{ChainID: ChainID(chainID)}
	for n := range SupportedRelays {
		if Network(network) == n {
			newID.Network = n
			break
		}
	}
	if newID.Network == "" {
		return fmt.Errorf("error unmarshaling identifier: did not find network in supported list %q", network)
	}
	i.ChainID = newID.ChainID
	i.Network = newID.Network
	return nil
}

// RelayerExt is a subset of [loop.Relayer] for adapting [types.Relayer], typically with a Chain. See [relayerAdapter].
type RelayerExt interface {
	types.ChainService
	ID() string
}

var _ loop.Relayer = (*relayerAdapter)(nil)

// relayerAdapter adapts a [types.Relayer] and [RelayerExt] to implement [loop.Relayer].
type relayerAdapter struct {
	types.Relayer
	RelayerExt
}

// NewRelayerAdapter returns a [loop.Relayer] adapted from a [types.Relayer] and [RelayerExt].
// Unlike NewRelayerServerAdapter which is used to adapt non-LOOPP relayers, this is used to adapt
// LOOPP-based relayer which are then server over GRPC (by the relayerServer).
func NewRelayerAdapter(r types.Relayer, e RelayerExt) loop.Relayer {
	return &relayerAdapter{Relayer: r, RelayerExt: e}
}

func (r *relayerAdapter) NewConfigProvider(ctx context.Context, rargs types.RelayArgs) (types.ConfigProvider, error) {
	return r.Relayer.NewConfigProvider(rargs)
}

func (r *relayerAdapter) NewMedianProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MedianProvider, error) {
	return r.Relayer.NewMedianProvider(rargs, pargs)
}

func (r *relayerAdapter) NewMercuryProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	return r.Relayer.NewMercuryProvider(rargs, pargs)
}

func (r *relayerAdapter) NewFunctionsProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.FunctionsProvider, error) {
	return r.Relayer.NewFunctionsProvider(rargs, pargs)
}

func (r *relayerAdapter) NewPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.PluginProvider, error) {
	return nil, fmt.Errorf("unexpected call to NewPluginProvider: did you forget to wrap relayerAdapter in a relayerServerAdapter?")
}

func (r *relayerAdapter) Start(ctx context.Context) error {
	var ms services.MultiStart
	return ms.Start(ctx, r.RelayerExt, r.Relayer)
}

func (r *relayerAdapter) Close() error {
	return services.CloseAll(r.Relayer, r.RelayerExt)
}

func (r *relayerAdapter) Name() string {
	return fmt.Sprintf("%s-%s", r.Relayer.Name(), r.RelayerExt.Name())
}

func (r *relayerAdapter) Ready() (err error) {
	return errors.Join(r.Relayer.Ready(), r.RelayerExt.Ready())
}

func (r *relayerAdapter) HealthReport() map[string]error {
	hr := make(map[string]error)
	maps.Copy(hr, r.Relayer.HealthReport())
	maps.Copy(hr, r.RelayerExt.HealthReport())
	return hr
}

func (r *relayerAdapter) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, total int, err error) {
	if len(chainIDs) > 1 {
		return nil, 0, fmt.Errorf("internal error: node statuses expects at most one chain id got %v", chainIDs)
	}
	if len(chainIDs) == 1 && chainIDs[0] != r.ID() {
		return nil, 0, fmt.Errorf("node statuses unexpected chain id got %s want %s", chainIDs[0], r.ID())
	}

	nodes, _, total, err = r.ListNodeStatuses(ctx, int32(limit), "")
	if err != nil {
		return nil, 0, err
	}
	if len(nodes) < offset {
		return []types.NodeStatus{}, 0, fmt.Errorf("out of range")
	}
	if limit <= 0 {
		limit = len(nodes)
	} else if len(nodes) < limit {
		limit = len(nodes)
	}
	return nodes[offset:limit], total, nil
}

type relayerServerAdapter struct {
	*relayerAdapter
}

func (r *relayerServerAdapter) NewPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.PluginProvider, error) {
	switch types.OCR2PluginType(rargs.ProviderType) {
	case types.Median:
		return r.NewMedianProvider(ctx, rargs, pargs)
	case types.Functions:
		return r.NewFunctionsProvider(ctx, rargs, pargs)
	case types.Mercury:
		return r.NewMercuryProvider(ctx, rargs, pargs)
	case types.DKG, types.OCR2VRF, types.OCR2Keeper, types.GenericPlugin:
		return r.relayerAdapter.NewPluginProvider(ctx, rargs, pargs)
	}

	return nil, fmt.Errorf("provider type not supported: %s", rargs.ProviderType)
}

// NewRelayerServerAdapter returns a [loop.Relayer] adapted from a [types.Relayer] and [RelayerExt].
// Unlike NewRelayerAdapter, this behaves like the loop `RelayerServer` and dispatches calls
// to `NewPluginProvider` according to the passed in `RelayArgs.ProviderType`.
// This should only be used to adapt relayers not running via GRPC in a LOOPP.
//
// nolint:staticcheck // SA1019
func NewRelayerServerAdapter(r types.Relayer, e RelayerExt) loop.Relayer {
	ra := &relayerAdapter{Relayer: r, RelayerExt: e}
	return &relayerServerAdapter{relayerAdapter: ra}
}
