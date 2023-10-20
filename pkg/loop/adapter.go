package loop

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-relay/pkg/services"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

// RelayerExt is a subset of [loop.Relayer] for adapting [types.Relayer], typically with a Chain. See [RelayerAdapter].
type RelayerExt interface {
	types.ChainService
	ID() string
}

var _ Relayer = (*RelayerAdapter)(nil)

// RelayerAdapter adapts a [types.Relayer] and [RelayerExt] to implement [Relayer].
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

func (r *RelayerAdapter) NewFunctionsProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.FunctionsProvider, error) {
	return r.Relayer.NewFunctionsProvider(rargs, pargs)
}

func (r *RelayerAdapter) NewPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.PluginProvider, error) {
	return nil, fmt.Errorf("unexpected call to NewPluginProvider: did you forget to wrap RelayerAdapter in a relayerServerAdapter?")
}

func (r *RelayerAdapter) Start(ctx context.Context) error {
	var ms services.MultiStart
	return ms.Start(ctx, r.RelayerExt, r.Relayer)
}

func (r *RelayerAdapter) Close() error {
	return services.CloseAll(r.Relayer, r.RelayerExt)
}

func (r *RelayerAdapter) Name() string {
	return fmt.Sprintf("%s-%s", r.Relayer.Name(), r.RelayerExt.Name())
}

func (r *RelayerAdapter) Ready() (err error) {
	return errors.Join(r.Relayer.Ready(), r.RelayerExt.Ready())
}

func (r *RelayerAdapter) HealthReport() map[string]error {
	hr := make(map[string]error)
	services.CopyHealth(hr, r.Relayer.HealthReport())
	return hr
}
