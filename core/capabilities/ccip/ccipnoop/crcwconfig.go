package ccipnoop

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// NoopChainRWProvider is a struct that implements the NoopChainRWProvider interface for Solana chains.
type NoopChainRWProvider struct{}

// GetChainWriter NoopChainRWProvider returns a new ContractWriter for Solana chains.
func (g NoopChainRWProvider) GetChainWriter(ctx context.Context, pararms ccipcommon.ChainWriterProviderOpts) (types.ContractWriter, error) {
	return pararms.Relayer.NewContractWriter(ctx, nil)
}

// GetChainReader returns a new ContractReader for Solana chains.
func (g NoopChainRWProvider) GetChainReader(ctx context.Context, params ccipcommon.ChainReaderProviderOpts) (types.ContractReader, error) {
	return params.Relayer.NewContractReader(ctx, nil)
}
