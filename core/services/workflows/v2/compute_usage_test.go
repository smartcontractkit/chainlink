package v2_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/resourcemanager"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	meteringpb "github.com/smartcontractkit/chainlink-protos/metering/go"

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

type staticTracker struct {
	confidential bool
	measured     bool
	duration     time.Duration
}

func (s staticTracker) TookExecution(string) (v2.ConfidentialExecution, bool) {
	return v2.ConfidentialExecution{Measured: s.measured, Duration: s.duration}, s.confidential
}

// runOneExecution starts an engine with cfg, fires one trigger event whose
// module execution returns execErr, and returns the execution id.
func runOneExecution(t *testing.T, cfg *v2.EngineConfig, execErr error) string {
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

	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(nil, execErr).Once()
	eventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "usage_event"}}
	executionID := <-finishedCh
	require.NoError(t, engine.Close())
	return executionID
}

func TestEngine_ComputeUsageMeterRecord(t *testing.T) {
	t.Parallel()

	newUsageMeter := func(t *testing.T) (*resourcemanager.ResourceManager, *recordingEmitter) {
		emitter := &recordingEmitter{}
		rm := resourcemanager.NewResourceManager(logger.Test(t), resourcemanager.ResourceManagerConfig{
			MeterRecordsEnabled: true,
			Emitter:             emitter,
		})
		return rm, emitter
	}
	identity := resourcemanager.WithWorkflowUsagePool(resourcemanager.ResourceIdentity{Product: "cre", Service: resourcemanager.EmittingServiceWorkflowEngine}, resourcemanager.ResourceTypeWorkflowCompute)

	t.Run("emits one compute record per execution and logs the contract line", func(t *testing.T) {
		t.Parallel()
		lggr, obs := logger.TestObserved(t, zapcore.InfoLevel)
		cfg := defaultTestConfig(t, nil)
		cfg.Lggr = lggr
		rm, emitter := newUsageMeter(t)
		cfg.UsageMeter = rm
		cfg.UsageIdentity = identity

		executionID := runOneExecution(t, cfg, nil)

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
		rm, emitter := newUsageMeter(t)
		cfg.UsageMeter = rm
		cfg.UsageIdentity = identity

		runOneExecution(t, cfg, errors.New("module failed"))

		require.Len(t, emitter.all(), 1)
	})

	t.Run("bills the enclave-measured duration for confidential executions", func(t *testing.T) {
		t.Parallel()
		cfg := defaultTestConfig(t, nil)
		rm, emitter := newUsageMeter(t)
		cfg.UsageMeter = rm
		cfg.UsageIdentity = identity
		cfg.ConfidentialExecutions = staticTracker{confidential: true, measured: true, duration: 1234 * time.Millisecond}

		runOneExecution(t, cfg, nil)

		records := emitter.all()
		require.Len(t, records, 1)
		require.Equal(t, "1234", records[0].GetUtilizations()[0].GetValue())
	})

	t.Run("skips confidential executions the enclave did not measure", func(t *testing.T) {
		t.Parallel()
		lggr, obs := logger.TestObserved(t, zapcore.WarnLevel)
		cfg := defaultTestConfig(t, nil)
		cfg.Lggr = lggr
		rm, emitter := newUsageMeter(t)
		cfg.UsageMeter = rm
		cfg.UsageIdentity = identity
		cfg.ConfidentialExecutions = staticTracker{confidential: true}

		runOneExecution(t, cfg, nil)

		require.Empty(t, emitter.all())
		require.Len(t, obs.FilterMessage("Compute usage meter record not emitted: enclave did not report execution duration").All(), 1)
	})

	t.Run("emits nothing when no usage meter is configured", func(t *testing.T) {
		t.Parallel()
		lggr, obs := logger.TestObserved(t, zapcore.InfoLevel)
		cfg := defaultTestConfig(t, nil)
		cfg.Lggr = lggr

		runOneExecution(t, cfg, nil)

		require.Empty(t, obs.FilterMessage("Emitted capability usage meter record").All())
	})
}
