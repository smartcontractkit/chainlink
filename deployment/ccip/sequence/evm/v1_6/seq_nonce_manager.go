package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipsharedops "github.com/smartcontractkit/chainlink/deployment/ccip/operation"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type NonceManagerApplyUpdatesSequenceInput struct {
	CallerUpdatesByChain []ccipops.NonceManagerUpdateAuthorizedCallerInput
	RampUpdatesByChain   []ccipops.NonceManagerApplyPreviousRampsUpdatesInput
	MCMSConfig           *proposalutils.TimelockConfig `json:"mcmsConfig,omitempty"`
}

var (
	UpdateNonceManagers = operations.NewSequence(
		"UpdateNonceManagers",
		semver.MustParse("1.0.0"),
		"Apply updates to the Nonce Manager contract across multiple EVM chains",
		func(b operations.Bundle, input NonceManagerApplyUpdatesSequenceInput, deps opsutil.ConfigureDependencies) (opsutil.OpOutput, error) {
			finalOutput := &opsutil.OpOutput{}

			// execute NonceManagerUpdateAuthorizedCallerOp
			for chainSel, update := range input.CallerUpdatesByChain {
				report, err := operations.ExecuteOperation(b, ccipops.NonceManagerUpdateAuthorizedCallerOp, deps, update)
				if err != nil {
					return report.Output, fmt.Errorf("failed to execute NonceManagerUpdateAuthorizedCallerOp on %d: %w", chainSel, err)
				}
				if err := finalOutput.Merge(report.Output); err != nil {
					return opsutil.OpOutput{}, fmt.Errorf("failed to merge output for chain %d: %w", chainSel, err)
				}
			}

			// execute NonceManagerPreviousRampsUpdatesOp
			for chainSel, update := range input.RampUpdatesByChain {
				report, err := operations.ExecuteOperation(b, ccipops.NonceManagerPreviousRampsUpdatesOp, deps, update)
				if err != nil {
					return report.Output, fmt.Errorf("failed to execute NonceManagerPreviousRampsUpdatesOp on %d: %w", chainSel, err)
				}
				if err := finalOutput.Merge(report.Output); err != nil {
					return opsutil.OpOutput{}, fmt.Errorf("failed to merge output for chain %d: %w", chainSel, err)
				}
			}

			// aggregate proposals
			if len(finalOutput.Proposals) > 0 {
				report, err := operations.ExecuteOperation(b, ccipsharedops.PostOpsAggregateProposals, deps, ccipsharedops.PostOpsInput{
					MCMSConfig: input.MCMSConfig,
					Proposals:  finalOutput.Proposals,
				})
				if err != nil {
					return opsutil.OpOutput{}, fmt.Errorf("failed to aggregate proposals: %w", err)
				}
				b.Logger.Infow("Generated proposal to Update NonceManagers")
				return opsutil.OpOutput{
					Proposals:                  report.Output,
					DescribedTimelockProposals: finalOutput.DescribedTimelockProposals,
				}, err
			}
			return *finalOutput, nil
		},
	)
)
