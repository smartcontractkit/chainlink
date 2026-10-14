package v2_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	confworkflowtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/resourcemanager"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	meteringpb "github.com/smartcontractkit/chainlink-protos/metering/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/confidentialrelay"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// recordingEmitter captures MeterRecords emitted by a ResourceManager.
type recordingEmitter struct {
	mu      sync.Mutex
	records []*meteringpb.MeterRecord
}

func (r *recordingEmitter) Emit(_ context.Context, body []byte, _ ...any) error {
	var rec meteringpb.MeterRecord
	if err := proto.Unmarshal(body, &rec); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, &rec)
	return nil
}

func (r *recordingEmitter) all() []*meteringpb.MeterRecord {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]*meteringpb.MeterRecord(nil), r.records...)
}

// moduleRun is what the mocked module does for the one trigger execution.
type moduleRun struct {
	result    *sdkpb.ExecutionResult
	err       error
	onExecute func()
}

// runOneExecution starts an engine with cfg, fires one trigger event, and
// returns the execution id.
func runOneExecution(t *testing.T, cfg *v2.EngineConfig, run moduleRun) string {
	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	initDoneCh := make(chan error)
	subscribedCh := make(chan []string, 1)
	finishedCh := make(chan string)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = setupMockBillingClient(t)
	cfg.OrgResolver = &mockOrgResolver{orgID: "org-123"}
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized:          func(err error) { initDoneCh <- err },
		OnSubscribedToTriggers: func(ids []string) { subscribedCh <- ids },
		OnExecutionFinished:    func(id string, _ string) { finishedCh <- id },
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedCh)

	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Run(func(context.Context, *sdkpb.ExecuteRequest, host.ExecutionHelper) {
		if run.onExecute != nil {
			run.onExecute()
		}
	}).Return(run.result, run.err).Once()
	eventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "usage_event"}}
	executionID := <-finishedCh
	require.NoError(t, engine.Close())
	return executionID
}

func TestEngine_ComputeUsageMeterRecord(t *testing.T) {
	t.Parallel()

	identity := resourcemanager.WithWorkflowUsagePool(resourcemanager.ResourceIdentity{Product: "cre", Service: resourcemanager.EmittingServiceWorkflowEngine}, resourcemanager.ResourceTypeWorkflowCompute)

	t.Run("emits one compute record per execution and logs the contract line", func(t *testing.T) {
		t.Parallel()
		lggr, obs := logger.TestObserved(t, zapcore.InfoLevel)
		cfg := defaultTestConfig(t, nil)
		cfg.Lggr = lggr
		emitter := &recordingEmitter{}
		cfg.MeterRecordsEnabled = true
		cfg.UsageEmitter = emitter
		cfg.UsageIdentity = identity

		executionID := runOneExecution(t, cfg, moduleRun{})

		records := emitter.all()
		require.Len(t, records, 1)
		rec := records[0]
		require.Equal(t, meteringpb.MeterAction_METER_ACTION_USAGE, rec.GetAction())
		require.Equal(t, resourcemanager.EmittingServiceWorkflowEngine, rec.GetIdentity().GetService())
		require.Equal(t, "cre:workflow:compute", rec.GetIdentity().GetResourcePool())
		require.Equal(t, "cre:workflow:compute", rec.GetIdentity().GetResourcePoolId())
		require.NotEmpty(t, rec.GetIdentity().GetDon().GetDonId(), "don id must be stamped from the local node")
		require.Len(t, rec.GetUtilizations(), 1)
		u := rec.GetUtilizations()[0]
		require.Equal(t, resourcemanager.ResourceTypeWorkflowCompute, u.GetResourceType())
		require.Equal(t, cfg.WorkflowID+":"+executionID, u.GetResourceId())
		require.Equal(t, executionID, u.GetEventId(), "compute capability event id is the execution id")
		require.Equal(t, "org-123", u.GetOrgId())
		require.NotEmpty(t, u.GetValue())

		logs := obs.FilterMessage("Emitted capability usage meter record").All()
		require.Len(t, logs, 1)
		fields := logs[0].ContextMap()
		require.Equal(t, executionID, fields["eventID"])
		require.Equal(t, resourcemanager.ResourceTypeWorkflowCompute, fields["resourceType"])
		require.Equal(t, "org-123", fields["orgID"])
		require.Contains(t, fields, "value")
	})

	t.Run("emits for failed executions too", func(t *testing.T) {
		t.Parallel()
		cfg := defaultTestConfig(t, nil)
		emitter := &recordingEmitter{}
		cfg.MeterRecordsEnabled = true
		cfg.UsageEmitter = emitter
		cfg.UsageIdentity = identity

		runOneExecution(t, cfg, moduleRun{err: errors.New("module failed")})

		require.Len(t, emitter.all(), 1)
	})

	t.Run("bills the compute duration the module reports", func(t *testing.T) {
		t.Parallel()
		cfg := defaultTestConfig(t, nil)
		emitter := &recordingEmitter{}
		cfg.MeterRecordsEnabled = true
		cfg.UsageEmitter = emitter
		cfg.UsageIdentity = identity
		runOneExecution(t, cfg, moduleRun{result: &sdkpb.ExecutionResult{ComputeDurationMs: proto.Int64(1234)}})

		records := emitter.all()
		require.Len(t, records, 1)
		require.Equal(t, "1234", records[0].GetUtilizations()[0].GetValue())
	})

	t.Run("falls back to wall-clock when the module reports no duration", func(t *testing.T) {
		t.Parallel()
		cfg := defaultTestConfig(t, nil)
		emitter := &recordingEmitter{}
		cfg.MeterRecordsEnabled = true
		cfg.UsageEmitter = emitter
		cfg.UsageIdentity = identity

		runOneExecution(t, cfg, moduleRun{onExecute: func() { time.Sleep(20 * time.Millisecond) }})

		records := emitter.all()
		require.Len(t, records, 1)
		ms, err := strconv.ParseInt(records[0].GetUtilizations()[0].GetValue(), 10, 64)
		require.NoError(t, err)
		require.GreaterOrEqual(t, ms, int64(20))
	})

	t.Run("emits nothing when metering is disabled", func(t *testing.T) {
		t.Parallel()
		lggr, obs := logger.TestObserved(t, zapcore.InfoLevel)
		cfg := defaultTestConfig(t, nil)
		cfg.Lggr = lggr

		runOneExecution(t, cfg, moduleRun{})

		require.Empty(t, obs.FilterMessage("Emitted capability usage meter record").All())
	})
}

