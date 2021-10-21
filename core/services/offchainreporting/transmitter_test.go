package offchainreporting_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Transmitter_CreateEthTransaction(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	gasLimit := uint64(1000)
	toAddress := cltest.NewAddress()
	payload := []byte{1, 2, 3}
	txm := new(bptxmmocks.TxManager)
	strategy := new(bptxmmocks.TxStrategy)

	transmitter := offchainreporting.NewTransmitter(txm, store.DB, fromAddress, gasLimit, strategy)

	txm.On("CreateEthTransaction", mock.Anything, bulletprooftxmanager.NewTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       gasLimit,
		Meta:           nil,
		Strategy:       strategy,
	}).Return(bulletprooftxmanager.EthTx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(context.Background(), toAddress, payload))

	txm.AssertExpectations(t)
}
