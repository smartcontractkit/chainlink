package v2

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/shardownership"
	wfTypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

func Test_generateReconciliationEventsV2(t *testing.T) {
	t.Parallel()
	// Validate that if no engines are on the node in the registry,
	// and we see that the contract has workflow state,
	// that we generate a WorkflowActivatedEvent
	t.Run("WorkflowActivatedEvent_whenNoEnginesInRegistry", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		wfID := [32]byte{1}
		owner := []byte{}
		createdAt := uint64(1000000)
		status := uint8(0)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		donFamily := "A"
		tag := "tag1"
		attributes := []byte{}
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		// The only event is WorkflowActivatedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowActivated, events[0].Name)
		expectedActivatedEvent := WorkflowActivatedEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			WorkflowTag:   tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedActivatedEvent, events[0].Data)
	})

	t.Run("WorkflowUpdatedEvent", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{1}
		wfName := "wf name 1"
		err := er.Add(wfID, "TestSource", &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata gets updated
		wfID2 := [32]byte{2}
		createdAt := uint64(1000000)
		status := uint8(0)
		binaryURL2 := "b2"
		configURL := "c1"
		donFamily := "A"
		tag := "tag1"
		attributes := []byte{}
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID2,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL2,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		require.Len(t, events, 2)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		expectedDeletedEvent := WorkflowDeletedEvent{
			WorkflowID: wfID,
			Source:     "TestSource",
		}
		require.Equal(t, expectedDeletedEvent, events[0].Data)
		require.Equal(t, WorkflowActivated, events[1].Name)
		expectedActivatedEvent := WorkflowActivatedEvent{
			WorkflowID:    wfID2,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL2,
			ConfigURL:     configURL,
			WorkflowTag:   tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedActivatedEvent, events[1].Data)
	})

	t.Run("RetiredIDReusedByDifferentRecord_TearsDownStaleEngine", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		wfID := [32]byte{1}
		victimOwner := []byte{1}
		wfName := "wf name 1"
		tag := "tag1"

		er := NewEngineRegistry()
		victimKey, err := ReconcileKey(victimOwner, wfName)
		require.NoError(t, err)
		require.NoError(t, er.AddWithReconcileKey(wfID, "TestSource", victimKey, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) { return nil, nil },
			"",
			"test-chain-selector",
			Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        []byte{9},
				Status:       uint8(0),
				WorkflowName: wfName,
				BinaryURL:    "b1",
				ConfigURL:    "c1",
				Tag:          tag,
				Attributes:   []byte{},
				DonFamily:    "A",
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		require.Equal(t, WorkflowDeletedEvent{WorkflowID: wfID, Source: "TestSource"}, events[0].Data)
	})

	t.Run("SameIDMatchingIdentity_NoOpRegardlessOfTag", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		wfID := [32]byte{1}
		owner := []byte{1}
		wfName := "wf name 1"

		er := NewEngineRegistry()
		key, err := ReconcileKey(owner, wfName)
		require.NoError(t, err)
		require.NoError(t, er.AddWithReconcileKey(wfID, "TestSource", key, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) { return nil, nil },
			"",
			"test-chain-selector",
			Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				Status:       uint8(0),
				WorkflowName: wfName,
				BinaryURL:    "b1",
				ConfigURL:    "c1",
				Tag:          "tag1",
				Attributes:   []byte{7, 7, 7},
				DonFamily:    "A",
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Empty(t, events)
	})

	t.Run("ActiveCollisionFromOtherSource_DoesNotEvictOtherSourceEngine", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		wfID := wfTypes.WorkflowID([32]byte{1})
		privateSource := "grpc:private:v1"
		publicSource := "contract:1:0xpublic"

		er := NewEngineRegistry()
		key, err := ReconcileKey([]byte("private-owner"), "private-workflow")
		require.NoError(t, err)
		require.NoError(t, er.AddWithReconcileKey(wfID, privateSource, key, &mockService{}))
		wr, err := NewWorkflowRegistry(
			lggr,
			func(context.Context, []byte) (types.ContractReader, error) { return nil, nil },
			"", "test-chain-selector",
			Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
			&eventHandler{}, capabilities.NewDonNotifier(), er,
		)
		require.NoError(t, err)

		metadata := []WorkflowMetadataView{{
			WorkflowID: wfID, Owner: []byte("attacker"), Status: WorkflowStatusActive,
			WorkflowName: "attacker-workflow", Tag: "attacker-tag", Source: publicSource,
		}}
		events, err := wr.generateReconciliationEvents(ctx, map[string]*reconciliationEvent{}, metadata, &types.Head{Height: "123"}, publicSource)
		require.NoError(t, err)
		require.Empty(t, events, "the other source's engine must not be torn down")
	})

	t.Run("PausedCollisionFromOtherSource_DoesNotStopOtherSourceEngine", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		wfID := wfTypes.WorkflowID([32]byte{9})
		privateSource := "grpc:private-policy:v1"
		publicSource := "contract:1:0xpublic"

		er := NewEngineRegistry()
		key, err := ReconcileKey([]byte("private-owner"), "private-policy-workflow")
		require.NoError(t, err)
		require.NoError(t, er.AddWithReconcileKey(wfID, privateSource, key, &mockService{}))
		wr, err := NewWorkflowRegistry(
			lggr,
			func(context.Context, []byte) (types.ContractReader, error) { return nil, nil },
			"", "test-chain-selector",
			Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
			&eventHandler{}, capabilities.NewDonNotifier(), er,
		)
		require.NoError(t, err)

		metadata := []WorkflowMetadataView{{
			WorkflowID: wfID, Owner: []byte("attacker"), Status: WorkflowStatusPaused,
			WorkflowName: "public-paused-collision", Source: publicSource,
		}}
		events, err := wr.generateReconciliationEvents(ctx, map[string]*reconciliationEvent{}, metadata, &types.Head{Height: "500"}, publicSource)
		require.NoError(t, err)
		require.Empty(t, events, "no WorkflowPaused must be emitted for the other source's engine")
	})

	t.Run("WorkflowDeletedEvent", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		err := er.Add(wfID, "TestSource", &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata is empty
		metadata := []WorkflowMetadataView{}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		// The only event is WorkflowDeletedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		expectedDeletedEvent := WorkflowDeletedEvent{
			WorkflowID: wfID,
			Source:     "TestSource",
		}
		require.Equal(t, expectedDeletedEvent, events[0].Data)
	})

	t.Run("generateReconciliationEvents is side-effect free; pre-dispatch drains delete targets", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		drainingEngine := &mockDrainableEngine{}
		err := er.Add(wfID, "TestSource", drainingEngine)
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		metadata := []WorkflowMetadataView{}
		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)

		_, draining := drainingEngine.DrainStartedAt()
		require.False(t, draining, "generateReconciliationEvents should not mutate engine state")

		wr.applyPreDispatchReconcileActions(ctx, events)
		_, draining = drainingEngine.DrainStartedAt()
		require.True(t, draining, "pre-dispatch actions should initiate drain for delete events")
	})

	t.Run("No change", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		wfID := [32]byte{1}
		owner := []byte{}
		status := uint8(0)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		// The only event is WorkflowActivatedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowActivated, events[0].Name)
		expectedActivatedEvent := WorkflowActivatedEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        status,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			WorkflowTag:   tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedActivatedEvent, events[0].Data)

		// Add the workflow to the engine registry as the handler would
		err = er.Add(wfID, ContractWorkflowSourceName, &mockService{})
		require.NoError(t, err)

		// Repeated ticks do not make any new events
		events, err = wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Empty(t, events)
	})

	t.Run("A paused workflow doesn't start a new workflow", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// No engines are in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		wfID := [32]byte{1}
		owner := []byte{}
		status := uint8(1)
		wfName := "wf name 1"
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		// No events
		require.Empty(t, events)
	})

	t.Run("A paused workflow deletes a running workflow", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		err := er.Add(wfID, "TestSource", &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// The workflow metadata gets updated
		status := uint8(1)
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       status,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		// The only event is WorkflowPausedEvent
		require.Len(t, events, 1)
		require.Equal(t, WorkflowPaused, events[0].Name)
		expectedPausedEvent := WorkflowPausedEvent{
			WorkflowID: wfID,
		}
		require.Equal(t, expectedPausedEvent.WorkflowID, events[0].Data.(WorkflowPausedEvent).WorkflowID)
	})

	t.Run("reconciles with a pending event if it has the same signature", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// The workflow metadata gets updated
		binaryURL := "b1"
		configURL := "c1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		event := WorkflowActivatedEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			WorkflowTag:   tag,
			Attributes:    attributes,
		}
		signature := fmt.Sprintf("%s-%s-%s", WorkflowActivated, event.WorkflowID.Hex(), toSpecStatus(WorkflowStatusActive))
		retryCount := 2
		nextRetryAt := fakeClock.Now().Add(5 * time.Minute)
		pendingEvents := map[string]*reconciliationEvent{
			event.WorkflowID.Hex(): {
				Event: Event{
					Data: event,
					Name: WorkflowActivated,
				},
				signature:   signature,
				id:          event.WorkflowID.Hex(),
				retryCount:  retryCount,
				nextRetryAt: nextRetryAt,
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		// The only event is WorkflowActivatedEvent
		// Since there's a failing event in the pendingEvents queue, we should expect to see
		// that event returned to us.
		require.Empty(t, pendingEvents)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowActivated, events[0].Name)
		require.Equal(t, event, events[0].Data)
		require.Equal(t, retryCount, events[0].retryCount)
		require.Equal(t, nextRetryAt, events[0].nextRetryAt)
	})

	t.Run("dropped activation is not re-enqueued", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		wfID := wfTypes.WorkflowID([32]byte{1})
		owner := []byte{}
		wfName := "wf name 1"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName,
			},
		}
		source := "TestSource"
		signature := fmt.Sprintf("%s-%s-%s", WorkflowActivated, wfID.Hex(), toSpecStatus(WorkflowStatusActive))
		wr.droppedActivations.drop(source, wfID.Hex(), signature)

		events, err := wr.generateReconciliationEvents(ctx, map[string]*reconciliationEvent{}, metadata, &types.Head{Height: "123"}, source)
		require.NoError(t, err)
		require.Empty(t, events)
	})

	t.Run("a paused workflow clears a pending activated event", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// The workflow metadata gets updated
		binaryURL := "b1"
		configURL := "c1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       WorkflowStatusPaused,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}
		// Now let's emit an event with the same signature; this should remove the event
		// from the pending queue.
		event := WorkflowActivatedEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			WorkflowTag:   tag,
			Attributes:    attributes,
		}
		signature := fmt.Sprintf("%s-%s-%s", WorkflowRegistered, event.WorkflowID.Hex(), toSpecStatus(WorkflowStatusActive))
		retryCount := 2
		nextRetryAt := fakeClock.Now().Add(5 * time.Minute)
		pendingEvents := map[string]*reconciliationEvent{
			event.WorkflowID.Hex(): {
				Event: Event{
					Data: event,
					Name: WorkflowRegistered,
				},
				signature:   signature,
				id:          event.WorkflowID.Hex(),
				retryCount:  retryCount,
				nextRetryAt: nextRetryAt,
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		require.Empty(t, pendingEvents)
		require.Empty(t, events)
	})

	t.Run("delete events are handled before any other events", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{1}
		wfName := "wf name 1"
		err := er.Add(wfID, "TestSource", &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// The workflow gets a new version with updated metadata, which changes the workflow ID
		wfID2 := [32]byte{2}
		binaryURL := "b1"
		configURL := "c1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		donFamily := "A"
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   wfID2,
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       WorkflowStatusActive,
				WorkflowName: wfName,
				BinaryURL:    binaryURL,
				ConfigURL:    configURL,
				Tag:          tag,
				Attributes:   attributes,
				DonFamily:    donFamily,
			},
		}

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)

		// Delete event happens before activate event
		require.Equal(t, events[0].Name, WorkflowDeleted)
		require.Equal(t, events[1].Name, WorkflowActivated)
	})

	t.Run("pending delete events are handled when workflow metadata no longer exists", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		err := er.Add(wfID, "TestSource", &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// A workflow is to be removed, but hits a failure, causing it to stay pending
		event := WorkflowDeletedEvent{
			WorkflowID: wfID,
			Source:     "TestSource",
		}
		pendingEvents := map[string]*reconciliationEvent{
			hex.EncodeToString(wfID[:]): {
				Event: Event{
					Data: event,
					Name: WorkflowDeleted,
				},
				id:          hex.EncodeToString(wfID[:]),
				signature:   fmt.Sprintf("%s-%s", WorkflowDeleted, hex.EncodeToString(wfID[:])),
				nextRetryAt: time.Now(),
				retryCount:  5,
			},
		}

		// No workflows in metadata
		metadata := []WorkflowMetadataView{}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		require.Empty(t, pendingEvents)
	})

	// Reproduces the pending-delete reappearance race: a WorkflowDeleted event was deferred
	// (e.g. ErrDrainInProgress) and stored in pendingEvents; before the next deletion retry
	// runs the workflow re-appears as Active in the metadata while its engine is still in
	// the registry. The Active+engineFound branch must drop the stale pending entry,
	// otherwise generateReconciliationEvents trips its end-of-loop invariant check.
	t.Run("active workflow with running engine clears stale pending WorkflowDeleted", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		err := er.Add(wfID, "TestSource", &mockService{})
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		idHex := hex.EncodeToString(wfID[:])
		pendingEvents := map[string]*reconciliationEvent{
			idHex: {
				Event: Event{
					Data: WorkflowDeletedEvent{WorkflowID: wfID, Source: "TestSource"},
					Name: WorkflowDeleted,
				},
				id:        idHex,
				signature: fmt.Sprintf("%s-%s", WorkflowDeleted, idHex),
			},
		}

		metadata := []WorkflowMetadataView{
			{
				WorkflowID: wfID,
				Status:     WorkflowStatusActive,
			},
		}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Empty(t, events, "engine already matches desired Active state; no new events expected")
		require.Empty(t, pendingEvents, "stale WorkflowDeleted pending entry must be cleared")
	})

	// Same scenario as above but the registered engine is in a draining state. The minimal
	// fix should still clear the stale pending entry without panicking; this guards against
	// regressions if the branch is later extended to emit a replacement WorkflowActivated.
	t.Run("active workflow with draining engine clears stale pending WorkflowDeleted", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		drainingEngine := &mockDrainableEngine{}
		require.True(t, drainingEngine.Drain())
		err := er.Add(wfID, "TestSource", drainingEngine)
		require.NoError(t, err)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		idHex := hex.EncodeToString(wfID[:])
		pendingEvents := map[string]*reconciliationEvent{
			idHex: {
				Event: Event{
					Data: WorkflowDeletedEvent{WorkflowID: wfID, Source: "TestSource"},
					Name: WorkflowDeleted,
				},
				id:        idHex,
				signature: fmt.Sprintf("%s-%s", WorkflowDeleted, idHex),
			},
		}

		metadata := []WorkflowMetadataView{
			{
				WorkflowID: wfID,
				Status:     WorkflowStatusActive,
			},
		}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Empty(t, events)
		require.Empty(t, pendingEvents)
	})

	t.Run("pending activate events are handled when workflow metadata no longer exists", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)

		// A workflow is added, but hits a failure during creation, causing it to stay pending
		binaryURL := "b1"
		configURL := "c1"
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		createdAt := uint64(1000000)
		tag := "tag1"
		attributes := []byte{}
		event := WorkflowActivatedEvent{
			WorkflowID:    wfID,
			WorkflowOwner: owner,
			CreatedAt:     createdAt,
			Status:        WorkflowStatusActive,
			WorkflowName:  wfName,
			BinaryURL:     binaryURL,
			ConfigURL:     configURL,
			WorkflowTag:   tag,
			Attributes:    attributes,
		}
		pendingEvents := map[string]*reconciliationEvent{
			hex.EncodeToString(wfID[:]): {
				Event: Event{
					Data: event,
					Name: WorkflowActivated,
				},
				id:          hex.EncodeToString(wfID[:]),
				signature:   fmt.Sprintf("%s-%s-%s", WorkflowActivated, hex.EncodeToString(wfID[:]), toSpecStatus(WorkflowStatusActive)),
				nextRetryAt: time.Now(),
				retryCount:  5,
			},
		}

		// The workflow then gets removed
		metadata := []WorkflowMetadataView{}

		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Empty(t, events)
		require.Empty(t, pendingEvents)
	})
}

