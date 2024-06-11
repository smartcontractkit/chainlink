package txmgr_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

var (
	tenGwei = assets.NewWeiI(10_000_000_000)
	oneGwei = assets.NewWeiI(1_000_000_000)
)

func TestStuckTxDetector_Disabled(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	lggr := logger.Test(t)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	autoPurgeCfg := testAutoPurgeConfig{
		enabled: false,
	}
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), autoPurgeCfg, feeEstimator, txStore, ethClient)

	t.Run("returns empty list if auto-purge feature is disabled", func(t *testing.T) {
		txs, err := stuckTxDetector.DetectStuckTransactions(tests.Context(t), []common.Address{fromAddress}, 100)
		require.NoError(t, err)
		require.Len(t, txs, 0)
	})
}

func TestStuckTxDetector_LoadPurgeBlockNumMap(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)
	blockNum := int64(100)

	lggr := logger.Test(t)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	marketGasPrice := assets.GWei(15)
	fee := gas.EvmFee{Legacy: marketGasPrice}
	feeEstimator.On("GetFee", mock.Anything, []byte{}, uint64(0), mock.Anything).Return(fee, uint64(0), nil)
	autoPurgeThreshold := uint32(5)
	autoPurgeMinAttempts := uint32(3)
	autoPurgeCfg := testAutoPurgeConfig{
		enabled:     true, // Enable auto-purge feature for testing
		threshold:   autoPurgeThreshold,
		minAttempts: autoPurgeMinAttempts,
	}
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), autoPurgeCfg, feeEstimator, txStore, ethClient)

	t.Run("purge num map loaded on startup rate limits new purges on startup", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		mustInsertFatalErrorTxWithError(t, txStore, 0, fromAddress, blockNum)

		err := stuckTxDetector.LoadPurgeBlockNumMap(ctx, []common.Address{fromAddress})
		require.NoError(t, err)

		enabledAddresses := []common.Address{fromAddress}
		// Create attempts broadcasted autoPurgeThreshold block ago to ensure broadcast block num check is not being triggered
		// Create autoPurgeMinAttempts number of attempts to ensure the broadcast attempt count check is not being triggered
		// Create attempts so that the latest has a higher gas price than the market to ensure the gas price check is not being triggered
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 1, fromAddress, autoPurgeMinAttempts, blockNum-int64(autoPurgeThreshold), marketGasPrice.Add(oneGwei))

		// Run detection logic on autoPurgeThreshold blocks past the latest broadcast attempt
		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, enabledAddresses, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 0)
	})
}

func TestStuckTxDetector_FindPotentialStuckTxs(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, config := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	lggr := logger.Test(t)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), config.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)

	t.Run("returns empty list if no unconfimed transactions found", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		stuckTxs, err := stuckTxDetector.FindUnconfirmedTxWithLowestNonce(ctx, []common.Address{fromAddress})
		require.NoError(t, err)
		require.Len(t, stuckTxs, 0)
	})

	t.Run("returns 1 unconfirmed transaction for each unique from address", func(t *testing.T) {
		_, fromAddress1 := cltest.MustInsertRandomKey(t, ethKeyStore)
		_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
		// Insert 2 txs for from address, should only return the lowest nonce txs
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress1)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress1)
		// Insert 1 tx for other from address, should return a tx
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress2)
		stuckTxs, err := stuckTxDetector.FindUnconfirmedTxWithLowestNonce(ctx, []common.Address{fromAddress1, fromAddress2})
		require.NoError(t, err)

		require.Len(t, stuckTxs, 2)
		var foundFromAddresses []common.Address
		for _, stuckTx := range stuckTxs {
			// Make sure lowest nonce tx is returned for both from addresses
			require.Equal(t, types.Nonce(0), *stuckTx.Sequence)
			// Make sure attempts are loaded into the tx
			require.Len(t, stuckTx.TxAttempts, 1)
			foundFromAddresses = append(foundFromAddresses, stuckTx.FromAddress)
		}
		require.Contains(t, foundFromAddresses, fromAddress1)
		require.Contains(t, foundFromAddresses, fromAddress2)
	})

	t.Run("excludes transactions already marked for purge", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		mustInsertUnconfirmedEthTxWithBroadcastPurgeAttempt(t, txStore, 0, fromAddress)
		stuckTxs, err := stuckTxDetector.FindUnconfirmedTxWithLowestNonce(ctx, []common.Address{fromAddress})
		require.NoError(t, err)
		require.Len(t, stuckTxs, 0)
	})
}

