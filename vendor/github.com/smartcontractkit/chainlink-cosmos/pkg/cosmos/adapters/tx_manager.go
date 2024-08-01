package adapters

import (
	"context"

	cosmosSDK "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
)

type Msg struct {
	db.Msg

	// In memory only
	DecodedMsg cosmosSDK.Msg
}

type Msgs []Msg

func (tms Msgs) GetSimMsgs() client.SimMsgs {
	var msgs []client.SimMsg
	for i := range tms {
		msgs = append(msgs, client.SimMsg{
			ID:  tms[i].ID,
			Msg: tms[i].DecodedMsg,
		})
	}
	return msgs
}

func (tms Msgs) GetIDs() []int64 {
	ids := make([]int64, len(tms))
	for i := range tms {
		ids[i] = tms[i].ID
	}
	return ids
}

type MsgEnqueuer interface {
	// Enqueue enqueues msg for broadcast and returns its id.
	// Returns ErrMsgUnsupported for unsupported message types.
	Enqueue(ctx context.Context, contractID string, msg cosmosSDK.Msg) (int64, error)
}

// TxManager manages txs composed of batches of queued messages.
type TxManager interface {
	MsgEnqueuer

	// GetMsgs returns any messages matching ids.
	GetMsgs(ctx context.Context, ids ...int64) (Msgs, error)
	// GasPrice returns the gas price in ucosm.
	GasPrice() (cosmosSDK.DecCoin, error)
}
