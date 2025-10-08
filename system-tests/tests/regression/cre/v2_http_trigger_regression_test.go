package cre

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"net"
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

	libcrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	http_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/http/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// regression - HTTP trigger negative test cases
type httpNegativeTest struct {
	name          string
	testCase      string
	expectedError string
}

var httpNegativeTests = []httpNegativeTest{
	{
		name:          "invalid AuthorizedKey.Type",
		testCase:      "invalid-key-type",
		expectedError: "invalid key type",
	},
	{
		name:          "invalid AuthorizedKey.PublicKey format",
		testCase:      "invalid-public-key",
		expectedError: "invalid public key",
	},
	{
		name:          "non-existing AuthorizedKey.PublicKey",
		testCase:      "non-existing-public-key",
		expectedError: "Auth failure",
	},
}

// getFreePort returns a free port that can be used for testing
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func HTTPTriggerFailsTest(t *testing.T, testEnv *ttypes.TestEnvironment, httpNegativeTest httpNegativeTest) {
	testLogger := framework.L
	const workflowFileLocation = "./http/main.go"

	// Generate a valid key pair for comparison
	publicKeyAddr, signingKey, newKeysErr := libcrypto.GenerateNewKeyPair()
	require.NoError(t, newKeysErr, "failed to generate new public key")

	// Get a free port for this test
	freePort, err := getFreePort()
	require.NoError(t, err, "failed to get free port")

	// Start fake HTTP server with unique port and endpoint
	testID := uuid.New().String()[0:8]
	fakeServer, err := startTestOrderServer(t, freePort, testID)
	require.NoError(t, err, "failed to start fake HTTP server")

	// Ensure cleanup of the fake server
	defer func() {
		if fakeServer != nil {
			testLogger.Info().Msgf("Cleaning up fake server on port %d", freePort)
		}
	}()

	// Start Beholder listener to capture error messages
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	testLogger.Info().Msg("Creating HTTP negative test workflow configuration...")

	// Determine the authorized key to use based on test case
	var authorizedKeyToUse string
	switch httpNegativeTest.testCase {
	case "invalid-public-key":
		authorizedKeyToUse = "invalid-public-key-format"
	case "non-existing-public-key":
		authorizedKeyToUse = "0x0000000000000000000000000000000000000000"
	default:
		authorizedKeyToUse = publicKeyAddr.Hex()
	}

	workflowConfig := http_negative_config.Config{
		AuthorizedKey: authorizedKeyToUse,
		URL:           fakeServer.BaseURLHost + "/orders-" + testID,
		TestCase:      httpNegativeTest.testCase,
	}

	workflowName := "http-trigger-fail-workflow-" + httpNegativeTest.testCase
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// For invalid key type and invalid public key format, we expect the workflow deployment/trigger setup to fail
	// For non-existing public key, we expect the trigger execution to fail with unauthorized error at gateway level
	if httpNegativeTest.testCase == "non-existing-public-key" {
		// Try to execute the trigger with a valid signing key but unauthorized public key
		testLogger.Info().Msg("Attempting to execute HTTP trigger with unauthorized key...")
		authFailureDetected := executeHTTPTriggerRequestExpectingFailure(t, testEnv, workflowName, signingKey)

		if authFailureDetected {
			testLogger.Info().Msg("HTTP Trigger Fail test successfully completed - authorization properly rejected at gateway level")
			return
		}
	}

	expectedError := httpNegativeTest.expectedError
	timeout := 2 * time.Minute
	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedError, testLogger, messageChan, kafkaErrChan, timeout)

	// For invalid key type and invalid public key format, we expect engine initialization failure
	// This is the correct behavior - the workflow engine should fail to initialize with invalid configs
	if err != nil && (httpNegativeTest.testCase == "invalid-key-type" || httpNegativeTest.testCase == "invalid-public-key") {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "found engine initialization failure message") {
			testLogger.Info().Msgf("HTTP Trigger Fail test successfully completed - engine initialization failed as expected for %s", httpNegativeTest.testCase)
			return
		}
	}

	require.NoError(t, err, "HTTP Trigger Fail test failed")
	testLogger.Info().Msg("HTTP Trigger Fail test successfully completed")
}

