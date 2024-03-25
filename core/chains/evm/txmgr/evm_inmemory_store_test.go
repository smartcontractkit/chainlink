package txmgr_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInMemoryStore_FindTxAttemptsRequiringResend(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// insert the transaction into the persistent store
	inTx_1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 1, fromAddress, time.Unix(1616509200, 0))
	inTx_3 := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, persistentStore, 3, fromAddress, time.Unix(1616509400, 0))
	inTx_0 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 0, fromAddress, time.Unix(1616509100, 0))
	inTx_2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 2, fromAddress, time.Unix(1616509300, 0))
	// modify the attempts
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, inTx_0.ID)
	attempt0_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt0_2))

	attempt2_2 := newInProgressLegacyEthTxAttempt(t, inTx_2.ID)
	attempt2_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt2_2))

	attempt3_2 := cltest.NewDynamicFeeEthTxAttempt(t, inTx_3.ID)
	attempt3_2.TxFee.DynamicTipCap = assets.NewWeiI(10)
	attempt3_2.TxFee.DynamicFeeCap = assets.NewWeiI(20)
	attempt3_2.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_2))
	attempt3_4 := cltest.NewDynamicFeeEthTxAttempt(t, inTx_3.ID)
	attempt3_4.TxFee.DynamicTipCap = assets.NewWeiI(30)
	attempt3_4.TxFee.DynamicFeeCap = assets.NewWeiI(40)
	attempt3_4.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_4))
	attempt3_3 := cltest.NewDynamicFeeEthTxAttempt(t, inTx_3.ID)
	attempt3_3.TxFee.DynamicTipCap = assets.NewWeiI(20)
	attempt3_3.TxFee.DynamicFeeCap = assets.NewWeiI(30)
	attempt3_3.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_3))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_2))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_3))

	tcs := []struct {
		name                      string
		inOlderThan               time.Time
		inMaxInFlightTransactions uint32
		inChainID                 *big.Int
		inFromAddress             common.Address

		hasErr        bool
		hasTxAttempts bool
	}{
		//{"finds nothing if transactions from a different key", time.Now(), 10, chainID, utils.RandomAddress(), false, false},
		//{"returns the highest price attempt for each transaction that was last broadcast before or on the given time", time.Unix(1616509200, 0), 0, chainID, fromAddress, false, true},
		//{"returns the highest price attempt for EIP-1559 transactions", time.Unix(1616509400, 0), 0, chainID, fromAddress, false, true},
		{"applies limit", time.Unix(1616509200, 0), 1, chainID, fromAddress, false, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxAttempts, actErr := inMemoryStore.FindTxAttemptsRequiringResend(ctx, tc.inOlderThan, tc.inMaxInFlightTransactions, tc.inChainID, tc.inFromAddress)
			expTxAttempts, expErr := persistentStore.FindTxAttemptsRequiringResend(ctx, tc.inOlderThan, tc.inMaxInFlightTransactions, tc.inChainID, tc.inFromAddress)
			require.Equal(t, expErr, actErr)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxAttempts {
				require.NotEqual(t, 0, len(expTxAttempts))
				assert.NotEqual(t, 0, len(actTxAttempts))
				require.Equal(t, len(expTxAttempts), len(actTxAttempts))
				for i := 0; i < len(expTxAttempts); i++ {
					assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
				}
			} else {
				require.Equal(t, 0, len(expTxAttempts))
				require.Equal(t, 0, len(actTxAttempts))
			}
		})
	}
}

