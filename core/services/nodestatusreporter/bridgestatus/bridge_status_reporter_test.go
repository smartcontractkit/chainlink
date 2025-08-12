package bridgestatus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	bridgeMocks "github.com/smartcontractkit/chainlink/v2/core/bridges/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	jobMocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/bridgestatus/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/bridgestatus/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Test constants and fixtures
const (
	testStatusPath      = "/status"
	testPollingInterval = 5 * time.Minute
	testBridgeName1     = "bridge1"
	testBridgeName2     = "bridge2"
	testBridgeURL1      = "http://bridge1.example.com"
	testBridgeURL2      = "http://bridge2.example.com"
)

// loadFixture loads a JSON fixture file
func loadFixture(t *testing.T, filename string) string {
	t.Helper()

	fixturePath := filepath.Join("fixtures", filename)
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err, "Failed to read fixture file: %s", fixturePath)

	return string(data)
}

// loadFixtureAsEAResponse loads and unmarshals fixture data
func loadFixtureAsEAResponse(t *testing.T, filename string) EAResponse {
	fixtureData := loadFixture(t, filename)

	var status EAResponse
	err := json.Unmarshal([]byte(fixtureData), &status)
	require.NoError(t, err, "Failed to unmarshal test fixture")

	return status
}

// parseWebURL creates WebURL from string
func parseWebURL(s string) models.WebURL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return models.WebURL(*u)
}

// Test fixtures
var (
	testBridge1 = bridges.BridgeType{
		Name: bridges.MustParseBridgeName(testBridgeName1),
		URL:  parseWebURL(testBridgeURL1),
	}
	testBridge2 = bridges.BridgeType{
		Name: bridges.MustParseBridgeName(testBridgeName2),
		URL:  parseWebURL(testBridgeURL2),
	}
	testBridges = []bridges.BridgeType{testBridge1, testBridge2}
)

// setupTestService creates a test service with mocks
func setupTestService(t *testing.T, enabled bool, pollingInterval time.Duration, httpClient *http.Client) (*Service, *bridgeMocks.ORM, *jobMocks.ORM, *mocks.BeholderEmitter) {
	t.Helper()

	bridgeStatusConfig := mocks.NewTestBridgeStatusReporterConfig(enabled, testStatusPath, pollingInterval)

	bridgeORM := bridgeMocks.NewORM(t)
	jobORM := jobMocks.NewORM(t)
	emitter := mocks.NewBeholderEmitter()
	lggr := logger.TestLogger(t)

	// Reduce log noise
	lggr.SetLogLevel(zapcore.ErrorLevel)

	service := NewBridgeStatusReporter(bridgeStatusConfig, bridgeORM, jobORM, httpClient, emitter, lggr)

	return service, bridgeORM, jobORM, emitter
}

// setupTestServiceWithIgnoreFlags creates a test service with custom ignore flag settings
func setupTestServiceWithIgnoreFlags(t *testing.T, enabled bool, pollingInterval time.Duration, httpClient *http.Client, ignoreInvalidBridges, ignoreJoblessBridges bool) (*Service, *bridgeMocks.ORM, *jobMocks.ORM, *mocks.BeholderEmitter) {
	t.Helper()

	bridgeStatusConfig := mocks.NewTestBridgeStatusReporterConfigWithSkip(enabled, testStatusPath, pollingInterval, ignoreInvalidBridges, ignoreJoblessBridges)

	bridgeORM := bridgeMocks.NewORM(t)
	jobORM := jobMocks.NewORM(t)
	emitter := mocks.NewBeholderEmitter()
	lggr := logger.TestLogger(t)

	// Reduce log noise
	lggr.SetLogLevel(zapcore.ErrorLevel)

	service := NewBridgeStatusReporter(bridgeStatusConfig, bridgeORM, jobORM, httpClient, emitter, lggr)

	return service, bridgeORM, jobORM, emitter
}

func TestNewBridgeStatusReporter(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	assert.NotNil(t, service)
	assert.Equal(t, ServiceName, service.Name())
}

