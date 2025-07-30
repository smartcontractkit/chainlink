package ccipaptos

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// AptosChainAccessorFactory implements cciptypes.ChainAccessorFactory for Aptos chains.
type AptosChainAccessorFactory struct{}

// NewChainAccessor creates a new chain accessor to be used for Aptos chains.
func (f AptosChainAccessorFactory) NewChainAccessor(
	params common.ChainAccessorFactoryParams,
) (ccipocr3.ChainAccessor, error) {
	return chainaccessor.NewDefaultAccessor(
		params.Lggr,
		params.ChainSelector,
		contractreader.NewExtendedContractReader(params.ContractReader),
		params.ContractWriter,
		params.AddrCodec,
	)
}
