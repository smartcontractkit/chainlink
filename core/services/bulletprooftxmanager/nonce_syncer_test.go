package bulletprooftxmanager_test

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_NonceSyncer_SyncAll(t *testing.T) {
	t.Parallel()

	t.Run("returns error if PendingNonceAt fails", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, from := cltest.MustAddRandomKeyToKeystore(t, store)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(0), errors.New("something exploded"))

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		err := ns.SyncAll(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")

		cltest.AssertCount(t, store, models.EthTx{}, 0)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 0)

		assertDatabaseNonce(t, store, from, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if chain nonce reflects local nonce", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, from := cltest.MustAddRandomKeyToKeystore(t, store)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(0), nil)

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))

		cltest.AssertCount(t, store, models.EthTx{}, 0)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 0)

		assertDatabaseNonce(t, store, from, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if chain nonce is behind local nonce", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, from := cltest.MustAddRandomKeyToKeystore(t, store, int64(32))

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			return from == addr
		})).Return(uint64(31), nil)

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))

		cltest.AssertCount(t, store, models.EthTx{}, 0)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 0)

		assertDatabaseNonce(t, store, from, 32)

		ethClient.AssertExpectations(t)
	})

	t.Run("fast forwards if chain nonce is ahead of local nonce and fills in recent transactions", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, acct1 := cltest.MustAddRandomKeyToKeystore(t, store, int64(0))
		_, acct2 := cltest.MustAddRandomKeyToKeystore(t, store, int64(32))

		accounts := store.KeyStore.Accounts()
		txes := makeRandomTransactions(t, store, 5, accounts, store.Config.ChainID())

		bPending := models.Block{
			Number:       0,
			Transactions: txes[4:],
		}

		b2 := models.Block{
			Number:       2,
			Hash:         cltest.NewHash(),
			ParentHash:   cltest.NewHash(),
			Transactions: txes[0:1],
		}
		b41 := models.Block{
			Number:       41,
			Hash:         cltest.NewHash(),
			ParentHash:   cltest.NewHash(),
			Transactions: txes[2:3],
		}
		bLatest := models.Block{
			Number:       42,
			Hash:         cltest.NewHash(),
			ParentHash:   b41.Hash,
			Transactions: txes[3:4],
		}

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// Nothing to do for acct2
			return acct2 == addr
		})).Return(uint64(32), nil)
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// acct1 has chain nonce of 5 which is ahead of local nonce 0
			return acct1 == addr
		})).Return(uint64(5), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "pending" && b[0].Args[1] == true &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "latest" && b[1].Args[1] == true
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &bPending
			elems[1].Result = &bLatest
		})
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 42 &&
				b[0].Args[0] == models.Int64ToHex(0) &&
				b[1].Args[0] == models.Int64ToHex(1) &&
				b[41].Args[0] == models.Int64ToHex(41)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[41].Result = &b41
			elems[2].Result = &b2
			elems[3].Error = errors.New("random error thrown in for fun")
		})

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))

		cltest.AssertCount(t, store, models.EthTx{}, 5)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 5)

		assertDatabaseNonce(t, store, acct1, 5)

		ethClient.AssertExpectations(t)
	})

	t.Run("only backfills to ETH_FINALITY_DEPTH", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		store.Config.Set("ETH_FINALITY_DEPTH", 2)

		_, acct1 := cltest.MustAddRandomKeyToKeystore(t, store, int64(0))

		accounts := store.KeyStore.Accounts()
		txes := makeRandomTransactions(t, store, 5, accounts, store.Config.ChainID())

		bPending := models.Block{
			Number:       0,
			Transactions: txes[4:],
		}

		b40 := models.Block{
			Number:     41,
			Hash:       cltest.NewHash(),
			ParentHash: cltest.NewHash(),
		}
		b41 := models.Block{
			Number:       41,
			Hash:         cltest.NewHash(),
			ParentHash:   b40.Hash,
			Transactions: txes[2:3],
		}
		bLatest := models.Block{
			Number:       42,
			Hash:         cltest.NewHash(),
			ParentHash:   b41.Hash,
			Transactions: txes[3:4],
		}

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(addr common.Address) bool {
			// acct1 has chain nonce of 5 which is ahead of local nonce 0
			return acct1 == addr
		})).Return(uint64(5), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "pending" && b[0].Args[1] == true &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "latest" && b[1].Args[1] == true
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &bPending
			elems[1].Result = &bLatest
		})
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == models.Int64ToHex(40) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == models.Int64ToHex(41)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b40
			elems[1].Result = &b41
		})

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		require.NoError(t, ns.SyncAll(context.Background()))

		cltest.AssertCount(t, store, models.EthTx{}, 3)
		cltest.AssertCount(t, store, models.EthTxAttempt{}, 3)

		assertDatabaseNonce(t, store, acct1, 5)

		ethClient.AssertExpectations(t)
	})
}

