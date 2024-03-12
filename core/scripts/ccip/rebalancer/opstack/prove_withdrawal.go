package opstack

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rebalancer/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/optimism_portal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge/opstack/withdrawprover"
)

func ProveWithdrawal(
	env multienv.Env,
	l1ChainID,
	l2ChainID uint64,
	l1BridgeAdapterAddress,
	l2OutputOracleAddress,
	optimismPortalAddress common.Address,
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
	proveMessage(env, l1Client, l2Client, l2TxHash, optimismPortalAddress, l2OutputOracleAddress)
}

func proveMessage(
	env multienv.Env,
	l1Client, l2Client *ethclient.Client,
	l2TxHash common.Hash,
	optimismPortalAddress common.Address,
	l2OutputOracleAddress common.Address,
) {
	prover, err := withdrawprover.New(
		&ethClient{l1Client},
		&ethClient{l2Client},
		optimismPortalAddress,
		l2OutputOracleAddress,
		common.Address{}, // disputeGameFactoryAddress unknown at the moment
	)
	helpers.PanicErr(err)

	messageProof, err := prover.Prove(context.Background(), l2TxHash)
	helpers.PanicErr(err)

	optimismPortal, err := optimism_portal.NewOptimismPortal(optimismPortalAddress, l1Client)
	helpers.PanicErr(err)

	l1ChainID, err := l1Client.ChainID(context.Background())
	helpers.PanicErr(err)

	fmt.Println("Calling proveWithdrawalTransaction on OptimismPortal, nonce:", messageProof.LowLevelMessage.Nonce, "\n",
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

	tx, err := optimismPortal.ProveWithdrawalTransaction(env.Transactors[l1ChainID.Uint64()],
		optimism_portal.TypesWithdrawalTransaction{
			Nonce:    messageProof.LowLevelMessage.Nonce,
			Sender:   messageProof.LowLevelMessage.Sender,
			Target:   messageProof.LowLevelMessage.Target,
			Value:    messageProof.LowLevelMessage.Value,
			GasLimit: messageProof.LowLevelMessage.GasLimit,
			Data:     messageProof.LowLevelMessage.Data,
		},
		messageProof.L2OutputIndex,
		optimism_portal.TypesOutputRootProof{
			Version:                  messageProof.OutputRootProof.Version,
			StateRoot:                messageProof.OutputRootProof.StateRoot,
			MessagePasserStorageRoot: messageProof.OutputRootProof.MessagePasserStorageRoot,
			LatestBlockhash:          messageProof.OutputRootProof.LatestBlockHash,
		},
		messageProof.WithdrawalProof,
	)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), l1Client, tx, int64(l1ChainID.Uint64()), "ProveWithdrawal")
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
