package relay

import (
	"context"
	"errors"
	"fmt"
	"math/big"

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

// RelayerExt is a subset of [loop.Relayer] for adapting [types.Relayer], typically with a ChainSet. See [RelayerAdapter].
type RelayerExt interface {
	services.ServiceCtx

	ChainStatus(ctx context.Context, id string) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error)

	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error)

	SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error
}

var _ loop.Relayer = (*RelayerAdapter)(nil)

// RelayerAdapter adapts a [types.Relayer] and [RelayerExt] to imlement [loop.Relayer].
type RelayerAdapter struct {
	types.Relayer
	RelayerExt
}

func (r *RelayerAdapter) NewConfigProvider(ctx context.Context, rargs types.RelayArgs) (types.ConfigProvider, error) {
	return r.Relayer.NewConfigProvider(rargs)
}

func (r *RelayerAdapter) NewMedianProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MedianProvider, error) {
	return r.Relayer.NewMedianProvider(rargs, pargs)
}

func (r *RelayerAdapter) NewMercuryProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	return r.Relayer.NewMercuryProvider(rargs, pargs)
}

func (r *RelayerAdapter) Start(ctx context.Context) error {
	var ms services.MultiStart
	return ms.Start(ctx, r.RelayerExt, r.Relayer)
}

func (r *RelayerAdapter) Close() error {
	return services.MultiClose{r.Relayer, r.RelayerExt}.Close()
}

func (r *RelayerAdapter) Name() string {
	return fmt.Sprintf("%s-%s", r.Relayer.Name(), r.RelayerExt.Name())
}

func (r *RelayerAdapter) Ready() (err error) {
	return errors.Join(r.Relayer.Ready(), r.RelayerExt.Ready())
}

func (r *RelayerAdapter) HealthReport() map[string]error {
	hr := make(map[string]error)
	maps.Copy(r.Relayer.HealthReport(), hr)
	maps.Copy(r.RelayerExt.HealthReport(), hr)
	return hr
}

type RelayerService interface {
	services.ServiceCtx
	Relayer() (loop.Relayer, error)
}

type RelayerServiceAdapter struct {
	*RelayerAdapter
}

func (a *RelayerServiceAdapter) Relayer() (loop.Relayer, error) {
	return a.RelayerAdapter, nil
}

// NewLocalRelayerService returns a RelayerService adapted from a [types.Relayer] and [RelayerExt].
func NewLocalRelayerService(r types.Relayer, e RelayerExt) RelayerService {
	return &RelayerServiceAdapter{&RelayerAdapter{Relayer: r, RelayerExt: e}}
}
