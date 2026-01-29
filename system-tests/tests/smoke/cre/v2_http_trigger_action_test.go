package cre

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	libcrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	workflowTag = "some-tag"
	requestID   = "http-trigger-validation-test"
)

func ExecuteHTTPTriggerActionTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	publicKeyAddr, signingKey, newKeysErr := libcrypto.GenerateNewKeyPair()
	require.NoError(t, newKeysErr, "failed to generate new public key")

	// Debug: compare the authorized key with the signing key address
	signerAddr := crypto.PubkeyToAddress(signingKey.PublicKey).Hex()
	testLogger.Info().
		Str("config_authorized_key", publicKeyAddr.Hex()).
		Str("signing_key_address", signerAddr).
		Msg("Generated key pair for workflow authorization")

	fakeServer, err := startTestOrderServer(t, testEnv.Config.Fake.Port)
	require.NoError(t, err, "failed to start fake HTTP server")

	uniqueWorkflowName := "http-trigger-action-test-" + uuid.New().String()[0:8]
	httpWorkflowConfig := t_helpers.HTTPWorkflowConfig{
		AuthorizedKey: publicKeyAddr,
		URL:           fakeServer.BaseURLHost,
	}

	pathToWorkflow := "../../../../core/scripts/cre/environment/examples/workflows/v2/http_simple/main.go"
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, uniqueWorkflowName, &httpWorkflowConfig, pathToWorkflow)

	testEnv.Logger.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.Dons.GatewayConnectors.Configurations, "expected at least one gateway configuration")

	gatewayConfig := testEnv.Dons.GatewayConnectors.Configurations[0].Incoming
	testEnv.Logger.Info().
		Str("protocol", gatewayConfig.Protocol).
		Str("host", gatewayConfig.Host).
		Int("port", gatewayConfig.ExternalPort).
		Str("path", gatewayConfig.Path).
		Msg("Gateway configuration details")

	newGatewayURL := gatewayConfig.Protocol + "://" + gatewayConfig.Host + ":" + strconv.Itoa(gatewayConfig.ExternalPort) + gatewayConfig.Path
	gatewayURL, err := url.Parse(newGatewayURL)
	require.NoError(t, err, "failed to parse gateway URL")

	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain type")
	workflowOwner := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient.MustGetRootPrivateKey()
	workflowOwnerAddress := strings.ToLower(crypto.PubkeyToAddress(workflowOwner.PublicKey).Hex())

	testEnv.Logger.Info().
		Str("owner", workflowOwnerAddress).
		Str("workflow_name", uniqueWorkflowName).
		Str("workflow_tag", "some-tag").
		Msg("Waiting for workflow to be loaded before triggering...")

	executeHTTPTriggerRequest(t, testEnv, gatewayURL, uniqueWorkflowName, workflowID, signingKey, workflowOwnerAddress)
	validateHTTPWorkflowRequest(t, testEnv)

	testEnv.Logger.Info().Msg("HTTP trigger and action test completed successfully")
}

// executeHTTPTriggerRequest executes an HTTP trigger request and waits for successful response
func executeHTTPTriggerRequest(t *testing.T, testEnv *ttypes.TestEnvironment, gatewayURL *url.URL, workflowName, workflowID string, singingKey *ecdsa.PrivateKey, workflowOwnerAddress string) {
	var finalResponse jsonrpc.Response[json.RawMessage]
	var triggerRequest jsonrpc.Request[json.RawMessage]

	tick := 5 * time.Second
	require.Eventually(t, func() bool {
		triggerRequest = createHTTPTriggerRequestWithKey(t, workflowName, workflowID, workflowOwnerAddress, singingKey)
		triggerRequestBody, err := json.Marshal(triggerRequest)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to marshal trigger request: %v", err)
			return false
		}

		testEnv.Logger.Info().Msgf("Gateway URL: %s", gatewayURL.String())
		testEnv.Logger.Info().Msg("Executing HTTP trigger request with retries until workflow is loaded...")

		req, err := http.NewRequestWithContext(t.Context(), "POST", gatewayURL.String(), bytes.NewBuffer(triggerRequestBody))
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to create request: %v", err)
			return false
		}
		req.Header.Set("Content-Type", "application/jsonrpc")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to execute request: %v", err)
			return false
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to read response body: %v", err)
			return false
		}

		if resp.StatusCode != http.StatusOK {
			testEnv.Logger.Warn().Msgf("HTTP trigger response (status %d): %s", resp.StatusCode, string(body))
			testEnv.Logger.Warn().Msgf("Gateway returned status %d, retrying...", resp.StatusCode)
			return false
		}

		testEnv.Logger.Info().Msgf("HTTP trigger response (status %d): %s", resp.StatusCode, string(body))

		err = json.Unmarshal(body, &finalResponse)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to unmarshal response: %v", err)
			return false
		}

		if finalResponse.Error != nil {
			testEnv.Logger.Warn().Msgf("JSON-RPC error in response: %v", finalResponse.Error)
			return false
		}

		testEnv.Logger.Info().Msg("Successfully received 200 OK response from gateway")
		return true
	}, 5*time.Minute, tick, "gateway should respond with 200 OK and valid response once workflow is loaded (ensure the workflow is loaded)")

	require.Equal(t, jsonrpc.JsonRpcVersion, finalResponse.Version, "expected JSON-RPC version %s, got %s", jsonrpc.JsonRpcVersion, finalResponse.Version)
	require.Equal(t, triggerRequest.ID, finalResponse.ID, "expected response ID %s, got %s", triggerRequest.ID, finalResponse.ID)
	require.Nil(t, finalResponse.Error, "unexpected error in response: %v", finalResponse.Error)
}

