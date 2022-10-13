package ocrcommon_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
)

func Test_DefaultTransmitter_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	gasLimit := uint32(1000)
	effectiveTransmitterAddress := fromAddress
	toAddress := testutils.NewAddress()
	payload := []byte{1, 2, 3}
	txm := txmmocks.NewTxManager(t)
	strategy := txmmocks.NewTxStrategy(t)

	transmitter := ocrcommon.NewTransmitter(txm, []common.Address{fromAddress}, gasLimit, effectiveTransmitterAddress, strategy, txmgr.TransmitCheckerSpec{})

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
