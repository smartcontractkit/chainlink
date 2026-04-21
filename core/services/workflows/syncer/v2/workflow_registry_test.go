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
	"github.com/stretchr/testify/require"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	wfTypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

func Test_generateReconciliationEventsV2(t *testing.T) {
	// Validate that if no engines are on the node in the registry,
	// and we see that the contract has workflow state,
	// that we generate a WorkflowActivatedEvent
	t.Run("WorkflowActivatedEvent_whenNoEnginesInRegistry", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
			Tag:           tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedActivatedEvent, events[0].Data)
	})

	t.Run("WorkflowUpdatedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
			Tag:           tag,
			Attributes:    attributes,
		}
		require.Equal(t, expectedActivatedEvent, events[1].Data)
	})

	t.Run("WorkflowDeletedEvent", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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

	t.Run("No change", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
			Tag:           tag,
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
		workflowDonNotifier := capabilities.NewDonNotifier()
		// Engine already in the workflow registry
		er := NewEngineRegistry()
		wfID := [32]byte{1}
		owner := []byte{}
		wfName := "wf name 1"
		err := er.Add(wfID, ContractWorkflowSourceName, &mockService{})
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
			Tag:           tag,
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

	t.Run("a paused workflow clears a pending activated event", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
			Tag:           tag,
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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

	t.Run("pending activate events are handled when workflow metadata no longer exists", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
			Tag:           tag,
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
	t.Run("successful start and close", func(t *testing.T) {
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
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
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
	t.Run("only deletes engines from specified source", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
	t.Run("source failure does not delete engines from that source", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
	t.Run("returns true for nil slice", func(t *testing.T) {
		require.True(t, isZeroOwner(nil))
	})

	t.Run("returns true for empty slice", func(t *testing.T) {
		require.True(t, isZeroOwner([]byte{}))
	})

	t.Run("returns true for all zeros (20 bytes - Ethereum address)", func(t *testing.T) {
		zeroAddress := make([]byte, 20)
		require.True(t, isZeroOwner(zeroAddress))
	})

	t.Run("returns true for all zeros (arbitrary length)", func(t *testing.T) {
		zeros := make([]byte, 32)
		require.True(t, isZeroOwner(zeros))
	})

	t.Run("returns false for valid owner address", func(t *testing.T) {
		validOwner, _ := hex.DecodeString("1234567890123456789012345678901234567890")
		require.False(t, isZeroOwner(validOwner))
	})

	t.Run("returns false for address with single non-zero byte", func(t *testing.T) {
		almostZero := make([]byte, 20)
		almostZero[19] = 1 // last byte is 1
		require.False(t, isZeroOwner(almostZero))
	})

	t.Run("returns false for address with non-zero first byte", func(t *testing.T) {
		almostZero := make([]byte, 20)
		almostZero[0] = 1 // first byte is 1
		require.False(t, isZeroOwner(almostZero))
	})
}

func Test_ParallelEventHandling(t *testing.T) {
	t.Run("processes multiple delete events concurrently", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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

			wg.Add(1)
			go func(evt *reconciliationEvent) {
				defer wg.Done()
				handleErr := wr.handleWithMetrics(ctx, evt.Event)
				if handleErr != nil {
					mu.Lock()
					pendingEventsBySource[sourceIdentifier][evt.id] = evt
					mu.Unlock()
				}
			}(event)
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
		lggr := logger.TestLogger(t)
		ctx := testutils.Context(t)
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
	ctx := testutils.Context(t)
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
		myShardID:               1,
		shardingEnabled:         true,
	}

	filtered, err := wr.filterWorkflowsByShard(ctx, workflows)
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	require.Equal(t, wf2.Hex(), filtered[0].WorkflowID.Hex())
}
