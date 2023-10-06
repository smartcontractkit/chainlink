package txmgr

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type SequenceSyncer[ADDR types.Hashable, TX_HASH types.Hashable, BLOCK_HASH types.Hashable, SEQ types.Sequence[SEQ]] interface {
	Sync(ctx context.Context, addr ADDR, localNonce SEQ) (SEQ, error)
}
