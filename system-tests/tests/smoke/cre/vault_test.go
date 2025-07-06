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

func TestVault_E2E(t *testing.T) {
	configErr := setCICtfConfigIfMissing("environment-gateway-vault-don.toml")
	require.NoError(t, configErr, "failed to set CTF config")

	c, err := framework.Load[Config](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(c.Blockchain)
	require.NoError(t, err)

	gatewayNodeSet, err := simple_node_set.NewSharedDBNodeSet(c.NodeSets[0], bc)
	require.NoError(t, err)

	// Create gateway job spec for the first nodeset
	gatewayJobSpec := `type = "gateway"
schemaVersion = 1
name = "vault_gateway"
forwardingAllowed = false

[gatewayConfig.ConnectionManagerConfig]
AuthChallengeLen = 10
AuthGatewayId = "vault_gateway"
AuthTimestampToleranceSec = 5
HeartbeatIntervalSec = 20

[gatewayConfig.HTTPClientConfig]
MaxResponseBytes = 100_000_000

[gatewayConfig.NodeServerConfig]
HandshakeTimeoutMillis = 1_000
MaxRequestBytes = 100_000
Path = "/node"
Port = 8080
ReadTimeoutMillis = 1_000
RequestTimeoutMillis = 10_000
WriteTimeoutMillis = 1_000

[gatewayConfig.UserServerConfig]
ContentTypeHeader = "application/jsonrpc"
MaxRequestBytes = 100_000
Path = "/"
Port = 5_002
ReadTimeoutMillis = 1_000
RequestTimeoutMillis = 10_000
WriteTimeoutMillis = 1_000
CORSEnabled = false
CORSAllowedOrigins = []

[[gatewayConfig.Dons]]
DonId = "vault"
HandlerName = "vault"
F = 0

  [gatewayConfig.Dons.HandlerConfig]
  request_timeout_sec = 30

  [[gatewayConfig.Dons.Members]]
  Name = "node_1"
  Address = "0x0000000000000000000000000000000000000001"` // Address is the Eth key of the node (signing key); you can specify the key that's used
	// Gateway URL hardcode in the delegate (or specify config)

	// Vault nodes get a URL of a gateway in config that they reach out to
	// The gateway job allowlists the nodes that can reach out to it

	// Validate and create the gateway job
	_, err = gateway.ValidatedGatewaySpec(gatewayJobSpec)
	require.NoError(t, err)

	gatewayNodeSetClients, err := clclient.New(gatewayNodeSet.CLNodes)
	require.NoError(t, err)

	// Add the gateway job to each node in the first nodeset
	for _, client := range gatewayNodeSetClients {
		job, resp, err := client.CreateJobRaw(gatewayJobSpec)
		require.NoError(t, err, "Gateway job creation request must not error")
		require.Empty(t, job.Errors, "Gateway job creation response must not return any errors")
		require.NotEmpty(t, job.Data.ID, "Gateway job creation response must return a job ID")
		require.Equal(t, http.StatusOK, resp.StatusCode, "Gateway job creation request must return 200 OK")
		fmt.Println(job.Data.ID)
	}

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
			parsedURL.Host = parsedURL.Hostname() + ":5002"
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
			require.Equal(t, "test-secret", response.Result.ID)
			require.Empty(t, response.Result.ErrorMessage)
		}
	})
}
