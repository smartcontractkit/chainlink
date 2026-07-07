package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/resourcemanager"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	meteringpb "github.com/smartcontractkit/chainlink-protos/metering/go"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// recordingEmitter is a fake resourcemanager.Emitter that decodes and stores
// every emitted MeterRecord. If err is set, Emit fails instead.
type recordingEmitter struct {
	mu      sync.Mutex
	err     error
	records []*meteringpb.MeterRecord
}

func (r *recordingEmitter) Emit(_ context.Context, body []byte, _ ...any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return r.err
	}
	var record meteringpb.MeterRecord
	if err := proto.Unmarshal(body, &record); err != nil {
		return err
	}
	r.records = append(r.records, &record)
	return nil
}

func (r *recordingEmitter) Records() []*meteringpb.MeterRecord {
	r.mu.Lock()
	defer r.mu.Unlock()
	records := make([]*meteringpb.MeterRecord, len(r.records))
	copy(records, r.records)
	return records
}

func newMeteringResourceManager(t *testing.T, enabled bool, emitter resourcemanager.Emitter) *resourcemanager.ResourceManager {
	t.Helper()
	return resourcemanager.NewResourceManager(commonlogger.Test(t), resourcemanager.ResourceManagerConfig{
		MeterRecordsEnabled: enabled,
		Emitter:             emitter,
	})
}

// recordingSnapshotEmitter is a fake resourcemanager.Emitter that decodes and
// stores every emitted MeterSnapshot (one per active resource). Used to drive
// the Meterable snapshot path.
type recordingSnapshotEmitter struct {
	mu        sync.Mutex
	snapshots []*meteringpb.MeterSnapshot
}

func (r *recordingSnapshotEmitter) Emit(_ context.Context, body []byte, _ ...any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var snapshot meteringpb.MeterSnapshot
	if err := proto.Unmarshal(body, &snapshot); err != nil {
		return err
	}
	r.snapshots = append(r.snapshots, &snapshot)
	return nil
}

func (r *recordingSnapshotEmitter) Snapshots() []*meteringpb.MeterSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	snapshots := make([]*meteringpb.MeterSnapshot, len(r.snapshots))
	copy(snapshots, r.snapshots)
	return snapshots
}

func newMeteringTestHandler(t *testing.T, artifactsStore WorkflowArtifactsStore, rm *resourcemanager.ResourceManager) *eventHandler {
	t.Helper()
	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}
	registry := capabilities.NewRegistry(lggr)
	registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
	limiters, err := v2.NewLimiters(lf, nil)
	require.NoError(t, err)
	rl, err := ratelimiter.NewRateLimiter(rlConfig)
	require.NoError(t, err)
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, wlConfig, lf)
	require.NoError(t, err)

	h, err := NewEventHandler(
		lggr,
		workflowstore.NewInMemoryStore(lggr, clockwork.NewFakeClock()),
		nil,
		true,
		registry,
		NewEngineRegistry(),
		custmsg.NewLabeler(),
		limiters,
		nil,
		rl,
		workflowLimits,
		artifactsStore,
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		&testDonNotifier{},
		WithResourceManager(rm),
		WithEngineFactoryFn(mockEngineFactory),
	)
	require.NoError(t, err)
	return h
}

func requireMeterRecord(t *testing.T, record *meteringpb.MeterRecord, action meteringpb.MeterAction, originatingEvent WorkflowRegistryEventName, workflowID string) {
	t.Helper()
	require.NotNil(t, record.Identity)
	assert.Equal(t, "workflow-syncer-v2", record.Identity.Service)
	assert.Equal(t, "workflow_specs_v2", record.Identity.ResourcePool)
	assert.Equal(t, action, record.Action)
	assert.NotNil(t, record.Timestamp)
	require.Len(t, record.Utilizations, 1)
	util := record.Utilizations[0]
	assert.Equal(t, "1", util.Value)
	assert.Equal(t, "operations", util.ResourceType)
	// resource_id = workflow_id for the syncer (no shared physical resource).
	assert.Equal(t, workflowID, util.ResourceId)
	assert.Equal(t, string(originatingEvent), util.EventId)
}

