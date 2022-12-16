//go:build integration

package soltxm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	relayutils "github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink/core/chains/solana/soltxm"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/mocks"

	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
)

func TestTxm_Integration(t *testing.T) {
	ctx := testutils.Context(t)
	url := solanaClient.SetupLocalSolNode(t)

	// setup key
	key, err := solkey.New()
	require.NoError(t, err)
	pubKey := key.PublicKey()

	// setup load test key
	loadTestKey, err := solkey.New()
	require.NoError(t, err)

	// setup receiver key
	privKeyReceiver, err := solana.NewRandomPrivateKey()
	pubKeyReceiver := privKeyReceiver.PublicKey()

	// fund keys
	solanaClient.FundTestAccounts(t, []solana.PublicKey{pubKey, loadTestKey.PublicKey()}, url)

	// setup mock keystore
	mkey := mocks.NewSolana(t)
	mkey.On("Get", key.ID()).Return(key, nil)
	mkey.On("Get", loadTestKey.ID()).Return(loadTestKey, nil)
	mkey.On("Get", pubKeyReceiver.String()).Return(solkey.Key{}, keystore.KeyNotFoundError{ID: pubKeyReceiver.String(), KeyType: "Solana"})

	// set up txm
	lggr := logger.TestLogger(t)
	confirmDuration, err := relayutils.NewDuration(500 * time.Millisecond)
	require.NoError(t, err)
	cfg := config.NewConfig(db.ChainCfg{
		ConfirmPollPeriod: &confirmDuration,
	}, lggr)
	client, err := solanaClient.NewClient(url, cfg, 2*time.Second, lggr)
	require.NoError(t, err)
	getClient := func() (solanaClient.ReaderWriter, error) {
		return client, nil
	}
	txm := soltxm.NewTxm("localnet", getClient, cfg, mkey, lggr)

	// track initial balance
	initBal, err := client.Balance(pubKey)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), initBal) // should be funded

	// start
	require.NoError(t, txm.Start(ctx))

	// already started
	assert.Error(t, txm.Start(ctx))

	createTx := func(signer solana.PublicKey, sender solana.PublicKey, receiver solana.PublicKey, amt uint64) *solana.Transaction {
		// create transfer tx
		hash, err := client.LatestBlockhash()
		assert.NoError(t, err)
		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					amt,
					sender,
					receiver,
				).Build(),
			},
			hash.Value.Blockhash,
			solana.TransactionPayer(signer),
		)
		require.NoError(t, err)
		return tx
	}

	// enqueue txs (must pass to move on to load test)
	require.NoError(t, txm.Enqueue("test_success_0", createTx(pubKey, pubKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL)))
	require.Error(t, txm.Enqueue("test_invalidSigner", createTx(pubKeyReceiver, pubKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL))) // cannot sign tx before enqueuing
	require.NoError(t, txm.Enqueue("test_invalidReceiver", createTx(pubKey, pubKey, solana.PublicKey{}, solana.LAMPORTS_PER_SOL)))
	time.Sleep(500 * time.Millisecond) // pause 0.5s for new blockhash
	require.NoError(t, txm.Enqueue("test_success_1", createTx(pubKey, pubKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL)))
	require.NoError(t, txm.Enqueue("test_txFail", createTx(pubKey, pubKey, pubKeyReceiver, 1000*solana.LAMPORTS_PER_SOL)))

	// load test: try to overload txs, confirm, or simulation
	for i := 0; i < 1000; i++ {
		assert.NoError(t, txm.Enqueue(fmt.Sprintf("load_%d", i), createTx(loadTestKey.PublicKey(), loadTestKey.PublicKey(), loadTestKey.PublicKey(), uint64(i))))
		time.Sleep(10 * time.Millisecond) // ~100 txs per second (note: have run 5ms delays for ~200tx/s succesfully)
	}

	// check to make sure all txs are closed out from inflight list (longest should last MaxConfirmTimeout)
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			assert.Equal(t, 0, txm.InflightTxs())
			break loop
		case <-ticker.C:
			if txm.InflightTxs() == 0 {
				cancel() // exit for loop
			}
		}
	}
	assert.NoError(t, txm.Close())

	// check balance changes
	senderBal, err := client.Balance(pubKey)
	assert.NoError(t, err)
	assert.Greater(t, initBal, senderBal)
	assert.Greater(t, initBal-senderBal, 2*solana.LAMPORTS_PER_SOL) // balance change = sent + fees

	receiverBal, err := client.Balance(pubKeyReceiver)
	assert.NoError(t, err)
	assert.Equal(t, 2*solana.LAMPORTS_PER_SOL, receiverBal)
}
