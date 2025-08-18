package ccipevm

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// EVMChainAccessorFactory implements cciptypes.ChainAccessorFactory for EVM chains.
type EVMChainAccessorFactory struct{}

// NewChainAccessor creates a new chain accessor to be used for EVM chains.
func (f EVMChainAccessorFactory) NewChainAccessor(
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
