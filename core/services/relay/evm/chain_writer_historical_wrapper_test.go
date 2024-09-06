package evm

import (
	"context"
	"math/big"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	interfacetesttypes "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
	primitives "github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

// This wrapper is required to enable the ChainReader to access historical data
// Since the geth simulated backend doesn't support historical data, we use this
// thin wrapper.
type ChainWriterHistoricalWrapper struct {
	commontypes.ChainWriter
	cwh *ClientWithContractHistory
}

func NewChainWriterHistoricalWrapper(cw commontypes.ChainWriter, cwh *ClientWithContractHistory) *ChainWriterHistoricalWrapper {
	return &ChainWriterHistoricalWrapper{ChainWriter: cw, cwh: cwh}
}

func (cwhw *ChainWriterHistoricalWrapper) SubmitTransaction(ctx context.Context, contractName, method string, args any, transactionID string, toAddress string, meta *commontypes.TxMeta, value *big.Int) error {
	if primArgs, ok := args.(interfacetesttypes.PrimitiveArgs); ok {
		callArgs := interfacetesttypes.ExpectedGetLatestValueArgs{
			ContractName:    contractName,
			ReadName:        "GetAlterablePrimitiveValue",
			ConfidenceLevel: primitives.Unconfirmed,
			Params:          nil,
			ReturnVal:       nil,
		}
		err := cwhw.cwh.SetUintLatestValue(ctx, primArgs.Value, callArgs)
		if err != nil {
			return err
		}
	}
	return cwhw.ChainWriter.SubmitTransaction(ctx, contractName, method, args, transactionID, toAddress, meta, value)
}
