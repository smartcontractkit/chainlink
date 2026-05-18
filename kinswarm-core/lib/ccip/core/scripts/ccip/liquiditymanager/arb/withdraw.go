package arb

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_l2_bridge_adapter"
)

// nolint
func WithdrawFromL2(
	env multienv.Env,
	l2ChainID uint64,
	l2BridgeAdapterAddress common.Address,
	amount *big.Int,
	l1ToAddress,
	l2TokenAddress,
	l1TokenAddress common.Address,
) {
	token, err := erc20.NewERC20(l2TokenAddress, env.Clients[l2ChainID])
	helpers.PanicErr(err)

	// check if we have enough balance
	balance, err := token.BalanceOf(nil, env.Transactors[l2ChainID].From)
	helpers.PanicErr(err)
	if balance.Cmp(amount) < 0 {
		panic(fmt.Errorf("not enough balance to withdraw, get more tokens or specify less amount, bal: %s, amount: %s",
			balance.String(), amount.String()))
	}

	l2Adapter, err := arbitrum_l2_bridge_adapter.NewArbitrumL2BridgeAdapter(l2BridgeAdapterAddress, env.Clients[l2ChainID])
	helpers.PanicErr(err)

	// Approve the adapter to receive the tokens
	tx, err := token.Approve(env.Transactors[l2ChainID], l2BridgeAdapterAddress, amount)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), env.Clients[l2ChainID], tx, int64(l2ChainID))

	// check the approval
	allowance, err := token.Allowance(nil, env.Transactors[l2ChainID].From, l2BridgeAdapterAddress)
	helpers.PanicErr(err)
	if allowance.Cmp(amount) < 0 {
		panic(fmt.Errorf("approval failed, allowance: %s, expected amount: %s", allowance.String(), amount.String()))
	}

	// at this point we should be able to withdraw the tokens to L1
	tx, err = l2Adapter.SendERC20(env.Transactors[l2ChainID],
		l2TokenAddress,
		l1TokenAddress,
		l1ToAddress,
		amount,
		[]byte{}, /* bridgeSpecificData, unused for arbitrum L2 adapter */
	)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), env.Clients[l2ChainID], tx, int64(l2ChainID))
}
