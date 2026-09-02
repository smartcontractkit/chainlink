package v2

import (
	"context"
	"errors"
	"math/big"
	"strconv"
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
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	meteringpb "github.com/smartcontractkit/chainlink-protos/metering/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/confidentialrelay"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// meteringwfID is the workflow ID that GenerateWorkflowID produces from the
// stub store's FetchWorkflowArtifacts return values (binary="binary",
// config="config") with owner 0xaabbccdd and name "wf-name". Tests that use
// persistUpserts and register as Active must use this ID so tryEngineCreate's
// validation passes.
var meteringwfID = func() types.WorkflowID {
	owner := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	wfIDBytes, _ := pkgworkflows.GenerateWorkflowID(owner, "wf-name", []byte("binary"), []byte("config"), "")
	return types.WorkflowID(wfIDBytes)
}()

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

// newTestSpecMeter builds a SpecMeter over rm and artifactsStore, or returns
// nil (the no-op meter) when rm is nil.
func newTestSpecMeter(t *testing.T, rm *resourcemanager.ResourceManager, artifactsStore WorkflowArtifactsStore, org orgresolver.OrgResolver) *SpecMeter {
	t.Helper()
	if rm == nil {
		return nil
	}
	sm, err := NewSpecMeter(logger.TestLogger(t), rm, resourcemanager.ResourceIdentity{}, artifactsStore, org)
	require.NoError(t, err)
	return sm
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
		&confidentialrelay.ExecutionHandlers{},
		NewEngineRegistry(),
		custmsg.NewLabeler(),
		limiters,
		nil,
		rl,
		workflowLimits,
		artifactsStore,
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		&testDonNotifier{},
		WithSpecMeter(newTestSpecMeter(t, rm, artifactsStore, nil)),
		WithEngineFactoryFn(mockEngineFactory),
	)
	require.NoError(t, err)
	return h
}

// requireSpecDelta asserts that record is a workflow-syncer-v2 spec delta on
// the workflow-count billing unit: a METER_ACTION_UPDATE carrying a single
// utilization with the given signed value and resource_id (= workflow_id).
// event_id must be the deterministic, cross-node-identical id wantEventID (an
// opaque string that is never format-validated).
func requireSpecDelta(t *testing.T, record *meteringpb.MeterRecord, value, workflowID, wantEventID string) {
	t.Helper()
	requireSpecDeltaRecord(t, record, "operations", value, workflowID, wantEventID)
}

// requireSpecBytesDelta asserts that record is a workflow-syncer-v2 spec delta
// on the storage billing unit: identical shape to requireSpecDelta except the
// resource_type is "storage_bytes" and the value is a signed byte count.
func requireSpecBytesDelta(t *testing.T, record *meteringpb.MeterRecord, value, workflowID, wantEventID string) {
	t.Helper()
	requireSpecDeltaRecord(t, record, "storage_bytes", value, workflowID, wantEventID)
}

func requireSpecDeltaRecord(t *testing.T, record *meteringpb.MeterRecord, resourceType, value, workflowID, wantEventID string) {
	t.Helper()
	require.NotNil(t, record.Identity)
	assert.Equal(t, "workflow-syncer-v2", record.Identity.Service)
	assert.Equal(t, "workflow_specs_v2", record.Identity.ResourcePool)
	// Every syncer meter record is a signed level delta (UPDATE), never a
	// RESERVE/RELEASE lifecycle edge.
	assert.Equal(t, meteringpb.MeterAction_METER_ACTION_UPDATE, record.Action)
	assert.NotNil(t, record.Timestamp)
	require.Len(t, record.Utilizations, 1)
	util := record.Utilizations[0]
	assert.Equal(t, value, util.Value)
	assert.Equal(t, resourceType, util.ResourceType)
	// resource_id = workflow_id for the syncer (no shared physical resource).
	assert.Equal(t, workflowID, util.ResourceId)
	// event_id is the deterministic reconciliation-derived id, identical on every
	// workflow-DON node; it is an opaque string with no validated format.
	assert.Equal(t, wantEventID, util.EventId)
}

