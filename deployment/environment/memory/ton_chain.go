package memory

import (
	//"crypto/secp256k1"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

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

func GenerateChainsTon(t *testing.T, numChains int) map[uint64]cldf_ton.Chain {
	testTonChainSelectors := getTestTonChainSelectors()
	if len(testTonChainSelectors) < numChains {
		t.Fatalf("not enough test ton chain selectors available")
	}
	chains := make(map[uint64]cldf_ton.Chain)
	for i := 0; i < numChains; i++ {
		chainID := testTonChainSelectors[i]

		nodeClient := tonChain(t, chainID)
		t.Logf("[TON-E2E] NodeClient %+v", nodeClient)
		// todo: configurable wallet version, we might need to use Highload wallet for some tests
		// todo: configurable wallet options
		wallet := createTonWallet(t, nodeClient, wallet.V3R2, wallet.WithWorkchain(0))
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
	ctx := context.Background()

	// initialize the docker network used by CTF
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	maxRetries := 10
	var networkConfigUrl string
	var containerName string

	// TODO: SKIP for now, taking too much time, remove when we get enough understanding in test environment
	// wget https://raw.githubusercontent.com/neodix42/mylocalton-docker/refs/heads/main/docker-compose.yaml
	// docker-compose up
	// if existing network error happens, run `docker network rm ton`
	useExistingTonlocalnet := false

	for i := 0; i < maxRetries; i++ {
		bcInput := &blockchain.Input{
			Image:   "ghcr.io/neodix42/mylocalton-docker:latest", // filled out by defaultTon function
			Type:    "ton",
			ChainID: strconv.FormatUint(chainID, 10),
		}

		// TODO: SKIP for now, taking too much time
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

	fmt.Printf("DEBUG: Mylocalton config url: %s\n", networkConfigUrl)

	connectionPool := liteclient.NewConnectionPool()

	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), networkConfigUrl)
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return nil
	}

	// connect to lite servers
	err = connectionPool.AddConnectionsFromConfig(context.Background(), cfg)
	require.NoError(t, err)

	// api client with full proof checks
	client := ton.NewAPIClient(connectionPool, ton.ProofCheckPolicyFast)
	client.SetTrustedBlockFromConfig(cfg)

	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		_, err := client.GetMasterchainInfo(ctx)
		require.NoError(t, err)
		if err != nil {
			t.Logf("API server not ready yet (attempt %d): %+v\n", i+1, err)
			continue
		}
		ready = true
		break
	}
	require.True(t, ready, "TON network not ready")
	time.Sleep(15 * time.Second) // we have slot errors that force retries if the chain is not given enough time to boot

	// TODO(ton): fund transmitter and default wallets
	//_, err = framework.ExecContainer(containerName, []string{"ton", "account", "fund-with-faucet", "--account", adminAddress.String(), "--amount", "100000000000"})
	// require.NoError(t, err)

	return client
}

func createTonChainConfig(chainID string, chain cldf_ton.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "localnet"
	chainConfig["NetworkNameFull"] = "ton-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			// TODO(ton): fill out URL correctly
			"URL": "http://localhost:8000/localhost.global.config.json",
		},
	}

	return chainConfig
}
