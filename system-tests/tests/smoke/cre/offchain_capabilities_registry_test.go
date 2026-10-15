package cre

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	job_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const cronCapabilityID = "cron-trigger@1.0.0"

// Test_CRE_V2_OffchainCapabilitiesRegistry runs on a topology where the workflow DON has
// [Capabilities.Local] UseOffchainRegistry = true, so the locally launched cron capability takes its
// spec_config only from a capabilities_registry cresettings job delivered via JD.
//
// The offchain config lowers cron's fastestScheduleIntervalSeconds from the 30s default to 5s. A
// workflow scheduled every 5s registers and executes only if that config reached the capability.
//
//nolint:paralleltest // the test owns the topology it runs on
func Test_CRE_V2_OffchainCapabilitiesRegistry(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(
		t,
		t_helpers.GetTestConfig(t, "/configs/workflow-gateway-capabilities-don-offchain-registry.toml"),
	)
	ExecuteOffchainCapabilitiesRegistryTest(t, testEnv)
}

func ExecuteOffchainCapabilitiesRegistryTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	workflowDON := testEnv.Dons.MustWorkflowDON()

	// Versions must strictly increase across runs against the same environment.
	version := uint64(time.Now().Unix()) //nolint:gosec // G115: positive
	payload := offchainRegistryPayload(t, testEnv, version, map[string]map[string]map[string]any{
		workflowDON.Name: {cronCapabilityID: {"fastestScheduleIntervalSeconds": 5}},
	})
	proposeAndApproveCapabilitiesRegistryJob(t, testEnv, workflowDON, payload, testLogger)

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowName := t_helpers.UniqueWorkflowName(testEnv, "offchaincapreg")
	workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName,
		&crontypes.WorkflowConfig{Schedule: "*/5 * * * * *"},
		"../../../../core/scripts/cre/environment/examples/workflows/cron/main.go",
	)

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog,
		"Amazing workflow user log", 4*time.Minute, t_helpers.WithUserLogWorkflowID(workflowID))
	testLogger.Info().Msg("Offchain capabilities registry test completed")
}

// offchainRegistryPayload builds the proto-JSON offchain_config from DON name -> capability ID -> spec_config.
func offchainRegistryPayload(t *testing.T, testEnv *ttypes.TestEnvironment, version uint64, specConfigs map[string]map[string]map[string]any) string {
	t.Helper()

	reg := &capabilitiespb.OffchainCapabilitiesRegistry{
		Domain:  offchain.ProductLabel,
		Env:     testEnv.CreEnvironment.CldfEnvironment.Name,
		Version: version,
		Dons:    map[string]*capabilitiespb.DONConfig{},
	}
	for donName, caps := range specConfigs {
		donCfg := &capabilitiespb.DONConfig{Capabilities: map[string]*capabilitiespb.CapabilityConfig{}}
		for capID, specConfig := range caps {
			m, err := values.WrapMap(specConfig)
			require.NoError(t, err, "failed to wrap spec_config of %s", capID)
			donCfg.Capabilities[capID] = &capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(m)}
		}
		reg.Dons[donName] = donCfg
	}

	out, err := protojson.Marshal(reg)
	require.NoError(t, err, "failed to marshal offchain capabilities registry")
	return string(out)
}

func proposeAndApproveCapabilitiesRegistryJob(t *testing.T, testEnv *ttypes.TestEnvironment, targetDON *cre.Don, offchainConfig string, testLogger zerolog.Logger) {
	t.Helper()

	report, err := operations.ExecuteOperation(
		testEnv.CreEnvironment.CldfEnvironment.OperationsBundle,
		job_ops.ProposeCapabilitiesRegistryJob,
		job_ops.ProposeCapabilitiesRegistryJobDeps{Env: *testEnv.CreEnvironment.CldfEnvironment},
		job_ops.ProposeCapabilitiesRegistryJobInput{
			Domain:         offchain.ProductLabel,
			Environment:    testEnv.CreEnvironment.CldfEnvironment.Name,
			DONName:        targetDON.Name,
			DONFilters:     []offchain.TargetDONFilter{{Key: offchain.FilterKeyDONName, Value: targetDON.Name}},
			ExtraLabels:    map[string]string{cre.CapabilityLabelKey: "capabilities-registry"},
			OffchainConfig: offchainConfig,
		},
	)
	require.NoError(t, err, "failed to propose capabilities registry job")

	// cresettings specs are auto-approved on proposal; jobs.Approve tolerates that.
	err = jobs.Approve(t.Context(), testEnv.CreEnvironment.CldfEnvironment.Offchain, testEnv.Dons, report.Output.Specs)
	require.NoError(t, err, "failed to approve capabilities registry job")
	testLogger.Info().Str("don", targetDON.Name).Msg("Capabilities registry job proposed and approved")
}
