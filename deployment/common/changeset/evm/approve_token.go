package evm

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonCS "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

var (
	ApproveTokenGeneric   = cldf.CreateChangeSet(logicApproveToken, preconditionsApproveToken)
	ApproveTokenForRouter = cldf.CreateChangeSet(logicApproveTokenForRouter, preconditionsApproveTokenForRouter)
)

type ApproveTokensConfig struct {
	ChainSelector    uint64
	AddressToApprove common.Address
	TokenAddress     common.Address
	Amount           *big.Int
}

func preconditionsApproveToken(e cldf.Environment, cfg ApproveTokensConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	_, found := state.Chains[cfg.ChainSelector]
	if !found {
		return fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	if cfg.TokenAddress.Cmp(common.Address{}) == 0 {
		return errors.New("token address cannot be empty or zero")
	}

	if cfg.AddressToApprove.Cmp(common.Address{}) == 0 {
		return errors.New("address to approve cannot be empty or zero")
	}

	return nil
}

func logicApproveToken(e cldf.Environment, cfg ApproveTokensConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infof("Approving Token(s) (approved account = '%s', token = '%s', amount = %d)",
		cfg.AddressToApprove.String(),
		cfg.TokenAddress.String(),
		cfg.Amount,
	)

	err := commonCS.ApproveToken(
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

type ApproveTokensForRouterConfig struct {
	ChainSelector uint64
	TokenAddress  common.Address
	Amount        *big.Int
}

func preconditionsApproveTokenForRouter(e cldf.Environment, cfg ApproveTokensForRouterConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, found := state.Chains[cfg.ChainSelector]
	if !found {
		return fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	if cfg.TokenAddress.Cmp(common.Address{}) == 0 {
		return errors.New("token address cannot be empty or zero")
	}

	if chainState.Router == nil {
		return fmt.Errorf("router not found for chain selector %d", cfg.ChainSelector)
	}

	return nil
}

func logicApproveTokenForRouter(e cldf.Environment, cfg ApproveTokensForRouterConfig) (cldf.ChangesetOutput, error) {
	state, _ := stateview.LoadOnchainState(e)

	chainState := state.Chains[cfg.ChainSelector]

	e.Logger.Infof("Approving Token(s) (approved account = '%s', token = '%s', amount = %d)",
		chainState.Router.Address().String(),
		cfg.TokenAddress.String(),
		cfg.Amount,
	)

	err := commonCS.ApproveToken(
		e,
		cfg.ChainSelector,
		cfg.TokenAddress,
		chainState.Router.Address(),
		cfg.Amount,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}