// meteringBytesEventID mirrors the production derivation of the storage-delta
// event_id: the count-delta event_id hashed under the storage namespace, so
// the two records of one transition never collide in the consumer dedup key
// space while remaining deterministic and cross-node-identical.
func meteringBytesEventID(countEventID string) string {
	return resourcemanager.EventID("workflow-spec-storage-bytes", countEventID)
}

// meteringTestBytes is the durable storage footprint the stub artifacts store
// reports: binary "binary" (6 bytes) + config "config" (6 bytes).
const meteringTestBytes = int64(len("binary") + len("config"))

func Test_meterRecords(t *testing.T) {
	t.Parallel()

	wfOwner := []byte{0xaa, 0xbb, 0xcc, 0xdd}

	t.Run("registered event persisting a new spec emits +1 and +storage bytes", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{}, newMeteringResourceManager(t, true, emitter))

		wfID := types.WorkflowID{1}
		err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    wfID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		})
		require.NoError(t, err)

		records := emitter.Records()
		require.Len(t, records, 2)
		wantCountID := resourcemanager.EventID("workflow-spec-register", wfID.Hex(), "0")
		requireSpecDelta(t, records[0], "1", wfID.Hex(), wantCountID)
		requireSpecBytesDelta(t, records[1], "12", wfID.Hex(), meteringBytesEventID(wantCountID))
	})

	t.Run("reprocessed registered event emits the IDENTICAL event_ids (cross-node dedup)", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		// The stub never returns a stored spec, so each call replays the
		// new-spec path exactly as a reprocessed event (or a second node
		// reconciling the same on-chain event) would.
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{}, newMeteringResourceManager(t, true, emitter))

		event := WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    types.WorkflowID{2},
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
			CreatedAt:     123,
		}
		require.NoError(t, h.workflowRegisteredEvent(t.Context(), event))
		require.NoError(t, h.workflowRegisteredEvent(t.Context(), event))

		// Each processing emits one count record + one storage record; the
		// first pair is [count, bytes], the second [count, bytes].
		records := emitter.Records()
		require.Len(t, records, 4)
		for _, r := range records {
			require.Len(t, r.Utilizations, 1)
		}
		// The same reconciliation event MUST yield the identical event_id per
		// billing unit, so the billing consumer dedups reprocessing and
		// cross-node duplicates. Both are derived deterministically from the
		// on-chain workflowID + CreatedAt.
		want := resourcemanager.EventID("workflow-spec-register", types.WorkflowID{2}.Hex(), "123")
		assert.Equal(t, want, records[0].Utilizations[0].GetEventId())
		assert.Equal(t, records[0].Utilizations[0].GetEventId(), records[2].Utilizations[0].GetEventId())
		wantBytes := meteringBytesEventID(want)
		assert.Equal(t, wantBytes, records[1].Utilizations[0].GetEventId())
		assert.Equal(t, records[1].Utilizations[0].GetEventId(), records[3].Utilizations[0].GetEventId())
		// The two billing units must never share an event_id, or consumer dedup
		// would collapse the pair.
		assert.NotEqual(t, records[0].Utilizations[0].GetEventId(), records[1].Utilizations[0].GetEventId())
	})

	t.Run("activating an already-stored spec emits nothing (status-only update)", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{3}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: "aabbccdd",
				WorkflowName:  "wf-name",
			},
		}, newMeteringResourceManager(t, true, emitter))

		// Same workflow ID already stored, only the status changes: the spec
		// stays stored, so no artifact-persistence transition and no delta.
		err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    wfID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		})
		require.NoError(t, err)

		assert.Empty(t, emitter.Records())
	})

	t.Run("paused event emits nothing", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{4}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: "aabbccdd",
			},
		}, newMeteringResourceManager(t, true, emitter))

		require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: wfID}))

		// Pause routes through the delete path, but the spec's release is
		// realized by its absence from subsequent snapshots, not a delta.
		assert.Empty(t, emitter.Records())
	})

	t.Run("deleted event emits -1 and -storage bytes after artifacts are deleted", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{5}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: "aabbccdd",
				StorageBytes:  meteringTestBytes,
			},
		}, newMeteringResourceManager(t, true, emitter))

		require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd"))

		records := emitter.Records()
		require.Len(t, records, 2)
		wantCountID := resourcemanager.EventID("workflow-spec-delete", wfID.Hex())
		requireSpecDelta(t, records[0], "-1", wfID.Hex(), wantCountID)
		requireSpecBytesDelta(t, records[1], "-12", wfID.Hex(), meteringBytesEventID(wantCountID))
	})

	t.Run("no record when persisting a new spec fails", func(t *testing.T) {
		t.Parallel()
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{upsertErr: assert.AnError}, newMeteringResourceManager(t, true, emitter))

		err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    types.WorkflowID{6},
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		})
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
				WorkflowOwner: "aabbccdd",
			},
			deleteErr: assert.AnError,
		}, newMeteringResourceManager(t, true, emitter))

		err := h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd")
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, emitter.Records())
	})

	t.Run("no record while a delete is deferred by drain; exactly one pair on the successful retry", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{8}
		emitter := &recordingEmitter{}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: "aabbccdd",
				StorageBytes:  meteringTestBytes,
			},
		}, newMeteringResourceManager(t, true, emitter))

		drainable := &mockDrainableEngine{}
		drainable.activeExecutions.Store(1)
		require.NoError(t, h.engineRegistry.Add(wfID, "test-source", drainable))

		err := h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd")
		require.ErrorIs(t, err, ErrDrainInProgress)
		assert.Empty(t, emitter.Records())

		drainable.activeExecutions.Store(0)
		require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd"))

		records := emitter.Records()
		require.Len(t, records, 2)
		wantCountID := resourcemanager.EventID("workflow-spec-delete", wfID.Hex())
		requireSpecDelta(t, records[0], "-1", wfID.Hex(), wantCountID)
		requireSpecBytesDelta(t, records[1], "-12", wfID.Hex(), meteringBytesEventID(wantCountID))
	})

	t.Run("emit failure never fails event handling", func(t *testing.T) {
		t.Parallel()
		wfID := types.WorkflowID{9}
		emitter := &recordingEmitter{err: errors.New("beholder unavailable")}
		h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
			spec: &job.WorkflowSpec{
				WorkflowID:    wfID.Hex(),
				Status:        job.WorkflowSpecStatusActive,
				WorkflowOwner: "aabbccdd",
			},
		}, newMeteringResourceManager(t, true, emitter))

		require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
			Status:        WorkflowStatusPaused,
			WorkflowID:    wfID,
			WorkflowOwner: wfOwner,
			WorkflowName:  "wf-name",
		}))
		require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd"))
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
		}))
		assert.Empty(t, emitter.Records())
	})

	t.Run("snapshot emits one MeterSnapshot per persisted spec", func(t *testing.T) {
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

		// The snapshot enumerates persisted specs from the store, not running
		// engines, so paused-but-stored specs are still accounted for.
		wfID1 := types.WorkflowID{20}
		wfID2 := types.WorkflowID{21}
		store := &stubWorkflowArtifactsStore{
			specs: []*job.WorkflowSpec{
				{WorkflowID: wfID1.Hex(), WorkflowOwner: "aabbccdd", StorageBytes: 100},
				{WorkflowID: wfID2.Hex(), WorkflowOwner: "aabbccdd", StorageBytes: 200},
			},
		}
		h := newMeteringTestHandler(t, store, rm)
		unregister := rm.Register(h.specMeter)
		t.Cleanup(unregister)

		servicetest.Run(t, rm)
		require.NoError(t, clock.BlockUntilContext(t.Context(), 1))
		clock.Advance(time.Minute)

		require.Eventually(t, func() bool {
			return len(emitter.Snapshots()) == 2
		}, time.Second, time.Millisecond)
		snapshots := emitter.Snapshots()
		require.Len(t, snapshots, 2)

		// Each snapshot carries both billed dimensions of its spec: the
		// workflow-count level (1) and the storage level (storage_bytes).
		wantBytes := map[string]string{wfID1.Hex(): "100", wfID2.Hex(): "200"}
		byWorkflowID := map[string]*meteringpb.MeterSnapshot{}
		for _, snap := range snapshots {
			require.NotNil(t, snap.Identity)
			assert.Equal(t, "workflow-syncer-v2", snap.Identity.Service)
			assert.Equal(t, "workflow_specs_v2", snap.Identity.ResourcePool)
			require.Len(t, snap.Utilization, 2)
			countUtil, bytesUtil := snap.Utilization[0], snap.Utilization[1]
			assert.Equal(t, "1", countUtil.Value)
			assert.Equal(t, "operations", countUtil.ResourceType)
			assert.Equal(t, "storage_bytes", bytesUtil.ResourceType)
			assert.Equal(t, wantBytes[bytesUtil.ResourceId], bytesUtil.Value)
			assert.Equal(t, countUtil.ResourceId, bytesUtil.ResourceId)
			// resource_id = workflow_id fully identifies the resource; no labels.
			byWorkflowID[countUtil.ResourceId] = snap
		}
		require.NotNil(t, byWorkflowID[wfID1.Hex()], "snapshot must contain an entry for the first persisted spec")
		require.NotNil(t, byWorkflowID[wfID2.Hex()], "snapshot must contain an entry for the second persisted spec")
	})
}