// validateHTTPWorkflowRequest validates that the workflow made the expected HTTP request
func validateHTTPWorkflowRequest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	tick := 5 * time.Second
	require.Eventually(t, func() bool {
		records, err := fake.R.Get("POST", "/orders")
		return err == nil && len(records) > 0
	}, 3*time.Minute, tick, "workflow should have made at least one HTTP request to mock server")

	records, err := fake.R.Get("POST", "/orders")
	require.NoError(t, err, "failed to get recorded requests")
	require.NotEmpty(t, records, "no requests recorded")

	recordedRequest := records[0]
	testEnv.Logger.Info().Msgf("Recorded request: %+v", recordedRequest)

	require.Equal(t, "POST", recordedRequest.Method, "expected POST method")
	require.Equal(t, "/orders", recordedRequest.Path, "expected /orders endpoint")
	require.Equal(t, "application/json", recordedRequest.Headers.Get("Content-Type"), "expected JSON content type")

	var workflowRequestBody map[string]any
	err = json.Unmarshal([]byte(recordedRequest.ReqBody), &workflowRequestBody)
	require.NoError(t, err, "request body should be valid JSON")

	require.Equal(t, "test-customer", workflowRequestBody["customer"], "expected customer field")
	require.Equal(t, "large", workflowRequestBody["size"], "expected size field")
	require.Contains(t, workflowRequestBody, "toppings", "expected toppings field")
}

func createHTTPTriggerRequestWithKey(t *testing.T, workflowName, workflowID, workflowOwner string, privateKey *ecdsa.PrivateKey) jsonrpc.Request[json.RawMessage] {
	signerAddress := strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())

	framework.L.Info().
		Str("workflow_owner", workflowOwner).
		Str("workflow_name", workflowName).
		Str("workflow_tag", workflowTag).
		Str("workflow_id", workflowID).
		Str("signer_address", signerAddress).
		Msg("Creating HTTP trigger request")

	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: workflowOwner,
			WorkflowName:  workflowName,
			WorkflowTag:   workflowTag,
			WorkflowID:    workflowID,
		},
		Input: json.RawMessage(`{
			"customer": "test-customer",
			"size": "large",
			"toppings": ["cheese", "pepperoni"],
			"dedupe": false
		}`),
	}

	payloadBytes, err := json.Marshal(triggerPayload)
	require.NoError(t, err)
	rawPayload := json.RawMessage(payloadBytes)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawPayload,
		ID:      requestID + uuid.New().String()[0:8],
	}

	token, err := utils.CreateRequestJWT(req)
	require.NoError(t, err)

	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)
	req.Auth = tokenString

	return req
}

// startTestOrderServer creates a fake HTTP server that records requests and returns proper responses for order endpoints
func startTestOrderServer(t *testing.T, port int) (*fake.Output, error) {
	fakeInput := &fake.Input{
		Port: port,
	}

	fakeOutput, err := fake.NewFakeDataProvider(fakeInput)
	if err != nil {
		return nil, err
	}

	// Set up the /orders endpoint
	response := map[string]any{
		"orderId": "test-order-" + uuid.New().String()[0:8],
		"status":  "success",
		"message": "Order processed successfully",
	}

	err = fake.JSON("POST", "/orders", response, 200)
	require.NoError(t, err, "failed to set up /orders endpoint")

	framework.L.Info().Msgf("Test order server started on port %d at: %s", port, fakeOutput.BaseURLHost)
	return fakeOutput, nil
}
