package txmgr_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	types "github.com/smartcontractkit/chainlink/v2/common/types"
)

func TestInMemoryStore_CreateTransaction(t *testing.T) {
	chainID := big.NewInt(1)
	idempotencyKey := "11"
	fromAddress := common.BytesToAddress(rand.Bytes(20))
	mockKeyStore := mocks.NewKeyStore[common.Address, *big.Int, types.Sequence](t)
	mockKeyStore.Mock.On("EnabledAddressesForChain", chainID).Return([]common.Address{fromAddress}, nil)
	mockEventRecorder := mocks.NewTxStore[
		common.Address, *big.Int, common.Hash, common.Hash, txmgrtypes.ChainReceipt[common.Hash, common.Hash], types.Sequence, feetypes.Fee](t)
	mockEventRecorder.Mock.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil, nil)
	ctx := context.Background()

	ims, err := txmgr.NewInMemoryStore[
		*big.Int, common.Address, common.Hash, common.Hash, txmgrtypes.ChainReceipt[common.Hash, common.Hash], types.Sequence, feetypes.Fee,
	](chainID, mockKeyStore, mockEventRecorder)
	ims.LegacyEnabled = false // TODO(jtw): this is just for initial testing, remove this
	require.NoError(t, err)

	tts := []struct {
		scenario                       string
		createTransactionInput         createTransactionInput
		createTransactionOutput        createTransactionOutput
		findTxWithIdempotencyKeyInput  findTxWithIdempotencyKeyInput
		findTxWithIdempotencyKeyOutput findTxWithIdempotencyKeyOutput
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
					Strategy:       nil, //TODO
				},
				chainID: chainID,
			},
			createTransactionOutput: createTransactionOutput{
				tx: txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, types.Sequence, feetypes.Fee]{
					IdempotencyKey: &idempotencyKey,
					CreatedAt:      time.Now().UTC(),
					State:          txmgr.TxUnstarted,
					ChainID:        chainID,
					FromAddress:    fromAddress,
					ToAddress:      common.BytesToAddress([]byte("test")),
					EncodedPayload: []byte{1, 2, 3},
					FeeLimit:       uint32(1000),
					Meta:           nil,
				},
				err: nil,
			},
			findTxWithIdempotencyKeyInput: findTxWithIdempotencyKeyInput{
				idempotencyKey: "11",
				chainID:        chainID,
			},
			findTxWithIdempotencyKeyOutput: findTxWithIdempotencyKeyOutput{
				tx: txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, types.Sequence, feetypes.Fee]{
					IdempotencyKey: &idempotencyKey,
					CreatedAt:      time.Now().UTC(),
					State:          txmgr.TxUnstarted,
					ChainID:        chainID,
					FromAddress:    fromAddress,
					ToAddress:      common.BytesToAddress([]byte("test")),
					EncodedPayload: []byte{1, 2, 3},
					FeeLimit:       uint32(1000),
					Meta:           nil,
				},
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
			require.Equal(t, tt.createTransactionOutput.err, actErr, "CreateTransaction: expected err to match actual err")
			// Check CreatedAt is within 1 second of now
			assert.WithinDuration(t, tt.createTransactionOutput.tx.CreatedAt, actTx.CreatedAt, time.Second, "CreateTransaction: expected time to be within 1 second of actual time")
			// Reset CreatedAt to avoid flaky test
			tt.createTransactionOutput.tx.CreatedAt = actTx.CreatedAt
			assert.Equal(t, tt.createTransactionOutput.tx, actTx, "CreateTransaction: expected tx to match actual tx")

			actTxPtr, actErr := ims.FindTxWithIdempotencyKey(ctx, tt.findTxWithIdempotencyKeyInput.idempotencyKey, tt.findTxWithIdempotencyKeyInput.chainID)
			require.Equal(t, tt.findTxWithIdempotencyKeyOutput.err, actErr, "FindTxWithIdempotencyKey: expected err to match actual err")
			// Check CreatedAt is within 1 second of now
			assert.WithinDuration(t, tt.findTxWithIdempotencyKeyOutput.tx.CreatedAt, actTx.CreatedAt, time.Second, "FindTxWithIdempotencyKey: expected time to be within 1 second of actual time")
			// Reset CreatedAt to avoid flaky test
			tt.findTxWithIdempotencyKeyOutput.tx.CreatedAt = actTxPtr.CreatedAt
			assert.Equal(t, tt.findTxWithIdempotencyKeyOutput.tx, actTx, "FindTxWithIdempotencyKey: expected tx to match actual tx")

			actErr = ims.CheckTxQueueCapacity(ctx, tt.checkTxQueueCapacityInput.fromAddress, tt.checkTxQueueCapacityInput.maxQueued, tt.checkTxQueueCapacityInput.chainID)
			require.Equal(t, tt.checkTxQueueCapacityOutput.err, actErr, "CheckTxQueueCapacity: expected err to match actual err")
		})
	}

}

type createTransactionInput struct {
	txRequest txmgrtypes.TxRequest[common.Address, common.Hash]
	chainID   *big.Int
}
type createTransactionOutput struct {
	tx  txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, types.Sequence, feetypes.Fee]
	err error
}
type findTxWithIdempotencyKeyInput struct {
	idempotencyKey string
	chainID        *big.Int
}
type findTxWithIdempotencyKeyOutput struct {
	tx  txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, types.Sequence, feetypes.Fee]
	err error
}
type checkTxQueueCapacityInput struct {
	fromAddress common.Address
	maxQueued   uint64
	chainID     *big.Int
}
type checkTxQueueCapacityOutput struct {
	err error
}
