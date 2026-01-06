package sequence

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"

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
	return opOutput.Output, nil
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
	return opOutput.Output, nil
}