func Test_NonceSyncer_MakeInsert(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethClient := new(mocks.Client)

	_, from := cltest.MustAddRandomKeyToKeystore(t, store)
	acct := accounts.Account{Address: from}

	kst := store.KeyStore
	kst.Unlock(cltest.Password)

	var blockNum int64 = 42
	var nonce uint64 = 1
	to := cltest.NewAddress()
	amount := big.NewInt(4200)
	var gasLimit uint64 = 120000
	gasPrice := big.NewInt(25000000000)
	data := cltest.MustRandomBytes(t, 72)
	unsigned := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	t.Run("falls back to zero insert if encodeRLP would fail for some reason, e.g. tx is zero struct", func(t *testing.T) {
		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		tx := types.NewTx(&types.LegacyTx{})

		ins, err := ns.MakeInsert(*tx, acct, blockNum, int64(nonce))
		require.NoError(t, err)

		assertZero(t, ins, store, from, int64(nonce), blockNum)
	})

	t.Run("makes insert with the given payload", func(t *testing.T) {
		tx, err := kst.SignTx(acct, unsigned, store.Config.ChainID())
		require.NoError(t, err)

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		ins, err := ns.MakeInsert(*tx, acct, blockNum, int64(nonce))
		require.NoError(t, err)

		assert.Equal(t, int64(0), ins.Etx.ID)
		assert.Equal(t, int64(nonce), *ins.Etx.Nonce)
		assert.Equal(t, from, ins.Etx.FromAddress)
		assert.Equal(t, *tx.To(), ins.Etx.ToAddress)
		assert.Equal(t, tx.Data(), ins.Etx.EncodedPayload)
		assert.Equal(t, *amount, (big.Int)(ins.Etx.Value))
		assert.Equal(t, tx.Gas(), ins.Etx.GasLimit)
		assert.Nil(t, ins.Etx.Error)
		assert.Nil(t, ins.Etx.BroadcastAt)
		assert.Equal(t, models.EthTxUnconfirmed, ins.Etx.State)

		rlp := new(bytes.Buffer)
		require.NoError(t, tx.EncodeRLP(rlp))

		assert.Equal(t, int64(0), ins.Attempt.ID)
		assert.Equal(t, int64(0), ins.Attempt.EthTxID)
		assert.Equal(t, tx.GasPrice().String(), ins.Attempt.GasPrice.String())
		assert.Equal(t, rlp.Bytes(), ins.Attempt.SignedRawTx)
		assert.NotEqual(t, common.Hash{}, ins.Attempt.Hash)
		assert.NotNil(t, ins.Attempt.BroadcastBeforeBlockNum)
		assert.Equal(t, blockNum, *ins.Attempt.BroadcastBeforeBlockNum)
		assert.Equal(t, models.EthTxAttemptBroadcast, ins.Attempt.State)
	})
}