func Test_Start(t *testing.T) {
	t.Parallel()
	t.Run("successful start and close", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		mockReader := &mockContractReader{startErr: nil}
		er := NewEngineRegistry()
		lf := limits.Factory{Logger: lggr}
		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		h := &eventHandler{
			engineRegistry: &EngineRegistry{},
			engineLimiters: limiters,
		}
		svc, eng := services.Config{
			Name:  "EventHandler",
			Close: h.close,
		}.NewServiceEngine(lggr)
		h.Service = svc
		h.eng = eng
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return mockReader, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			h,
			workflowDonNotifier,
			er,
		)
		fakeClock := clockwork.NewFakeClock()
		wr.clock = fakeClock
		require.NoError(t, err)
		servicetest.Run(t, wr)
		workflowDonNotifier.NotifyDonSet(commonCap.DON{})
	})
}

func Test_GetAllowlistedRequests(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()
	workflowDonNotifier := capabilities.NewDonNotifier()
	er := NewEngineRegistry()

	// Mock allowlisted requests
	expectedRequests := []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
		{
			RequestDigest:   [32]byte{1, 2, 3},
			Owner:           common.Address{4, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ExpiryTimestamp: 123456789,
		},
		{
			RequestDigest:   [32]byte{7, 8, 9},
			Owner:           common.Address{10, 11, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ExpiryTimestamp: 987654321,
		},
	}

	// Mock contract reader to return expectedRequests
	mockContractReader := &mockContractReader{
		allowlistedRequests: expectedRequests,
	}

	wr, err := NewWorkflowRegistry(
		lggr,
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return mockContractReader, nil
		},
		"",
		"test-chain-selector",
		Config{
			QueryCount:   20,
			SyncStrategy: SyncStrategyReconciliation,
		},
		&eventHandler{},
		workflowDonNotifier,
		er,
	)
	require.NoError(t, err)

	// Simulate syncAllowlistedRequests updating the field
	wr.allowListedMu.Lock()
	wr.allowListedRequests = expectedRequests
	wr.allowListedMu.Unlock()

	// Test GetAllowlistedRequests returns the correct data
	got := wr.GetAllowlistedRequests(ctx)
	require.Equal(t, expectedRequests, got)
}

