package v1_6

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type ApproveTokenEVMConfig struct {
	ChainSelector uint64

	Amount *big.Int
}

func TokenApproveTransferEVMChangeset(e deployment.Environment, cfg ApproveTokenEVMConfig) (cldf.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
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

func TokenApproveTransferSolChangeset(e deployment.Environment, cfg ApproveTokenSolConfig) (cldf.ChangesetOutput, error) {
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

func TokenApproveFeeBillingSigner(e deployment.Environment, cfg ApproveTokenFeeBillingSignerSolConfig) (cldf.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	tokenPool, _ := getActiveTokenPool(&e, solTestTokenPool.LockAndRelease_PoolType, cfg.ChainSelector, "")

	solTokenPubKey := state.SolChains[cfg.ChainSelector].SPLTokens[0]

	poolSigner, err := solTokenUtil.TokenPoolSignerAddress(solTokenPubKey, tokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get token pool signer address: %w", err)
	}

	err = doApproveTokenTransfer(
		e,
		cfg.ChainSelector,
		poolSigner,
		cfg.Amount,
		cfg.Decimals,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func doApproveTokenTransfer(
	e deployment.Environment,
	chainSelector uint64,
	addressToApprove solana.PublicKey,
	amount uint64,
	decimals uint8,
) error {
	solChain, found := e.SolChains[chainSelector]
	if !found {
		return fmt.Errorf("failed to get chain for selector %d", chainSelector)
	}

	state, err := changeset.LoadOnchainState(e)
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
	e *deployment.Environment,
	poolType solTestTokenPool.PoolType,
	selector uint64,
	metadata string,
) (solana.PublicKey, cldf.ContractType) {
	state, _ := ccipChangeset.LoadOnchainState(*e)
	chainState := state.SolChains[selector]
	switch poolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		if metadata == "" {
			return chainState.BurnMintTokenPools[ccipChangeset.CLLMetadata], ccipChangeset.BurnMintTokenPool
		}
		return chainState.BurnMintTokenPools[metadata], ccipChangeset.BurnMintTokenPool
	case solTestTokenPool.LockAndRelease_PoolType:
		if metadata == "" {
			return chainState.LockReleaseTokenPools[ccipChangeset.CLLMetadata], ccipChangeset.LockReleaseTokenPool
		}
		return chainState.LockReleaseTokenPools[metadata], ccipChangeset.LockReleaseTokenPool
	default:
		return solana.PublicKey{}, ""
	}
}