func Test_NonceSyncer_MakeZeroInsert(t *testing.T) {
	t.Parallel()

	var blockNum int64 = 42
	var nonce int64 = 1

	t.Run("errors if keystore.SignTx errors", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)
		kst := new(mocks.KeyStoreInterface)
		store.KeyStore = kst
		acct := accounts.Account{Address: cltest.NewAddress()}

		kst.On("SignTx", acct, mock.Anything, store.Config.ChainID()).Return(nil, errors.New("something exploded"))

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		_, err := ns.MakeZeroInsert(acct, blockNum, nonce)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})

	t.Run("returns insert with zero EthTx and Attempt", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethClient := new(mocks.Client)

		_, from := cltest.MustAddRandomKeyToKeystore(t, store)
		acct := accounts.Account{Address: from}

		oldKst := store.KeyStore
		oldKst.Unlock(cltest.Password)

		kst := new(mocks.KeyStoreInterface)
		store.KeyStore = kst

		kst.On("SignTx", acct, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == nonce && *tx.To() == from && big.NewInt(0).Cmp(tx.Value()) == 0
		}), store.Config.ChainID()).Return(
			func(acct accounts.Account, unsigned *types.Transaction, chainID *big.Int) *types.Transaction {
				signed, err := oldKst.SignTx(acct, unsigned, chainID)
				if err != nil {
					t.Fatal(err)
				}
				return signed
			},
			func(accounts.Account, *types.Transaction, *big.Int) error { return nil },
		)

		ns := bulletprooftxmanager.NewNonceSyncer(store, store.Config, ethClient)

		ins, err := ns.MakeZeroInsert(acct, blockNum, nonce)
		require.NoError(t, err)

		assertZero(t, ins, store, from, nonce, blockNum)
	})
}

func assertZero(t *testing.T, ins bulletprooftxmanager.NSinserttx, store *store.Store, from common.Address, nonce, blockNum int64) {
	t.Helper()

	assert.Equal(t, int64(0), ins.Etx.ID)
	assert.Equal(t, nonce, *ins.Etx.Nonce)
	assert.Equal(t, from, ins.Etx.FromAddress)
	assert.Equal(t, from, ins.Etx.ToAddress)
	assert.Equal(t, []byte{}, ins.Etx.EncodedPayload)
	assert.Equal(t, assets.NewEthValue(0), ins.Etx.Value)
	assert.Equal(t, store.Config.EthGasLimitDefault(), ins.Etx.GasLimit)
	assert.Nil(t, ins.Etx.Error)
	assert.Nil(t, ins.Etx.BroadcastAt)
	assert.Equal(t, models.EthTxUnconfirmed, ins.Etx.State)

	assert.Equal(t, int64(0), ins.Attempt.ID)
	assert.Equal(t, int64(0), ins.Attempt.EthTxID)
	assert.Equal(t, store.Config.EthGasPriceDefault().String(), ins.Attempt.GasPrice.String())
	assert.Len(t, ins.Attempt.SignedRawTx, 103)
	assert.NotEqual(t, common.Hash{}, ins.Attempt.Hash)
	assert.NotNil(t, ins.Attempt.BroadcastBeforeBlockNum)
	assert.Equal(t, blockNum, *ins.Attempt.BroadcastBeforeBlockNum)
	assert.Equal(t, models.EthTxAttemptBroadcast, ins.Attempt.State)
}

func assertDatabaseNonce(t *testing.T, store *store.Store, from common.Address, nonce int64) {
	t.Helper()

	k, err := store.KeyByAddress(from)
	require.NoError(t, err)
	assert.Equal(t, nonce, k.NextNonce)
}

func makeRandomTransactions(t *testing.T, store *store.Store, n int, accounts []accounts.Account, chainID *big.Int) (txes []types.Transaction) {
	for i := 0; i < n; i++ {
		unsigned := types.NewTransaction(uint64(i), cltest.NewAddress(), big.NewInt(int64(100+i)), uint64(100000+i), big.NewInt(int64(1000000000+i)), cltest.MustRandomBytes(t, 100+i))

		// rotate accounts
		acct := accounts[i%len(accounts)]
		signed, err := store.KeyStore.SignTx(acct, unsigned, chainID)
		require.NoError(t, err)
		txes = append(txes, *signed)
	}
	return
}
