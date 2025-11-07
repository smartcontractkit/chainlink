package memory

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_sui_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui/provider"
	sui_common "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func getTestSuiChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.SUI_LOCALNET.Selector}
}

func randomSeed() []byte {
	seed := make([]byte, ed25519.SeedSize)
	_, err := rand.Read(seed)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random seed: %+v", err))
	}

	return seed
}

func GenerateChainsSui(t *testing.T, numChains int) []cldf_chain.BlockChain {
	testSuiChainSelectors := getTestSuiChainSelectors()
	if len(testSuiChainSelectors) < numChains {
		t.Fatalf("not enough test sui chain selectors available")
	}
	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testSuiChainSelectors[i]

		seeded := ed25519.NewKeyFromSeed(randomSeed()) // 64 bytes: seed||pub
		seed := seeded[:32]                            // or: seeded.Seed() if available
		hexKey := hex.EncodeToString(seed)             // 64 hex chars

		platform := "linux/amd64"
		img := "mysten/sui-tools:devnet"
		// generate adhoc sui privKey
		c, err := cldf_sui_provider.NewCTFChainProvider(t, selector,
			cldf_sui_provider.CTFChainProviderConfig{
				Once:              once,
				DeployerSignerGen: cldf_sui_provider.AccountGenPrivateKey(hexKey),
				Image:             &img,
				Platform:          &platform,
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}

	t.Logf("Created %d Sui chains: %+v", len(chains), chains)
	return chains
}

func createSuiChainConfig(chainID string, chain cldf_sui.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "sui-localnet"
	chainConfig["NetworkNameFull"] = "sui-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}

// func FundSuiAccount(faucetURL string, rpcURL, address string) error {
// 	r := resty.New().SetBaseURL(faucetURL)
// 	b := &models.FaucetRequest{
// 		FixedAmountRequest: &models.FaucetFixedAmountRequest{
// 			Recipient: address,
// 		},
// 	}
// 	resp, err := r.R().SetBody(b).SetHeader("Content-Type", "application/json").Post("/gas")
// 	if err != nil {
// 		return err
// 	}
// 	framework.L.Info().Any("Resp", resp).Msg("Address is funded!")

// 	return nil
// }

func fundSuiNodes(t *testing.T, suiChain cldf_sui.Chain, nodes []*Node) {

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
		FundSuiAccount(t, suiChain.Signer, "0x"+transmitter.Account(), suiChain.Client, coins[i])
	}
}

func FundSuiAccount(t *testing.T, signer cldf_sui.SuiSigner, to string, client sui.ISuiAPI, coin models.CoinData) {
	ctx := t.Context()

	signerAddrStr, err := signer.GetAddress()
	require.NoError(t, err)

	balance, _ := strconv.ParseUint(coin.Balance, 10, 64)
	gas := uint64(100_000_000)
	if balance <= gas {
		t.Logf("Skipping coin %s (too small: %d)", coin.CoinObjectId, balance)
		return
	}

	transferAmount := balance - gas

	t.Logf("Transferring coin %s to %s (amount=%d)...", coin.CoinObjectId, to, transferAmount)

	unsignedReq := models.TransferSuiRequest{
		Signer:      signerAddrStr,
		SuiObjectId: coin.CoinObjectId,
		GasBudget:   fmt.Sprintf("%d", gas),
		Recipient:   to,
		Amount:      fmt.Sprintf("%d", transferAmount),
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
