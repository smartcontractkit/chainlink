package v1_6

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipsharedops "github.com/smartcontractkit/chainlink/deployment/ccip/operation"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type UpdateNonceManagerConfig struct {
	UpdatesByChain map[uint64]NonceManagerUpdate // source -> dest -> update
	MCMS           *proposalutils.TimelockConfig
}

type NonceManagerUpdate struct {
	AddedAuthCallers   []common.Address
	RemovedAuthCallers []common.Address
	PreviousRampsArgs  []PreviousRampCfg
}

type PreviousRampCfg struct {
	RemoteChainSelector uint64
	OverrideExisting    bool
	// Set these only if the prevOnRamp or prevOffRamp addresses are not required to be in nonce manager.
	// If one of the onRamp or OffRamp is set with non-zero address and other is set with zero address,
	// it will not be possible to update the previous ramps later unless OverrideExisting is set to true.
	AllowEmptyOnRamp  bool // If true, the prevOnRamp address can be 0x0.
	AllowEmptyOffRamp bool // If true, the prevOffRamp address can be 0x0.
}

func (cfg UpdateNonceManagerConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	for sourceSel, update := range cfg.UpdatesByChain {
		sourceChainState, ok := state.Chains[sourceSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", sourceSel)
		}
		if sourceChainState.NonceManager == nil {
			return fmt.Errorf("missing nonce manager for chain %d", sourceSel)
		}
		sourceChain, ok := e.BlockChains.EVMChains()[sourceSel]
		if !ok {
			return fmt.Errorf("missing chain %d in environment", sourceSel)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, sourceChain.DeployerKey.From, sourceChainState.Timelock.Address(), sourceChainState.OnRamp); err != nil {
			return fmt.Errorf("chain %s: %w", sourceChain.String(), err)
		}
		for _, prevRamp := range update.PreviousRampsArgs {
			if prevRamp.RemoteChainSelector == sourceSel {
				return errors.New("source and dest chain cannot be the same")
			}
			if _, ok := state.Chains[prevRamp.RemoteChainSelector]; !ok {
				return fmt.Errorf("dest chain %d not found in onchain state for chain %d", prevRamp.RemoteChainSelector, sourceSel)
			}
			// If one of the onRamp or OffRamp is set with non-zero address and other is set with zero address,
			// it will not be possible to update the previous ramps later.
			// Allow blank onRamp or offRamp only if AllowEmptyOnRamp or AllowEmptyOffRamp is set to true.
			// see https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/ccip/NonceManager.sol#L139-L142
			if !prevRamp.AllowEmptyOnRamp {
				if prevOnRamp := state.Chains[sourceSel].EVM2EVMOnRamp; prevOnRamp == nil ||
					prevOnRamp[prevRamp.RemoteChainSelector] == nil ||
					prevOnRamp[prevRamp.RemoteChainSelector].Address() == (common.Address{}) {
					return fmt.Errorf("no previous onramp for source chain %d and dest chain %d, "+
						"If you want to set zero address for onRamp, set AllowEmptyOnRamp to true", sourceSel, prevRamp.RemoteChainSelector)
				}
			}
			if !prevRamp.AllowEmptyOffRamp {
				if prevOffRamp := state.Chains[sourceSel].EVM2EVMOffRamp; prevOffRamp == nil ||
					prevOffRamp[prevRamp.RemoteChainSelector] == nil ||
					prevOffRamp[prevRamp.RemoteChainSelector].Address() == (common.Address{}) {
					return fmt.Errorf("no previous offramp for source chain %d and dest chain %d"+
						"If you want to set zero address for offRamp, set AllowEmptyOffRamp to true", prevRamp.RemoteChainSelector, sourceSel)
				}
			}
		}
	}
	return nil
}

