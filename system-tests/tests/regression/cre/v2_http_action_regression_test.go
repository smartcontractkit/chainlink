package cre

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	httpaction_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/httpaction-negative/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// HTTP Action test cases for failure scenarios
type httpActionFailureTest struct {
	name          string
	testCase      string
	method        string
	url           string
	headers       map[string]string
	body          string
	timeout       int
	expectedError string
}

var httpActionFailureTests = []httpActionFailureTest{
	// Invalid URL tests
	{
		name:          "invalid URL format",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "not-a-valid-url",
		expectedError: "HTTP Action failure test completed: invalid-url-format",
	},
	{
		name:          "non-existing URL",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "http://non-existing-domain-12345.com/api/test",
		expectedError: "HTTP Action failure test completed: non-existing-domain",
	},
	// Invalid method tests
	{
		name:          "invalid HTTP method",
		testCase:      "crud-failure",
		method:        "INVALID",
		url:           "http://host.docker.internal:8080/test",
		expectedError: "HTTP Action failure test completed: invalid-http-method",
	},
	// Invalid headers tests
	{
		name:     "invalid headers",
		testCase: "crud-failure",
		method:   "GET",
		url:      "http://host.docker.internal:8080/test",
		headers: map[string]string{
			"Invalid\nHeader": "value",
		},
		expectedError: "HTTP Action failure test completed: invalid-headers",
	},
	// Invalid body tests
	{
		name:          "corrupt JSON body",
		testCase:      "crud-failure",
		method:        "POST",
		url:           "http://host.docker.internal:8080/test",
		body:          `{"invalid": json}`,
		expectedError: "HTTP Action failure test completed: corrupt-json-body",
	},
	// Size limit tests
	{
		name:          "oversized request body",
		testCase:      "crud-failure",
		method:        "POST",
		url:           "http://host.docker.internal:8080/test",
		body:          strings.Repeat("a", 10*1024*1024), // 10MB body
		expectedError: "HTTP Action failure test completed: oversized-request-body",
	},
	{
		name:          "oversized URL",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "http://host.docker.internal:8080/test?" + strings.Repeat("param=value&", 10000), // Very long URL
		expectedError: "HTTP Action failure test completed: oversized-url",
	},
	// Timeout tests
	{
		name:          "request timeout",
		testCase:      "crud-failure",
		method:        "GET",
		url:           "http://host.docker.internal:8080/delay/10", // Endpoint that delays response
		timeout:       1,                                           // 1 second timeout
		expectedError: "HTTP Action failure test completed: unknown-failure-type",
	},
}

func HTTPActionFailureTest(t *testing.T, testEnv *ttypes.TestEnvironment, httpActionTest httpActionFailureTest) {
	testLogger := framework.L
	const workflowFileLocation = "./httpaction-negative/main.go"

	testLogger.Info().Msg("Creating HTTP Action failure test workflow configuration...")

	workflowConfig := httpaction_config.Config{
		URL:       httpActionTest.url,
		TestCase:  httpActionTest.testCase,
		Method:    httpActionTest.method,
		Headers:   httpActionTest.headers,
		Body:      httpActionTest.body,
		TimeoutMs: httpActionTest.timeout,
	}

	workflowName := "http-action-fail-workflow-" + httpActionTest.method + "-" + uuid.New().String()[0:8]

	// Start Beholder listener BEFORE registering workflow to avoid missing messages
	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, testEnv)

	// Now register and deploy the workflow
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// Wait for specific error message in Beholder based on test case
	testLogger.Info().Msgf("Waiting for expected HTTP Action failure: '%s' in Beholder...", httpActionTest.expectedError)
	timeout := 60 * time.Second

	// Expect exact error message for this test case - no fallbacks
	err := t_helpers.AssertBeholderMessage(listenerCtx, t, httpActionTest.expectedError, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "Expected HTTP Action failure message '%s' not found in Beholder logs", httpActionTest.expectedError)
	testLogger.Info().Msg("HTTP Action failure test completed successfully")

	// Note: Workflow cleanup happens via t.Cleanup() after this function returns
	// The delay below ensures cleanup completes before the next test starts
	testLogger.Info().Msg("Waiting for workflow cleanup before next test...")
	time.Sleep(5 * time.Second)
}
