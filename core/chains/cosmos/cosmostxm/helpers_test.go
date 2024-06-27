package cosmostxm

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
)

func (txm *Txm) ORM() *ORM {
	return txm.orm
}

func (txm *Txm) ConfirmTx(ctx context.Context, tc cosmosclient.Reader, txHash string, broadcasted []int64, maxPolls int, pollPeriod time.Duration) error {
	return txm.confirmTx(ctx, tc, txHash, broadcasted, maxPolls, pollPeriod)
}

func (txm *Txm) ConfirmAnyUnconfirmed(ctx context.Context) {
	txm.confirmAnyUnconfirmed(ctx)
}

func (txm *Txm) MarshalMsg(msg sdk.Msg) (string, []byte, error) {
	return txm.marshalMsg(msg)
}

func (txm *Txm) SendMsgBatch(ctx context.Context) {
	txm.sendMsgBatch(ctx)
}