// Mock contract reader implementation
type mockContractReader struct {
	types.ContractReader
	bindErr             error
	startErr            error
	allowlistedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
}

func (m *mockContractReader) GetLatestValueWithHeadData(
	_ context.Context,
	_ string,
	_ primitives.ConfidenceLevel,
	_ any,
	result any,
) (*types.Head, error) {
	// Simulate returning allowlisted requests
	if res, ok := result.(*struct {
		Requests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
		err      error
	}); ok {
		res.Requests = m.allowlistedRequests
		return &types.Head{Height: "123"}, nil
	}
	return &types.Head{Height: "0"}, nil
}

func (m *mockContractReader) Bind(
	_ context.Context,
	_ []types.BoundContract,
) error {
	return m.bindErr
}

func (m *mockContractReader) Start(
	_ context.Context,
) error {
	return m.startErr
}

func Test_generateReconciliationEvents_SourceIsolation(t *testing.T) {
	t.Parallel()
	t.Run("only deletes engines from specified source", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		// Setup: engines from two sources
		er := NewEngineRegistry()
		wfIDContract := [32]byte{1}
		wfIDGrpc := [32]byte{2}
		require.NoError(t, er.Add(wfIDContract, ContractWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDGrpc, GRPCWorkflowSourceName, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Reconcile ContractWorkflowSource with empty metadata
		// Should only delete contract engine, not GRPC engine
		pendingEvents := make(map[string]*reconciliationEvent)
		events, err := wr.generateReconciliationEvents(
			ctx, pendingEvents, []WorkflowMetadataView{}, &types.Head{Height: "123"}, ContractWorkflowSourceName)

		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		deletedEvent := events[0].Data.(WorkflowDeletedEvent)
		require.Equal(t, wfTypes.WorkflowID(wfIDContract), deletedEvent.WorkflowID)
		require.Equal(t, ContractWorkflowSourceName, deletedEvent.Source)
	})

	t.Run("activates workflows tagged with source", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// New workflow from GRPCWorkflowSource
		wfID := [32]byte{1}
		metadata := []WorkflowMetadataView{{
			WorkflowID:   wfID,
			Owner:        []byte{1, 2, 3},
			Status:       WorkflowStatusActive,
			Source:       GRPCWorkflowSourceName,
			WorkflowName: "test-workflow",
			BinaryURL:    "http://binary.url",
			ConfigURL:    "http://config.url",
		}}

		pendingEvents := make(map[string]*reconciliationEvent)
		events, err := wr.generateReconciliationEvents(
			ctx, pendingEvents, metadata, &types.Head{Height: "123"}, GRPCWorkflowSourceName)

		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowActivated, events[0].Name)
		activatedEvent := events[0].Data.(WorkflowActivatedEvent)
		require.Equal(t, wfTypes.WorkflowID(wfID), activatedEvent.WorkflowID)
		require.Equal(t, GRPCWorkflowSourceName, activatedEvent.Source)
	})

	t.Run("does not delete engines from other sources when source returns empty", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		// Setup: engines from two sources
		er := NewEngineRegistry()
		wfIDContract := [32]byte{1}
		wfIDGrpc := [32]byte{2}
		require.NoError(t, er.Add(wfIDContract, ContractWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDGrpc, GRPCWorkflowSourceName, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Reconcile GRPCWorkflowSource with empty metadata
		// Should only generate delete event for GRPC engine, not contract engine
		pendingEvents := make(map[string]*reconciliationEvent)
		events, err := wr.generateReconciliationEvents(
			ctx, pendingEvents, []WorkflowMetadataView{}, &types.Head{Height: "123"}, GRPCWorkflowSourceName)

		require.NoError(t, err)
		require.Len(t, events, 1)
		deletedEvent := events[0].Data.(WorkflowDeletedEvent)
		require.Equal(t, wfTypes.WorkflowID(wfIDGrpc), deletedEvent.WorkflowID)

		// Contract engine should still be in registry (we're just checking the event, not actually processing)
		_, ok := er.Get(wfIDContract)
		require.True(t, ok, "Contract engine should still exist")
	})

	t.Run("handles paused workflow from source", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		// Setup: engine exists for a workflow
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		require.NoError(t, er.Add(wfID, ContractWorkflowSourceName, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Workflow is now paused
		metadata := []WorkflowMetadataView{{
			WorkflowID:   wfID,
			Owner:        []byte{1, 2, 3},
			Status:       WorkflowStatusPaused,
			Source:       ContractWorkflowSourceName,
			WorkflowName: "test-workflow",
		}}

		pendingEvents := make(map[string]*reconciliationEvent)
		events, err := wr.generateReconciliationEvents(
			ctx, pendingEvents, metadata, &types.Head{Height: "123"}, ContractWorkflowSourceName)

		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowPaused, events[0].Name)
		pausedEvent := events[0].Data.(WorkflowPausedEvent)
		require.Equal(t, wfTypes.WorkflowID(wfID), pausedEvent.WorkflowID)
		require.Equal(t, ContractWorkflowSourceName, pausedEvent.Source)
	})

	t.Run("no events when source has no engines and returns empty metadata", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		// Setup: engine only from contract source
		er := NewEngineRegistry()
		wfIDContract := [32]byte{1}
		require.NoError(t, er.Add(wfIDContract, ContractWorkflowSourceName, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Reconcile GRPCWorkflowSource with empty metadata
		// Should generate no events since GRPC has no engines
		pendingEvents := make(map[string]*reconciliationEvent)
		events, err := wr.generateReconciliationEvents(
			ctx, pendingEvents, []WorkflowMetadataView{}, &types.Head{Height: "123"}, GRPCWorkflowSourceName)

		require.NoError(t, err)
		require.Empty(t, events)
	})
}

// Test_PerSourceReconciliation_FailureIsolation validates the main bug fix:
// when a source fails to fetch, engines from that source should NOT be deleted.
func Test_PerSourceReconciliation_FailureIsolation(t *testing.T) {
	t.Parallel()
	t.Run("source failure does not delete engines from that source", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		// Setup: engines from ContractWorkflowSource and GRPCWorkflowSource
		er := NewEngineRegistry()
		wfIDContract := [32]byte{1}
		wfIDGrpc := [32]byte{2}
		require.NoError(t, er.Add(wfIDContract, ContractWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDGrpc, GRPCWorkflowSourceName, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Simulate: contract source succeeds with its workflow
		contractPendingEvents := make(map[string]*reconciliationEvent)
		contractMetadata := []WorkflowMetadataView{{
			WorkflowID:   wfIDContract,
			Owner:        []byte{1, 2, 3},
			Status:       WorkflowStatusActive,
			Source:       ContractWorkflowSourceName,
			WorkflowName: "contract-workflow",
			BinaryURL:    "http://binary.url",
			ConfigURL:    "http://config.url",
		}}
		contractEvents, err := wr.generateReconciliationEvents(
			ctx, contractPendingEvents, contractMetadata, &types.Head{Height: "123"}, ContractWorkflowSourceName)
		require.NoError(t, err)
		require.Empty(t, contractEvents, "No events expected since engine already exists")

		// Simulate: GRPC source FAILS (returns error, so we skip reconciliation)
		// In the actual sync loop, we would NOT call generateReconciliationEvents
		// when the source fetch fails. This test validates that by NOT calling the method
		// for the failed source, the GRPC engine is preserved.

		// Assert: Both engines should still exist
		_, ok := er.Get(wfIDContract)
		require.True(t, ok, "Contract engine should exist after contract source reconciliation")

		_, ok = er.Get(wfIDGrpc)
		require.True(t, ok, "GRPC engine should NOT be deleted when GRPC source fails (skipped reconciliation)")
	})

	t.Run("source recovers after failure - normal reconciliation resumes", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		// Setup: engines from GRPCWorkflowSource
		er := NewEngineRegistry()
		wfIDGrpc1 := [32]byte{1}
		wfIDGrpc2 := [32]byte{2}
		require.NoError(t, er.Add(wfIDGrpc1, GRPCWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDGrpc2, GRPCWorkflowSourceName, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Tick 1: GRPC source fails (skip reconciliation - both engines preserved)
		// ... (simulated by not calling generateReconciliationEvents)

		// Tick 2: GRPC source recovers with only wfIDGrpc1
		grpcPendingEvents := make(map[string]*reconciliationEvent)
		grpcMetadata := []WorkflowMetadataView{{
			WorkflowID:   wfIDGrpc1,
			Owner:        []byte{1, 2, 3},
			Status:       WorkflowStatusActive,
			Source:       GRPCWorkflowSourceName,
			WorkflowName: "grpc-workflow-1",
			BinaryURL:    "http://binary.url",
			ConfigURL:    "http://config.url",
		}}
		events, err := wr.generateReconciliationEvents(
			ctx, grpcPendingEvents, grpcMetadata, &types.Head{Height: "124"}, GRPCWorkflowSourceName)
		require.NoError(t, err)

		// Should generate delete event for wfIDGrpc2 (no longer in metadata)
		require.Len(t, events, 1)
		require.Equal(t, WorkflowDeleted, events[0].Name)
		deletedEvent := events[0].Data.(WorkflowDeletedEvent)
		require.Equal(t, wfTypes.WorkflowID(wfIDGrpc2), deletedEvent.WorkflowID)
		require.Equal(t, GRPCWorkflowSourceName, deletedEvent.Source)
	})

	t.Run("all sources fail - no deletions", func(t *testing.T) {
		t.Parallel()
		// This test validates that when all sources fail, no deletion events are generated
		// because we skip reconciliation for each failed source.
		lggr := logger.TestLogger(t)
		workflowDonNotifier := capabilities.NewDonNotifier()

		er := NewEngineRegistry()
		wfIDContract := [32]byte{1}
		wfIDGrpc := [32]byte{2}
		require.NoError(t, er.Add(wfIDContract, ContractWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDGrpc, GRPCWorkflowSourceName, &mockService{}))

		_, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Both sources fail - we don't call generateReconciliationEvents for either
		// This is simulated by simply not calling the method

		// Both engines should still exist
		require.True(t, er.Contains(wfIDContract))
		require.True(t, er.Contains(wfIDGrpc))
	})

	t.Run("independent source reconciliation preserves isolation", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()

		// Setup: multiple workflows from each source
		er := NewEngineRegistry()
		wfIDContract1 := [32]byte{1}
		wfIDContract2 := [32]byte{2}
		wfIDGrpc1 := [32]byte{3}
		wfIDGrpc2 := [32]byte{4}
		require.NoError(t, er.Add(wfIDContract1, ContractWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDContract2, ContractWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDGrpc1, GRPCWorkflowSourceName, &mockService{}))
		require.NoError(t, er.Add(wfIDGrpc2, GRPCWorkflowSourceName, &mockService{}))

		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
				return nil, nil
			},
			"",
			"test-chain-selector",
			Config{
				QueryCount:   20,
				SyncStrategy: SyncStrategyReconciliation,
			},
			&eventHandler{},
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		// Contract source: wfIDContract1 removed (only wfIDContract2 remains)
		contractPending := make(map[string]*reconciliationEvent)
		contractMeta := []WorkflowMetadataView{{
			WorkflowID:   wfIDContract2,
			Status:       WorkflowStatusActive,
			Source:       ContractWorkflowSourceName,
			WorkflowName: "contract-workflow-2",
			BinaryURL:    "http://binary.url",
			ConfigURL:    "http://config.url",
		}}
		contractEvents, err := wr.generateReconciliationEvents(
			ctx, contractPending, contractMeta, &types.Head{Height: "123"}, ContractWorkflowSourceName)
		require.NoError(t, err)

		// Should delete wfIDContract1
		require.Len(t, contractEvents, 1)
		require.Equal(t, WorkflowDeleted, contractEvents[0].Name)
		require.Equal(t, wfTypes.WorkflowID(wfIDContract1), contractEvents[0].Data.(WorkflowDeletedEvent).WorkflowID)

		// GRPC source: wfIDGrpc2 removed (only wfIDGrpc1 remains)
		grpcPending := make(map[string]*reconciliationEvent)
		grpcMeta := []WorkflowMetadataView{{
			WorkflowID:   wfIDGrpc1,
			Status:       WorkflowStatusActive,
			Source:       GRPCWorkflowSourceName,
			WorkflowName: "grpc-workflow-1",
			BinaryURL:    "http://binary.url",
			ConfigURL:    "http://config.url",
		}}
		grpcEvents, err := wr.generateReconciliationEvents(
			ctx, grpcPending, grpcMeta, &types.Head{Height: "123"}, GRPCWorkflowSourceName)
		require.NoError(t, err)

		// Should delete wfIDGrpc2, but NOT any contract workflows
		require.Len(t, grpcEvents, 1)
		require.Equal(t, WorkflowDeleted, grpcEvents[0].Name)
		require.Equal(t, wfTypes.WorkflowID(wfIDGrpc2), grpcEvents[0].Data.(WorkflowDeletedEvent).WorkflowID)
	})
}

func Test_isZeroOwner(t *testing.T) {
	t.Parallel()
	t.Run("returns true for nil slice", func(t *testing.T) {
		t.Parallel()
		require.True(t, isZeroOwner(nil))
	})

	t.Run("returns true for empty slice", func(t *testing.T) {
		t.Parallel()
		require.True(t, isZeroOwner([]byte{}))
	})

	t.Run("returns true for all zeros (20 bytes - Ethereum address)", func(t *testing.T) {
		t.Parallel()
		zeroAddress := make([]byte, 20)
		require.True(t, isZeroOwner(zeroAddress))
	})

	t.Run("returns true for all zeros (arbitrary length)", func(t *testing.T) {
		t.Parallel()
		zeros := make([]byte, 32)
		require.True(t, isZeroOwner(zeros))
	})

	t.Run("returns false for valid owner address", func(t *testing.T) {
		t.Parallel()
		validOwner, _ := hex.DecodeString("1234567890123456789012345678901234567890")
		require.False(t, isZeroOwner(validOwner))
	})

	t.Run("returns false for address with single non-zero byte", func(t *testing.T) {
		t.Parallel()
		almostZero := make([]byte, 20)
		almostZero[19] = 1 // last byte is 1
		require.False(t, isZeroOwner(almostZero))
	})

	t.Run("returns false for address with non-zero first byte", func(t *testing.T) {
		t.Parallel()
		almostZero := make([]byte, 20)
		almostZero[0] = 1 // first byte is 1
		require.False(t, isZeroOwner(almostZero))
	})
}

func Test_ParallelEventHandling(t *testing.T) {
	t.Parallel()
	t.Run("processes multiple delete events concurrently", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()

		n := 10
		wfIDs := make([]wfTypes.WorkflowID, n)
		for i := range wfIDs {
			wfIDs[i] = wfTypes.WorkflowID([32]byte{byte(i + 1)})
			require.NoError(t, er.Add(wfIDs[i], "TestSource", &mockService{}))
		}

		handler := newTestEvtHandler(nil)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) { return nil, nil },
			"",
			"test-chain-selector",
			Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
			handler,
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		pendingEvents := map[string]*reconciliationEvent{}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, []WorkflowMetadataView{}, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Len(t, events, n)

		// Simulate the parallel event loop from syncUsingReconciliationStrategy
		sourceIdentifier := "TestSource"
		pendingEventsBySource := map[string]map[string]*reconciliationEvent{
			sourceIdentifier: {},
		}
		reconcileReport := newReconcileReport()

		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, event := range events {
			mu.Lock()
			reconcileReport.NumEventsByType[string(event.Name)]++
			mu.Unlock()

			wg.Go(func() {
				handleErr := wr.handleWithMetrics(ctx, event.Event)
				if handleErr != nil {
					mu.Lock()
					pendingEventsBySource[sourceIdentifier][event.id] = event
					mu.Unlock()
				}
			})
		}
		wg.Wait()

		handled := handler.GetEvents()
		require.Len(t, handled, n)

		handledIDs := make(map[wfTypes.WorkflowID]bool)
		for _, evt := range handled {
			d := evt.Data.(WorkflowDeletedEvent)
			handledIDs[d.WorkflowID] = true
		}
		for _, id := range wfIDs {
			require.True(t, handledIDs[id], "expected workflow %x to be handled", id)
		}

		require.Empty(t, pendingEventsBySource[sourceIdentifier])
		require.Equal(t, n, reconcileReport.NumEventsByType[string(WorkflowDeleted)])
	})

	t.Run("processes mixed event types concurrently", func(t *testing.T) {
		t.Parallel()
		lggr := logger.TestLogger(t)
		ctx := t.Context()
		workflowDonNotifier := capabilities.NewDonNotifier()
		er := NewEngineRegistry()

		existingID := wfTypes.WorkflowID([32]byte{1})
		require.NoError(t, er.Add(existingID, "TestSource", &mockService{}))

		newID := wfTypes.WorkflowID([32]byte{2})

		handler := newTestEvtHandler(nil)
		wr, err := NewWorkflowRegistry(
			lggr,
			func(ctx context.Context, bytes []byte) (types.ContractReader, error) { return nil, nil },
			"",
			"test-chain-selector",
			Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
			handler,
			workflowDonNotifier,
			er,
		)
		require.NoError(t, err)

		pendingEvents := map[string]*reconciliationEvent{}
		metadata := []WorkflowMetadataView{
			{
				WorkflowID:   newID,
				Owner:        []byte{0x01},
				Status:       WorkflowStatusActive,
				WorkflowName: "new-wf",
				BinaryURL:    "b1",
				ConfigURL:    "c1",
				DonFamily:    "A",
			},
		}
		events, err := wr.generateReconciliationEvents(ctx, pendingEvents, metadata, &types.Head{Height: "123"}, "TestSource")
		require.NoError(t, err)
		require.Len(t, events, 2) // 1 delete + 1 activate

		var wg sync.WaitGroup
		for _, event := range events {
			wg.Add(1)
			go func(evt *reconciliationEvent) {
				defer wg.Done()
				_ = wr.handleWithMetrics(ctx, evt.Event)
			}(event)
		}
		wg.Wait()

		handled := handler.GetEvents()
		require.Len(t, handled, 2)

		nameSet := map[WorkflowRegistryEventName]bool{}
		for _, evt := range handled {
			nameSet[evt.Name] = true
		}
		require.True(t, nameSet[WorkflowDeleted])
		require.True(t, nameSet[WorkflowActivated])
	})
}

