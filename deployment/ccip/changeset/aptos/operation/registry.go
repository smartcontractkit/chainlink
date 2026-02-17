package operation

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/operation/aptos"
)

func GetAptosOperations() []*cld_ops.Operation[any, any, any] {
	// Since go has no Union types we cannot typecheck the output of operations
	// but only register operations returning:
	// - mcmstypes.Transaction
	// - []mcmstypes.Transaction
	// - mcmstypes.Operation
	// - []mcmstypes.Operation
	// - mcmstypes.BatchOperation
	// Otherwise dynamic sequence won't know how to handle the output
	var operations []*cld_ops.Operation[any, any, any]

	operations = append(operations, CCIPOperations...)
	operations = append(operations, FeeQuoterOperations...)
	operations = append(operations, MCMSOperations...)
	operations = append(operations, OffRampOperations...)
	operations = append(operations, OnRampOperations...)
	operations = append(operations, RouterOperations...)
	operations = append(operations, TokenAdminRegistryOperations...)
	operations = append(operations, TokenPoolOperations...)
	operations = append(operations, TokenOperations...)
	operations = append(operations, aptos.CurseMultipleOp.AsUntypedRelaxed())
	operations = append(operations, aptos.UncurseMultipleOp.AsUntypedRelaxed())

	return operations
}