func TestStuckTxDetector_DetectStuckTransactionsHeuristic(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	lggr := logger.Test(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	// Return 10 gwei as market gas price
	marketGasPrice := tenGwei
	fee := gas.EvmFee{Legacy: marketGasPrice}
	feeEstimator.On("GetFee", mock.Anything, []byte{}, uint64(0), mock.Anything).Return(fee, uint64(0), nil)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	autoPurgeThreshold := uint32(5)
	autoPurgeMinAttempts := uint32(3)
	autoPurgeCfg := testAutoPurgeConfig{
		enabled:     true, // Enable auto-purge feature for testing
		threshold:   autoPurgeThreshold,
		minAttempts: autoPurgeMinAttempts,
	}
	blockNum := int64(100)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), autoPurgeCfg, feeEstimator, txStore, ethClient)

	t.Run("not stuck, Threshold amount of blocks have not passed since broadcast", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		enabledAddresses := []common.Address{fromAddress}
		// Create attempts broadcasted at the current broadcast number to test the block num threshold check
		// Create autoPurgeMinAttempts number of attempts to ensure the broadcast attempt count check is not being triggered
		// Create attempts so that the latest has a higher gas price than the market to ensure the gas price check is not being triggered
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, autoPurgeMinAttempts, blockNum, marketGasPrice.Add(oneGwei))

		// Run detection logic on the same block number as the latest broadcast attempt to stay within the autoPurgeThreshold
		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, enabledAddresses, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 0)
	})

	t.Run("not stuck, Threshold amount of blocks have not passed since last purge", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		enabledAddresses := []common.Address{fromAddress}
		// Create attempts broadcasted autoPurgeThreshold block ago to ensure broadcast block num check is not being triggered
		// Create autoPurgeMinAttempts number of attempts to ensure the broadcast attempt count check is not being triggered
		// Create attempts so that the latest has a higher gas price than the market to ensure the gas price check is not being triggered
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, autoPurgeMinAttempts, blockNum-int64(autoPurgeThreshold), marketGasPrice.Add(oneGwei))

		// Set the last purge block num as the current block num to test rate limiting condition
		stuckTxDetector.SetPurgeBlockNum(fromAddress, blockNum)

		// Run detection logic on autoPurgeThreshold blocks past the latest broadcast attempt
		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, enabledAddresses, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 0)
	})

	t.Run("not stuck, MinAttempts amount of attempts have not been broadcasted", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		enabledAddresses := []common.Address{fromAddress}
		// Create attempts broadcasted autoPurgeThreshold block ago to ensure broadcast block num check is not being triggered
		// Create fewer attempts than autoPurgeMinAttempts to test min attempt check
		// Create attempts so that the latest has a higher gas price than the market to ensure the gas price check is not being triggered
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, autoPurgeMinAttempts-1, blockNum-int64(autoPurgeThreshold), marketGasPrice.Add(oneGwei))

		// Run detection logic on autoPurgeThreshold blocks past the latest broadcast attempt
		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, enabledAddresses, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 0)
	})

	t.Run("not stuck, transaction gas price is lower than market gas price", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		enabledAddresses := []common.Address{fromAddress}
		// Create attempts broadcasted autoPurgeThreshold block ago to ensure broadcast block num check is not being triggered
		// Create autoPurgeMinAttempts number of attempts to ensure the broadcast attempt count check is not being triggered
		// Create attempts so that the latest has a lower gas price than the market to test the gas price check
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, autoPurgeMinAttempts, blockNum-int64(autoPurgeThreshold), marketGasPrice.Sub(oneGwei))

		// Run detection logic on autoPurgeThreshold blocks past the latest broadcast attempt
		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, enabledAddresses, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 0)
	})

	t.Run("detects stuck transaction", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		enabledAddresses := []common.Address{fromAddress}
		// Create attempts so that the oldest broadcast attempt's block num is what meets the threshold check
		// Create autoPurgeMinAttempts number of attempts to ensure the broadcast attempt count check is not being triggered
		// Create attempts broadcasted autoPurgeThreshold block ago to ensure broadcast block num check is not being triggered
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, autoPurgeMinAttempts, blockNum-int64(autoPurgeThreshold)+int64(autoPurgeMinAttempts-1), marketGasPrice.Add(oneGwei))

		// Run detection logic on autoPurgeThreshold blocks past the latest broadcast attempt
		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, enabledAddresses, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 1)
	})
}