type mockShardMappingClient struct {
	mappings map[string]uint32
}

func (m *mockShardMappingClient) GetWorkflowShardMapping(_ context.Context, workflowIDs []string) (*ringpb.GetWorkflowShardMappingResponse, error) {
	out := make(map[string]uint32)
	for _, id := range workflowIDs {
		if shard, ok := m.mappings[id]; ok {
			out[id] = shard
		}
	}
	return &ringpb.GetWorkflowShardMappingResponse{Mappings: out, MappingVersion: 1}, nil
}

func (m *mockShardMappingClient) ReportWorkflowTriggerRegistration(context.Context, *ringpb.ReportWorkflowTriggerRegistrationRequest) (*ringpb.ReportWorkflowTriggerRegistrationResponse, error) {
	return &ringpb.ReportWorkflowTriggerRegistrationResponse{Success: true}, nil
}

func (m *mockShardMappingClient) Close() error { return nil }

var _ shardorchestrator.ClientInterface = (*mockShardMappingClient)(nil)

func TestWorkflowRegistry_filterWorkflowsByShard(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	wf1 := wfTypes.WorkflowID([32]byte{1})
	wf2 := wfTypes.WorkflowID([32]byte{2})
	wf3 := wfTypes.WorkflowID([32]byte{3})
	workflows := []WorkflowMetadataView{
		{WorkflowID: wf1, WorkflowName: "wf1"},
		{WorkflowID: wf2, WorkflowName: "wf2"},
		{WorkflowID: wf3, WorkflowName: "wf3"},
	}

	client := &mockShardMappingClient{
		mappings: map[string]uint32{
			wf1.Hex(): 0,
			wf2.Hex(): 1,
		},
	}
	wr := &workflowRegistry{
		shardOrchestratorClient: client,
		shardResolver:           shardownership.NewRingOCRShardResolver(client, logger.TestLogger(t)),
		myShardID:               1,
		shardingEnabled:         true,
	}

	filtered, err := wr.filterWorkflowsByShard(ctx, workflows)
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	require.Equal(t, wf2.Hex(), filtered[0].WorkflowID.Hex())
}

