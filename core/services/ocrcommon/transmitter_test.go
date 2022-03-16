package ocrcommon_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
)

func Test_Transmitter_CreateEthTransaction(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	gasLimit := uint64(1000)
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := new(txmmocks.TxManager)
	strategy := new(txmmocks.TxStrategy)

	transmitter := ocrcommon.NewTransmitter(txm, fromAddress, gasLimit, strategy, txmgr.TransmitCheckerSpec{})

	txm.On("CreateEthTransaction", txmgr.NewTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Meta:           nil,
		Strategy:       strategy,
	}, mock.Anything).Return(txmgr.EthTx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(context.Background(), toAddress, payload))

	txm.AssertExpectations(t)
}