func TestInMemoryStore_FindTxesWithMetaFieldByReceiptBlockNum(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize the Meta field which is sqlutil.JSON
	subID := uint64(123)
	b, err := json.Marshal(txmgr.TxMeta{SubID: &subID})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	timeNow := time.Now()
	nonce := evmtypes.Nonce(123)
	blockNum := int64(3)
	broadcastBeforeBlockNum := int64(3)
	// initialize transactions
	inTx_0 := cltest.NewEthTx(fromAddress)
	inTx_0.BroadcastAt = &timeNow
	inTx_0.InitialBroadcastAt = &timeNow
	inTx_0.Sequence = &nonce
	inTx_0.State = commontxmgr.TxConfirmed
	inTx_0.MinConfirmations.SetValid(6)
	inTx_0.Meta = &meta
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx_0))
	attempt := cltest.NewLegacyEthTxAttempt(t, inTx_0.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt))
	inTx_0.TxAttempts = append(inTx_0.TxAttempts, attempt)
	// insert the transaction receipt into the persistent store
	rec_0 := mustInsertEthReceipt(t, persistentStore, 3, utils.NewHash(), inTx_0.TxAttempts[0].Hash)
	inTx_0.TxAttempts[0].Receipts = append(inTx_0.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec_0))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

	tcs := []struct {
		name        string
		inMetaField string
		inBlockNum  int64
		inChainID   *big.Int

		hasErr bool
		hasTxs bool
	}{
		{"successfully finds tx", "SubId", blockNum, chainID, false, true},
		{"unknown meta_field: finds no txs", "unknown", blockNum, chainID, false, false},
		{"incorrect meta_field: finds no txs", "MaxLink", blockNum, chainID, false, false},
		{"incorrect blockNum: finds no txs", "SubId", 12, chainID, false, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxs, actErr := inMemoryStore.FindTxesWithMetaFieldByReceiptBlockNum(ctx, tc.inMetaField, tc.inBlockNum, tc.inChainID)
			expTxs, expErr := persistentStore.FindTxesWithMetaFieldByReceiptBlockNum(ctx, tc.inMetaField, tc.inBlockNum, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxs {
				require.NotEqual(t, 0, len(expTxs))
				assert.NotEqual(t, 0, len(actTxs))
				require.Equal(t, len(expTxs), len(actTxs))
				for i := 0; i < len(expTxs); i++ {
					assertTxEqual(t, *expTxs[i], *actTxs[i])
				}
			} else {
				require.Equal(t, 0, len(expTxs))
				require.Equal(t, 0, len(actTxs))
			}
		})
	}
}

