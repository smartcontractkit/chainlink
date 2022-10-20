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
	"gopkg.in/guregu/null.v4"

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

const (
	LOADTEST_N            = 1000
	LOADTEST_TICKSPERSLOT = solanaClient.DEFAULT_TICKS_PER_SLOT * 10
)

type CreateTx func(signer, sender, receiver solana.PublicKey, amt uint64) *solana.Transaction

// helper for building tx
func XXXTxWithBlockhash(t *testing.T, signer solana.PublicKey, sender solana.PublicKey, receiver solana.PublicKey, amt uint64, hash solana.Hash) *solana.Transaction {
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amt,
				sender,
				receiver,
			).Build(),
		},
		hash,
		solana.TransactionPayer(signer),
	)
	require.NoError(t, err)
	return tx
}

// helper for setting up txm
func XXXSetupTxm(t *testing.T, url string, ks keystore.Solana) (*soltxm.Txm, solanaClient.ReaderWriter, config.Config, CreateTx) {
	lggr := logger.TestLogger(t)
	cfg := config.NewConfig(db.ChainCfg{
		MaxRetries: null.IntFrom(0),
	}, lggr)
	client, err := solanaClient.NewClient(url, cfg, 2*time.Second, lggr)
	require.NoError(t, err)
	getClient := func() (solanaClient.ReaderWriter, error) {
		return client, nil
	}
	txm := soltxm.NewTxm("localnet", getClient, cfg, ks, lggr)

	var createTx = func(signer solana.PublicKey, sender solana.PublicKey, receiver solana.PublicKey, amt uint64) *solana.Transaction {
		// create transfer tx
		hash, err := client.LatestBlockhash()
		require.NoError(t, err)
		return XXXTxWithBlockhash(t, signer, sender, receiver, amt, hash.Value.Blockhash)
	}

	return txm, client, cfg, createTx
}

func XXXConfirmDone(t *testing.T, ctx context.Context, txm *soltxm.Txm) {
	// check to make sure all txs are closed out from inflight list
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second) // closes test after 30s
	t.Cleanup(cancel)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			idCount, sigCount := txm.InflightTxs()
			require.Equal(t, 0, idCount)  // no unique TXs pending
			require.Equal(t, 0, sigCount) // no signatures pending confirmation
			break loop
		case <-ticker.C:
			id, sig := txm.InflightTxs()
			if id == 0 && sig == 0 {
				cancel() // exit for loop
				break loop
			}
			t.Logf("tx count: IDs - %d, sigs - %d", id, sig)
		}
	}
}

func XXXLoadTest(t *testing.T, ctx context.Context, txm *soltxm.Txm, createTx CreateTx, key solana.PublicKey) time.Duration {
	start := time.Now()
	for i := 0; i < LOADTEST_N; i++ {
		assert.NoError(t, txm.Enqueue(fmt.Sprintf("load_%d", i), createTx(key, key, key, uint64(i))))
		time.Sleep(5 * time.Millisecond) // ~100 txs per second (note: have run 5ms delays for ~200tx/s succesfully)
	}

	XXXConfirmDone(t, ctx, txm)
	return time.Since(start)
}

