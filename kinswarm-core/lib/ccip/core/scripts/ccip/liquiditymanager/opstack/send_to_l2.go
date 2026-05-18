package opstack

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	l1AdapterABI  = abihelpers.MustParseABI(optimism_l1_bridge_adapter.OptimismL1BridgeAdapterMetaData.ABI)
	rebalancerABI = abihelpers.MustParseABI(liquiditymanager.LiquidityManagerMetaData.ABI)
)

// nolint
func SendToL2(
	env multienv.Env,
	l1ChainID uint64,
	l1BridgeAdapterAddress,
	l1TokenAddress,
	l2TokenAddress,
	l2Recipient common.Address,
	amount *big.Int,
) {
	// do some basic checks before proceeding
	l1Token, err := erc20.NewERC20(l1TokenAddress, env.Clients[l1ChainID])
	helpers.PanicErr(err)

	// check if we have enough balance otherwise approve will fail
	balance, err := l1Token.BalanceOf(nil, env.Transactors[l1ChainID].From)
	helpers.PanicErr(err)
	if balance.Cmp(amount) < 0 {
		panic(fmt.Sprintf("Insufficient balance, get more tokens or specify smaller amount: %s < %s", balance, amount))
	}

	// call the L1 adapter to send the funds to L2
	// first approve the L1 adapter to spend the tokens
	// check allowance so we don't approve unnecessarily
	allowance, err := l1Token.Allowance(nil, env.Transactors[l1ChainID].From, l1BridgeAdapterAddress)
	helpers.PanicErr(err)
	if allowance.Cmp(amount) < 0 {
		tx, err2 := l1Token.Approve(env.Transactors[l1ChainID], l1BridgeAdapterAddress, amount)
		helpers.PanicErr(err2)
		helpers.ConfirmTXMined(context.Background(), env.Clients[l1ChainID], tx, int64(l1ChainID),
			"Approve", amount.String(), "to", l1BridgeAdapterAddress.String())

		// check allowance
		allowance, err2 = l1Token.Allowance(nil, env.Transactors[l1ChainID].From, l1BridgeAdapterAddress)
		helpers.PanicErr(err2)
		if allowance.Cmp(amount) < 0 {
			panic(fmt.Sprintf("Allowance failed, expected %s, got %s", amount, allowance))
		}
	} else {
		fmt.Println("Allowance already set to", allowance, "for", l1BridgeAdapterAddress.String())
	}

	fmt.Println("Sending ERC20 to Optimism L2:", "\n",
		"l1TokenAddress:", l1TokenAddress.String(), "\n",
		"l2Recipient:", l2Recipient.String(), "\n",
		"amount:", amount)

	calldata, err := l1AdapterABI.Pack(
		"sendERC20",
		l1TokenAddress,
		l2TokenAddress,
		l2Recipient,
		amount,
		[]byte{}, // no bridge specific payload for L1 to OP stack L2
	)
	helpers.PanicErr(err)

	gasPrice, err := env.Clients[l1ChainID].SuggestGasPrice(context.Background())
	helpers.PanicErr(err)

	// Estimate gas of the bridging operation and multiply that by 1.8 since
	// optimism bridging costs are paid for in gas.
	gasCost, err := env.Clients[l1ChainID].EstimateGas(context.Background(), ethereum.CallMsg{
		From:     env.Transactors[l1ChainID].From,
		Gas:      1e6,
		GasPrice: gasPrice,
		To:       &l1BridgeAdapterAddress,
		Data:     calldata,
	})
	helpers.PanicErr(err)

	gasCost = scaleGasCost(gasCost)

	fmt.Println("Estimated gas cost for bridging operation, after scaling by 1.8x:", gasCost)

	l1Adapter, err := optimism_l1_bridge_adapter.NewOptimismL1BridgeAdapter(l1BridgeAdapterAddress, env.Clients[l1ChainID])
	helpers.PanicErr(err)

	tx, err := l1Adapter.SendERC20(&bind.TransactOpts{
		From:     env.Transactors[l1ChainID].From,
		Signer:   env.Transactors[l1ChainID].Signer,
		GasLimit: gasCost,
	}, l1TokenAddress, l2TokenAddress, l2Recipient, amount, []byte{})
	helpers.PanicErr(err)

	helpers.ConfirmTXMined(context.Background(), env.Clients[l1ChainID], tx, int64(l1ChainID),
		"SendERC20", amount.String(), "to", l2Recipient.String())
}

func scaleGasCost(gasCost uint64) uint64 {
	return gasCost * 18 / 10
}

// nolint
func SendToL2ViaRebalancer(
	env multienv.Env,
	l1ChainID,
	remoteChainID uint64,
	l1RebalancerAddress common.Address,
	amount *big.Int,
) {
	remoteChain, ok := chainsel.ChainByEvmChainID(remoteChainID)
	if !ok {
		panic(fmt.Sprintf("Chain ID %d not found in chain selectors", remoteChainID))
	}

	l1Rebalancer, err := liquiditymanager.NewLiquidityManager(l1RebalancerAddress, env.Clients[l1ChainID])
	helpers.PanicErr(err)

	// check if there is enough liquidity to transfer the provided amount.
	liquidity, err := l1Rebalancer.GetLiquidity(nil)
	helpers.PanicErr(err)
	if liquidity.Cmp(amount) < 0 {
		panic(fmt.Sprintf("Insufficient liquidity, add more tokens to the liquidity container or specify smaller amount: %s < %s", liquidity, amount))
	}

	// Estimate gas of the bridging operation and multiply that by 1.8 since
	// optimism bridging costs are paid for in gas.
	calldata, err := rebalancerABI.Pack("rebalanceLiquidity",
		remoteChain.Selector,
		amount,
		big.NewInt(0), // no eth fee for L1 -> L2
		[]byte{},      // no bridge specific payload for L1 to OP stack L2
	)
	helpers.PanicErr(err)

	gasPrice, err := env.Clients[l1ChainID].SuggestGasPrice(context.Background())
	helpers.PanicErr(err)

	// Estimate gas of the bridging operation and multiply that by 1.8 since
	// optimism bridging costs are paid for in gas.
	gasCost, err := env.Clients[l1ChainID].EstimateGas(context.Background(), ethereum.CallMsg{
		From:     env.Transactors[l1ChainID].From,
		Gas:      1e6,
		GasPrice: gasPrice,
		To:       &l1RebalancerAddress,
		Data:     calldata,
	})
	helpers.PanicErr(err)

	gasCost = scaleGasCost(gasCost)

	fmt.Println("Estimated gas cost for bridging operation, after scaling by 1.8x:", gasCost)

	tx, err := l1Rebalancer.RebalanceLiquidity(&bind.TransactOpts{
		From:     env.Transactors[l1ChainID].From,
		Signer:   env.Transactors[l1ChainID].Signer,
		GasLimit: gasCost,
	}, remoteChain.Selector, amount, big.NewInt(0), []byte{})
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), env.Clients[l1ChainID], tx, int64(l1ChainID),
		"RebalanceLiquidity", amount.String(), "to", remoteChain.Name)
}