func TestInMemoryStore_FindTxesWithMetaFieldByStates(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize the Meta field which is sqlutil.JSON
	subID := uint64(123)
	b, err := json.Marshal(txmgr.TxMeta{SubID: &subID})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	// initialize transactions
	inTx_0 := cltest.NewEthTx(fromAddress)
	inTx_0.Meta = &meta
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx_0))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

	tcs := []struct {
		name        string
		inMetaField string
		inStates    []txmgrtypes.TxState
		inChainID   *big.Int

		hasErr bool
		hasTxs bool
	}{
		{"successfully finds tx", "SubId", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, true},
		{"incorrect state: finds no txs", "SubId", []txmgrtypes.TxState{commontxmgr.TxConfirmed}, chainID, false, false},
		{"unknown meta_field: finds no txs", "unknown", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
		{"incorrect meta_field: finds no txs", "MaxLink", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxs, actErr := inMemoryStore.FindTxesWithMetaFieldByStates(ctx, tc.inMetaField, tc.inStates, tc.inChainID)
			expTxs, expErr := persistentStore.FindTxesWithMetaFieldByStates(ctx, tc.inMetaField, tc.inStates, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if !tc.hasErr {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxs {
				require.NotEqual(t, 0, len(expTxs))
				assert.NotEqual(t, 0, len(actTxs))
				require.Equal(t, len(expTxs), len(actTxs))
				for i := 0; i < len(expTxs); i++ {
					assertTxEqual(t, *expTxs[i], *actTxs[i])
				}
			} else {
				require.Equal(t, 0, len(expTxs))
				require.Equal(t, 0, len(actTxs))
			}
		})
	}
}

func TestInMemoryStore_FindTxesByMetaFieldAndStates(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize the Meta field which is sqlutil.JSON
	subID := uint64(123)
	b, err := json.Marshal(txmgr.TxMeta{SubID: &subID})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	// initialize transactions
	inTx_0 := cltest.NewEthTx(fromAddress)
	inTx_0.Meta = &meta
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx_0))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

	tcs := []struct {
		name        string
		inMetaField string
		inMetaValue string
		inStates    []txmgrtypes.TxState
		inChainID   *big.Int

		hasErr bool
		hasTxs bool
	}{
		{"successfully finds tx", "SubId", "123", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, true},
		{"incorrect state: finds no txs", "SubId", "123", []txmgrtypes.TxState{commontxmgr.TxConfirmed}, chainID, false, false},
		{"incorrect meta_value: finds no txs", "SubId", "incorrect", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
		{"unknown meta_field: finds no txs", "unknown", "123", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
		{"incorrect meta_field: finds no txs", "JobID", "123", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxs, actErr := inMemoryStore.FindTxesByMetaFieldAndStates(ctx, tc.inMetaField, tc.inMetaValue, tc.inStates, tc.inChainID)
			expTxs, expErr := persistentStore.FindTxesByMetaFieldAndStates(ctx, tc.inMetaField, tc.inMetaValue, tc.inStates, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if !tc.hasErr {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxs {
				require.NotEqual(t, 0, len(expTxs))
				assert.NotEqual(t, 0, len(actTxs))
				require.Equal(t, len(expTxs), len(actTxs))
				for i := 0; i < len(expTxs); i++ {
					assertTxEqual(t, *expTxs[i], *actTxs[i])
				}
			} else {
				require.Equal(t, 0, len(expTxs))
				require.Equal(t, 0, len(actTxs))
			}
		})
	}
}

func TestInMemoryStore_FindTxWithIdempotencyKey(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	idempotencyKey := "777"
	inTx := cltest.NewEthTx(fromAddress)
	inTx.IdempotencyKey = &idempotencyKey
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

	tcs := []struct {
		name             string
		inIdempotencyKey string
		inChainID        *big.Int

		hasErr bool
		hasTx  bool
	}{
		{"no idempotency key", "", chainID, false, false},
		{"wrong idempotency key", "wrong", chainID, false, false},
		{"finds tx with idempotency key", idempotencyKey, chainID, false, true},
		{"wrong chain", idempotencyKey, big.NewInt(999), false, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTx, actErr := inMemoryStore.FindTxWithIdempotencyKey(ctx, tc.inIdempotencyKey, tc.inChainID)
			expTx, expErr := persistentStore.FindTxWithIdempotencyKey(ctx, tc.inIdempotencyKey, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if !tc.hasErr {
				require.Nil(t, actErr)
				require.Nil(t, expErr)
			}
			if tc.hasTx {
				require.NotNil(t, actTx)
				require.NotNil(t, expTx)
				assertTxEqual(t, *expTx, *actTx)
			} else {
				require.Nil(t, actTx)
				require.Nil(t, expTx)
			}
		})
	}
}

func TestInMemoryStore_CheckTxQueueCapacity(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	inTxs := []evmtxmgr.Tx{
		cltest.NewEthTx(fromAddress),
		cltest.NewEthTx(fromAddress),
	}
	for _, inTx := range inTxs {
		// insert the transaction into the persistent store
		require.NoError(t, persistentStore.InsertTx(ctx, &inTx))
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name           string
		inFromAddress  common.Address
		inMaxQueuedTxs uint64
		inChainID      *big.Int

		hasErr bool
	}{
		{"capacity reached", fromAddress, 2, chainID, true},
		{"above capacity", fromAddress, 1, chainID, true},
		{"below capacity", fromAddress, 3, chainID, false},
		{"wrong chain", fromAddress, 2, big.NewInt(999), false},
		{"wrong address", common.Address{}, 2, chainID, false},
		{"max queued txs is 0", fromAddress, 0, chainID, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actErr := inMemoryStore.CheckTxQueueCapacity(ctx, tc.inFromAddress, tc.inMaxQueuedTxs, tc.inChainID)
			expErr := persistentStore.CheckTxQueueCapacity(ctx, tc.inFromAddress, tc.inMaxQueuedTxs, tc.inChainID)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.NoError(t, expErr)
				require.NoError(t, actErr)
			}
		})
	}
}

func TestInMemoryStore_CountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize unstarted transactions
	inUnstartedTxs := []evmtxmgr.Tx{
		cltest.NewEthTx(fromAddress),
		cltest.NewEthTx(fromAddress),
	}
	for _, inTx := range inUnstartedTxs {
		// insert the transaction into the persistent store
		require.NoError(t, persistentStore.InsertTx(ctx, &inTx))
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name          string
		inFromAddress common.Address
		inChainID     *big.Int

		expUnstartedCount uint32
		hasErr            bool
	}{
		{"return correct total transactions", fromAddress, chainID, 2, false},
		{"invalid chain id", fromAddress, big.NewInt(999), 0, false},
		{"invalid address", common.Address{}, chainID, 0, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actMemoryCount, actErr := inMemoryStore.CountUnstartedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			actPersistentCount, expErr := persistentStore.CountUnstartedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.NoError(t, expErr)
				require.NoError(t, actErr)
			}
			assert.Equal(t, tc.expUnstartedCount, actMemoryCount)
			assert.Equal(t, tc.expUnstartedCount, actPersistentCount)
		})
	}
}