func XXXNetworkSpam(t *testing.T, close chan struct{}, client solanaClient.ReaderWriter, k solana.PrivateKey) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	key := k.PublicKey()

	data, err := soltxm.ComputeUnitPrice(1).Data()
	// data, err := soltxm.MAX_COMPUTE_UNIT_LIMIT.Data()
	require.NoError(t, err)

	// get initial hash
	hash, err := client.LatestBlockhash()
	require.NoError(t, err)

	i := 0
	for {
		select {
		case <-close:
			t.Log("NetworkSpam stopped")
			return
		case <-ticker.C:
			// update latest blockhash every 10 sends
			if i%10 == 0 {
				hash, err = client.LatestBlockhash()
				require.NoError(t, err)
			}

			// build tx with max compute unit limit
			tx, err := solana.NewTransaction(
				[]solana.Instruction{
					system.NewTransferInstruction(
						uint64(i),
						key,
						key,
					).Build(),
				},
				hash.Value.Blockhash,
				solana.TransactionPayer(key),
			)
			require.NoError(t, err)
			tx.Message.AccountKeys = append(tx.Message.AccountKeys, soltxm.MAX_COMPUTE_UNIT_LIMIT.ProgramID())
			tx.Message.Instructions = append(tx.Message.Instructions, solana.CompiledInstruction{
				ProgramIDIndex: uint16(len(tx.Message.AccountKeys) - 1),
				Data:           data,
			})

			// sign & send tx
			tx.Sign(func(_ solana.PublicKey) *solana.PrivateKey {
				return &k
			})

			_, txErr := client.SendTx(ctx, tx)
			require.NoError(t, txErr)
		}
		i++
	}
}

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
	mkey.On("Get", pubKeyReceiver.String()).Return(solkey.Key{}, keystore.KeyNotFoundError{ID: pubKeyReceiver.String(), KeyType: "Solana"})

	// set up txm
	txm, client, cfg, createTx := XXXSetupTxm(t, url, mkey)

	// track initial balance
	initBal, err := client.Balance(pubKey)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), initBal) // should be funded

	// start
	require.NoError(t, txm.Start(ctx))
	t.Cleanup(func() {
		require.NoError(t, txm.Close())
	})

	// already started
	assert.Error(t, txm.Start(ctx))

	// assert fee change
	t.Run("fees", func(t *testing.T) {
		// enqueue txs (must pass to move on to load test)
		require.NoError(t, txm.Enqueue("test_success_0", createTx(pubKey, pubKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL)))
		time.Sleep(500 * time.Millisecond) // wait for balance to change, new blockhash
		balance0, err := client.Balance(pubKey)
		require.NoError(t, err)
		fee0 := initBal - balance0 - solana.LAMPORTS_PER_SOL // fee used for first tx
		assert.Equal(t, uint64(5000), fee0)                  // base fee for 1 signature

		// change fee
		cfg.Update(db.ChainCfg{
			DefaultComputeUnitPrice: null.IntFrom(1000), // change fee
		})

		require.NoError(t, txm.Enqueue("test_success_1", createTx(pubKey, pubKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL)))
		time.Sleep(500 * time.Millisecond) // wait for balance to change
		balance1, err := client.Balance(pubKey)
		require.NoError(t, err)
		fee1 := balance0 - balance1 - solana.LAMPORTS_PER_SOL // fee used for second tx
		require.Greater(t, fee1, fee0)                        // second tx should have higher fee

		// check balance changes
		assert.NoError(t, err)
		assert.Greater(t, initBal, balance1)
		assert.Greater(t, initBal-balance1, 2*solana.LAMPORTS_PER_SOL) // balance change = sent + fees

		receiverBal, err := client.Balance(pubKeyReceiver)
		assert.NoError(t, err)
		assert.Equal(t, 2*solana.LAMPORTS_PER_SOL, receiverBal)

		XXXConfirmDone(t, ctx, txm) // confirm inflight txs are complete
	})

	// assert duplicate tx handled by dropping the duplicate
	t.Run("duplicate", func(t *testing.T) {
		tx := createTx(pubKey, pubKey, pubKeyReceiver, 0)
		require.NoError(t, txm.Enqueue("test_duplicate_0", tx))
		require.NoError(t, txm.Enqueue("test_duplicate_1", tx))

		XXXConfirmDone(t, ctx, txm)
	})

	t.Run("invalid", func(t *testing.T) {
		// TXM should error immediately
		require.Error(t, txm.Enqueue("test_invalidSigner", createTx(pubKeyReceiver, pubKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL))) // cannot sign tx before enqueuing

		// TXM will see onchain error (accept initial tx)
		require.NoError(t, txm.Enqueue("test_invalidReceiver", createTx(pubKey, pubKey, solana.PublicKey{}, solana.LAMPORTS_PER_SOL)))
		require.NoError(t, txm.Enqueue("test_invalidAmount", createTx(pubKey, pubKey, pubKeyReceiver, 1000*solana.LAMPORTS_PER_SOL)))
		// invalid or outdated blockhash is simply dropped by network and can never be confirmed
		require.NoError(t, txm.Enqueue("test_invalidBlockhash", XXXTxWithBlockhash(t, pubKey, pubKey, pubKeyReceiver, solana.LAMPORTS_PER_SOL, solana.Hash{})))

		XXXConfirmDone(t, ctx, txm)
	})

}

