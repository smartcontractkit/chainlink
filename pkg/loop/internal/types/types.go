package internal

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type PluginRelayer interface {
	NewRelayer(ctx context.Context, config string, keystore core.Keystore) (Relayer, error)
}

type MedianProvider interface {
	NewMedianProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.MedianProvider, error)
}

type MercuryProvider interface {
	NewMercuryProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.MercuryProvider, error)
}

type FunctionsProvider interface {
	NewFunctionsProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.FunctionsProvider, error)
}

type AutomationProvider interface {
	NewAutomationProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.AutomationProvider, error)
}

type CCIPExecProvider interface {
	NewExecutionProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.CCIPExecProvider, error)
}

type CCIPCommitProvider interface {
	NewCommitProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.CCIPCommitProvider, error)
}

type OCR3CapabilityProvider interface {
	NewOCR3CapabilityProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.OCR3CapabilityProvider, error)
}

// Relayer is like types.Relayer, but with a dynamic NewPluginProvider method.
//
//go:generate mockery --quiet --name Relayer --output ./mocks/ --case=underscore
type Relayer interface {
	types.ChainService
	NewContractReader(ctx context.Context, contractReaderConfig []byte) (types.ContractReader, error)
	NewConfigProvider(context.Context, types.RelayArgs) (types.ConfigProvider, error)
	NewPluginProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.PluginProvider, error)
	NewLLOProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.LLOProvider, error)
}