// fakeOrgResolver is a no-network OrgResolver that maps workflow owners (hex)
// to organization IDs. It asserts that metering org resolution is exercised
// (previously org was entirely untested).
type fakeOrgResolver struct {
	orgs map[string]string
}

func (f *fakeOrgResolver) Get(_ context.Context, owner string) (string, error) {
	if org, ok := f.orgs[owner]; ok {
		return org, nil
	}
	return "", errors.New("owner not found")
}

func (f *fakeOrgResolver) Start(context.Context) error    { return nil }
func (f *fakeOrgResolver) Close() error                   { return nil }
func (f *fakeOrgResolver) HealthReport() map[string]error { return nil }
func (f *fakeOrgResolver) Ready() error                   { return nil }
func (f *fakeOrgResolver) Name() string                   { return "fake-org-resolver" }

// newMeteringTestHandlerWithOrg is newMeteringTestHandler plus a wired
// OrgResolver, so tests can assert OrgId on emitted records.
func newMeteringTestHandlerWithOrg(t *testing.T, artifactsStore WorkflowArtifactsStore, rm *resourcemanager.ResourceManager, org orgresolver.OrgResolver) *eventHandler {
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
		&confidentialrelay.ExecutionHandlers{},
		NewEngineRegistry(),
		custmsg.NewLabeler(),
		limiters,
		nil,
		rl,
		workflowLimits,
		artifactsStore,
		workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		&testDonNotifier{},
		WithSpecMeter(newTestSpecMeter(t, rm, artifactsStore, org)),
		WithEngineFactoryFn(mockEngineFactory),
		WithOrgResolver(org),
	)
	require.NoError(t, err)
	return h
}