// func TestTxm_Congestion(t *testing.T) {
// 	ctx := testutils.Context(t)
// 	// url := solanaClient.SetupLocalSolNodeOpts(t, LOADTEST_TICKSPERSLOT)
// 	url := "http://localhost:8899"

// 	// txm key
// 	key, err := solkey.New()
// 	require.NoError(t, err)

// 	// spam keys
// 	spamN := 15
// 	var spam []solana.PrivateKey
// 	var spamFund []solana.PublicKey
// 	for i := 0; i < spamN; i++ {
// 		spamKey, err := solana.NewRandomPrivateKey()
// 		require.NoError(t, err)
// 		spam = append(spam, spamKey)
// 		spamFund = append(spamFund, spamKey.PublicKey())
// 	}

// 	// fund keys
// 	solanaClient.FundTestAccounts(t, append(spamFund, key.PublicKey()), url)

// 	// setup mock keystore
// 	mkey := mocks.NewSolana(t)
// 	mkey.On("Get", key.ID()).Return(key, nil)

// 	// set up txm
// 	txm, client, cfg, createTx := XXXSetupTxm(t, url, mkey)

// 	// start
// 	require.NoError(t, txm.Start(ctx))
// 	t.Cleanup(func() {
// 		require.NoError(t, txm.Close())
// 	})

// 	// already started
// 	assert.Error(t, txm.Start(ctx))

// 	// track times
// 	var noCongestion time.Duration
// 	var congestedNoFees time.Duration
// 	var congestedWithFees time.Duration

// 	t.Log("Benchmarking: No Congestion")
// 	// load test: try to overload txs, confirm
// 	// benchmark for congestion testing
// 	// set fees to no bumping
// 	cfg.Update(db.ChainCfg{
// 		DefaultComputeUnitPrice: null.IntFrom(0),
// 		MinComputeUnitPrice:     null.IntFrom(0),
// 		MaxComputeUnitPrice:     null.IntFrom(0),
// 	})

// 	noCongestion = XXXLoadTest(t, ctx, txm, createTx, key.PublicKey())

// 	// start spammers
// 	stop := make(chan struct{})
// 	var wg sync.WaitGroup
// 	wg.Add(spamN)
// 	for i := 0; i < spamN; i++ {
// 		go func(ind int) {
// 			go XXXNetworkSpam(t, stop, client, spam[ind])
// 			wg.Done()
// 		}(i)
// 	}

// 	t.Log("Benchmarking: Congested + No Fee Txm")
// 	// note this will show errors regarding duplicate signatures because the fee cannot be bumped so the rebroadcast tx is identical => identical signature
// 	congestedNoFees = XXXLoadTest(t, ctx, txm, createTx, key.PublicKey())

// 	t.Log("Benchmarking: Congested + Fee Enhanced Txm")
// 	// set fees to no bumping
// 	cfg.Update(db.ChainCfg{
// 		DefaultComputeUnitPrice: null.IntFrom(2),
// 		MinComputeUnitPrice:     null.IntFrom(0),
// 		MaxComputeUnitPrice:     null.IntFrom(1_000_000),
// 	})
// 	congestedWithFees = XXXLoadTest(t, ctx, txm, createTx, key.PublicKey())

// 	// log accounts used
// 	fmt.Println("Sender", key.PublicKey())
// 	fmt.Println("Spammers", spamFund)

// 	// stop spammers
// 	close(stop)
// 	wg.Wait()

// 	t.Logf("Benchmark:\n- No Congestion = %d\n- Congested (No Fees) = %d\n- Congested (With Fees) = %d", noCongestion, congestedNoFees, congestedWithFees)
// 	assert.True(t, noCongestion < congestedNoFees, "congestedNoFees should take longer than noCongestion")
// 	assert.True(t, congestedNoFees > congestedWithFees, "congestedNoFees should take longer than congestedWithFees")
// }