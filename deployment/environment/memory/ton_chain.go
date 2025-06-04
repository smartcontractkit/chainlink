package memory

import (
	//"crypto/secp256k1"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"

	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

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
		t.Log(ton)
		chains[chainID] = ton
	}
	return chains
}

func tonChain(t *testing.T, chainID uint64) *ton.APIClient {
	t.Helper()
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	bcInput := &blockchain.Input{
		Image:   "ghcr.io/neodix42/mylocalton-docker:latest", // filled out by defaultTon function
		Type:    "ton",
		ChainID: strconv.FormatUint(chainID, 10),
	}

	var networkConfigUrl string
	var containerName string

	maxRetries := 10
	useExistingTonlocalnet := false

	for i := 0; i < maxRetries; i++ {
		if !useExistingTonlocalnet {
			output, err := blockchain.NewBlockchainNetwork(bcInput)
			if err != nil {
				t.Logf("Error creating TON network: %v", err)
				time.Sleep(time.Second)
				maxRetries -= 1
				continue
			}
			require.NoError(t, err)
			containerName = output.ContainerName

			// todo: ctf-configured clean up?
			testcontainers.CleanupContainer(t, output.Container)
			networkConfigUrl = fmt.Sprintf("http://%s/localhost.global.config.json", output.Nodes[0].ExternalHTTPUrl)
		} else {
			networkConfigUrl = fmt.Sprintf("http://%s/localhost.global.config.json", "localhost:8000")
		}
		break
	}
	_ = containerName

	cfg, err := liteclient.GetConfigFromUrl(t.Context(), networkConfigUrl)
	require.NoError(t, err, "Failed to get config from URL: %s", networkConfigUrl)

	connectionPool := liteclient.NewConnectionPool()
	err = connectionPool.AddConnectionsFromConfig(t.Context(), cfg)
	require.NoError(t, err)

	client := ton.NewAPIClient(connectionPool, ton.ProofCheckPolicyFast)
	client.SetTrustedBlockFromConfig(cfg)

	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		_, err := client.GetMasterchainInfo(t.Context())
		require.NoError(t, err)
		if err != nil {
			t.Logf("API server not ready yet (attempt %d): %+v\n", i+1, err)
			continue
		}
		ready = true
		break
	}
	require.True(t, ready, "TON network not ready")

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