func TestInMemoryStore_CountUnconfirmedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize unconfirmed transactions
	inNonces := []int64{1, 2, 3}
	for _, inNonce := range inNonces {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, inNonce, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name          string
		inFromAddress common.Address
		inChainID     *big.Int

		expUnconfirmedCount uint32
		hasErr              bool
	}{
		{"return correct total transactions", fromAddress, chainID, 3, false},
		{"invalid chain id", fromAddress, big.NewInt(999), 0, false},
		{"invalid address", common.Address{}, chainID, 0, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actMemoryCount, actErr := inMemoryStore.CountUnconfirmedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			actPersistentCount, expErr := persistentStore.CountUnconfirmedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.NoError(t, expErr)
				require.NoError(t, actErr)
			}
			assert.Equal(t, tc.expUnconfirmedCount, actMemoryCount)
			assert.Equal(t, tc.expUnconfirmedCount, actPersistentCount)
		})
	}
}

func TestInMemoryStore_FindTxAttemptsConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize transactions
	inTxDatas := []struct {
		nonce                   int64
		broadcastBeforeBlockNum int64
		broadcastAt             time.Time
	}{
		{0, 1, time.Unix(1616509300, 0)},
		{1, 1, time.Unix(1616509400, 0)},
		{2, 1, time.Unix(1616509500, 0)},
	}
	for _, inTxData := range inTxDatas {
		// insert the transaction into the persistent store
		inTx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
			t, persistentStore, inTxData.nonce, inTxData.broadcastBeforeBlockNum,
			inTxData.broadcastAt, fromAddress,
		)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name      string
		inChainID *big.Int

		expTxAttemptsCount int
		hasError           bool
	}{
		{"finds tx attempts confirmed missing receipt", chainID, 3, false},
		{"wrong chain", big.NewInt(999), 0, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxAttempts, actErr := inMemoryStore.FindTxAttemptsConfirmedMissingReceipt(ctx, tc.inChainID)
			expTxAttempts, expErr := persistentStore.FindTxAttemptsConfirmedMissingReceipt(ctx, tc.inChainID)
			if tc.hasError {
				require.NotNil(t, actErr)
				require.NotNil(t, expErr)
			} else {
				require.NoError(t, actErr)
				require.NoError(t, expErr)
				require.Equal(t, tc.expTxAttemptsCount, len(expTxAttempts))
				require.Equal(t, tc.expTxAttemptsCount, len(actTxAttempts))
				for i := 0; i < len(expTxAttempts); i++ {
					assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
				}
			}
		})
	}
}

func TestInMemoryStore_FindTxAttemptsRequiringReceiptFetch(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize transactions
	inTxDatas := []struct {
		nonce                   int64
		broadcastBeforeBlockNum int64
		broadcastAt             time.Time
	}{
		{0, 1, time.Unix(1616509300, 0)},
		{1, 1, time.Unix(1616509400, 0)},
		{2, 1, time.Unix(1616509500, 0)},
	}
	for _, inTxData := range inTxDatas {
		// insert the transaction into the persistent store
		inTx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
			t, persistentStore, inTxData.nonce, inTxData.broadcastBeforeBlockNum,
			inTxData.broadcastAt, fromAddress,
		)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name      string
		inChainID *big.Int

		expTxAttemptsCount int
		hasError           bool
	}{
		{"finds tx attempts requiring receipt fetch", chainID, 3, false},
		{"wrong chain", big.NewInt(999), 0, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			expTxAttempts, expErr := persistentStore.FindTxAttemptsRequiringReceiptFetch(ctx, tc.inChainID)
			actTxAttempts, actErr := inMemoryStore.FindTxAttemptsRequiringReceiptFetch(ctx, tc.inChainID)
			if tc.hasError {
				require.NotNil(t, actErr)
				require.NotNil(t, expErr)
			} else {
				require.NoError(t, actErr)
				require.NoError(t, expErr)
				require.Equal(t, tc.expTxAttemptsCount, len(expTxAttempts))
				require.Equal(t, tc.expTxAttemptsCount, len(actTxAttempts))
				for i := 0; i < len(expTxAttempts); i++ {
					assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
				}
			}
		})
	}
}

func TestInMemoryStore_GetInProgressTxAttempts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("gets 0 in progress transaction", func(t *testing.T) {
		expTxAttempts, expErr := persistentStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		actTxAttempts, actErr := inMemoryStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		require.NoError(t, actErr)
		require.NoError(t, expErr)
		assert.Equal(t, len(expTxAttempts), len(actTxAttempts))
	})

	t.Run("gets 1 in progress transaction", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertUnconfirmedEthTxWithAttemptState(t, persistentStore, int64(7), fromAddress, txmgrtypes.TxAttemptInProgress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTxAttempts, expErr := persistentStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		actTxAttempts, actErr := inMemoryStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxAttempts), len(actTxAttempts))
		for i := 0; i < len(expTxAttempts); i++ {
			assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
		}
	})
}

