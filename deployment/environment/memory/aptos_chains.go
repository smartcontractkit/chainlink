package memory

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"

	"github.com/smartcontractkit/freeport"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func getTestAptosChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.APTOS_LOCALNET.Selector}
}

func createAptosAccount(t *testing.T, useDefault bool) *aptos.Account {
	if useDefault {
		addressStr := blockchain.DefaultAptosAccount
		var defaultAddress aptos.AccountAddress
		err := defaultAddress.ParseStringRelaxed(addressStr)
		require.NoError(t, err)

		privateKeyStr := blockchain.DefaultAptosPrivateKey
		privateKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privateKeyStr, "0x"))
		require.NoError(t, err)
		privateKey := ed25519.NewKeyFromSeed(privateKeyBytes)

		t.Logf("Using default Aptos account: %s %+v", addressStr, privateKeyBytes)

		account, err := aptos.NewAccountFromSigner(&crypto.Ed25519PrivateKey{Inner: privateKey}, defaultAddress)
		require.NoError(t, err)
		return account
	} else {
		account, err := aptos.NewEd25519SingleSenderAccount()
		require.NoError(t, err)
		return account
	}
}

func GenerateChainsAptos(t *testing.T, numChains int) map[uint64]deployment.AptosChain {
	testAptosChainSelectors := getTestAptosChainSelectors()
	if len(testAptosChainSelectors) < numChains {
		t.Fatalf("not enough test aptos chain selectors available")
	}
	chains := make(map[uint64]deployment.AptosChain)
	for i := 0; i < numChains; i++ {
		selector := testAptosChainSelectors[i]
		chainID, err := chainsel.GetChainIDFromSelector(selector)
		require.NoError(t, err)
		account := createAptosAccount(t, true)

		url, nodeClient := aptosChain(t, chainID, account.Address)
		chains[selector] = deployment.AptosChain{
			Selector:       selector,
			Client:         nodeClient,
			DeployerSigner: account,
			URL:            url,
			Confirm: func(txHash string, opts ...any) error {
				userTx, err := nodeClient.WaitForTransaction(txHash, opts...)
				if err != nil {
					return err
				}
				if !userTx.Success {
					return fmt.Errorf("transaction failed: %s", userTx.VmStatus)
				}
				return nil
			},
		}
	}
	t.Logf("Created %d Aptos chains: %+v", len(chains), chains)
	return chains
}

func aptosChain(t *testing.T, chainID string, adminAddress aptos.AccountAddress) (string, *aptos.NodeClient) {
	t.Helper()

	// initialize the docker network used by CTF
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	maxRetries := 10
	var url string
	var containerName string
	for i := 0; i < maxRetries; i++ {
		// reserve all the ports we need explicitly to avoid port conflicts in other tests
		ports := freeport.GetN(t, 2)

		bcInput := &blockchain.Input{
			Image:       "", // filled out by defaultAptos function
			Type:        "aptos",
			ChainID:     chainID,
			PublicKey:   adminAddress.String(),
			CustomPorts: []string{fmt.Sprintf("%d:8080", ports[0]), fmt.Sprintf("%d:8081", ports[1])},
		}
		output, err := blockchain.NewBlockchainNetwork(bcInput)
		if err != nil {
			t.Logf("Error creating Aptos network: %v", err)
			freeport.Return(ports)
			time.Sleep(time.Second)
			maxRetries -= 1
			continue
		}
		require.NoError(t, err)
		containerName = output.ContainerName
		testcontainers.CleanupContainer(t, output.Container)
		url = output.Nodes[0].ExternalHTTPUrl + "/v1"
		break
	}

	client, err := aptos.NewNodeClient(url, 0)
	require.NoError(t, err)

	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		_, err := client.GetChainId()
		if err != nil {
			t.Logf("API server not ready yet (attempt %d): %+v\n", i+1, err)
			continue
		}
		ready = true
		break
	}
	require.True(t, ready, "Aptos network not ready")
	time.Sleep(15 * time.Second) // we have slot errors that force retries if the chain is not given enough time to boot

	dc, err := framework.NewDockerClient()
	require.NoError(t, err)
	// incase we didn't use the default account above
	_, err = dc.ExecContainer(containerName, []string{"aptos", "account", "fund-with-faucet", "--account", adminAddress.String(), "--amount", "100000000000"})
	require.NoError(t, err)

	return url, client
}

func createAptosChainConfig(chainID string, chain deployment.AptosChain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "aptos-local"
	chainConfig["NetworkNameFull"] = "aptos-local"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