func TestService_Start_Disabled(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _, _ := setupTestService(t, false, testPollingInterval, httpClient)

	ctx := context.Background()
	err := service.Start(ctx)
	require.NoError(t, err)

	err = service.Close()
	require.NoError(t, err)
}

func TestService_Start_Enabled(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, jobORM, emitter := setupTestService(t, true, 100*time.Millisecond, httpClient)

	// Mock the calls that will be triggered by the polling ticker
	bridgeORM.On("BridgeTypes", mock.Anything, mock.Anything, mock.Anything).Return([]bridges.BridgeType{}, 0, nil).Maybe()
	jobORM.On("FindJobIDsWithBridge", mock.Anything, mock.AnythingOfType("string")).Return([]int32{}, nil).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()
	err := service.Start(ctx)
	require.NoError(t, err)

	err = service.Close()
	require.NoError(t, err)
}

func TestService_HealthReport(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	health := service.HealthReport()
	assert.Contains(t, health, service.Name())
}

func TestService_pollAllBridges_NoBridges(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, nil)

	ctx := context.Background()

	// Should handle no bridges gracefully
	assert.NotPanics(t, func() {
		service.pollAllBridges(ctx)
	})

	bridgeORM.AssertExpectations(t)
}

func TestService_pollAllBridges_WithBridges(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "bridge_status_response.json"), http.StatusOK)
	service, bridgeORM, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return(testBridges, len(testBridges), nil)

	// Mock job ORM calls for finding external job IDs
	jobORM.On("FindJobIDsWithBridge", mock.Anything, mock.AnythingOfType("string")).Return([]int32{}, nil)

	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollAllBridges(ctx)

	bridgeORM.AssertExpectations(t)
	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollAllBridges_FetchError(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, assert.AnError)

	ctx := context.Background()

	// Should handle bridge ORM error gracefully
	assert.NotPanics(t, func() {
		service.pollAllBridges(ctx)
	})

	bridgeORM.AssertExpectations(t)
}

func TestService_pollBridge_Success(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "bridge_status_response.json"), http.StatusOK)
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM calls for finding external job IDs
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)

	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "test-bridge", "http://example.com")

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_HTTPError(t *testing.T) {
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM call that now happens at the start of pollBridge
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)

	ctx := context.Background()

	// Should handle HTTP error gracefully
	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://invalid.invalid:8080")
	})

	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)
	jobORM.AssertExpectations(t)
}

func TestService_pollBridge_InvalidJSON(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("invalid json", http.StatusOK)
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM call that now happens at the start of pollBridge
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)

	ctx := context.Background()

	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://invalid.invalid:8080")
	})
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)

	jobORM.AssertExpectations(t)
}

func TestService_pollBridge_InvalidURL(t *testing.T) {
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM call that now happens at the start of pollBridge
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)

	ctx := context.Background()

	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "://invalid-url")
	})

	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)

	jobORM.AssertExpectations(t)
}

func TestService_pollBridge_EmptyURL(t *testing.T) {
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM call that now happens at the start of pollBridge
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)

	ctx := context.Background()

	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "")
	})

	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)

	jobORM.AssertExpectations(t)
}

func TestService_pollBridge_URLPathPreservation(t *testing.T) {
	// Create mock that expects the exact URL with preserved path + status
	httpClient := mocks.NewMockHTTPClientWithExpectedURL(
		loadFixture(t, "bridge_status_response.json"),
		http.StatusOK,
		"http://localhost:8080/bridge/v1/status",
	)
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()

	// If URL path joining is broken, this will fail with "unexpected URL" error
	service.pollBridge(ctx, "test-bridge", "http://localhost:8080/bridge/v1")

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_Non200Status(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("Not Found", http.StatusNotFound)
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM call that now happens at the start of pollBridge
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)

	ctx := context.Background()

	assert.NotPanics(t, func() {
		service.pollBridge(ctx, "test-bridge", "http://example.com")
	})

	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)
	jobORM.AssertExpectations(t)
}

func TestService_emitBridgeStatus_Success(t *testing.T) {
	httpClient := &http.Client{}
	service, _, _, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.emitBridgeStatus(ctx, "test-bridge", loadFixtureAsEAResponse(t, "bridge_status_response.json"), []JobInfo{})

	emitter.AssertExpectations(t)
}

