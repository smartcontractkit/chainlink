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

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
)

func TestInMemoryStore_CreateTransaction(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	idempotencyKey := "11"
	lggr := logger.TestLogger(t)
	txStore := evmtxmgr.NewTxStore(db, lggr, cfg.Database())
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	chainID := ethClient.ConfiguredChainID()

	subject := uuid.New()
	strategy := commontxmmocks.NewTxStrategy(t)
	strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
	strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("*txmgr.evmTxStore")).Return(int64(0), nil)
	ctx := context.Background()

	ims, err := txmgr.NewInMemoryStore[
		*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee,
	](chainID, keyStore.Eth(), txStore)
	require.NoError(t, err)

	tts := []struct {
		scenario                       string
		createTransactionInput         createTransactionInput
		createTransactionOutputCheck   func(*testing.T, txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], error)
		findTxWithIdempotencyKeyInput  findTxWithIdempotencyKeyInput
		findTxWithIdempotencyKeyOutput func(*testing.T, txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], error)
		checkTxQueueCapacityInput      checkTxQueueCapacityInput
		checkTxQueueCapacityOutput     checkTxQueueCapacityOutput
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
				chainID: chainID,
			},
			createTransactionOutputCheck: func(t *testing.T, tx txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
				funcName := "CreateTransaction"
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
			checkTxQueueCapacityInput: checkTxQueueCapacityInput{
				fromAddress: fromAddress,
				maxQueued:   uint64(16),
				chainID:     chainID,
			},
			checkTxQueueCapacityOutput: checkTxQueueCapacityOutput{
				err: nil,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.scenario, func(t *testing.T) {
			actTx, actErr := ims.CreateTransaction(ctx, tt.createTransactionInput.txRequest, tt.createTransactionInput.chainID)
			tt.createTransactionOutputCheck(t, actTx, actErr)

			actTxPtr, actErr := ims.FindTxWithIdempotencyKey(ctx, tt.findTxWithIdempotencyKeyInput.idempotencyKey, tt.findTxWithIdempotencyKeyInput.chainID)
			tt.findTxWithIdempotencyKeyOutput(t, *actTxPtr, actErr)

			actErr = ims.CheckTxQueueCapacity(ctx, tt.checkTxQueueCapacityInput.fromAddress, tt.checkTxQueueCapacityInput.maxQueued, tt.checkTxQueueCapacityInput.chainID)
			require.Equal(t, tt.checkTxQueueCapacityOutput.err, actErr, "CheckTxQueueCapacity: expected err to match actual err")
		})
	}

}

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