// A TEE-routed trigger is metered from the duration the enclave set on the result.
func TestEngine_ComputeUsageMeterRecord_ConfidentialRun(t *testing.T) {
	t.Parallel()
	const capID = "confidential-workflows@1.0.0-alpha"
	lggr := logger.Test(t)
	cfg := defaultTestConfig(t, nil)
	emitter := &recordingEmitter{}
	cfg.MeterRecordsEnabled = true
	cfg.UsageEmitter = emitter
	cfg.UsageIdentity = resourcemanager.WithWorkflowUsagePool(resourcemanager.ResourceIdentity{Product: "cre", Service: resourcemanager.EmittingServiceWorkflowEngine}, resourcemanager.ResourceTypeWorkflowCompute)

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	execCap := capmocks.NewExecutableCapability(t)
	capreg.EXPECT().GetExecutable(matches.AnyContext, capID).Return(execCap, nil)

	tees, err := anypb.New(&confworkflowtypes.ProvidedTeesResponse{Tee: []*sdkpb.TeeTypeAndRegions{
		{Type: sdkpb.TeeType_TEE_TYPE_AWS_NITRO, Regions: []string{"us-east-1"}},
	}})
	require.NoError(t, err)
	execCap.EXPECT().Execute(matches.AnyContext, mock.MatchedBy(func(req capabilities.CapabilityRequest) bool {
		return req.Method == "ProvidedTees"
	})).Return(capabilities.CapabilityResponse{Payload: tees}, nil).Once()
	// The enclave runs the same WASM host, which sets compute_duration_ms.
	enclaveResp, err := anypb.New(&confworkflowtypes.ConfidentialWorkflowResponse{SdkExecutionResult: &sdkpb.ExecutionResult{
		Result:            &sdkpb.ExecutionResult_Value{Value: valuespb.NewStringValue("enclave-output")},
		ComputeDurationMs: proto.Int64(1234),
	}})
	require.NoError(t, err)
	execCap.EXPECT().Execute(matches.AnyContext, mock.MatchedBy(func(req capabilities.CapabilityRequest) bool {
		return req.Method == "Execute"
	})).Run(func(context.Context, capabilities.CapabilityRequest) {
		time.Sleep(20 * time.Millisecond) // round trip the engine must not bill
	}).Return(capabilities.CapabilityResponse{Payload: enclaveResp}, nil).Once()

	conf, err := v2.NewConfidentialModule(capreg, &confidentialrelay.ExecutionHandlers{}, "https://example.com/binary.wasm", []byte("hash"),
		cfg.WorkflowID, cfg.WorkflowOwner, cfg.WorkflowName.String(), cfg.WorkflowTag,
		func(context.Context, string) (string, error) { return "org-123", nil },
		limits.NewGateLimiter(true), cfg.LocalLimiters.Settings, lggr)
	require.NoError(t, err)

	local := modulemocks.NewModuleV2(t)
	local.EXPECT().Start().Maybe()
	local.EXPECT().Close().Maybe()
	teeSubs := newTriggerSubs(1)
	teeSubs.GetTriggerSubscriptions().Subscriptions[0].Requirements = &sdkpb.Requirements{Tee: &sdkpb.Tee{
		Item: &sdkpb.Tee_AnyRegions{AnyRegions: &sdkpb.Regions{Regions: []string{"us-east-1"}}},
	}}
	// Subscribe runs locally; the trigger is routed to the enclave.
	local.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(teeSubs, nil).Once()
	cfg.Module = host.NewRequirementSelectingModule(host.ModuleAndHandler{Module: local}, []host.ModuleAndHandler{{
		Module:              conf,
		RequirementsHandler: host.RequirementsHandler{Tee: conf.Tee},
	}})
	cfg.CapRegistry = capreg
	cfg.BillingClient = setupMockBillingClient(t)
	cfg.OrgResolver = &mockOrgResolver{orgID: "org-123"}

	initDoneCh := make(chan error)
	subscribedCh := make(chan []string, 1)
	finishedCh := make(chan string)
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized:          func(err error) { initDoneCh <- err },
		OnSubscribedToTriggers: func(ids []string) { subscribedCh <- ids },
		OnExecutionFinished:    func(id string, _ string) { finishedCh <- id },
	}
	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedCh)

	eventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "confidential_event"}}
	executionID := <-finishedCh
	require.NoError(t, engine.Close())

	records := emitter.all()
	require.Len(t, records, 1)
	u := records[0].GetUtilizations()[0]
	require.Equal(t, cfg.WorkflowID+":"+executionID, u.GetResourceId())
	require.Equal(t, "1234", u.GetValue(), "enclave-measured duration, not the round trip")
}