func TestService_pollAllBridges_RefreshError(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, _, _ := setupTestService(t, true, testPollingInterval, httpClient)

	// Setup bridge ORM mock to return error
	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return([]bridges.BridgeType{}, 0, assert.AnError)

	ctx := context.Background()

	// Should handle bridge refresh error gracefully (no panic)
	assert.NotPanics(t, func() {
		service.pollAllBridges(ctx)
	})

	bridgeORM.AssertExpectations(t)
}

func TestService_pollAllBridges_MultipleBridges(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "bridge_status_response.json"), http.StatusOK)
	service, bridgeORM, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Setup bridge ORM mock to return our test bridges
	bridgeORM.On("BridgeTypes", mock.Anything, 0, 1000).Return(testBridges, len(testBridges), nil)

	// Mock job ORM calls for finding external job IDs
	jobORM.On("FindJobIDsWithBridge", mock.Anything, mock.AnythingOfType("string")).Return([]int32{}, nil)

	// Track emitted bridge names from protobuf events
	var emittedBridgeNamesMutex sync.Mutex
	emittedBridgeNames := []string{}

	// Setup emitter mock to capture protobuf events and extract bridge names
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Unmarshal protobuf to extract bridge name
		protobufBytes := args.Get(1).([]byte)
		var event events.BridgeStatusEvent
		if err := proto.Unmarshal(protobufBytes, &event); err == nil {
			emittedBridgeNamesMutex.Lock()
			emittedBridgeNames = append(emittedBridgeNames, event.BridgeName)
			emittedBridgeNamesMutex.Unlock()
		}
	})

	ctx := context.Background()
	service.pollAllBridges(ctx)

	bridgeORM.AssertExpectations(t)
	jobORM.AssertExpectations(t)

	// Verify we emitted events for both bridges
	expectedBridgeNames := []string{testBridgeName1, testBridgeName2}
	assert.ElementsMatch(t, expectedBridgeNames, emittedBridgeNames, "Should emit telemetry for each bridge")

	emitter.AssertExpectations(t)
}

func TestService_emitBridgeStatus_CaptureOutput(t *testing.T) {
	emitter := mocks.NewBeholderEmitter()
	var capturedProtobufBytes []byte

	// Capture protobuf metadata labels
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedProtobufBytes = args.Get(1).([]byte)
	})

	config := mocks.NewTestBridgeStatusReporterConfig(true, "/status", 5*time.Minute)
	service := NewBridgeStatusReporter(
		config,
		nil, // bridgeORM not needed for this test
		nil, // jobORM not needed for this test
		nil, // httpClient not needed for this test
		emitter,
		logger.TestLogger(t),
	)

	// Load fixture and emit
	ctx := context.Background()
	status := loadFixtureAsEAResponse(t, "bridge_status_response.json")
	service.emitBridgeStatus(ctx, "test-bridge", status, []JobInfo{})

	// Unmarshal and verify protobuf matches fixture values
	require.NotEmpty(t, capturedProtobufBytes)
	var event events.BridgeStatusEvent
	err := proto.Unmarshal(capturedProtobufBytes, &event)
	require.NoError(t, err)

	// Verify key fields match fixture
	assert.Equal(t, "test-bridge", event.BridgeName)
	assert.Equal(t, status.Adapter.Name, event.AdapterName)
	assert.Equal(t, status.Adapter.Version, event.AdapterVersion)
	assert.InDelta(t, status.Adapter.UptimeSeconds, event.AdapterUptimeSeconds, 0.001)

	// Verify Endpoints
	for i, endpoint := range status.Endpoints {
		assert.Equal(t, endpoint.Name, event.Endpoints[i].Name)
		assert.Equal(t, endpoint.Aliases, event.Endpoints[i].Aliases)
		assert.Equal(t, endpoint.Transports, event.Endpoints[i].Transports)
	}

	// Verify Default Endpoint
	assert.Equal(t, status.DefaultEndpoint, event.DefaultEndpoint)

	// Verify configuration
	// Helper function to safely convert values to strings, handling nil (same as in production code)
	safeString := func(v any) string {
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%v", v)
	}

	for i, configuration := range status.Configuration {
		assert.Equal(t, configuration.Name, event.Configuration[i].Name)
		assert.Equal(t, safeString(configuration.Value), event.Configuration[i].Value) // Values are converted to strings
		assert.Equal(t, configuration.Type, event.Configuration[i].Type)
		assert.Equal(t, configuration.Description, event.Configuration[i].Description)
		assert.Equal(t, configuration.Required, event.Configuration[i].Required)
		assert.Equal(t, safeString(configuration.Default), event.Configuration[i].DefaultValue) // Defaults converted to strings
		assert.Equal(t, configuration.CustomSetting, event.Configuration[i].CustomSetting)
		assert.Equal(t, safeString(configuration.EnvDefaultOverride), event.Configuration[i].EnvDefaultOverride) // Overrides converted to strings
	}

	// Verify Runtime
	assert.Equal(t, status.Runtime.NodeVersion, event.Runtime.NodeVersion)
	assert.Equal(t, status.Runtime.Platform, event.Runtime.Platform)
	assert.Equal(t, status.Runtime.Architecture, event.Runtime.Architecture)
	assert.Equal(t, status.Runtime.Hostname, event.Runtime.Hostname)

	// Verify Metrics
	assert.Equal(t, status.Metrics.Enabled, event.Metrics.Enabled)

	emitter.AssertExpectations(t)
}

