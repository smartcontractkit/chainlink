package arb

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/abstract_arbitrum_token_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arb_node_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_gateway_router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_inbox"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_token_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	l1AdapterABI     = abihelpers.MustParseABI(arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterMetaData.ABI)
	nodeInterfaceABI = abihelpers.MustParseABI(arb_node_interface.NodeInterfaceMetaData.ABI)
)

// nolint
func SendToL2(
	env multienv.Env,
	l1ChainID,
	l2ChainID uint64,
	l1BridgeAdapterAddress,
	l1TokenAddress,
	l1RefundAddress,
	l2RefundAddress,
	l2Recipient common.Address,
	amount *big.Int,
) {
	if l1RefundAddress == (common.Address{}) {
		l1RefundAddress = env.Transactors[l1ChainID].From
	}
	if l2RefundAddress == (common.Address{}) {
		l2RefundAddress = l2Recipient
	}

	// do some basic checks before proceeding
	l1Token, err := erc20.NewERC20(l1TokenAddress, env.Clients[l1ChainID])
	helpers.PanicErr(err)

	// check if we have enough balance otherwise approve will fail
	balance, err := l1Token.BalanceOf(nil, env.Transactors[l1ChainID].From)
	helpers.PanicErr(err)
	if balance.Cmp(amount) < 0 {
		panic(fmt.Sprintf("Insufficient balance, get more tokens or specify smaller amount: %s < %s", balance, amount))
	}

	l1GatewayRouter, err := arbitrum_gateway_router.NewArbitrumGatewayRouter(ArbitrumContracts[l1ChainID]["L1GatewayRouter"], env.Clients[l1ChainID])
	helpers.PanicErr(err)

	params := populateFunctionParams(
		env,
		l1ChainID,
		l2ChainID,
		l1GatewayRouter,
		l1Token.Address(),
		l1RefundAddress,
		l2RefundAddress,
		l2Recipient,
		l1BridgeAdapterAddress,
		amount,
	)

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

	// transact with the bridge adapter to send funds cross-chain
	bridgeCalldata, err := l1AdapterABI.Pack("exposeSendERC20Params", arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterSendERC20Params{
		GasLimit:          params.GasLimit,
		MaxSubmissionCost: params.MaxSubmissionCost,
		MaxFeePerGas:      params.MaxFeePerGas,
	})
	helpers.PanicErr(err)
	bridgeCalldata = bridgeCalldata[4:] // remove the method id
	sendERC20Calldata, err := l1AdapterABI.Pack("sendERC20", l1TokenAddress, common.HexToAddress("0x0"), l2Recipient, amount, bridgeCalldata)
	helpers.PanicErr(err)

	fmt.Println("Sending ERC20 to Arbitrum L2:", "\n",
		"l1TokenAddress:", l1TokenAddress.String(), "\n",
		"l2Recipient:", l2Recipient.String(), "\n",
		"amount:", amount, "\n",
		"calldata:", hexutil.Encode(bridgeCalldata), "\n",
		"value:", params.Deposit)

	gasPrice, err := env.Clients[l1ChainID].SuggestGasPrice(context.Background())
	helpers.PanicErr(err)
	nonce, err := env.Clients[l1ChainID].PendingNonceAt(context.Background(), env.Transactors[l1ChainID].From)
	helpers.PanicErr(err)
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &l1BridgeAdapterAddress,
		GasPrice: gasPrice,
		Gas:      1e6,
		Value:    params.Deposit,
		Data:     sendERC20Calldata,
	})
	signedTx, err := env.Transactors[l1ChainID].Signer(env.Transactors[l1ChainID].From, rawTx)
	helpers.PanicErr(err)
	err = env.Clients[l1ChainID].SendTransaction(context.Background(), signedTx)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), env.Clients[l1ChainID], signedTx, int64(l1ChainID),
		"Calling SendERC20, amount:", amount.String(), ", on adapter at:", l1BridgeAdapterAddress.String())
}

