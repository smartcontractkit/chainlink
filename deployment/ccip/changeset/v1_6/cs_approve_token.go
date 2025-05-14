package v1_6

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type ApproveTokenEVMConfig struct {
	ChainSelector uint64

	Amount *big.Int
}

func TokenApproveTransferEVMChangeset(e cldf.Environment, cfg ApproveTokenEVMConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainState, found := state.Chains[cfg.ChainSelector]
	if !found {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	routerAddress := chainState.Router.Address()
	tokenAddress := chainState.LinkToken.Address()

	err = ccipChangeset.ApproveToken(e, cfg.ChainSelector, tokenAddress, routerAddress, cfg.Amount)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

type ApproveTokenSolConfig struct {
	ChainSelector uint64

	AddressToApprove solana.PublicKey

	Amount   uint64
	Decimals uint8
}

func TokenApproveTransferSolChangeset(e cldf.Environment, cfg ApproveTokenSolConfig) (cldf.ChangesetOutput, error) {
	err := doApproveTokenTransfer(
		e,
		cfg.ChainSelector,
		cfg.AddressToApprove,
		cfg.Amount,
		cfg.Decimals,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

type ApproveTokenFeeBillingSignerSolConfig struct {
	ChainSelector uint64

	Amount   uint64
	Decimals uint8
}

func TokenApproveFeeBillingSigner(e cldf.Environment, cfg ApproveTokenFeeBillingSignerSolConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	billingSignerPDA, _, err := solState.FindFeeBillingSignerPDA(state.SolChains[cfg.ChainSelector].Router)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get billing signer PDA: %w", err)
	}

	err = doApproveTokenTransfer(
		e,
		cfg.ChainSelector,
		billingSignerPDA,
		cfg.Amount,
		cfg.Decimals,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func doApproveTokenTransfer(
	e cldf.Environment,
	chainSelector uint64,
	addressToApprove solana.PublicKey,
	amount uint64,
	decimals uint8,
) error {
	solChain, found := e.SolChains[chainSelector]
	if !found {
		return fmt.Errorf("failed to get chain for selector %d", chainSelector)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}

	chainState, found := state.SolChains[chainSelector]
	if !found {
		return fmt.Errorf("failed to get chain state for selector %d", chainSelector)
	}

	solTokenPubKey := state.SolChains[chainSelector].SPLTokens[0]

	tokenProgram, err := chainState.TokenToTokenProgram(solTokenPubKey)
	if err != nil {
		return fmt.Errorf("failed to get token program: %w", err)
	}

	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		tokenProgram,
		solTokenPubKey,
		solChain.DeployerKey.PublicKey(),
	)
	if err != nil {
		return fmt.Errorf("failed to find associated token address: %w", err)
	}

	ix1, err := solTokenUtil.TokenApproveChecked(
		amount,
		decimals,
		tokenProgram,
		deployerATA,
		solTokenPubKey,
		addressToApprove,
		solChain.DeployerKey.PublicKey(),
		solana.PublicKeySlice{},
	)
	if err != nil {
		return fmt.Errorf("failed to TokenApproveChecked: %w", err)
	}
	if err = solChain.Confirm([]solana.Instruction{ix1}); err != nil {
		e.Logger.Errorw("Failed to confirm instructions for TokenApproveChecked", "chain", solChain.String(), "err", err)
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to create approve instruction: %w", err)
	}

	return nil
}

func getActiveTokenPool(
	e *cldf.Environment,
	poolType solTestTokenPool.PoolType,
	selector uint64,
	metadata string,
) (solana.PublicKey, cldf.ContractType) {
	state, _ := stateview.LoadOnchainState(*e)
	chainState := state.SolChains[selector]
	switch poolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		if metadata == "" {
			return chainState.BurnMintTokenPools[shared.CLLMetadata], shared.BurnMintTokenPool
		}
		return chainState.BurnMintTokenPools[metadata], shared.BurnMintTokenPool
	case solTestTokenPool.LockAndRelease_PoolType:
		if metadata == "" {
			return chainState.LockReleaseTokenPools[shared.CLLMetadata], shared.LockReleaseTokenPool
		}
		return chainState.LockReleaseTokenPools[metadata], shared.LockReleaseTokenPool
	default:
		return solana.PublicKey{}, ""
	}
}
