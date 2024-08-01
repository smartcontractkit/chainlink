package reader

import (
	readerinternal "github.com/smartcontractkit/chainlink-ccip/internal/reader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type CCIPReader = readerinternal.CCIP

func NewCCIPReader(
	lggr logger.Logger,
	contractReaders map[cciptypes.ChainSelector]types.ContractReader,
	contractWriters map[cciptypes.ChainSelector]types.ChainWriter,
	destChain cciptypes.ChainSelector,
) CCIPReader {
	return readerinternal.NewCCIPChainReader(lggr, contractReaders, contractWriters, destChain)
}

// NewCCIPReaderWithExtendedContractReaders can be used when you want to directly provide contractreader.Extended
func NewCCIPReaderWithExtendedContractReaders(
	lggr logger.Logger,
	contractReaders map[cciptypes.ChainSelector]contractreader.Extended,
	contractWriters map[cciptypes.ChainSelector]types.ChainWriter,
	destChain cciptypes.ChainSelector,
) CCIPReader {
	cr := readerinternal.NewCCIPChainReader(lggr, nil, contractWriters, destChain)
	for ch, extendedCr := range contractReaders {
		cr.WithExtendedContractReader(ch, extendedCr)
	}
	return cr
}
