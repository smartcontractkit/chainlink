package common

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type ChainAccessorFactoryParams struct {
	Lggr           logger.Logger
	Relayer        loop.Relayer
	ChainSelector  ccipocr3.ChainSelector
	ContractReader types.ContractReader
	ContractWriter types.ContractWriter
	AddrCodec      ccipocr3.AddressCodec
}

type ChainAccessorFactory interface {
	NewChainAccessor(ChainAccessorFactoryParams) (ccipocr3.ChainAccessor, error)
}