func Test_meterRecords(t *testing.T) {
	t.Parallel()

	wfOwner := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	wfOwnerHex := hex.EncodeToString(wfOwner)

	t.Run("registered event creating a new spec emits RESERVE", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{}, newMeteringResourceManager(t, true, emitter))

		wfID := types.WorkflowID{1}
		err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    wfID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		}, WorkflowRegistered)
		require.NoError(t, err)

		records := emitter.Records()
		require.Len(t, records, 1)
		requireMeterRecord(t, records[0], meteringpb.MeterAction_METER_ACTION_RESERVE, WorkflowRegistered, wfID.Hex())
	})

	t.Run("retried registered event emits equivalent utilization identity", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		// The stub never returns a stored spec, so each call replays the
		// new-spec path exactly as a reprocessed event would.
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{}, newMeteringResourceManager(t, true, emitter))

		event := WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    types.WorkflowID{2},
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		}
		require.NoError(t, h.workflowRegisteredEvent(t.Context(), event, WorkflowRegistered))
		require.NoError(t, h.workflowRegisteredEvent(t.Context(), event, WorkflowRegistered))

		records := emitter.Records()
		require.Len(t, records, 2)
		require.Len(t, records[0].Utilizations, 1)
		require.Len(t, records[1].Utilizations, 1)
		assert.Equal(t, records[0].Action, records[1].Action)
		assert.Equal(t, records[0].Identity.GetService(), records[1].Identity.GetService())
		assert.Equal(t, records[0].Identity.GetResourcePool(), records[1].Identity.GetResourcePool())
		assert.Equal(t, records[0].Utilizations[0].GetResourceId(), records[1].Utilizations[0].GetResourceId())
		assert.Equal(t, records[0].Utilizations[0].GetEventId(), records[1].Utilizations[0].GetEventId())
	})

	t.Run("activated event with existing spec emits UPDATE", func(t *testing.T) {
		t.Parallel()
		binary := []byte("binary-data")
		config := []byte("")
		giveWFID, err := pkgworkflows.GenerateWorkflowID(wfOwner, "wf-name", binary, config, "")
		require.NoError(t, err)
		wfID := types.WorkflowID(giveWFID)

		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				Workflow:      hex.EncodeToString(binary),
				Config:        string(config),
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusPaused,
				WorkflowOwner: wfOwnerHex,
				WorkflowName:  "wf-name",
			},
		}, newMeteringResourceManager(t, true, emitter))

		err = h.workflowActivatedEvent(t.Context(), WorkflowActivatedEvent{
			Status:        WorkflowStatusActive,
			WorkflowID:    wfID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		})
		require.NoError(t, err)

		records := emitter.Records()
		require.Len(t, records, 1)
		requireMeterRecord(t, records[0], meteringpb.MeterAction_METER_ACTION_UPDATE, WorkflowActivated, wfID.Hex())
	})

	t.Run("paused event emits RELEASE after artifacts are deleted", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{3}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: wfOwnerHex,
			},
		}, newMeteringResourceManager(t, true, emitter))

		require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: wfID}))

		records := emitter.Records()
		require.Len(t, records, 1)
		requireMeterRecord(t, records[0], meteringpb.MeterAction_METER_ACTION_RELEASE, WorkflowPaused, wfID.Hex())
	})

	t.Run("deleted event emits RELEASE after artifacts are deleted", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{4}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: wfOwnerHex,
			},
		}, newMeteringResourceManager(t, true, emitter))

		require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, WorkflowDeleted))

		records := emitter.Records()
		require.Len(t, records, 1)
		requireMeterRecord(t, records[0], meteringpb.MeterAction_METER_ACTION_RELEASE, WorkflowDeleted, wfID.Hex())
	})

	t.Run("no record when creating a new spec fails", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{upsertErr: assert.AnError}, newMeteringResourceManager(t, true, emitter))

		err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    types.WorkflowID{5},
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		}, WorkflowRegistered)
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, emitter.Records())
	})

	t.Run("no record when a status-only update fails", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{6}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: wfOwnerHex,
			},
			upsertErr: assert.AnError,
		}, newMeteringResourceManager(t, true, emitter))

		err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    wfID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		}, WorkflowRegistered)
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, emitter.Records())
	})

	t.Run("no record when deleting artifacts fails", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{7}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: wfOwnerHex,
			},
			deleteErr: assert.AnError,
		}, newMeteringResourceManager(t, true, emitter))

		err := h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, WorkflowDeleted)
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, emitter.Records())
	})

	t.Run("no record while a delete is deferred by drain; exactly one on the successful retry", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{8}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: wfOwnerHex,
			},
		}, newMeteringResourceManager(t, true, emitter))

		drainable := &mockDrainableEngine{}
		drainable.activeExecutions.Store(1)
		require.NoError(t, h.engineRegistry.Add(wfID, "test-source", drainable))

		err := h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, WorkflowDeleted)
		require.ErrorIs(t, err, ErrDrainInProgress)
		assert.Empty(t, emitter.Records())

		drainable.activeExecutions.Store(0)
		require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, WorkflowDeleted))

		records := emitter.Records()
		require.Len(t, records, 1)
		requireMeterRecord(t, records[0], meteringpb.MeterAction_METER_ACTION_RELEASE, WorkflowDeleted, wfID.Hex())
	})

	t.Run("emit failure never fails event handling", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{9}
		emitter := &recordingEmitter{err: errors.New("beholder unavailable")}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: wfOwnerHex,
			},
		}, newMeteringResourceManager(t, true, emitter))

		require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    wfID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		}, WorkflowRegistered))
		require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, WorkflowDeleted))
		assert.Empty(t, emitter.Records())
	})

	t.Run("disabled resource manager emits nothing", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{}, newMeteringResourceManager(t, false, emitter))

		require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    types.WorkflowID{10},
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		}, WorkflowRegistered))
		assert.Empty(t, emitter.Records())
	})

	t.Run("snapshot emits one MeterSnapshot per running engine", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingSnapshotEmitter{}
		clock := clockwork.NewFakeClockAt(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
		rm := resourcemanager.NewResourceManager(commonlogger.Test(t), resourcemanager.ResourceManagerConfig{
			MeterRecordsEnabled:   true,
			MeterSnapshotsEnabled: true,
			Emitter:               emitter,
			SnapshotInterval:      time.Minute,
			Clock:                 clock,
		})
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{}, rm)
		unregister := rm.Register(h)
		t.Cleanup(unregister)

		// Register two engines; the running engine registry is the in-memory source
		// the snapshot reads (no per-tick GetWorkflowSpec).
		wfID1 := types.WorkflowID{20}
		wfID2 := types.WorkflowID{21}
		require.NoError(t, h.engineRegistry.Add(wfID1, "test-source", &fakeService{}))
		require.NoError(t, h.engineRegistry.Add(wfID2, "test-source", &fakeService{}))

		servicetest.Run(t, rm)
		require.NoError(t, clock.BlockUntilContext(t.Context(), 1))
		clock.Advance(time.Minute)

		require.Eventually(t, func() bool {
			return len(emitter.Snapshots()) == 2
		}, time.Second, time.Millisecond)
		snapshots := emitter.Snapshots()
		require.Len(t, snapshots, 2)

		byWorkflowID := map[string]*meteringpb.MeterSnapshot{}
		for _, snap := range snapshots {
			require.NotNil(t, snap.Identity)
			assert.Equal(t, "workflow-syncer-v2", snap.Identity.Service)
			assert.Equal(t, "workflow_specs_v2", snap.Identity.ResourcePool)
			require.Len(t, snap.Utilization, 1)
			assert.Equal(t, "1", snap.Utilization[0].Value)
			// resource_id = workflow_id fully identifies the resource; no labels.
			byWorkflowID[snap.Utilization[0].ResourceId] = snap
		}
		require.NotNil(t, byWorkflowID[wfID1.Hex()], "snapshot must contain an entry for the first running engine")
		require.NotNil(t, byWorkflowID[wfID2.Hex()], "snapshot must contain an entry for the second running engine")
	})

	t.Run("graceful close emits a RELEASE per running engine", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{}, newMeteringResourceManager(t, true, emitter))

		wfID := types.WorkflowID{30}
		require.NoError(t, h.engineRegistry.Add(wfID, "test-source", &fakeService{}))

		es := h.engineRegistry.GetAll()
		h.emitGracefulCloseReleases(t.Context(), es)

		records := emitter.Records()
		require.Len(t, records, 1)
		requireMeterRecord(t, records[0], meteringpb.MeterAction_METER_ACTION_RELEASE, WorkflowDeleted, wfID.Hex())
	})
}
