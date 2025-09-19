package ccipnoop

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// chainRWProvider is a struct that implements the chainRWProvider.
type chainRWProvider struct{}

// GetChainWriter chainRWProvider returns a new noop ContractWriter.
func (n chainRWProvider) GetChainWriter(ctx context.Context, pararms common.ChainWriterProviderOpts) (types.ContractWriter, error) {
	return pararms.Relayer.NewContractWriter(ctx, nil)
}

// GetChainReader returns a new ContractReader for Solana chains.
func (n chainRWProvider) GetChainReader(ctx context.Context, params common.ChainReaderProviderOpts) (types.ContractReader, error) {
	return params.Relayer.NewContractReader(ctx, nil)
}
