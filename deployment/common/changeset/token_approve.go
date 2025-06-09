package changeset

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
)

// ApproveToken approves the router to spend the given amount of tokens
func ApproveToken(env cldf.Environment, src uint64, tokenAddress common.Address, routerAddress common.Address, amount *big.Int) error {
	evmChains := env.BlockChains.EVMChains()
	token, err := erc20.NewERC20(tokenAddress, evmChains[src].Client)
	if err != nil {
		return err
	}

	tx, err := token.Approve(evmChains[src].DeployerKey, routerAddress, amount)
	if err != nil {
		return err
	}

	_, err = evmChains[src].Confirm(tx)
	if err != nil {
		return err
	}

	return nil
}
