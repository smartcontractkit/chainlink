package chains_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
)

type TestingTxStore[
	ADDR types.Hashable,
	CHAIN_ID types.ID,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH], chainID CHAIN_ID) (tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, gas.EvmFee], err error)
}

type txStoreFunc func(t *testing.T) (TestingTxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee], common.Address)

func evmTxStore(t *testing.T) (TestingTxStore[common.Address, *big.Int, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee], common.Address) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

	return cltest.NewTxStore(t, db, cfg.Database()), fromAddress
}

var txStoresFuncs = []txStoreFunc{
	evmTxStore,
	/*
		ims, err := txmgr.NewInMemoryStore[
			*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee,
		](chainID, keyStore.Eth(), txStore)
	*/
}

func TestTxStore_CreateTransaction(t *testing.T) {
	for _, f := range txStoresFuncs {
		txStore, fromAddress := f(t)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		subject := uuid.New()
		strategy := commontxmmocks.NewTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("*txmgr.evmTxStore")).Return(int64(0), nil)
		ctx := context.Background()
		idempotencyKey := "11"

		tts := []struct {
			scenario                     string
			createTransactionInput       createTransactionInput
			createTransactionOutputCheck func(*testing.T, txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], error)
		}{
			{
				scenario: "success",
				createTransactionInput: createTransactionInput{
					txRequest: txmgrtypes.TxRequest[common.Address, common.Hash]{
						IdempotencyKey: &idempotencyKey,
						FromAddress:    fromAddress,
						ToAddress:      common.BytesToAddress([]byte("test")),
						EncodedPayload: []byte{1, 2, 3},
						FeeLimit:       uint32(1000),
						Meta:           nil,
						Strategy:       strategy,
					},
					chainID: ethClient.ConfiguredChainID(),
				},
				createTransactionOutputCheck: func(t *testing.T, tx txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
					funcName := "CreateTransaction"
					require.NoError(t, err, fmt.Sprintf("%s: expected err to be nil", funcName))
					assert.Equal(t, &idempotencyKey, tx.IdempotencyKey, fmt.Sprintf("%s: expected idempotencyKey to match actual idempotencyKey", funcName))
					// Check CreatedAt is within 1 second of now
					assert.WithinDuration(t, time.Now().UTC(), tx.CreatedAt, time.Second, fmt.Sprintf("%s: expected time to be within 1 second of actual time", funcName))
					assert.Equal(t, txmgr.TxUnstarted, tx.State, fmt.Sprintf("%s: expected state to match actual state", funcName))
					assert.Equal(t, ethClient.ConfiguredChainID(), tx.ChainID, fmt.Sprintf("%s: expected chainID to match actual chainID", funcName))
					assert.Equal(t, fromAddress, tx.FromAddress, fmt.Sprintf("%s: expected fromAddress to match actual fromAddress", funcName))
					assert.Equal(t, common.BytesToAddress([]byte("test")), tx.ToAddress, fmt.Sprintf("%s: expected toAddress to match actual toAddress", funcName))
					assert.Equal(t, []byte{1, 2, 3}, tx.EncodedPayload, fmt.Sprintf("%s: expected encodedPayload to match actual encodedPayload", funcName))
					assert.Equal(t, uint32(1000), tx.FeeLimit, fmt.Sprintf("%s: expected feeLimit to match actual feeLimit", funcName))
					var expMeta *datatypes.JSON
					assert.Equal(t, expMeta, tx.Meta, fmt.Sprintf("%s: expected meta to match actual meta", funcName))
					assert.Equal(t, uuid.NullUUID{UUID: subject, Valid: true}, tx.Subject, fmt.Sprintf("%s: expected subject to match actual subject", funcName))
				},
			},
		}

		for _, tt := range tts {
			t.Run(tt.scenario, func(t *testing.T) {
				actTx, actErr := txStore.CreateTransaction(ctx, tt.createTransactionInput.txRequest, tt.createTransactionInput.chainID)
				tt.createTransactionOutputCheck(t, actTx, actErr)

				// TODO(jtw): Check that the transaction was persisted
			})
		}
	}
}

