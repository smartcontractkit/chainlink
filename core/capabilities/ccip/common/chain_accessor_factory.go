package common

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type ChainAccessorFactoryParams struct {
	Lggr                   logger.Logger
	Relayer                loop.Relayer
	ChainSelector          ccipocr3.ChainSelector
	ExtendedContractReader contractreader.Extended
	ContractWriter         types.ContractWriter
	AddrCodec              ccipocr3.AddressCodec
}

type ChainAccessorFactory interface {
	NewChainAccessor(ChainAccessorFactoryParams) (ccipocr3.ChainAccessor, error)
}