func TestService_Start_AlreadyStarted(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock the calls that will be triggered by the polling ticker
	bridgeORM.On("BridgeTypes", mock.Anything, mock.Anything, mock.Anything).Return([]bridges.BridgeType{}, 0, nil).Maybe()
	jobORM.On("FindJobIDsWithBridge", mock.Anything, mock.AnythingOfType("string")).Return([]int32{}, nil).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	err := service.Start(ctx)
	require.NoError(t, err)
	err = service.Start(ctx)
	// services.StateMachine prevents double start, should return error
	require.Error(t, err)

	// Clean up
	err = service.Close()
	require.NoError(t, err)
}

func TestService_Close_AlreadyClosed(t *testing.T) {
	httpClient := &http.Client{}
	service, bridgeORM, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock the calls that will be triggered by the polling ticker
	bridgeORM.On("BridgeTypes", mock.Anything, mock.Anything, mock.Anything).Return([]bridges.BridgeType{}, 0, nil).Maybe()
	jobORM.On("FindJobIDsWithBridge", mock.Anything, mock.AnythingOfType("string")).Return([]int32{}, nil).Maybe()
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := context.Background()

	err := service.Start(ctx)
	require.NoError(t, err)

	err = service.Close()
	require.NoError(t, err)
	err = service.Close()

	// services.StateMachine prevents double close, should return error
	require.Error(t, err)
}