// redeliveredDeleteSpecs returns a fresh stub seeded with a stored spec and the
// owner the metering tests use, so a delete + redelivered-delete pair can be
// driven through Handle (which pre-reads the spec itself).
func redeliveredDeleteStore() *stubWorkflowArtifactsStore {
	return &stubWorkflowArtifactsStore{
		spec: &job.WorkflowSpec{
			WorkflowID:    types.WorkflowID{5}.Hex(),
			Status:        job.WorkflowSpecStatusActive,
			WorkflowOwner: "aabbccdd",
			StorageBytes:  meteringTestBytes,
		},
	}
}

// a redelivered delete (after the spec was already removed) must NOT emit a
// second -1. DeleteWorkflowArtifacts maps ErrNoRows→nil, so without gating on
// the pre-read the redelivery would double-emit.
func Test_meterRecords_RedeliveredDeleteEmitsExactlyOneDelta(t *testing.T) {
	t.Parallel()
	wfID := types.WorkflowID{5}
	emitter := &recordingEmitter{}
	store := redeliveredDeleteStore()
	h := newMeteringTestHandler(t, store, newMeteringResourceManager(t, true, emitter))

	deleteEvent := Event{Name: WorkflowDeleted, Data: WorkflowDeletedEvent{WorkflowID: wfID}}
	require.NoError(t, h.Handle(t.Context(), deleteEvent))
	require.NoError(t, h.Handle(t.Context(), deleteEvent)) // redelivery

	records := emitter.Records()
	require.Len(t, records, 2, "redelivered delete must not emit a second -1 pair")
	wantCountID := resourcemanager.EventID("workflow-spec-delete", wfID.Hex())
	requireSpecDelta(t, records[0], "-1", wfID.Hex(), wantCountID)
	requireSpecBytesDelta(t, records[1], "-12", wfID.Hex(), meteringBytesEventID(wantCountID))
}

