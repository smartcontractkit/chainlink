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

	key := cltest.MustInsertRandomKey(t, store.DB, 0)

	gasLimit := uint64(1000)
	fromAddress := key.Address.Address()
	toAddress := cltest.NewAddress()
	payload := []byte{1, 2, 3}
	txm := new(bptxmmocks.TxManager)
	strategy := new(bptxmmocks.TxStrategy)

	transmitter := offchainreporting.NewTransmitter(txm, store.DB, fromAddress, gasLimit, strategy)

	txm.On("CreateEthTransaction", mock.Anything, fromAddress, toAddress, payload, gasLimit, nil, strategy).Return(bulletprooftxmanager.EthTx{}, nil).Once()
	require.NoError(t, transmitter.CreateEthTransaction(context.Background(), toAddress, payload))

	txm.AssertExpectations(t)
}