func TestService_PollAllBridges_3000Bridges(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "bridge_status_response.json"), http.StatusOK)
	service, mockORM, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	numBridges := 3000
	var allBridges []bridges.BridgeType
	for i := 0; i < numBridges; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://bridge%d.example.com", i))
		bridge := bridges.BridgeType{
			Name: bridges.MustParseBridgeName(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		}
		allBridges = append(allBridges, bridge)
	}

	// Page 1: bridges 0-999 (1000 bridges)
	page1 := allBridges[0:1000]
	mockORM.On("BridgeTypes", mock.Anything, 0, bridgePollPageSize).Return(page1, 3000, nil).Once()

	// Page 2: bridges 1000-1999 (1000 bridges)
	page2 := allBridges[1000:2000]
	mockORM.On("BridgeTypes", mock.Anything, 1000, bridgePollPageSize).Return(page2, 3000, nil).Once()

	// Page 3: bridges 2000-2999 (1000 bridges)
	page3 := allBridges[2000:3000]
	mockORM.On("BridgeTypes", mock.Anything, 2000, bridgePollPageSize).Return(page3, 3000, nil).Once()

	// Page 4: empty (end of results)
	mockORM.On("BridgeTypes", mock.Anything, 3000, bridgePollPageSize).Return([]bridges.BridgeType{}, 3000, nil).Once()

	// Mock job ORM calls for finding external job IDs for all bridges
	jobORM.On("FindJobIDsWithBridge", mock.Anything, mock.AnythingOfType("string")).Return([]int32{}, nil).Times(numBridges)

	// Expect 3000 telemetry emissions
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(numBridges)

	ctx := context.Background()

	service.pollAllBridges(ctx)
	mockORM.AssertExpectations(t)
	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_PollAllBridges_ContextTimeout(t *testing.T) {
	httpClient := &http.Client{}
	service, mockORM, jobORM, _ := setupTestService(t, true, testPollingInterval, httpClient)

	numBridges := 5
	var allBridges []bridges.BridgeType
	for i := 0; i < numBridges; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://bridge%d.example.com", i))
		bridge := bridges.BridgeType{
			Name: bridges.MustParseBridgeName(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		}
		allBridges = append(allBridges, bridge)
	}

	mockORM.On("BridgeTypes", mock.Anything, 0, bridgePollPageSize).Return(allBridges, numBridges, nil).Once()

	// Mock job ORM calls for each bridge
	for i := 0; i < numBridges; i++ {
		bridgeName := fmt.Sprintf("bridge%d", i)
		jobORM.On("FindJobIDsWithBridge", mock.Anything, bridgeName).Return([]int32{}, nil)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "HTTP handler should not complete due to context cancellation")
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	for i := range allBridges {
		allBridges[i].URL = models.WebURL(*serverURL)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	service.pollAllBridges(ctx)
	mockORM.AssertExpectations(t)
}

func TestService_emitBridgeStatus_EmptyFields(t *testing.T) {
	emitter := mocks.NewBeholderEmitter()
	var capturedProtobufBytes []byte

	// Capture protobuf metadata labels
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedProtobufBytes = args.Get(1).([]byte)
	})

	config := mocks.NewTestBridgeStatusReporterConfig(true, "/status", 5*time.Minute)
	service := NewBridgeStatusReporter(
		config,
		nil, // bridgeORM not needed for this test
		nil, // jobORM not needed for this test
		nil, // httpClient not needed for this test
		emitter,
		logger.TestLogger(t),
	)

	// Load empty fixture and emit
	ctx := context.Background()
	status := loadFixtureAsEAResponse(t, "bridge_status_empty.json")
	service.emitBridgeStatus(ctx, "empty-bridge", status, []JobInfo{})

	// Unmarshal and verify protobuf handles empty values correctly
	require.NotEmpty(t, capturedProtobufBytes)
	var event events.BridgeStatusEvent
	err := proto.Unmarshal(capturedProtobufBytes, &event)
	require.NoError(t, err)

	// Verify empty/minimal values are handled correctly
	assert.Equal(t, "empty-bridge", event.BridgeName)
	assert.Empty(t, event.AdapterName)
	assert.Empty(t, event.AdapterVersion)
	assert.InDelta(t, float64(0), event.AdapterUptimeSeconds, 0.001)
	assert.Empty(t, event.DefaultEndpoint)

	// Verify empty runtime info
	require.NotNil(t, event.Runtime)
	assert.Empty(t, event.Runtime.NodeVersion)
	assert.Empty(t, event.Runtime.Platform)
	assert.Empty(t, event.Runtime.Architecture)
	assert.Empty(t, event.Runtime.Hostname)

	// Verify metrics with false enabled
	require.NotNil(t, event.Metrics)
	assert.False(t, event.Metrics.Enabled)

	// Verify empty arrays
	assert.Empty(t, event.Endpoints)
	assert.Empty(t, event.Configuration)

	emitter.AssertExpectations(t)
}

