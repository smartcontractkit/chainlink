package cre

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	stvault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	httpactionconfig "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	multiGatewayConfigPath               = "/configs/workflow-gateway-capabilities-multi-gateway-don.toml"
	multiGatewayTopologyMarker           = "multi-gateway"
	multiDonTestOrgID                    = "multi-don-test-org"
	multiGatewayProxyDonEU               = "gateway_don_eu"
	multiGatewayEUEndpoint               = "/api/multi-gateway-eu-routed"
	multiGatewayEUTestCase               = "mgw-eu-get"
	multiGatewayEUSuccessMessage         = "HTTP Action CRUD success test completed: " + multiGatewayEUTestCase
	multiGatewayOutgoingMessageLogNeedle = "handling outgoing message"
)

func isMultiGatewayTopology(topologyName string) bool {
	return strings.Contains(topologyName, multiGatewayTopologyMarker)
}

func getMultiGatewayTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()
	return t_helpers.GetTestConfig(t, multiGatewayConfigPath)
}

func TestMultiGatewayTopology_LoadExpectedConfig(t *testing.T) {
	t.Parallel()

	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")

	cfg := &envconfig.Config{}
	require.NoError(t, cfg.Load(getMultiGatewayTestConfig(t).EnvironmentConfigPath))

	var foundWorkflow, foundGatewayUS, foundGatewayEU bool
	for _, nodeSet := range cfg.NodeSets {
		switch nodeSet.Name {
		case "workflow":
			foundWorkflow = true
		case "bootstrap-gateway-us":
			foundGatewayUS = true
			require.Equal(t, "gateway_don_us", nodeSet.GatewayDonID)
		case "gateway-eu":
			foundGatewayEU = true
			require.Equal(t, multiGatewayProxyDonEU, nodeSet.GatewayDonID)
		}
	}
	require.True(t, foundWorkflow)
	require.True(t, foundGatewayUS)
	require.True(t, foundGatewayEU)

	for _, nodeSet := range cfg.NodeSets {
		if nodeSet.Name != "workflow" {
			continue
		}

		settingsRaw := nodeSet.EnvVars["CL_CRE_SETTINGS"]
		require.NotEmpty(t, settingsRaw)

		var settings struct {
			Global struct {
				PropagateOrgIDInRequestMetadata string `json:"PropagateOrgIDInRequestMetadata"`
			} `json:"global"`
			Org map[string]struct {
				PerWorkflow struct {
					HTTPAction struct {
						GatewayProxyDonID string `json:"GatewayProxyDonID"`
					} `json:"HTTPAction"`
				} `json:"PerWorkflow"`
			} `json:"org"`
		}
		require.NoError(t, json.Unmarshal([]byte(settingsRaw), &settings))
		require.Equal(t, "true", settings.Global.PropagateOrgIDInRequestMetadata)

		orgSettings, ok := settings.Org[multiDonTestOrgID]
		require.True(t, ok, "expected org override for %q", multiDonTestOrgID)
		require.Equal(t, multiGatewayProxyDonEU, orgSettings.PerWorkflow.HTTPAction.GatewayProxyDonID)

		for _, nodeSpec := range nodeSet.NodeSpecs {
			require.Contains(t, nodeSpec.Node.UserConfigOverrides, "[CRE.Linking]")
			require.Contains(t, nodeSpec.Node.UserConfigOverrides, dockerHost+":18124")
		}
	}
}

// Test_CRE_V2_HTTP_Action_Multi_Gateway checks that HTTP action outbound requests use the gateway
// DON named by PerWorkflow.HTTPAction.GatewayProxyDonID for the workflow owner's org.
//
// Two gateway DONs (US and EU) are registered on the workflow connector. The test binds the owner
// to multi-don-test-org (linking service + CL_CRE_SETTINGS override to gateway_don_eu), runs an
// HTTP GET workflow, and expects EU gateway outbound logs while the US gateway stays idle.
// Skips unless TOPOLOGY_NAME contains "multi-gateway".
func Test_CRE_V2_HTTP_Action_Multi_Gateway(t *testing.T) {
	if !isMultiGatewayTopology(topology) {
		t.Skipf("skipping multi-gateway HTTP action test: TOPOLOGY_NAME=%q does not match %q", topology, multiGatewayTopologyMarker)
	}

	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, getMultiGatewayTestConfig(t))
	ExecuteHTTPActionMultiGatewayRoutingTest(t, testEnv)
}

// ExecuteHTTPActionMultiGatewayRoutingTest runs the multi-gateway routing scenario (standalone test and smoke suite).
func ExecuteHTTPActionMultiGatewayRoutingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()

	testLogger := framework.L
	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err)

	sc := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient
	workflowOwner := sc.MustGetRootKeyAddress().Hex()
	linkingService.SetOwnerOrg(workflowOwner, multiDonTestOrgID)

	fakeHTTP, err := fake.NewFakeDataProvider(testEnv.Config.FakeHTTP)
	require.NoError(t, err)

	response := map[string]any{"status": "success", "gateway": "eu"}
	require.NoError(t, fake.JSON("GET", multiGatewayEUEndpoint, response, 200))

	targetURL := fakeHTTP.BaseURLDocker + multiGatewayEUEndpoint
	testLogger.Info().
		Str("workflow_owner", workflowOwner).
		Str("org_id", multiDonTestOrgID).
		Str("url", targetURL).
		Msg("starting multi-gateway HTTP action routing test")

	const workflowFileLocation = "./httpaction/main.go"
	workflowConfig := httpactionconfig.Config{
		URL:      targetURL,
		TestCase: multiGatewayEUTestCase,
		Method:   "GET",
	}

	testID := uuid.New().String()[0:8]
	workflowName := "http-action-mgw-eu-" + testID
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	t_helpers.WatchWorkflowLogs(
		t,
		testLogger,
		userLogsCh,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		multiGatewayEUSuccessMessage,
		4*time.Minute,
		t_helpers.WithUserLogWorkflowID(workflowID),
	)

	t_helpers.AssertContainerLogsForNodeset(t, testEnv, "gateway-eu", multiGatewayOutgoingMessageLogNeedle)
	t_helpers.AssertContainerLogsAbsentForNodeset(t, testEnv, "bootstrap-gateway-us", multiGatewayOutgoingMessageLogNeedle)

	testLogger.Info().Msg("multi-gateway HTTP action routing test completed")
}
