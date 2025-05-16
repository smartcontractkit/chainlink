package v1_6

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	evmStateView "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

type ApproveTokensForRouterConfig struct {
	ChainSelector uint64
	TokenAddress  common.Address
	Amount        *big.Int
}

func (cfg ApproveTokensForRouterConfig) Validate(e cldf.Environment) (evmStateView.CCIPChainState, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return evmStateView.CCIPChainState{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, found := state.Chains[cfg.ChainSelector]
	if !found {
		return evmStateView.CCIPChainState{}, fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	if cfg.TokenAddress.Cmp(common.Address{}) == 0 {
		return evmStateView.CCIPChainState{}, errors.New("token address cannot be empty or zero")
	}

	if chainState.Router == nil {
		return evmStateView.CCIPChainState{}, fmt.Errorf("router not found for chain selector %d", cfg.ChainSelector)
	}

	return chainState, nil
}

func ApproveTokensForRouter(e cldf.Environment, cfg ApproveTokensForRouterConfig) (cldf.ChangesetOutput, error) {
	state, err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate config: %w", err)
	}

	e.Logger.Infow("Approving Token(s) (approved account = '%s', token = '%s', amount = %d)",
		state.Router.Address().String(),
		cfg.TokenAddress.String(),
		cfg.Amount,
	)

	err = ccipChangeset.ApproveToken(
		e,
		cfg.ChainSelector,
		cfg.TokenAddress,
		state.Router.Address(),
		cfg.Amount,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

type ApproveTokensConfig struct {
	ChainSelector    uint64
	AddressToApprove common.Address
	TokenAddress     common.Address
	Amount           *big.Int
}

func (cfg ApproveTokensConfig) Validate(e cldf.Environment) (evmStateView.CCIPChainState, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return evmStateView.CCIPChainState{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, found := state.Chains[cfg.ChainSelector]
	if !found {
		return evmStateView.CCIPChainState{}, fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	if cfg.TokenAddress.Cmp(common.Address{}) == 0 {
		return evmStateView.CCIPChainState{}, errors.New("token address cannot be empty or zero")
	}

	if cfg.AddressToApprove.Cmp(common.Address{}) == 0 {
		return evmStateView.CCIPChainState{}, errors.New("address to approve cannot be empty or zero")
	}

	return chainState, nil
}

func ApproveTokens(e cldf.Environment, cfg ApproveTokensConfig) (cldf.ChangesetOutput, error) {
	_, err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate config: %w", err)
	}

	e.Logger.Infow("Approving Token(s) (approved account = '%s', token = '%s', amount = %d)",
		cfg.AddressToApprove.String(),
		cfg.TokenAddress.String(),
		cfg.Amount,
	)

	err = ccipChangeset.ApproveToken(
		e,
		cfg.ChainSelector,
		cfg.TokenAddress,
		cfg.AddressToApprove,
		cfg.Amount,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}
