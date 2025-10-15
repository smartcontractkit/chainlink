package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"

	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/transmitter/ocr"
)

func newMockTxStrategy(t *testing.T) *commontxmmocks.TxStrategy {
	return commontxmmocks.NewTxStrategy(t)
}

func Test_DefaultTransmitter_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	ethKeyStore := keystest.Addresses{fromAddress}

	gasLimit := uint64(1000)
	effectiveTransmitterAddress := fromAddress
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocr.NewTransmitter(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
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

	memKeys := keystest.NewMemoryChainStore()
	fromAddress := memKeys.MustCreate(t)
	fromAddress2 := memKeys.MustCreate(t)
	ethKeyStore := keys.NewStore(memKeys)

	gasLimit := uint64(1000)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocr.NewTransmitter(
		txm,
		[]common.Address{fromAddress, fromAddress2},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
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

	fromAddress := common.Address{}

	gasLimit := uint64(1000)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	transmitter, err := ocr.NewTransmitter(
		txm,
		[]common.Address{fromAddress},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		keystest.Addresses{},
	)
	require.NoError(t, err)
	require.Error(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_No_Keystore_Error(t *testing.T) {
	t.Parallel()

	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	fromAddress2 := memKS.MustCreate(t)

	gasLimit := uint64(1000)
	effectiveTransmitterAddress := common.Address{}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)

	_, err := ocr.NewTransmitter(
		txm,
		[]common.Address{fromAddress, fromAddress2},
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		nil,
	)
	require.Error(t, err)
}

func Test_Transmitter_With_StatusChecker_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	ethKeyStore := keys.NewChainStore(memKS, testutils.FixtureChainID)

	gasLimit := uint64(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	idempotencyKey := "1-0"
	txMeta := &txmgr.TxMeta{MessageIDs: []string{"1"}}

	transmitter, err := NewTransmitterWithStatusChecker(
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

	// This case is for when the message ID was not found in the status checker
	txm.On("GetTransactionStatus", mock.Anything, idempotencyKey).Return(types.Unknown, errors.New("dummy")).Once()

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
	txm.AssertExpectations(t)
}