// Test for external job IDs and job names functionality
func TestService_pollBridge_WithJobInfo(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "bridge_status_response.json"), http.StatusOK)
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Create test job IDs and external job UUIDs
	testJobIDs := []int32{1, 2}

	testJob1 := job.Job{
		ID:            1,
		ExternalJobID: uuid.New(),
		Name:          null.StringFrom("BTC/USD Price Feed"),
	}
	testJob2 := job.Job{
		ID:            2,
		ExternalJobID: uuid.New(),
		Name:          null.StringFrom("ETH/USD Price Feed"),
	}

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return(testJobIDs, nil)
	jobORM.On("FindJob", mock.Anything, int32(1)).Return(testJob1, nil)
	jobORM.On("FindJob", mock.Anything, int32(2)).Return(testJob2, nil)

	// Capture the emitted protobuf to verify job information
	var capturedProtobufBytes []byte
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedProtobufBytes = args.Get(1).([]byte)
	})

	ctx := context.Background()
	service.pollBridge(ctx, "test-bridge", "http://example.com")

	// Verify the job information (IDs and names) were included in the protobuf
	require.NotEmpty(t, capturedProtobufBytes)
	var event events.BridgeStatusEvent
	err := proto.Unmarshal(capturedProtobufBytes, &event)
	require.NoError(t, err)

	require.Len(t, event.Jobs, 2)
	expectedJobs := []*events.JobInfo{
		{ExternalJobId: testJob1.ExternalJobID.String(), JobName: "BTC/USD Price Feed"},
		{ExternalJobId: testJob2.ExternalJobID.String(), JobName: "ETH/USD Price Feed"},
	}
	assert.ElementsMatch(t, expectedJobs, event.Jobs)

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_JobORMError(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "bridge_status_response.json"), http.StatusOK)
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM to return error
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, assert.AnError)

	// Should still emit telemetry with empty external job IDs
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "test-bridge", "http://example.com")

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

// Test ignoreJoblessBridges functionality
func TestService_pollBridge_IgnoreJoblessBridges_Enabled(t *testing.T) {
	// Use a nil httpClient since no HTTP request should be made when bridge is skipped for having no jobs
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, true, true)

	// Mock job ORM to return no job IDs (jobless bridge)
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "jobless-bridge").Return([]int32{}, nil)

	// Should NOT emit telemetry for jobless bridge when ignoreJoblessBridges is true
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)

	ctx := context.Background()
	service.pollBridge(ctx, "jobless-bridge", "http://example.com")

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_IgnoreJoblessBridges_Disabled(t *testing.T) {
	// Use valid HTTP client with successful response since we want the full flow when ignoreJoblessBridges is false
	httpClient := mocks.NewMockHTTPClient(loadFixture(t, "bridge_status_response.json"), http.StatusOK)
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, true, false)

	// Mock job ORM to return no job IDs (jobless bridge)
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "jobless-bridge").Return([]int32{}, nil)

	// Should emit telemetry for jobless bridge when ignoreJoblessBridges is false
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "jobless-bridge", "http://example.com")

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

// Test ignoreInvalidBridges functionality - HTTP error
func TestService_pollBridge_IgnoreInvalidBridges_HTTPError_Enabled(t *testing.T) {
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, true, false)

	// Create test job and external job UUID for FindJob mock
	testJob := job.Job{ID: 1, ExternalJobID: uuid.New(), Name: null.StringFrom("BTC/USD Price Feed")}

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "invalid-bridge").Return([]int32{1}, nil)
	jobORM.On("FindJob", mock.Anything, int32(1)).Return(testJob, nil)

	ctx := context.Background()
	service.pollBridge(ctx, "invalid-bridge", "http://invalid.invalid:8080")

	// Should NOT emit telemetry for invalid bridge when ignoreInvalidBridges is true
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_IgnoreInvalidBridges_HTTPError_Disabled(t *testing.T) {
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, false, false)

	// Create test job and external job UUID for FindJob mock
	testJob := job.Job{ID: 1, ExternalJobID: uuid.New(), Name: null.StringFrom("BTC/USD Price Feed")}

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "invalid-bridge").Return([]int32{1}, nil)
	jobORM.On("FindJob", mock.Anything, int32(1)).Return(testJob, nil)

	// Should emit empty telemetry for invalid bridge when ignoreInvalidBridges is false
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "invalid-bridge", "http://invalid.invalid:8080") // This will fail with HTTP error

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