// Pause and activate are level-neutral: neither emits a metering delta on
// either billing unit. The register pair (count + bytes) stays as the sole
// records through a pause→activate cycle.
func Test_meterRecords_PauseActivateCycleIsLevelNeutral(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{persistUpserts: true}, newMeteringResourceManager(t, true, emitter))

	wfID := meteringwfID
	createdAt := uint64(123)
	wantEventID := resourcemanager.EventID("workflow-spec-register", wfID.Hex(), strconv.FormatUint(createdAt, 10))
	payload := WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    wfID,
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     createdAt,
	}

	require.NoError(t, h.workflowRegisteredEvent(t.Context(), payload))
	require.Len(t, emitter.Records(), 2)
	requireSpecDelta(t, emitter.Records()[0], "1", wfID.Hex(), wantEventID)
	requireSpecBytesDelta(t, emitter.Records()[1], "12", wfID.Hex(), meteringBytesEventID(wantEventID))

	require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: wfID}))
	require.Len(t, emitter.Records(), 2, "pause must not emit a delta")

	require.NoError(t, h.workflowActivatedEvent(t.Context(), WorkflowActivatedEvent(payload)))
	require.Len(t, emitter.Records(), 2, "activate-after-pause must not emit a delta")
}

// OrgId is resolved from the workflow owner and stamped on every emitted
// record, on both billing units.
func Test_meterRecords_OrgIdResolvedOnRecords(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	orgR := &fakeOrgResolver{orgs: map[string]string{"aabbccdd": "org-42"}}
	h := newMeteringTestHandlerWithOrg(t, &stubWorkflowArtifactsStore{}, newMeteringResourceManager(t, true, emitter), orgR)

	wfID := types.WorkflowID{50}
	require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:        WorkflowStatusPaused,
		WorkflowID:    wfID,
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     1,
	}))

	records := emitter.Records()
	require.Len(t, records, 2)
	require.Len(t, records[0].Utilizations, 1)
	require.Len(t, records[1].Utilizations, 1)
	assert.Equal(t, "org-42", records[0].Utilizations[0].OrgId, "OrgId must be resolved and stamped on the count record")
	assert.Equal(t, "org-42", records[1].Utilizations[0].OrgId, "OrgId must be resolved and stamped on the storage record")
}

// don_id is folded into the metering identity once resolved and stamped on both
// emitted records and snapshot entries.
func Test_meterRecords_DonIDOnRecordAndSnapshot(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	wfID := types.WorkflowID{60}
	store := &stubWorkflowArtifactsStore{
		specs: []*job.WorkflowSpec{{WorkflowID: wfID.Hex(), WorkflowOwner: "aabbccdd"}},
	}
	h := newMeteringTestHandler(t, store, newMeteringResourceManager(t, true, emitter))

	// Simulate a resolved workflow DON id (as SetWorkflowDon would store).
	resolvedDon := "7"
	h.specMeter.resolvedDonID.Store(&resolvedDon)

	require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:        WorkflowStatusPaused,
		WorkflowID:    wfID,
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     1,
	}))
	records := emitter.Records()
	require.Len(t, records, 2)
	require.NotNil(t, records[0].Identity.GetDon())
	assert.Equal(t, "7", records[0].Identity.GetDon().GetDonId(), "record identity must carry resolved don_id")
	assert.Equal(t, "7", records[1].Identity.GetDon().GetDonId(), "storage record identity must carry resolved don_id")

	entries := h.specMeter.GetUtilization(t.Context())
	require.Len(t, entries, 1)
	require.NotNil(t, entries[0].Identity.Don)
	assert.Equal(t, "7", entries[0].Identity.Don.DonID, "snapshot identity must carry resolved don_id")
}

