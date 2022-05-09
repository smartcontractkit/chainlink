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
	cfg := config.NewConfig(db.ChainCfg{}, lggr)
	client, err := solanaClient.NewClient(url, cfg, 2*time.Second, lggr)
	require.NoError(t, err)
	getClient := func() (solanaClient.ReaderWriter, error) {
		return client, nil
	}
	txm := soltxm.NewTxm(getClient, cfg, lggr)

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

	// create transfer tx
	hash, err := client.LatestBlockhash()
	assert.NoError(t, err)
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				solana.LAMPORTS_PER_SOL,
				pubKey,
				pubKeyReceiver,
			).Build(),
		},
		hash.Value.Blockhash,
		solana.TransactionPayer(pubKey),
	)
	assert.NoError(t, err)

	// sign tx
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privKey.PublicKey().Equals(key) {
				return &privKey
			}
			return nil
		},
	)
	assert.NoError(t, err)

	// enqueue tx
	assert.NoError(t, txm.Enqueue("testTransmission", tx))
	time.Sleep(time.Second) // wait for tx

	// check balance changes
	senderBal, err := client.Balance(pubKey)
	assert.NoError(t, err)
	assert.Greater(t, initBal, senderBal)
	assert.Greater(t, initBal-senderBal, solana.LAMPORTS_PER_SOL) // balance change = sent + fees

	receiverBal, err := client.Balance(pubKeyReceiver)
	assert.NoError(t, err)
	assert.Equal(t, solana.LAMPORTS_PER_SOL, receiverBal)
}
