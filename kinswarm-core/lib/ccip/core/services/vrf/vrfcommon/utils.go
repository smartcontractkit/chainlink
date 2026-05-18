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

type RespCountEntry struct {
	RequestID string
	Count     int
}

func GetRespCounts(ctx context.Context, txm txmgr.TxManager, chainID *big.Int, confirmedBlockNum int64) (
	[]RespCountEntry,
	error,
) {
	counts := []RespCountEntry{}
	metaField := "RequestID"
	states := []txmgrtypes.TxState{txmgrcommon.TxUnconfirmed, txmgrcommon.TxUnstarted, txmgrcommon.TxInProgress}
	// Search for txes with a non-null meta field in the provided states
	unconfirmedTxes, err := txm.FindTxesWithMetaFieldByStates(ctx, metaField, states, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "getRespCounts failed due to error in FindTxesWithMetaFieldByStates")
	}
	// Fetch completed transactions only as far back as the given confirmedBlockNum. This avoids
	// a table scan of the whole table, which could be large if it is unpruned.
	var confirmedTxes []*txmgr.Tx
	confirmedTxes, err = txm.FindTxesWithMetaFieldByReceiptBlockNum(ctx, metaField, confirmedBlockNum, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "getRespCounts failed due to error in FindTxesWithMetaFieldByReceiptBlockNum")
	}
	txes := DedupeTxList(append(unconfirmedTxes, confirmedTxes...))
	respCountMap := make(map[string]int)
	// Consolidate the number of txes for each meta RequestID
	for _, tx := range txes {
		var meta *txmgrtypes.TxMeta[common.Address, common.Hash]
		meta, err = tx.GetMeta()
		if err != nil {
			return nil, errors.Wrap(err, "getRespCounts failed parsing tx meta field")
		}
		if meta != nil && meta.RequestID != nil {
			requestId := meta.RequestID.String()
			if _, exists := respCountMap[requestId]; !exists {
				respCountMap[requestId] = 0
			}
			respCountMap[requestId]++
		}
	}

	// Parse response count map into output
	for key, value := range respCountMap {
		respCountEntry := RespCountEntry{
			RequestID: key,
			Count:     value,
		}
		counts = append(counts, respCountEntry)
	}
	return counts, nil
}

func DedupeTxList(txes []*txmgr.Tx) []*txmgr.Tx {
	txIdMap := make(map[string]bool)
	dedupedTxes := []*txmgr.Tx{}
	for _, tx := range txes {
		if _, found := txIdMap[tx.GetID()]; !found {
			txIdMap[tx.GetID()] = true
			dedupedTxes = append(dedupedTxes, tx)
		}
	}
	return dedupedTxes
}
