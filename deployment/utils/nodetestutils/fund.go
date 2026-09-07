package nodetestutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/avast/retry-go/v4"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/internal/aptostestutils"

	sui_common "github.com/smartcontractkit/chainlink-sui/bindings/bind"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
)

// fundNodesAptos funds the given nodes with the given amount of APT.
func fundNodesAptos(t *testing.T, aptosChain cldf_aptos.Chain, nodes []*Node) {
	for _, node := range nodes {
		aptoskeys, err := node.App.GetKeyStore().Aptos().GetAll()
		require.NoError(t, err)
		require.Len(t, aptoskeys, 1)
		transmitter := aptoskeys[0]
		transmitterAccountAddress := aptos.AccountAddress{}
		require.NoError(t, transmitterAccountAddress.ParseStringRelaxed(transmitter.Account()))
		aptostestutils.FundAccount(t, aptosChain.DeployerSigner, transmitterAccountAddress, 100*1e8, aptosChain.Client)
	}
}

// solanaAirdropAttempts is the number of times we retry a node SOL airdrop before giving up. The
// Solana test-validator container can briefly refuse its RPC port (e.g. under runner memory pressure
// the validator process may die and get restarted), so a single attempt flakes the whole test suite.
// A handful of attempts with a fixed delay gives the container a grace window to come back.
const solanaAirdropAttempts = uint(6)

// fundNodesSol funds the given nodes with the given amount of SOL.
func fundNodesSol(t *testing.T, solChain cldf_solana.Chain, nodes []*Node) {
	for _, node := range nodes {
		solkeys, err := node.App.GetKeyStore().Solana().GetAll()
		require.NoError(t, err)
		require.Len(t, solkeys, 1)
		transmitter := solkeys[0]
		err = retry.Do(func() error {
			_, e := solChain.Client.RequestAirdrop(
				t.Context(), transmitter.PublicKey(), 1000*solana.LAMPORTS_PER_SOL, solRpc.CommitmentConfirmed,
			)
			return e
		},
			retry.Context(t.Context()),
			retry.Attempts(solanaAirdropAttempts),
			retry.Delay(2*time.Second),
			retry.DelayType(retry.FixedDelay),
			retry.LastErrorOnly(true),
		)
		require.NoError(t, err, enrichSolanaAirdropErr(t.Context(), solChain, err))
		// we don't wait for confirmation so we don't block the tests, it'll take a while before nodes start transmitting
	}
}

// enrichSolanaAirdropErr turns a bare airdrop failure into a self-diagnosing error by probing the
// validator's /health endpoint. If the RPC is itself unreachable, the validator container likely died
// (OOM / runner resource pressure) rather than the airdrop RPC being transiently flaky.
func enrichSolanaAirdropErr(ctx context.Context, solChain cldf_solana.Chain, airdropErr error) string {
	if airdropErr == nil {
		return ""
	}
	if _, herr := solChain.Client.GetHealth(ctx); herr != nil {
		return fmt.Sprintf(
			"Solana validator RPC unreachable at %s after %d airdrop attempts: %v — validator container likely died (check runner OOM/resource pressure)",
			solChain.URL, solanaAirdropAttempts, airdropErr,
		)
	}
	return fmt.Sprintf(
		"Solana airdrop failed after %d attempts at %s although RPC /health is ok: %v",
		solanaAirdropAttempts, solChain.URL, airdropErr,
	)
}

// suiGasCoinPoolSize is how many gas coins each Sui transmitter is funded with. A Sui transmitter
// signs with a single key shared by the commit and exec plugins, and the relayer selects one whole
// coin per transmission with no coin splitting. A single gas coin therefore serialises all of a
// node's transmits and starves under concurrent commit/exec + retries ("no coins available for gas
// budget" / "coin already reserved"). Funding a pool lets concurrent transmissions each take their
// own coin.
const suiGasCoinPoolSize = uint64(30)