/*
func TestTxStore_FindTxWithIdempotencyKey(t *testing.T) {
	txStore := evmtxmgr.NewTxStore(nil, nil, nil)
	ctx := context.Background()

	tts := []struct {
		scenario                       string
		findTxWithIdempotencyKeyInput  findTxWithIdempotencyKeyInput
		findTxWithIdempotencyKeyOutput func(*testing.T, txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], error)
	}{
		{
			findTxWithIdempotencyKeyInput: findTxWithIdempotencyKeyInput{
				idempotencyKey: "11",
				chainID:        chainID,
			},
			findTxWithIdempotencyKeyOutput: func(t *testing.T, tx txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
				funcName := "FindTxWithIdempotencyKey"
				require.NoError(t, err, fmt.Sprintf("%s: expected err to be nil", funcName))
				assert.Equal(t, &idempotencyKey, tx.IdempotencyKey, fmt.Sprintf("%s: expected idempotencyKey to match actual idempotencyKey", funcName))
				// Check CreatedAt is within 1 second of now
				assert.WithinDuration(t, time.Now().UTC(), tx.CreatedAt, time.Second, fmt.Sprintf("%s: expected time to be within 1 second of actual time", funcName))
				assert.Equal(t, txmgr.TxUnstarted, tx.State, fmt.Sprintf("%s: expected state to match actual state", funcName))
				assert.Equal(t, chainID, tx.ChainID, fmt.Sprintf("%s: expected chainID to match actual chainID", funcName))
				assert.Equal(t, fromAddress, tx.FromAddress, fmt.Sprintf("%s: expected fromAddress to match actual fromAddress", funcName))
				assert.Equal(t, common.BytesToAddress([]byte("test")), tx.ToAddress, fmt.Sprintf("%s: expected toAddress to match actual toAddress", funcName))
				assert.Equal(t, []byte{1, 2, 3}, tx.EncodedPayload, fmt.Sprintf("%s: expected encodedPayload to match actual encodedPayload", funcName))
				assert.Equal(t, uint32(1000), tx.FeeLimit, fmt.Sprintf("%s: expected feeLimit to match actual feeLimit", funcName))
				var expMeta *datatypes.JSON
				assert.Equal(t, expMeta, tx.Meta, fmt.Sprintf("%s: expected meta to match actual meta", funcName))
				assert.Equal(t, uuid.NullUUID{UUID: subject, Valid: true}, tx.Subject, fmt.Sprintf("%s: expected subject to match actual subject", funcName))
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.scenario, func(t *testing.T) {
			actTxPtr, actErr := txStore.FindTxWithIdempotencyKey(ctx, tt.findTxWithIdempotencyKeyInput.idempotencyKey, tt.findTxWithIdempotencyKeyInput.chainID)
			tt.findTxWithIdempotencyKeyOutput(t, *actTxPtr, actErr)
		})
	}
}

func TestTxStore_CheckTxQueueCapacity(t *testing.T) {
	txStore := evmtxmgr.NewTxStore(nil, nil, nil)
	ctx := context.Background()

	tts := []struct {
		scenario                  string
		checkTxQueueCapacityInput checkTxQueueCapacityInput
		expErr                    error
	}{
		{
			checkTxQueueCapacityInput: checkTxQueueCapacityInput{
				fromAddress: fromAddress,
				maxQueued:   uint64(16),
				chainID:     chainID,
			},
			expErr: nil,
		},
	}

	for _, tt := range tts {
		t.Run(tt.scenario, func(t *testing.T) {
			actErr := txStore.CheckTxQueueCapacity(ctx, tt.checkTxQueueCapacityInput.fromAddress, tt.checkTxQueueCapacityInput.maxQueued, tt.checkTxQueueCapacityInput.chainID)
			require.Equal(t, tt.expErr, actErr, "CheckTxQueueCapacity: expected err to match actual err")
		})
	}
}

func TestTxStore_FindLatestSequence(t *testing.T) {
	txStore := evmtxmgr.NewTxStore(nil, nil, nil)
	ctx := context.Background()

	tts := []struct {
		scenario                 string
		findLatestSequenceInput  findLatestSequenceInput
		findLatestSequenceOutput func(*testing.T, evmtypes.Nonce, error)
	}{
		{
			findLatestSequenceInput: findLatestSequenceInput{
				fromAddress: fromAddress,
				chainID:     chainID,
			},
			findLatestSequenceOutput: func(t *testing.T, seq evmtypes.Nonce, err error) {
				funcName := "FindLatestSequence"
				require.NoError(t, err, fmt.Sprintf("%s: expected err to be nil", funcName))
				assert.Equal(t, uint64(0), seq, fmt.Sprintf("%s: expected seq to match actual seq", funcName))
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.scenario, func(t *testing.T) {
			actSeq, actErr := txStore.FindLatestSequence(ctx, tt.findLatestSequenceInput.fromAddress, tt.findLatestSequenceInput.chainID)
			tt.findLatestSequenceOutput(t, actSeq, actErr)
		})
	}
}
*/

type createTransactionInput struct {
	txRequest txmgrtypes.TxRequest[common.Address, common.Hash]
	chainID   *big.Int
}
type findTxWithIdempotencyKeyInput struct {
	idempotencyKey string
	chainID        *big.Int
}
type checkTxQueueCapacityInput struct {
	fromAddress common.Address
	maxQueued   uint64
	chainID     *big.Int
}
type checkTxQueueCapacityOutput struct {
	err error
}
type findLatestSequenceInput struct {
	fromAddress common.Address
	chainID     *big.Int
}
