package sequence

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	aptosops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/aptos"
)

type AptosCurseUncurseInput struct {
	CCIPAddress   aptos.AccountAddress
	ChainSelector uint64
	Subjects      []fastcurse.Subject
}

var AptosCurseSequence = cldf_ops.NewSequence(
	"aptos-curse-sequence",
	operation.Version1_0_0,
	"Curse sequence for Aptos",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in AptosCurseUncurseInput) (output sequences.OnChainOutput, err error) {
		chain, ok := chains.AptosChains()[in.ChainSelector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
		}
		subjectBytes := make([][]byte, len(in.Subjects))
		for i, subject := range in.Subjects {
			subjectBytes[i] = subject[:]
		}
		curseInput := aptosops.CurseMultipleInput{
			CCIPAddress: in.CCIPAddress,
			Subjects:    subjectBytes,
		}
		deps := dependency.AptosDeps{
			AptosChain: chain,
		}
		report, err := cldf_ops.ExecuteOperation(b, aptosops.CurseMultipleOp, deps, curseInput)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to execute curse operation on Aptos chain %d: %w", chain.Selector, err)
		}
		batchOperation := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chain.Selector),
			Transactions:  []mcmstypes.Transaction{report.Output},
		}
		return sequences.OnChainOutput{
			BatchOps: []mcmstypes.BatchOperation{batchOperation},
		}, nil
	},
)

var AptosUncurseSequence = cldf_ops.NewSequence(
	"aptos-uncurse-sequence",
	operation.Version1_0_0,
	"Uncurse sequence for Aptos",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in AptosCurseUncurseInput) (output sequences.OnChainOutput, err error) {
		chain, ok := chains.AptosChains()[in.ChainSelector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
		}
		subjectBytes := make([][]byte, len(in.Subjects))
		for i, subject := range in.Subjects {
			subjectBytes[i] = subject[:]
		}
		uncurseInput := aptosops.UncurseMultipleInput{
			CCIPAddress: in.CCIPAddress,
			Subjects:    subjectBytes,
		}
		deps := dependency.AptosDeps{
			AptosChain: chain,
		}
		report, err := cldf_ops.ExecuteOperation(b, aptosops.UncurseMultipleOp, deps, uncurseInput)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to execute uncurse operation on Aptos chain %d: %w", chain.Selector, err)
		}
		batchOperation := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chain.Selector),
			Transactions:  []mcmstypes.Transaction{report.Output},
		}
		return sequences.OnChainOutput{
			BatchOps: []mcmstypes.BatchOperation{batchOperation},
		}, nil
	},
)
