package sequence

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	rmn_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type SuiCurseUncurseInput struct {
	CCIPAddress          string
	CCIPObjectRef        string
	CCIPOwnerCapObjectId string
	ChainSelector        uint64
	Subjects             []fastcurse.Subject
}

var SuiCurseSequence = cldf_ops.NewSequence(
	"sui-curse-sequence",
	semver.MustParse("1.0.0"),
	"Curse sequence for SUI",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in SuiCurseUncurseInput) (sequences.OnChainOutput, error) {
		chain, ok := chains.SuiChains()[in.ChainSelector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("SUI chain with selector %d not found in environment", in.ChainSelector)
		}

		subjectBytes := make([][]byte, len(in.Subjects))
		for i, subject := range in.Subjects {
			s := subject
			subjectBytes[i] = s[:]
		}

		deps := sui_ops.OpTxDeps{
			Client: chain.Client,
			Signer: chain.Signer,
			GetCallOpts: func() *bind.CallOpts {
				gasBudget := uint64(400_000_000)
				return &bind.CallOpts{WaitForExecution: true, GasBudget: &gasBudget}
			},
			SuiRPC: chain.URL,
		}

		curseInput := rmn_ops.CurseUncurseChainInput{
			CCIPPackageId:    in.CCIPAddress,
			StateObjectId:    in.CCIPObjectRef,
			OwnerCapObjectId: in.CCIPOwnerCapObjectId,
			Subjects:         subjectBytes,
		}

		report, err := cldf_ops.ExecuteOperation(b, rmn_ops.CurseChainOp, deps, curseInput)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to execute curse operation on SUI chain %d: %w", in.ChainSelector, err)
		}

		call := report.Output.Call
		tx, err := suisdk.NewTransactionWithStateObj(
			call.Module, call.Function, call.PackageID,
			call.Data, call.Module, []string{},
			call.StateObjID, call.TypeArgs,
		)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to create MCMS transaction: %w", err)
		}

		return sequences.OnChainOutput{
			BatchOps: []mcmstypes.BatchOperation{{
				ChainSelector: mcmstypes.ChainSelector(in.ChainSelector),
				Transactions:  []mcmstypes.Transaction{tx},
			}},
		}, nil
	},
)

var SuiUncurseSequence = cldf_ops.NewSequence(
	"sui-uncurse-sequence",
	semver.MustParse("1.0.0"),
	"Uncurse sequence for SUI",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in SuiCurseUncurseInput) (sequences.OnChainOutput, error) {
		chain, ok := chains.SuiChains()[in.ChainSelector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("SUI chain with selector %d not found in environment", in.ChainSelector)
		}

		subjectBytes := make([][]byte, len(in.Subjects))
		for i, subject := range in.Subjects {
			s := subject
			subjectBytes[i] = s[:]
		}

		deps := sui_ops.OpTxDeps{
			Client: chain.Client,
			Signer: chain.Signer,
			GetCallOpts: func() *bind.CallOpts {
				gasBudget := uint64(400_000_000)
				return &bind.CallOpts{WaitForExecution: true, GasBudget: &gasBudget}
			},
			SuiRPC: chain.URL,
		}

		uncurseInput := rmn_ops.CurseUncurseChainInput{
			CCIPPackageId:    in.CCIPAddress,
			StateObjectId:    in.CCIPObjectRef,
			OwnerCapObjectId: in.CCIPOwnerCapObjectId,
			Subjects:         subjectBytes,
		}

		report, err := cldf_ops.ExecuteOperation(b, rmn_ops.UncurseChainOp, deps, uncurseInput)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to execute uncurse operation on SUI chain %d: %w", in.ChainSelector, err)
		}

		call := report.Output.Call
		tx, err := suisdk.NewTransactionWithStateObj(
			call.Module, call.Function, call.PackageID,
			call.Data, call.Module, []string{},
			call.StateObjID, call.TypeArgs,
		)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to create MCMS transaction: %w", err)
		}

		return sequences.OnChainOutput{
			BatchOps: []mcmstypes.BatchOperation{{
				ChainSelector: mcmstypes.ChainSelector(in.ChainSelector),
				Transactions:  []mcmstypes.Transaction{tx},
			}},
		}, nil
	},
)
