package changeset

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
)

// ApproveToken approves the router to spend the given amount of tokens
func ApproveToken(env cldf.Environment, src uint64, tokenAddress common.Address, routerAddress common.Address, amount *big.Int) error {
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