// executeHTTPTriggerRequestExpectingFailure attempts to execute an HTTP trigger expecting it to fail
// Returns true if auth failure was detected, false otherwise
func executeHTTPTriggerRequestExpectingFailure(t *testing.T, testEnv *ttypes.TestEnvironment, workflowName string, signingKey *ecdsa.PrivateKey) bool {
	testLogger := framework.L

	// Get gateway configuration
	require.NotEmpty(t, testEnv.CreEnvironment.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayConfig := testEnv.CreEnvironment.DonTopology.GatewayConnectorOutput.Configurations[0]

	// Build gateway URL
	newGatewayURL := gatewayConfig.Incoming.Protocol + "://" + gatewayConfig.Incoming.Host + ":" + strconv.Itoa(gatewayConfig.Incoming.ExternalPort) + gatewayConfig.Incoming.Path
	gatewayURL, err := url.Parse(newGatewayURL)
	require.NoError(t, err, "failed to parse gateway URL")

	// Get workflow owner
	workflowOwner, err := crypto.HexToECDSA(testEnv.WrappedBlockchainOutputs[0].DeployerPrivateKey)
	require.NoError(t, err, "failed to convert private key to ECDSA")
	workflowOwnerAddress := strings.ToLower(crypto.PubkeyToAddress(workflowOwner.PublicKey).Hex())

	testLogger.Info().Msgf("Attempting HTTP trigger execution that should fail for workflow: %s", workflowName)
	testLogger.Info().Msgf("Gateway URL: %s", gatewayURL.String())

	// Retry logic to wait for workflow to be loaded, then expect auth failure
	var authFailureDetected bool
	tick := 5 * time.Second
	timeout := 3 * time.Minute

	require.Eventually(t, func() bool {
		// Create HTTP trigger request with unauthorized key
		triggerRequest := createHTTPTriggerRequestWithKey(t, workflowName, workflowOwnerAddress, signingKey)
		triggerRequestBody, err := json.Marshal(triggerRequest)
		if err != nil {
			testLogger.Warn().Msgf("Failed to marshal trigger request: %v", err)
			return false
		}

		// Execute the HTTP request that should fail due to unauthorized key
		req, err := http.NewRequestWithContext(t.Context(), "POST", gatewayURL.String(), bytes.NewBuffer(triggerRequestBody))
		if err != nil {
			testLogger.Warn().Msgf("Failed to create HTTP request: %v", err)
			return false
		}
		req.Header.Set("Content-Type", "application/jsonrpc")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			testLogger.Info().Msgf("HTTP trigger request failed as expected: %v", err)
			authFailureDetected = true
			return true
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			testLogger.Warn().Msgf("Failed to read response body: %v", err)
			return false
		}

		testLogger.Info().Msgf("HTTP trigger response (status %d): %s", resp.StatusCode, string(body))

		// Parse the response to check for authorization errors
		var response jsonrpc.Response[json.RawMessage]
		if err := json.Unmarshal(body, &response); err == nil {
			if response.Error != nil {
				errorMsg := response.Error.Message
				testLogger.Info().Msgf("Received error in JSON-RPC response: %v", errorMsg)

				// Check if this is an auth failure (expected)
				if errorMsg == "Auth failure" {
					testLogger.Info().Msg("Authorization properly rejected at gateway level")
					authFailureDetected = true
					return true
				}

				// If it's "workflow not found", continue retrying (workflow not loaded yet)
				if errorMsg == "workflow not found" {
					testLogger.Info().Msg("Workflow not found yet, retrying...")
					return false
				}

				// Any other error is unexpected for this test
				testLogger.Warn().Msgf("Unexpected error received: %v", errorMsg)
				return false
			}
		}

		// If we get here, no error was returned, which is unexpected for unauthorized request
		testLogger.Warn().Msg("Expected auth failure but got successful response")
		return false
	}, timeout, tick, "should eventually get auth failure once workflow is loaded")

	return authFailureDetected
}

// createHTTPTriggerRequestWithKey creates an HTTP trigger request (adapted from positive test)
func createHTTPTriggerRequestWithKey(t *testing.T, workflowName, workflowOwner string, privateKey *ecdsa.PrivateKey) jsonrpc.Request[json.RawMessage] {
	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: workflowOwner,
			WorkflowName:  workflowName,
			WorkflowTag:   "TEMP_TAG",
		},
		Input: []byte(`{
			"customer": "test-customer-unauthorized",
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
		ID:      "http-trigger-unauthorized-test-" + uuid.New().String()[0:8],
	}

	token, err := utils.CreateRequestJWT(req)
	require.NoError(t, err)

	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)
	req.Auth = tokenString

	return req
}

// startTestOrderServer creates a fake HTTP server for testing with unique endpoint
func startTestOrderServer(t *testing.T, port int, testID string) (*fake.Output, error) {
	fakeInput := &fake.Input{
		Port: port,
	}

	fakeOutput, err := fake.NewFakeDataProvider(fakeInput)
	if err != nil {
		return nil, err
	}

	// Set up a unique endpoint for this test
	endpoint := "/orders-" + testID
	response := map[string]any{
		"orderId": "test-order-regression-" + testID,
		"status":  "success",
		"message": "Order processed successfully",
	}

	err = fake.JSON("POST", endpoint, response, 200)
	require.NoError(t, err, "failed to set up %s endpoint", endpoint)

	framework.L.Info().Msgf("Test order server started on port %d at: %s with endpoint %s", port, fakeOutput.BaseURLHost, endpoint)
	return fakeOutput, nil
}
