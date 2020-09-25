package offchainreporting_test

import (
	"context"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
)

func Test_Transmitter_CreateEthTransaction(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	db := store.DB.DB()

	gasLimit := uint64(1000)
	fromAddress := gethCommon.HexToAddress(cltest.DefaultKey)
	toAddress := cltest.NewAddress()
	payload := []byte{1, 2, 3}

	transmitter := offchainreporting.NewTransmitter(db, fromAddress, gasLimit)

	require.NoError(t, transmitter.CreateEthTransaction(context.Background(), toAddress, payload))

	etx := models.EthTx{}
	require.NoError(t, store.ORM.DB.First(&etx).Error)

	require.Equal(t, gasLimit, etx.GasLimit)
	require.Equal(t, fromAddress, etx.FromAddress)
	require.Equal(t, toAddress, etx.ToAddress)
	require.Equal(t, payload, etx.EncodedPayload)
	require.Equal(t, assets.NewEthValue(0), etx.Value)
}
