package opstack

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_bridge_adapter_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/opstack/withdrawprover"
)

const (
	FinalizationActionProveWithdrawal    uint8 = 0
	FinalizationActionFinalizeWithdrawal uint8 = 1
)

// nolint
func FinalizeL1(
	env multienv.Env,
	l1ChainID,
	l2ChainID uint64,
	l1BridgeAdapterAddress common.Address,
	l2TxHash common.Hash,
) {
	l1Client, ok := env.Clients[l1ChainID]
	if !ok {
		panic(fmt.Sprintf("No L1 client found for chain ID %d", l1ChainID))
	}

	l2Client, ok := env.Clients[l2ChainID]
	if !ok {
		panic(fmt.Sprintf("No L2 client found for chain ID %d", l2ChainID))
	}

	l1Adapter, err := optimism_l1_bridge_adapter.NewOptimismL1BridgeAdapter(l1BridgeAdapterAddress, l1Client)
	helpers.PanicErr(err)

	encodedPayload := finalizationPayload(l2Client, l2TxHash)

	tx, err := l1Adapter.FinalizeWithdrawERC20(
		env.Transactors[l1ChainID],
		common.Address{}, // not used
		common.Address{}, // not used
		encodedPayload,   // finalization payload
	)
	helpers.PanicErr(err)

	helpers.ConfirmTXMined(context.Background(), env.Clients[l1ChainID], tx, int64(l1ChainID), "FinalizeWithdrawalTransaction")
}

// nolint
func FinalizeWithdrawalViaRebalancer(
	env multienv.Env,
	l1ChainID,
	l2ChainID,
	remoteChainID uint64,
	amount *big.Int,
	l1RebalancerAddress common.Address,
	l2TxHash common.Hash,
) {
	remoteChain, ok := chainsel.ChainByEvmChainID(remoteChainID)
	if !ok {
		panic(fmt.Sprintf("Chain ID %d not found in chain selectors", remoteChainID))
	}

	l1Client, ok := env.Clients[l1ChainID]
	if !ok {
		panic(fmt.Sprintf("No L1 client found for chain ID %d", l1ChainID))
	}

	l2Client, ok := env.Clients[l2ChainID]
	if !ok {
		panic(fmt.Sprintf("No L2 client found for chain ID %d", l2ChainID))
	}

	l1Rebalancer, err := liquiditymanager.NewLiquidityManager(l1RebalancerAddress, l1Client)
	helpers.PanicErr(err)

	encodedPayload := finalizationPayload(l2Client, l2TxHash)

	tx, err := l1Rebalancer.ReceiveLiquidity(
		env.Transactors[l1ChainID],
		remoteChain.Selector,
		amount,
		true, // shouldWrapNative
		encodedPayload,
	)
	helpers.PanicErr(err)

	helpers.ConfirmTXMined(context.Background(), env.Clients[l1ChainID], tx, int64(l1ChainID), "FinalizeWithdrawalTransaction")
}

func finalizationPayload(l2Client *ethclient.Client, l2TxHash common.Hash) []byte {
	receipt, err := l2Client.TransactionReceipt(context.Background(), l2TxHash)
	helpers.PanicErr(err)

	messagePassedLog := withdrawprover.GetMessagePassedLog(receipt.Logs)
	if messagePassedLog == nil {
		panic(fmt.Sprintf("No message passed log found in receipt %s", receipt.TxHash.String()))
	}

	messagePassed, err := withdrawprover.ParseMessagePassedLog(messagePassedLog)
	helpers.PanicErr(err)

	encodedFinalizeWithdrawal, err := encoderABI.Methods["encodeOptimismFinalizationPayload"].Inputs.Pack(
		optimism_l1_bridge_adapter_encoder.OptimismL1BridgeAdapterOptimismFinalizationPayload{
			WithdrawalTransaction: optimism_l1_bridge_adapter_encoder.TypesWithdrawalTransaction{
				Nonce:    messagePassed.Nonce,
				Sender:   messagePassed.Sender,
				Target:   messagePassed.Target,
				Value:    messagePassed.Value,
				GasLimit: messagePassed.GasLimit,
				Data:     messagePassed.Data,
			},
		},
	)
	helpers.PanicErr(err)

	// then encode the finalize withdraw erc20 payload next.
	encodedPayload, err := encoderABI.Methods["encodeFinalizeWithdrawalERC20Payload"].Inputs.Pack(
		optimism_l1_bridge_adapter_encoder.OptimismL1BridgeAdapterFinalizeWithdrawERC20Payload{
			Action: FinalizationActionFinalizeWithdrawal,
			Data:   encodedFinalizeWithdrawal,
		},
	)
	helpers.PanicErr(err)

	return encodedPayload
}