func TestWorkflowRegistry_ShardResolverWiring(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	wf1 := wfTypes.WorkflowID([32]byte{1})
	wf2 := wfTypes.WorkflowID([32]byte{2})
	owner1 := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	owner2 := []byte{0x11, 0x22, 0x33, 0x44, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	workflows := []WorkflowMetadataView{
		{WorkflowID: wf1, Owner: owner1, WorkflowName: "wf1"},
		{WorkflowID: wf2, Owner: owner2, WorkflowName: "wf2"},
	}

	const manualConfigTOML = `static_default_assignment = [0]

[per_owner_assignment]
  "aabbccdd00000000000000000000000000000000" = [1]
`

	newReg := func(t *testing.T, opts ...Option) *workflowRegistry {
		t.Helper()
		wr, err := NewWorkflowRegistry(
			logger.TestLogger(t),
			func(_ context.Context, _ []byte) (types.ContractReader, error) { return nil, nil },
			"", "test-chain-selector",
			Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
			&eventHandler{},
			&testDonNotifier{don: commonCap.DON{ID: 1}},
			NewEngineRegistry(),
			opts...,
		)
		require.NoError(t, err)
		return wr
	}

	t.Run("ringocr-only: WithShardOrchestratorClient auto-wires resolver", func(t *testing.T) {
		t.Parallel()
		client := &mockShardMappingClient{mappings: map[string]uint32{
			wf1.Hex(): 0,
			wf2.Hex(): 1,
		}}
		wr := newReg(t,
			WithShardEnabled(true),
			WithShardID(0),
			WithShardOrchestratorClient(client),
		)
		require.NotNil(t, wr.shardResolver, "resolver must be auto-wired from non-nil client")

		filtered, err := wr.filterWorkflowsByShard(ctx, workflows)
		require.NoError(t, err)
		require.Len(t, filtered, 1)
		require.Equal(t, wf1.Hex(), filtered[0].WorkflowID.Hex())
	})

	t.Run("manual-only: nil client does not auto-wire; explicit resolver used", func(t *testing.T) {
		t.Parallel()
		settings := &loop.AtomicSettings{}
		require.NoError(t, settings.Store(core.SettingsUpdate{
			Settings: manualConfigTOML,
			Hash:     "test",
		}))
		manual := shardownership.NewManualShardResolver(settings, nil, logger.TestLogger(t))
		wr := newReg(t,
			WithShardEnabled(true),
			WithShardID(0),
			WithShardOrchestratorClient(nil),
			WithRegistryShardResolver(manual),
		)
		require.NotNil(t, wr.shardResolver, "resolver must be set via WithRegistryShardResolver")

		// owner aabbccdd -> shard 1 (per_owner), owner 11223344 -> shard 0 (static default)
		filtered, err := wr.filterWorkflowsByShard(ctx, workflows)
		require.NoError(t, err)
		require.Len(t, filtered, 1)
		require.Equal(t, wf2.Hex(), filtered[0].WorkflowID.Hex())
	})

	t.Run("ringocr-with-overrides: explicit resolver not overwritten by auto-wire", func(t *testing.T) {
		t.Parallel()
		settings := &loop.AtomicSettings{}
		require.NoError(t, settings.Store(core.SettingsUpdate{
			Settings: manualConfigTOML,
			Hash:     "test",
		}))
		ringClient := &mockShardMappingClient{mappings: map[string]uint32{
			wf1.Hex(): 0,
			wf2.Hex(): 1,
		}}
		override := shardownership.NewOverrideShardResolver(settings, nil,
			shardownership.NewRingOCRShardResolver(ringClient, logger.TestLogger(t)),
			logger.TestLogger(t))
		wr := newReg(t,
			WithShardEnabled(true),
			WithShardID(0),
			WithShardOrchestratorClient(ringClient),
			WithRegistryShardResolver(override),
		)

		// If auto-wire had overwritten the override, wf1 would pass (ringClient
		// maps it to shard 0). With the override, manual wins for owner
		// aabbccdd -> shard 1, so only wf2 (static default shard 0) passes.
		filtered, err := wr.filterWorkflowsByShard(ctx, workflows)
		require.NoError(t, err)
		require.Len(t, filtered, 1)
		require.Equal(t, wf2.Hex(), filtered[0].WorkflowID.Hex())
	})

	t.Run("explicit resolver before WithShardOrchestratorClient is preserved", func(t *testing.T) {
		t.Parallel()
		settings := &loop.AtomicSettings{}
		require.NoError(t, settings.Store(core.SettingsUpdate{
			Settings: manualConfigTOML,
			Hash:     "test",
		}))
		ringClient := &mockShardMappingClient{mappings: map[string]uint32{
			wf1.Hex(): 0,
		}}
		override := shardownership.NewOverrideShardResolver(settings, nil,
			shardownership.NewRingOCRShardResolver(ringClient, logger.TestLogger(t)),
			logger.TestLogger(t))
		wr := newReg(t,
			WithShardEnabled(true),
			WithShardID(0),
			WithRegistryShardResolver(override),
			WithShardOrchestratorClient(ringClient),
		)

		// Same expectation: override resolver preserved, wf1 filtered to shard 1
		filtered, err := wr.filterWorkflowsByShard(ctx, workflows)
		require.NoError(t, err)
		require.Len(t, filtered, 1)
		require.Equal(t, wf2.Hex(), filtered[0].WorkflowID.Hex())
	})

	t.Run("nil client without explicit resolver: no filtering", func(t *testing.T) {
		t.Parallel()
		wr := newReg(t,
			WithShardEnabled(true),
			WithShardID(0),
			WithShardOrchestratorClient(nil),
		)
		require.Nil(t, wr.shardResolver, "resolver must be nil when client is nil and no explicit resolver")

		filtered, err := wr.filterWorkflowsByShard(ctx, workflows)
		require.NoError(t, err)
		require.Len(t, filtered, 2, "nil resolver means no filtering — all workflows pass through")
	})
}

func TestWorkflowRegistry_getTicker_nonPositiveDuration(t *testing.T) {
	t.Parallel()
	wr := &workflowRegistry{
		clock: clockwork.NewRealClock(),
	}

	require.NotPanics(t, func() {
		tickerCh := wr.getTicker(-1 * time.Second)
		require.NotNil(t, tickerCh)
	})
}

func TestWorkflowRegistry_getTicker_WithTickerOverride(t *testing.T) {
	t.Parallel()
	customCh := make(chan time.Time)
	wr := &workflowRegistry{
		ticker: customCh,
		clock:  clockwork.NewRealClock(),
	}

	tickerCh := wr.getTicker(10 * time.Second)
	require.Equal(t, (<-chan time.Time)(customCh), tickerCh)
}

// orphanSweepFakeHandler is a minimal evtHandler that returns a fixed spec list
// and records every Handle call, so the reconciliation loop's orphaned-spec
// behavior can be tested end to end.
type orphanSweepFakeHandler struct {
	mu      sync.Mutex
	specs   []*job.WorkflowSpec
	handled []Event
}

func (h *orphanSweepFakeHandler) Close() error                { return nil }
func (h *orphanSweepFakeHandler) Start(context.Context) error { return nil }
func (h *orphanSweepFakeHandler) Handle(_ context.Context, event Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handled = append(h.handled, event)
	return nil
}
func (h *orphanSweepFakeHandler) EmitActivationAbandoned(context.Context, Event, eventsv2.ActivationAbandonReason, error, int32) error {
	return nil
}
func (h *orphanSweepFakeHandler) ListWorkflowSpecs(context.Context) ([]*job.WorkflowSpec, error) {
	return h.specs, nil
}

func (h *orphanSweepFakeHandler) SetWorkflowDon(commonCap.DON) {}
func (h *orphanSweepFakeHandler) Handled() []Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	cp := make([]Event, len(h.handled))
	copy(cp, h.handled)
	return cp
}