type L1ToL2MessageGasParams struct {
	GasLimit          *big.Int
	MaxSubmissionCost *big.Int
	MaxFeePerGas      *big.Int
	Deposit           *big.Int
}

func populateFunctionParams(
	env multienv.Env,
	l1ChainID,
	l2ChainID uint64,
	l1GatewayRouter *arbitrum_gateway_router.ArbitrumGatewayRouter,
	l1TokenAddress,
	l1RefundAddress,
	l2RefundAddress,
	l2RecipientAddress,
	l1BridgeAdapterAddress common.Address,
	amount *big.Int,
) L1ToL2MessageGasParams {
	l1Client := env.Clients[l1ChainID]

	l1Gateway, err := l1GatewayRouter.GetGateway(nil, l1TokenAddress)
	helpers.PanicErr(err)

	// get the counterpart gateway on L2 from the L1 gateway
	// unfortunately we need to instantiate a new wrapper because the counterpartGateway field,
	// although it is public, is not accessible via a getter function on the token gateway interface
	abstractGateway, err := abstract_arbitrum_token_gateway.NewAbstractArbitrumTokenGateway(l1Gateway, l1Client)
	helpers.PanicErr(err)
	l2Gateway, err := abstractGateway.CounterpartGateway(nil)
	helpers.PanicErr(err)

	l1TokenGateway, err := arbitrum_token_gateway.NewArbitrumTokenGateway(l1Gateway, l1Client)
	helpers.PanicErr(err)

	retryableData := RetryableData{
		From:                l1Gateway,
		To:                  l2Gateway,
		ExcessFeeRefundAddr: l2RefundAddress,
		CallValueRefundAddr: l1RefundAddress,
		// this is the amount - see the arbitrum SDK.
		// https://github.com/OffchainLabs/arbitrum-sdk/blob/4c0d43abd5fcc5d219b20bc55e9d0ee152c01309/src/lib/assetBridger/ethBridger.ts#L318
		L2CallValue: amount,
		// 3 seems to work, but not sure if it's the best value
		// you definitely need a non-nil deposit for the NodeInterface call to succeed
		Deposit: big.NewInt(3),
		// MaxSubmissionCost: , // To be filled in
		// GasLimit: , // To be filled in
		// MaxFeePerGas: , // To be filled in
		// Data: , // To be filled in
	}

	// determine the finalizeInboundTransfer calldata
	finalizeInboundTransferCalldata, err := l1TokenGateway.GetOutboundCalldata(
		nil,
		l1TokenAddress,         // L1 token address
		l1BridgeAdapterAddress, // L1 sender address
		l2RecipientAddress,     // L2 recipient address
		amount,                 // token amount
		[]byte{},               // extra data (unused here)
	)
	helpers.PanicErr(err)
	retryableData.Data = finalizeInboundTransferCalldata

	fmt.Println("Constructed RetryableData", "\n",
		"From:", retryableData.From.String(), "\n",
		"To:", retryableData.To.String(), "\n",
		"L2CallValue:", retryableData.L2CallValue, "\n",
		"ExcessFeeRefundAddr:", retryableData.ExcessFeeRefundAddr.String(), "\n",
		"CallValueRefundAddr:", retryableData.CallValueRefundAddr.String(), "\n",
		"Data (finalizeInboundTransfer call):", hexutil.Encode(retryableData.Data))

	l1BaseFee, err := l1Client.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)
	estimates := estimateAll(env, l1ChainID, l2ChainID, retryableData, l1BaseFee)

	return estimates
}

