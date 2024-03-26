package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type CCIPCommitProvider interface {
	PluginProvider

	NewCommitStoreReader(ctx context.Context, addr ccip.Address) (ccip.CommitStoreReader, error)
	NewOffRampReader(ctx context.Context, addr ccip.Address) (ccip.OffRampReader, error)
	NewOnRampReader(ctx context.Context, addr ccip.Address) (ccip.OnRampReader, error)
	NewPriceGetter(ctx context.Context) (ccip.PriceGetter, error)
	NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error)
	SourceNativeToken(ctx context.Context) (ccip.Address, error)
}

type CCIPExecProvider interface {
	PluginProvider

	NewCommitStoreReader(ctx context.Context, addr ccip.Address) (ccip.CommitStoreReader, error)
	NewOffRampReader(ctx context.Context, addr ccip.Address) (ccip.OffRampReader, error)
	NewOnRampReader(ctx context.Context, addr ccip.Address) (ccip.OnRampReader, error)
	NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error)
	NewTokenDataReader(ctx context.Context, tokenAddress ccip.Address) (ccip.TokenDataReader, error)
	NewTokenPoolBatchedReader(ctx context.Context) (ccip.TokenPoolBatchedReader, error)
	SourceNativeToken(ctx context.Context) (ccip.Address, error)
}

type CCIPCommitFactoryGenerator interface {
	NewCommitFactory(ctx context.Context, provider CCIPCommitProvider) (ReportingPluginFactory, error)
}

type CCIPExecutionFactoryGenerator interface {
	NewExecutionFactory(ctx context.Context, provider CCIPExecProvider) (ReportingPluginFactory, error)
}
type CCIPFactoryGenerator interface {
	CCIPCommitFactoryGenerator
	CCIPExecutionFactoryGenerator
}
