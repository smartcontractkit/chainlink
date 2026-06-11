package nodetestutils

import (
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
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

// fundNodesSol funds the given nodes with the given amount of SOL.
func fundNodesSol(t *testing.T, solChain cldf_solana.Chain, nodes []*Node) {
	for _, node := range nodes {
		solkeys, err := node.App.GetKeyStore().Solana().GetAll()
		require.NoError(t, err)
		require.Len(t, solkeys, 1)
		transmitter := solkeys[0]
		_, err = solChain.Client.RequestAirdrop(t.Context(), transmitter.PublicKey(), 1000*solana.LAMPORTS_PER_SOL, solRpc.CommitmentConfirmed)
		require.NoError(t, err)
		// we don't wait for confirmation so we don't block the tests, it'll take a while before nodes start transmitting
	}
}

// fundNodesSol funds the given nodes with the given amount of SUI.
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

		t.Logf("Transferring coin %s (balance=%d) to %s...", coin.GetObjectId(), coin.GetBalance(), to)

		// Resolve the coin into an ImmOrOwnedObject input; ptb.Object would leave it
		// Unresolved (BCS variant 3), which the network rejects.
		coinIDBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(coin.GetObjectId()))
		require.NoError(t, err, "failed to convert coin object id %s", coin.GetObjectId())

		resolvedCoin, err := resolver.ResolveCallArg(ctx, &transaction.CallArg{
			UnresolvedObject: &transaction.UnresolvedObject{ObjectId: *coinIDBytes},
		}, "")
		require.NoError(t, err, "failed to resolve coin object %s", coin.GetObjectId())

		// Transfer the whole coin to the node. Gas is paid from the reserved coin so the
		// transferred coin object is not consumed for gas.
		ptb := transaction.NewTransaction()
		coinArg := ptb.Data.V1.AddInput(*resolvedCoin)
		ptb.TransferObjects([]transaction.Argument{coinArg}, ptb.Pure(to))

		callOpts := &sui_common.CallOpts{
			Signer:           signer,
			GasObject:        gasCoinID,
			WaitForExecution: true,
		}

		tx, err := sui_common.ExecutePTB(ctx, callOpts, client, ptb)
		require.NoError(t, err, "failed to execute transfer for coin %s", coin.GetObjectId())

		t.Logf("Transferred coin %s to %s, Digest: %s, Status: %s",
			coin.GetObjectId(), to, tx.Digest, tx.Effects.Status.Status)

		time.Sleep(300 * time.Millisecond)
	}
}