// a transient DB error on GetWorkflowSpec (not a missing row) must surface
// an error and emit no +1, so the event is retried instead of creating a
// spurious spec + delta.
func Test_meterRecords_TransientDBErrorOnGetSpecEmitsNoDelta(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{getSpecErr: errors.New("db connection lost")},
		newMeteringResourceManager(t, true, emitter))

	err := h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:        WorkflowStatusPaused,
		WorkflowID:    types.WorkflowID{70},
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     1,
	})
	require.Error(t, err)
	assert.Empty(t, emitter.Records(), "transient DB error must not take the new-spec path or emit a +1")
}

// A full register→pause→activate→delete cycle emits exactly four records:
// a count+bytes +pair at register and a count+bytes -pair at delete. Pause and
// activate emit nothing. The delete event_id carries the registered_at that
// flowed insert→RETURNING, proving the generation-scoped id is DON-consistent.
// The delete -bytes pair releases exactly the register +bytes value, so both
// delta streams self-balance across the full lifecycle.
func Test_meterRecords_FullLifecycleEmitsExactlyFourRecords(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{persistUpserts: true}, newMeteringResourceManager(t, true, emitter))

	wfID := meteringwfID
	createdAt := uint64(123)
	payload := WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    wfID,
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     createdAt,
	}

	require.NoError(t, h.workflowRegisteredEvent(t.Context(), payload))
	require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: wfID}))
	require.NoError(t, h.workflowActivatedEvent(t.Context(), WorkflowActivatedEvent(payload)))
	require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd"))

	records := emitter.Records()
	require.Len(t, records, 4, "exactly four records: +pair at register, -pair at delete")
	wantRegisterID := resourcemanager.EventID("workflow-spec-register", wfID.Hex(), strconv.FormatUint(createdAt, 10))
	requireSpecDelta(t, records[0], "1", wfID.Hex(), wantRegisterID)
	requireSpecBytesDelta(t, records[1], "12", wfID.Hex(), meteringBytesEventID(wantRegisterID))
	wantDeleteID := resourcemanager.EventID("workflow-spec-delete", wfID.Hex(), strconv.FormatUint(createdAt, 10))
	requireSpecDelta(t, records[2], "-1", wfID.Hex(), wantDeleteID)
	requireSpecBytesDelta(t, records[3], "-12", wfID.Hex(), meteringBytesEventID(wantDeleteID))
}

// A transient delete error returns before emission (no -1); a successful retry
// emits exactly one -pair. This fixes the lost-−1 hazard from the old pre-read gate.
func Test_meterRecords_TransientDeleteErrorThenRetryEmitsOnce(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	store := &stubWorkflowArtifactsStore{
		persistUpserts: true,
		deleteErr:      assert.AnError,
	}
	h := newMeteringTestHandler(t, store, newMeteringResourceManager(t, true, emitter))

	wfID := types.WorkflowID{77}
	wantRegisterID := resourcemanager.EventID("workflow-spec-register", wfID.Hex(), "1")
	require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:        WorkflowStatusPaused,
		WorkflowID:    wfID,
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     1,
	}))
	require.Len(t, emitter.Records(), 2, "register emits +pair")
	requireSpecDelta(t, emitter.Records()[0], "1", wfID.Hex(), wantRegisterID)
	requireSpecBytesDelta(t, emitter.Records()[1], "12", wfID.Hex(), meteringBytesEventID(wantRegisterID))

	err := h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd")
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, emitter.Records()[2:], "failed delete must not emit a -pair")

	store.deleteErr = nil
	require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd"))

	records := emitter.Records()
	require.Len(t, records, 4, "exactly one -pair after successful retry")
	wantDeleteID := resourcemanager.EventID("workflow-spec-delete", wfID.Hex(), "1")
	requireSpecDelta(t, records[2], "-1", wfID.Hex(), wantDeleteID)
	requireSpecBytesDelta(t, records[3], "-12", wfID.Hex(), meteringBytesEventID(wantDeleteID))
}

