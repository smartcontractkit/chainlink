package txmgr_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	evmassets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInMemoryStore_CreateTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	persistentStore := cltest.NewTestTxStore(t, db, cfg.Database())
	kst := cltest.NewKeyStore(t, db, cfg.Database())

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	toAddress := testutils.NewAddress()
	gasLimit := uint32(1000)
	payload := []byte{1, 2, 3}

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := context.Background()

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore)
	require.NoError(t, err)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		actTx, err := inMemoryStore.CreateTransaction(testutils.Context(t), evmtxmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		}, chainID)
		require.NoError(t, err)

		assert.Greater(t, actTx.ID, int64(0))
		assert.Equal(t, commontxmgr.TxUnstarted, actTx.State)
		assert.Equal(t, gasLimit, actTx.FeeLimit)
		assert.Equal(t, fromAddress, actTx.FromAddress)
		assert.Equal(t, toAddress, actTx.ToAddress)
		assert.Equal(t, payload, actTx.EncodedPayload)
		assert.Equal(t, big.Int(evmassets.NewEthValue(0)), actTx.Value)
		assert.Equal(t, subject, actTx.Subject.UUID)

		cltest.AssertCount(t, db, "evm.txes", 1)

		var dbEthTx evmtxmgr.DbEthTx
		require.NoError(t, db.Get(&dbEthTx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, commontxmgr.TxUnstarted, dbEthTx.State)
		assert.Equal(t, gasLimit, dbEthTx.GasLimit)
		assert.Equal(t, fromAddress, dbEthTx.FromAddress)
		assert.Equal(t, toAddress, dbEthTx.ToAddress)
		assert.Equal(t, payload, dbEthTx.EncodedPayload)
		assert.Equal(t, evmassets.NewEthValue(0), dbEthTx.Value)
		assert.Equal(t, subject, dbEthTx.Subject.UUID)
	})
}