// Test ignoreInvalidBridges functionality - Non-200 status
func TestService_pollBridge_IgnoreInvalidBridges_Non200Status_Enabled(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("Not Found", http.StatusNotFound)
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, true, false)

	// Create test job and external job UUID for FindJob mock
	testJob := job.Job{ID: 1, ExternalJobID: uuid.New(), Name: null.StringFrom("BTC/USD Price Feed")}

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "invalid-bridge").Return([]int32{1}, nil)
	jobORM.On("FindJob", mock.Anything, int32(1)).Return(testJob, nil)

	ctx := context.Background()
	service.pollBridge(ctx, "invalid-bridge", "http://example.com")

	// Should NOT emit telemetry for invalid bridge when ignoreInvalidBridges is true
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_IgnoreInvalidBridges_Non200Status_Disabled(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("Not Found", http.StatusNotFound)
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, false, false)

	// Create test job and external job UUID for FindJob mock
	testJob := job.Job{ID: 1, ExternalJobID: uuid.New(), Name: null.StringFrom("BTC/USD Price Feed")}

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "invalid-bridge").Return([]int32{1}, nil)
	jobORM.On("FindJob", mock.Anything, int32(1)).Return(testJob, nil)

	// Should emit empty telemetry for invalid bridge when ignoreInvalidBridges is false
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "invalid-bridge", "http://example.com")

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

// Test ignoreInvalidBridges functionality - Invalid JSON
func TestService_pollBridge_IgnoreInvalidBridges_InvalidJSON_Enabled(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("invalid json", http.StatusOK)
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, true, false)

	// Create test job and external job UUID for FindJob mock
	testJob := job.Job{ID: 1, ExternalJobID: uuid.New(), Name: null.StringFrom("BTC/USD Price Feed")}

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "invalid-bridge").Return([]int32{1}, nil)
	jobORM.On("FindJob", mock.Anything, int32(1)).Return(testJob, nil)

	ctx := context.Background()
	service.pollBridge(ctx, "invalid-bridge", "http://example.com")

	// Should NOT emit telemetry for invalid bridge when ignoreInvalidBridges is true
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_IgnoreInvalidBridges_InvalidJSON_Disabled(t *testing.T) {
	httpClient := mocks.NewMockHTTPClient("invalid json", http.StatusOK)
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, false, false)

	// Create test job and external job UUID for FindJob mock
	testJob := job.Job{ID: 1, ExternalJobID: uuid.New(), Name: null.StringFrom("BTC/USD Price Feed")}

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "invalid-bridge").Return([]int32{1}, nil)
	jobORM.On("FindJob", mock.Anything, int32(1)).Return(testJob, nil)

	// Should emit empty telemetry for invalid bridge when ignoreInvalidBridges is false
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "invalid-bridge", "http://example.com")

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

// Test combined functionality
func TestService_pollBridge_BothIgnoreFlags_Enabled(t *testing.T) {
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, true, true)

	// Mock job ORM to return no job IDs (jobless bridge)
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "jobless-invalid-bridge").Return([]int32{}, nil)

	ctx := context.Background()
	service.pollBridge(ctx, "jobless-invalid-bridge", "http://invalid.invalid:8080") // This would fail with HTTP error too

	// Should NOT emit telemetry - skipped because of no jobs (ignoreJoblessBridges)
	emitter.AssertNotCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
	emitter.AssertNotCalled(t, "With", mock.Anything)

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

func TestService_pollBridge_BothIgnoreFlags_Disabled(t *testing.T) {
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestServiceWithIgnoreFlags(t, true, testPollingInterval, httpClient, false, false)

	// Mock job ORM to return no job IDs (jobless bridge)
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "jobless-invalid-bridge").Return([]int32{}, nil)

	// Should emit empty telemetry even for jobless invalid bridge when both flags are false
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	service.pollBridge(ctx, "jobless-invalid-bridge", "http://invalid.invalid:8080") // This will fail with HTTP error

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}

