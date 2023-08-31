package relay

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type Network string

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
	return fmt.Sprintf("%s.%s", i.Network, i.ChainID.String())
}

func (i *ID) String() string {
	return i.Name()
}
func NewID(n Network, c ChainID) (ID, error) {
	id := ID{Network: n, ChainID: c}
	err := id.validate()
	if err != nil {
		return ID{}, err
	}
	return id, nil
}
func (i *ID) validate() error {
	// the only validation is to ensure that EVM chain ids are compatible with int64
	if i.Network == EVM {
		_, err := i.ChainID.Int64()
		if err != nil {
			return fmt.Errorf("RelayIdentifier invalid: EVM relayer must have integer-compatible chain ID: %w", err)
		}
	}
	return nil
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

type ChainID string

func (c ChainID) String() string {
	return string(c)
}
func (c ChainID) Int64() (int64, error) {
	i, err := strconv.Atoi(c.String())
	if err != nil {
		return int64(0), err
	}
	return int64(i), nil
}

// RelayerExt is a subset of [loop.Relayer] for adapting [types.Relayer], typically with a ChainSet. See [relayerAdapter].
type RelayerExt interface {
	types.ChainService
	// TODO remove after BFC-2441
	ID() string
	GetChainStatus(ctx context.Context) (types.ChainStatus, error)
	ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error)
	// choose different name than SendTx to avoid collison during refactor.
	Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error
}

var _ loop.Relayer = (*relayerAdapter)(nil)

// relayerAdapter adapts a [types.Relayer] and [RelayerExt] to implement [loop.Relayer].
type relayerAdapter struct {
	types.Relayer
	// TODO we can un-embedded `ext` once BFC-2441 is merged. Right now that's not possible
	// because this are conflicting definitions of SendTx
	ext RelayerExt
}

// NewRelayerAdapter returns a [loop.Relayer] adapted from a [types.Relayer] and [RelayerExt].
func NewRelayerAdapter(r types.Relayer, e RelayerExt) loop.Relayer {
	return &relayerAdapter{Relayer: r, ext: e}
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

func (r *relayerAdapter) Start(ctx context.Context) error {
	var ms services.MultiStart
	return ms.Start(ctx, r.ext, r.Relayer)
}

func (r *relayerAdapter) Close() error {
	return services.CloseAll(r.Relayer, r.ext)
}

func (r *relayerAdapter) Name() string {
	return fmt.Sprintf("%s-%s", r.Relayer.Name(), r.ext.Name())
}

func (r *relayerAdapter) Ready() (err error) {
	return errors.Join(r.Relayer.Ready(), r.ext.Ready())
}

func (r *relayerAdapter) HealthReport() map[string]error {
	hr := make(map[string]error)
	maps.Copy(r.Relayer.HealthReport(), hr)
	maps.Copy(r.ext.HealthReport(), hr)
	return hr
}

// Implement the existing [loop.Relayer] interface using the underlaying chain service
// TODO Delete this code after BFC-2441

func (r *relayerAdapter) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
	if id != r.ext.ID() {
		return types.ChainStatus{}, fmt.Errorf("unexpected chain id. got %s want %s", id, r.ID())
	}
	return r.ext.GetChainStatus(ctx)
}
func (r *relayerAdapter) ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error) {
	stat, err := r.ext.GetChainStatus(ctx)
	if err != nil {
		return nil, -1, err
	}
	return []types.ChainStatus{stat}, 1, nil
}

func (r *relayerAdapter) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, total int, err error) {
	if len(chainIDs) > 1 {
		return nil, 0, fmt.Errorf("internal error: node statuses expects at most one chain id got %v", chainIDs)
	}
	if len(chainIDs) == 1 && chainIDs[0] != r.ext.ID() {
		return nil, 0, fmt.Errorf("node statuses unexpected chain id got %s want %s", chainIDs[0], r.ID())
	}

	nodes, _, total, err = r.ext.ListNodeStatuses(ctx, int32(limit), "")
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

func (r *relayerAdapter) SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error {
	if chainID != r.ext.ID() {
		return fmt.Errorf("send tx unexpected chain id. got %s want %s", chainID, r.ext.ID())
	}
	return r.ext.Transact(ctx, from, to, amount, balanceCheck)
}

func (r *relayerAdapter) ID() string {
	return r.ext.ID()
}

func (r *relayerAdapter) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	return r.ext.GetChainStatus(ctx)
}

func (r *relayerAdapter) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	return r.ext.ListNodeStatuses(ctx, pageSize, pageToken)
}

// choose different name than SendTx to avoid collison during refactor.
func (r *relayerAdapter) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return r.ext.Transact(ctx, from, to, amount, balanceCheck)
}