func TestInMemoryStore_HasInProgressTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no in progress transaction", func(t *testing.T) {
		expExists, expErr := persistentStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		actExists, actErr := inMemoryStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		require.NoError(t, actErr)
		require.NoError(t, expErr)
		assert.Equal(t, expExists, actExists)
	})

	t.Run("has an in progress transaction", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 7, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expExists, expErr := persistentStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		actExists, actErr := inMemoryStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, expExists, actExists)
	})
}

func TestInMemoryStore_GetTxByID(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no transaction", func(t *testing.T) {
		expTx, expErr := persistentStore.GetTxByID(ctx, 0)
		actTx, actErr := inMemoryStore.GetTxByID(ctx, 0)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Nil(t, expTx)
		assert.Nil(t, actTx)
	})

	t.Run("successfully get transaction by ID", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 7, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTx, expErr := persistentStore.GetTxByID(ctx, inTx.ID)
		actTx, actErr := inMemoryStore.GetTxByID(ctx, inTx.ID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.NotNil(t, expTx)
		require.NotNil(t, actTx)
		assertTxEqual(t, *expTx, *actTx)
	})
}

func TestInMemoryStore_FindTxWithSequence(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTx, expErr := persistentStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		actTx, actErr := inMemoryStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Nil(t, expTx)
		assert.Nil(t, actTx)
	})

	t.Run("successfully get transaction by ID", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 666, 1, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTx, expErr := persistentStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		actTx, actErr := inMemoryStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.NotNil(t, expTx)
		require.NotNil(t, actTx)
		assertTxEqual(t, *expTx, *actTx)
	})

	t.Run("incorrect from address", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 777, 7, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		wrongFromAddress := common.Address{}
		expTx, expErr := persistentStore.FindTxWithSequence(ctx, wrongFromAddress, evmtypes.Nonce(777))
		actTx, actErr := inMemoryStore.FindTxWithSequence(ctx, wrongFromAddress, evmtypes.Nonce(777))
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Nil(t, expTx)
		require.Nil(t, actTx)
	})
}

func TestInMemoryStore_CountTransactionsByState(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expCount, expErr := persistentStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		actCount, actErr := inMemoryStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expCount, actCount)
	})
	t.Run("wrong chain id", func(t *testing.T) {
		wrongChainID := big.NewInt(999)
		expCount, expErr := persistentStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, wrongChainID)
		actCount, actErr := inMemoryStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, wrongChainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expCount, actCount)
	})
	t.Run("3 unconfirmed transactions", func(t *testing.T) {
		for i := int64(0); i < 3; i++ {
			// insert the transaction into the persistent store
			inTx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, i, fromAddress)
			// insert the transaction into the in-memory store
			require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
		}

		expCount, expErr := persistentStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		actCount, actErr := inMemoryStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expCount, actCount)
	})
}

func TestInMemoryStore_FindTxsRequiringResubmissionDueToInsufficientEth(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	_, otherAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTxs, expErr := persistentStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})

	// Insert order is mixed up to test sorting
	// insert the transaction into the persistent store
	inTx_2 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, persistentStore, 1, fromAddress)
	inTx_3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 2, fromAddress)
	attempt3_2 := cltest.NewLegacyEthTxAttempt(t, inTx_3.ID)
	attempt3_2.State = txmgrtypes.TxAttemptInsufficientFunds
	attempt3_2.TxFee.Legacy = assets.NewWeiI(100)
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_2))
	inTx_1 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, persistentStore, 0, fromAddress)
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_2))
	inTx_3.TxAttempts = append([]evmtxmgr.TxAttempt{attempt3_2}, inTx_3.TxAttempts...)
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_3))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

	// These should never be returned
	// insert the transaction into the persistent store
	otx_1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 3, fromAddress)
	otx_2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 4, 100, fromAddress)
	otx_3 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, persistentStore, 0, otherAddress)
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &otx_1))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &otx_2))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(otherAddress, &otx_3))

	t.Run("return all eth_txes with at least one attempt that is in insufficient_eth state", func(t *testing.T) {
		expTxs, expErr := persistentStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)

		assert.Equal(t, len(expTxs), len(actTxs))
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})

	t.Run("does not return txes with different chainID", func(t *testing.T) {
		wrongChainID := big.NewInt(999)
		expTxs, expErr := persistentStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, wrongChainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, wrongChainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})

	t.Run("does not return txes with different fromAddress", func(t *testing.T) {
		anotherFromAddress := common.Address{}
		expTxs, expErr := persistentStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, anotherFromAddress, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, anotherFromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})
}