func TestStuckTxDetector_DetectStuckTransactionsZkEVM(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	lggr := logger.Test(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	autoPurgeCfg := testAutoPurgeConfig{
		enabled: true,
	}
	blockNum := int64(100)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, chaintype.ChainZkEvm, assets.NewWei(assets.NewEth(100).ToInt()), autoPurgeCfg, feeEstimator, txStore, ethClient)
	t.Run("returns empty list if no stuck transactions identified", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		tx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, blockNum, tenGwei)
		attempts := tx.TxAttempts[0]
		// Request still returns transaction by hash, transaction not discarded by network and not considered stuck
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempts.Hash, "eth_getTransactionByHash")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			resp, err := json.Marshal(types.Transaction{})
			require.NoError(t, err)
			elems[0].Error = json.Unmarshal(resp, elems[0].Result)
		}).Once()

		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, []common.Address{fromAddress}, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 0)
	})

	t.Run("returns stuck transactions discarded by chain", func(t *testing.T) {
		// Insert tx that will be mocked as stuck
		_, fromAddress1 := cltest.MustInsertRandomKey(t, ethKeyStore)
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress1, 1, blockNum, tenGwei)

		// Insert tx that will still be valid
		_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
		mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress2, 1, blockNum, tenGwei)

		// Return nil response for a tx and a normal response for the other
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = nil // Return nil to signal discarded tx
			resp, err := json.Marshal(types.Transaction{})
			require.NoError(t, err)
			elems[1].Error = json.Unmarshal(resp, elems[1].Result) // Return non-nil result to signal a valid tx
		}).Once()

		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, []common.Address{fromAddress1, fromAddress2}, blockNum)
		require.NoError(t, err)
		// Expect only 1 tx to return as stuck due to nil eth_getTransactionByHash response
		require.Len(t, txs, 1)
	})
}

func TestStuckTxDetector_DetectStuckTransactionsScroll(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	lggr := logger.Test(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	blockNum := int64(100)

	t.Run("returns stuck tx identified using the custom scroll API", func(t *testing.T) {
		// Insert tx that will be mocked as stuck
		_, fromAddress1 := cltest.MustInsertRandomKey(t, ethKeyStore)
		tx1 := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress1, 1, blockNum, tenGwei)
		attempts1 := tx1.TxAttempts[0]

		// Insert tx that will still be valid
		_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
		tx2 := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress2, 1, blockNum, tenGwei)
		attempts2 := tx2.TxAttempts[0]

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			_, err := res.Write([]byte(fmt.Sprintf(`{"errcode": 0,"errmsg": "","data": {"%s": 1, "%s": 0}}`, attempts1.Hash, attempts2.Hash)))
			require.NoError(t, err)
		}))
		defer func() { testServer.Close() }()
		testUrl, err := url.Parse(testServer.URL)
		require.NoError(t, err)

		autoPurgeCfg := testAutoPurgeConfig{
			enabled:         true,
			detectionApiUrl: testUrl,
		}
		stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, chaintype.ChainScroll, assets.NewWei(assets.NewEth(100).ToInt()), autoPurgeCfg, feeEstimator, txStore, ethClient)

		txs, err := stuckTxDetector.DetectStuckTransactions(ctx, []common.Address{fromAddress1, fromAddress2}, blockNum)
		require.NoError(t, err)
		require.Len(t, txs, 1)
		require.Equal(t, tx1.ID, txs[0].ID)
	})
}

func mustInsertUnconfirmedTxWithBroadcastAttempts(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, numAttempts uint32, latestBroadcastBlockNum int64, latestGasPrice *assets.Wei) txmgr.Tx {
	ctx := tests.Context(t)
	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, nonce, fromAddress)
	// Insert attempts from oldest to newest
	for i := int64(numAttempts - 1); i >= 0; i-- {
		blockNum := latestBroadcastBlockNum - i
		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

		attempt.State = txmgrtypes.TxAttemptBroadcast
		attempt.BroadcastBeforeBlockNum = &blockNum
		attempt.TxFee = gas.EvmFee{Legacy: latestGasPrice.Sub(assets.NewWeiI(i))}
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	}
	etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func mustInsertFatalErrorTxWithError(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, blockNum int64) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	etx.State = txmgrcommon.TxFatalError
	etx.Error = null.StringFrom("fatal error")
	broadcastAt := time.Now()
	etx.BroadcastAt = &broadcastAt
	etx.InitialBroadcastAt = &broadcastAt
	n := types.Nonce(nonce)
	etx.Sequence = &n
	etx.ChainID = testutils.FixtureChainID
	require.NoError(t, txStore.InsertTx(tests.Context(t), &etx))

	attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
	ctx := tests.Context(t)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.IsPurgeAttempt = true
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

	receipt := newTxReceipt(attempt.Hash, int(blockNum), 0)
	_, err := txStore.InsertReceipt(ctx, &receipt)
	require.NoError(t, err)

	etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func mustInsertUnconfirmedEthTxWithBroadcastPurgeAttempt(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address) txmgr.Tx {
	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, nonce, fromAddress)
	attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
	ctx := tests.Context(t)

	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.IsPurgeAttempt = true
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

type testAutoPurgeConfig struct {
	enabled         bool
	threshold       uint32
	minAttempts     uint32
	detectionApiUrl *url.URL
}

func (t testAutoPurgeConfig) Enabled() bool             { return t.enabled }
func (t testAutoPurgeConfig) Threshold() uint32         { return t.threshold }
func (t testAutoPurgeConfig) MinAttempts() uint32       { return t.minAttempts }
func (t testAutoPurgeConfig) DetectionApiUrl() *url.URL { return t.detectionApiUrl }
