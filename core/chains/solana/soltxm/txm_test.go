//go:build integration

package soltxm_test

import (
	"context"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink/core/chains/solana/soltxm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxm_Integration(t *testing.T) {
	url := solanaClient.SetupLocalSolNode(t)
	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PublicKey()
	solanaClient.FundTestAccounts(t, []solana.PublicKey{pubKey}, url)

	// set up txm
	lggr := logger.TestLogger(t)
	confirmDuration := models.MustMakeDuration(500 * time.Millisecond)
	cfg := config.NewConfig(db.ChainCfg{
		ConfirmPollPeriod: &confirmDuration,
	}, lggr)
	client, err := solanaClient.NewClient(url, cfg, 2*time.Second, lggr)
	require.NoError(t, err)
	getClient := func() (solanaClient.ReaderWriter, error) {
		return client, nil
	}
	txm := soltxm.NewTxm("localnet", getClient, cfg, lggr)

	// track initial balance
	initBal, err := client.Balance(pubKey)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), initBal) // should be funded

	// start
	require.NoError(t, txm.Start(context.Background()))

	// already started
	assert.Error(t, txm.Start(context.Background()))

	// create receiver
	privKeyReceiver, err := solana.NewRandomPrivateKey()
	pubKeyReceiver := privKeyReceiver.PublicKey()

	createTx := func(signer *solana.PrivateKey, receiver solana.PublicKey, amt uint64) *solana.Transaction {
		// create transfer tx
		hash, err := client.LatestBlockhash()
		assert.NoError(t, err)
		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					amt,
					pubKey,
					receiver,
				).Build(),
			},
			hash.Value.Blockhash,
			solana.TransactionPayer(pubKey),
		)
		require.NoError(t, err)

		// sign tx
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				return signer
			},
		)
		require.NoError(t, err)

		return tx
	}

	// enqueue txs
	assert.NoError(t, txm.Enqueue("test_success_0", createTx(&privKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL)))
	assert.NoError(t, txm.Enqueue("test_invalidSigner", createTx(&privKeyReceiver, pubKeyReceiver, solana.LAMPORTS_PER_SOL)))
	assert.NoError(t, txm.Enqueue("test_invalidReceiver", createTx(&privKey, solana.PublicKey{}, solana.LAMPORTS_PER_SOL)))
	time.Sleep(500 * time.Millisecond) // pause 0.5s for new blockhash
	assert.NoError(t, txm.Enqueue("test_success_1", createTx(&privKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL)))
	assert.NoError(t, txm.Enqueue("test_txFail", createTx(&privKey, pubKeyReceiver, 1000*solana.LAMPORTS_PER_SOL)))

	// check to make sure all txs are closed out from cache (longest tx should last 5s)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
