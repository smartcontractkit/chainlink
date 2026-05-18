package opstack

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l1_bridge_adapter_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/opstack/withdrawprover"
)

var (
	encoderABI = abihelpers.MustParseABI(optimism_l1_bridge_adapter_encoder.OptimismL1BridgeAdapterEncoderMetaData.ABI)
)

// nolint
func ProveWithdrawal(
	env multienv.Env,
	l1ChainID,
	l2ChainID uint64,
	l1BridgeAdapterAddress,
	optimismPortalAddress,
	l2OutputOracleAddress common.Address,
	l2TxHash common.Hash,
) {
	l2Client, ok := env.Clients[l2ChainID]
	if !ok {
		panic(fmt.Sprintf("No client found for chain %d, map: %+v", l2ChainID, env.Clients))
	}
	l1Client, ok := env.Clients[l1ChainID]
	if !ok {
		panic(fmt.Sprintf("No client found for chain %d, map: %+v", l1ChainID, env.Clients))
	}

	encodedPayload := proveMessagePayload(l1Client, l2Client, l2TxHash, optimismPortalAddress, l2OutputOracleAddress)

	l1BridgeAdapter, err := optimism_l1_bridge_adapter.NewOptimismL1BridgeAdapter(l1BridgeAdapterAddress, l1Client)
	helpers.PanicErr(err)

	tx, err := l1BridgeAdapter.FinalizeWithdrawERC20(env.Transactors[l1ChainID],
		common.HexToAddress("0x0"), // doesn't matter
		common.HexToAddress("0x0"), // doesn't matter
		encodedPayload,             // all the data needed to prove withdrawal onchain.
	)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), l1Client, tx, int64(l1ChainID), "ProveWithdrawal")
}

// nolint
func ProveWithdrawalViaRebalancer(
	env multienv.Env,
	l1ChainID,
	l2ChainID,
	remoteChainID uint64,
	amount *big.Int,
	l1LiquidityManagerAddress,
	optimismPortalAddress,
	l2OutputOracleAddress common.Address,
	l2TxHash common.Hash,
) {
	remoteChain, ok := chainsel.ChainByEvmChainID(remoteChainID)
	if !ok {
		panic(fmt.Sprintf("Chain ID %d not found in chain selectors", remoteChainID))
	}

	l2Client, ok := env.Clients[l2ChainID]
	if !ok {
		panic(fmt.Sprintf("No client found for chain %d, map: %+v", l2ChainID, env.Clients))
	}
	l1Client, ok := env.Clients[l1ChainID]
	if !ok {
		panic(fmt.Sprintf("No client found for chain %d, map: %+v", l1ChainID, env.Clients))
	}

	encodedPayload := proveMessagePayload(l1Client, l2Client, l2TxHash, optimismPortalAddress, l2OutputOracleAddress)

	l1LiquidityManager, err := liquiditymanager.NewLiquidityManager(l1LiquidityManagerAddress, l1Client)
	helpers.PanicErr(err)

	tx, err := l1LiquidityManager.ReceiveLiquidity(
		env.Transactors[l1ChainID],
		remoteChain.Selector,
		amount,
		false,
		encodedPayload,
	)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), l1Client, tx, int64(l1ChainID),
		"ProveWithdrawal", amount.String(), "from", remoteChain.Name)
}

func CallGetFPACEnabled(
	env multienv.Env,
	l1ChainID,
	l2ChainID uint64,
) {
	l1Client, ok := env.Clients[l1ChainID]
	if !ok {
		panic(fmt.Sprintf("No L1 client found for chain ID %d", l1ChainID))
	}

	l2Client, ok := env.Clients[l2ChainID]
	if !ok {
		panic(fmt.Sprintf("No L2 client found for chain ID %d", l2ChainID))
	}

	prover, err := withdrawprover.New(
		&ethClient{l1Client},
		&ethClient{l2Client},
		OptimismContractsByChainID[l1ChainID]["OptimismPortalProxy"],
		OptimismContractsByChainID[l1ChainID]["L2OutputOracle"],
	)
	helpers.PanicErr(err)

	fpacEnabled, err := prover.GetFPAC(context.Background())
	helpers.PanicErr(err)
	fmt.Println("FPAC enabled:", fpacEnabled)
}