// suiMinGasCoinBalance is the smallest gas coin we will create (0.2 SUI == one transmission's gas
// budget). The pool is shrunk rather than minting dust coins if the source coin can't back the full
// pool at this size.
const suiMinGasCoinBalance = uint64(200_000_000)

// fundNodesSui funds the given nodes with a pool of SUI gas coins.
func fundNodesSui(t *testing.T, suiChain cldf_sui.Chain, nodes []*Node) {
	ctx := t.Context()
	signer := suiChain.Signer
	client := suiChain.Client
	signerAddr, err := signer.GetAddress()
	require.NoError(t, err)

	allCoins, err := client.GetCoinsByAddress(ctx, signerAddr)
	require.NoError(t, err)
	require.NotEmpty(t, allCoins)

	// Reserve the first coin to pay for gas; distribute the remaining coins to the nodes.
	gasCoinID := allCoins[0].GetObjectId()
	coins := allCoins[1:]
	require.GreaterOrEqual(t, len(coins), len(nodes))

	resolver := sui_common.NewObjectResolver(client)

	for i, node := range nodes {
		suiKeys, err := node.App.GetKeyStore().Sui().GetAll()
		require.NoError(t, err)
		require.Len(t, suiKeys, 1)

		transmitter := suiKeys[0]
		coin := coins[i]
		to := "0x" + transmitter.Account()

		// Size the pool to the source coin: never mint coins smaller than one gas budget, so a
		// small source coin just yields a smaller pool instead of unusable dust.
		balance := coin.GetBalance()
		poolSize := suiGasCoinPoolSize
		if balance/poolSize < suiMinGasCoinBalance {
			poolSize = balance / suiMinGasCoinBalance
		}
		require.Positive(t, poolSize,
			"source coin %s balance %d is too small to fund a gas pool", coin.GetObjectId(), balance)
		perCoin := balance / poolSize

		t.Logf("Splitting coin %s (balance=%d) into %d gas coins of %d for %s...",
			coin.GetObjectId(), balance, poolSize, perCoin, to)

		// Resolve the coin into an ImmOrOwnedObject input; ptb.Object would leave it
		// Unresolved (BCS variant 3), which the network rejects.
		coinIDBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(coin.GetObjectId()))
		require.NoError(t, err, "failed to convert coin object id %s", coin.GetObjectId())

		resolvedCoin, err := resolver.ResolveCallArg(ctx, &transaction.CallArg{
			UnresolvedObject: &transaction.UnresolvedObject{ObjectId: *coinIDBytes},
		}, "")
		require.NoError(t, err, "failed to resolve coin object %s", coin.GetObjectId())

		// Split the source coin into a pool of gas coins and transfer each to the node. One
		// single-amount SplitCoins per coin (the SDK can't index multi-amount split results); all
		// commands draw from the same input coin, which the PTB threads through sequentially. Gas is
		// paid from the deployer's reserved coin, so the source coin's full balance is available.
		ptb := transaction.NewTransaction()
		coinArg := ptb.Data.V1.AddInput(*resolvedCoin)
		recipientArg := ptb.Pure(to)
		amountArg := ptb.Pure(perCoin)
		for range poolSize {
			splitCoin := ptb.SplitCoins(coinArg, []transaction.Argument{amountArg})
			ptb.TransferObjects([]transaction.Argument{splitCoin}, recipientArg)
		}

		callOpts := &sui_common.CallOpts{
			Signer:           signer,
			GasObject:        gasCoinID,
			WaitForExecution: true,
		}

		tx, err := sui_common.ExecutePTB(ctx, callOpts, client, ptb)
		require.NoError(t, err, "failed to split/transfer gas coins for node %s", to)

		t.Logf("Funded %s with %d gas coins (coin %s), Digest: %s, Status: %s",
			to, poolSize, coin.GetObjectId(), tx.Digest, tx.Effects.Status.Status)

		time.Sleep(300 * time.Millisecond)
	}
}