// A legacy row with RegisteredAt == 0 produces a delete event_id with no
// timestamp part (the fallback for pre-migration rows).
func Test_meterRecords_LegacyRowDeleteUsesFallbackID(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	h := newMeteringTestHandler(t, &stubWorkflowArtifactsStore{
		spec: &job.WorkflowSpec{
			WorkflowID:    types.WorkflowID{88}.Hex(),
			Status:        job.WorkflowSpecStatusActive,
			WorkflowOwner: "aabbccdd",
			RegisteredAt:  0,
			StorageBytes:  meteringTestBytes,
		},
	}, newMeteringResourceManager(t, true, emitter))

	require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: types.WorkflowID{88}}, "aabbccdd"))

	records := emitter.Records()
	require.Len(t, records, 2)
	wantCountID := resourcemanager.EventID("workflow-spec-delete", types.WorkflowID{88}.Hex())
	requireSpecDelta(t, records[0], "-1", types.WorkflowID{88}.Hex(), wantCountID)
	requireSpecBytesDelta(t, records[1], "-12", types.WorkflowID{88}.Hex(), meteringBytesEventID(wantCountID))
}

// Identical event sequences with RM nil / disabled / erroring emitter produce
// identical row+engine outcomes and zero handler errors (fail-open equivalence).
func Test_meterRecords_FailOpenEquivalence(t *testing.T) {
	t.Parallel()
	wfID := meteringwfID
	owner := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	payload := WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    wfID,
		WorkflowOwner: owner,
		WorkflowName:  "wf-name",
		CreatedAt:     1,
	}
	delEvt := WorkflowDeletedEvent{WorkflowID: wfID}

	runCycle := func(t *testing.T, rm *resourcemanager.ResourceManager, emitter resourcemanager.Emitter) (*stubWorkflowArtifactsStore, *eventHandler) {
		store := &stubWorkflowArtifactsStore{persistUpserts: true}
		h := newMeteringTestHandler(t, store, rm)
		require.NoError(t, h.workflowRegisteredEvent(t.Context(), payload))
		require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: wfID}))
		require.NoError(t, h.workflowActivatedEvent(t.Context(), WorkflowActivatedEvent(payload)))
		require.NoError(t, h.workflowDeletedEvent(t.Context(), delEvt, "aabbccdd"))
		return store, h
	}

	// No spec meter (metering disabled at construction)
	nilStore, nilH := runCycle(t, nil, nil)
	assert.Nil(t, nilH.specMeter)
	assert.Nil(t, nilStore.spec, "spec should be deleted by the cycle")

	// RM disabled
	disabledEmitter := &recordingEmitter{}
	disabledRM := newMeteringResourceManager(t, false, disabledEmitter)
	disabledStore, disabledH := runCycle(t, disabledRM, disabledEmitter)
	assert.NotNil(t, disabledH.specMeter)
	assert.Empty(t, disabledEmitter.Records(), "disabled RM emits nothing")
	assert.Nil(t, disabledStore.spec, "spec should be deleted by the cycle")

	// RM enabled but erroring emitter
	errEmitter := &recordingEmitter{err: errors.New("beholder unavailable")}
	errRM := newMeteringResourceManager(t, true, errEmitter)
	errStore, errH := runCycle(t, errRM, errEmitter)
	assert.NotNil(t, errH.specMeter)
	assert.Empty(t, errEmitter.Records(), "erroring emitter stores nothing")
	assert.Nil(t, errStore.spec, "spec should be deleted by the cycle")
}

