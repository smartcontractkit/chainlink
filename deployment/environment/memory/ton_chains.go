package memory

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type TonChain struct {
	Client         *ton.APIClient
	DeployerWallet *wallet.Wallet
}

func getTestTonChainSelectors() []uint64 {
	return []uint64{chainsel.TON_LOCALNET.Selector}
}

// Note: this utility functions can be replaced once we have in the chainlink-ton utils package
func createTonWallet(t *testing.T, client ton.APIClientWrapped, version wallet.Version, option wallet.Option) *wallet.Wallet {
	seed := wallet.NewSeed()
	rw, err := wallet.FromSeed(client, seed, version)
	require.NoError(t, err)
	pw, perr := wallet.FromPrivateKeyWithOptions(client, rw.PrivateKey(), version, option)
	require.NoError(t, perr)
	return pw
}

func fundTonWallets(t *testing.T, client ton.APIClientWrapped, recipients []*address.Address, amounts []tlb.Coins) {
	require.Len(t, amounts, len(recipients), "recipients and amounts must have the same length")
	// initialize the prefunded wallet(Highload-V2), for other wallets, see https://github.com/neodix42/mylocalton-docker#pre-installed-wallets
	version := wallet.HighloadV2Verified //nolint:staticcheck // SA1019: only available option in mylocalton-docker
	rawHlWallet, err := wallet.FromSeed(client, strings.Fields(blockchain.DefaultTonHlWalletMnemonic), version)
	require.NoError(t, err)
	mcFunderWallet, err := wallet.FromPrivateKeyWithOptions(client, rawHlWallet.PrivateKey(), version, wallet.WithWorkchain(-1))
	require.NoError(t, err)
	funder, err := mcFunderWallet.GetSubwallet(uint32(42))
	require.NoError(t, err)
	// double check funder address
	require.Equal(t, blockchain.DefaultTonHlWalletAddress, funder.Address().StringRaw(), "funder address mismatch")
	// create transfer messages for each recipient
	messages := make([]*wallet.Message, len(recipients))
	for i, addr := range recipients {
		transfer, terr := funder.BuildTransfer(addr, amounts[i], false, "")
		require.NoError(t, terr)
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

	bcInput := &blockchain.Input{
		ChainID: strconv.FormatUint(chainID, 10),
		Type:    "ton",
		Port:    strconv.Itoa(freeport.GetOne(t)),
	}
	var bc *blockchain.Output
	// spin up mylocalton with CTFv2
	bc, err := blockchain.NewBlockchainNetwork(bcInput)
	require.NoError(t, err, "Failed to create TON blockchain")

	// get local config from simple http server in genesis node
	cfg, err := liteclient.GetConfigFromUrl(t.Context(), fmt.Sprintf("http://%s/localhost.global.config.json", bc.Nodes[0].ExternalHTTPUrl))
	require.NoError(t, err)

	// establish connection to the TON node
	connectionPool := liteclient.NewConnectionPool()
	err = connectionPool.AddConnectionsFromConfig(t.Context(), cfg)
	require.NoError(t, err)
	client := ton.NewAPIClient(connectionPool, ton.ProofCheckPolicyFast)
	client.SetTrustedBlockFromConfig(cfg)

	// check connection, CTFv2 handles the readiness
	_, err = client.GetMasterchainInfo(t.Context())
	require.NoError(t, err, "TON network not ready")
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
