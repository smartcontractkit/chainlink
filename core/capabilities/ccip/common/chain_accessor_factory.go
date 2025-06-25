package common

import (
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ChainAccessorFactory interface {
	NewChainAccessor(
		lggr logger.Logger,
		relayer loop.Relayer,
		chainSelector cciptypes.ChainSelector,
		contractReader types.ContractReader,
		contractWriter types.ContractWriter,
		addrCodec cciptypes.AddressCodec,
	) (cciptypes.ChainAccessor, error)
}