func proveMessagePayload(
	l1Client, l2Client *ethclient.Client,
	l2TxHash common.Hash,
	optimismPortalAddress,
	l2OutputOracleAddress common.Address,
) []byte {
	prover, err := withdrawprover.New(
		&ethClient{l1Client},
		&ethClient{l2Client},
		optimismPortalAddress,
		l2OutputOracleAddress,
	)
	helpers.PanicErr(err)

	messageProof, err := prover.Prove(context.Background(), l2TxHash)
	helpers.PanicErr(err)

	fmt.Println("Calling proveWithdrawalTransaction on bridge adapter, nonce:", messageProof.LowLevelMessage.Nonce, "\n",
		"sender:", messageProof.LowLevelMessage.Sender.String(), "\n",
		"target:", messageProof.LowLevelMessage.Target.String(), "\n",
		"value:", messageProof.LowLevelMessage.Value.String(), "\n",
		"gasLimit:", messageProof.LowLevelMessage.GasLimit.String(), "\n",
		"data:", hexutil.Encode(messageProof.LowLevelMessage.Data), "\n",
		"l2OutputIndex:", messageProof.L2OutputIndex, "\n",
		"outputRootProof version:", hexutil.Encode(messageProof.OutputRootProof.Version[:]), "\n",
		"outputRootProof stateRoot:", hexutil.Encode(messageProof.OutputRootProof.StateRoot[:]), "\n",
		"outputRootProof messagePasserStorageRoot:", hexutil.Encode(messageProof.OutputRootProof.MessagePasserStorageRoot[:]), "\n",
		"outputRootProof latestBlockHash:", hexutil.Encode(messageProof.OutputRootProof.LatestBlockHash[:]), "\n",
		"withdrawalProof:", formatWithdrawalProof(messageProof.WithdrawalProof))

	// encode the prove withdrawal payload first.
	encodedProveWithdrawal, err := encoderABI.Methods["encodeOptimismProveWithdrawalPayload"].Inputs.Pack(
		optimism_l1_bridge_adapter_encoder.OptimismL1BridgeAdapterOptimismProveWithdrawalPayload{
			WithdrawalTransaction: optimism_l1_bridge_adapter_encoder.TypesWithdrawalTransaction{
				Nonce:    messageProof.LowLevelMessage.Nonce,
				Sender:   messageProof.LowLevelMessage.Sender,
				Target:   messageProof.LowLevelMessage.Target,
				Value:    messageProof.LowLevelMessage.Value,
				GasLimit: messageProof.LowLevelMessage.GasLimit,
				Data:     messageProof.LowLevelMessage.Data,
			},
			L2OutputIndex: messageProof.L2OutputIndex,
			OutputRootProof: optimism_l1_bridge_adapter_encoder.TypesOutputRootProof{
				Version:                  messageProof.OutputRootProof.Version,
				StateRoot:                messageProof.OutputRootProof.StateRoot,
				MessagePasserStorageRoot: messageProof.OutputRootProof.MessagePasserStorageRoot,
				LatestBlockhash:          messageProof.OutputRootProof.LatestBlockHash,
			},
			WithdrawalProof: messageProof.WithdrawalProof,
		},
	)
	helpers.PanicErr(err)

	// then encode the finalize withdraw erc20 payload next.
	encodedPayload, err := encoderABI.Methods["encodeFinalizeWithdrawalERC20Payload"].Inputs.Pack(
		optimism_l1_bridge_adapter_encoder.OptimismL1BridgeAdapterFinalizeWithdrawERC20Payload{
			Action: FinalizationActionProveWithdrawal,
			Data:   encodedProveWithdrawal,
		},
	)
	helpers.PanicErr(err)

	return encodedPayload
}

func formatWithdrawalProof(proof [][]byte) string {
	var builder strings.Builder
	builder.WriteString("{")
	for i, p := range proof {
		builder.WriteString(hexutil.Encode(p))
		if i < len(proof)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("}")
	return builder.String()
}

type ethClient struct {
	*ethclient.Client
}

func (e *ethClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return e.Client.Client().CallContext(ctx, result, method, args...)
}
