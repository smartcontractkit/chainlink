package relay

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// ServerAdapter extends [loop.RelayerAdapter] by overriding NewPluginProvider to dispatches calls according to `RelayArgs.ProviderType`.
// This should only be used to adapt relayers not running via GRPC in a LOOPP.
type ServerAdapter struct {
	RelayerAdapter
}

// NewServerAdapter returns a new ServerAdapter.
func NewServerAdapter(r types.Relayer, e RelayerExt) *ServerAdapter { //nolint:staticcheck
	return &ServerAdapter{RelayerAdapter: RelayerAdapter{Relayer: r, RelayerExt: e}}
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
	case types.DKG, types.OCR2VRF, types.GenericPlugin:
		return r.RelayerAdapter.NewPluginProvider(ctx, rargs, pargs)
	case types.LLO, types.CCIPCommit, types.CCIPExecution:
		return nil, fmt.Errorf("provider type not supported: %s", rargs.ProviderType)
	}
	return nil, fmt.Errorf("provider type not recognized: %s", rargs.ProviderType)
}
