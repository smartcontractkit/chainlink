package ccipsolana

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// SolanaChainAccessorFactory implements cciptypes.ChainAccessorFactory for Solana chains.
type SolanaChainAccessorFactory struct{}

// NewChainAccessor creates a new chain accessor to be used for Solana chains.
func (f SolanaChainAccessorFactory) NewChainAccessor(
	params common.ChainAccessorFactoryParams,
) (ccipocr3.ChainAccessor, error) {
	return chainaccessor.NewDefaultAccessor(
		params.Lggr,
		params.ChainSelector,
		params.ExtendedContractReader,
		params.ContractWriter,
		params.AddrCodec,
	)
}
