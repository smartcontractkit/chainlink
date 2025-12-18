package sequence

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
)

// CurseSubjectsSeqInput curses multiple subjects on Aptos RMN Remote.
type CurseSubjectsSeqInput struct {
	CCIPAddress aptos.AccountAddress
	Subjects    [][]byte
}

// UncurseSubjectsSeqInput lifts curse for multiple subjects on Aptos RMN Remote.
type UncurseSubjectsSeqInput struct {
	CCIPAddress aptos.AccountAddress
	Subjects    [][]byte
}

var CurseSubjectsSequence = operations.NewSequence(
	"curse-subjects-sequence",
	operation.Version1_0_0,
	"Curse multiple subjects via RMN Remote on Aptos",
	curseSubjectsSequence,
)

var UncurseSubjectsSequence = operations.NewSequence(
	"uncurse-subjects-sequence",
	operation.Version1_0_0,
	"Uncurse multiple subjects via RMN Remote on Aptos",
	uncurseSubjectsSequence,
)

func curseSubjectsSequence(b operations.Bundle, deps operation.AptosDeps, in CurseSubjectsSeqInput) (mcmstypes.BatchOperation, error) {
	report, err := operations.ExecuteOperation(
		b,
		operation.CurseMultipleOp,
		deps,
		operation.CurseMultipleInput{
			CCIPAddress: in.CCIPAddress,
			Subjects:    in.Subjects,
		},
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to execute curse-multiple-op: %w", err)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{report.Output},
	}, nil
}

func uncurseSubjectsSequence(b operations.Bundle, deps operation.AptosDeps, in UncurseSubjectsSeqInput) (mcmstypes.BatchOperation, error) {
	report, err := operations.ExecuteOperation(
		b,
		operation.UncurseMultipleOp,
		deps,
		operation.UncurseMultipleInput{
			CCIPAddress: in.CCIPAddress,
			Subjects:    in.Subjects,
		},
	)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to execute uncurse-multiple-op: %w", err)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{report.Output},
	}, nil
}
