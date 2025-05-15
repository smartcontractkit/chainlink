package txmgrtest

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

func NewTestTxStore(t testing.TB, ds sqlutil.DataSource) txmgr.TestEvmTxStore {
	return txmgr.NewTxStore(ds, logger.Test(t))
}

func NewEthTx(fromAddress common.Address) txmgr.Tx {
	return txmgr.Tx{
		FromAddress:    fromAddress,
		ToAddress:      evmtestutils.NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		Value:          big.Int(assets.NewEthValue(142)),
		FeeLimit:       uint64(1000000000),
		State:          txmgrcommon.TxUnstarted,
	}
}

func MustInsertUnconfirmedEthTx(t testing.TB, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.Tx {
	broadcastAt := time.Now()
	chainID := evmtestutils.FixtureChainID
	for _, opt := range opts {
		switch v := opt.(type) {
		case time.Time:
			broadcastAt = v
		case *big.Int:
			chainID = v
		}
	}
	etx := NewEthTx(fromAddress)

	etx.BroadcastAt = &broadcastAt
	etx.InitialBroadcastAt = &broadcastAt
	n := types.Nonce(nonce)
	etx.Sequence = &n
	etx.State = txmgrcommon.TxUnconfirmed
	etx.ChainID = chainID
	require.NoError(t, txStore.InsertTx(evmtestutils.Context(t), &etx))
	return etx
}

func MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.Tx {
	etx := MustInsertUnconfirmedEthTx(t, txStore, nonce, fromAddress, opts...)
	attempt := NewLegacyEthTxAttempt(t, etx.ID)
	ctx := evmtestutils.Context(t)

	tx := evmtestutils.NewLegacyTransaction(uint64(nonce), evmtestutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertConfirmedEthTxWithLegacyAttempt(t testing.TB, txStore txmgr.TestEvmTxStore, nonce int64, broadcastBeforeBlockNum int64, fromAddress common.Address) txmgr.Tx {
	timeNow := time.Now()
	etx := NewEthTx(fromAddress)
	ctx := evmtestutils.Context(t)

	etx.BroadcastAt = &timeNow
	etx.InitialBroadcastAt = &timeNow
	n := types.Nonce(nonce)
	etx.Sequence = &n
	etx.State = txmgrcommon.TxConfirmed
	etx.MinConfirmations.SetValid(6)
	require.NoError(t, txStore.InsertTx(ctx, &etx))
	attempt := NewLegacyEthTxAttempt(t, etx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	etx.TxAttempts = append(etx.TxAttempts, attempt)
	return etx
}

func NewLegacyEthTxAttempt(t testing.TB, etxID int64) txmgr.TxAttempt {
	gasPrice := assets.NewWeiI(1)
	return txmgr.TxAttempt{
		ChainSpecificFeeLimit: 42,
		TxID:                  etxID,
		TxFee:                 gas.EvmFee{GasPrice: gasPrice},
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx: hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:        utils.NewHash(),
		State:       txmgrtypes.TxAttemptInProgress,
	}
}

func NewDynamicFeeEthTxAttempt(t *testing.T, etxID int64) txmgr.TxAttempt {
	gasTipCap := assets.NewWeiI(1)
	gasFeeCap := assets.NewWeiI(1)
	return txmgr.TxAttempt{
		TxType: 0x2,
		TxID:   etxID,
		TxFee: gas.EvmFee{
			DynamicFee: gas.DynamicFee{GasTipCap: gasTipCap, GasFeeCap: gasFeeCap},
		},
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx:           hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:                  utils.NewHash(),
		State:                 txmgrtypes.TxAttemptInProgress,
		ChainSpecificFeeLimit: 42,
	}
}

func AssertCount(t testing.TB, ds sqlutil.DataSource, tableName string, expected int64) {
	t.Helper()
	ctx := evmtestutils.Context(t)
	var count int64
	err := ds.GetContext(ctx, &count, fmt.Sprintf(`SELECT count(*) FROM %s;`, tableName))
	require.NoError(t, err)
	require.Equal(t, expected, count)
}