func TestInMemoryStore_GetNonFatalTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTxs, expErr := persistentStore.GetNonFatalTransactions(ctx, chainID)
		actTxs, actErr := inMemoryStore.GetNonFatalTransactions(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})

	t.Run("get in progress, unstarted, and unconfirmed transactions", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx_0 := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 123, fromAddress)
		inTx_1 := mustCreateUnstartedGeneratedTx(t, persistentStore, fromAddress, chainID)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		expTxs, expErr := persistentStore.GetNonFatalTransactions(ctx, chainID)
		actTxs, actErr := inMemoryStore.GetNonFatalTransactions(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxs), len(actTxs))

		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})

	t.Run("wrong chain ID", func(t *testing.T) {
		wrongChainID := big.NewInt(999)
		expTxs, expErr := persistentStore.GetNonFatalTransactions(ctx, wrongChainID)
		actTxs, actErr := inMemoryStore.GetNonFatalTransactions(ctx, wrongChainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})
}

func TestInMemoryStore_FindTransactionsConfirmedInBlockRange(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTxs, expErr := persistentStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		actTxs, actErr := inMemoryStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})

	t.Run("find all transactions confirmed in range", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx_0 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 700, 8, fromAddress)
		rec_0 := mustInsertEthReceipt(t, persistentStore, 8, utils.NewHash(), inTx_0.TxAttempts[0].Hash)
		inTx_1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 777, 9, fromAddress)
		rec_1 := mustInsertEthReceipt(t, persistentStore, 9, utils.NewHash(), inTx_1.TxAttempts[0].Hash)
		// insert the transaction into the in-memory store
		inTx_0.TxAttempts[0].Receipts = append(inTx_0.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec_0))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		inTx_1.TxAttempts[0].Receipts = append(inTx_1.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec_1))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		expTxs, expErr := persistentStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		actTxs, actErr := inMemoryStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxs), len(actTxs))
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})

	t.Run("wrong chain ID", func(t *testing.T) {
		wrongChainID := big.NewInt(999)
		expTxs, expErr := persistentStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, wrongChainID)
		actTxs, actErr := inMemoryStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, wrongChainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})
}

func TestInMemoryStore_FindEarliestUnconfirmedBroadcastTime(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expBroadcastAt, expErr := persistentStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		actBroadcastAt, actErr := inMemoryStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.False(t, expBroadcastAt.Valid)
		assert.False(t, actBroadcastAt.Valid)
	})
	t.Run("find broadcast at time", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 123, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expBroadcastAt, expErr := persistentStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		actBroadcastAt, actErr := inMemoryStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.True(t, expBroadcastAt.Valid)
		require.True(t, actBroadcastAt.Valid)
		assert.Equal(t, expBroadcastAt.Time.Unix(), actBroadcastAt.Time.Unix())
	})
	t.Run("wrong chain ID", func(t *testing.T) {
		wrongChainID := big.NewInt(999)
		expBroadcastAt, expErr := persistentStore.FindEarliestUnconfirmedBroadcastTime(ctx, wrongChainID)
		actBroadcastAt, actErr := inMemoryStore.FindEarliestUnconfirmedBroadcastTime(ctx, wrongChainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.False(t, expBroadcastAt.Valid)
		assert.False(t, actBroadcastAt.Valid)
	})

}

func TestInMemoryStore_FindEarliestUnconfirmedTxAttemptBlock(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expBlock, expErr := persistentStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		actBlock, actErr := inMemoryStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.False(t, expBlock.Valid)
		assert.False(t, actBlock.Valid)
	})

	t.Run("find earliest unconfirmed tx block", func(t *testing.T) {
		broadcastBeforeBlockNum := int64(2)
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 123, fromAddress)
		attempt := cltest.NewLegacyEthTxAttempt(t, inTx.ID)
		attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
		attempt.State = txmgrtypes.TxAttemptBroadcast
		require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt))
		inTx.TxAttempts = append(inTx.TxAttempts, attempt)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expBlock, expErr := persistentStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		actBlock, actErr := inMemoryStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.True(t, expBlock.Valid)
		assert.True(t, actBlock.Valid)
		assert.Equal(t, expBlock.Int64, actBlock.Int64)
	})

	t.Run("wrong chain ID", func(t *testing.T) {
		wrongChainID := big.NewInt(999)
		expBlock, expErr := persistentStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, wrongChainID)
		actBlock, actErr := inMemoryStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, wrongChainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.False(t, expBlock.Valid)
		assert.False(t, actBlock.Valid)
	})
}

