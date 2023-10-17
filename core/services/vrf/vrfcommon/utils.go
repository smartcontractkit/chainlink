package vrfcommon

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

type respCountEntry struct {
	RequestID string
	Count     int
}

func GetRespCounts(ctx context.Context, txm txmgr.TxManager, chainID *big.Int, confirmedBlockNum int64) (
	[]respCountEntry,
	error,
) {
	counts := []respCountEntry{}
	metaField := "RequestID"
	states := []txmgrtypes.TxState{txmgrcommon.TxUnconfirmed, txmgrcommon.TxUnstarted, txmgrcommon.TxInProgress}
	// Search for txes with a non-null meta field in the provided states
	unconfirmedTxes, err := txm.FindTxesWithMetaFieldByStates(ctx, metaField, states, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "getRespCounts failed due to error in FindTxesWithMetaFieldByStates")
	}
	// Fetch completed transactions only as far back as the given cutoffBlockNumber. This avoids
	// a table scan of the whole table, which could be large if it is unpruned.
	var confirmedTxes []*txmgr.Tx
	confirmedTxes, err = txm.FindTxesWithMetaFieldByReceiptBlockNum(ctx, metaField, confirmedBlockNum, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "getRespCounts failed due to error in FindTxesWithMetaFieldByReceiptBlockNum")
	}
	txes := append(unconfirmedTxes, confirmedTxes...)
	respCountMap := make(map[string]int)
	// Consolidate the number of txes for each meta RequestID
	for _, tx := range txes {
		var meta *txmgrtypes.TxMeta[common.Address, common.Hash]
		meta, err = tx.GetMeta()
		if err != nil {
			return nil, errors.Wrap(err, "getRespCounts failed parsing tx meta field")
		}
		// Query ensures the field is non-nil in the tx
		requestId := meta.RequestID.String()
		if _, exists := respCountMap[requestId]; !exists {
			respCountMap[requestId] = 0
		}
		count := respCountMap[requestId]
		respCountMap[requestId] = count + 1
	}

	// Parse response count map into output
	for key, value := range respCountMap {
		respCountEntry := respCountEntry{
			RequestID: key,
			Count:     value,
		}
		counts = append(counts, respCountEntry)
	}
	return counts, nil
}
