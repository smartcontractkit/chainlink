package v1_6

import (
	"fmt"
	"math/big"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type ApproveTokenEVMConfig struct {
	ChainSelector uint64

	Amount *big.Int
}

func (cfg ApproveTokenEVMConfig) Validate(e cldf.Environment) (stateview.CCIPOnChainState, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return stateview.CCIPOnChainState{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, found := state.Chains[cfg.ChainSelector]
	if !found {
		return stateview.CCIPOnChainState{}, fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	if chainState.Router == nil {
		return stateview.CCIPOnChainState{}, fmt.Errorf("router not found for chain selector %d", cfg.ChainSelector)
	}

	if chainState.StaticLinkToken == nil {
		return stateview.CCIPOnChainState{}, fmt.Errorf("link token not found for chain selector %d", cfg.ChainSelector)
	}

	return state, nil
}

func TokenApproveTransferEVMChangeset(e cldf.Environment, cfg ApproveTokenEVMConfig) (cldf.ChangesetOutput, error) {
	state, err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate config: %w", err)
	}

	chainState := state.Chains[cfg.ChainSelector]

	err = ccipChangeset.ApproveToken(
		e,
		cfg.ChainSelector,
		chainState.Router.Address(),
		chainState.StaticLinkToken.Address(),
		cfg.Amount,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}