func TestInMemoryStore_LoadTxAttempts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("load tx attempt", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, persistentStore, 1, 7, time.Now(), fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTx := evmtxmgr.Tx{ID: inTx.ID, TxAttempts: []evmtxmgr.TxAttempt{}, FromAddress: fromAddress} // empty tx attempts for test
		expErr := persistentStore.LoadTxAttempts(ctx, &expTx)
		require.Equal(t, 1, len(expTx.TxAttempts))
		expAttempt := expTx.TxAttempts[0]

		actTx := evmtxmgr.Tx{ID: inTx.ID, TxAttempts: []evmtxmgr.TxAttempt{}, FromAddress: fromAddress} // empty tx attempts for test
		actErr := inMemoryStore.LoadTxAttempts(ctx, &actTx)
		require.Equal(t, 1, len(actTx.TxAttempts))
		actAttempt := actTx.TxAttempts[0]

		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assertTxAttemptEqual(t, expAttempt, actAttempt)
	})
}

func TestInMemoryStore_PreloadTxes(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("load transaction", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, int64(7), fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expAttempts := []evmtxmgr.TxAttempt{{ID: 0, TxID: inTx.ID}}
		expErr := persistentStore.PreloadTxes(ctx, expAttempts)
		require.Equal(t, 1, len(expAttempts))
		expAttempt := expAttempts[0]

		actAttempts := []evmtxmgr.TxAttempt{{ID: 0, TxID: inTx.ID}}
		actErr := inMemoryStore.PreloadTxes(ctx, actAttempts)
		require.Equal(t, 1, len(actAttempts))
		actAttempt := actAttempts[0]

		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assertTxAttemptEqual(t, expAttempt, actAttempt)
	})
}

func TestInMemoryStore_IsTxFinalized(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("tx not past finality depth", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 111, 1, fromAddress)
		rec := mustInsertEthReceipt(t, persistentStore, 1, utils.NewHash(), inTx.TxAttempts[0].Hash)
		// insert the transaction into the in-memory store
		inTx.TxAttempts[0].Receipts = append(inTx.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		blockHeight := int64(2)
		expIsFinalized, expErr := persistentStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		actIsFinalized, actErr := inMemoryStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expIsFinalized, actIsFinalized)
	})

	t.Run("tx is past finality depth", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 122, 2, fromAddress)
		rec := mustInsertEthReceipt(t, persistentStore, 2, utils.NewHash(), inTx.TxAttempts[0].Hash)
		// insert the transaction into the in-memory store
		inTx.TxAttempts[0].Receipts = append(inTx.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		blockHeight := int64(10)
		expIsFinalized, expErr := persistentStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		actIsFinalized, actErr := inMemoryStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expIsFinalized, actIsFinalized)
	})

	t.Run("wrong chain ID", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 133, 3, fromAddress)
		rec := mustInsertEthReceipt(t, persistentStore, 3, utils.NewHash(), inTx.TxAttempts[0].Hash)
		// insert the transaction into the in-memory store
		inTx.TxAttempts[0].Receipts = append(inTx.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		blockHeight := int64(10)
		wrongChainID := big.NewInt(999)
		expIsFinalized, expErr := persistentStore.IsTxFinalized(ctx, blockHeight, inTx.ID, wrongChainID)
		actIsFinalized, actErr := inMemoryStore.IsTxFinalized(ctx, blockHeight, inTx.ID, wrongChainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expIsFinalized, actIsFinalized)
	})
}

func TestInMemoryStore_FindTxsRequiringGasBump(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("gets transactions requiring gas bumping", func(t *testing.T) {
		currentBlockNum := int64(10)

		// insert the transaction into the persistent store
		inTx_0 := mustInsertUnconfirmedEthTxWithAttemptState(t, persistentStore, 1, fromAddress, txmgrtypes.TxAttemptBroadcast)
		require.NoError(t, persistentStore.SetBroadcastBeforeBlockNum(ctx, currentBlockNum, chainID))
		inTx_1 := mustInsertUnconfirmedEthTxWithAttemptState(t, persistentStore, 2, fromAddress, txmgrtypes.TxAttemptBroadcast)
		require.NoError(t, persistentStore.SetBroadcastBeforeBlockNum(ctx, currentBlockNum+1, chainID))
		// insert the transaction into the in-memory store
		inTx_0.TxAttempts[0].BroadcastBeforeBlockNum = &currentBlockNum
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		tempCurrentBlockNum := currentBlockNum + 1
		inTx_1.TxAttempts[0].BroadcastBeforeBlockNum = &tempCurrentBlockNum
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		newBlock := int64(12)
		gasBumpThreshold := int64(2)
		expTxs, expErr := persistentStore.FindTxsRequiringGasBump(ctx, fromAddress, newBlock, gasBumpThreshold, 0, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringGasBump(ctx, fromAddress, newBlock, gasBumpThreshold, 0, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxs), len(actTxs))
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})
}

