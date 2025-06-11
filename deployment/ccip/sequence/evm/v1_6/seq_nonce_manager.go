package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type NonceManagerUpdatesSequenceInput struct {
	UpdatesByChain map[uint64]NonceManagerUpdateInput `json:"updatesByChain"`
}

type NonceManagerUpdateInput struct {
	AuthorizedCallerArgs *opsutil.EVMCallInput[nonce_manager.AuthorizedCallersAuthorizedCallerArgs]
	PreviousRampsArgs    *opsutil.EVMCallInput[[]nonce_manager.NonceManagerPreviousRampsArgs]
}

var (
	UpdateNonceManagerSequence = operations.NewSequence(
		"UpdateNonceManagerSequence",
		semver.MustParse("1.0.0"),
		"Apply updates to the Nonce Manager contract across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input NonceManagerUpdatesSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {

			/*
				// load on-chain state
				s, err := stateview.LoadOnchainState(cldf.Environment)
				if err != nil {
					return opsutil.OpOutput{}, err
				}

				for chainSel, update := range input.UpdatesByChain {
					// build NonceManagerUpdateAuthorizedCallerOp input
					var callerOpInput *nonce_manager.AuthorizedCallersAuthorizedCallerArgs
					if len(update.AddedAuthCallers) > 0 || len(update.RemovedAuthCallers) > 0 {
						callerOpInput = &(nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
							AddedCallers:   update.AddedAuthCallers,
							RemovedCallers: update.RemovedAuthCallers,
						})
					}

					// build NonceManagerApplyPreviousRampsUpdatesOp input
					var rampUpdatesOpInput *[]nonce_manager.NonceManagerPreviousRampsArgs
					if len(update.PreviousRampsArgs) > 0 {
						rampUpdatesOpInput := make([]nonce_manager.NonceManagerPreviousRampsArgs, 0)
						for _, prevRamp := range update.PreviousRampsArgs {
							var onRamp, offRamp common.Address
							if !prevRamp.AllowEmptyOnRamp {
								onRamp = s.Chains[chainSel].EVM2EVMOnRamp[prevRamp.RemoteChainSelector].Address()
							}
							if !prevRamp.AllowEmptyOffRamp {
								offRamp = s.Chains[chainSel].EVM2EVMOffRamp[prevRamp.RemoteChainSelector].Address()
							}
							rampUpdatesOpInput = append(rampUpdatesOpInput, nonce_manager.NonceManagerPreviousRampsArgs{
								RemoteChainSelector:   prevRamp.RemoteChainSelector,
								OverrideExistingRamps: prevRamp.OverrideExisting,
								PrevRamps: nonce_manager.NonceManagerPreviousRamps{
									PrevOnRamp:  onRamp,
									PrevOffRamp: offRamp,
								},
							})
						}
					}
			*/
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))

			for chainSel, update := range input.UpdatesByChain {
				chainOutputs := []opsutil.EVMCallOutput{}
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}

				callerOpInput := update.AuthorizedCallerArgs
				rampUpdatesOpInput := update.PreviousRampsArgs

				// execute NonceManagerUpdateAuthorizedCallerOp
				if callerOpInput != nil {
					report, err := operations.ExecuteOperation(b, ccipops.NonceManagerUpdateAuthorizedCallerOp, chain, *callerOpInput)
					if err != nil {
						return nil, fmt.Errorf("failed to execute NonceManagerUpdateAuthorizedCallerOp on %s: %w", chain, err)
					}
					chainOutputs = append(chainOutputs, report.Output)
				}

				// execute NonceManagerPreviousRampsUpdatesOp
				if rampUpdatesOpInput != nil {
					report, err := operations.ExecuteOperation(b, ccipops.NonceManagerPreviousRampsUpdatesOp, chain, *rampUpdatesOpInput)
					if err != nil {
						return nil, fmt.Errorf("failed to execute NonceManagerPreviousRampsUpdatesOp on %s: %w", chain, err)
					}
					chainOutputs = append(chainOutputs, report.Output)
				}
				opOutputs[chainSel] = chainOutputs
			}

			return opOutputs, nil
		},
	)
)
