package relay

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

const (
	NetworkEVM      = "evm"
	NetworkCosmos   = "cosmos"
	NetworkSolana   = "solana"
	NetworkStarkNet = "starknet"
	NetworkAptos    = "aptos"
	NetworkTron     = "tron"

	NetworkDummy = "dummy"
)

var SupportedNetworks = map[string]struct{}{
	NetworkEVM:      {},
	NetworkCosmos:   {},
	NetworkSolana:   {},
	NetworkStarkNet: {},
	NetworkAptos:    {},
	NetworkTron:     {},

	NetworkDummy: {},
}

var _ loop.Relayer = (*ServerAdapter)(nil)

// ServerAdapter extends [loop.RelayerAdapter] by overriding NewPluginProvider to dispatches calls according to `RelayArgs.ProviderType`.
// This should only be used to adapt relayers not running via GRPC in a LOOPP.
type ServerAdapter struct {
	types.Relayer
}

// NewServerAdapter returns a new ServerAdapter.
func NewServerAdapter(r types.Relayer) *ServerAdapter {
	return &ServerAdapter{Relayer: r}
}

func (r *ServerAdapter) NewPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.PluginProvider, error) {
	switch types.OCR2PluginType(rargs.ProviderType) {
	case types.Median:
		return r.NewMedianProvider(ctx, rargs, pargs)
	case types.Functions:
		return r.NewFunctionsProvider(ctx, rargs, pargs)
	case types.Mercury:
		return r.NewMercuryProvider(ctx, rargs, pargs)
	case types.OCR2Keeper:
		return r.NewAutomationProvider(ctx, rargs, pargs)
	case types.OCR3Capability:
		return r.NewOCR3CapabilityProvider(ctx, rargs, pargs)
	case types.CCIPCommit:
		return r.NewCCIPCommitProvider(ctx, rargs, pargs)
	case types.CCIPExecution:
		return r.NewCCIPExecProvider(ctx, rargs, pargs)
	case types.DKG, types.OCR2VRF, types.GenericPlugin:
		return r.Relayer.NewPluginProvider(ctx, rargs, pargs)
	case types.LLO:
		return nil, fmt.Errorf("provider type not supported: %s", rargs.ProviderType)
	}
	return nil, fmt.Errorf("provider type not recognized: %s", rargs.ProviderType)
}
