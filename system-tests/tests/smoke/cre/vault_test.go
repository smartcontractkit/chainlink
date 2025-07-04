package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
)

type Config struct {
	Blockchain *blockchain.Input        `toml:"blockchain" validate:"required"`
	NodeSets   []*simple_node_set.Input `toml:"nodesets" validate:"required"`
}

const VAULT_DON_ID = "vault"
const VAULT_HANDLER_NAME = "vault"
const VAULT_GATEWAY_ID = "vault_gateway"
const VAULT_NODE_1_NAME = "node_1"
const GATEWAY_PORT_FOR_NODES = "18080"
const GATEWAY_PORT_FOR_USERS = "5002"
const NODE_REQUEST_PATH = "/node"

func TestVault_E2E(t *testing.T) {
	// configErr := setCICtfConfigIfMissing("environment-gateway-vault-don.toml")
	// require.NoError(t, configErr, "failed to set CTF config")

	c, err := framework.Load[Config](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(c.Blockchain)
	require.NoError(t, err)

	gatewayNodeSet, err := simple_node_set.NewSharedDBNodeSet(c.NodeSets[0], bc)
	require.NoError(t, err)

	gatewayNodeSetClients, err := clclient.New(gatewayNodeSet.CLNodes)
	require.NoError(t, err)

	// Vault node configuration
	// 1. [Capabilities.GatewayConnector] must include the following:
	// DonID, which must match the DonId in the gateway job spec
	// ChainIDForNodeKey, which must match the ChainID in the gateway and vault job specs
	// NodeAddress, which must match the address of the node key used to sign the gateway job
	//
	// 2. [[Capabilities.GatewayConnector.Gateways]] must include the following:
	// Id, which must match the AuthGatewayId in the gateway job spec
	// URL, which is the WS URL of the gateway node (outputted after the node is configured)

	vaultNodeSetConfig := c.NodeSets[1]
	vaultNodeSet, err := simple_node_set.NewSharedDBNodeSet(vaultNodeSetConfig, bc)
	require.NoError(t, err)

	vaultNodeSetClients, err := clclient.New(vaultNodeSet.CLNodes)
	require.NoError(t, err)

	// Retrieve the ETH addresses of the vault nodes
	var ethAddresses []string
	for _, client := range vaultNodeSetClients {
		nodeEthAddresses, err := client.EthAddresses()
		require.NoError(t, err)
		require.NotEmpty(t, nodeEthAddresses)
		ethAddresses = append(ethAddresses, nodeEthAddresses[0])
	}

	// Update the vault node config to include the gateway connector configuration
	for _, node := range vaultNodeSetConfig.NodeSpecs {

		// Parse the gateway node internal URL to extract the hostname
		parsedURL, err := url.Parse(gatewayNodeSet.CLNodes[0].Node.InternalP2PUrl)
		require.NoError(t, err)
		gatewayUrl := fmt.Sprintf("ws://%s:%s%s", parsedURL.Hostname(), GATEWAY_PORT_FOR_NODES, NODE_REQUEST_PATH)

		node.Node.UserConfigOverrides += fmt.Sprintf(`
		[Capabilities.GatewayConnector]
		DonID = "%s"
		ChainIDForNodeKey = "%s"
		NodeAddress = "%s"

		[[Capabilities.GatewayConnector.Gateways]]
		Id = "%s"
		URL = "%s"
		`,
			VAULT_DON_ID,
			c.Blockchain.ChainID,
			ethAddresses[0],
			VAULT_GATEWAY_ID,
			gatewayUrl,
		)

		fmt.Println("gatewayUrl: ", gatewayUrl)
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

		[[gatewayConfig.Dons.Members]]
		Name = "%s"
		Address = "%s"`,
		VAULT_GATEWAY_ID,
		NODE_REQUEST_PATH,
		GATEWAY_PORT_FOR_NODES,
		GATEWAY_PORT_FOR_USERS,
		VAULT_DON_ID,
		VAULT_HANDLER_NAME,
		VAULT_NODE_1_NAME,
		ethAddresses[0],
	)

	_, err = gateway.ValidatedGatewaySpec(gatewayJobSpec)
	require.NoError(t, err)

	// Add the gateway job to each node in the first nodeset
	for _, client := range gatewayNodeSetClients {
		job, resp, err := client.CreateJobRaw(gatewayJobSpec)
		require.NoError(t, err, "Gateway job creation request must not error")
		require.Empty(t, job.Errors, "Gateway job creation response must not return any errors")
		require.NotEmpty(t, job.Data.ID, fmt.Sprintf("Gateway job creation response must return a job ID: %v.", job))
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
			contractID = "0x0000000000000000000000000000000000000000"
			ocrKeyBundleID = "%s"
			transmitterID = "%s"
			relay = "evm"
			p2pv2Bootstrappers = []
			allowNoBootstrappers = true

			[relayConfig]
            chainID = "%s"
		`, nodeOCRKeyID, nodeTransmitterAddresses[0], c.Blockchain.ChainID)

		job, resp, err := client.CreateJobRaw(vaultJobSpec)
		fmt.Println(job)
		fmt.Println(resp)
		require.NoError(t, err, "Vault job creation request must not error")
		require.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("Vault job creation response must return 200 OK: %v", resp))
		require.NotEmpty(t, job.Data.ID, fmt.Sprintf("Vault job creation response must return a job ID: %v.", job))
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
			secretRequest := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "vault.secrets.create",
				"params": map[string]interface{}{
					"id":    "test-secret",
					"value": "test-secret-value",
				},
				"id":   "1",
				"auth": "jwt-token",
			}

			requestBody, err := json.Marshal(secretRequest)
			require.NoError(t, err)

			// Make HTTP request to gateway endpoint
			parsedURL, err := url.Parse(n.Node.ExternalURL)
			require.NoError(t, err)
			parsedURL.Host = parsedURL.Hostname() + ":" + GATEWAY_PORT_FOR_USERS
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
			var response jsonrpc.Response
			err = json.Unmarshal(body, &response)
			require.NoError(t, err)

			// Verify JSON-RPC response structure
			require.Equal(t, jsonrpc.JsonRpcVersion, response.Version)
			require.Equal(t, "1", response.ID)
			var result vault.SecretsCreateResponse
			err = json.Unmarshal(response.Result, &result)
			require.NoError(t, err)
			require.True(t, result.Success)
			require.Equal(t, "test-secret", result.SecretID)
			require.Empty(t, result.ErrorMessage)
		}
	})
}
