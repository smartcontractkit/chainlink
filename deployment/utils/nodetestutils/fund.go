package nodetestutils

import (
	"encoding/base64"
	"strconv"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/smartcontractkit/chainlink/deployment/internal/aptostestutils"

	sui_common "github.com/smartcontractkit/chainlink-sui/bindings/bind"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
)

var (
	// tonTransmitterFundAmount is the amount of TON to fund the transmitter address.
	tonTransmitterFundAmount = tlb.MustFromTON("200")
)

// fundNodesTon funds the given nodes with the given amount of TON.
func fundNodesTon(t *testing.T, tonChain cldf_ton.Chain, nodes []*Node) {
	messages := make([]*wallet.Message, 0, len(nodes))
	for _, node := range nodes {
		tonkeys, err := node.App.GetKeyStore().TON().GetAll()
		require.NoError(t, err)
		require.Len(t, tonkeys, 1)
		transmitter := tonkeys[0].PubkeyToAddress()
		msg, err := tonChain.Wallet.BuildTransfer(transmitter, tonTransmitterFundAmount, false, "")
		require.NoError(t, err)
		messages = append(messages, msg)
	}
	_, _, err := tonChain.Wallet.SendManyWaitTransaction(t.Context(), messages)
	require.NoError(t, err)
}

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
	signerAddr, _ := signer.GetAddress()

	getCoinsReq := models.SuiXGetAllCoinsRequest{Owner: signerAddr, Limit: 50}
	allCoins, _ := client.SuiXGetAllCoins(ctx, getCoinsReq)

	coins := allCoins.Data[1:]

	require.GreaterOrEqual(t, len(coins), len(nodes))

	for i, node := range nodes {
		suiKeys, err := node.App.GetKeyStore().Sui().GetAll()
		require.NoError(t, err)
		require.Len(t, suiKeys, 1)

		transmitter := suiKeys[0]
		coin := coins[i]
		to := "0x" + transmitter.Account()
		client := suiChain.Client

		balance, _ := strconv.ParseUint(coin.Balance, 10, 64)
		gas := uint64(100_000_000)
		if balance <= gas {
			t.Logf("Skipping coin %s (too small: %d)", coin.CoinObjectId, balance)
			return
		}

		transferAmount := balance - gas

		t.Logf("Transferring coin %s to %s (amount=%d)...", coin.CoinObjectId, to, transferAmount)

		unsignedReq := models.TransferSuiRequest{
			Signer:      signerAddr,
			SuiObjectId: coin.CoinObjectId,
			GasBudget:   strconv.FormatUint(gas, 10),
			Recipient:   to,
			Amount:      strconv.FormatUint(transferAmount, 10),
		}

		txnMeta, err := client.TransferSui(ctx, unsignedReq)
		require.NoError(t, err, "failed to create unsigned transfer txn for %s", coin.CoinObjectId)

		decodedTx, err := base64.StdEncoding.DecodeString(txnMeta.TxBytes)
		require.NoError(t, err, "failed to decode tx bytes for %s", coin.CoinObjectId)

		tx, err := sui_common.SignAndSendTx(ctx, signer, client, decodedTx, true)
		require.NoError(t, err, "failed to execute transfer for coin %s", coin.CoinObjectId)

		t.Logf("Transferred coin %s to %s, Digest: %s, Status: %s",
			coin.CoinObjectId, to, tx.Digest, tx.Effects.Status.Status)

		time.Sleep(300 * time.Millisecond)
	}
}
