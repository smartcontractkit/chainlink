package testhelpers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/multicall3"
)

// CCIPSendCalldata packs the calldata for the Router's ccipSend method.
// This is expected to be used in Multicall scenarios (i.e multiple ccipSend calls
// in a single transaction).
func CCIPSendCalldata(
	destChainSelector uint64,
	evm2AnyMessage router.ClientEVM2AnyMessage,
) ([]byte, error) {
	calldata, err := routerABI.Methods["ccipSend"].Inputs.Pack(
		destChainSelector,
		evm2AnyMessage,
	)
	if err != nil {
		return nil, fmt.Errorf("pack ccipSend calldata: %w", err)
	}

	calldata = append(routerABI.Methods["ccipSend"].ID, calldata...)
	return calldata, nil
}

// GenMessagesForMulticall3 generates the calls and total value for the messages so that they can be used
// with a multicall3 transaction.
// Note that this is EVM-specific.
func GenMessagesForMulticall3(
	ctx context.Context,
	sourceRouter *router.Router,
	destChainSelector uint64,
	count int,
	baseMsg router.ClientEVM2AnyMessage,
) (calls []multicall3.Multicall3Call3Value, totalValue *big.Int, err error) {
	totalValue = big.NewInt(0)
	for range count {
		msg := router.ClientEVM2AnyMessage{
			Receiver:     baseMsg.Receiver,
			Data:         baseMsg.Data,
			TokenAmounts: baseMsg.TokenAmounts,
			FeeToken:     baseMsg.FeeToken,
			ExtraArgs:    baseMsg.ExtraArgs,
		}

		fee, err := sourceRouter.GetFee(&bind.CallOpts{Context: ctx}, destChainSelector, msg)
		if err != nil {
			return nil, nil, fmt.Errorf("router get fee: %w", err)
		}

		totalValue.Add(totalValue, fee)

		calldata, err := CCIPSendCalldata(destChainSelector, msg)
		if err != nil {
			return nil, nil, fmt.Errorf("generate calldata: %w", err)
		}

		calls = append(calls, multicall3.Multicall3Call3Value{
			Target:       sourceRouter.Address(),
			AllowFailure: false,
			CallData:     calldata,
			Value:        fee,
		})
	}

	return calls, totalValue, nil
}
