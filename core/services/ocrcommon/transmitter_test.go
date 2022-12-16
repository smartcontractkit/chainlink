package ocrcommon_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
)

func Test_DefaultTransmitter_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	gasLimit := uint32(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := fromAddress
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewTxManager(t)
	strategy := txmmocks.NewTxStrategy(t)

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

	txm.On("CreateEthTransaction", txmgr.NewTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Meta:           nil,
		Strategy:       strategy,
	}, mock.Anything).Return(txmgr.EthTx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	gasLimit := uint32(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewTxManager(t)
	strategy := txmmocks.NewTxStrategy(t)

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

	txm.On("CreateEthTransaction", txmgr.NewTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Meta:           nil,
		Strategy:       strategy,
	}, mock.Anything).Return(txmgr.EthTx{}, nil).Once()
	txm.On("CreateEthTransaction", txmgr.NewTx{
		FromAddress:    fromAddress2,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Meta:           nil,
		Strategy:       strategy,
	}, mock.Anything).Return(txmgr.EthTx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload))
	require.NoError(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_Round_Robin_Error(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	fromAddress := common.Address{}

	gasLimit := uint32(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewTxManager(t)
	strategy := txmmocks.NewTxStrategy(t)

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
	require.Error(t, transmitter.CreateEthTransaction(testutils.Context(t), toAddress, payload))
}

func Test_DefaultTransmitter_Forwarding_Enabled_CreateEthTransaction_No_Keystore_Error(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	gasLimit := uint32(1000)
	chainID := big.NewInt(0)
	effectiveTransmitterAddress := common.Address{}
	txm := txmmocks.NewTxManager(t)
	strategy := txmmocks.NewTxStrategy(t)

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
