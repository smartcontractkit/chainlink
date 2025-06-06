package memory

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type TonChain struct {
	Client         *ton.APIClient
	DeployerWallet *wallet.Wallet
}

func getTestTonChainSelectors() []uint64 {
	return []uint64{chainsel.TON_LOCALNET.Selector}
}

func createTonWallet(t *testing.T, client ton.APIClientWrapped, version wallet.Version, option wallet.Option) *wallet.Wallet {
	seed := wallet.NewSeed()
	rw, err := wallet.FromSeed(client, seed, version)
	require.NoError(t, err, fmt.Errorf("Failed to generate random wallet: %v", err))
	pw, perr := wallet.FromPrivateKeyWithOptions(client, rw.PrivateKey(), version, option)
	require.NoError(t, perr)
	require.NoError(t, perr, fmt.Errorf("Failed to generate random wallet: %v", err))
	return pw
}

func fundTonWallets(t *testing.T, client ton.APIClientWrapped, recipients []*address.Address, amounts []tlb.Coins) {
	rawHlWallet, err := wallet.FromSeed(client, strings.Fields(blockchain.DefaultTonHlWalletMnemonic), wallet.HighloadV2Verified)
	require.NoError(t, err, "failed to create highload wallet")
	mcFunderWallet, err := wallet.FromPrivateKeyWithOptions(client, rawHlWallet.PrivateKey(), wallet.HighloadV2Verified, wallet.WithWorkchain(-1))
	require.NoError(t, err, "failed to create highload wallet")
	subWalletID := uint32(42)
	funder, err := mcFunderWallet.GetSubwallet(subWalletID)
	require.NoError(t, err, "failed to get highload subwallet")
	// double check funder address
	require.Equal(t, funder.Address().StringRaw(), blockchain.DefaultTonHlWalletAddress, "funder address mismatch")

	if len(recipients) != len(amounts) {
		t.Fatalf("number of recipients (%d) does not match number of amounts (%d)", len(recipients), len(amounts))
	}

	messages := make([]*wallet.Message, len(recipients))
	for i, addr := range recipients {
		transfer, terr := funder.BuildTransfer(addr, amounts[i], false, "")
		require.NoError(t, terr, fmt.Sprintf("failed to build transfer for %s", addr.String()))
		messages[i] = transfer
	}
	_, _, txerr := funder.SendManyWaitTransaction(t.Context(), messages)
	require.NoError(t, txerr, "airdrop transaction failed")
	// we don't wait for the transaction to be confirmed here, as it may take some time
}

func GenerateChainsTon(t *testing.T, numChains int) map[uint64]cldf_ton.Chain {
	testTonChainSelectors := getTestTonChainSelectors()
	if numChains > 1 {
		t.Fatalf("only one ton chain is supported for now, got %d", numChains)
	}
	if len(testTonChainSelectors) < numChains {
		t.Fatalf("not enough test ton chain selectors available")
	}

	chains := make(map[uint64]cldf_ton.Chain)
	for i := 0; i < numChains; i++ {
		chainID := testTonChainSelectors[i]
		nodeClient := tonChain(t, chainID)
		wallet := createTonWallet(t, nodeClient, wallet.V3R2, wallet.WithWorkchain(0))
		// airdrop the deployer wallet
		fundTonWallets(t, nodeClient, []*address.Address{wallet.Address()}, []tlb.Coins{tlb.MustFromTON("1000")})
		ton := cldf_ton.Chain{
			ChainMetadata: cldf_ton.ChainMetadata{Selector: chainID},
			Client:        nodeClient,
			Wallet:        wallet,
			WalletAddress: wallet.Address(),
		}
		chains[chainID] = ton
	}
	return chains
}

func tonChain(t *testing.T, chainID uint64) *ton.APIClient {
	t.Helper()
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	bcInput := &blockchain.Input{
		ChainID: strconv.FormatUint(chainID, 10),
		Type:    "ton",
		Image:   "ghcr.io/neodix42/mylocalton-docker:latest",
		Port:    strconv.Itoa(freeport.GetOne(t)),
	}
	var bcOut *blockchain.Output

	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		bcOut, err = blockchain.NewBlockchainNetwork(bcInput)
		if err == nil {
			break
		}
		t.Logf("Error creating TON network (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second)
	}
	require.NoError(t, err, "Failed to create blockchain network after %d attempts", maxRetries)
	networkCfg := fmt.Sprintf("http://%s/localhost.global.config.json", bcOut.Nodes[0].ExternalHTTPUrl)

	cfg, err := liteclient.GetConfigFromUrl(t.Context(), networkCfg)
	require.NoError(t, err, "Failed to get config from URL: %s", networkCfg)

	connectionPool := liteclient.NewConnectionPool()
	err = connectionPool.AddConnectionsFromConfig(t.Context(), cfg)
	require.NoError(t, err)

	client := ton.NewAPIClient(connectionPool, ton.ProofCheckPolicyFast)
	client.SetTrustedBlockFromConfig(cfg)

	const readinessRetries = 30
	var lastErr error
	for i := 0; i < readinessRetries; i++ {
		_, lastErr = client.GetMasterchainInfo(t.Context())
		if lastErr == nil {
			break
		}
		t.Logf("API server not ready yet (attempt %d/%d): %v", i+1, readinessRetries, lastErr)
		time.Sleep(time.Second)
	}
	require.NoError(t, lastErr, "TON network not ready after %d attempts", readinessRetries)
	return client
}

func createTonChainConfig(chainID string, chain cldf_ton.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "ton-local"
	chainConfig["NetworkNameFull"] = "ton-local"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
