package sequence

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/utils/operations/contract"
	fastcurseapi "github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
)

// CurseUncurseSeqInput groups subjects to curse and uncurse on the Aptos RMN Remote.
type AptosCurseInput struct {
	fastcurseapi.CurseInput
	CCIPAddress aptos.AccountAddress
}

// CurseUncurseSequence executes curse/uncurse operations on Aptos RMN Remote.
var SeqCurse = cldf_ops.NewSequence(
	"aptos-curse-uncurse-sequence",
	operation.Version1_0_0,
	"Curses or uncurses subjects on Aptos RMN Remote",
	curseUncurseSequence,
)

func curseUncurseSequence(b cldf_ops.Bundle, chain cldf_aptos.Chain, in AptosCurseInput) (output sequences.OnChainOutput, err error) {
	curseArgs := operation.CurseArgs{
		CCIPAddress: in.CCIPAddress,
		Subjects:    in.Subjects,
	}
	opOutput, err := cldf_ops.ExecuteOperation(b, operation.CurseOp, chain, curseArgs)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to curse with RMNRemote on chain %d: %w", chain.Selector, err)
	}
	batchOp, err := contract.NewBatchOperationFromWrites([]WriteOutput{opOutput.Output})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to create batch operation from writes: %w", err)
	}

	output.BatchOps = append(output.BatchOps, batchOp)
	return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batchOp}}, nil
}

var SeqUncurse = cldf_ops.NewSequence(
	"aptos-uncurse-sequence",
	operation.Version1_0_0,
	"Uncurses subjects on Aptos RMN Remote",
	uncurseSequence,
)

func uncurseSequence(b cldf_ops.Bundle, chain cldf_aptos.Chain, in AptosCurseInput) (output sequences.OnChainOutput, err error) {
	curseArgs := operation.CurseArgs{
		CCIPAddress: in.CCIPAddress,
		Subjects:    in.Subjects,
	}
	opOutput, err := cldf_ops.ExecuteOperation(b, operation.UncurseOp, chain, curseArgs)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to uncurse with RMNRemote on chain %d: %w", chain.Selector, err)
	}
	batchOp, err := NewBatchOperationFromWrites([]WriteOutput{opOutput.Output})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to create batch operation from writes: %w", err)
	}

	output.BatchOps = append(output.BatchOps, batchOp)
	return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batchOp}}, nil
}

// ExecInfo contains information about an executed transaction.
// Defined as a struct in case we want to add more fields in the future without breaking existing usage.
type ExecInfo struct {
	// Hash is the transaction hash.
	Hash string
}

// WriteOutput is the output of a write operation.
type WriteOutput struct {
	// ChainSelector is the selector of the target chain.
	ChainSelector uint64 `json:"chainSelector"`
	// Tx is the prepared transaction (in MCMS format).
	Tx mcms_types.Transaction `json:"tx"`
	// ExecInfo is populated if the write was executed, contains info about the executed transaction.
	ExecInfo *ExecInfo `json:"execInfo,omitempty"`
}

func (o WriteOutput) Executed() bool {
	return o.ExecInfo != nil
}

// NewBatchOperation constructs an MCMS BatchOperation from a slice of WriteOutputs.
// It filters out any WriteOutputs that have already been executed.
// Returns an error if the WriteOutputs target multiple chains.
// If all WriteOutputs are executed, it returns an empty BatchOperation and no error.
func NewBatchOperationFromWrites(outs []WriteOutput) (mcms_types.BatchOperation, error) {
	if len(outs) == 0 {
		return mcms_types.BatchOperation{}, nil
	}

	batchOps := make(map[uint64]mcms_types.BatchOperation)
	var chainSelector uint64
	for i, out := range outs {
		if out.Executed() {
			continue // Skip executed transactions, they should not be included.
		}
		if batchOp, exists := batchOps[out.ChainSelector]; !exists {
			if i != 0 {
				return mcms_types.BatchOperation{}, errors.New("failed to make batch operation: writes target multiple chains")
			}
			batchOps[out.ChainSelector] = mcms_types.BatchOperation{
				ChainSelector: mcms_types.ChainSelector(out.ChainSelector),
				Transactions:  []mcms_types.Transaction{out.Tx},
			}
			chainSelector = out.ChainSelector
		} else {
			batchOp.Transactions = append(batchOp.Transactions, out.Tx)
			batchOps[out.ChainSelector] = batchOp
		}
	}

	// If there are no unexecuted writes, return an empty BatchOperation.
	if len(batchOps) == 0 {
		return mcms_types.BatchOperation{}, nil
	}

	return batchOps[chainSelector], nil
}
