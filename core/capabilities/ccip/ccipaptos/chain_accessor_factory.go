package ccipaptos

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// AptosChainAccessorFactory implements cciptypes.ChainAccessorFactory for Aptos chains.
type AptosChainAccessorFactory struct{}

// NewChainAccessor creates a new chain accessor to be used for Aptos chains.
func (f AptosChainAccessorFactory) NewChainAccessor(
	params common.ChainAccessorFactoryParams,
) (ccipocr3.ChainAccessor, error) {
	reader, err := common.WrapContractReaderForDefaultAccessor(
		params.ContractReader,
		params.ChainSelector,
		params.Lggr,
	)
	if err != nil {
		return nil, err
	}
	return chainaccessor.NewDefaultAccessor(
		params.Lggr,
		params.ChainSelector,
		reader,
		params.ContractWriter,
		params.AddrCodec,
	)
}
