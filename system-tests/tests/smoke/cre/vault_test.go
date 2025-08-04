package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
)

func TestVault_SecretsCreate(t *testing.T) {
	// TODO this cache file needs to be an env var maybe with a default (so that it can also be used in CI, where paths are different)
	// TODO add to other tests
	createErr := createEnvironmentIfNotExists("../../../../core/scripts/cre/environment/configs/workflow-don-cache.toml", "workflow")
	require.NoError(t, createErr, "failed to create environment")

	confErr := setConfigurationIfMissing("../../../../core/scripts/cre/environment/configs/workflow-don-cache.toml")
	require.NoError(t, confErr, "failed to set configuration")

	/*
		LOAD ENVIRONMENT STATE
	*/
	in, err := framework.Load[environment.Config](nil)
	require.NoError(t, err, "couldn't load environment state")
	validateEnvVars(t)
	require.Len(t, in.NodeSets, 1, "expected 1 nodeset in the environment")

	var envArtifact environment.EnvArtifact
	artFile, err := os.ReadFile(os.Getenv("ENV_ARTIFACT_PATH"))
	require.NoError(t, err, "failed to read artifact file")
	err = json.Unmarshal(artFile, &envArtifact)
	require.NoError(t, err, "failed to unmarshal artifact file")

	fullCldEnvOutput, _, loadErr := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in, envArtifact)
	require.NoError(t, loadErr, "failed to load environment")

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
	gatewayConfiguration, handlerErr := credon.GatewayConfigurationForHandler(coregateway.VaultHandlerType, fullCldEnvOutput.DonTopology.GatewayConnectorOutput)
	require.NoError(t, handlerErr)
	gatewayURL, err := url.Parse(gatewayConfiguration.Incoming.Protocol + "://" + gatewayConfiguration.Incoming.Host + ":" + fmt.Sprint(gatewayConfiguration.Incoming.ExternalPort) + gatewayConfiguration.Incoming.Path)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(context.Background(), "POST", gatewayURL.String(), bytes.NewBuffer(requestBody))
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