var (
	UpdateNonceManagerSequence = operations.NewSequence(
		"UpdateNonceManagerSequence",
		semver.MustParse("1.0.0"),
		"Apply updates to the Nonce Manager contract across multiple EVM chains",
		func(b operations.Bundle, input UpdateNonceManagerConfig, deps opsutil.ConfigureDependencies) (opsutil.OpOutput, error) {
			finalOutput := &opsutil.OpOutput{}

			// validate input
			if err := input.Validate(deps.Env); err != nil {
				return opsutil.OpOutput{}, err
			}

			// load on-chain state
			s, err := stateview.LoadOnchainState(deps.Env)
			if err != nil {
				return opsutil.OpOutput{}, err
			}

			for chainSel, update := range input.UpdatesByChain {
				// build inputs
				callerOpInput := ccipops.NonceManagerUpdateAuthorizedCallerInput{
					ChainSelector: chainSel,
					Callers: nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
						AddedCallers:   update.AddedAuthCallers,
						RemovedCallers: update.RemovedAuthCallers,
					},
					MCMS: input.MCMS,
				}

				var rampUpdatesOpInput *ccipops.NonceManagerApplyPreviousRampsUpdatesInput

				if len(update.PreviousRampsArgs) > 0 {
					previousRampsArgs := make([]nonce_manager.NonceManagerPreviousRampsArgs, 0)
					for _, prevRamp := range update.PreviousRampsArgs {
						var onRamp, offRamp common.Address
						if !prevRamp.AllowEmptyOnRamp {
							onRamp = s.Chains[chainSel].EVM2EVMOnRamp[prevRamp.RemoteChainSelector].Address()
						}
						if !prevRamp.AllowEmptyOffRamp {
							offRamp = s.Chains[chainSel].EVM2EVMOffRamp[prevRamp.RemoteChainSelector].Address()
						}
						previousRampsArgs = append(previousRampsArgs, nonce_manager.NonceManagerPreviousRampsArgs{
							RemoteChainSelector:   prevRamp.RemoteChainSelector,
							OverrideExistingRamps: prevRamp.OverrideExisting,
							PrevRamps: nonce_manager.NonceManagerPreviousRamps{
								PrevOnRamp:  onRamp,
								PrevOffRamp: offRamp,
							},
						})
					}
					rampUpdatesOpInput = &(ccipops.NonceManagerApplyPreviousRampsUpdatesInput{
						ChainSelector:     chainSel,
						PreviousRampsArgs: previousRampsArgs,
						MCMS:              input.MCMS,
					})
				}

				// execute NonceManagerUpdateAuthorizedCallerOp
				report, err := operations.ExecuteOperation(b, ccipops.NonceManagerUpdateAuthorizedCallerOp, deps, callerOpInput)
				if err != nil {
					return report.Output, fmt.Errorf("failed to execute NonceManagerUpdateAuthorizedCallerOp on %d: %w", chainSel, err)
				}
				if err := finalOutput.Merge(report.Output); err != nil {
					return opsutil.OpOutput{}, fmt.Errorf("failed to merge output for chain %d: %w", chainSel, err)
				}

				// execute NonceManagerPreviousRampsUpdatesOp
				if rampUpdatesOpInput != nil {
					report, err := operations.ExecuteOperation(b, ccipops.NonceManagerPreviousRampsUpdatesOp, deps, *rampUpdatesOpInput)
					if err != nil {
						return report.Output, fmt.Errorf("failed to execute NonceManagerPreviousRampsUpdatesOp on %d: %w", chainSel, err)
					}
					if err := finalOutput.Merge(report.Output); err != nil {
						return opsutil.OpOutput{}, fmt.Errorf("failed to merge output for chain %d: %w", chainSel, err)
					}
				}
			}
			// if the MCMSConfig is not nil, we need to aggregate the proposals
			if len(finalOutput.Proposals) > 0 {
				report, err := operations.ExecuteOperation(b, ccipsharedops.PostOpsAggregateProposals, deps, ccipsharedops.PostOpsInput{
					MCMSConfig: input.MCMS,
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
