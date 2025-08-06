package common

import (
	"fmt"

	sel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func WrapContractReaderForDefaultAccessor(
	contractReader types.ContractReader,
	chainSelector ccipocr3.ChainSelector,
	lggr logger.Logger,
) (contractreader.Extended, error) {
	chainFamily, err1 := sel.GetSelectorFamily(uint64(chainSelector))
	if err1 != nil {
		return nil, fmt.Errorf("failed to get chain family from selector: %w", err1)
	}
	chainID, err1 := sel.GetChainIDFromSelector(uint64(chainSelector))
	if err1 != nil {
		return nil, fmt.Errorf("failed to get chain id from selector: %w", err1)
	}
	reader := contractreader.NewExtendedContractReader(
		contractreader.NewObserverReader(contractReader, lggr, chainFamily, chainID),
	)
	if reader == nil {
		return nil, fmt.Errorf("failed to create extended contract reader for chain selector %d", chainSelector)
	}
	return reader, nil
}