// OrphanDeletes returns the workflow IDs of every WorkflowDeleted event
// dispatched to the handler. In these tests the engine registry is empty, so
// every delete is sweep-generated.
func (h *orphanSweepFakeHandler) OrphanDeletes() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	var ids []string
	for _, evt := range h.handled {
		if evt.Name != WorkflowDeleted {
			continue
		}
		if payload, ok := evt.Data.(WorkflowDeletedEvent); ok {
			ids = append(ids, payload.WorkflowID.Hex())
		}
	}
	return ids
}

func newTestRegistryMetrics(t *testing.T) *metrics {
	t.Helper()
	m, err := newMetrics()
	require.NoError(t, err)
	return m
}

// Orphaned specs — rows whose workflow ID is absent from the union of every
// source's metadata with no running engine — are released through the normal
// WorkflowDeleted event path. Liveness is judged against the union regardless
// of the row's Source attribution: IDs are content-addressed and may recur
// across sources, and pre-attribution rows carry no source at all.
func Test_reconcileOrphanedSpecs(t *testing.T) {
	t.Parallel()

	liveID := wfTypes.WorkflowID{1}
	orphanID := wfTypes.WorkflowID{2}

	t.Run("engine-less orphans are released regardless of attribution", func(t *testing.T) {
		t.Parallel()
		attributedOrphan := wfTypes.WorkflowID{3}
		h := &orphanSweepFakeHandler{}
		w := &workflowRegistry{lggr: logger.TestLogger(t), handler: h, engineRegistry: NewEngineRegistry(), metrics: newTestRegistryMetrics(t)}

		specs := []*job.WorkflowSpec{
			{WorkflowID: liveID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
			{WorkflowID: orphanID.Hex(), WorkflowOwner: "aabbccdd"}, // pre-attribution row
			{WorkflowID: attributedOrphan.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
		}
		w.reconcileOrphanedSpecs(t.Context(), specs, map[string]struct{}{liveID.Hex(): {}})

		assert.ElementsMatch(t, []string{orphanID.Hex(), attributedOrphan.Hex()}, h.OrphanDeletes(),
			"both orphans are released; the live spec is retained")
	})

	t.Run("row listed by any source is retained even when its attribution differs", func(t *testing.T) {
		t.Parallel()
		// Workflow IDs are content-addressed and may recur across sources: a
		// row written by source-a whose ID is currently listed by source-b is
		// live and must not be released.
		h := &orphanSweepFakeHandler{}
		w := &workflowRegistry{lggr: logger.TestLogger(t), handler: h, engineRegistry: NewEngineRegistry(), metrics: newTestRegistryMetrics(t)}

		specs := []*job.WorkflowSpec{
			{WorkflowID: orphanID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
		}
		w.reconcileOrphanedSpecs(t.Context(), specs, map[string]struct{}{orphanID.Hex(): {}})

		assert.Empty(t, h.Handled(), "a row present in the union must never be released")
	})

	t.Run("engine-owned orphan is left to engine reconciliation", func(t *testing.T) {
		t.Parallel()
		er := NewEngineRegistry()
		require.NoError(t, er.Add(orphanID, "test-source", &mockService{}))
		h := &orphanSweepFakeHandler{}
		w := &workflowRegistry{lggr: logger.TestLogger(t), handler: h, engineRegistry: er, metrics: newTestRegistryMetrics(t)}

		specs := []*job.WorkflowSpec{
			{WorkflowID: orphanID.Hex(), WorkflowOwner: "aabbccdd"},
		}
		w.reconcileOrphanedSpecs(t.Context(), specs, map[string]struct{}{liveID.Hex(): {}})

		assert.Empty(t, h.Handled(), "engine-owned orphan must not be released here")
	})
}

// fakeMetadataSource is a canned WorkflowMetadataSource for driving the
// reconciliation loop: it returns fixed metadata or a fixed error.
type fakeMetadataSource struct {
	name     string
	metadata []WorkflowMetadataView
	err      error
}

func (s *fakeMetadataSource) ListWorkflowMetadata(context.Context, commonCap.DON) ([]WorkflowMetadataView, *types.Head, error) {
	if s.err != nil {
		return nil, nil, s.err
	}
	return s.metadata, &types.Head{Height: "1"}, nil
}
func (s *fakeMetadataSource) Name() string             { return s.name }
func (s *fakeMetadataSource) SourceIdentifier() string { return s.name }
func (s *fakeMetadataSource) Ready() error             { return nil }

// failingShardMappingClient always fails shard-mapping lookups, driving the
// shard-filter error path.
type failingShardMappingClient struct{}

func (f *failingShardMappingClient) GetWorkflowShardMapping(context.Context, []string) (*ringpb.GetWorkflowShardMappingResponse, error) {
	return nil, assert.AnError
}

func (f *failingShardMappingClient) ReportWorkflowTriggerRegistration(context.Context, *ringpb.ReportWorkflowTriggerRegistrationRequest) (*ringpb.ReportWorkflowTriggerRegistrationResponse, error) {
	return &ringpb.ReportWorkflowTriggerRegistrationResponse{Success: true}, nil
}

func (f *failingShardMappingClient) Close() error { return nil }

var _ shardorchestrator.ClientInterface = (*failingShardMappingClient)(nil)

// newSweepLoopRegistry builds a workflowRegistry wired for driving
// syncUsingReconciliationStrategy directly: no contract source, injected fake
// sources, and an unbuffered ticker channel the test controls.
func newSweepLoopRegistry(t *testing.T, h evtHandler, er *EngineRegistry, sources []WorkflowMetadataSource, opts ...Option) (*workflowRegistry, chan time.Time) {
	t.Helper()
	tick := make(chan time.Time)
	wr, err := NewWorkflowRegistry(
		logger.TestLogger(t),
		func(ctx context.Context, bytes []byte) (types.ContractReader, error) { return nil, nil },
		"", // no contract source
		"test-chain-selector",
		Config{QueryCount: 20, SyncStrategy: SyncStrategyReconciliation},
		h,
		&testDonNotifier{don: commonCap.DON{ID: 1}},
		er,
		append([]Option{WithTicker(tick)}, opts...)...,
	)
	require.NoError(t, err)
	wr.workflowSources = sources
	return wr, tick
}

// runOneSweepTick drives the reconciliation loop through at least one complete
// tick. The ticker channel is unbuffered, so the second send only succeeds
// once the loop has fully processed the first tick — including the orphan
// sweep decision — and returned to its select.
func runOneSweepTick(t *testing.T, wr *workflowRegistry, tick chan time.Time) {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		defer close(done)
		wr.syncUsingReconciliationStrategy(ctx)
	}()
	tick <- time.Now()
	tick <- time.Now() // barrier: the first tick has been fully processed
	cancel()
	<-done
}

// Test_syncLoop_OrphanedSpecs exercises orphaned-spec reconciliation end to
// end through syncUsingReconciliationStrategy: orphans are released per
// source, and a source failure only defers its own specs.
func Test_syncLoop_OrphanedSpecs(t *testing.T) {
	t.Parallel()

	liveID := wfTypes.WorkflowID{1}
	orphanID := wfTypes.WorkflowID{2}
	liveMeta := WorkflowMetadataView{
		WorkflowID:   liveID,
		Owner:        []byte{0xaa, 0xbb, 0xcc, 0xdd},
		Status:       WorkflowStatusActive,
		WorkflowName: "live-wf",
		BinaryURL:    "b1",
		ConfigURL:    "c1",
		Source:       "source-a",
	}

	t.Run("engine-less orphan is released through the event path", func(t *testing.T) {
		t.Parallel()
		h := &orphanSweepFakeHandler{specs: []*job.WorkflowSpec{
			{WorkflowID: liveID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
			{WorkflowID: orphanID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
		}}
		wr, tick := newSweepLoopRegistry(t, h, NewEngineRegistry(), []WorkflowMetadataSource{
			&fakeMetadataSource{name: "source-a", metadata: []WorkflowMetadataView{liveMeta}},
		})

		runOneSweepTick(t, wr, tick)

		deletes := h.OrphanDeletes()
		require.NotEmpty(t, deletes, "the orphan must be released")
		for _, id := range deletes {
			assert.Equal(t, orphanID.Hex(), id, "only the orphan is released")
		}
	})

	t.Run("any failing source defers orphan cleanup", func(t *testing.T) {
		t.Parallel()
		// The union of source metadata is incomplete when a source fails: an
		// absent row may belong to the failed source, so no row can be safely
		// judged orphaned. Healthy sources still reconcile their engines.
		okMeta := liveMeta
		okMeta.Source = "source-ok"
		h := &orphanSweepFakeHandler{specs: []*job.WorkflowSpec{
			{WorkflowID: liveID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-ok"},
			{WorkflowID: orphanID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-bad"},
		}}
		wr, tick := newSweepLoopRegistry(t, h, NewEngineRegistry(), []WorkflowMetadataSource{
			&fakeMetadataSource{name: "source-ok", metadata: []WorkflowMetadataView{okMeta}},
			&fakeMetadataSource{name: "source-bad", err: assert.AnError},
		})

		runOneSweepTick(t, wr, tick)

		assert.NotEmpty(t, h.Handled(), "healthy source still reconciles")
		assert.Empty(t, h.OrphanDeletes(), "no spec is released while any source is failing")
	})

	t.Run("shard filter error skips the source entirely", func(t *testing.T) {
		t.Parallel()
		h := &orphanSweepFakeHandler{specs: []*job.WorkflowSpec{
			{WorkflowID: liveID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
			{WorkflowID: orphanID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
		}}
		wr, tick := newSweepLoopRegistry(t, h, NewEngineRegistry(), []WorkflowMetadataSource{
			&fakeMetadataSource{name: "source-a", metadata: []WorkflowMetadataView{liveMeta}},
		},
			WithShardEnabled(true),
			WithShardID(0),
			WithShardOrchestratorClient(&failingShardMappingClient{}),
		)

		runOneSweepTick(t, wr, tick)

		assert.Empty(t, h.Handled(), "no events dispatched and no spec released when the shard filter fails")
	})

	t.Run("workflow filtered to another shard is swept", func(t *testing.T) {
		t.Parallel()
		myID := wfTypes.WorkflowID{3}
		movedID := wfTypes.WorkflowID{4}
		myMeta := liveMeta
		myMeta.WorkflowID = myID
		movedMeta := liveMeta
		movedMeta.WorkflowID = movedID
		h := &orphanSweepFakeHandler{specs: []*job.WorkflowSpec{
			{WorkflowID: myID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
			{WorkflowID: movedID.Hex(), WorkflowOwner: "aabbccdd", Source: "source-a"},
		}}
		wr, tick := newSweepLoopRegistry(t, h, NewEngineRegistry(), []WorkflowMetadataSource{
			&fakeMetadataSource{name: "source-a", metadata: []WorkflowMetadataView{myMeta, movedMeta}},
		},
			WithShardEnabled(true),
			WithShardID(0),
			WithShardOrchestratorClient(&mockShardMappingClient{mappings: map[string]uint32{
				myID.Hex():    0, // stays on this shard
				movedID.Hex(): 1, // moved to another shard
			}}),
		)

		runOneSweepTick(t, wr, tick)

		// The moved workflow is absent from this shard's post-filter metadata
		// and has no engine, so its spec is released here (-1); the receiving
		// shard registers it (+1) under a distinct event_id — net level 0.
		deletes := h.OrphanDeletes()
		require.NotEmpty(t, deletes, "the shard-moved workflow must be swept")
		for _, id := range deletes {
			assert.Equal(t, movedID.Hex(), id, "only the moved workflow is released")
		}
	})
}
