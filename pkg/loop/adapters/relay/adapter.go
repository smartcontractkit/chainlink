package relay

import (
	"context"

	looptypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/types"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// RelayerExt is a subset of [loop.Relayer] for adapting [types.Relayer], typically with a Chain. See [RelayerAdapter].
type RelayerExt interface {
	types.ChainService
	ID() string
}

var _ looptypes.Relayer = (*RelayerAdapter)(nil)

// RelayerAdapter adapts a [types.Relayer] and [RelayerExt] to implement [Relayer].
type RelayerAdapter struct {
	types.Relayer
	RelayerExt
}

func (r *RelayerAdapter) NewContractReader(_ context.Context, contractReaderConfig []byte) (types.ContractReader, error) {
	return r.Relayer.NewContractReader(contractReaderConfig)
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

func (r *RelayerAdapter) NewAutomationProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.AutomationProvider, error) {
	return r.Relayer.NewAutomationProvider(rargs, pargs)
}

func (r *RelayerAdapter) NewLLOProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.LLOProvider, error) {
	return r.Relayer.NewLLOProvider(rargs, pargs)
}

func (r *RelayerAdapter) NewPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.PluginProvider, error) {
	return r.Relayer.NewPluginProvider(rargs, pargs)
}

func (r *RelayerAdapter) NewOCR3CapabilityProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.OCR3CapabilityProvider, error) {
	return r.Relayer.NewOCR3CapabilityProvider(rargs, pargs)
}

func (r *RelayerAdapter) Start(ctx context.Context) error {
	var ms services.MultiStart
	return ms.Start(ctx, r.RelayerExt, r.Relayer)
}

func (r *RelayerAdapter) Close() error {
	return services.CloseAll(r.Relayer, r.RelayerExt)
}

func (r *RelayerAdapter) Name() string {
	return r.Relayer.Name()
}

func (r *RelayerAdapter) Ready() (err error) {
	return r.Relayer.Ready()
}

func (r *RelayerAdapter) HealthReport() map[string]error {
	hr := make(map[string]error)
	services.CopyHealth(hr, r.Relayer.HealthReport())
	return hr
}
