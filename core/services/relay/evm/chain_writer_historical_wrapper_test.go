package evm

import (
	"context"
	"math/big"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	interfacetesttypes "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/bindings"
)

// This wrapper is required to enable the ChainReader to access historical data
// Since the geth simulated backend doesn't support historical data, we use this
// thin wrapper.
type ChainWriterHistoricalWrapper struct {
	commontypes.ContractWriter
	cwh *ClientWithContractHistory
}

func NewChainWriterHistoricalWrapper(cw commontypes.ContractWriter, cwh *ClientWithContractHistory) *ChainWriterHistoricalWrapper {
	return &ChainWriterHistoricalWrapper{ContractWriter: cw, cwh: cwh}
}

func (cwhw *ChainWriterHistoricalWrapper) SubmitTransaction(ctx context.Context, contractName, method string, args any, transactionID string, toAddress string, meta *commontypes.TxMeta, value *big.Int) error {
	alterablePrimitiveCall, newValue := cwhw.getPrimitiveValueIfPossible(args)
	if alterablePrimitiveCall {
		callArgs := interfacetesttypes.ExpectedGetLatestValueArgs{
			ContractName:    contractName,
			ReadName:        "GetAlterablePrimitiveValue",
			ConfidenceLevel: primitives.Unconfirmed,
			Params:          nil,
			ReturnVal:       nil,
		}
		err := cwhw.cwh.SetUintLatestValue(ctx, newValue, callArgs)
		if err != nil {
			return err
		}
	}
	return cwhw.ContractWriter.SubmitTransaction(ctx, contractName, method, args, transactionID, toAddress, meta, value)
}

func (cwhw *ChainWriterHistoricalWrapper) getPrimitiveValueIfPossible(args any) (bool, uint64) {
	primitiveArgs, alterablePrimitiveCall := args.(interfacetesttypes.PrimitiveArgs)
	var newValue uint64
	var alterablePrimitiveValue bindings.SetAlterablePrimitiveValueInput
	if alterablePrimitiveCall {
		newValue = primitiveArgs.Value
	} else {
		alterablePrimitiveValue, alterablePrimitiveCall = args.(bindings.SetAlterablePrimitiveValueInput)
		if alterablePrimitiveCall {
			newValue = alterablePrimitiveValue.Value
		}
	}
	return alterablePrimitiveCall, newValue
}
