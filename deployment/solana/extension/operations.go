package deployment_solana

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
)

// Common Solana dependencies
type SolanaDeps struct {
	Client         *solRpc.Client
	Auth           solana.PrivateKey
	SendAndConfirm func(
		ctx context.Context,
		rpcClient *rpc.Client,
		instructions []solana.Instruction,
		signer solana.PrivateKey,
		commitment rpc.CommitmentType,
		opts ...solCommonUtil.TxModifier,
	) (*solRpc.GetTransactionResult, error)
}

type SolanaTxResult struct {
	Address solana.PublicKey
	// TODO: Find the type
	Receipt solRpc.GetTransactionResult

	// TODO: Add any relevant info for tx ops
}