// assertTxEqual asserts that two transactions are equal
func assertTxEqual(t *testing.T, exp, act evmtxmgr.Tx) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.IdempotencyKey, act.IdempotencyKey)
	assert.Equal(t, exp.Sequence, act.Sequence)
	assert.Equal(t, exp.FromAddress, act.FromAddress)
	assert.Equal(t, exp.ToAddress, act.ToAddress)
	assert.Equal(t, exp.EncodedPayload, act.EncodedPayload)
	assert.Equal(t, exp.Value, act.Value)
	assert.Equal(t, exp.FeeLimit, act.FeeLimit)
	assert.Equal(t, exp.Error, act.Error)
	if exp.BroadcastAt != nil {
		require.NotNil(t, act.BroadcastAt)
		assert.Equal(t, exp.BroadcastAt.Unix(), act.BroadcastAt.Unix())
	} else {
		assert.Equal(t, exp.BroadcastAt, act.BroadcastAt)
	}
	if exp.InitialBroadcastAt != nil {
		require.NotNil(t, act.InitialBroadcastAt)
		assert.Equal(t, exp.InitialBroadcastAt.Unix(), act.InitialBroadcastAt.Unix())
	} else {
		assert.Equal(t, exp.InitialBroadcastAt, act.InitialBroadcastAt)
	}
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.Meta, act.Meta)
	assert.Equal(t, exp.Subject, act.Subject)
	assert.Equal(t, exp.ChainID, act.ChainID)
	assert.Equal(t, exp.PipelineTaskRunID, act.PipelineTaskRunID)
	assert.Equal(t, exp.MinConfirmations, act.MinConfirmations)
	assert.Equal(t, exp.TransmitChecker, act.TransmitChecker)
	assert.Equal(t, exp.SignalCallback, act.SignalCallback)
	assert.Equal(t, exp.CallbackCompleted, act.CallbackCompleted)

	if len(exp.TxAttempts) == 0 {
		return
	}
	require.Equal(t, len(exp.TxAttempts), len(act.TxAttempts))
	for i := 0; i < len(exp.TxAttempts); i++ {
		assertTxAttemptEqual(t, exp.TxAttempts[i], act.TxAttempts[i])
	}
}

func assertTxAttemptEqual(t *testing.T, exp, act evmtxmgr.TxAttempt) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.TxID, act.TxID)
	assert.Equal(t, exp.TxFee, act.TxFee)
	assert.Equal(t, exp.ChainSpecificFeeLimit, act.ChainSpecificFeeLimit)
	assert.Equal(t, exp.SignedRawTx, act.SignedRawTx)
	assert.Equal(t, exp.Hash, act.Hash)
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.BroadcastBeforeBlockNum, act.BroadcastBeforeBlockNum)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.TxType, act.TxType)

	if len(exp.Receipts) == 0 {
		return
	}
	require.Equal(t, len(exp.Receipts), len(act.Receipts))
	for i := 0; i < len(exp.Receipts); i++ {
		assertChainReceiptEqual(t, exp.Receipts[i], act.Receipts[i])
	}
}

func assertChainReceiptEqual(t *testing.T, exp, act evmtxmgr.ChainReceipt) {
	assert.Equal(t, exp.GetStatus(), act.GetStatus())
	assert.Equal(t, exp.GetTxHash(), act.GetTxHash())
	assert.Equal(t, exp.GetBlockNumber(), act.GetBlockNumber())
	assert.Equal(t, exp.IsZero(), act.IsZero())
	assert.Equal(t, exp.IsUnmined(), act.IsUnmined())
	assert.Equal(t, exp.GetFeeUsed(), act.GetFeeUsed())
	assert.Equal(t, exp.GetTransactionIndex(), act.GetTransactionIndex())
	assert.Equal(t, exp.GetBlockHash(), act.GetBlockHash())
}
