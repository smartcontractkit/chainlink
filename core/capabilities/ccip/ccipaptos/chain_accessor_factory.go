package ccipaptos

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// AptosChainAccessorFactory implements cciptypes.ChainAccessorFactory for Aptos chains.
type AptosChainAccessorFactory struct{}

// NewChainAccessor creates a new chain accessor to be used for Aptos chains.
func (f AptosChainAccessorFactory) NewChainAccessor(
	lggr logger.Logger,
	_ loop.Relayer,
	chainSelector cciptypes.ChainSelector,
	contractReader types.ContractReader,
	contractWriter types.ContractWriter,
	addrCodec cciptypes.AddressCodec,
) (cciptypes.ChainAccessor, error) {
	return chainaccessor.NewDefaultAccessor(
		lggr,
		chainSelector,
		contractreader.NewExtendedContractReader(contractReader),
		contractWriter,
		addrCodec,
	)
}
