package ocrcommon_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/statuschecker"
)

func newMockTxStrategy(t *testing.T) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

func Test_DefaultTransmitter_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)

	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress, fromAddress2},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)

	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		FromAddress:      fromAddress2,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             nil,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_Round_Robin_Error(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	fromAddress := common.Address{}

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
	)
	require.NoError(t, err)
	require.Error(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_No_Keystore_Error(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	_, err := ocrcommon.NewTransmitter(
		txm,
		[]common.Address{fromAddress, fromAddress2},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		nil,
	)
	require.Error(t, err)
}

func Test_Transmitter_With_StatusChecker_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)
	statusChecker := mocks.NewTransactionStatusChecker(t)
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	idempotencyKey := "1-0"
	txMeta := &txmgr.TxMeta{MessageIDs: []string{"1"}}

	transmitter, err := ocrcommon.NewTransmitterWithStatusChecker(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		chainID,
		ethKeyStore,
		statusChecker,
	)
	require.NoError(t, err)

	statusChecker.On("CheckMessageStatus", mock.Anything, "1").Return([]statuschecker.TransactionStatus{}, -1, nil).Once()

	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		IdempotencyKey:   &idempotencyKey,
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         gasLimit,
		ForwarderAddress: common.Address{},
		Meta:             txMeta,
		Strategy:         strategy,
	}).Return(txmgr.Tx{}, nil).Once()

	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, txMeta))
	// check that the status checker was called
	statusChecker.AssertExpectations(t)
}
