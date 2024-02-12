package txm_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	txm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/optimistictxm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func newReaperWithChainID(t *testing.T, txStore txm.ReaperTxStore, cfg txm.ReaperConfig, cid *big.Int, client txm.ReaperClient, ks txm.KeyStore) *txm.Reaper {
	return txm.NewReaper(logger.Test(t), txStore, cfg, cid, client, ks)
}

func newReaper(t *testing.T, txStore txm.ReaperTxStore, cfg txm.ReaperConfig, client txm.ReaperClient, ks txm.KeyStore) *txm.Reaper {
	return newReaperWithChainID(t, txStore, cfg, &cltest.FixtureChainID, client, ks)
}

func TestReaper_ReapTxs(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	keyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	client := evmtest.NewEthClientMockWithDefaultChain(t)

	oneDayAgo := time.Now().Add(-24 * time.Hour)

	encodedPayload := []byte{1, 2, 3}
	value := big.Int(assets.NewEthValue(142))
	gasLimit := uint32(242)

	t.Run("with nothing in the database, doesn't error", func(t *testing.T) {
		rc := txm.ReaperConfig{ReaperThreshold: 1 * time.Hour}
		r := newReaper(t, txStore, rc, client, keyStore)

		err := r.ReapTxs()
		assert.NoError(t, err)
	})

	t.Run("skips if threshold=0", func(t *testing.T) {
		rc := txm.ReaperConfig{ReaperThreshold: 0 * time.Second}
		r := newReaper(t, txStore, rc, client, keyStore)

		err := r.ReapTxs()
		assert.NoError(t, err)
	})

	t.Run("doesn't touch txs with different chain ID", func(t *testing.T) {
		rc := txm.ReaperConfig{ReaperThreshold: 1 * time.Hour}
		key1, addr1 := cltest.MustInsertRandomKey(t, keyStore)
		chainID := big.NewInt(42)
		r := newReaperWithChainID(t, txStore, rc, chainID, client, keyStore)

		nonce0 := evmtypes.Nonce(0)
		txConfirmed := txm.Tx{
			ChainID:            big.NewInt(99),
			Sequence:           &nonce0,
			FromAddress:        addr1,
			ToAddress:          utils.RandomAddress(),
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &oneDayAgo,
			InitialBroadcastAt: &oneDayAgo,
			Error:              null.String{},
			State:              commontxmgr.TxConfirmed,
		}

		require.NoError(t, txStore.InsertTx(&txConfirmed))

		err := r.ReapTxs()
		assert.NoError(t, err)
		// Didn't delete because eth_tx has chain ID of 0
		cltest.AssertCount(t, db, "evm.txes", 1)
		keyStore.Delete(key1.ID())
		pgtest.MustExec(t, db, `DELETE FROM evm.txes`)
	})

	t.Run("deletes confirmed evm.txes that exceed the age threshold", func(t *testing.T) {
		rc := txm.ReaperConfig{ReaperThreshold: 1 * time.Hour}
		_, addr2 := cltest.MustInsertRandomKey(t, keyStore)
		r := newReaper(t, txStore, rc, client, keyStore)

		nonce0 := evmtypes.Nonce(0)
		timeNow := time.Now()
		txUnconfirmed := txm.Tx{
			ChainID:            &cltest.FixtureChainID,
			Sequence:           &nonce0,
			FromAddress:        addr2,
			ToAddress:          utils.RandomAddress(),
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              commontxmgr.TxConfirmed,
		}

		require.NoError(t, txStore.InsertTx(&txUnconfirmed))

		client.On("SequenceAt", mock.Anything, addr2, mock.Anything).Return(evmtypes.Nonce(0), nil).Twice()
		err := r.ReapTxs()
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, db, "evm.txes", 1)

		pgtest.MustExec(t, db, `UPDATE evm.txes SET created_at=$1, broadcast_at=$1`, oneDayAgo)

		err = r.ReapTxs()
		assert.NoError(t, err)
		// Didn't delete because eth_tx although old enough, wasn't mined on-chain
		cltest.AssertCount(t, db, "evm.txes", 1)

		client.On("SequenceAt", mock.Anything, addr2, mock.Anything).Return(evmtypes.Nonce(1), nil).Once()
		pgtest.MustExec(t, db, `UPDATE evm.txes SET state= 'confirmed'`)
		err = r.ReapTxs()
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was old enough and nonce seen used on-chain
		cltest.AssertCount(t, db, "evm.txes", 0)
	})
}