// The orphan sweep releases a paused tombstone (no engine) with exactly one
// -pair carrying the generation delete-ids. This is the sweep's new
// load-bearing case: workflows deleted on-chain while paused have a tombstone
// but no engine. The sweep dispatches an ordinary WorkflowDeleted event
// through Handle — the same path as reconciliation-generated deletes. The
// -bytes record releases the ORIGINAL registered byte count even though pause
// cleared the payload: the persisted storage_bytes survives tombstoning.
func Test_meterRecords_SweepReleasesPausedTombstone(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	store := &stubWorkflowArtifactsStore{persistUpserts: true}
	h := newMeteringTestHandler(t, store, newMeteringResourceManager(t, true, emitter))

	wfID := meteringwfID
	createdAt := uint64(456)
	require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    wfID,
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     createdAt,
	}))
	require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: wfID}))
	// The tombstone must keep the registered byte count despite the cleared
	// payload: that is what makes the delete-time -bytes delta exact.
	require.NotNil(t, store.spec)
	require.Empty(t, store.spec.Workflow)
	require.Empty(t, store.spec.Config)
	require.Equal(t, meteringTestBytes, store.spec.StorageBytes)
	// After pause the engine is popped; the tombstone has no engine → sweep path.
	require.NoError(t, h.Handle(t.Context(), Event{Name: WorkflowDeleted, Data: WorkflowDeletedEvent{WorkflowID: wfID}}))

	records := emitter.Records()
	require.Len(t, records, 4, "register +pair and sweep -pair")
	wantRegisterID := resourcemanager.EventID("workflow-spec-register", wfID.Hex(), strconv.FormatUint(createdAt, 10))
	requireSpecDelta(t, records[0], "1", wfID.Hex(), wantRegisterID)
	requireSpecBytesDelta(t, records[1], "12", wfID.Hex(), meteringBytesEventID(wantRegisterID))
	wantDeleteID := resourcemanager.EventID("workflow-spec-delete", wfID.Hex(), strconv.FormatUint(createdAt, 10))
	requireSpecDelta(t, records[2], "-1", wfID.Hex(), wantDeleteID)
	requireSpecBytesDelta(t, records[3], "-12", wfID.Hex(), meteringBytesEventID(wantDeleteID))
}

// A paused tombstone keeps reporting its registered storage level in snapshots
// (billing covers the registration lifetime), and the eventual delete releases
// exactly those bytes: +N at register, 0 at pause/activate, -N at delete.
func Test_meterRecords_TombstoneSnapshotAndDeleteReleaseRegisteredBytes(t *testing.T) {
	t.Parallel()
	emitter := &recordingEmitter{}
	// specs simulates what the DB list query returns for the paused tombstone:
	// identity columns plus storage_bytes (the payload columns are cleared).
	store := &stubWorkflowArtifactsStore{
		persistUpserts: true,
		specs: []*job.WorkflowSpec{
			{WorkflowID: meteringwfID.Hex(), WorkflowOwner: "aabbccdd", StorageBytes: meteringTestBytes},
		},
	}
	h := newMeteringTestHandler(t, store, newMeteringResourceManager(t, true, emitter))

	wfID := meteringwfID
	createdAt := uint64(789)
	require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    wfID,
		WorkflowOwner: []byte{0xaa, 0xbb, 0xcc, 0xdd},
		WorkflowName:  "wf-name",
		CreatedAt:     createdAt,
	}))
	require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: wfID}))

	// While paused, the snapshot path still reports both dimensions of the
	// held registration: count 1 and the registered byte count.
	entries := h.specMeter.GetUtilization(t.Context())
	require.Len(t, entries, 1)
	require.Len(t, entries[0].Utilizations, 2)
	assert.Equal(t, "operations", entries[0].Utilizations[0].ResourceType)
	assert.Equal(t, "1", entries[0].Utilizations[0].Value)
	assert.Equal(t, "storage_bytes", entries[0].Utilizations[1].ResourceType)
	assert.Equal(t, "12", entries[0].Utilizations[1].Value)

	require.NoError(t, h.workflowDeletedEvent(t.Context(), WorkflowDeletedEvent{WorkflowID: wfID}, "aabbccdd"))

	records := emitter.Records()
	require.Len(t, records, 4)
	wantRegisterID := resourcemanager.EventID("workflow-spec-register", wfID.Hex(), strconv.FormatUint(createdAt, 10))
	requireSpecDelta(t, records[0], "1", wfID.Hex(), wantRegisterID)
	requireSpecBytesDelta(t, records[1], "12", wfID.Hex(), meteringBytesEventID(wantRegisterID))
	wantDeleteID := resourcemanager.EventID("workflow-spec-delete", wfID.Hex(), strconv.FormatUint(createdAt, 10))
	requireSpecDelta(t, records[2], "-1", wfID.Hex(), wantDeleteID)
	requireSpecBytesDelta(t, records[3], "-12", wfID.Hex(), meteringBytesEventID(wantDeleteID))
}