func estimateAll(env multienv.Env, l1ChainID, l2ChainID uint64, rd RetryableData, l1BaseFee *big.Int) L1ToL2MessageGasParams {
	l2Client := env.Clients[l2ChainID]
	l2MaxFeePerGas := estimateMaxFeePerGasOnL2(l2Client)

	maxSubmissionFee := estimateSubmissionFee(env.Clients[l1ChainID], l1ChainID, l1BaseFee, uint64(len(rd.Data)))

	gasLimit := estimateRetryableGasLimit(l2Client, l2ChainID, rd)

	deposit := new(big.Int).Mul(gasLimit, l2MaxFeePerGas)
	deposit = deposit.Add(deposit, maxSubmissionFee)

	fmt.Println("estimated L1 -> L2 fees:", "\n",
		"gasLimit:", gasLimit, "\n",
		"maxSubmissionFee:", maxSubmissionFee, "\n",
		"maxFeePerGas (on L2):", l2MaxFeePerGas, "\n",
		"deposit:", deposit)

	return L1ToL2MessageGasParams{
		GasLimit:          gasLimit,
		MaxSubmissionCost: maxSubmissionFee,
		MaxFeePerGas:      l2MaxFeePerGas,
		Deposit:           deposit,
	}
}

// nolint
func estimateRetryableGasLimit(l2Client *ethclient.Client, l2ChainID uint64, rd RetryableData) *big.Int {
	packed, err := nodeInterfaceABI.Pack("estimateRetryableTicket",
		rd.From,
		assets.Ether(1).ToInt(), // this is what is done in the SDK, not sure why yet
		rd.To,
		rd.L2CallValue,
		rd.ExcessFeeRefundAddr,
		rd.CallValueRefundAddr,
		rd.Data)
	helpers.PanicErr(err)

	fmt.Println("calling node interface with calldata:", hexutil.Encode(packed), "value:", rd.Deposit)
	nodeInterfaceAddr := ArbitrumContracts[l2ChainID]["NodeInterface"]
	gasLimit, err := l2Client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &nodeInterfaceAddr,
		Data: packed,
	})
	helpers.PanicErr(err)

	// no percent increase on gas limit
	// should be pretty accurate
	return big.NewInt(int64(gasLimit))
}

func estimateMaxFeePerGasOnL2(l2Client *ethclient.Client) *big.Int {
	l2BaseFee, err := l2Client.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)
	// base fee on L2 is bumped by 200% by the arbitrum sdk (i.e 3x)
	l2BaseFee = new(big.Int).Mul(l2BaseFee, big.NewInt(3))
	return l2BaseFee
}

// nolint
func estimateSubmissionFee(l1Client *ethclient.Client, l1ChainID uint64, l1BaseFee *big.Int, calldataSize uint64) *big.Int {
	inbox, err := arbitrum_inbox.NewArbitrumInbox(ArbitrumContracts[l1ChainID]["L1Inbox"], l1Client)
	helpers.PanicErr(err)

	submissionFee, err := inbox.CalculateRetryableSubmissionFee(nil, big.NewInt(int64(calldataSize)), l1BaseFee)
	helpers.PanicErr(err)

	// submission fee is bumped by 300% (i.e 4x) by the arbitrum sdk
	// do the same here
	submissionFee = submissionFee.Mul(submissionFee, big.NewInt(4))

	return submissionFee
}

type RetryableData struct {
	// From is the gateway on L1 that will be sending the funds to the L2 gateway.
	From common.Address
	// To is the gateway on L2 that will be receiving the funds and eventually
	// sending them to the final recipient.
	To                common.Address
	L2CallValue       *big.Int
	Deposit           *big.Int
	MaxSubmissionCost *big.Int
	// ExcessFeeRefundAddr is an address on L2 that will be receiving excess fees
	ExcessFeeRefundAddr common.Address
	// CallValueRefundAddr is an address on L1 that will be receiving excess fees
	CallValueRefundAddr common.Address
	GasLimit            *big.Int
	MaxFeePerGas        *big.Int
	// Data is the calldata for the L2 gateway's `finalizeInboundTransfer` method.
	// The final recipient on L2 is specified in this calldata.
	Data []byte
}
