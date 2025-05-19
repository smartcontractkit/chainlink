package ocrcommon_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
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

	transmitter, err := ocrcommon.NewTransmitter(
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

	transmitter, err := ocrcommon.NewTransmitter(
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

	fromAddress := testutils.NewAddress()

	gasLimit := uint64(1000)
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
		keystest.Addresses{},
	)
	require.NoError(t, err)
	require.Error(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_No_Keystore_Error(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	fromAddress2 := testutils.NewAddress()

	gasLimit := uint64(1000)
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
		nil,
	)
	require.Error(t, err)
}

func Test_DualTransmitter(t *testing.T) {
	t.Parallel()

	memoryKeystore := keystest.NewMemoryChainStore()
	fromAddress := memoryKeystore.MustCreate(t)
	secondaryFromAddress := memoryKeystore.MustCreate(t)

	contractAddress := utils.RandomAddress()
	secondaryContractAddress := utils.RandomAddress()

	gasLimit := uint64(1000)
	effectiveTransmitterAddress := fromAddress
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewMockEvmTxManager(t)
	strategy := newMockTxStrategy(t)
	dualTransmissionConfig := &types.DualTransmissionConfig{
		ContractAddress:    secondaryContractAddress,
		TransmitterAddress: secondaryFromAddress,
		Meta: map[string][]string{
			"key1": {"value1"},
			"key2": {"value2", "value3"},
			"key3": {"value4", "value5", "value6"},
		},
	}

	transmitter, err := ocrcommon.NewOCR2FeedsTransmitter(
		txm,
		[]common.Address{fromAddress},
		contractAddress,
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		txmgr.TransmitCheckerSpec{},
		keys.NewStore(memoryKeystore),
		dualTransmissionConfig,
	)
	require.NoError(t, err)

	primaryTxConfirmed := false
	secondaryTxConfirmed := false

	txm.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx txmgr.TxRequest) bool {
		switch tx.FromAddress {
		case fromAddress:
			// Primary transmission
			assert.Equal(t, tx.ToAddress, toAddress, "unexpected primary toAddress")
			assert.Nil(t, tx.Meta, "Meta should be empty")
			primaryTxConfirmed = true
		case secondaryFromAddress:
			// Secondary transmission
			assert.Equal(t, tx.ToAddress, secondaryContractAddress, "unexpected secondary toAddress")
			assert.True(t, *tx.Meta.DualBroadcast, "DualBroadcast should be true")
			assert.Equal(t, "key1=value1&key2=value2&key2=value3&key3=value4&key3=value5&key3=value6", *tx.Meta.DualBroadcastParams, "DualBroadcastParams not equal")
			secondaryTxConfirmed = true
		default:
			// Should never be reached
			return false
		}

		return true
	})).Twice().Return(txmgr.Tx{}, nil)

	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload, nil))
	require.NoError(t, transmitter.CreateSecondaryEthTransaction(testutils.Context(t), payload, nil))

	require.True(t, primaryTxConfirmed)
	require.True(t, secondaryTxConfirmed)
}
