package relayerset

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type RelayerSet interface {
	Get(ctx context.Context, relayID types.RelayID) (Relayer, error)
	GetAll(ctx context.Context) ([]Relayer, error)
}

type PluginArgs struct {
	TransmitterID string
	PluginConfig  []byte
}

type RelayArgs struct {
	ContractID         string
	RelayConfig        []byte
	ProviderType       string
	MercuryCredentials *types.MercuryCredentials
}

type Relayer interface {
	services.Service
	NewPluginProvider(rargs RelayArgs, pargs PluginArgs) (types.PluginProvider, error)
}
