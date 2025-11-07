package cre

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"

	httpactionconfig "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	thelpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// HTTP Action test cases for successful CRUD operations
type httpActionSuccessTest struct {
	name       string
	testCase   string
	method     string
	body       string
	endpoint   string
	statusCode int
	url        string
}

var httpActionSuccessTests = []httpActionSuccessTest{
	{
		name:       "POST operation",
		testCase:   "crud-post-success",
		method:     "POST",
		body:       `{"name": "Test Resource", "type": "test"}`,
		endpoint:   "/api/resources/",
		statusCode: 200,
		url:        "",
	},
	{
		name:       "GET operation",
		testCase:   "crud-get-success",
		method:     "GET",
		body:       ``,
		endpoint:   "/api/resources/test-resource-1",
		statusCode: 200,
		url:        "",
	},
	{
		name:       "PUT operation",
		testCase:   "crud-put-success",
		method:     "PUT",
		body:       `{"name": "Updated Test Resource", "type": "test"}`,
		endpoint:   "/api/resources/test-resource-2",
		statusCode: 201,
		url:        "",
	},
	{
		name:       "DELETE operation",
		testCase:   "crud-delete-success",
		method:     "DELETE",
		body:       ``,
		statusCode: 200,
		endpoint:   "/api/resources/test-resource-3",
		url:        "",
	},
}

// ExecuteHTTPActionCRUDSuccessTest executes HTTP Action CRUD operations success test
func ExecuteHTTPActionCRUDSuccessTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	fakeHTTP, err := fake.NewFakeDataProvider(testEnv.Config.FakeHTTP)
	if err != nil {
		testLogger.Error().Err(err).Msg("Failed to start fake HTTP")
	} else {
		testLogger.Info().Msg("Fake HTTP started successfully")
	}

	// Set up a unique endpoint for this test
	response := map[string]any{
		"status": "success",
	}

	for _, testCase := range httpActionSuccessTests {
		err = fake.JSON(testCase.method, testCase.endpoint, response, testCase.statusCode)
		require.NoError(t, err, "failed to set up %s endpoint for %s", testCase.endpoint, testCase.method)
	}

	framework.L.Info().Msgf("Test HTTP server started on port %d at: %s (%s)", testEnv.Config.FakeHTTP.Port, fakeHTTP.BaseURLHost, fakeHTTP.BaseURLDocker)

	defer func() {
		if fakeHTTP != nil {
			testLogger.Info().Msgf("Cleaning up fake server on port %d", testEnv.Config.FakeHTTP.Port)
		}
	}()

	for _, testCase := range httpActionSuccessTests {
		testCase.url = fakeHTTP.BaseURLDocker + testCase.endpoint

		testName := "[v2] HTTP Action " + testCase.name
		t.Run(testName, func(t *testing.T) {
			HTTPActionSuccessTest(t, testEnv, testCase)
		})
	}
}

// HTTPActionSuccessTest executes a single HTTP Action success test case
func HTTPActionSuccessTest(t *testing.T, testEnv *ttypes.TestEnvironment, httpActionTest httpActionSuccessTest) {
	testLogger := framework.L
	const workflowFileLocation = "./httpaction/main.go"

	testLogger.Info().Msg("Creating HTTP Action success test workflow configuration...")

	workflowConfig := httpactionconfig.Config{
		URL:      httpActionTest.url,
		TestCase: httpActionTest.testCase,
		Method:   httpActionTest.method,
		Body:     httpActionTest.body,
	}

	testID := uuid.New().String()[0:8]
	workflowName := "http-action-success-workflow-" + httpActionTest.testCase + "-" + testID
	thelpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// Start Beholder listener to capture workflow execution messages
	listenerCtx, messageChan, kafkaErrChan := thelpers.StartBeholder(t, testLogger, testEnv)

	// Wait for workflow execution to complete and verify success
	testLogger.Info().Msg("Waiting for HTTP Action CRUD operations to complete...")
	timeout := 60 * time.Second

	// Expect exact success message for this test case
	expectedMessage := "HTTP Action CRUD success test completed: " + httpActionTest.testCase
	err := thelpers.AssertBeholderMessage(listenerCtx, t, expectedMessage, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "HTTP Action CRUD success test failed")

	testLogger.Info().Msg("HTTP Action CRUD success test completed")
}
