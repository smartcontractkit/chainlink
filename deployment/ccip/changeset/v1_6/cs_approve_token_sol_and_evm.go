package v1_6

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"

	"github.com/smartcontractkit/chainlink/deployment"
)

type ApproveTokenEVMConfig struct {
	ChainSelector uint64
	TokenAddress  string
	RouterAddress string
	Amount        *big.Int
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
	tokenAddress := common.HexToAddress(cfg.TokenAddress)
	routerAddress := common.HexToAddress(cfg.RouterAddress)

	err := ApproveToken(e, cfg.ChainSelector, tokenAddress, routerAddress, cfg.Amount)

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

type ApproveTokenSolConfig struct {
	ChainSelector uint64
	TokenAddress  string
	RouterAddress string
	Amount        *big.Int
}

// func ApproveTokenTransferSolChangeset(e deployment.Environment, cfg ApproveTokenEVMConfig) (cldf.ChangesetOutput, error) {
// 	ixApprove, err := soltokens.TokenApproveChecked(1e9, 9, tokenProgram, deployerWSOL, wSOL, billingSignerPDA, deployer.PublicKey(), []solana.PublicKey{})

// 	if err != nil {
// 		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create approve instruction: %w", err)
// 	}

// 	return cldf.ChangesetOutput{}, nil
// }
