package relay

import (
	"context"
	"fmt"
	"regexp"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type Network = string
type ChainID = string

const (
	EVM      = "evm"
	Cosmos   = "cosmos"
	Solana   = "solana"
	StarkNet = "starknet"
)

var SupportedRelays = map[Network]struct{}{
	EVM:      {},
	Cosmos:   {},
	Solana:   {},
	StarkNet: {},
}

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
	newID := &ID{ChainID: chainID}
	for n := range SupportedRelays {
		if network == n {
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

// ServerAdapter extends [loop.RelayerAdapter] by overriding NewPluginProvider to dispatches calls according to `RelayArgs.ProviderType`.
// This should only be used to adapt relayers not running via GRPC in a LOOPP.
type ServerAdapter struct {
	loop.RelayerAdapter
}

// NewServerAdapter returns a new ServerAdapter.
func NewServerAdapter(r types.Relayer, e loop.RelayerExt) *ServerAdapter { //nolint:staticcheck
	return &ServerAdapter{RelayerAdapter: loop.RelayerAdapter{Relayer: r, RelayerExt: e}}
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
	case types.DKG, types.OCR2VRF, types.GenericPlugin:
		return r.RelayerAdapter.NewPluginProvider(ctx, rargs, pargs)
	case types.CCIPCommit, types.CCIPExecution:
		return nil, fmt.Errorf("provider type not supported: %s", rargs.ProviderType)
	}
	return nil, fmt.Errorf("provider type not recognized: %s", rargs.ProviderType)
}
