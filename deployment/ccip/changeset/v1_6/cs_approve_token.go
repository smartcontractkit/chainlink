package v1_6

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"

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

// ApproveToken approves the router to spend the given amount of tokens
func ApproveToken(env deployment.Environment, src uint64, tokenAddress common.Address, routerAddress common.Address, amount *big.Int) error {
	token, err := erc20.NewERC20(tokenAddress, env.Chains[src].Client)
	if err != nil {
		return err
	}

	tx, err := token.Approve(env.Chains[src].DeployerKey, routerAddress, amount)
	if err != nil {
		return err
	}

	_, err = env.Chains[src].Confirm(tx)
	if err != nil {
		return err
	}

	return nil
}

func ApproveTokenTransferEVMChangeset(e deployment.Environment, cfg ApproveTokenEVMConfig) (cldf.ChangesetOutput, error) {
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

	err = ApproveToken(e, cfg.ChainSelector, tokenAddress, routerAddress, cfg.Amount)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

type ApproveTokenSolConfig struct {
	ChainSelector uint64
	SolChain      cldf.SolChain

	SolTokenPubKey     solana.PublicKey
	RemoteTokenAccount solana.PublicKey
	TokenPubKey        solana.PublicKey

	Amount   uint64
	Decimals uint8
}

func ApproveTokenTransferSolChangeset(e deployment.Environment, cfg ApproveTokenSolConfig) (cldf.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainState, found := state.SolChains[cfg.ChainSelector]
	if !found {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain state for selector %d", cfg.ChainSelector)
	}

	tokenPool, _ := getActiveTokenPool(&e, solTestTokenPool.LockAndRelease_PoolType, cfg.ChainSelector, "")

	tokenProgram, err := chainState.TokenToTokenProgram(cfg.SolTokenPubKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get token program: %w", err)
	}
	poolSigner, err := solTokenUtil.TokenPoolSignerAddress(cfg.SolTokenPubKey, tokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get token pool signer address: %w", err)
	}

	ix1, err := solTokenUtil.TokenApproveChecked(
		cfg.Amount,
		cfg.Decimals,
		tokenProgram,
		cfg.RemoteTokenAccount,
		cfg.TokenPubKey,
		poolSigner,
		cfg.SolChain.DeployerKey.PublicKey(),
		solana.PublicKeySlice{},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to TokenApproveChecked: %w", err)
	}
	if err = cfg.SolChain.Confirm([]solana.Instruction{ix1}); err != nil {
		e.Logger.Errorw("Failed to confirm instructions for TokenApproveChecked", "chain", cfg.SolChain.String(), "err", err)
		return cldf.ChangesetOutput{}, err
	}

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create approve instruction: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
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