// TestService_pollBridge_EndToEnd_RealWebServer tests the complete flow with a real HTTP server
func TestService_pollBridge_EndToEnd_RealWebServer(t *testing.T) {
	// Create a test HTTP server that serves fixture data
	fixtureData := loadFixture(t, "bridge_status_response.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path
		assert.Equal(t, testStatusPath, r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fixtureData))
	}))
	defer server.Close()

	// Use real HTTP client (not mock)
	httpClient := &http.Client{}
	service, _, jobORM, emitter := setupTestService(t, true, testPollingInterval, httpClient)

	// Mock job ORM calls
	jobORM.On("FindJobIDsWithBridge", mock.Anything, "test-bridge").Return([]int32{}, nil)

	// Capture the emitted protobuf to verify end-to-end flow
	var capturedProtobufBytes []byte
	emitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedProtobufBytes = args.Get(1).([]byte)
	})

	ctx := context.Background()
	service.pollBridge(ctx, "test-bridge", server.URL)

	// Verify the complete end-to-end flow worked
	require.NotEmpty(t, capturedProtobufBytes, "Should have emitted protobuf data")

	var event events.BridgeStatusEvent
	err := proto.Unmarshal(capturedProtobufBytes, &event)
	require.NoError(t, err, "Should be able to unmarshal protobuf")

	// Verify the data from fixture made it through the complete pipeline
	expectedStatus := loadFixtureAsEAResponse(t, "bridge_status_response.json")

	// Verify basic bridge info
	assert.Equal(t, "test-bridge", event.BridgeName)
	assert.Equal(t, expectedStatus.Adapter.Name, event.AdapterName)
	assert.Equal(t, expectedStatus.Adapter.Version, event.AdapterVersion)
	assert.InDelta(t, expectedStatus.Adapter.UptimeSeconds, event.AdapterUptimeSeconds, 0.001)
	assert.Equal(t, expectedStatus.DefaultEndpoint, event.DefaultEndpoint)

	// Verify endpoints - loop through fixture and compare with protobuf
	require.Len(t, event.Endpoints, len(expectedStatus.Endpoints))
	for i, expectedEndpoint := range expectedStatus.Endpoints {
		actualEndpoint := event.Endpoints[i]
		assert.Equal(t, expectedEndpoint.Name, actualEndpoint.Name)
		assert.Equal(t, expectedEndpoint.Aliases, actualEndpoint.Aliases)
		assert.Equal(t, expectedEndpoint.Transports, actualEndpoint.Transports)
	}

	// Verify configuration - loop through fixture and compare with protobuf
	// Helper function to safely convert values to strings, handling nil (same as in production code)
	safeString := func(v any) string {
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%v", v)
	}

	require.Len(t, event.Configuration, len(expectedStatus.Configuration))
	for i, expectedConfig := range expectedStatus.Configuration {
		actualConfig := event.Configuration[i]
		assert.Equal(t, expectedConfig.Name, actualConfig.Name)
		assert.Equal(t, safeString(expectedConfig.Value), actualConfig.Value)
		assert.Equal(t, expectedConfig.Type, actualConfig.Type)
		assert.Equal(t, expectedConfig.Description, actualConfig.Description)
		assert.Equal(t, expectedConfig.Required, actualConfig.Required)
		assert.Equal(t, safeString(expectedConfig.Default), actualConfig.DefaultValue)
		assert.Equal(t, expectedConfig.CustomSetting, actualConfig.CustomSetting)
		assert.Equal(t, safeString(expectedConfig.EnvDefaultOverride), actualConfig.EnvDefaultOverride)
	}

	// Verify runtime info
	require.NotNil(t, event.Runtime)
	assert.Equal(t, expectedStatus.Runtime.NodeVersion, event.Runtime.NodeVersion)
	assert.Equal(t, expectedStatus.Runtime.Platform, event.Runtime.Platform)
	assert.Equal(t, expectedStatus.Runtime.Architecture, event.Runtime.Architecture)
	assert.Equal(t, expectedStatus.Runtime.Hostname, event.Runtime.Hostname)

	// Verify metrics info
	require.NotNil(t, event.Metrics)
	assert.Equal(t, expectedStatus.Metrics.Enabled, event.Metrics.Enabled)

	// Verify job info is included
	assert.Empty(t, event.Jobs)

	// Verify timestamp is set
	assert.NotEmpty(t, event.Timestamp)

	jobORM.AssertExpectations(t)
	emitter.AssertExpectations(t)
}
