package sequence

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type DeployMCMSSeqInput struct {
	ChainSelector uint64
	MCMSConfig    types.MCMSWithTimelockConfigV2
}

type DeployMCMSSeqOutput struct {
	PackageID               string
	Objects                 mcmsops.DeployMCMSObjects
	AcceptOwnershipProposal mcms.TimelockProposal
}

var DeployMCMSSequence = operations.NewSequence(
	"deploy-sui-mcms-sequence",
	semver.MustParse("1.0.0"),
	"Deploy SUI MCMS contract and configure it",
	deploySuiMCMS,
)

func deploySuiMCMS(b operations.Bundle, deps sui_ops.OpTxDeps, in DeployMCMSSeqInput) (DeployMCMSSeqOutput, error) {
	suiInput := mcmsops.DeployMCMSSeqInput{
		ChainSelector: in.ChainSelector,
		Bypasser:      &in.MCMSConfig.Bypasser,
		Proposer:      &in.MCMSConfig.Proposer,
		Canceller:     &in.MCMSConfig.Canceller,
	}

	report, err := operations.ExecuteSequence(b, mcmsops.DeployMCMSSequence, deps, suiInput)
	if err != nil {
		return DeployMCMSSeqOutput{}, fmt.Errorf("failed to deploy SUI MCMS: %w", err)
	}

	return DeployMCMSSeqOutput{
		PackageID:               report.Output.PackageId,
		Objects:                 report.Output.Objects,
		AcceptOwnershipProposal: report.Output.AcceptOwnershipProposal,
	}, nil
}
