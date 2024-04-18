package types

import (
	"context"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type StuckTxDetector[
	CHAIN_ID types.ID, // CHAIN_ID - chain id type
	ADDR types.Hashable, // ADDR - chain address type
	TX_HASH, BLOCK_HASH types.Hashable, // various chain hash types
	SEQ types.Sequence, // SEQ - chain sequence type (nonce, utxo, etc)
	FEE feetypes.Fee, // FEE - chain fee type
] interface {
	DetectStuckTransactions(ctx context.Context, enabledAddresses []ADDR, blockNum int64) ([]Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], error)
	SetPurgeBlockNum(fromAddress ADDR, blockNum int64)
}
