package cre

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	httpactionconfig "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// HTTP Action multi-headers test: workflow asserts response MultiHeaders contain multiple Set-Cookie values.
const (
	multiHeadersTestCase             = "multi-headers"
	multiHeadersSuccessMsg           = "HTTP Action multi-headers test completed"
	multiHeadersRegressionTestCase   = "mh-regression-both" // short to stay under workflow name length limit (64)
	multiHeadersRegressionSuccessMsg = "HTTP Action multi-headers regression completed"
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
	{
		name:       "multi-headers response",
		testCase:   "multi-headers",
		method:     "GET",
		body:       ``,
		statusCode: 200,
		endpoint:   "/api/multi-headers",
		url:        "",
	},
}

// ExecuteHTTPActionRegressionTest runs HTTP Action regression tests (e.g. both Headers and MultiHeaders set rejected).
func ExecuteHTTPActionRegressionTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L

	fakeHTTP, err := fake.NewFakeDataProvider(testEnv.Config.FakeHTTP)
	require.NoError(t, err, "Failed to start fake HTTP")
	testLogger.Info().Msg("Fake HTTP started for regression test")
	defer func() {
		testLogger.Info().Msgf("Cleaning up fake server on port %d", testEnv.Config.FakeHTTP.Port)
	}()

	response := map[string]any{"status": "success"}
	err = fake.JSON("GET", "/api/resources/", response, 200)
	require.NoError(t, err, "failed to set up regression endpoint")

	testLogger.Info().Msgf("Test HTTP server started on port %d at: %s (%s)", testEnv.Config.FakeHTTP.Port, fakeHTTP.BaseURLHost, fakeHTTP.BaseURLDocker)

	t.Run("[v2] HTTP Action multi-headers regression: both Headers and MultiHeaders set rejected", func(t *testing.T) {
		HTTPActionRegressionTest(t, testEnv, fakeHTTP.BaseURLDocker+"/api/resources/")
	})
}

// HTTPActionRegressionTest runs a single HTTP Action regression test (deploy workflow, expect success message).
func HTTPActionRegressionTest(t *testing.T, testEnv *ttypes.TestEnvironment, url string) {
	testLogger := framework.L
	const workflowFileLocation = "./httpaction/main.go"

	workflowConfig := httpactionconfig.Config{
		URL:      url,
		TestCase: multiHeadersRegressionTestCase,
		Method:   "GET",
		Body:     ``,
	}

	testID := uuid.New().String()[0:8]
	workflowName := "http-action-regression-workflow-" + multiHeadersRegressionTestCase + "-" + testID
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	testLogger.Info().Msg("Waiting for HTTP Action regression workflow to complete...")
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, multiHeadersRegressionSuccessMsg, 4*time.Minute)
	testLogger.Info().Msg("HTTP Action regression test completed")
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
		if testCase.testCase == multiHeadersTestCase {
			err = fake.Func("GET", testCase.endpoint, func(c *gin.Context) {
				for name, values := range c.Request.Header {
					for _, value := range values {
						c.Writer.Header().Add(name, value)
					}
				}
				c.Writer.Header().Add("Set-Cookie", "sessionid=multi-e2e-1; Path=/")
				c.Writer.Header().Add("Set-Cookie", "csrf=multi-e2e-2; Path=/")
				c.Writer.Header().Add("Set-Cookie", "pref=multi-e2e-3; Path=/")
				c.JSON(http.StatusOK, response)
			})
			require.NoError(t, err, "failed to set up %s endpoint for %s", testCase.endpoint, testCase.method)
			continue
		}

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
			if parallelEnabled && fanoutEnabled {
				t.Parallel()
			}
			// Each case gets its own per-test execution context to avoid shared-signer nonce collisions
			// while still reusing the shared environment cache (sync.Once) for admin sessions.
			perCaseEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, testEnv.TestConfig)
			HTTPActionSuccessTest(t, perCaseEnv, testCase)
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
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))

	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	// Wait for workflow execution to complete and verify success
	testLogger.Info().Msg("Waiting for HTTP Action CRUD operations to complete...")

	var expectedMessage string
	switch httpActionTest.testCase {
	case multiHeadersTestCase:
		expectedMessage = multiHeadersSuccessMsg
	default:
		expectedMessage = "HTTP Action CRUD success test completed: " + httpActionTest.testCase
	}

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedMessage, 4*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))

	testLogger.Info().Msg("HTTP Action CRUD success test completed")
}
