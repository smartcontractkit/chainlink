package common

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type ChainAccessorFactory interface {
	NewChainAccessor(
		lggr logger.Logger,
		relayer loop.Relayer,
		chainSelector ccipocr3.ChainSelector,
		contractReader types.ContractReader,
		contractWriter types.ContractWriter,
		addrCodec ccipocr3.AddressCodec,
	) (ccipocr3.ChainAccessor, error)
}
