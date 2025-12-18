package sequence

import (
	"fmt"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
)

type DynamicSeqInput struct {
	Defs          []operations.Definition
	Inputs        []any // Each element should be the corresponding input type for its operation
	ChainSelector uint64
}

// DynamicSequence allows operations to be retrieved and executed at runtime
var DynamicSequence = operations.NewSequence(
	"aptos-dynamic-sequence",
	operation.Version1_0_0,
	"Allows operations to be retrieved and executed at runtime",
	dynamicSequence,
)

func dynamicSequence(b operations.Bundle, deps dependency.AptosDeps, in DynamicSeqInput) (mcmstypes.BatchOperation, error) {
	var txs []mcmstypes.Transaction

	for i, def := range in.Defs {
		op, err := b.OperationRegistry.Retrieve(def)
		if err != nil {
			return mcmstypes.BatchOperation{}, fmt.Errorf("failed to retrieve operation %s: %w", def.ID, err)
		}

		res, err := operations.ExecuteOperation(b, op, any(deps), in.Inputs[i])
		if err != nil {
			return mcmstypes.BatchOperation{}, fmt.Errorf("failed to execute operation %s: %w", def.ID, err)
		}

		// Handle different return types from operations
		switch output := res.Output.(type) {
		case mcmstypes.Transaction:
			txs = append(txs, output)
		case []mcmstypes.Transaction:
			txs = append(txs, output...)
		case mcmstypes.Operation:
			txs = append(txs, output.Transaction)
		case []mcmstypes.Operation:
			for _, op := range output {
				txs = append(txs, op.Transaction)
			}
		case mcmstypes.BatchOperation:
			txs = append(txs, output.Transactions...)
		default:
			return mcmstypes.BatchOperation{}, fmt.Errorf("operation %s returned unexpected type %T", def.ID, res.Output)
		}
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}, nil
}
