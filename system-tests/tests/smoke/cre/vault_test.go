package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
)

type Config struct {
	Blockchain *blockchain.Input        `toml:"blockchain" validate:"required"`
	NodeSets   []*simple_node_set.Input `toml:"nodesets" validate:"required"`
}

const VaultDonID = "vault"
const VaultHandlerName = "vault"
const VaultGatewayID = "vault_gateway"
const VaultNode1Name = "node_1"
const GatewayPortForNodes = "18080"
const GatewayPortForUsers = "5002"
const NodeRequestPath = "/node"

// This key is taken from https://smartcontractkit.github.io/chainlink-testing-framework/framework/components/blockchains/evm.html#test-private-keys
// It couldn't find a way to read keys from the blockchain node output.
const DefaultAnvilPrivateKey = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

func TestVault_E2E(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	configErr := setDefaultConfig("environment-gateway-vault-don.toml")
	require.NoError(t, configErr, "failed to set default CTF config")

	c, err := framework.Load[Config](t)
	require.NoError(t, err)

	bootstrapNodeSetConfig := c.NodeSets[0]
	gatewayNodeSetConfig := c.NodeSets[1]
	vaultNodeSetConfig := c.NodeSets[2]

	// BLOCKCHAIN SETUP - start
	bc, err := blockchain.NewBlockchainNetwork(c.Blockchain)
	require.NoError(t, err)

	// Create seth client for deployments
	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].ExternalWSUrl).
		WithPrivateKeys([]string{DefaultAnvilPrivateKey}).
		// do not check if there's a pending nonce nor check node's health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	require.NoError(t, err)

	chainsConfigs := make([]devenv.ChainConfig, 0)
	chainsConfigs = append(chainsConfigs, devenv.ChainConfig{
		ChainID:   c.Blockchain.ChainID,
		ChainName: sethClient.Cfg.Network.Name,
		ChainType: strings.ToUpper(bc.Family),
		WSRPCs: []devenv.CribRPCs{{
			External: bc.Nodes[0].ExternalWSUrl,
			Internal: bc.Nodes[0].InternalWSUrl,
		}},
		HTTPRPCs: []devenv.CribRPCs{{
			External: bc.Nodes[0].ExternalHTTPUrl,
			Internal: bc.Nodes[0].InternalHTTPUrl,
		}},
		DeployerKey: sethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
	})
	blockChains, err := devenv.NewChains(lggr, chainsConfigs)
	require.NoError(t, err)

	// Create CLDF environment
	cldEnv := &cldf.Environment{
		Logger:            lggr,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         datastore.NewMemoryDataStore().Seal(),
		GetContext:        func() context.Context { return t.Context() },
		BlockChains:       cldf_chain.NewBlockChains(maps.Collect(blockChains.All())),
	}

	registrySel := cldEnv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

	ocrDeploymentOutput, err := changeset.DeployOCR3V2(*cldEnv, &changeset.DeployRequestV2{
		ChainSel: registrySel,
	})
	require.NoError(t, err)

	addrs, err := ocrDeploymentOutput.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, addrs, 1)
	ocr3Addr := addrs[0].Address

	// BLOCKCHAIN SETUP - end

	bootstrapNodeSet, err := simple_node_set.NewSharedDBNodeSet(bootstrapNodeSetConfig, bc)
	require.NoError(t, err)

	bootstrapNodeSetClients, err := clclient.New(bootstrapNodeSet.CLNodes)
	require.NoError(t, err)
	bootstrapNodeClient := bootstrapNodeSetClients[0]

	bootstrapJobSpec := fmt.Sprintf(`
		type = "bootstrap"
		schemaVersion = 1
		name = "Bootstrap"
		forwardingAllowed = false
		contractID = "%s"
		relay = "evm"

		[relayConfig]
		chainID = "%s"
		providerType = "ocr3-capability"
	`, ocr3Addr, c.Blockchain.ChainID)

	job, resp, err := bootstrapNodeClient.CreateJobRaw(bootstrapJobSpec)
	require.NoError(t, err, "Bootstrap job creation request must not error")
	require.Empty(t, job.Errors, "Bootstrap job creation response must not return any errors")
	require.NotEmpty(t, job.Data.ID, "Bootstrap job creation response must return a job ID: %v.", job)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Bootstrap job creation request must return 200 OK")

	keys, err := bootstrapNodeClient.MustReadP2PKeys()
	require.NoError(t, err)

	parsedBootstrapURL, err := url.Parse(bootstrapNodeSet.CLNodes[0].Node.InternalP2PUrl)
	require.NoError(t, err)
	bootstrapP2PLocator := fmt.Sprintf("%s@%s", keys.Data[0].Attributes.PeerID, parsedBootstrapURL.Host)

	gatewayNodeSet, err := simple_node_set.NewSharedDBNodeSet(gatewayNodeSetConfig, bc)
	require.NoError(t, err)

	gatewayNodeSetClients, err := clclient.New(gatewayNodeSet.CLNodes)
	require.NoError(t, err)

	// Vault node job spec:
	// 1. [Capabilities.GatewayConnector] must include the following:
	// DonID, which must match the DonId in the gateway job spec
	// ChainIDForNodeKey, which must match the ChainID in the gateway and vault job specs
	// NodeAddress, which must match the address of the node key used to sign the gateway job
	//
	// 2. [[Capabilities.GatewayConnector.Gateways]] must include the following:
	// Id, which must match the AuthGatewayId in the gateway job spec
	// URL, which is the WS URL of the gateway node (outputted after the node is configured)

	vaultNodeSet, err := simple_node_set.NewSharedDBNodeSet(vaultNodeSetConfig, bc)
	require.NoError(t, err)

	vaultNodeSetClients, err := clclient.New(vaultNodeSet.CLNodes)
	require.NoError(t, err)

	// Retrieve the ETH addresses of the vault nodes
	ethAddresses := []string{}
	for _, client := range vaultNodeSetClients {
		nodeEthAddresses, err2 := client.EthAddresses()
		require.NoError(t, err2)
		require.NotEmpty(t, nodeEthAddresses)
		ethAddresses = append(ethAddresses, nodeEthAddresses[0])
	}

	// Update the vault node config to include the gateway connector configuration
	for _, node := range vaultNodeSetConfig.NodeSpecs {
		parsedURL, err2 := url.Parse(gatewayNodeSet.CLNodes[0].Node.InternalP2PUrl)
		require.NoError(t, err2)
		internalGatewayURL := fmt.Sprintf("ws://%s:%s%s", parsedURL.Hostname(), GatewayPortForNodes, NodeRequestPath)
		node.Node.UserConfigOverrides += fmt.Sprintf(`
		[Feature]
		LogPoller = true

		[Capabilities.GatewayConnector]
		DonID = "%s"
		ChainIDForNodeKey = "%s"
		NodeAddress = "%s"

		[[Capabilities.GatewayConnector.Gateways]]
		Id = "%s"
		URL = "%s"
		`,
			VaultDonID,
			c.Blockchain.ChainID,
			ethAddresses[0],
			VaultGatewayID,
			internalGatewayURL,
		)
	}

	vaultNodeSet, err = simple_node_set.UpgradeNodeSet(t, vaultNodeSetConfig, bc, 3*time.Second)
	require.NoError(t, err)
	vaultNodeSetClients, err = clclient.New(vaultNodeSet.CLNodes)
	require.NoError(t, err)

	// Create gateway job spec for the first nodeset
	gatewayJobSpec := fmt.Sprintf(`type = "gateway"
		schemaVersion = 1
		name = "gateway_node"
		forwardingAllowed = false

		[gatewayConfig.ConnectionManagerConfig]
		AuthChallengeLen = 10
		AuthGatewayId = "%s"
		AuthTimestampToleranceSec = 5
		HeartbeatIntervalSec = 20

		[gatewayConfig.HTTPClientConfig]
		MaxResponseBytes = 100_000_000

		[gatewayConfig.NodeServerConfig]
		HandshakeTimeoutMillis = 1_000
		MaxRequestBytes = 100_000
		Path = "%s"
		Port = %s
		ReadTimeoutMillis = 1_000
		RequestTimeoutMillis = 10_000
		WriteTimeoutMillis = 1_000

		[gatewayConfig.UserServerConfig]
		ContentTypeHeader = "application/jsonrpc"
		MaxRequestBytes = 100_000
		Path = "/"
		Port = %s
		ReadTimeoutMillis = 1_000
		RequestTimeoutMillis = 10_000
		WriteTimeoutMillis = 1_000
		CORSEnabled = false
		CORSAllowedOrigins = []

		[[gatewayConfig.Dons]]
		DonId = "%s"
		HandlerName = "%s"
		F = 0

		[gatewayConfig.Dons.HandlerConfig]
		request_timeout_sec = 30
		node_rate_limiter = {
			globalRPS = 100,
			globalBurst = 100,
			perSenderRPS = 10,
			perSenderBurst = 10
		}

		[[gatewayConfig.Dons.Members]]
		Name = "%s"
		Address = "%s"`,
		VaultGatewayID,
		NodeRequestPath,
		GatewayPortForNodes,
		GatewayPortForUsers,
		VaultDonID,
		VaultHandlerName,
		VaultNode1Name,
		ethAddresses[0],
	)

	_, err = gateway.ValidatedGatewaySpec(gatewayJobSpec)
	require.NoError(t, err)

	// Add the gateway job to each node in the first nodeset
	for _, client := range gatewayNodeSetClients {
		job, resp, err := client.CreateJobRaw(gatewayJobSpec)
		require.NoError(t, err, "Gateway job creation request must not error")
		require.Empty(t, job.Errors, "Gateway job creation response must not return any errors")
		require.NotEmpty(t, job.Data.ID, "Gateway job creation response must return a job ID: %v.", job)
		require.Equal(t, http.StatusOK, resp.StatusCode, "Gateway job creation request must return 200 OK")
	}
	fmt.Println("✅ Gateway jobs created successfully.")

	// Add the vault job to each node in the second nodeset
	for _, client := range vaultNodeSetClients {
		// Get the actual OCR key bundle ID and transmitter address for this node
		nodeTransmitterAddresses, err := client.EthAddresses()
		require.NoError(t, err, "Should be able to get ETH addresses from vault node")
		require.NotEmpty(t, nodeTransmitterAddresses, "Vault node should have at least one ETH address")

		nodeOCRKeys, err := client.MustReadOCR2Keys()
		require.NoError(t, err, "Should be able to get OCR2 keys from vault node")

		var nodeOCRKeyID string
		for _, key := range nodeOCRKeys.Data {
			if key.Attributes.ChainType == "evm" {
				nodeOCRKeyID = key.ID
				break
			}
		}
		require.NotEmpty(t, nodeOCRKeyID, "Vault node should have an EVM OCR2 key")

		// Create vault job spec without relayConfig since EVM configuration is provided by node boot config
		vaultJobSpec := fmt.Sprintf(`type = "offchainreporting2"
			schemaVersion = 1
			pluginType = "vault-plugin"
			name = "vault_node"
			forwardingAllowed = false
			maxTaskDuration = "30s"
			contractID = "%s"
			ocrKeyBundleID = "%s"
			transmitterID = "%s"
			relay = "evm"
			p2pv2Bootstrappers = ["%s"]

			[relayConfig]
            chainID = "%s"

			[pluginConfig]
			requestExpiryDuration = "60s"
		`, ocr3Addr, nodeOCRKeyID, nodeTransmitterAddresses[0], bootstrapP2PLocator, c.Blockchain.ChainID)

		job, resp, err := client.CreateJobRaw(vaultJobSpec)
		require.NoError(t, err, "Vault job creation request must not error")
		require.Equal(t, http.StatusOK, resp.StatusCode, "Vault job creation response must return 200 OK: %v", resp)
		require.NotEmpty(t, job.Data.ID, "Vault job creation response must return a job ID: %v.", job)
		fmt.Println(job.Data.ID)
	}
	fmt.Println("✅ Vault jobs created successfully.")

	fmt.Println("⏳ Waiting for a connection between gateway and vault to be established...")
	// TODO: Make this more robust
	time.Sleep(15 * time.Second)
	fmt.Println("Proceeding to test...")

	t.Run("vault secrets create", func(t *testing.T) {
		for _, n := range gatewayNodeSet.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
			require.NotEmpty(t, n.Node.InternalP2PUrl)

			// Prepare the JSON-RPC request to create a secret
			secretsRequest := jsonrpc.Request[vault.SecretsCreateRequest]{
				Version: jsonrpc.JsonRpcVersion,
				Method:  vault.MethodSecretsCreate,
				Params: &vault.SecretsCreateRequest{
					ID:    "test-secret",
					Value: "test-secret-value",
				},
				ID: "1",
			}
			requestBody, err := json.Marshal(secretsRequest)
			require.NoError(t, err)

			// Make HTTP request to gateway endpoint
			parsedURL, err := url.Parse(n.Node.ExternalURL)
			require.NoError(t, err)
			parsedURL.Host = parsedURL.Hostname() + ":" + GatewayPortForUsers
			gatewayURL := parsedURL.String() + "/"
			req, err := http.NewRequestWithContext(context.Background(), "POST", gatewayURL, bytes.NewBuffer(requestBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/jsonrpc")
			req.Header.Set("Accept", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Print response body
			body, err := io.ReadAll(resp.Body)
			fmt.Println("Response Body:", string(body))
			require.NoError(t, err)

			// Check response status
			require.Equal(t, http.StatusOK, resp.StatusCode, "Gateway endpoint should respond with 200 OK")

			// Parse response
			var response jsonrpc.Response[vault.SecretsCreateResponse]
			err = json.Unmarshal(body, &response)
			require.NoError(t, err)

			// Verify JSON-RPC response structure
			require.Equal(t, jsonrpc.JsonRpcVersion, response.Version)
			require.Equal(t, "1", response.ID)
			require.NoError(t, err)
			require.True(t, response.Result.Success)
			require.Equal(t, "test-secret", response.Result.SecretID)
			require.Empty(t, response.Result.ErrorMessage)
		}
	})
}

func setDefaultConfig(configName string) error {
	if os.Getenv("CTF_CONFIGS") == "" {
		return os.Setenv("CTF_CONFIGS", configName)
	}

	return nil
}
